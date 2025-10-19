package save

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSaveManager(t *testing.T) {
	// Create temp directory for testing
	tempDir := t.TempDir()
	
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	if sm.saveDir != tempDir {
		t.Errorf("Expected saveDir %s, got %s", tempDir, sm.saveDir)
	}
	
	if sm.currentSlot != 1 {
		t.Errorf("Expected currentSlot 1, got %d", sm.currentSlot)
	}
	
	// Verify directory was created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Save directory was not created")
	}
}

func TestNewSaveManagerWithEmptyDir(t *testing.T) {
	// Test with empty string (should use default home directory)
	sm, err := NewSaveManager("")
	if err != nil {
		t.Fatalf("Failed to create SaveManager with default dir: %v", err)
	}
	
	// Should have created a directory
	if sm.saveDir == "" {
		t.Error("SaveManager did not set a save directory")
	}
	
	// Clean up
	os.RemoveAll(filepath.Join(sm.saveDir, ".."))
}

func TestSaveAndLoadGame(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Create test save data
	originalData := &SaveData{
		Seed:           42,
		PlayTime:       1800,
		PlayerX:        100.5,
		PlayerY:        200.3,
		PlayerHealth:   75,
		PlayerMaxHealth: 100,
		PlayerAbilities: map[string]bool{
			"double_jump": true,
			"dash":        true,
		},
		CurrentRoomID:   5,
		VisitedRooms:    []int{0, 1, 2, 3, 5},
		DefeatedEnemies: map[int]bool{1: true, 2: true},
		CollectedItems:  map[int]bool{10: true, 11: true},
		UnlockedDoors:   map[string]bool{"door_1": true},
		BossesDefeated:  []int{1},
		CheckpointID:    3,
	}
	
	// Save the game
	slotID := 1
	if err := sm.SaveGame(originalData, slotID); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}
	
	// Load the game
	loadedData, err := sm.LoadGame(slotID)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}
	
	// Verify data
	if loadedData.Seed != originalData.Seed {
		t.Errorf("Expected seed %d, got %d", originalData.Seed, loadedData.Seed)
	}
	if loadedData.PlayTime != originalData.PlayTime {
		t.Errorf("Expected play time %d, got %d", originalData.PlayTime, loadedData.PlayTime)
	}
	if loadedData.PlayerX != originalData.PlayerX {
		t.Errorf("Expected player X %f, got %f", originalData.PlayerX, loadedData.PlayerX)
	}
	if loadedData.PlayerHealth != originalData.PlayerHealth {
		t.Errorf("Expected player health %d, got %d", originalData.PlayerHealth, loadedData.PlayerHealth)
	}
	if len(loadedData.VisitedRooms) != len(originalData.VisitedRooms) {
		t.Errorf("Expected %d visited rooms, got %d", len(originalData.VisitedRooms), len(loadedData.VisitedRooms))
	}
	if !loadedData.PlayerAbilities["double_jump"] {
		t.Error("Expected double_jump ability to be unlocked")
	}
}

func TestSaveGameInvalidSlot(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	data := &SaveData{Seed: 42}
	
	// Test invalid slot IDs
	invalidSlots := []int{-1, 5, 10, 100}
	for _, slotID := range invalidSlots {
		err := sm.SaveGame(data, slotID)
		if err == nil {
			t.Errorf("Expected error for invalid slot %d, got nil", slotID)
		}
	}
}

func TestLoadGameNonexistent(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Try to load from empty slot
	_, err = sm.LoadGame(1)
	if err == nil {
		t.Error("Expected error when loading nonexistent save")
	}
}

func TestAutoSave(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	data := &SaveData{
		Seed:         123,
		PlayerHealth: 80,
	}
	
	// Auto-save
	if err := sm.AutoSave(data); err != nil {
		t.Fatalf("Failed to auto-save: %v", err)
	}
	
	// Load auto-save
	loadedData, err := sm.LoadGame(autoSaveID)
	if err != nil {
		t.Fatalf("Failed to load auto-save: %v", err)
	}
	
	if loadedData.Seed != data.Seed {
		t.Errorf("Expected seed %d, got %d", data.Seed, loadedData.Seed)
	}
}

func TestDeleteSave(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Save a game
	data := &SaveData{Seed: 42}
	slotID := 2
	if err := sm.SaveGame(data, slotID); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}
	
	// Verify it exists
	_, err = sm.LoadGame(slotID)
	if err != nil {
		t.Fatal("Save should exist before deletion")
	}
	
	// Delete it
	if err := sm.DeleteSave(slotID); err != nil {
		t.Fatalf("Failed to delete save: %v", err)
	}
	
	// Verify it's gone
	_, err = sm.LoadGame(slotID)
	if err == nil {
		t.Error("Save should not exist after deletion")
	}
}

func TestDeleteNonexistentSave(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Delete nonexistent save (should not error)
	if err := sm.DeleteSave(3); err != nil {
		t.Errorf("Unexpected error deleting nonexistent save: %v", err)
	}
}

func TestListSaves(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Initially all empty
	saves, err := sm.ListSaves()
	if err != nil {
		t.Fatalf("Failed to list saves: %v", err)
	}
	
	if len(saves) != maxSlots {
		t.Errorf("Expected %d slots, got %d", maxSlots, len(saves))
	}
	
	emptyCount := 0
	for _, save := range saves {
		if save.IsEmpty {
			emptyCount++
		}
	}
	if emptyCount != maxSlots {
		t.Errorf("Expected all %d slots empty, got %d", maxSlots, emptyCount)
	}
	
	// Save to some slots
	data1 := &SaveData{Seed: 100, PlayerHealth: 50}
	data2 := &SaveData{Seed: 200, PlayerHealth: 75}
	
	sm.SaveGame(data1, 1)
	sm.SaveGame(data2, 3)
	
	// List again
	saves, err = sm.ListSaves()
	if err != nil {
		t.Fatalf("Failed to list saves: %v", err)
	}
	
	existCount := 0
	for _, save := range saves {
		if save.Exists && !save.IsEmpty {
			existCount++
		}
	}
	if existCount != 2 {
		t.Errorf("Expected 2 existing saves, got %d", existCount)
	}
}

func TestGetSaveInfo(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Save a game
	data := &SaveData{
		Seed:         999,
		PlayTime:     3600,
		PlayerHealth: 90,
		CurrentRoomID: 7,
	}
	slotID := 2
	if err := sm.SaveGame(data, slotID); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}
	
	// Get info
	info, err := sm.GetSaveInfo(slotID)
	if err != nil {
		t.Fatalf("Failed to get save info: %v", err)
	}
	
	if !info.Exists {
		t.Error("Save should exist")
	}
	if info.IsEmpty {
		t.Error("Save should not be empty")
	}
	if info.Corrupt {
		t.Error("Save should not be corrupt")
	}
	if info.Seed != data.Seed {
		t.Errorf("Expected seed %d, got %d", data.Seed, info.Seed)
	}
	if info.PlayerHealth != data.PlayerHealth {
		t.Errorf("Expected health %d, got %d", data.PlayerHealth, info.PlayerHealth)
	}
	if info.FileSize == 0 {
		t.Error("File size should be greater than 0")
	}
}

func TestSaveMetadata(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	data := &SaveData{Seed: 42}
	slotID := 1
	
	beforeSave := time.Now()
	if err := sm.SaveGame(data, slotID); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}
	afterSave := time.Now()
	
	loadedData, err := sm.LoadGame(slotID)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}
	
	// Check version was set
	if loadedData.Version != saveVersion {
		t.Errorf("Expected version %s, got %s", saveVersion, loadedData.Version)
	}
	
	// Check save time was set
	if loadedData.SaveTime.Before(beforeSave) || loadedData.SaveTime.After(afterSave) {
		t.Error("Save time not within expected range")
	}
	
	// Check slot ID was set
	if loadedData.SlotID != slotID {
		t.Errorf("Expected slot ID %d, got %d", slotID, loadedData.SlotID)
	}
}

func TestGetSlotFilename(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	// Test auto-save filename
	autoFilename := sm.getSlotFilename(autoSaveID)
	expectedAuto := filepath.Join(tempDir, "autosave.json")
	if autoFilename != expectedAuto {
		t.Errorf("Expected auto-save filename %s, got %s", expectedAuto, autoFilename)
	}
	
	// Test regular slot filenames
	for i := 1; i < 5; i++ {
		filename := sm.getSlotFilename(i)
		if filepath.Base(filename) != "save_"+string(rune('0'+i))+".json" {
			t.Errorf("Unexpected filename format for slot %d: %s", i, filename)
		}
	}
}

func TestGetters(t *testing.T) {
	tempDir := t.TempDir()
	sm, err := NewSaveManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create SaveManager: %v", err)
	}
	
	if sm.GetSaveDir() != tempDir {
		t.Errorf("Expected save dir %s, got %s", tempDir, sm.GetSaveDir())
	}
	
	if sm.GetCurrentSlot() != 1 {
		t.Errorf("Expected current slot 1, got %d", sm.GetCurrentSlot())
	}
	
	// Change slot by saving
	data := &SaveData{Seed: 42}
	sm.SaveGame(data, 3)
	
	if sm.GetCurrentSlot() != 3 {
		t.Errorf("Expected current slot 3, got %d", sm.GetCurrentSlot())
	}
}
