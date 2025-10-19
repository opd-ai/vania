// Package engine provides the game runner that integrates Ebiten rendering
// with the procedural generation system, handling the game loop, player
// movement, and visual display.
package engine

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/vania/internal/audio"
	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/graphics"
	"github.com/opd-ai/vania/internal/input"
	"github.com/opd-ai/vania/internal/particle"
	"github.com/opd-ai/vania/internal/physics"
	"github.com/opd-ai/vania/internal/render"
	"github.com/opd-ai/vania/internal/save"
	"github.com/opd-ai/vania/internal/world"
)

// GameRunner wraps the Game with Ebiten rendering
type GameRunner struct {
	game              *Game
	renderer          *render.Renderer
	inputHandler      *input.InputHandler
	playerBody        *physics.Body
	combatSystem      *CombatSystem
	transitionHandler *RoomTransitionHandler
	enemyInstances    []*entity.EnemyInstance
	itemInstances     []*entity.ItemInstance
	particleSystem    *particle.ParticleSystem
	particlePresets   *particle.ParticlePresets
	doubleJumpUsed    bool
	dashCooldown      int
	playerFacingDir   float64
	paused            bool
	saveManager       *save.SaveManager
	checkpointManager *save.CheckpointManager
	startTime         time.Time
	visitedRooms      map[int]bool
	defeatedEnemies   map[int]bool
	collectedItems    map[int]bool
	unlockedDoors     map[string]bool
	lockedDoorMessage string
	lockedDoorTimer   int
	itemMessage       string
	itemMessageTimer  int
	musicContext      *audio.MusicContext
}

// NewGameRunner creates a new game runner
func NewGameRunner(game *Game) *GameRunner {
	// Initialize player at starting position
	playerX := float64(render.ScreenWidth / 2)
	playerY := float64(render.ScreenHeight / 2)

	if game.World != nil && game.World.StartRoom != nil {
		// Position player in the start room
		playerX = 100.0
		playerY = 500.0
	}

	// Create enemy instances for current room
	var enemyInstances []*entity.EnemyInstance
	if game.CurrentRoom != nil && len(game.Entities) > 0 {
		// Spawn enemies in room (3-5 enemies per combat room)
		enemyCount := 3
		for i := 0; i < enemyCount && i < len(game.Entities); i++ {
			enemy := game.Entities[i]
			// Position enemies across the room
			enemyX := 300.0 + float64(i*150)
			enemyY := 500.0
			enemyInstances = append(enemyInstances, entity.NewEnemyInstance(enemy, enemyX, enemyY))
		}
	}

	transitionHandler := NewRoomTransitionHandler(game)

	// Initialize save system
	saveManager, err := save.NewSaveManager("")
	if err != nil {
		// Fall back to no save system if there's an error
		saveManager = nil
	}

	var checkpointManager *save.CheckpointManager
	if saveManager != nil {
		checkpointManager = save.NewCheckpointManager(saveManager)
	}

	// Create item instances for current room
	itemInstances := createItemInstancesForRoom(game.CurrentRoom, game.Items)

	return &GameRunner{
		game:              game,
		renderer:          render.NewRenderer(),
		inputHandler:      input.NewInputHandler(),
		playerBody:        physics.NewBody(playerX, playerY, physics.PlayerWidth, physics.PlayerHeight),
		combatSystem:      NewCombatSystem(),
		transitionHandler: transitionHandler,
		enemyInstances:    enemyInstances,
		itemInstances:     itemInstances,
		particleSystem:    particle.NewParticleSystem(1000), // Max 1000 particles
		particlePresets:   &particle.ParticlePresets{},
		doubleJumpUsed:    false,
		dashCooldown:      0,
		playerFacingDir:   1.0,
		paused:            false,
		saveManager:       saveManager,
		checkpointManager: checkpointManager,
		startTime:         time.Now(),
		visitedRooms:      make(map[int]bool),
		defeatedEnemies:   make(map[int]bool),
		collectedItems:    make(map[int]bool),
		unlockedDoors:     make(map[string]bool),
		lockedDoorMessage: "",
		lockedDoorTimer:   0,
		itemMessage:       "",
		itemMessageTimer:  0,
		musicContext:      audio.NewMusicContext(),
	}
}

// Update implements ebiten.Game interface
func (gr *GameRunner) Update() error {
	// Check for quit
	if gr.inputHandler.IsQuitRequested() {
		return ebiten.Termination
	}

	// Get input state
	inputState := gr.inputHandler.Update()

	// Handle pause
	if inputState.PausePress {
		gr.paused = !gr.paused
	}

	if gr.paused {
		return nil
	}

	// Update transition handler
	if gr.transitionHandler.Update() {
		// Transition completed - spawn new enemies and items
		gr.enemyInstances = gr.transitionHandler.SpawnEnemiesForRoom(gr.game.CurrentRoom)
		gr.itemInstances = gr.transitionHandler.SpawnItemsForRoom(gr.game.CurrentRoom)
	}

	// Don't update game logic during transition
	if gr.transitionHandler.IsTransitioning() {
		return nil
	}

	// Check for door collision and transition
	door := gr.transitionHandler.CheckDoorCollision(
		gr.game.Player.X,
		gr.game.Player.Y,
		physics.PlayerWidth,
		physics.PlayerHeight,
		gr.unlockedDoors,
	)
	if door != nil {
		// Player touched a door - start transition
		gr.transitionHandler.StartTransition(door)
		return nil
	}

	// Check if player is near a locked door
	if gr.lockedDoorTimer > 0 {
		gr.lockedDoorTimer--
	}
	gr.checkLockedDoorInteraction()

	// Apply physics
	gr.playerBody.ApplyGravity()

	// Update combat system
	gr.combatSystem.Update()

	// Update particle system
	gr.particleSystem.Update()

	// Track previous ground state for landing particles
	wasOnGround := gr.playerBody.OnGround

	// Handle player movement
	if inputState.MoveLeft {
		gr.playerBody.MoveHorizontal(-1)
		gr.playerFacingDir = -1.0
	} else if inputState.MoveRight {
		gr.playerBody.MoveHorizontal(1)
		gr.playerFacingDir = 1.0
	} else {
		gr.playerBody.ApplyFriction()
	}

	// Handle attack
	if inputState.AttackPress {
		gr.combatSystem.PlayerAttack()
	}

	// Handle jump
	if inputState.JumpPress {
		hasDoubleJump := gr.game.Player.Abilities["double_jump"]
		jumped := gr.playerBody.Jump(hasDoubleJump, &gr.doubleJumpUsed)
		if jumped {
			// Create jump dust particles
			emitter := gr.particlePresets.CreateJumpDust(gr.game.Player.X+16, gr.game.Player.Y+32)
			emitter.Burst(8)
			gr.particleSystem.AddEmitter(emitter)
		}
	}

	// Handle dash
	if inputState.DashPress && gr.dashCooldown <= 0 {
		if gr.game.Player.Abilities["dash"] {
			direction := 0.0
			if inputState.MoveRight {
				direction = 1.0
			} else if inputState.MoveLeft {
				direction = -1.0
			}
			gr.playerBody.Dash(direction)
			gr.dashCooldown = 30 // 30 frames cooldown

			// Create dash trail particles
			emitter := gr.particlePresets.CreateDashTrail(gr.game.Player.X+16, gr.game.Player.Y+16)
			emitter.Start()
			gr.particleSystem.AddEmitter(emitter)

			// Stop dash trail after dash duration (let it fade naturally)
			// The emitter will be removed automatically as particles die
		}
	}

	if gr.dashCooldown > 0 {
		gr.dashCooldown--
	}

	// Apply knockback from damage
	knockbackX, knockbackY := gr.combatSystem.GetKnockback()
	if knockbackX != 0 || knockbackY != 0 {
		gr.playerBody.Velocity.X += knockbackX
		gr.playerBody.Velocity.Y += knockbackY
	}

	// Update position
	gr.playerBody.Update()

	// Resolve collisions with platforms
	if gr.game.CurrentRoom != nil {
		gr.playerBody.ResolveCollisionWithPlatforms(gr.game.CurrentRoom.Platforms)
	}

	// Check for landing (player just touched ground)
	if !wasOnGround && gr.playerBody.OnGround {
		// Create landing dust particles
		emitter := gr.particlePresets.CreateLandDust(gr.game.Player.X+16, gr.game.Player.Y+32)
		emitter.Burst(12)
		gr.particleSystem.AddEmitter(emitter)
	}

	// Update player position in game
	gr.game.Player.X = gr.playerBody.Position.X
	gr.game.Player.Y = gr.playerBody.Position.Y
	gr.game.Player.VelX = gr.playerBody.Velocity.X
	gr.game.Player.VelY = gr.playerBody.Velocity.Y

	// Update player animation based on state
	if gr.game.Player.AnimController != nil {
		// Update animation
		gr.game.Player.AnimController.Update()

		// Determine which animation to play
		currentAnim := gr.game.Player.AnimController.GetCurrentAnimation()

		// Attack animation has priority
		if gr.combatSystem.IsPlayerAttacking() {
			if currentAnim != "attack" {
				gr.game.Player.AnimController.Play("attack", true)
			}
		} else if !gr.playerBody.OnGround {
			// In air - jump animation
			if currentAnim != "jump" {
				gr.game.Player.AnimController.Play("jump", true)
			}
		} else if inputState.MoveLeft || inputState.MoveRight {
			// Moving - walk animation
			if currentAnim != "walk" {
				gr.game.Player.AnimController.Play("walk", false)
			}
		} else {
			// Standing still - idle animation
			if currentAnim != "idle" {
				gr.game.Player.AnimController.Play("idle", false)
			}
		}
	}

	// Update enemies
	for _, enemy := range gr.enemyInstances {
		if enemy.IsDead() {
			continue
		}

		// Update enemy AI
		enemy.Update(gr.game.Player.X, gr.game.Player.Y)

		// Apply gravity to ground-based enemies
		if enemy.Enemy.Behavior != entity.FlyingBehavior {
			if !enemy.OnGround {
				enemy.VelY += 0.5
				if enemy.VelY > 10.0 {
					enemy.VelY = 10.0
				}
			}
		}

		// Update enemy position
		enemy.X += enemy.VelX
		enemy.Y += enemy.VelY

		// Check enemy collisions with platforms
		if gr.game.CurrentRoom != nil {
			enemy.OnGround = false
			for _, platform := range gr.game.CurrentRoom.Platforms {
				px := float64(platform.X)
				py := float64(platform.Y)
				pw := float64(platform.Width)
				ph := float64(platform.Height)

				ex, ey, ew, eh := enemy.GetBounds()

				// Check collision with platform
				if ex < px+pw && ex+ew > px && ey < py+ph && ey+eh > py {
					// Resolve collision - simple top collision for now
					if enemy.VelY > 0 && ey+eh-ph < py {
						enemy.Y = py - eh
						enemy.VelY = 0
						enemy.OnGround = true
					}
				}
			}
		}

		// Check player attack hitting enemy
		if gr.combatSystem.IsPlayerAttacking() {
			attackX, attackY, attackW, attackH := gr.combatSystem.GetAttackHitbox(
				gr.game.Player.X, gr.game.Player.Y, gr.playerFacingDir,
			)
			wasAlive := !enemy.IsDead()
			if gr.combatSystem.CheckEnemyHit(attackX, attackY, attackW, attackH, enemy) {
				// Create hit effect particles
				ex, ey, _, _ := enemy.GetBounds()
				hitEmitter := gr.particlePresets.CreateHitEffect(ex+16, ey+16, gr.playerFacingDir)
				hitEmitter.Burst(10)
				gr.particleSystem.AddEmitter(hitEmitter)

				// Create blood splatter particles
				bloodEmitter := gr.particlePresets.CreateBloodSplatter(ex+16, ey+16, gr.playerFacingDir)
				bloodEmitter.Burst(6)
				gr.particleSystem.AddEmitter(bloodEmitter)

				gr.combatSystem.ApplyDamageToEnemy(enemy, gr.game.Player.Damage, gr.game.Player.X)

				// Track defeated enemy (use position as ID for now)
				if wasAlive && enemy.IsDead() {
					enemyKey := int(enemy.X*1000 + enemy.Y)
					gr.defeatedEnemies[enemyKey] = true

					// Create explosion effect on death
					explosionEmitter := gr.particlePresets.CreateExplosion(ex+16, ey+16, 1.0)
					explosionEmitter.Burst(20)
					gr.particleSystem.AddEmitter(explosionEmitter)
				}
			}
		}

		// Check enemy collision with player
		if gr.combatSystem.CheckPlayerEnemyCollision(
			gr.game.Player.X, gr.game.Player.Y, physics.PlayerWidth, physics.PlayerHeight, enemy,
		) {
			damage := enemy.Enemy.Damage
			if enemy.State == entity.AttackState {
				damage = enemy.GetAttackDamage()
			}
			gr.combatSystem.ApplyDamageToPlayer(gr.game.Player, damage, enemy.X)
		}
	}

	// Update music context based on game state
	gr.updateMusicContext()

	// Check for item collection
	if gr.itemMessageTimer > 0 {
		gr.itemMessageTimer--
	}
	gr.checkItemCollection()

	// Update camera
	gr.renderer.UpdateCamera(gr.game.Player.X, gr.game.Player.Y)

	// Check for auto-save
	gr.CheckAutoSave()

	// Track current room as visited
	if gr.game.CurrentRoom != nil {
		gr.visitedRooms[gr.game.CurrentRoom.ID] = true
	}

	return nil
}

// Draw implements ebiten.Game interface
func (gr *GameRunner) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Render world
	if gr.game.CurrentRoom != nil && gr.game.Graphics != nil {
		gr.renderer.RenderWorld(screen, gr.game.CurrentRoom, gr.game.Graphics.Tilesets)
	}

	// Render items
	for _, item := range gr.itemInstances {
		if !item.Collected && !gr.collectedItems[item.ID] {
			itemX, itemY, itemW, itemH := item.GetBounds()
			gr.renderer.RenderItem(screen, itemX, itemY, itemW, itemH, item.Collected, nil)
		}
	}

	// Render enemies
	for _, enemy := range gr.enemyInstances {
		if !enemy.IsDead() {
			ex, ey, ew, eh := enemy.GetBounds()
			
			// Get current animation frame if available
			var spriteToRender *graphics.Sprite
			if enemy.AnimController != nil {
				spriteToRender = enemy.AnimController.GetCurrentFrame()
			}
			// Fallback to base sprite if no animation frame
			if spriteToRender == nil {
				if sprite, ok := enemy.Enemy.SpriteData.(*graphics.Sprite); ok {
					spriteToRender = sprite
				}
			}
			
			gr.renderer.RenderEnemy(screen, ex, ey, ew, eh, enemy.CurrentHealth, enemy.Enemy.Health, false, spriteToRender)
		}
	}

	// Render attack effect
	if gr.combatSystem.IsPlayerAttacking() {
		attackX, attackY, attackW, attackH := gr.combatSystem.GetAttackHitbox(
			gr.game.Player.X, gr.game.Player.Y, gr.playerFacingDir,
		)
		if attackW > 0 && attackH > 0 {
			gr.renderer.RenderAttackEffect(screen, attackX, attackY, attackW, attackH)
		}
	}

	// Render particles (before player so they appear behind)
	allParticles := gr.particleSystem.GetAllParticles()
	gr.renderer.RenderParticles(screen, allParticles)

	// Render player
	if gr.game.Player != nil {
		// Use animated sprite if available, otherwise fall back to base sprite
		spriteToRender := gr.game.Player.Sprite
		if gr.game.Player.AnimController != nil {
			animFrame := gr.game.Player.AnimController.GetCurrentFrame()
			if animFrame != nil {
				spriteToRender = animFrame
			}
		}
		gr.renderer.RenderPlayer(screen, gr.game.Player.X, gr.game.Player.Y, spriteToRender)
	}

	// Render UI
	if gr.game.Player != nil {
		gr.renderer.RenderUI(screen, gr.game.Player.Health, gr.game.Player.MaxHealth, gr.game.Player.Abilities)
	}

	// Render transition effect if transitioning
	if gr.transitionHandler.IsTransitioning() {
		progress := gr.transitionHandler.GetTransitionProgress()
		gr.renderer.RenderTransitionEffect(screen, progress)
	}

	// Show locked door message if active
	if gr.lockedDoorTimer > 0 && gr.lockedDoorMessage != "" {
		// Draw message in center of screen with background
		messageX := render.ScreenWidth/2 - 100
		messageY := render.ScreenHeight/2 - 20

		// Draw semi-transparent background
		messageImg := ebiten.NewImage(200, 40)
		messageImg.Fill(color.RGBA{0, 0, 0, 180})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(messageX), float64(messageY))
		screen.DrawImage(messageImg, op)

		// Draw text
		ebitenutil.DebugPrintAt(screen, gr.lockedDoorMessage, messageX+10, messageY+12)
	}

	// Show item collection message if active
	if gr.itemMessageTimer > 0 && gr.itemMessage != "" {
		// Draw message in center-top of screen with background
		messageX := render.ScreenWidth/2 - 100
		messageY := 80

		// Draw semi-transparent background
		messageImg := ebiten.NewImage(200, 40)
		messageImg.Fill(color.RGBA{255, 215, 0, 200}) // Golden background
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(messageX), float64(messageY))
		screen.DrawImage(messageImg, op)

		// Draw text
		ebitenutil.DebugPrintAt(screen, gr.itemMessage, messageX+10, messageY+12)
	}

	// Show debug info
	aliveEnemies := 0
	for _, enemy := range gr.enemyInstances {
		if !enemy.IsDead() {
			aliveEnemies++
		}
	}

	debugInfo := fmt.Sprintf("Seed: %d | Room: %s | FPS: %.2f | Enemies: %d/%d | Items: %d/%d\nPosition: (%.0f, %.0f) | Velocity: (%.1f, %.1f)\nHealth: %d/%d | OnGround: %v | Invuln: %v\nControls: WASD/Arrows=Move, Space=Jump, J=Attack, K=Dash, P=Pause, Ctrl+Q=Quit",
		gr.game.Seed,
		gr.getCurrentRoomName(),
		ebiten.ActualTPS(),
		aliveEnemies,
		len(gr.enemyInstances),
		len(gr.collectedItems),
		len(gr.game.Items),
		gr.game.Player.X,
		gr.game.Player.Y,
		gr.game.Player.VelX,
		gr.game.Player.VelY,
		gr.game.Player.Health,
		gr.game.Player.MaxHealth,
		gr.playerBody.OnGround,
		gr.combatSystem.IsInvulnerable(),
	)

	if gr.paused {
		debugInfo = "PAUSED\nPress P to resume\n\n" + debugInfo
	}

	ebitenutil.DebugPrint(screen, debugInfo)
}

// Layout implements ebiten.Game interface
func (gr *GameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return render.ScreenWidth, render.ScreenHeight
}

// Run starts the game with rendering
func (gr *GameRunner) Run() error {
	ebiten.SetWindowSize(render.ScreenWidth, render.ScreenHeight)
	ebiten.SetWindowTitle("VANIA - Procedural Metroidvania")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(gr); err != nil {
		return err
	}

	return nil
}

// getCurrentRoomName returns a friendly name for the current room
func (gr *GameRunner) getCurrentRoomName() string {
	if gr.game.CurrentRoom == nil {
		return "None"
	}

	if gr.game.CurrentRoom.Biome != nil {
		return gr.game.CurrentRoom.Biome.Name
	}

	return fmt.Sprintf("Room %d", gr.game.CurrentRoom.ID)
}

// CreateSaveData generates a SaveData struct from current game state
func (gr *GameRunner) CreateSaveData() *save.SaveData {
	// Build visited rooms list
	visitedRoomsList := make([]int, 0, len(gr.visitedRooms))
	for roomID := range gr.visitedRooms {
		visitedRoomsList = append(visitedRoomsList, roomID)
	}

	// Calculate play time
	playTime := int64(time.Since(gr.startTime).Seconds())

	// Current room ID
	currentRoomID := 0
	if gr.game.CurrentRoom != nil {
		currentRoomID = gr.game.CurrentRoom.ID
	}

	return &save.SaveData{
		Seed:            gr.game.Seed,
		PlayTime:        playTime,
		PlayerX:         gr.game.Player.X,
		PlayerY:         gr.game.Player.Y,
		PlayerHealth:    gr.game.Player.Health,
		PlayerMaxHealth: gr.game.Player.MaxHealth,
		PlayerAbilities: gr.game.Player.Abilities,
		CurrentRoomID:   currentRoomID,
		VisitedRooms:    visitedRoomsList,
		DefeatedEnemies: gr.defeatedEnemies,
		CollectedItems:  gr.collectedItems,
		UnlockedDoors:   gr.unlockedDoors,
		BossesDefeated:  gr.getBossesDefeated(),
		CheckpointID:    currentRoomID,
	}
}

// getBossesDefeated returns a list of defeated boss IDs
func (gr *GameRunner) getBossesDefeated() []int {
	bossesDefeated := make([]int, 0)
	// Check defeated enemies for bosses
	for enemyID := range gr.defeatedEnemies {
		// Boss enemies have IDs >= 1000 (convention in enemy generation)
		if enemyID >= 1000 {
			bossesDefeated = append(bossesDefeated, enemyID)
		}
	}
	return bossesDefeated
}

// checkLockedDoorInteraction checks if player is trying to use a locked door
func (gr *GameRunner) checkLockedDoorInteraction() {
	if gr.game.CurrentRoom == nil {
		return
	}

	playerX := gr.game.Player.X
	playerY := gr.game.Player.Y

	// Check collision with each door
	for i := range gr.game.CurrentRoom.Doors {
		door := &gr.game.CurrentRoom.Doors[i]

		// Simple AABB collision check
		doorX := float64(door.X)
		doorY := float64(door.Y)
		doorW := float64(door.Width)
		doorH := float64(door.Height)

		if playerX < doorX+doorW &&
			playerX+physics.PlayerWidth > doorX &&
			playerY < doorY+doorH &&
			playerY+physics.PlayerHeight > doorY {

			// Check if door is locked
			doorKey := gr.transitionHandler.GetDoorKey(door)
			if door.Locked && !gr.unlockedDoors[doorKey] {
				// Check if player can unlock this door
				if gr.transitionHandler.CanUnlockDoor(door, gr.game.Player.Abilities, gr.collectedItems) {
					// Automatically unlock the door
					gr.UnlockDoor(door)
				} else {
					// Show locked message
					if door.LeadsTo != nil {
						requirement := gr.transitionHandler.findEdgeRequirement(gr.game.CurrentRoom.ID, door.LeadsTo.ID)
						if requirement != "" {
							gr.lockedDoorMessage = fmt.Sprintf("Requires: %s", requirement)
						} else {
							gr.lockedDoorMessage = "Door is locked"
						}
					} else {
						gr.lockedDoorMessage = "Door is locked"
					}
					gr.lockedDoorTimer = 120 // Show for 2 seconds
				}
			}
		}
	}
}

// UnlockDoor unlocks a door and adds particle effect
func (gr *GameRunner) UnlockDoor(door *world.Door) {
	if door == nil {
		return
	}

	doorKey := gr.transitionHandler.GetDoorKey(door)
	gr.unlockedDoors[doorKey] = true

	// Show unlock message
	gr.lockedDoorMessage = "Door unlocked!"
	gr.lockedDoorTimer = 120 // Show for 2 seconds

	// Create sparkle particle effect at door position
	doorCenterX := float64(door.X) + float64(door.Width)/2
	doorCenterY := float64(door.Y) + float64(door.Height)/2
	sparkleEmitter := gr.particlePresets.CreateSparkles(doorCenterX, doorCenterY)
	sparkleEmitter.Burst(15)
	gr.particleSystem.AddEmitter(sparkleEmitter)
}

// checkItemCollection checks for item collision and collection
func (gr *GameRunner) checkItemCollection() {
	playerX := gr.game.Player.X
	playerY := gr.game.Player.Y
	playerW := float64(physics.PlayerWidth)
	playerH := float64(physics.PlayerHeight)

	for _, item := range gr.itemInstances {
		// Skip already collected items
		if item.Collected || gr.collectedItems[item.ID] {
			continue
		}

		// Check collision with player
		itemX, itemY, itemW, itemH := item.GetBounds()
		
		if playerX < itemX+itemW &&
			playerX+playerW > itemX &&
			playerY < itemY+itemH &&
			playerY+playerH > itemY {
			
			// Collect the item!
			gr.collectItem(item)
		}
	}
}

// collectItem handles item collection
func (gr *GameRunner) collectItem(item *entity.ItemInstance) {
	if item == nil || item.Collected {
		return
	}

	// Mark as collected
	item.Collected = true
	gr.collectedItems[item.ID] = true

	// Show message
	gr.itemMessage = fmt.Sprintf("Collected: %s", item.Item.Name)
	gr.itemMessageTimer = 120 // Show for 2 seconds

	// Create sparkle particle effect at item position
	sparkleEmitter := gr.particlePresets.CreateSparkles(item.X, item.Y)
	sparkleEmitter.Burst(20)
	gr.particleSystem.AddEmitter(sparkleEmitter)

	// Apply item effect if needed
	switch item.Item.Effect {
	case "heal":
		gr.game.Player.Health += item.Item.Value
		if gr.game.Player.Health > gr.game.Player.MaxHealth {
			gr.game.Player.Health = gr.game.Player.MaxHealth
		}
	case "increase_damage":
		gr.game.Player.Damage += item.Item.Value / 10
	}
}

// SaveGame saves the current game state to a slot
func (gr *GameRunner) SaveGame(slotID int) error {
	if gr.saveManager == nil {
		return fmt.Errorf("save system not initialized")
	}

	saveData := gr.CreateSaveData()
	return gr.saveManager.SaveGame(saveData, slotID)
}

// LoadGame loads game state from a slot
func (gr *GameRunner) LoadGame(slotID int) error {
	if gr.saveManager == nil {
		return fmt.Errorf("save system not initialized")
	}

	saveData, err := gr.saveManager.LoadGame(slotID)
	if err != nil {
		return err
	}

	return gr.RestoreFromSaveData(saveData)
}

// RestoreFromSaveData restores game state from save data
func (gr *GameRunner) RestoreFromSaveData(saveData *save.SaveData) error {
	// Verify seed matches
	if saveData.Seed != gr.game.Seed {
		return fmt.Errorf("save file seed mismatch: expected %d, got %d", gr.game.Seed, saveData.Seed)
	}

	// Restore player state
	gr.game.Player.X = saveData.PlayerX
	gr.game.Player.Y = saveData.PlayerY
	gr.game.Player.Health = saveData.PlayerHealth
	gr.game.Player.MaxHealth = saveData.PlayerMaxHealth
	gr.game.Player.Abilities = saveData.PlayerAbilities

	// Update player body position
	gr.playerBody.Position.X = saveData.PlayerX
	gr.playerBody.Position.Y = saveData.PlayerY

	// Restore world state
	gr.visitedRooms = make(map[int]bool)
	for _, roomID := range saveData.VisitedRooms {
		gr.visitedRooms[roomID] = true
	}
	gr.defeatedEnemies = saveData.DefeatedEnemies
	gr.collectedItems = saveData.CollectedItems
	if gr.collectedItems == nil {
		gr.collectedItems = make(map[int]bool)
	}
	gr.unlockedDoors = saveData.UnlockedDoors
	if gr.unlockedDoors == nil {
		gr.unlockedDoors = make(map[string]bool)
	}

	// Find and set current room
	for _, room := range gr.game.World.Rooms {
		if room.ID == saveData.CurrentRoomID {
			gr.game.CurrentRoom = room
			// No need to call SetCurrentRoom - the game already has CurrentRoom set
			break
		}
	}

	// Adjust start time to account for saved play time
	gr.startTime = time.Now().Add(-time.Duration(saveData.PlayTime) * time.Second)

	return nil
}

// CheckAutoSave checks if an auto-save should be triggered
func (gr *GameRunner) CheckAutoSave() {
	if gr.checkpointManager == nil {
		return
	}

	if gr.checkpointManager.ShouldCheckpoint() {
		saveData := gr.CreateSaveData()
		if err := gr.checkpointManager.CreateCheckpoint(saveData); err == nil {
			// Successfully auto-saved (could show notification to player)
		}
	}
}

// createItemInstancesForRoom creates item instances for a room
func createItemInstancesForRoom(room *world.Room, allItems []*entity.Item) []*entity.ItemInstance {
	var instances []*entity.ItemInstance
	
	if room == nil || room.Type != world.TreasureRoom {
		return instances
	}
	
	// Place 2-4 items in treasure rooms
	itemCount := 2 + (room.ID % 3) // 2-4 items based on room ID
	if itemCount > len(allItems) {
		itemCount = len(allItems)
	}
	
	for i := 0; i < itemCount && i < len(allItems); i++ {
		// Generate unique item ID based on room and position
		itemID := room.ID*1000 + i
		
		// Position items across the room (spread horizontally)
		itemX := 200.0 + float64(i*150)
		itemY := 500.0 // Ground level
		
		instance := entity.NewItemInstance(allItems[i%len(allItems)], itemID, itemX, itemY)
		instances = append(instances, instance)
	}
	
	return instances
}


// updateMusicContext updates the music context based on current game state
func (gr *GameRunner) updateMusicContext() {
// Count nearby enemies (alive enemies within aggro range)
nearbyCount := 0
inCombat := false

for _, enemy := range gr.enemyInstances {
if enemy.IsDead() {
continue
}

// Calculate distance to player
dx := gr.game.Player.X - enemy.X
dy := gr.game.Player.Y - enemy.Y
distance := dx*dx + dy*dy

// Check if enemy is nearby (within ~300 pixels)
if distance < 90000 { // 300^2
nearbyCount++
}

// Check if any enemy is actively chasing or attacking
if enemy.State == entity.ChaseState || enemy.State == entity.AttackState {
inCombat = true
}
}

	// Determine if this is a boss fight
	isBossFight := false
	if gr.game.CurrentRoom != nil && gr.game.CurrentRoom.Type == world.BossRoom {
		isBossFight = true
	}

	// Calculate player health percentage
	healthPct := 1.0 // Default to full health
	if gr.game.Player.MaxHealth > 0 {
		healthPct = float64(gr.game.Player.Health) / float64(gr.game.Player.MaxHealth)
	}

// Get room danger level
dangerLevel := 0
if gr.game.CurrentRoom != nil && gr.game.CurrentRoom.Biome != nil {
dangerLevel = gr.game.CurrentRoom.Biome.DangerLevel
}

// Update music context
gr.musicContext.InCombat = inCombat
gr.musicContext.IsBossFight = isBossFight
gr.musicContext.NearbyEnemyCount = nearbyCount
gr.musicContext.PlayerHealthPct = healthPct
gr.musicContext.RoomDangerLevel = dangerLevel

// Calculate intensity and update adaptive music track
intensity := gr.musicContext.CalculateIntensity()

// Get current biome's adaptive track and update it
if gr.game.CurrentRoom != nil && gr.game.CurrentRoom.Biome != nil {
if track, exists := gr.game.Audio.AdaptiveTracks[gr.game.CurrentRoom.Biome.Name]; exists {
track.SetIntensity(intensity)
track.Update()
}
}
}
