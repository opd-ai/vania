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
}

// Note: Full input testing requires ebiten context which needs X11/graphics libraries.
// These tests verify the data structures are correctly defined.
// Integration tests with actual keyboard input should be run in a graphical environment.
