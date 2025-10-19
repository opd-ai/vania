package achievement

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewAchievementPersistence tests creating persistence manager
func TestNewAchievementPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	if persistence == nil {
		t.Fatal("Expected persistence manager to be created")
	}
	
	if persistence.GetSaveDir() != tmpDir {
		t.Errorf("Expected save dir %s, got %s", tmpDir, persistence.GetSaveDir())
	}
}

// TestSaveAndLoad tests saving and loading achievements
func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create tracker and unlock some achievements
	tracker := NewAchievementTracker()
	tracker.RecordEnemyKill(false)
	tracker.RecordBossKill(60, false)
	tracker.RecordRoomVisit(false)
	
	// Save achievements
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	err = persistence.Save(tracker)
	if err != nil {
		t.Fatalf("Failed to save achievements: %v", err)
	}
	
	// Verify file was created
	if !persistence.Exists() {
		t.Error("Expected save file to exist")
	}
	
	// Create new tracker and load achievements
	newTracker := NewAchievementTracker()
	err = persistence.Load(newTracker)
	if err != nil {
		t.Fatalf("Failed to load achievements: %v", err)
	}
	
	// Verify unlocked achievements were restored
	if !newTracker.IsUnlocked("first_blood") {
		t.Error("Expected 'first_blood' to be unlocked after load")
	}
	
	if !newTracker.IsUnlocked("boss_hunter") {
		t.Error("Expected 'boss_hunter' to be unlocked after load")
	}
	
	// Verify statistics were restored
	stats := newTracker.GetStatistics()
	if stats.EnemiesDefeated < 1 {
		t.Error("Expected enemies defeated to be restored")
	}
	
	if stats.BossesDefeated < 1 {
		t.Error("Expected bosses defeated to be restored")
	}
}

// TestLoadNonExistent tests loading when no save file exists
func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	tracker := NewAchievementTracker()
	
	// Should not error when file doesn't exist
	err = persistence.Load(tracker)
	if err != nil {
		t.Errorf("Expected no error loading non-existent file, got: %v", err)
	}
	
	// Tracker should be in initial state
	if len(tracker.GetUnlockedAchievements()) != 0 {
		t.Error("Expected no unlocked achievements")
	}
}

// TestDelete tests deleting achievement data
func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create and save achievements
	tracker := NewAchievementTracker()
	tracker.RecordEnemyKill(false)
	
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	err = persistence.Save(tracker)
	if err != nil {
		t.Fatalf("Failed to save achievements: %v", err)
	}
	
	// Verify file exists
	if !persistence.Exists() {
		t.Error("Expected save file to exist")
	}
	
	// Delete achievements
	err = persistence.Delete()
	if err != nil {
		t.Errorf("Failed to delete achievements: %v", err)
	}
	
	// Verify file is gone
	if persistence.Exists() {
		t.Error("Expected save file to be deleted")
	}
}

// TestDeleteNonExistent tests deleting when no file exists
func TestDeleteNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	// Should not error when file doesn't exist
	err = persistence.Delete()
	if err != nil {
		t.Errorf("Expected no error deleting non-existent file, got: %v", err)
	}
}

// TestPersistencePreservesTimestamps tests that unlock timestamps are preserved
func TestPersistencePreservesTimestamps(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create tracker and unlock achievement
	tracker := NewAchievementTracker()
	tracker.RecordEnemyKill(false)
	
	// Get original unlock time
	originalUnlock := tracker.unlocked["first_blood"]
	if originalUnlock == nil {
		t.Fatal("Expected achievement to be unlocked")
	}
	originalTime := originalUnlock.UnlockedAt
	
	// Save and load
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	err = persistence.Save(tracker)
	if err != nil {
		t.Fatalf("Failed to save achievements: %v", err)
	}
	
	newTracker := NewAchievementTracker()
	err = persistence.Load(newTracker)
	if err != nil {
		t.Fatalf("Failed to load achievements: %v", err)
	}
	
	// Verify timestamp was preserved
	newUnlock := newTracker.unlocked["first_blood"]
	if newUnlock == nil {
		t.Fatal("Expected achievement to be unlocked after load")
	}
	
	// Allow for small time differences due to serialization
	timeDiff := newUnlock.UnlockedAt.Sub(originalTime)
	if timeDiff < -time.Second || timeDiff > time.Second {
		t.Errorf("Expected unlock time to be preserved, got difference of %v", timeDiff)
	}
}

// TestPersistencePreservesProgress tests that progress is updated correctly
func TestPersistencePreservesProgress(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create tracker with partial progress
	tracker := NewAchievementTracker()
	for i := 0; i < 25; i++ {
		tracker.RecordEnemyKill(false)
	}
	
	// Save
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	err = persistence.Save(tracker)
	if err != nil {
		t.Fatalf("Failed to save achievements: %v", err)
	}
	
	// Load into new tracker
	newTracker := NewAchievementTracker()
	err = persistence.Load(newTracker)
	if err != nil {
		t.Fatalf("Failed to load achievements: %v", err)
	}
	
	// Verify statistics were restored
	stats := newTracker.GetStatistics()
	if stats.EnemiesDefeated != 25 {
		t.Errorf("Expected 25 enemies defeated, got %d", stats.EnemiesDefeated)
	}
	
	// Continue progress
	for i := 0; i < 25; i++ {
		newTracker.RecordEnemyKill(false)
	}
	
	// Should unlock "slayer" now
	if !newTracker.IsUnlocked("slayer") {
		t.Error("Expected 'slayer' to be unlocked after continuing progress")
	}
}

// TestMultipleSaves tests multiple save operations
func TestMultipleSaves(t *testing.T) {
	tmpDir := t.TempDir()
	
	persistence, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	tracker := NewAchievementTracker()
	
	// First save
	tracker.RecordEnemyKill(false)
	err = persistence.Save(tracker)
	if err != nil {
		t.Fatalf("Failed first save: %v", err)
	}
	
	// Second save with more progress
	tracker.RecordBossKill(60, false)
	err = persistence.Save(tracker)
	if err != nil {
		t.Fatalf("Failed second save: %v", err)
	}
	
	// Load and verify both achievements are unlocked
	newTracker := NewAchievementTracker()
	err = persistence.Load(newTracker)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}
	
	if !newTracker.IsUnlocked("first_blood") {
		t.Error("Expected 'first_blood' to be unlocked")
	}
	
	if !newTracker.IsUnlocked("boss_hunter") {
		t.Error("Expected 'boss_hunter' to be unlocked")
	}
}

// TestPersistenceCreatesSaveDirectory tests that save directory is created
func TestPersistenceCreatesSaveDirectory(t *testing.T) {
	tmpDir := filepath.Join(t.TempDir(), "nested", "directory")
	
	// Directory should not exist yet
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Fatal("Expected directory to not exist")
	}
	
	// Create persistence manager
	_, err := NewAchievementPersistence(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create persistence manager: %v", err)
	}
	
	// Directory should now exist
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}
}
