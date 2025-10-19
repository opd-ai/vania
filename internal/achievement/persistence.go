package achievement

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AchievementSaveData represents achievement data for persistence
type AchievementSaveData struct {
	Version    string                         `json:"version"`
	SaveTime   time.Time                      `json:"save_time"`
	Unlocked   map[AchievementID]time.Time    `json:"unlocked"`
	Statistics Statistics                     `json:"statistics"`
}

// AchievementPersistence manages achievement save/load operations
type AchievementPersistence struct {
	saveDir string
}

const (
	achievementVersion = "1.0.0"
	achievementFile    = "achievements.json"
)

// NewAchievementPersistence creates a new achievement persistence manager
func NewAchievementPersistence(saveDir string) (*AchievementPersistence, error) {
	// Default save directory in user's home config
	if saveDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		saveDir = filepath.Join(homeDir, ".vania", "achievements")
	}
	
	// Create save directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}
	
	return &AchievementPersistence{
		saveDir: saveDir,
	}, nil
}

// Save saves achievement data to disk
func (ap *AchievementPersistence) Save(tracker *AchievementTracker) error {
	// Prepare save data
	unlockedMap := make(map[AchievementID]time.Time)
	for id, unlocked := range tracker.unlocked {
		unlockedMap[id] = unlocked.UnlockedAt
	}
	
	saveData := AchievementSaveData{
		Version:    achievementVersion,
		SaveTime:   time.Now(),
		Unlocked:   unlockedMap,
		Statistics: tracker.stats,
	}
	
	// Convert to JSON
	jsonData, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal achievement data: %w", err)
	}
	
	// Write to file
	filename := filepath.Join(ap.saveDir, achievementFile)
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write achievement file: %w", err)
	}
	
	return nil
}

// Load loads achievement data from disk
func (ap *AchievementPersistence) Load(tracker *AchievementTracker) error {
	filename := filepath.Join(ap.saveDir, achievementFile)
	
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // No saved data, start fresh
	}
	
	// Read file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read achievement file: %w", err)
	}
	
	// Parse JSON
	var saveData AchievementSaveData
	if err := json.Unmarshal(jsonData, &saveData); err != nil {
		return fmt.Errorf("failed to parse achievement file: %w", err)
	}
	
	// Validate version (for future compatibility)
	if saveData.Version != achievementVersion {
		// Could implement migration logic here
		return fmt.Errorf("incompatible achievement version: %s (expected %s)", saveData.Version, achievementVersion)
	}
	
	// Restore unlocked achievements
	for id, unlockedAt := range saveData.Unlocked {
		if tracker.achievements[id] != nil {
			tracker.unlocked[id] = &UnlockedAchievement{
				AchievementID: id,
				UnlockedAt:    unlockedAt,
				Progress:      1.0,
			}
			
			// Update progress to complete
			if progress := tracker.progress[id]; progress != nil {
				progress.Progress = 1.0
				progress.CurrentValue = progress.TargetValue
			}
		}
	}
	
	// Restore statistics
	tracker.stats = saveData.Statistics
	
	return nil
}

// Delete removes saved achievement data
func (ap *AchievementPersistence) Delete() error {
	filename := filepath.Join(ap.saveDir, achievementFile)
	
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // Already deleted
	}
	
	// Delete file
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete achievement file: %w", err)
	}
	
	return nil
}

// Exists checks if saved achievement data exists
func (ap *AchievementPersistence) Exists() bool {
	filename := filepath.Join(ap.saveDir, achievementFile)
	_, err := os.Stat(filename)
	return err == nil
}

// GetSaveDir returns the save directory path
func (ap *AchievementPersistence) GetSaveDir() string {
	return ap.saveDir
}
