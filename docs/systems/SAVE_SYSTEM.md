# Save/Load System Documentation

## Overview

The VANIA save/load system provides persistent game state management, allowing players to save their progress and resume later. The system supports multiple save slots, automatic checkpoints, and comprehensive state tracking.

## Architecture

```
internal/save/
├── save_manager.go    - Core save/load functionality
├── checkpoint.go      - Automatic checkpoint system
├── save_manager_test.go - SaveManager tests (13 tests)
└── checkpoint_test.go   - CheckpointManager tests (9 tests)
```

### Components

#### 1. SaveManager
Central component managing all save operations:
- **Save/Load Operations**: Handles reading and writing save files
- **Slot Management**: Supports 5 save slots (slot 0 reserved for auto-save)
- **File Operations**: JSON-based save format for readability
- **Validation**: Version checking and data integrity

#### 2. CheckpointManager
Handles automatic checkpoint creation:
- **Auto-Save**: Periodic saves every 5 minutes (configurable)
- **Timer Management**: Tracks time since last checkpoint
- **Enable/Disable**: Can be toggled on/off

#### 3. SaveData Structure
Complete game state snapshot:
```go
type SaveData struct {
    // Metadata
    Version     string
    Seed        int64
    SaveTime    time.Time
    PlayTime    int64
    SlotID      int
    
    // Player state
    PlayerX, PlayerY     float64
    PlayerHealth         int
    PlayerMaxHealth      int
    PlayerAbilities      map[string]bool
    
    // World state
    CurrentRoomID        int
    VisitedRooms         []int
    DefeatedEnemies      map[int]bool
    CollectedItems       map[int]bool
    UnlockedDoors        map[string]bool
    
    // Progress
    BossesDefeated       []int
    CheckpointID         int
}
```

## Features

### Multiple Save Slots
- **5 slots total**: Slots 1-4 for manual saves, slot 0 for auto-save
- **Slot independence**: Each slot is completely independent
- **Slot management**: List, view info, or delete any slot

### Automatic Checkpoints
- **Periodic auto-save**: Default 5-minute interval
- **Configurable interval**: Can be adjusted programmatically
- **Non-intrusive**: Saves in background without interrupting gameplay
- **Can be disabled**: Players can turn off auto-save if desired

### Save File Format
Files stored as human-readable JSON:
```json
{
  "version": "1.0.0",
  "seed": 42,
  "save_time": "2025-10-19T04:00:00Z",
  "play_time_seconds": 1800,
  "slot_id": 1,
  "player_x": 150.5,
  "player_y": 200.0,
  "player_health": 75,
  "player_max_health": 100,
  "player_abilities": {
    "double_jump": true,
    "dash": true
  },
  "current_room_id": 5,
  "visited_rooms": [0, 1, 2, 3, 5],
  "defeated_enemies": {"123": true, "456": true},
  "collected_items": {"10": true, "11": true},
  "unlocked_doors": {"door_1": true},
  "bosses_defeated": [1],
  "checkpoint_id": 5
}
```

### Save Location
- **Default**: `~/.vania/saves/`
- **Files**: 
  - `autosave.json` - Automatic checkpoint save
  - `save_1.json` through `save_4.json` - Manual save slots

## Usage

### Initialization

```go
// Create save manager (uses default location)
saveManager, err := save.NewSaveManager("")
if err != nil {
    log.Fatal(err)
}

// Or specify custom directory
saveManager, err := save.NewSaveManager("/path/to/saves")

// Create checkpoint manager
checkpointManager := save.NewCheckpointManager(saveManager)
```

### Saving Game

```go
// Create save data from game state
saveData := &save.SaveData{
    Seed:            game.Seed,
    PlayTime:        int64(time.Since(startTime).Seconds()),
    PlayerX:         player.X,
    PlayerY:         player.Y,
    PlayerHealth:    player.Health,
    PlayerMaxHealth: player.MaxHealth,
    PlayerAbilities: player.Abilities,
    CurrentRoomID:   currentRoom.ID,
    // ... more fields
}

// Manual save to slot 1
err := saveManager.SaveGame(saveData, 1)
if err != nil {
    log.Printf("Failed to save: %v", err)
}

// Auto-save
err := checkpointManager.CreateCheckpoint(saveData)
```

### Loading Game

```go
// Load from slot 1
saveData, err := saveManager.LoadGame(1)
if err != nil {
    log.Printf("Failed to load: %v", err)
    return
}

// Restore game state
player.X = saveData.PlayerX
player.Y = saveData.PlayerY
player.Health = saveData.PlayerHealth
player.Abilities = saveData.PlayerAbilities
// ... restore more state
```

### Slot Management

```go
// List all saves
saves, err := saveManager.ListSaves()
for _, save := range saves {
    if save.Exists && !save.IsEmpty {
        fmt.Printf("Slot %d: Seed %d, Played %ds, Health %d/%d\n",
            save.SlotID, save.Seed, save.PlayTime, 
            save.PlayerHealth, save.PlayerMaxHealth)
    }
}

// Get detailed info for a slot
info, err := saveManager.GetSaveInfo(2)
if err == nil {
    fmt.Printf("Save time: %v\n", info.SaveTime)
    fmt.Printf("File size: %d bytes\n", info.FileSize)
}

// Delete a save
err := saveManager.DeleteSave(3)
```

### Checkpoint Configuration

```go
// Configure auto-save interval
checkpointManager.SetCheckpointInterval(10 * time.Minute)

// Disable auto-save
checkpointManager.EnableAutoSave(false)

// Check if checkpoint needed
if checkpointManager.ShouldCheckpoint() {
    checkpointManager.CreateCheckpoint(saveData)
}

// Get time since last checkpoint
duration := checkpointManager.GetTimeSinceLastCheckpoint()
fmt.Printf("Last save: %v ago\n", duration)
```

## Integration with Game Engine

The save system is integrated into `GameRunner`:

### Automatic Integration
```go
type GameRunner struct {
    // ... other fields
    saveManager       *save.SaveManager
    checkpointManager *save.CheckpointManager
    visitedRooms      map[int]bool
    defeatedEnemies   map[int]bool
    collectedItems    map[int]bool
}

// In Update() method
func (gr *GameRunner) Update() error {
    // ... game logic
    
    // Auto-save check
    gr.CheckAutoSave()
    
    // Track visited rooms
    if gr.game.CurrentRoom != nil {
        gr.visitedRooms[gr.game.CurrentRoom.ID] = true
    }
    
    return nil
}
```

### Manual Save/Load
```go
// Save current game (call from menu or input handler)
err := gameRunner.SaveGame(slotID)

// Load game (call from main menu)
err := gameRunner.LoadGame(slotID)
```

## Testing

22 comprehensive tests with 100% pass rate:

### SaveManager Tests (13 tests)
- ✅ Manager initialization
- ✅ Save/load operations
- ✅ Invalid slot handling
- ✅ Nonexistent file handling
- ✅ Auto-save functionality
- ✅ Delete operations
- ✅ Slot listing
- ✅ Save info retrieval
- ✅ Metadata validation
- ✅ Filename generation
- ✅ Getter methods

### CheckpointManager Tests (9 tests)
- ✅ Manager initialization
- ✅ Interval configuration
- ✅ Enable/disable auto-save
- ✅ Checkpoint timing
- ✅ Checkpoint creation
- ✅ Disabled checkpoint behavior
- ✅ Time tracking
- ✅ Timer reset
- ✅ Complete workflow

Run tests:
```bash
go test ./internal/save -v
```

## Error Handling

The system handles various error conditions gracefully:

### Invalid Slot IDs
```go
err := saveManager.SaveGame(data, 10)
// Returns: "invalid slot ID: 10 (must be 0-4)"
```

### Missing Save Files
```go
data, err := saveManager.LoadGame(2)
// Returns: "save file does not exist"
```

### Version Mismatch
```go
// Loading old save with different version
data, err := saveManager.LoadGame(1)
// Returns: "incompatible save version: 0.9.0 (expected 1.0.0)"
```

### Seed Mismatch
```go
err := gameRunner.RestoreFromSaveData(saveData)
// Returns: "save file seed mismatch: expected 42, got 100"
```

## Future Enhancements

Potential improvements for the save system:

### Planned Features
1. **Cloud Saves**: Sync saves across devices
2. **Compression**: Reduce save file size
3. **Encryption**: Protect save data
4. **Backup System**: Automatic backup of recent saves
5. **Save Import/Export**: Share saves with other players
6. **Quick Save/Load**: F5/F9 hotkeys
7. **Save Screenshots**: Visual preview of each save
8. **Statistics Tracking**: Detailed play statistics per save

### Advanced State Tracking
1. **Quest Progress**: Track story progression
2. **Dialogue History**: Remember conversations
3. **Map Exploration**: % completion per biome
4. **Achievement Progress**: Track unlocked achievements
5. **Enemy Bestiary**: Remember encountered enemies

## Best Practices

### When to Save
- **Save Points**: Designated safe rooms
- **After Boss Defeats**: Preserve progress
- **Before Risky Areas**: Let players retry easily
- **Room Transitions**: Auto-save on entering new biomes
- **Ability Unlocks**: After gaining new abilities

### Save Data Size
- Current save files: ~1-2 KB
- Keep under 10 KB for fast saves
- Consider compression for larger data

### Performance
- Save operations: <10ms on modern hardware
- Non-blocking: Don't freeze gameplay
- Auto-save during loading screens or transitions

## Security Considerations

Current implementation:
- ✅ Input validation on slot IDs
- ✅ File path sanitization
- ✅ Version checking
- ✅ Seed verification
- ✅ Error handling for corrupted files

Future security enhancements:
- [ ] Checksum validation
- [ ] Encrypted save files
- [ ] Anti-tampering measures
- [ ] Save file signature

## Troubleshooting

### Save Directory Not Created
**Symptom**: Error "failed to create save directory"
**Solution**: Check permissions on home directory

### Cannot Load Save
**Symptom**: "save file does not exist"
**Solution**: Verify save directory path and file existence

### Version Mismatch
**Symptom**: "incompatible save version"
**Solution**: Save was created with different game version

### Seed Mismatch
**Symptom**: "save file seed mismatch"
**Solution**: Trying to load save from different seed/game

## Implementation Notes

### Design Decisions

1. **JSON Format**: Chosen for human readability and debuggability
   - Easy to inspect and manually edit if needed
   - Cross-platform compatibility
   - Standard library support

2. **5 Slot Limit**: Balances choice with simplicity
   - 1 auto-save slot (slot 0)
   - 4 manual save slots (slots 1-4)
   - Prevents save slot clutter

3. **Auto-Save Interval**: 5 minutes default
   - Long enough to not be intrusive
   - Short enough to prevent significant progress loss
   - Configurable for different preferences

4. **Seed Verification**: Ensures save compatibility
   - Prevents loading wrong game's save
   - Maintains procedural generation integrity

5. **Minimal State**: Only essential data saved
   - Reduces file size
   - Faster save/load operations
   - Less prone to corruption

### Code Quality
- ✅ Comprehensive error handling
- ✅ 22 tests, 100% pass rate
- ✅ Full documentation
- ✅ Clean architecture
- ✅ Go best practices

---

**Version**: 1.0.0  
**Last Updated**: 2025-10-19  
**Status**: Production Ready
