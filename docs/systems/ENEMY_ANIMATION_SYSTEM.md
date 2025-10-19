# Enemy Animation System

## Overview

The enemy animation system extends the existing animation framework to provide procedurally generated animation frames for all enemies in the game. Each enemy receives a complete set of animations (idle, patrol, attack, death, hit) generated from their base sprite, creating dynamic and visually interesting enemy behaviors.

## Architecture

### Core Components

1. **Animation Generator Extensions** (`internal/animation/generator.go`)
   - `GenerateEnemyIdleFrames()` - Subtle breathing animations for idle enemies
   - `GenerateEnemyPatrolFrames()` - Walking/bobbing animations for moving enemies
   - `GenerateEnemyAttackFrames()` - Forward-leaning attack animations
   - `GenerateEnemyDeathFrames()` - Fade-out death animations
   - `GenerateHitFrames()` - Flash/tint animations for taking damage

2. **Enemy Instance** (`internal/entity/ai.go`)
   - `AnimController` field for managing animation state
   - `CreateEnemyAnimController()` - Initializes animations for an enemy
   - Animation state transitions based on AI behavior

3. **Rendering** (`internal/render/renderer.go`)
   - `RenderEnemy()` updated to render animated sprites
   - Supports fallback to colored rectangles if no sprite available

## Animation Types

### Idle Animation
- **Purpose**: Default animation when enemy is stationary
- **Frames**: 4 frames
- **Frame Time**: 15 ticks
- **Loops**: Yes
- **Effect**: Subtle vertical "breathing" movement

### Patrol Animation
- **Purpose**: Movement animation during patrol/chase/flee behaviors
- **Frames**: 4 frames
- **Frame Time**: 8 ticks
- **Loops**: Yes
- **Effect**: Bobbing motion simulating walking/movement

### Attack Animation
- **Purpose**: Played during enemy attack state
- **Frames**: 3 frames
- **Frame Time**: 5 ticks
- **Loops**: No
- **Effect**: Forward lean suggesting aggressive motion

### Death Animation
- **Purpose**: Played when enemy health reaches 0
- **Frames**: 4 frames
- **Frame Time**: 10 ticks
- **Loops**: No
- **Effect**: Fade out with downward shift

### Hit Animation
- **Purpose**: Brief flash when enemy takes damage
- **Frames**: 2 frames
- **Frame Time**: 3 ticks
- **Loops**: No
- **Effect**: Red tint flash

## Implementation Details

### Enemy Sprite Generation

Enemy sprites are now generated during the `generateEntities()` phase in `internal/engine/game.go`:

```go
// Generate sprite for this enemy
enemySize := 32
if enemy.Size == entity.SmallEnemy {
    enemySize = 16
} else if enemy.Size == entity.LargeEnemy {
    enemySize = 48
} else if enemy.Size == entity.BossEnemy {
    enemySize = 64
}
enemySpriteGen := graphics.NewSpriteGenerator(enemySize, enemySize, graphics.VerticalSymmetry)
enemy.SpriteData = enemySpriteGen.Generate(seed)
```

### Animation Controller Initialization

When an `EnemyInstance` is created, the animation controller is automatically initialized:

```go
// Create animation controller if sprite data is available
var animController *animation.AnimationController
if sprite, ok := enemy.SpriteData.(*graphics.Sprite); ok && sprite != nil {
    animController = CreateEnemyAnimController(sprite, enemy)
}
```

### Animation State Management

The enemy's `Update()` function manages animation state transitions:

```go
switch ei.State {
case AttackState:
    if currentAnim != "attack" {
        ei.AnimController.Play("attack", true)
    }
case PatrolState, ChaseState, FleeState:
    if currentAnim != "patrol" && currentAnim != "attack" {
        ei.AnimController.Play("patrol", false)
    }
case IdleState:
    if currentAnim != "idle" && currentAnim != "attack" {
        ei.AnimController.Play("idle", false)
    }
case DeadState:
    if currentAnim != "death" {
        ei.AnimController.Play("death", true)
    }
}
```

### Rendering with Animations

The renderer uses the current animation frame when drawing enemies:

```go
// Get current animation frame if available
var spriteToRender *graphics.Sprite
if enemy.AnimController != nil {
    spriteToRender = enemy.AnimController.GetCurrentFrame()
}
// Fallback to base sprite if no animation frame
if spriteToRender == nil {
    if sprite, ok := enemy.Enemy.SpriteData.(*graphics.Sprite); ok {
        spriteToRender = sprite
    }
}

gr.renderer.RenderEnemy(screen, ex, ey, ew, eh, 
    enemy.CurrentHealth, enemy.Enemy.Health, false, spriteToRender)
```

## Deterministic Generation

Enemy animations are deterministically generated using a seed derived from:
- Enemy's danger level
- Enemy's biome type (as character codes)

This ensures:
- Same seed always produces same animations
- Different enemies have varied animations
- Reproducible across game sessions

```go
seed := int64(enemy.DangerLevel * 12345)
if enemy.BiomeType != "" {
    for _, c := range enemy.BiomeType {
        seed += int64(c)
    }
}
animGen := animation.NewAnimationGenerator(seed)
```

## Performance Considerations

### Frame Generation
- Animation frames are generated once when enemy instance is created
- Frames are cached in the AnimationController
- No per-frame generation overhead

### Memory Usage
- Each enemy stores 17 total animation frames:
  - 4 idle frames
  - 4 patrol frames
  - 3 attack frames
  - 4 death frames
  - 2 hit frames
- Memory scales linearly with number of active enemies
- Dead enemies can have their animation controllers released

### Update Efficiency
- Animation state transitions use string comparisons
- Update only occurs once per frame per enemy
- No animation updates for dead enemies after death animation completes

## Testing

The enemy animation system includes comprehensive tests:

### Animation Generator Tests
- `TestGenerateEnemyIdleFrames` - Validates idle frame generation
- `TestGenerateEnemyPatrolFrames` - Validates patrol frame generation
- `TestGenerateEnemyAttackFrames` - Validates attack frame generation
- `TestGenerateEnemyDeathFrames` - Validates death frame generation
- `TestEnemyAnimationNilSprite` - Ensures nil sprite handling
- `TestEnemyAnimationZeroFrames` - Ensures zero frame count handling
- `TestEnemyAnimationDeterminism` - Validates deterministic generation

### Enemy Instance Tests
- `TestCreateEnemyAnimController` - Validates controller creation
- `TestEnemyInstanceWithAnimController` - Validates integration
- `TestEnemyAnimationStateTransitions` - Validates state transitions
- `TestEnemyHitAnimation` - Validates damage response
- `TestEnemyInstanceNoSpriteData` - Validates graceful degradation

## Integration with Existing Systems

### AI System
- Enemy AI states (Idle, Patrol, Chase, Attack, Dead) drive animation selection
- No changes needed to existing AI behavior logic
- Animations automatically reflect AI decisions

### Combat System
- `TakeDamage()` triggers hit animation
- Death state triggers death animation
- Attack state triggers attack animation

### Rendering System
- Renderer seamlessly handles animated vs. static sprites
- Health bars remain positioned correctly
- Camera system unaffected

## Future Enhancements

Potential improvements for the enemy animation system:

1. **Boss-Specific Animations**
   - Unique animation sets for boss enemies
   - Phase-specific animations
   - Special attack animations

2. **Directional Sprites**
   - Different animations for left/right facing
   - Sprite flipping based on movement direction

3. **Variable Frame Rates**
   - Different frame rates based on enemy speed
   - Slow-motion effects for powerful attacks

4. **Transition Animations**
   - Smooth interpolation between animation states
   - Blend frames for state changes

5. **Particle Integration**
   - Spawn particles during attack frames
   - Death particle effects synchronized with death animation

## Usage Example

```go
// Enemy with animations is created automatically during game generation
enemy := &entity.Enemy{
    Name:        "Shadow Crawler",
    Health:      100,
    Damage:      15,
    Speed:       2.5,
    Size:        entity.MediumEnemy,
    Behavior:    entity.PatrolBehavior,
    BiomeType:   "cave",
    SpriteData:  generatedSprite, // Created during generation
}

// Instance automatically gets animation controller
instance := entity.NewEnemyInstance(enemy, x, y)

// Animation plays automatically during update
instance.Update(playerX, playerY)

// Current frame used during rendering
frame := instance.AnimController.GetCurrentFrame()
```

## Conclusion

The enemy animation system provides a complete, procedurally-generated animation solution for all enemies in VANIA. By extending the existing animation framework, it maintains consistency with the player animation system while adding visual variety and polish to enemy behaviors. The system is fully deterministic, well-tested, and integrates seamlessly with existing AI, combat, and rendering systems.
