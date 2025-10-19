# VANIA Particle Effects System

## Overview

The particle effects system provides visual feedback for combat, movement, and environmental effects in the VANIA procedural Metroidvania game. It's designed to be lightweight, flexible, and easy to integrate with the existing game engine.

## Architecture

### Core Components

1. **Particle** - Individual particle with position, velocity, size, color, and lifetime
2. **ParticleEmitter** - Generates and manages groups of particles
3. **ParticleSystem** - Global manager for all particles and emitters
4. **ParticlePresets** - Factory methods for common particle effects

### Package Structure

```
internal/particle/
├── particle.go       # Core particle and emitter classes
├── presets.go        # Preset effect generators
└── particle_test.go  # Comprehensive test suite
```

## Features

### Particle Types

The system supports various particle types for different effects:

#### Combat Particles
- **HitSpark** - Orange sparks when attacks land
- **BloodSplatter** - Red particles when enemies are hit
- **Explosion** - Large burst effect when enemies die

#### Movement Particles
- **DashTrail** - Blue trail particles during dash
- **JumpDust** - Brown dust when jumping
- **LandDust** - Larger dust burst when landing
- **WalkDust** - Subtle dust for walking (not yet implemented)

#### Environmental Particles
- **Rain** - Downward-falling water droplets
- **Snow** - Slowly falling snow particles
- **Embers** - Rising fire particles
- **Sparkles** - Floating yellow sparkles
- **Bubbles** - Rising underwater bubbles

#### Effect Particles
- **Smoke** - Rising gray smoke
- **Lightning** - Brief electric flash
- **DamageNumber** - Floating damage text (data stored in particle)

### Particle Properties

Each particle has:
- **Position** (X, Y) - World coordinates
- **Velocity** (VelX, VelY) - Movement speed
- **Acceleration** (AccelX, AccelY) - Gravity and forces
- **Life/MaxLife** - Lifetime in frames (60 FPS)
- **Size** - Particle radius
- **Color** - RGBA color
- **Alpha** - Transparency (automatically fades)
- **Rotation/RotationSpeed** - Angular properties
- **Type** - Particle type identifier
- **Data** - Custom data (e.g., damage numbers)

### Emitter Configuration

Emitters are highly configurable:
- **EmitRate** - Particles per frame
- **Spread** - Angle spread in radians
- **Speed/SpeedVariance** - Initial velocity range
- **Life/LifeVariance** - Particle lifetime range
- **Size/SizeVariance** - Particle size range
- **Gravity** - Vertical acceleration
- **Color** - Base particle color
- **OneShot** - Emit once vs continuous

## Usage

### Basic Usage

```go
// Create particle system (in game initialization)
particleSystem := particle.NewParticleSystem(1000) // Max 1000 particles
presets := &particle.ParticlePresets{}

// Create a hit effect
emitter := presets.CreateHitEffect(x, y, direction)
emitter.Burst(10) // Emit 10 particles
particleSystem.AddEmitter(emitter)

// Update every frame
particleSystem.Update()

// Get all particles for rendering
allParticles := particleSystem.GetAllParticles()
```

### Preset Effects

```go
// Jump dust when player jumps
emitter := presets.CreateJumpDust(playerX, playerY)
emitter.Burst(8)
particleSystem.AddEmitter(emitter)

// Dash trail (continuous)
emitter := presets.CreateDashTrail(playerX, playerY)
emitter.Start()
particleSystem.AddEmitter(emitter)
// Later: emitter.Stop()

// Enemy hit
hitEmitter := presets.CreateHitEffect(enemyX, enemyY, attackDirection)
hitEmitter.Burst(10)
particleSystem.AddEmitter(hitEmitter)

bloodEmitter := presets.CreateBloodSplatter(enemyX, enemyY, attackDirection)
bloodEmitter.Burst(6)
particleSystem.AddEmitter(bloodEmitter)

// Enemy death explosion
explosionEmitter := presets.CreateExplosion(enemyX, enemyY, 1.0)
explosionEmitter.Burst(20)
particleSystem.AddEmitter(explosionEmitter)

// Landing dust
emitter := presets.CreateLandDust(playerX, playerY)
emitter.Burst(12)
particleSystem.AddEmitter(emitter)
```

### Custom Particles

```go
// Create a custom particle
p := particle.NewParticle(
    x, y,           // Position
    velX, velY,     // Velocity
    60,             // Life (frames)
    3.0,            // Size
    color.RGBA{255, 200, 100, 255}, // Color
    particle.HitSpark, // Type
)
p.AccelY = 0.2 // Add gravity
particleSystem.AddParticle(p)
```

## Integration

### Rendering Integration

The particle system is integrated with the Ebiten renderer:

```go
// In internal/render/renderer.go
func (r *Renderer) RenderParticles(screen *ebiten.Image, particles []*particle.Particle) {
    for _, p := range particles {
        // Render particle as simple colored square
        // Position relative to camera
        // Skip if outside screen bounds
        // Apply alpha and rotation
    }
}
```

### Game Engine Integration

Particles are triggered at key moments in `internal/engine/runner.go`:

1. **Jump** - Dust particles when leaving ground
2. **Land** - Larger dust burst when touching ground
3. **Dash** - Trail effect during dash
4. **Attack Hit** - Sparks and blood when hitting enemies
5. **Enemy Death** - Explosion effect

## Performance

- **Max Particles**: Configurable limit (default 1000)
- **Max Emitters**: Limited to 1/10 of max particles (default 100)
- **Auto Cleanup**: Dead particles and one-shot emitters automatically removed
- **Culling**: Particles outside screen bounds not rendered
- **Efficiency**: Simple square rendering, minimal overhead

## Testing

Comprehensive test suite covers:
- Particle creation and lifecycle (19 tests)
- Emitter behavior (burst, continuous, one-shot)
- ParticleSystem management
- Automatic cleanup
- Edge cases and limits

Run tests:
```bash
go test ./internal/particle -v
```

All tests pass: ✅ 19/19

## API Reference

### ParticleSystem

```go
func NewParticleSystem(maxParticles int) *ParticleSystem
func (ps *ParticleSystem) Update()
func (ps *ParticleSystem) AddEmitter(emitter *ParticleEmitter)
func (ps *ParticleSystem) AddParticle(particle *Particle)
func (ps *ParticleSystem) GetAllParticles() []*Particle
func (ps *ParticleSystem) Clear()
func (ps *ParticleSystem) GetParticleCount() int
```

### ParticleEmitter

```go
func NewParticleEmitter(x, y float64, ptype ParticleType) *ParticleEmitter
func (e *ParticleEmitter) Update()
func (e *ParticleEmitter) EmitParticles(count int)
func (e *ParticleEmitter) Start()
func (e *ParticleEmitter) Stop()
func (e *ParticleEmitter) SetPosition(x, y float64)
func (e *ParticleEmitter) Burst(count int)
```

### Particle

```go
func NewParticle(x, y, velX, velY float64, life int, size float64, col color.RGBA, ptype ParticleType) *Particle
func (p *Particle) Update()
func (p *Particle) IsAlive() bool
```

### ParticlePresets

```go
func (pp *ParticlePresets) CreateHitEffect(x, y, direction float64) *ParticleEmitter
func (pp *ParticlePresets) CreateDashTrail(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateJumpDust(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateLandDust(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateWalkDust(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateBloodSplatter(x, y, direction float64) *ParticleEmitter
func (pp *ParticlePresets) CreateExplosion(x, y, size float64) *ParticleEmitter
func (pp *ParticlePresets) CreateSmoke(x, y float64, continuous bool) *ParticleEmitter
func (pp *ParticlePresets) CreateRain(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateSnow(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateEmbers(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateSparkles(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateBubbles(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateLightning(x, y float64) *ParticleEmitter
func (pp *ParticlePresets) CreateDamageNumber(x, y float64, damage int) *Particle
```

## Future Enhancements

Potential improvements for future versions:

1. **Walk Dust** - Periodic dust particles while walking
2. **Biome Particles** - Environmental particles per biome (rain, snow, embers)
3. **Damage Numbers** - Render floating damage text
4. **Trail Optimization** - Better trail effect management
5. **Particle Sprites** - Use actual sprites instead of squares
6. **Blend Modes** - Additive blending for glowing effects
7. **Collision** - Particles interact with world geometry
8. **Sound Integration** - Play sounds with particle effects

## Credits

Particle system designed and implemented for VANIA procedural Metroidvania game engine.
Follows Go best practices with comprehensive testing and documentation.
