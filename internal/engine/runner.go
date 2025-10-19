// Package engine provides the game runner that integrates Ebiten rendering
// with the procedural generation system, handling the game loop, player
// movement, and visual display.
package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/vania/internal/input"
	"github.com/opd-ai/vania/internal/physics"
	"github.com/opd-ai/vania/internal/render"
)

// GameRunner wraps the Game with Ebiten rendering
type GameRunner struct {
	game          *Game
	renderer      *render.Renderer
	inputHandler  *input.InputHandler
	playerBody    *physics.Body
	doubleJumpUsed bool
	dashCooldown  int
	paused        bool
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
	
	return &GameRunner{
		game:          game,
		renderer:      render.NewRenderer(),
		inputHandler:  input.NewInputHandler(),
		playerBody:    physics.NewBody(playerX, playerY, physics.PlayerWidth, physics.PlayerHeight),
		doubleJumpUsed: false,
		dashCooldown:  0,
		paused:        false,
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
	
	// Handle player movement
	if inputState.MoveLeft {
		gr.playerBody.MoveHorizontal(-1)
	} else if inputState.MoveRight {
		gr.playerBody.MoveHorizontal(1)
	} else {
		gr.playerBody.ApplyFriction()
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
	
	// Render player
	if gr.game.Player != nil {
		gr.renderer.RenderPlayer(screen, gr.game.Player.X, gr.game.Player.Y, gr.game.Player.Sprite)
	}
	
	// Render UI
	if gr.game.Player != nil {
		gr.renderer.RenderUI(screen, gr.game.Player.Health, gr.game.Player.MaxHealth, gr.game.Player.Abilities)
	}
	
	// Show debug info
	debugInfo := fmt.Sprintf("Seed: %d | Room: %s | FPS: %.2f\nPosition: (%.0f, %.0f) | Velocity: (%.1f, %.1f)\nOnGround: %v | Controls: WASD/Arrows=Move, Space=Jump, K=Dash, P=Pause, Ctrl+Q=Quit",
		gr.game.Seed,
		gr.getCurrentRoomName(),
		ebiten.ActualTPS(),
		gr.game.Player.X,
		gr.game.Player.Y,
		gr.game.Player.VelX,
		gr.game.Player.VelY,
		gr.playerBody.OnGround,
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
