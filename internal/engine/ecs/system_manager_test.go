package ecs

import (
	"errors"
	"testing"
)

func TestNewSystemManager(t *testing.T) {
	sm := NewSystemManager()

	if sm == nil {
		t.Fatal("Expected non-nil SystemManager")
	}

	if sm.GetSystemCount() != 0 {
		t.Errorf("Expected 0 systems, got %d", sm.GetSystemCount())
	}

	if !sm.sorted {
		t.Error("Expected new SystemManager to be sorted")
	}
}

func TestSystemManagerRegister(t *testing.T) {
	sm := NewSystemManager()
	system := &MockSystem{}

	sm.Register(system, 10)

	if sm.GetSystemCount() != 1 {
		t.Errorf("Expected 1 system, got %d", sm.GetSystemCount())
	}

	if sm.sorted {
		t.Error("Expected SystemManager to be unsorted after registration")
	}
}

func TestSystemManagerUnregister(t *testing.T) {
	sm := NewSystemManager()
	system1 := &MockSystem{}
	system2 := &MockSystem{}

	sm.Register(system1, 10)
	sm.Register(system2, 20)

	if sm.GetSystemCount() != 2 {
		t.Errorf("Expected 2 systems, got %d", sm.GetSystemCount())
	}

	removed := sm.Unregister(system1)
	if !removed {
		t.Error("Expected system to be removed")
	}

	if sm.GetSystemCount() != 1 {
		t.Errorf("Expected 1 system after removal, got %d", sm.GetSystemCount())
	}
}

func TestSystemManagerUnregisterNonExistent(t *testing.T) {
	sm := NewSystemManager()
	system := &MockSystem{}

	removed := sm.Unregister(system)
	if removed {
		t.Error("Expected false when removing non-existent system")
	}
}

func TestSystemManagerPriorityOrdering(t *testing.T) {
	sm := NewSystemManager()

	// Register systems in reverse priority order
	system1 := &MockSystem{}
	system2 := &MockSystem{}
	system3 := &MockSystem{}

	sm.Register(system3, 30) // Should execute last
	sm.Register(system1, 10) // Should execute first
	sm.Register(system2, 20) // Should execute second

	// Call Update to trigger sorting
	err := sm.Update(0.016)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all systems were updated
	if !system1.updateCalled || !system2.updateCalled || !system3.updateCalled {
		t.Error("Expected all systems to be updated")
	}
}

func TestSystemManagerUpdate(t *testing.T) {
	sm := NewSystemManager()
	system := &MockSystem{}

	sm.Register(system, 10)

	err := sm.Update(0.016)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !system.updateCalled {
		t.Error("Expected system Update to be called")
	}
}

func TestSystemManagerUpdateError(t *testing.T) {
	sm := NewSystemManager()
	expectedError := errors.New("test error")
	system := &MockSystem{updateError: expectedError}

	sm.Register(system, 10)

	err := sm.Update(0.016)
	if err == nil {
		t.Error("Expected error from Update")
	}
}

func TestSystemManagerDraw(t *testing.T) {
	sm := NewSystemManager()
	system := &MockSystem{}

	sm.Register(system, 10)

	sm.Draw(nil)

	if !system.drawCalled {
		t.Error("Expected system Draw to be called")
	}
}

func TestSystemManagerSetGenre(t *testing.T) {
	sm := NewSystemManager()
	system1 := &MockSystem{}
	system2 := &MockSystem{}

	sm.Register(system1, 10)
	sm.Register(system2, 20)

	sm.SetGenre("scifi")

	if system1.currentGenre != "scifi" {
		t.Errorf("Expected system1 genre 'scifi', got %s", system1.currentGenre)
	}

	if system2.currentGenre != "scifi" {
		t.Errorf("Expected system2 genre 'scifi', got %s", system2.currentGenre)
	}
}

func TestSystemManagerClear(t *testing.T) {
	sm := NewSystemManager()
	system1 := &MockSystem{}
	system2 := &MockSystem{}

	sm.Register(system1, 10)
	sm.Register(system2, 20)

	if sm.GetSystemCount() != 2 {
		t.Errorf("Expected 2 systems, got %d", sm.GetSystemCount())
	}

	sm.Clear()

	if sm.GetSystemCount() != 0 {
		t.Errorf("Expected 0 systems after clear, got %d", sm.GetSystemCount())
	}

	if !sm.sorted {
		t.Error("Expected SystemManager to be sorted after clear")
	}
}

func TestSystemManagerMultipleUpdates(t *testing.T) {
	sm := NewSystemManager()
	system := &MockSystem{}

	sm.Register(system, 10)

	// First update
	err := sm.Update(0.016)
	if err != nil {
		t.Errorf("Expected no error on first update, got %v", err)
	}

	// Reset flag
	system.updateCalled = false

	// Second update
	err = sm.Update(0.016)
	if err != nil {
		t.Errorf("Expected no error on second update, got %v", err)
	}

	if !system.updateCalled {
		t.Error("Expected system Update to be called on second update")
	}
}

func TestSystemManagerDeterministicOrdering(t *testing.T) {
	sm := NewSystemManager()

	// Create systems with same priority
	system1 := &MockSystem{}
	system2 := &MockSystem{}
	system3 := &MockSystem{}

	sm.Register(system1, 10)
	sm.Register(system2, 10)
	sm.Register(system3, 10)

	// Multiple updates should maintain order
	for i := 0; i < 5; i++ {
		err := sm.Update(0.016)
		if err != nil {
			t.Errorf("Update %d failed: %v", i, err)
		}
	}

	// All should have been called
	if !system1.updateCalled || !system2.updateCalled || !system3.updateCalled {
		t.Error("Expected all systems to be updated")
	}
}
