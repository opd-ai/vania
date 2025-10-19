// Package save provides game state persistence functionality, allowing
// players to save and load their progress across multiple save slots with
// automatic checkpoints and manual save points.
package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SaveData represents a complete game save
type SaveData struct {
	// Metadata
	Version     string    `json:"version"`
	Seed        int64     `json:"seed"`
	SaveTime    time.Time `json:"save_time"`
	PlayTime    int64     `json:"play_time_seconds"`
	SlotID      int       `json:"slot_id"`
	
	// Player state
	PlayerX        float64           `json:"player_x"`
	PlayerY        float64           `json:"player_y"`
	PlayerHealth   int               `json:"player_health"`
	PlayerMaxHealth int              `json:"player_max_health"`
	PlayerAbilities map[string]bool  `json:"player_abilities"`
	
	// World state
	CurrentRoomID  int              `json:"current_room_id"`
	VisitedRooms   []int            `json:"visited_rooms"`
	DefeatedEnemies map[int]bool    `json:"defeated_enemies"`
	CollectedItems  map[int]bool    `json:"collected_items"`
	UnlockedDoors   map[string]bool `json:"unlocked_doors"`
	
	// Progress tracking
	BossesDefeated []int            `json:"bosses_defeated"`
	CheckpointID   int              `json:"checkpoint_id"`
}

// SaveManager handles all save/load operations
type SaveManager struct {
	saveDir      string
	currentSlot  int
	autoSaveSlot int
}

const (
	saveVersion = "1.0.0"
	maxSlots    = 5
	autoSaveID  = 0 // Slot 0 is reserved for auto-save
)

// NewSaveManager creates a new save manager
func NewSaveManager(saveDir string) (*SaveManager, error) {
	// Default save directory in user's home config
	if saveDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		saveDir = filepath.Join(homeDir, ".vania", "saves")
	}
	
	// Create save directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}
	
	return &SaveManager{
		saveDir:      saveDir,
		currentSlot:  1,
		autoSaveSlot: autoSaveID,
	}, nil
}

// SaveGame saves the game state to a specific slot
func (sm *SaveManager) SaveGame(data *SaveData, slotID int) error {
	if slotID < 0 || slotID >= maxSlots {
		return fmt.Errorf("invalid slot ID: %d (must be 0-%d)", slotID, maxSlots-1)
	}
	
	// Set metadata
	data.Version = saveVersion
	data.SaveTime = time.Now()
	data.SlotID = slotID
	
	// Convert to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}
	
	// Write to file
	filename := sm.getSlotFilename(slotID)
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}
	
	sm.currentSlot = slotID
	return nil
}

// LoadGame loads game state from a specific slot
func (sm *SaveManager) LoadGame(slotID int) (*SaveData, error) {
	if slotID < 0 || slotID >= maxSlots {
		return nil, fmt.Errorf("invalid slot ID: %d (must be 0-%d)", slotID, maxSlots-1)
	}
	
	filename := sm.getSlotFilename(slotID)
	
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("save file does not exist")
	}
	
	// Read file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}
	
	// Parse JSON
	var data SaveData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse save file: %w", err)
	}
	
	// Validate version
	if data.Version != saveVersion {
		return nil, fmt.Errorf("incompatible save version: %s (expected %s)", data.Version, saveVersion)
	}
	
	sm.currentSlot = slotID
	return &data, nil
}

// AutoSave saves to the auto-save slot
func (sm *SaveManager) AutoSave(data *SaveData) error {
	return sm.SaveGame(data, sm.autoSaveSlot)
}

// DeleteSave removes a save file
func (sm *SaveManager) DeleteSave(slotID int) error {
	if slotID < 0 || slotID >= maxSlots {
		return fmt.Errorf("invalid slot ID: %d (must be 0-%d)", slotID, maxSlots-1)
	}
	
	filename := sm.getSlotFilename(slotID)
	
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // Already deleted
	}
	
	// Delete file
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete save file: %w", err)
	}
	
	return nil
}

// ListSaves returns information about all save slots
func (sm *SaveManager) ListSaves() ([]SaveInfo, error) {
	saves := make([]SaveInfo, 0, maxSlots)
	
	for i := 0; i < maxSlots; i++ {
		info, err := sm.GetSaveInfo(i)
		if err != nil {
			// Slot is empty or corrupt
			saves = append(saves, SaveInfo{
				SlotID:  i,
				Exists:  false,
				IsEmpty: true,
			})
		} else {
			saves = append(saves, info)
		}
	}
	
	return saves, nil
}

// GetSaveInfo returns information about a specific save slot
func (sm *SaveManager) GetSaveInfo(slotID int) (SaveInfo, error) {
	if slotID < 0 || slotID >= maxSlots {
		return SaveInfo{}, fmt.Errorf("invalid slot ID: %d", slotID)
	}
	
	filename := sm.getSlotFilename(slotID)
	
	// Check if file exists
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return SaveInfo{
			SlotID:  slotID,
			Exists:  false,
			IsEmpty: true,
		}, fmt.Errorf("save does not exist")
	}
	
	// Read minimal data for info
	data, err := sm.LoadGame(slotID)
	if err != nil {
		return SaveInfo{
			SlotID:  slotID,
			Exists:  true,
			IsEmpty: false,
			Corrupt: true,
		}, err
	}
	
	return SaveInfo{
		SlotID:      slotID,
		Exists:      true,
		IsEmpty:     false,
		Corrupt:     false,
		Seed:        data.Seed,
		SaveTime:    data.SaveTime,
		PlayTime:    data.PlayTime,
		PlayerHealth: data.PlayerHealth,
		RoomID:      data.CurrentRoomID,
		FileSize:    stat.Size(),
	}, nil
}

// SaveInfo contains metadata about a save slot
type SaveInfo struct {
	SlotID      int
	Exists      bool
	IsEmpty     bool
	Corrupt     bool
	Seed        int64
	SaveTime    time.Time
	PlayTime    int64
	PlayerHealth int
	RoomID      int
	FileSize    int64
}

// getSlotFilename returns the filename for a given slot
func (sm *SaveManager) getSlotFilename(slotID int) string {
	if slotID == sm.autoSaveSlot {
		return filepath.Join(sm.saveDir, "autosave.json")
	}
	return filepath.Join(sm.saveDir, fmt.Sprintf("save_%d.json", slotID))
}

// GetSaveDir returns the save directory path
func (sm *SaveManager) GetSaveDir() string {
	return sm.saveDir
}

// GetCurrentSlot returns the currently active slot
func (sm *SaveManager) GetCurrentSlot() int {
	return sm.currentSlot
}
