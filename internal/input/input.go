// Package input handles player input from keyboard and game controllers,
// providing a unified interface for movement, actions, and menu navigation.
package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// BufferFrames is the window (in frames at 60fps) during which actions
	// are buffered and will execute when conditions are met. Industry standard: ~100ms = 6 frames.
	BufferFrames = 6
)

// InputState represents the current input state
type InputState struct {
	MoveLeft    bool
	MoveRight   bool
	Jump        bool
	JumpPress   bool // True only on the frame jump was pressed
	JumpRelease bool // True only on the frame jump was released
	Attack      bool
	AttackPress bool
	Dash        bool
	DashPress   bool
	UseAbility  bool
	Block       bool // Hold to block/parry
	BlockPress  bool // True only on the frame block was pressed
	Pause       bool
	PausePress  bool
}

// BufferedInput tracks buffered action inputs
type BufferedInput struct {
	AttackBuffer int // Countdown timer for buffered attack
	DashBuffer   int // Countdown timer for buffered dash
}

// InputHandler manages input processing with rebindable controls
type InputHandler struct {
	prevState  InputState
	buffered   BufferedInput
	keyMapping *KeyMapping
}

// KeyMapping defines rebindable key bindings for all actions
type KeyMapping struct {
	MoveLeft   []ebiten.Key
	MoveRight  []ebiten.Key
	Jump       []ebiten.Key
	Attack     []ebiten.Key
	Dash       []ebiten.Key
	UseAbility []ebiten.Key
	Block      []ebiten.Key
	Pause      []ebiten.Key
}

// DefaultKeyMapping returns the default key configuration
func DefaultKeyMapping() *KeyMapping {
	return &KeyMapping{
		MoveLeft:   []ebiten.Key{ebiten.KeyA, ebiten.KeyArrowLeft},
		MoveRight:  []ebiten.Key{ebiten.KeyD, ebiten.KeyArrowRight},
		Jump:       []ebiten.Key{ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp},
		Attack:     []ebiten.Key{ebiten.KeyJ, ebiten.KeyZ},
		Dash:       []ebiten.Key{ebiten.KeyK, ebiten.KeyX},
		UseAbility: []ebiten.Key{ebiten.KeyL, ebiten.KeyC},
		Block:      []ebiten.Key{ebiten.KeyS, ebiten.KeyArrowDown, ebiten.KeyShiftLeft},
		Pause:      []ebiten.Key{ebiten.KeyEscape, ebiten.KeyP},
	}
}

// NewInputHandler creates a new input handler with default key bindings
func NewInputHandler() *InputHandler {
	return &InputHandler{
		prevState:  InputState{},
		buffered:   BufferedInput{},
		keyMapping: DefaultKeyMapping(),
	}
}

// NewInputHandlerWithMapping creates a new input handler with custom key bindings
func NewInputHandlerWithMapping(mapping *KeyMapping) *InputHandler {
	return &InputHandler{
		prevState:  InputState{},
		buffered:   BufferedInput{},
		keyMapping: mapping,
	}
}

// SetKeyMapping updates the key bindings
func (ih *InputHandler) SetKeyMapping(mapping *KeyMapping) {
	ih.keyMapping = mapping
}

// isAnyKeyPressed checks if any key in the list is pressed
func (ih *InputHandler) isAnyKeyPressed(keys []ebiten.Key) bool {
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	return false
}

// isAnyKeyJustPressed checks if any key in the list was just pressed
func (ih *InputHandler) isAnyKeyJustPressed(keys []ebiten.Key) bool {
	for _, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			return true
		}
	}
	return false
}

// isAnyKeyJustReleased checks if any key in the list was just released
func (ih *InputHandler) isAnyKeyJustReleased(keys []ebiten.Key) bool {
	for _, key := range keys {
		if inpututil.IsKeyJustReleased(key) {
			return true
		}
	}
	return false
}

// Update reads current input state using configured key bindings
func (ih *InputHandler) Update() InputState {
	state := InputState{}

	// Movement (not buffered)
	state.MoveLeft = ih.isAnyKeyPressed(ih.keyMapping.MoveLeft)
	state.MoveRight = ih.isAnyKeyPressed(ih.keyMapping.MoveRight)

	// Jump (physics system handles buffering)
	state.Jump = ih.isAnyKeyPressed(ih.keyMapping.Jump)
	state.JumpPress = ih.isAnyKeyJustPressed(ih.keyMapping.Jump)
	state.JumpRelease = ih.isAnyKeyJustReleased(ih.keyMapping.Jump)

	// Attack (buffered)
	state.Attack = ih.isAnyKeyPressed(ih.keyMapping.Attack)
	state.AttackPress = ih.isAnyKeyJustPressed(ih.keyMapping.Attack)

	// Dash (buffered)
	state.Dash = ih.isAnyKeyPressed(ih.keyMapping.Dash)
	state.DashPress = ih.isAnyKeyJustPressed(ih.keyMapping.Dash)

	// Use ability (not buffered)
	state.UseAbility = ih.isAnyKeyPressed(ih.keyMapping.UseAbility)

	// Block/Parry (not buffered - requires precise timing)
	state.Block = ih.isAnyKeyPressed(ih.keyMapping.Block)
	state.BlockPress = ih.isAnyKeyJustPressed(ih.keyMapping.Block)

	// Pause (not buffered)
	state.Pause = ih.isAnyKeyPressed(ih.keyMapping.Pause)
	state.PausePress = ih.isAnyKeyJustPressed(ih.keyMapping.Pause)

	// Update previous state
	ih.prevState = state

	return state
}

// BufferAttack buffers an attack input for execution when conditions are met
func (ih *InputHandler) BufferAttack() {
	ih.buffered.AttackBuffer = BufferFrames
}

// BufferDash buffers a dash input for execution when conditions are met
func (ih *InputHandler) BufferDash() {
	ih.buffered.DashBuffer = BufferFrames
}

// GetBufferedAttack returns true if attack is buffered and consumes the buffer
func (ih *InputHandler) GetBufferedAttack() bool {
	if ih.buffered.AttackBuffer > 0 {
		ih.buffered.AttackBuffer = 0
		return true
	}
	return false
}

// GetBufferedDash returns true if dash is buffered and consumes the buffer
func (ih *InputHandler) GetBufferedDash() bool {
	if ih.buffered.DashBuffer > 0 {
		ih.buffered.DashBuffer = 0
		return true
	}
	return false
}

// UpdateBuffers decrements buffer timers (call once per frame)
func (ih *InputHandler) UpdateBuffers() {
	if ih.buffered.AttackBuffer > 0 {
		ih.buffered.AttackBuffer--
	}
	if ih.buffered.DashBuffer > 0 {
		ih.buffered.DashBuffer--
	}
}

// IsQuitRequested checks if the user wants to quit
func (ih *InputHandler) IsQuitRequested() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyQ) && ebiten.IsKeyPressed(ebiten.KeyControl)
}
