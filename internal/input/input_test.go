package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/vania/internal/settings"
)

func TestNewInputHandler(t *testing.T) {
	ih := NewInputHandler()

	if ih == nil {
		t.Fatal("NewInputHandler returned nil")
	}

	// Verify default key mapping is set
	if ih.keyMapping == nil {
		t.Error("Default key mapping should be initialized")
	}

	// Verify buffers are initialized
	if ih.buffered.AttackBuffer != 0 || ih.buffered.DashBuffer != 0 {
		t.Error("Buffers should initialize to 0")
	}
}

func TestNewInputHandlerWithMapping(t *testing.T) {
	customMapping := &KeyMapping{
		Jump: []ebiten.Key{ebiten.KeyA},
	}

	ih := NewInputHandlerWithMapping(customMapping)

	if ih == nil {
		t.Fatal("NewInputHandlerWithMapping returned nil")
	}

	if ih.keyMapping != customMapping {
		t.Error("Custom key mapping should be set")
	}
}

func TestSetKeyMapping(t *testing.T) {
	ih := NewInputHandler()

	customMapping := &KeyMapping{
		Jump: []ebiten.Key{ebiten.KeyB},
	}

	ih.SetKeyMapping(customMapping)

	if ih.keyMapping != customMapping {
		t.Error("SetKeyMapping should update the key mapping")
	}
}

func TestDefaultKeyMapping(t *testing.T) {
	mapping := DefaultKeyMapping()

	if mapping == nil {
		t.Fatal("DefaultKeyMapping returned nil")
	}

	// Verify all actions have key bindings
	if len(mapping.MoveLeft) == 0 {
		t.Error("MoveLeft should have default keys")
	}
	if len(mapping.Jump) == 0 {
		t.Error("Jump should have default keys")
	}
	if len(mapping.Attack) == 0 {
		t.Error("Attack should have default keys")
	}
	if len(mapping.Dash) == 0 {
		t.Error("Dash should have default keys")
	}
}

func TestInputStateInitialization(t *testing.T) {
	state := InputState{}

	// All fields should be false initially
	if state.MoveLeft || state.MoveRight || state.Jump || state.Attack {
		t.Error("InputState should initialize with all false values")
	}

	if state.JumpPress || state.AttackPress || state.DashPress || state.PausePress {
		t.Error("InputState press flags should initialize as false")
	}

	if state.JumpRelease {
		t.Error("InputState JumpRelease should initialize as false")
	}
}

func TestInputStateFields(t *testing.T) {
	// Verify all expected fields exist and can be set
	state := InputState{
		MoveLeft:    true,
		MoveRight:   false,
		Jump:        true,
		JumpPress:   true,
		JumpRelease: true,
		Attack:      false,
		AttackPress: false,
		Dash:        false,
		DashPress:   false,
		UseAbility:  false,
		Pause:       false,
		PausePress:  false,
	}

	if !state.MoveLeft {
		t.Error("MoveLeft field should be settable")
	}
	if !state.Jump {
		t.Error("Jump field should be settable")
	}
	if !state.JumpPress {
		t.Error("JumpPress field should be settable")
	}
	if !state.JumpRelease {
		t.Error("JumpRelease field should be settable")
	}
}

func TestBufferAttack(t *testing.T) {
	ih := NewInputHandler()

	// Buffer an attack
	ih.BufferAttack()

	if ih.buffered.AttackBuffer != BufferFrames {
		t.Errorf("Attack buffer should be %d, got %d", BufferFrames, ih.buffered.AttackBuffer)
	}
}

func TestBufferDash(t *testing.T) {
	ih := NewInputHandler()

	// Buffer a dash
	ih.BufferDash()

	if ih.buffered.DashBuffer != BufferFrames {
		t.Errorf("Dash buffer should be %d, got %d", BufferFrames, ih.buffered.DashBuffer)
	}
}

func TestGetBufferedAttack(t *testing.T) {
	ih := NewInputHandler()

	// No buffered attack initially
	if ih.GetBufferedAttack() {
		t.Error("GetBufferedAttack should return false when no attack is buffered")
	}

	// Buffer an attack
	ih.BufferAttack()

	// Should return true and consume the buffer
	if !ih.GetBufferedAttack() {
		t.Error("GetBufferedAttack should return true when attack is buffered")
	}

	// Buffer should be consumed
	if ih.buffered.AttackBuffer != 0 {
		t.Error("Attack buffer should be consumed after GetBufferedAttack")
	}

	// Second call should return false
	if ih.GetBufferedAttack() {
		t.Error("GetBufferedAttack should return false after buffer is consumed")
	}
}

func TestGetBufferedDash(t *testing.T) {
	ih := NewInputHandler()

	// No buffered dash initially
	if ih.GetBufferedDash() {
		t.Error("GetBufferedDash should return false when no dash is buffered")
	}

	// Buffer a dash
	ih.BufferDash()

	// Should return true and consume the buffer
	if !ih.GetBufferedDash() {
		t.Error("GetBufferedDash should return true when dash is buffered")
	}

	// Buffer should be consumed
	if ih.buffered.DashBuffer != 0 {
		t.Error("Dash buffer should be consumed after GetBufferedDash")
	}

	// Second call should return false
	if ih.GetBufferedDash() {
		t.Error("GetBufferedDash should return false after buffer is consumed")
	}
}

func TestUpdateBuffers(t *testing.T) {
	ih := NewInputHandler()

	// Buffer both actions
	ih.BufferAttack()
	ih.BufferDash()

	// Initial state
	if ih.buffered.AttackBuffer != BufferFrames {
		t.Errorf("Initial attack buffer should be %d", BufferFrames)
	}
	if ih.buffered.DashBuffer != BufferFrames {
		t.Errorf("Initial dash buffer should be %d", BufferFrames)
	}

	// Update buffers once
	ih.UpdateBuffers()

	if ih.buffered.AttackBuffer != BufferFrames-1 {
		t.Errorf("Attack buffer should decrement to %d, got %d", BufferFrames-1, ih.buffered.AttackBuffer)
	}
	if ih.buffered.DashBuffer != BufferFrames-1 {
		t.Errorf("Dash buffer should decrement to %d, got %d", BufferFrames-1, ih.buffered.DashBuffer)
	}

	// Update until buffers expire
	for i := 0; i < BufferFrames; i++ {
		ih.UpdateBuffers()
	}

	if ih.buffered.AttackBuffer != 0 {
		t.Error("Attack buffer should expire to 0")
	}
	if ih.buffered.DashBuffer != 0 {
		t.Error("Dash buffer should expire to 0")
	}

	// Further updates should not go negative
	ih.UpdateBuffers()
	if ih.buffered.AttackBuffer < 0 || ih.buffered.DashBuffer < 0 {
		t.Error("Buffers should not go negative")
	}
}

func TestBufferFramesConstant(t *testing.T) {
	// Verify the buffer window is set to industry standard
	if BufferFrames != 6 {
		t.Errorf("BufferFrames should be 6 (100ms at 60fps), got %d", BufferFrames)
	}
}

func TestKeyMappingFromSettings(t *testing.T) {
	// Create settings with custom key bindings
	controls := &settings.ControlSettings{
		KeyBindings: map[settings.ControlAction]ebiten.Key{
			settings.ActionMoveLeft:  ebiten.KeyA,
			settings.ActionMoveRight: ebiten.KeyD,
			settings.ActionJump:      ebiten.KeySpace,
			settings.ActionAttack:    ebiten.KeyJ,
			settings.ActionDash:      ebiten.KeyK,
			settings.ActionInteract:  ebiten.KeyF,
			settings.ActionPause:     ebiten.KeyEscape,
		},
	}

	mapping := KeyMappingFromSettings(controls)

	// Verify mapping was created
	if mapping == nil {
		t.Fatal("KeyMappingFromSettings returned nil")
	}

	// Verify each action was mapped correctly
	if len(mapping.MoveLeft) != 1 || mapping.MoveLeft[0] != ebiten.KeyA {
		t.Error("MoveLeft should map to KeyA")
	}
	if len(mapping.MoveRight) != 1 || mapping.MoveRight[0] != ebiten.KeyD {
		t.Error("MoveRight should map to KeyD")
	}
	if len(mapping.Jump) != 1 || mapping.Jump[0] != ebiten.KeySpace {
		t.Error("Jump should map to KeySpace")
	}
	if len(mapping.Attack) != 1 || mapping.Attack[0] != ebiten.KeyJ {
		t.Error("Attack should map to KeyJ")
	}
	if len(mapping.Dash) != 1 || mapping.Dash[0] != ebiten.KeyK {
		t.Error("Dash should map to KeyK")
	}
	if len(mapping.UseAbility) != 1 || mapping.UseAbility[0] != ebiten.KeyF {
		t.Error("UseAbility should map to KeyF (from Interact)")
	}
	if len(mapping.Pause) != 1 || mapping.Pause[0] != ebiten.KeyEscape {
		t.Error("Pause should map to KeyEscape")
	}
}

func TestKeyMappingFromSettingsEmptyBindings(t *testing.T) {
	// Create settings with no key bindings
	controls := &settings.ControlSettings{
		KeyBindings: map[settings.ControlAction]ebiten.Key{},
	}

	mapping := KeyMappingFromSettings(controls)

	// Verify mapping was created even with empty bindings
	if mapping == nil {
		t.Fatal("KeyMappingFromSettings should handle empty bindings")
	}

	// All arrays should be empty
	if len(mapping.MoveLeft) != 0 || len(mapping.Jump) != 0 || len(mapping.Attack) != 0 {
		t.Error("Empty bindings should produce empty key arrays")
	}
}

// Note: Full input testing with actual keyboard state requires ebiten context
// which needs X11/graphics libraries. These tests verify the data structures,
// buffering logic, and settings integration are correctly implemented.
// Integration tests with actual keyboard input should be run in a graphical environment.
