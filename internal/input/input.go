// Package input handles player input from keyboard and game controllers,
// providing a unified interface for movement, actions, and menu navigation.
package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputState represents the current input state
type InputState struct {
	MoveLeft   bool
	MoveRight  bool
	Jump       bool
	JumpPress  bool // True only on the frame jump was pressed
	Attack     bool
	AttackPress bool
	Dash       bool
	DashPress  bool
	UseAbility bool
	Pause      bool
	PausePress bool
}

// InputHandler manages input processing
type InputHandler struct {
	prevState InputState
}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{
		prevState: InputState{},
	}
}

// Update reads current input state
func (ih *InputHandler) Update() InputState {
	state := InputState{}
	
	// Movement
	state.MoveLeft = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	state.MoveRight = ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	
	// Jump
	state.Jump = ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)
	state.JumpPress = inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp)
	
	// Attack
	state.Attack = ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyZ)
	state.AttackPress = inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	
	// Dash
	state.Dash = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyX)
	state.DashPress = inpututil.IsKeyJustPressed(ebiten.KeyK) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	
	// Use ability
	state.UseAbility = ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyC)
	
	// Pause
	state.Pause = ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyP)
	state.PausePress = inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP)
	
	// Update previous state
	ih.prevState = state
	
	return state
}

// IsQuitRequested checks if the user wants to quit
func (ih *InputHandler) IsQuitRequested() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyQ) && ebiten.IsKeyPressed(ebiten.KeyControl)
}
