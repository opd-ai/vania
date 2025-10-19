package save

import (
	"time"
)

// CheckpointManager handles automatic checkpoint saves
type CheckpointManager struct {
	saveManager     *SaveManager
	lastCheckpoint  time.Time
	checkpointInterval time.Duration
	autoSaveEnabled bool
}

// NewCheckpointManager creates a new checkpoint manager
func NewCheckpointManager(saveManager *SaveManager) *CheckpointManager {
	return &CheckpointManager{
		saveManager:        saveManager,
		lastCheckpoint:     time.Now(),
		checkpointInterval: 5 * time.Minute, // Auto-save every 5 minutes
		autoSaveEnabled:    true,
	}
}

// SetCheckpointInterval sets how often auto-saves occur
func (cm *CheckpointManager) SetCheckpointInterval(interval time.Duration) {
	cm.checkpointInterval = interval
}

// EnableAutoSave enables or disables auto-save
func (cm *CheckpointManager) EnableAutoSave(enabled bool) {
	cm.autoSaveEnabled = enabled
}

// ShouldCheckpoint returns true if it's time for an auto-save
func (cm *CheckpointManager) ShouldCheckpoint() bool {
	if !cm.autoSaveEnabled {
		return false
	}
	return time.Since(cm.lastCheckpoint) >= cm.checkpointInterval
}

// CreateCheckpoint creates a checkpoint save
func (cm *CheckpointManager) CreateCheckpoint(data *SaveData) error {
	if !cm.autoSaveEnabled {
		return nil
	}
	
	err := cm.saveManager.AutoSave(data)
	if err == nil {
		cm.lastCheckpoint = time.Now()
	}
	return err
}

// GetTimeSinceLastCheckpoint returns duration since last checkpoint
func (cm *CheckpointManager) GetTimeSinceLastCheckpoint() time.Duration {
	return time.Since(cm.lastCheckpoint)
}

// ResetCheckpointTimer resets the checkpoint timer
func (cm *CheckpointManager) ResetCheckpointTimer() {
	cm.lastCheckpoint = time.Now()
}
