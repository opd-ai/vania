# Door System Documentation

## Overview

The Door System implements ability-gated progression, a core mechanic of Metroidvania games. Players must acquire specific abilities to unlock doors that lead to new areas, creating a non-linear exploration experience with gradual world access.

## Features

### Core Functionality

1. **Locked Doors**: Doors can be locked and require specific abilities to pass through
2. **Automatic Unlocking**: When a player acquires the required ability and approaches a locked door, it automatically unlocks
3. **Visual Feedback**: 
   - On-screen messages when approaching locked doors
   - Sparkle particle effects when doors unlock
   - Clear indication of required abilities
4. **Persistent State**: Door unlock states are saved and restored with game saves
5. **Ability-Based Progression**: Doors are tied to the world graph's ability requirements (e.g., "double_jump", "dash")

## Architecture

### Components

#### 1. Door Structure (`internal/world/graph_gen.go`)

```go
type Door struct {
    X, Y          int    // Position in the room
    Width, Height int    // Size of the door
    Direction     string // "north", "south", "east", "west"
    LeadsTo       *Room  // Connected room
    Locked        bool   // Whether door requires ability/key
}
```

#### 2. RoomTransitionHandler (`internal/engine/transitions.go`)

Manages all door-related functionality:

**Key Methods:**

- `CheckDoorCollision(playerX, playerY, playerW, playerH, unlockedDoors) *Door`
  - Checks if player is colliding with a door
  - Returns nil if door is locked and not unlocked
  - Returns door pointer if accessible

- `GetDoorKey(door) string`
  - Generates unique identifier for a door
  - Format: `"room_{roomID}_door_{x}_{y}_{direction}"`
  - Used for tracking unlock state

- `CanUnlockDoor(door, playerAbilities, collectedItems) bool`
  - Checks if player has required ability to unlock door
  - Looks up ability requirement from world graph edges
  - Returns true if player can unlock the door

- `findEdgeRequirement(fromRoomID, toRoomID) string`
  - Searches world graph for connection requirements
  - Checks both forward and reverse directions
  - Returns ability name (e.g., "double_jump") or empty string

#### 3. GameRunner (`internal/engine/runner.go`)

Integrates door system with game loop:

**New Fields:**
```go
unlockedDoors     map[string]bool  // Tracks which doors have been unlocked
lockedDoorMessage string           // Message to display
lockedDoorTimer   int              // How long to show message (frames)
```

**Key Methods:**

- `checkLockedDoorInteraction()`
  - Called every frame to check for locked door proximity
  - Automatically unlocks doors when player has required ability
  - Shows appropriate messages

- `UnlockDoor(door)`
  - Marks door as unlocked in state
  - Creates sparkle particle effect
  - Shows "Door unlocked!" message

- `getBossesDefeated() []int`
  - Helper to extract boss IDs from defeated enemies
  - Used for save system

#### 4. Save System Integration

**SaveData Structure** includes:
```go
UnlockedDoors   map[string]bool `json:"unlocked_doors"`
BossesDefeated  []int           `json:"bosses_defeated"`
```

Door unlock states persist across:
- Manual saves
- Auto-saves
- Checkpoints
- Game restarts

## Usage

### For Players

When approaching a locked door:

1. **Without Required Ability**: Message displays "Requires: {ability_name}"
2. **With Required Ability**: Door automatically unlocks with sparkle effect
3. **Already Unlocked**: Door opens normally, allowing passage

### For Developers

#### Creating Locked Doors

Doors are automatically locked based on world graph edge requirements:

```go
// In world generation
edge := world.GraphEdge{
    From:        roomID1,
    To:          roomID2,
    Requirement: "double_jump", // This ability is required
}
```

The corresponding door will be locked until player has "double_jump" ability.

#### Adding New Abilities

1. Add ability to entity generation system
2. Set unlock order for progression
3. Reference ability name in world graph edges
4. Doors will automatically require and unlock based on these

#### Checking Door State

```go
// Check if specific door is unlocked
doorKey := transitionHandler.GetDoorKey(door)
isUnlocked := runner.unlockedDoors[doorKey]

// Check if player can unlock a door
canUnlock := transitionHandler.CanUnlockDoor(
    door, 
    player.Abilities, 
    collectedItems,
)
```

## Technical Details

### Door Key Format

Door keys uniquely identify doors across the game world:
```
Format: "room_{roomID}_door_{x}_{y}_{direction}"
Example: "room_42_door_100_200_east"
```

This ensures:
- Each door has a unique identifier
- Keys are consistent across save/load
- No collisions between similar doors in different rooms

### Collision Detection

Doors use AABB (Axis-Aligned Bounding Box) collision:

```go
playerX < doorX+doorW &&
playerX+playerW > doorX &&
playerY < doorY+doorH &&
playerY+playerH > doorY
```

This provides pixel-perfect collision with doors.

### Message Display

Locked door messages:
- Display for 2 seconds (120 frames at 60 FPS)
- Shown in center of screen with semi-transparent background
- Includes ability requirement or generic "Door is locked" message

### Particle Effects

Unlock effects use the particle system:
```go
sparkleEmitter := particlePresets.CreateSparkles(doorX, doorY, 1.0)
sparkleEmitter.Burst(15)
particleSystem.AddEmitter(sparkleEmitter)
```

Creates visual feedback with 15 sparkle particles.

## Integration with Other Systems

### World Generation
- Doors are placed by world generator
- Lock state determined by graph edge requirements
- Ensures progression gates are properly placed

### Save System
- Door states persist in SaveData
- Restored on game load
- Includes auto-save and checkpoints

### Particle System
- Unlock effects use existing particle system
- No new particle types needed
- Sparkles provide clear visual feedback

### Combat System
- Independent of combat
- Can be extended to require boss defeats
- Currently based purely on abilities

## Testing

### Test Coverage

New tests in `internal/engine/transitions_test.go`:

1. `TestRoomTransitionHandler_CheckDoorCollision_Locked`
   - Verifies locked doors block passage
   - Tests unlocked state bypass

2. `TestRoomTransitionHandler_GetDoorKey`
   - Validates door key generation
   - Ensures consistency

3. `TestRoomTransitionHandler_CanUnlockDoor`
   - Tests ability checking
   - Validates unlock logic

4. `TestRoomTransitionHandler_findEdgeRequirement`
   - Tests requirement lookup
   - Checks bidirectional edges

### Manual Testing

To test the door system:

1. Generate a game with seed
2. Navigate to a locked door
3. Verify "Requires: {ability}" message
4. Acquire required ability
5. Return to door
6. Verify automatic unlock with sparkles
7. Save game
8. Load game
9. Verify door remains unlocked

## Performance Considerations

### Efficiency
- Door state is O(1) lookup with map
- Key generation only happens when needed
- No performance impact on game loop

### Memory
- Each unlocked door: ~50 bytes (string key + bool)
- Typical game: 50-100 doors
- Total memory: ~5KB (negligible)

## Future Enhancements

### Potential Additions

1. **Key Items**: Extend to support key item requirements (not just abilities)
2. **One-Way Doors**: Doors that can only be used in one direction
3. **Timed Doors**: Doors that open/close based on timer
4. **Switch-Activated Doors**: Doors controlled by room switches/puzzles
5. **Multi-Requirement Doors**: Doors requiring multiple abilities
6. **Visual Door States**: Different sprites for locked/unlocked doors
7. **Door Animations**: Animated opening/closing sequences
8. **Sound Effects**: Audio feedback for locked doors and unlocking

### Extensibility

The system is designed to be easily extended:

```go
// Example: Add key item requirement
type Door struct {
    // ... existing fields
    RequiredKeyItem int  // Item ID required to unlock
}

// Example: Add switch requirement
type Door struct {
    // ... existing fields
    RequiredSwitchID string  // Switch that must be activated
}
```

## Troubleshooting

### Common Issues

**Issue**: Door not unlocking despite having ability
- **Check**: Ability name matches exactly (case-sensitive)
- **Check**: World graph edge exists with correct requirement
- **Check**: Door.LeadsTo is properly set

**Issue**: Door state not persisting after save/load
- **Check**: SaveManager initialized properly
- **Check**: UnlockedDoors field in SaveData populated
- **Check**: LoadGame correctly restores unlockedDoors map

**Issue**: Message not showing
- **Check**: lockedDoorTimer > 0
- **Check**: lockedDoorMessage not empty string
- **Check**: Player actually near locked door

## Code Examples

### Example 1: Check if player near locked door

```go
func (gr *GameRunner) checkLockedDoorInteraction() {
    if gr.game.CurrentRoom == nil {
        return
    }
    
    for i := range gr.game.CurrentRoom.Doors {
        door := &gr.game.CurrentRoom.Doors[i]
        
        // Check collision
        if isPlayerNearDoor(door) {
            doorKey := gr.transitionHandler.GetDoorKey(door)
            
            if door.Locked && !gr.unlockedDoors[doorKey] {
                // Handle locked door
                if canUnlock {
                    gr.UnlockDoor(door)
                } else {
                    gr.showLockedMessage(door)
                }
            }
        }
    }
}
```

### Example 2: Custom unlock logic

```go
// Unlock door only after defeating boss
func (gr *GameRunner) UnlockBossDoor(door *world.Door, bossID int) {
    if !gr.defeatedEnemies[bossID] {
        gr.lockedDoorMessage = "Defeat the boss to proceed"
        return
    }
    
    gr.UnlockDoor(door)
}
```

## Conclusion

The Door System provides core Metroidvania progression mechanics with:
- ✅ Ability-gated exploration
- ✅ Visual and textual feedback
- ✅ Save/load persistence
- ✅ Extensible architecture
- ✅ Comprehensive testing
- ✅ Performance optimized

This system enables non-linear world exploration while maintaining clear progression gates, a hallmark of the Metroidvania genre.
