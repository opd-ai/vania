package save

import (
	"testing"
	"time"
)

func TestNewCheckpointManager(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	cm := NewCheckpointManager(sm)
	
	if cm.saveManager != sm {
		t.Error("SaveManager not set correctly")
	}
	if !cm.autoSaveEnabled {
		t.Error("Auto-save should be enabled by default")
	}
	if cm.checkpointInterval != 5*time.Minute {
		t.Errorf("Expected 5 minute interval, got %v", cm.checkpointInterval)
	}
}

func TestSetCheckpointInterval(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	newInterval := 10 * time.Second
	cm.SetCheckpointInterval(newInterval)
	
	if cm.checkpointInterval != newInterval {
		t.Errorf("Expected interval %v, got %v", newInterval, cm.checkpointInterval)
	}
}

func TestEnableAutoSave(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	// Initially enabled
	if !cm.autoSaveEnabled {
		t.Error("Auto-save should be enabled initially")
	}
	
	// Disable
	cm.EnableAutoSave(false)
	if cm.autoSaveEnabled {
		t.Error("Auto-save should be disabled")
	}
	
	// Re-enable
	cm.EnableAutoSave(true)
	if !cm.autoSaveEnabled {
		t.Error("Auto-save should be enabled")
	}
}

func TestShouldCheckpoint(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	// Set short interval for testing
	cm.SetCheckpointInterval(100 * time.Millisecond)
	
	// Should not need checkpoint immediately
	if cm.ShouldCheckpoint() {
		t.Error("Should not need checkpoint immediately after creation")
	}
	
	// Wait for interval
	time.Sleep(150 * time.Millisecond)
	
	// Should need checkpoint now
	if !cm.ShouldCheckpoint() {
		t.Error("Should need checkpoint after interval")
	}
	
	// Disable auto-save
	cm.EnableAutoSave(false)
	if cm.ShouldCheckpoint() {
		t.Error("Should not checkpoint when auto-save is disabled")
	}
}

func TestCreateCheckpoint(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	data := &SaveData{
		Seed:         42,
		PlayerHealth: 100,
	}
	
	// Create checkpoint
	beforeCheckpoint := time.Now()
	err := cm.CreateCheckpoint(data)
	if err != nil {
		t.Fatalf("Failed to create checkpoint: %v", err)
	}
	
	// Verify checkpoint was saved
	loadedData, err := sm.LoadGame(autoSaveID)
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}
	
	if loadedData.Seed != data.Seed {
		t.Errorf("Expected seed %d, got %d", data.Seed, loadedData.Seed)
	}
	
	// Verify checkpoint time was updated
	if cm.lastCheckpoint.Before(beforeCheckpoint) {
		t.Error("Checkpoint time not updated")
	}
}

func TestCreateCheckpointDisabled(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	// Disable auto-save
	cm.EnableAutoSave(false)
	
	data := &SaveData{Seed: 42}
	
	// Try to create checkpoint (should not error but do nothing)
	err := cm.CreateCheckpoint(data)
	if err != nil {
		t.Errorf("Unexpected error when auto-save disabled: %v", err)
	}
	
	// Verify no checkpoint was saved
	_, err = sm.LoadGame(autoSaveID)
	if err == nil {
		t.Error("Checkpoint should not have been saved when auto-save disabled")
	}
}

func TestGetTimeSinceLastCheckpoint(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	// Initially should be near zero
	duration := cm.GetTimeSinceLastCheckpoint()
	if duration > 100*time.Millisecond {
		t.Error("Time since checkpoint should be minimal initially")
	}
	
	// Wait a bit
	time.Sleep(200 * time.Millisecond)
	
	// Should have increased
	duration = cm.GetTimeSinceLastCheckpoint()
	if duration < 150*time.Millisecond {
		t.Error("Time since checkpoint should have increased")
	}
}

func TestResetCheckpointTimer(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	// Wait a bit
	time.Sleep(100 * time.Millisecond)
	
	// Verify time has passed
	if cm.GetTimeSinceLastCheckpoint() < 50*time.Millisecond {
		t.Error("Expected time to have passed")
	}
	
	// Reset timer
	cm.ResetCheckpointTimer()
	
	// Time should be near zero again
	if cm.GetTimeSinceLastCheckpoint() > 50*time.Millisecond {
		t.Error("Timer should have been reset")
	}
}

func TestCheckpointWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	sm, _ := NewSaveManager(tempDir)
	cm := NewCheckpointManager(sm)
	
	// Set very short interval for testing
	cm.SetCheckpointInterval(50 * time.Millisecond)
	
	data := &SaveData{
		Seed:         123,
		PlayerHealth: 75,
	}
	
	// Initially no checkpoint needed
	if cm.ShouldCheckpoint() {
		t.Error("Should not need checkpoint initially")
	}
	
	// Wait for interval
	time.Sleep(100 * time.Millisecond)
	
	// Should need checkpoint
	if !cm.ShouldCheckpoint() {
		t.Error("Should need checkpoint after interval")
	}
	
	// Create checkpoint
	if err := cm.CreateCheckpoint(data); err != nil {
		t.Fatalf("Failed to create checkpoint: %v", err)
	}
	
	// Should not need checkpoint immediately after
	if cm.ShouldCheckpoint() {
		t.Error("Should not need checkpoint immediately after creating one")
	}
	
	// Verify checkpoint was saved
	loadedData, err := sm.LoadGame(autoSaveID)
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}
	if loadedData.Seed != data.Seed {
		t.Errorf("Expected seed %d, got %d", data.Seed, loadedData.Seed)
	}
}
