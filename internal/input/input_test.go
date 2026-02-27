package input

import (
	"testing"
)

func TestNewInputHandler(t *testing.T) {
	ih := NewInputHandler()

	if ih == nil {
		t.Fatal("NewInputHandler returned nil")
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

// Note: Full input testing requires ebiten context which needs X11/graphics libraries.
// These tests verify the data structures are correctly defined.
// Integration tests with actual keyboard input should be run in a graphical environment.
