# Room Transition System

## Overview

The room transition system provides polished visual transitions when the player moves between connected rooms. Transitions can be customized by type, duration, and direction.

## Transition Types

### Fade (Default)
- Fades to black during the first half of the transition
- Fades from black during the second half
- Smooth and universal, works for all room types
- Best for general use

### Slide
- Slides a black overlay across the screen
- Direction matches the door the player is entering
- Provides directional feedback to the player
- Best for horizontal or vertical room connections

### Iris
- Circular wipe that closes and then opens
- Creates a dramatic, focused transition
- Best for significant room changes (boss rooms, secret areas)
- More performance-intensive than fade or slide

## Usage

### Basic Configuration

```go
// Create transition handler (done automatically in GameRunner)
handler := engine.NewRoomTransitionHandler(game)

// Set transition type
handler.SetTransitionType(engine.TransitionFade)

// Set transition duration (0.3-0.8 seconds)
handler.SetTransitionDuration(0.5) // 0.5 seconds
```

### Per-Room or Per-Door Configuration

To configure transitions per room type or door:

```go
// Example: Use different transitions for different room types
switch room.Type {
case world.BossRoom:
    handler.SetTransitionType(engine.TransitionIris)
    handler.SetTransitionDuration(0.8) // Slower for dramatic effect
case world.CombatRoom:
    handler.SetTransitionType(engine.TransitionSlide)
    handler.SetTransitionDuration(0.4)
default:
    handler.SetTransitionType(engine.TransitionFade)
    handler.SetTransitionDuration(0.5)
}
```

## Technical Details

### Duration Constraints
- Minimum: 0.3 seconds (18 frames at 60 FPS)
- Maximum: 0.8 seconds (48 frames at 60 FPS)
- Default: 0.5 seconds (30 frames at 60 FPS)

### Gameplay Freezing
During transitions, all gameplay is frozen:
- Player input is ignored
- Physics simulation is paused
- Enemies don't move or attack
- Item collection is disabled

This prevents unwanted interactions during screen transitions.

### Direction Mapping
Slide transitions automatically use the door's direction:
- `"east"` or `"right"` → slides left to right
- `"west"` or `"left"` → slides right to left
- `"north"` or `"up"` → slides bottom to top
- `"south"` or `"down"` → slides top to bottom

### Performance Considerations

**Fade**: Fastest, single image overlay  
**Slide**: Fast, single image overlay with translation  
**Iris**: Slower, requires per-pixel calculations for circular effect

For most games, the performance difference is negligible, but on lower-end systems, prefer fade or slide transitions.

## Future Enhancements

Potential improvements for future versions:
- Custom transition curves (ease-in, ease-out, elastic)
- Procedurally synthesized transition-specific sound effects (no bundled audio files)
- Per-door transition type overrides stored in world data
- Procedurally generated animated sprite-based transitions (curtain, dissolve, etc.) — no pre-rendered image assets
- GPU-accelerated shader-based transitions

## Implementation Notes

See `internal/engine/transitions.go` for the core transition logic and `internal/render/renderer.go` for the rendering implementation.
