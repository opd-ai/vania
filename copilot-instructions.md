# Copilot Instructions for VANIA

## Project Overview

VANIA is a procedural Metroidvania game engine written in pure Go that generates ALL game content (graphics, audio, narrative, levels) algorithmically at runtime from a single seed value. The system uses deterministic random generation, cryptographic seed derivation, and specialized algorithms (cellular automata for sprites, additive synthesis for audio, graph theory for worlds) to create infinite unique, playable experiences. The architecture is organized around a master seed that derives independent subsystem seeds through SHA-256 hashing, ensuring reproducibility across runs while maintaining subsystem independence.

## Code Organization

### Directory Structure
The project follows Go's standard layout with clear separation between entry points and internal packages:
- **cmd/game/**: CLI entry point with flag parsing and display formatting
- **internal/pcg/**: Core PCG framework (seed management, caching, validation)
- **internal/graphics/**: Sprite, tileset, and palette generation
- **internal/audio/**: Waveform synthesis, sound effects, and music generation
- **internal/narrative/**: Story, character, and faction generation
- **internal/world/**: Graph-based level generation with biome system
- **internal/entity/**: Enemy, boss, item, and ability generation
- **internal/engine/**: Game state integration and orchestration

### Package Naming
- Use lowercase, single-word package names: `pcg`, `graphics`, `audio`, `world`, `entity`
- Package names should be nouns describing functionality
- Avoid generic names like `util`, `common`, or `helpers`
- Package names should match directory names exactly

### File Naming
- Use snake_case: `sprite_gen.go`, `music_gen.go`, `enemy_gen.go`
- Test files: `*_test.go` suffix (e.g., `seed_test.go`)
- Generator files: `*_gen.go` suffix for procedural generators
- Related functionality in same file (e.g., all sprite generation in `sprite_gen.go`)

## Coding Standards

### Error Handling

**Pattern**: Return errors, don't panic. Only panic for truly unrecoverable situations.

We generally do **implicit error handling** - functions that generate content return values without errors unless validation fails. This is acceptable because:
- Generation is deterministic from seeds
- Validation happens at the system level
- Invalid generation can be retried with different seeds

**Example - Typical generator pattern:**
```go
// Good - Simple return for deterministic generation
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rng := rand.New(rand.NewSource(seed))
    // Generation logic
    return sprite
}

// Good - Error when validation is needed
func (gg *GameGenerator) GenerateCompleteGame() (*Game, error) {
    // Generate all systems
    if !gg.validate(worldData, entities, narrative) {
        return nil, fmt.Errorf("failed validation: world has no start room")
    }
    return game, nil
}
```

**Anti-pattern:**
```go
// Bad - Unnecessary error for deterministic generation
func (sg *SpriteGenerator) Generate(seed int64) (*Sprite, error) {
    // This always succeeds given valid inputs
    return sprite, nil
}
```

### Naming Conventions

**Variables**: 
- Use camelCase for local variables: `masterSeed`, `biomeIdx`, `roomCount`
- Short names for short scopes: `i`, `j` for loops; `x`, `y` for coordinates
- Descriptive names for wider scopes: `currentBranch`, `criticalPathLength`
- Acronyms in caps: `RNG`, `PCG`, `SFX`, but `rng`, `sfx` when lowercase

**Functions**:
- Use PascalCase for exported: `NewSpriteGenerator`, `Generate`, `DeriveSeeds`
- Use camelCase for unexported: `generateGraph`, `selectMusicGenerator`, `applySymmetry`
- Use verbs for actions: `Generate`, `Create`, `Apply`, `Calculate`, `Validate`
- Constructor pattern: `New<Type>` returns initialized `*<Type>`

**Constants**:
- Use PascalCase for exported constants: `SineWave`, `JumpSFX`, `FantasyTheme`
- Group related constants with `iota` for enums:
```go
type WaveType int

const (
    SineWave WaveType = iota
    SquareWave
    SawtoothWave
    TriangleWave
)
```

**Interfaces**:
- Name with `-er` suffix when possible: `Generator`, `Validator`, `Synthesizer`
- Prefer small, focused interfaces (single method when possible)
- We don't heavily use interfaces in this codebase - concrete types are preferred for generators

### Testing

**Test file naming**: `*_test.go` in the same package directory

**Table-driven tests**: Use for testing multiple similar cases
```go
func TestSymmetryTypes(t *testing.T) {
    testCases := []struct {
        name     string
        symmetry SymmetryType
    }{
        {"NoSymmetry", NoSymmetry},
        {"HorizontalSymmetry", HorizontalSymmetry},
        {"VerticalSymmetry", VerticalSymmetry},
        {"RadialSymmetry", RadialSymmetry},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            gen := NewSpriteGenerator(16, 16, tc.symmetry)
            sprite := gen.Generate(555)
            
            if sprite == nil {
                t.Errorf("Failed to generate sprite with %s", tc.name)
            }
        })
    }
}
```

**Determinism tests**: Critical for PCG - always test that same seed produces same output
```go
func TestSpriteDeterminism(t *testing.T) {
    gen := NewSpriteGenerator(32, 32, HorizontalSymmetry)
    seed := int64(999)
    
    sprite1 := gen.Generate(seed)
    sprite2 := gen.Generate(seed)
    
    // Compare pixels to verify determinism
    for y := 0; y < 5; y++ {
        for x := 0; x < 5; x++ {
            c1 := sprite1.Image.RGBAAt(x, y)
            c2 := sprite2.Image.RGBAAt(x, y)
            
            if c1.R != c2.R || c1.G != c2.G || c1.B != c2.B || c1.A != c2.A {
                t.Errorf("Sprites not deterministic at (%d,%d)", x, y)
                return
            }
        }
    }
}
```

**Mocking approach**: We don't use mocking libraries. Instead:
- Inject seeds for deterministic testing
- Use concrete types for generators
- Test integration at generator level, not individual functions

**Coverage expectations**: 
- Core PCG framework: 100% coverage (seed.go, cache.go)
- Generator packages: Test main generation path and determinism
- Not all helper functions need tests if they're tested through generators
- Current coverage: pcg ✅, graphics ✅, audio ✅, others have no tests yet

### Concurrency

**Thread-safe caching**: Use `sync.RWMutex` for concurrent read access
```go
type AssetCache struct {
    Sprites map[string]interface{}
    mu      sync.RWMutex
}

func (c *AssetCache) GetSprite(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.Sprites[key]
    return val, ok
}

func (c *AssetCache) SetSprite(key string, sprite interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Sprites[key] = sprite
}
```

**Goroutines**: Currently not used in generation pipeline (sequential is simpler and fast enough)
- Future work: Parallel generation of independent subsystems (graphics, audio, narrative)
- Would require: WaitGroup, error channels, proper synchronization

**Context**: Not currently used but should be added for:
- Cancellable generation
- Timeout support
- Request-scoped values

### Dependencies

**Current dependencies**: Only Go standard library
- `math/rand`: Deterministic RNG
- `crypto/sha256`: Seed derivation
- `image`, `image/color`: Sprite generation
- `sync`: Thread-safe caching
- No external dependencies!

**Adding new dependencies**:
1. Evaluate if standard library can solve the problem
2. For graphics/audio, consider pure Go implementations first
3. Update `go.mod` with `go get <package>`
4. Document why the dependency is needed

**Preferred patterns**:
- Pure functions that take RNG as parameter
- Builders/generators as structs with configuration
- Avoid global state

## Architecture Patterns

### Seed Derivation Pattern

**Purpose**: Generate independent subsystem seeds from a master seed for reproducibility

**When to use**: Any time you need deterministic randomness that should be independent from other systems

**Implementation**:
```go
// HashSeed derives a subsystem seed from master seed + identifier
func HashSeed(masterSeed int64, identifier string) int64 {
    h := sha256.New()
    binary.Write(h, binary.LittleEndian, masterSeed)
    h.Write([]byte(identifier))
    sum := h.Sum(nil)
    return int64(binary.LittleEndian.Uint64(sum[:8]))
}

// DeriveSeeds generates all subsystem seeds from master
func DeriveSeeds(masterSeed int64) map[string]int64 {
    return map[string]int64{
        "graphics":  HashSeed(masterSeed, "graphics"),
        "audio":     HashSeed(masterSeed, "audio"),
        "narrative": HashSeed(masterSeed, "narrative"),
        "world":     HashSeed(masterSeed, "world"),
        "entity":    HashSeed(masterSeed, "entity"),
    }
}
```

**Why**: Ensures each subsystem's random generation is independent but reproducible. Same master seed always produces same derived seeds.

### Generator Pattern

**Purpose**: Encapsulate procedural generation logic with configuration

**Structure**:
```go
// Generator struct holds configuration
type SpriteGenerator struct {
    Width       int
    Height      int
    Symmetry    SymmetryType
    Constraints SpriteConstraints
}

// Constructor with sensible defaults
func NewSpriteGenerator(width, height int, symmetry SymmetryType) *SpriteGenerator {
    return &SpriteGenerator{
        Width:    width,
        Height:   height,
        Symmetry: symmetry,
        Constraints: SpriteConstraints{
            MinDensity:     0.3,
            MaxDensity:     0.7,
            RequireOutline: true,
            ColorCount:     6,
        },
    }
}

// Generate method takes seed, returns generated content
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rng := rand.New(rand.NewSource(seed))
    // Generation logic
    return sprite
}
```

**Key principles**:
- Constructor (`New*`) creates generator with configuration
- `Generate` method takes seed, returns content
- Generator is reusable (can call Generate multiple times)
- RNG is created per-generation for determinism

### Composition Over Inheritance

**Pattern**: Build complex generators by composing simpler ones

**Example**:
```go
type GraphicsSystem struct {
    SpriteGen  *graphics.SpriteGenerator
    TilesetGen *graphics.TilesetGenerator
    PaletteGen *graphics.PaletteGenerator
    Tilesets   map[string]*graphics.Tileset
    Sprites    map[string]*graphics.Sprite
}

func (gg *GameGenerator) generateGraphics(narrative *narrative.WorldContext) *GraphicsSystem {
    system := &GraphicsSystem{
        SpriteGen:  graphics.NewSpriteGenerator(32, 32, graphics.VerticalSymmetry),
        TilesetGen: graphics.NewTilesetGenerator(16, string(narrative.Theme)),
        PaletteGen: graphics.NewPaletteGenerator(graphics.AnalogousScheme),
        Tilesets:   make(map[string]*graphics.Tileset),
        Sprites:    make(map[string]*graphics.Sprite),
    }
    
    // Use generators to populate maps
    system.Sprites["player"] = system.SpriteGen.Generate(gg.GraphicsGen.Seed)
    return system
}
```

### Validation Pattern

**Purpose**: Ensure generated content meets quality standards

**Implementation**:
```go
type Validator struct {
    minQualityScore float64
}

func NewValidator(minScore float64) *Validator {
    return &Validator{minQualityScore: minScore}
}

func (v *Validator) MeetsThreshold(metrics *QualityMetrics) bool {
    score := v.CalculateQualityScore(metrics)
    return score >= v.minQualityScore
}
```

**When to use**:
- After generating complete game to ensure playability
- Before caching assets
- When generation quality affects gameplay

## Common Patterns and Idioms

### Deterministic RNG

**Always create local RNG from seed**:
```go
// Good - deterministic
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rng := rand.New(rand.NewSource(seed))
    // Use rng for all random decisions
}

// Bad - uses global rand, not deterministic
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rand.Seed(seed) // Modifies global state!
    x := rand.Intn(10) // Uses global rand
}
```

### Cellular Automata for Organic Shapes

**Pattern used for sprite generation**:
```go
func (sg *SpriteGenerator) cellularAutomataStep(grid [][]bool) [][]bool {
    newGrid := make([][]bool, len(grid))
    for y := range newGrid {
        newGrid[y] = make([]bool, len(grid[y]))
        for x := range newGrid[y] {
            neighbors := sg.countNeighbors(grid, x, y)
            if grid[y][x] {
                newGrid[y][x] = neighbors >= 4  // Stay alive
            } else {
                newGrid[y][x] = neighbors >= 5  // Birth
            }
        }
    }
    return newGrid
}
```

### Procedural Color from HSV

**Converting HSV to RGB for harmonious palettes**:
```go
func (sg *SpriteGenerator) generatePalette(rng *rand.Rand, count int) []color.RGBA {
    palette := make([]color.RGBA, count)
    baseHue := rng.Float64() * 360.0
    
    for i := range palette {
        hue := baseHue + float64(i)*30.0  // Analogous colors
        for hue >= 360.0 {
            hue -= 360.0
        }
        saturation := 0.6 + rng.Float64()*0.3
        value := 0.4 + float64(i)*0.1
        palette[i] = hsvToRGB(hue, saturation, value)
    }
    return palette
}
```

### Graph-Based World Generation

**Ensuring playable, connected worlds**:
```go
func (wg *WorldGenerator) generateGraph(world *World) {
    roomID := 0
    
    // Create critical path (guaranteed progression)
    criticalPathLength := 15 + wg.rng.Intn(10)
    currentNode := 0
    
    for i := 0; i < criticalPathLength; i++ {
        world.Graph.Nodes[roomID] = &GraphNode{
            RoomID:   roomID,
            Depth:    i + 1,
            Required: true,  // On critical path
        }
        
        // Ability-gate progression every 5 rooms
        requirement := ""
        if i > 0 && i%5 == 0 {
            abilities := []string{"double_jump", "dash", "wall_climb", "glide"}
            requirement = abilities[wg.rng.Intn(len(abilities))]
        }
        
        world.Graph.Edges = append(world.Graph.Edges, GraphEdge{
            From:        currentNode,
            To:          roomID,
            Requirement: requirement,
        })
        
        currentNode = roomID
        roomID++
    }
    
    // Add optional side branches for exploration
    // ...
}
```

## Common Pitfalls to Avoid

### 1. Breaking Determinism

**Problem**: Using global `rand` or system time during generation

**Why harmful**: Same seed won't produce same game, breaking core feature

**Solution**:
```go
// Good - local RNG from seed
func Generate(seed int64) {
    rng := rand.New(rand.NewSource(seed))
    value := rng.Intn(100)
}

// Bad - global rand is not deterministic
func Generate(seed int64) {
    rand.Seed(seed)
    value := rand.Intn(100)  // Uses global state
}

// Bad - using system time
func Generate(seed int64) {
    value := time.Now().UnixNano() % 100  // Different every time!
}
```

### 2. Shared State Between Generators

**Problem**: One generator modifying state that affects another

**Why harmful**: Subsystems aren't independent, generation order matters

**Solution**: Each generator gets its own derived seed and RNG
```go
// Good - independent generators
seeds := pcg.DeriveSeeds(masterSeed)
graphicsGen := NewGraphicsGenerator(seeds["graphics"])
audioGen := NewAudioGenerator(seeds["audio"])

// Bad - sharing RNG
sharedRNG := rand.New(rand.NewSource(masterSeed))
graphicsGen.UseRNG(sharedRNG)  // Graphics affects audio's randomness!
audioGen.UseRNG(sharedRNG)
```

### 3. Ignoring Thread Safety

**Problem**: Multiple goroutines accessing cache without synchronization

**Why harmful**: Data races, cache corruption, unpredictable behavior

**Solution**: Always use mutex for shared data structures
```go
// Good - thread-safe cache
type AssetCache struct {
    sprites map[string]interface{}
    mu      sync.RWMutex
}

func (c *AssetCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.sprites[key]
    return val, ok
}

// Bad - no synchronization
type AssetCache struct {
    sprites map[string]interface{}
}

func (c *AssetCache) Get(key string) interface{} {
    return c.sprites[key]  // Race condition!
}
```

### 4. Premature Optimization

**Problem**: Adding complex optimization before measuring performance

**Why harmful**: Code complexity increases, bugs hide, gains often minimal

**Example**: Current system generates complete game in ~0.3 seconds. This is fast enough. Don't optimize until you measure a real problem.

### 5. Magic Numbers

**Problem**: Hardcoded values without explanation

**Solution**: Use named constants with comments
```go
// Bad
if temperature < -10 {
    return "freezing"
}

// Good
const FreezingThreshold = -10  // Temperature in Celsius below which player movement slows

if temperature < FreezingThreshold {
    return "freezing"
}
```

### 6. Testing Without Determinism

**Problem**: Tests that pass randomly due to seed variation

**Solution**: Always use fixed seeds in tests
```go
// Good - deterministic test
func TestGeneration(t *testing.T) {
    seed := int64(42)  // Fixed seed
    result := Generate(seed)
    if result != expectedValue {
        t.Error("Generation failed")
    }
}

// Bad - random seed makes test flaky
func TestGeneration(t *testing.T) {
    seed := time.Now().UnixNano()  // Different every run!
    result := Generate(seed)
    // How do you know what to expect?
}
```

## Documentation Requirements

### Public APIs (Exported Functions/Types)

**Requirement**: Every exported symbol must have a doc comment

**Format**: Start with the symbol name, use complete sentences
```go
// SpriteGenerator generates procedural pixel art sprites using cellular automata.
// Sprites can have various symmetry types for visual coherence.
type SpriteGenerator struct {
    Width    int
    Height   int
    Symmetry SymmetryType
}

// NewSpriteGenerator creates a sprite generator with the specified dimensions
// and symmetry type. Default constraints are applied for quality.
func NewSpriteGenerator(width, height int, symmetry SymmetryType) *SpriteGenerator {
    return &SpriteGenerator{
        Width:    width,
        Height:   height,
        Symmetry: symmetry,
    }
}

// Generate creates a deterministic sprite from the given seed.
// The same seed will always produce the same sprite.
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    // Implementation
}
```

### Complex Logic

**When to add inline comments**:
- Non-obvious algorithms (cellular automata rules, music theory)
- Magic numbers that need explanation
- Performance-critical sections
- Workarounds or edge case handling

**Example**:
```go
// Run cellular automata iterations to create organic shapes
for i := 0; i < 3; i++ {
    grid = sg.cellularAutomataStep(grid)
}

// Apply 4-way symmetry by copying top-left quadrant to other quadrants
// This ensures visual balance for enemies and items
for y := 0; y < midY; y++ {
    for x := 0; x < midX; x++ {
        val := grid[y][x]
        grid[y][sg.Width-1-x] = val           // Mirror horizontally
        grid[sg.Height-1-y][x] = val          // Mirror vertically
        grid[sg.Height-1-y][sg.Width-1-x] = val // Mirror both
    }
}
```

### Package Documentation

**Not currently used**: No `doc.go` files exist

**When to add**: If package purpose isn't clear from package name
```go
// Package pcg provides core procedural content generation utilities.
//
// This package implements seed derivation, asset caching, and quality
// validation for deterministic game content generation. All generators
// in other packages depend on the seed management provided here.
package pcg
```

## Project-Specific Conventions

### Seed Management

- Master seed: `int64` from user input or timestamp
- Subsystem seeds: Derived via `HashSeed(masterSeed, "subsystem_name")`
- Local RNG: Always create with `rand.New(rand.NewSource(seed))`
- Never modify `masterSeed` - it's the reproducibility key

### Generator Initialization

Pattern for all generators:
```go
// 1. Define generator struct with config
type Generator struct {
    Config ConfigStruct
}

// 2. Constructor with required params, sensible defaults
func NewGenerator(required int, params string) *Generator {
    return &Generator{
        Config: ConfigStruct{
            Required: required,
            Params:   params,
            Defaults: defaultValue,
        },
    }
}

// 3. Generate method takes seed, returns result
func (g *Generator) Generate(seed int64) *Result {
    rng := rand.New(rand.NewSource(seed))
    // Use g.Config for configuration
    // Use rng for random decisions
    return result
}
```

### Constants for Enums

Use `iota` for sequential enums:
```go
type RoomType int

const (
    CombatRoom RoomType = iota  // 0
    PuzzleRoom                   // 1
    TreasureRoom                 // 2
    CorridorRoom                 // 3
    BossRoom                     // 4
    StartRoom                    // 5
    SaveRoom                     // 6
)
```

### Map Initialization

Always initialize maps before use:
```go
// Good
sprites := make(map[string]*Sprite)
sprites["player"] = playerSprite

// Bad - panic on nil map
var sprites map[string]*Sprite
sprites["player"] = playerSprite  // Panic!
```

## Development Workflow

### Before Submitting Code

- [x] Run `go fmt ./...` - Format all code
- [x] Run `go vet ./...` - Check for common mistakes
- [x] Run `go test ./...` - Ensure all tests pass
- [x] Run `go build ./cmd/game` - Verify builds successfully
- [x] Test with multiple seeds - Verify determinism
- [x] Check determinism - Same seed should produce same output
- [x] Update documentation if adding public APIs
- [x] Add tests for new generators or critical functions

### Running the Game

```bash
# Build
go build -o vania ./cmd/game

# Run with random seed
./vania

# Run with specific seed
./vania --seed 42

# Test generation time
time ./vania --seed 12345
```

### Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/pcg

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Adding a New Generator

1. Create new file in appropriate package: `internal/<package>/<name>_gen.go`
2. Define result struct (e.g., `Sprite`, `AudioSample`, `Biome`)
3. Define generator struct with configuration
4. Implement `New<Generator>` constructor
5. Implement `Generate(seed int64)` method
6. Add tests in `<name>_gen_test.go`
7. Integrate into `engine/game.go` generation pipeline
8. Test with multiple seeds for determinism

### Example: Adding Particle Generator

```go
// internal/graphics/particle_gen.go
package graphics

import "math/rand"

// Particle represents a visual effect particle
type Particle struct {
    X, Y      float64
    VelocityX float64
    VelocityY float64
    Color     color.RGBA
    Lifetime  float64
}

// ParticleGenerator generates particle effects
type ParticleGenerator struct {
    MaxParticles int
    EmitRate     float64
}

// NewParticleGenerator creates a particle generator
func NewParticleGenerator(maxParticles int) *ParticleGenerator {
    return &ParticleGenerator{
        MaxParticles: maxParticles,
        EmitRate:     10.0,
    }
}

// Generate creates a particle system from a seed
func (pg *ParticleGenerator) Generate(seed int64) []Particle {
    rng := rand.New(rand.NewSource(seed))
    
    count := rng.Intn(pg.MaxParticles) + 1
    particles := make([]Particle, count)
    
    for i := range particles {
        particles[i] = Particle{
            X:         rng.Float64() * 100,
            Y:         rng.Float64() * 100,
            VelocityX: rng.Float64()*2 - 1,
            VelocityY: rng.Float64()*2 - 1,
            Color:     color.RGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 255},
            Lifetime:  rng.Float64() * 2.0,
        }
    }
    
    return particles
}
```

## Performance Considerations

### Current Performance
- Complete game generation: ~300ms
- Sprite generation: ~1ms per sprite
- Tileset generation: ~5ms per tileset
- Music generation: ~50ms per track
- World graph generation: ~20ms

### Optimization Guidelines
1. **Profile first**: Use `go test -bench` and `pprof` before optimizing
2. **Cache expensive operations**: Use `AssetCache` for reused content
3. **Prefer simplicity**: Current performance is acceptable for most use cases
4. **Parallel generation**: Future optimization - generate subsystems concurrently
5. **Memory allocation**: Reuse slices when possible with `make([]T, 0, capacity)`

### When to Optimize
- Generation time exceeds 1 second
- Memory usage exceeds 100MB for single game
- Users report noticeable lag
- Benchmarks show regression

## Future Architecture Considerations

### Planned Enhancements

**Save/Load System**:
- Serialize game state to JSON
- Store seed with save for regeneration
- Track player progress and state

**Adaptive Music**:
- Layer system that responds to gameplay
- Transition between biome tracks
- Combat intensity affects music

**Advanced AI**:
- Behavior trees for enemy AI
- Pathfinding with A*
- Context-aware decision making

**Animation System**:
- Sprite frame generation
- Interpolation between frames
- State-based animation

### Architectural Principles to Maintain

1. **Determinism**: Always preserve seed-based reproducibility
2. **Independence**: Keep subsystems loosely coupled
3. **Simplicity**: Prefer clear code over clever optimizations
4. **Testing**: Maintain high test coverage for PCG core
5. **Documentation**: Keep this guide updated with new patterns

## Quick Reference

### Common Type Patterns

```go
// Generator pattern
type XGenerator struct {
    Config ConfigStruct
}

func NewXGenerator(params) *XGenerator { ... }
func (g *XGenerator) Generate(seed int64) *Result { ... }

// Enum pattern
type XType int
const (
    TypeA XType = iota
    TypeB
    TypeC
)

// Result struct pattern
type Result struct {
    Data    []DataType
    Metadata MetaStruct
}
```

### Common Function Signatures

```go
// Generators
func (g *Generator) Generate(seed int64) *Result
func NewGenerator(config) *Generator

// Seed management
func HashSeed(masterSeed int64, identifier string) int64
func DeriveSeeds(masterSeed int64) map[string]int64

// Caching
func (c *Cache) Get(key string) (interface{}, bool)
func (c *Cache) Set(key string, value interface{})

// Validation
func (v *Validator) Validate(content interface{}) bool
```

### Import Organization

Standard library first, then local packages:
```go
import (
    "crypto/sha256"
    "encoding/binary"
    "math/rand"
    
    "github.com/opd-ai/vania/internal/pcg"
    "github.com/opd-ai/vania/internal/graphics"
)
```

---

**Last Updated**: 2024-10-18  
**Maintained by**: VANIA Project Contributors  
**Questions**: See README.md and IMPLEMENTATION.md for more details
