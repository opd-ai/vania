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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/opd-ai/vania/internal/audio"
	"github.com/opd-ai/vania/internal/engine/ecs"
	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/graphics"
	"github.com/opd-ai/vania/internal/input"
	"github.com/opd-ai/vania/internal/particle"
	"github.com/opd-ai/vania/internal/physics"
	"github.com/opd-ai/vania/internal/render"
	"github.com/opd-ai/vania/internal/save"
	"github.com/opd-ai/vania/internal/world"
)

const (
	// Message timing constants
	lockedDoorMessageDuration = 120 // 2 seconds at 60 FPS
	itemMessageDuration       = 120 // 2 seconds at 60 FPS
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
	grappleCooldown   int
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
	showDebugInfo     bool
	playerStatus      *StatusManager // active status effects on the player
	systemManager     *ecs.SystemManager
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
		// Find ground platform Y for spawning
		groundY := findGroundY(game.CurrentRoom)

		// Spawn enemies in room (3-5 enemies per combat room)
		enemyCount := 3
		for i := 0; i < enemyCount && i < len(game.Entities); i++ {
			enemy := game.Entities[i]
			// Position enemies across the room on the ground platform
			enemyX := 300.0 + float64(i*150)
			_, _, _, eh := entity.GetEnemySizeBounds(enemy)
			enemyY := groundY - eh
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

	// Create the renderer here so we can pass it to ECS systems
	renderer := render.NewRenderer()
	ps := particle.NewParticleSystem(1000) // Max 1000 particles

	// Initialize ECS SystemManager and register subsystem wrappers
	sm := ecs.NewSystemManager()
	sm.Register(NewAudioECSSystem(game.Audio), 10)
	sm.Register(NewParticleECSSystem(ps, renderer), 20)

	return &GameRunner{
		game:              game,
		renderer:          renderer,
		inputHandler:      input.NewInputHandler(),
		playerBody:        physics.NewBody(playerX, playerY, physics.PlayerWidth, physics.PlayerHeight),
		combatSystem:      NewCombatSystem(),
		transitionHandler: transitionHandler,
		enemyInstances:    enemyInstances,
		itemInstances:     itemInstances,
		particleSystem:    ps,
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
		showDebugInfo:     false, // Debug info starts hidden
		playerStatus:      NewStatusManager(),
		systemManager:     sm,
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

	// Handle debug toggle (F3 key)
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		gr.showDebugInfo = !gr.showDebugInfo
	}

	if gr.paused {
		return nil
	}

	return gr.updatePlaying(inputState)
}

// updatePlaying runs the main game-logic update when not paused.
func (gr *GameRunner) updatePlaying(inputState input.InputState) error {
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
		gr.transitionHandler.StartTransition(door)
		return nil
	}

	// Locked-door message countdown
	if gr.lockedDoorTimer > 0 {
		gr.lockedDoorTimer--
	}
	gr.checkLockedDoorInteraction()

	// Apply physics gravity (glide ability modifies fall rate)
	hasGlide := gr.game.Player.Abilities["glide"]
	isGliding := hasGlide && inputState.UseAbility && !gr.playerBody.OnGround && gr.playerBody.Velocity.Y > 0
	gr.playerBody.ApplyGravity(isGliding)

	gr.combatSystem.Update()
	gr.updateStatusEffects()

	// Delegate particle and audio updates through the ECS SystemManager
	if err := gr.systemManager.Update(1.0 / 60.0); err != nil {
		return err
	}

	wasOnGround := gr.playerBody.OnGround

	gr.updatePlayerInput(inputState)
	gr.updatePlayerPhysics(wasOnGround)
	gr.updatePlayerAnimation(inputState)
	gr.updateEnemies()

	gr.updateMusicContext()

	if gr.itemMessageTimer > 0 {
		gr.itemMessageTimer--
	}
	gr.checkItemCollection()
	gr.renderer.UpdateCamera(gr.game.Player.X, gr.game.Player.Y)
	gr.CheckAutoSave()
	gr.updateRoomTracking()

	return nil
}

// updateStatusEffects ticks active player status effects and applies damage.
func (gr *GameRunner) updateStatusEffects() {
	if statusDmg := gr.playerStatus.Update(1.0 / 60.0); statusDmg > 0 {
		gr.game.Player.Health -= statusDmg
		if gr.game.Player.Health < 0 {
			gr.game.Player.Health = 0
		}
		gr.combatSystem.AddDamageNumber(statusDmg, gr.game.Player.X, gr.game.Player.Y-10, false)
	}
}

// updatePlayerInput processes movement, attack, jump, dash, and grapple inputs.
func (gr *GameRunner) updatePlayerInput(inputState input.InputState) {
	speedMult := gr.playerStatus.SpeedMultiplier()
	if inputState.MoveLeft {
		gr.playerBody.MoveHorizontalScaled(-1, speedMult)
		gr.playerFacingDir = -1.0
	} else if inputState.MoveRight {
		gr.playerBody.MoveHorizontalScaled(1, speedMult)
		gr.playerFacingDir = 1.0
	} else {
		gr.playerBody.ApplyFriction()
	}

	gr.updatePlayerAttacks(inputState)
	gr.updatePlayerJump(inputState)
	gr.updatePlayerDash(inputState)
	gr.updatePlayerGrapple(inputState)

	gr.inputHandler.UpdateBuffers()
}

// updatePlayerAttacks handles melee and ranged attack input with buffering.
func (gr *GameRunner) updatePlayerAttacks(inputState input.InputState) {
	if inputState.AttackPress {
		if !gr.combatSystem.PlayerAttack() {
			gr.inputHandler.BufferAttack()
		}
	}
	if gr.inputHandler.GetBufferedAttack() && gr.combatSystem.CanAttack() {
		gr.combatSystem.PlayerAttack()
	}
	if inputState.RangedAttackPress && gr.game.Player.Abilities["ranged"] {
		gr.combatSystem.PlayerRangedAttack(
			gr.game.Player.X, gr.game.Player.Y,
			gr.playerFacingDir, gr.game.Player.Damage,
		)
	}
}

// updatePlayerJump handles jump input and buffering.
func (gr *GameRunner) updatePlayerJump(inputState input.InputState) {
	if inputState.JumpPress {
		hasDoubleJump := gr.game.Player.Abilities["double_jump"]
		if gr.playerBody.Jump(hasDoubleJump, &gr.doubleJumpUsed) {
			emitter := gr.particlePresets.CreateJumpDust(gr.game.Player.X+16, gr.game.Player.Y+32)
			emitter.Burst(8)
			gr.particleSystem.AddEmitter(emitter)
		} else {
			gr.playerBody.BufferJump()
		}
	}
	if inputState.JumpRelease {
		gr.playerBody.ReleaseJump()
	}
}

// updatePlayerDash handles dash input with cooldown and buffering.
func (gr *GameRunner) updatePlayerDash(inputState input.InputState) {
	if gr.dashCooldown > 0 {
		gr.dashCooldown--
	}
	hasDash := gr.game.Player.Abilities["dash"]
	if !hasDash {
		return
	}

	dir := dashDirection(inputState)

	if inputState.DashPress {
		if gr.dashCooldown <= 0 {
			gr.executeDash(dir)
		} else {
			gr.inputHandler.BufferDash()
		}
	}
	if gr.inputHandler.GetBufferedDash() && gr.dashCooldown <= 0 {
		gr.executeDash(dir)
	}
}

// dashDirection returns the dash direction derived from the current input state.
func dashDirection(inputState input.InputState) float64 {
	switch {
	case inputState.MoveRight:
		return 1.0
	case inputState.MoveLeft:
		return -1.0
	default:
		return 0.0
	}
}

// executeDash performs the dash and emits trail particles.
func (gr *GameRunner) executeDash(direction float64) {
	gr.playerBody.Dash(direction)
	gr.dashCooldown = 30
	emitter := gr.particlePresets.CreateDashTrail(gr.game.Player.X+16, gr.game.Player.Y+16)
	emitter.Start()
	gr.particleSystem.AddEmitter(emitter)
}

// updatePlayerGrapple handles grapple hook activation and release.
func (gr *GameRunner) updatePlayerGrapple(inputState input.InputState) {
	if gr.grappleCooldown > 0 {
		gr.grappleCooldown--
	}
	if inputState.UseAbility && gr.grappleCooldown <= 0 && !gr.playerBody.Grappling {
		if gr.game.Player.Abilities["grapple"] && gr.game.CurrentRoom != nil {
			if anchor, found := physics.FindNearestAnchor(gr.playerBody.Position, gr.game.CurrentRoom.Anchors); found {
				gr.playerBody.StartGrapple(anchor)
				gr.grappleCooldown = 15
			}
		}
	}
	if gr.playerBody.Grappling {
		gr.playerBody.UpdateGrapple()
		if !inputState.UseAbility {
			gr.playerBody.ReleaseGrapple()
		}
	}
}

// updatePlayerPhysics applies knockback, integrates velocity, resolves
// collisions, and syncs the game-state player position.
func (gr *GameRunner) updatePlayerPhysics(wasOnGround bool) {
	knockbackX, knockbackY := gr.combatSystem.GetKnockback()
	if knockbackX != 0 || knockbackY != 0 {
		gr.playerBody.Velocity.X += knockbackX
		gr.playerBody.Velocity.Y += knockbackY
	}

	gr.playerBody.Update()

	if gr.game.CurrentRoom != nil {
		gr.playerBody.ResolveCollisionWithPlatforms(gr.game.CurrentRoom.Platforms)
	}

	if !wasOnGround && gr.playerBody.OnGround {
		emitter := gr.particlePresets.CreateLandDust(gr.game.Player.X+16, gr.game.Player.Y+32)
		emitter.Burst(12)
		gr.particleSystem.AddEmitter(emitter)
	}

	gr.game.Player.X = gr.playerBody.Position.X
	gr.game.Player.Y = gr.playerBody.Position.Y
	gr.game.Player.VelX = gr.playerBody.Velocity.X
	gr.game.Player.VelY = gr.playerBody.Velocity.Y
}

// updatePlayerAnimation drives the animation state machine based on movement
// and combat state.
func (gr *GameRunner) updatePlayerAnimation(inputState input.InputState) {
	if gr.game.Player.AnimController == nil {
		return
	}
	gr.game.Player.AnimController.Update()
	currentAnim := gr.game.Player.AnimController.GetCurrentAnimation()

	switch {
	case gr.combatSystem.IsPlayerAttacking():
		if currentAnim != "attack" {
			gr.game.Player.AnimController.Play("attack", true)
		}
	case !gr.playerBody.OnGround:
		if currentAnim != "jump" {
			gr.game.Player.AnimController.Play("jump", true)
		}
	case inputState.MoveLeft || inputState.MoveRight:
		if currentAnim != "walk" {
			gr.game.Player.AnimController.Play("walk", false)
		}
	default:
		if currentAnim != "idle" {
			gr.game.Player.AnimController.Play("idle", false)
		}
	}
}

// updateEnemies runs AI, physics, and combat interactions for all enemies.
func (gr *GameRunner) updateEnemies() {
	for _, enemy := range gr.enemyInstances {
		if enemy.IsDead() {
			continue
		}
		gr.updateSingleEnemy(enemy)
	}
}

// updateSingleEnemy handles AI, physics, and combat for one enemy instance.
func (gr *GameRunner) updateSingleEnemy(enemy *entity.EnemyInstance) {
	enemy.Update(gr.game.Player.X, gr.game.Player.Y)
	gr.applyEnemyGravity(enemy)
	enemy.X += enemy.VelX
	enemy.Y += enemy.VelY
	gr.resolveEnemyPlatformCollisions(enemy)
	gr.checkMeleeHitEnemy(enemy)
	gr.checkProjectileHitEnemy(enemy)
	gr.checkEnemyHitPlayer(enemy)
}

// applyEnemyGravity applies gravity to ground-based (non-flying) enemies.
func (gr *GameRunner) applyEnemyGravity(enemy *entity.EnemyInstance) {
	if enemy.Enemy.Behavior == entity.FlyingBehavior || enemy.OnGround {
		return
	}
	enemy.VelY += 0.5
	if enemy.VelY > 10.0 {
		enemy.VelY = 10.0
	}
}

// resolveEnemyPlatformCollisions pushes an enemy out of any overlapping platforms.
func (gr *GameRunner) resolveEnemyPlatformCollisions(enemy *entity.EnemyInstance) {
	if gr.game.CurrentRoom == nil {
		return
	}
	enemy.OnGround = false
	ex, ey, ew, eh := enemy.GetBounds()
	for _, platform := range gr.game.CurrentRoom.Platforms {
		px := float64(platform.X)
		py := float64(platform.Y)
		pw := float64(platform.Width)
		ph := float64(platform.Height)
		if ex < px+pw && ex+ew > px && ey < py+ph && ey+eh > py {
			if enemy.VelY > 0 && ey+eh-ph < py {
				enemy.Y = py - eh
				enemy.VelY = 0
				enemy.OnGround = true
			}
		}
	}
}

// checkMeleeHitEnemy tests whether the current player melee attack hits the
// given enemy and applies damage and particle effects if so.
func (gr *GameRunner) checkMeleeHitEnemy(enemy *entity.EnemyInstance) {
	if !gr.combatSystem.IsPlayerAttacking() {
		return
	}
	attackX, attackY, attackW, attackH := gr.combatSystem.GetAttackHitbox(
		gr.game.Player.X, gr.game.Player.Y, gr.playerFacingDir,
	)
	wasAlive := !enemy.IsDead()
	if !gr.combatSystem.CheckEnemyHit(attackX, attackY, attackW, attackH, enemy) {
		return
	}

	ex, ey, _, _ := enemy.GetBounds()
	hitEmitter := gr.particlePresets.CreateHitEffect(ex+16, ey+16, gr.playerFacingDir)
	hitEmitter.Burst(10)
	gr.particleSystem.AddEmitter(hitEmitter)

	bloodEmitter := gr.particlePresets.CreateBloodSplatter(ex+16, ey+16, gr.playerFacingDir)
	bloodEmitter.Burst(6)
	gr.particleSystem.AddEmitter(bloodEmitter)

	gr.combatSystem.ApplyDamageToEnemy(enemy, gr.game.Player.Damage, gr.game.Player.X)

	if gr.game.Achievements != nil {
		gr.game.Achievements.RecordDamage(gr.game.Player.Damage, 0)
	}
	if wasAlive && enemy.IsDead() {
		gr.recordEnemyDeath(enemy)
	}
}

// checkProjectileHitEnemy tests whether any active projectile hits the given
// enemy and applies damage and particle effects if so.
func (gr *GameRunner) checkProjectileHitEnemy(enemy *entity.EnemyInstance) {
	projDmg := gr.combatSystem.CheckProjectileEnemyHit(enemy)
	if projDmg <= 0 {
		return
	}
	ex, ey, _, _ := enemy.GetBounds()
	hitEmitter := gr.particlePresets.CreateHitEffect(ex+16, ey+16, gr.playerFacingDir)
	hitEmitter.Burst(6)
	gr.particleSystem.AddEmitter(hitEmitter)
	if gr.game.Achievements != nil {
		gr.game.Achievements.RecordDamage(projDmg, 0)
	}
	if enemy.IsDead() {
		gr.recordEnemyDeath(enemy)
	}
}

// recordEnemyDeath updates tracking maps and fires the death explosion.
func (gr *GameRunner) recordEnemyDeath(enemy *entity.EnemyInstance) {
	ex, ey, _, _ := enemy.GetBounds()
	enemyKey := int(enemy.X*1000 + enemy.Y)
	gr.defeatedEnemies[enemyKey] = true
	if gr.game.Achievements != nil {
		wasPerfect := gr.combatSystem.GetInvulnerableFrames() == 0
		gr.game.Achievements.RecordEnemyKill(wasPerfect)
	}
	explosionEmitter := gr.particlePresets.CreateExplosion(ex+16, ey+16, 1.0)
	explosionEmitter.Burst(20)
	gr.particleSystem.AddEmitter(explosionEmitter)
}

// checkEnemyHitPlayer tests whether the given enemy is colliding with the
// player and applies damage if so.
func (gr *GameRunner) checkEnemyHitPlayer(enemy *entity.EnemyInstance) {
	if !gr.combatSystem.CheckPlayerEnemyCollision(
		gr.game.Player.X, gr.game.Player.Y, physics.PlayerWidth, physics.PlayerHeight, enemy,
	) {
		return
	}
	damage := enemy.Enemy.Damage
	if enemy.State == entity.AttackState {
		damage = enemy.GetAttackDamage()
	}
	if gr.game.Achievements != nil {
		gr.game.Achievements.RecordDamage(0, damage)
	}
	gr.combatSystem.ApplyDamageToPlayer(gr.game.Player, damage, enemy.X)
	if gr.game.Player.Health <= 0 && gr.game.Achievements != nil {
		gr.game.Achievements.RecordDeath()
	}
}

// updateRoomTracking marks the current room as visited and fires the first-
// visit achievement event.
func (gr *GameRunner) updateRoomTracking() {
	if gr.game.CurrentRoom == nil {
		return
	}
	wasVisited := gr.visitedRooms[gr.game.CurrentRoom.ID]
	gr.visitedRooms[gr.game.CurrentRoom.ID] = true
	if !wasVisited && gr.game.Achievements != nil {
		isPerfect := !gr.combatSystem.IsInvulnerable()
		gr.game.Achievements.RecordRoomVisit(isPerfect)
	}
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

	// Render particles and other ECS-managed visuals
	gr.systemManager.Draw(screen)

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
		transitionType := string(gr.transitionHandler.GetTransitionType())
		slideDirection := gr.transitionHandler.GetSlideDirection()
		gr.renderer.RenderTransitionEffect(screen, progress, transitionType, slideDirection)
	}

	// Show locked door message if active (centered on screen)
	if gr.lockedDoorTimer > 0 && gr.lockedDoorMessage != "" {
		msgX := (render.ScreenWidth - render.MessageWidth) / 2
		msgY := (render.ScreenHeight - render.MessageHeight) / 2
		gr.renderMessageWithProgress(screen, gr.lockedDoorMessage, gr.lockedDoorTimer, lockedDoorMessageDuration,
			msgX, msgY, color.RGBA{0, 0, 0, 180})
	}

	// Show item collection message if active (below HUD area)
	if gr.itemMessageTimer > 0 && gr.itemMessage != "" {
		msgX := (render.ScreenWidth - render.MessageWidth) / 2
		msgY := render.AbilityIconY + render.AbilityIconSize + render.UIMargin
		gr.renderMessageWithProgress(screen, gr.itemMessage, gr.itemMessageTimer, itemMessageDuration,
			msgX, msgY, color.RGBA{255, 215, 0, 200})
	}

	// Show debug info if enabled (positioned below UI to avoid overlap)
	if gr.showDebugInfo {
		aliveEnemies := 0
		for _, enemy := range gr.enemyInstances {
			if !enemy.IsDead() {
				aliveEnemies++
			}
		}

		debugInfo := fmt.Sprintf("Seed: %d | Room: %s | FPS: %.2f | Enemies: %d/%d | Items: %d/%d\nPosition: (%.0f, %.0f) | Velocity: (%.1f, %.1f)\nHealth: %d/%d | OnGround: %v | Invuln: %v\nControls: WASD/Arrows=Move, Space=Jump, J=Attack, K=Dash, P=Pause, F3=Debug, Ctrl+Q=Quit",
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

		// Position debug info below health bar and abilities
		debugX := render.UIMargin
		debugY := render.AbilityIconY + render.AbilityIconSize + render.UIMargin + render.MessageHeight + render.UIMargin
		// Use text rendering abstraction with fallback to debug text
		if gr.renderer != nil {
			gr.renderer.RenderText(screen, debugInfo, debugX, debugY, color.RGBA{255, 255, 255, 255})
		} else {
			ebitenutil.DebugPrintAt(screen, debugInfo, debugX, debugY)
		}
	}

	// Always show minimal controls hint in top-right corner if debug is off
	if !gr.showDebugInfo {
		controlsHint := "F3=Debug Info"
		hintWidth := len(controlsHint) * 8 // 8px per char (standardized font metrics)
		hintX := render.ScreenWidth - hintWidth - render.UIMargin
		hintY := render.UIMargin
		if gr.renderer != nil {
			gr.renderer.RenderText(screen, controlsHint, hintX, hintY, color.RGBA{200, 200, 200, 255})
		} else {
			ebitenutil.DebugPrintAt(screen, controlsHint, hintX, hintY)
		}
	}
}

// Layout implements ebiten.Game interface
func (gr *GameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return render.ScreenWidth, render.ScreenHeight
}

// SetGenre propagates a genre change to all ECS-registered subsystems.
// genreID must be one of: "fantasy", "scifi", "horror", "cyberpunk", "postapoc".
// The SystemManager broadcasts the change to all registered System implementations.
func (gr *GameRunner) SetGenre(genreID string) {
	gr.game.Genre = genreID
	gr.renderer.SetGenre(genreID)
	gr.systemManager.SetGenre(genreID)
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

	// Build achievement statistics
	var achievementStats *save.AchievementStatistics
	if gr.game.Achievements != nil {
		stats := gr.game.Achievements.GetStatistics()
		achievementStats = &save.AchievementStatistics{
			EnemiesDefeated:   stats.EnemiesDefeated,
			BossesDefeated:    stats.BossesDefeated,
			TotalDamageDealt:  stats.TotalDamageDealt,
			DamageTaken:       stats.DamageTaken,
			PerfectKills:      stats.PerfectKills,
			RoomsVisited:      stats.RoomsVisited,
			BiomesExplored:    stats.BiomesExplored,
			SecretsFound:      stats.SecretsFound,
			ItemsCollected:    stats.ItemsCollected,
			AbilitiesUnlocked: stats.AbilitiesUnlocked,
			DeathCount:        stats.DeathCount,
			PerfectRooms:      stats.PerfectRooms,
			ConsecutiveKills:  stats.ConsecutiveKills,
			LongestCombo:      stats.LongestCombo,
		}
	}

	return &save.SaveData{
		Seed:             gr.game.Seed,
		PlayTime:         playTime,
		PlayerX:          gr.game.Player.X,
		PlayerY:          gr.game.Player.Y,
		PlayerHealth:     gr.game.Player.Health,
		PlayerMaxHealth:  gr.game.Player.MaxHealth,
		PlayerAbilities:  gr.game.Player.Abilities,
		CurrentRoomID:    currentRoomID,
		VisitedRooms:     visitedRoomsList,
		DefeatedEnemies:  gr.defeatedEnemies,
		CollectedItems:   gr.collectedItems,
		UnlockedDoors:    gr.unlockedDoors,
		BossesDefeated:   gr.getBossesDefeated(),
		CheckpointID:     currentRoomID,
		AchievementStats: achievementStats,
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
					gr.lockedDoorTimer = lockedDoorMessageDuration
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
	gr.lockedDoorTimer = lockedDoorMessageDuration

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

	// Record item collection for achievements
	if gr.game.Achievements != nil {
		gr.game.Achievements.RecordItemCollected()
	}

	// Show message
	gr.itemMessage = fmt.Sprintf("Collected: %s", item.Item.Name)
	gr.itemMessageTimer = itemMessageDuration

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

	// Check if item grants an ability (for key items)
	if item.Item.Type == entity.KeyItem {
		// Key items grant abilities
		abilityName := item.Item.Name // Use item name as ability identifier
		if !gr.game.Player.Abilities[abilityName] {
			gr.game.Player.Abilities[abilityName] = true

			// Record ability unlock for achievements
			if gr.game.Achievements != nil {
				gr.game.Achievements.RecordAbilityUnlocked()
			}
		}
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

	// Restore achievement statistics if available
	if saveData.AchievementStats != nil && gr.game.Achievements != nil {
		stats := gr.game.Achievements.GetStatistics()
		stats.EnemiesDefeated = saveData.AchievementStats.EnemiesDefeated
		stats.BossesDefeated = saveData.AchievementStats.BossesDefeated
		stats.TotalDamageDealt = saveData.AchievementStats.TotalDamageDealt
		stats.DamageTaken = saveData.AchievementStats.DamageTaken
		stats.PerfectKills = saveData.AchievementStats.PerfectKills
		stats.RoomsVisited = saveData.AchievementStats.RoomsVisited
		stats.BiomesExplored = saveData.AchievementStats.BiomesExplored
		stats.SecretsFound = saveData.AchievementStats.SecretsFound
		stats.ItemsCollected = saveData.AchievementStats.ItemsCollected
		stats.AbilitiesUnlocked = saveData.AchievementStats.AbilitiesUnlocked
		stats.DeathCount = saveData.AchievementStats.DeathCount
		stats.PerfectRooms = saveData.AchievementStats.PerfectRooms
		stats.ConsecutiveKills = saveData.AchievementStats.ConsecutiveKills
		stats.LongestCombo = saveData.AchievementStats.LongestCombo
		stats.PlayTime = saveData.PlayTime
		gr.game.Achievements.UpdateStatistics(stats)
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

	// Find ground platform Y for spawning
	groundY := findGroundY(room)

	// Place 2-4 items in treasure rooms
	itemCount := 2 + (room.ID % 3) // 2-4 items based on room ID
	if itemCount > len(allItems) {
		itemCount = len(allItems)
	}

	for i := 0; i < itemCount && i < len(allItems); i++ {
		// Generate unique item ID based on room and position
		itemID := room.ID*1000 + i

		// Position items across the room on the ground platform
		itemX := 200.0 + float64(i*150)
		itemY := groundY - 16.0 // Items are 16px tall

		instance := entity.NewItemInstance(allItems[i%len(allItems)], itemID, itemX, itemY)
		instances = append(instances, instance)
	}

	return instances
}

// findGroundY returns the Y coordinate of the ground platform surface in a room.
// Falls back to screen-bottom floor if no ground platform is found.
func findGroundY(room *world.Room) float64 {
	if room == nil {
		return float64(render.ScreenHeight - 32)
	}
	// Look for the lowest full-width platform (ground platform)
	groundY := float64(render.ScreenHeight) // Default: screen bottom
	for _, p := range room.Platforms {
		if p.Width >= render.ScreenWidth/2 && float64(p.Y) < groundY {
			// Found a wide platform — use its top surface
			groundY = float64(p.Y)
		}
	}
	return groundY
}

// renderMessageWithProgress renders a message with a progress bar showing remaining time
func (gr *GameRunner) renderMessageWithProgress(screen *ebiten.Image, message string, currentTimer, maxDuration, x, y int, bgColor color.Color) {
	messageWidth := render.MessageWidth
	messageHeight := render.MessageHeight
	progressHeight := render.ProgressBarHeight

	// Calculate progress (0.0 to 1.0)
	progress := float64(currentTimer) / float64(maxDuration)
	if progress > 1.0 {
		progress = 1.0
	}
	if progress < 0.0 {
		progress = 0.0
	}

	// Draw in correct order: border → background → progress → text

	// 1. Draw subtle border for better definition
	borderImg := ebiten.NewImage(messageWidth+2, messageHeight+2)
	borderImg.Fill(color.RGBA{255, 255, 255, 50})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x-1), float64(y-1))
	screen.DrawImage(borderImg, opts)

	// 2. Draw main message background
	messageImg := ebiten.NewImage(messageWidth, messageHeight)
	messageImg.Fill(bgColor)
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(messageImg, opts)

	// 3. Draw progress bar background (darker)
	progressBgImg := ebiten.NewImage(messageWidth-4, progressHeight)
	progressBgImg.Fill(color.RGBA{0, 0, 0, 100})
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x+2), float64(y+messageHeight-progressHeight-2))
	screen.DrawImage(progressBgImg, opts)

	// 4. Draw progress bar fill (progress remaining)
	progressWidth := int(float64(messageWidth-4) * progress)
	if progressWidth > 0 {
		progressImg := ebiten.NewImage(progressWidth, progressHeight)

		// Color progress bar based on remaining time
		var progressColor color.Color
		if progress > 0.5 {
			// Green when lots of time left
			progressColor = color.RGBA{100, 255, 100, 200}
		} else if progress > 0.2 {
			// Yellow when half time left
			progressColor = color.RGBA{255, 255, 100, 200}
		} else {
			// Red when almost out of time
			progressColor = color.RGBA{255, 100, 100, 200}
		}

		progressImg.Fill(progressColor)
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x+2), float64(y+messageHeight-progressHeight-2))
		screen.DrawImage(progressImg, opts)
	}

	// 5. Draw text last so it's visible on top of background
	if gr.renderer != nil {
		gr.renderer.RenderText(screen, message, x+10, y+12, color.RGBA{255, 255, 255, 255})
	} else {
		ebitenutil.DebugPrintAt(screen, message, x+10, y+12)
	}
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
