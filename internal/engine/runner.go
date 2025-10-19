// Package engine provides the game runner that integrates Ebiten rendering
// with the procedural generation system, handling the game loop, player
// movement, and visual display.
package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/input"
	"github.com/opd-ai/vania/internal/physics"
	"github.com/opd-ai/vania/internal/render"
)

// GameRunner wraps the Game with Ebiten rendering
type GameRunner struct {
	game           *Game
	renderer       *render.Renderer
	inputHandler   *input.InputHandler
	playerBody     *physics.Body
	combatSystem   *CombatSystem
	enemyInstances []*entity.EnemyInstance
	doubleJumpUsed bool
	dashCooldown   int
	playerFacingDir float64
	paused         bool
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
	
	return &GameRunner{
		game:           game,
		renderer:       render.NewRenderer(),
		inputHandler:   input.NewInputHandler(),
		playerBody:     physics.NewBody(playerX, playerY, physics.PlayerWidth, physics.PlayerHeight),
		combatSystem:   NewCombatSystem(),
		enemyInstances: enemyInstances,
		doubleJumpUsed: false,
		dashCooldown:   0,
		playerFacingDir: 1.0,
		paused:         false,
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
	
	// Apply physics
	gr.playerBody.ApplyGravity()
	
	// Update combat system
	gr.combatSystem.Update()
	
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
		gr.playerBody.Jump(hasDoubleJump, &gr.doubleJumpUsed)
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
	
	// Update player position in game
	gr.game.Player.X = gr.playerBody.Position.X
	gr.game.Player.Y = gr.playerBody.Position.Y
	gr.game.Player.VelX = gr.playerBody.Velocity.X
	gr.game.Player.VelY = gr.playerBody.Velocity.Y
	
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
			if gr.combatSystem.CheckEnemyHit(attackX, attackY, attackW, attackH, enemy) {
				gr.combatSystem.ApplyDamageToEnemy(enemy, gr.game.Player.Damage, gr.game.Player.X)
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
	
	// Update camera
	gr.renderer.UpdateCamera(gr.game.Player.X, gr.game.Player.Y)
	
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
	
	// Render enemies
	for _, enemy := range gr.enemyInstances {
		if !enemy.IsDead() {
			ex, ey, ew, eh := enemy.GetBounds()
			gr.renderer.RenderEnemy(screen, ex, ey, ew, eh, enemy.CurrentHealth, enemy.Enemy.Health, false)
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
	
	// Render player
	if gr.game.Player != nil {
		gr.renderer.RenderPlayer(screen, gr.game.Player.X, gr.game.Player.Y, gr.game.Player.Sprite)
	}
	
	// Render UI
	if gr.game.Player != nil {
		gr.renderer.RenderUI(screen, gr.game.Player.Health, gr.game.Player.MaxHealth, gr.game.Player.Abilities)
	}
	
	// Show debug info
	aliveEnemies := 0
	for _, enemy := range gr.enemyInstances {
		if !enemy.IsDead() {
			aliveEnemies++
		}
	}
	
	debugInfo := fmt.Sprintf("Seed: %d | Room: %s | FPS: %.2f | Enemies: %d/%d\nPosition: (%.0f, %.0f) | Velocity: (%.1f, %.1f)\nHealth: %d/%d | OnGround: %v | Invuln: %v\nControls: WASD/Arrows=Move, Space=Jump, J=Attack, K=Dash, P=Pause, Ctrl+Q=Quit",
		gr.game.Seed,
		gr.getCurrentRoomName(),
		ebiten.ActualTPS(),
		aliveEnemies,
		len(gr.enemyInstances),
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
