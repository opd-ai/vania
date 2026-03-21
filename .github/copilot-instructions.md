# Project Overview

VANIA is a fully procedurally-generated retro Metroidvania game engine written in pure Go, where **ALL assets** — graphics, audio, narrative, and levels — are generated algorithmically at runtime from a single seed value. The project eliminates traditional asset creation entirely: no pre-rendered images, bundled audio files, or static narrative content exist in the repository. Instead, pixel art sprites are generated via cellular automata with symmetry transforms, sound effects and music through waveform synthesis with ADSR envelopes, and stories through procedural narrative assembly.

The target audience includes developers interested in procedural content generation (PCG) techniques and players seeking infinite unique Metroidvania experiences. Given the same seed, the game produces identical output across all platforms and runs — this determinism is a core architectural constraint, not merely a feature. The single-binary philosophy means the compiled executable contains everything needed to generate complete, playable games without external dependencies or asset files.

The engine supports five genre themes (fantasy, scifi, horror, cyberpunk, postapoc) with a `SetGenre()` interface that propagates thematic changes through rendering, audio, menu, and eventually all systems. Current development focuses on completing v2.0 features including status effects, full inventory UI, and remaining AI enhancements.

## Sibling Repository Context

VANIA is part of the **opd-ai Procedural Game Suite** — eight sibling repositories sharing architectural patterns, code conventions, and the zero-external-assets philosophy. All games are built with Go + Ebiten and generate 100% of content at runtime.

| Repo | Genre | Description |
|------|-------|-------------|
| `opd-ai/venture` | Co-op action-RPG | Top-down roguelike with 35 class system, guilds, trading |
| `opd-ai/vania` | Metroidvania platformer | Side-scrolling ability-gated exploration (this repo) |
| `opd-ai/velocity` | Galaga-like shooter | Space shooter with procedural enemy waves |
| `opd-ai/violence` | Raycasting FPS | First-person shooter with multiplayer, libp2p networking |
| `opd-ai/way` | Battle-cart racer | Racing game with procedural tracks |
| `opd-ai/wyrm` | First-person survival RPG | Survival mechanics with crafting |
| `opd-ai/where` | Wilderness survival | Open-world survival simulation |
| `opd-ai/whack` | Arena battle game | Combat arena with procedural enemies |

Code patterns should remain compatible across all sibling repos to enable future extraction of shared libraries.

## Technical Stack

- **Primary Language**: Go 1.24.9
- **Game Framework**: Ebiten v2.6.3 — 2D game engine with cross-platform + WASM support
- **Audio Backend**: ebitengine/oto v3.1.0 — Low-latency audio output
- **Image Processing**: golang.org/x/image v0.12.0 — Extended image format support
- **Testing**: Go standard `testing` package with table-driven tests and benchmarks
- **Build/Deploy**: `go build -o vania ./cmd/game`, verification via `verify.sh` (includes xvfb for headless CI)

### Key Dependency Versions (from go.mod)
```
github.com/hajimehoshi/ebiten/v2 v2.6.3
github.com/ebitengine/oto/v3 v3.1.0  (indirect)
github.com/ebitengine/purego v0.5.0  (indirect)
golang.org/x/image v0.12.0           (indirect)
golang.org/x/mobile v0.0.0-20230922142353-e2f452493d57 (indirect)
```

## Project Structure

VANIA uses the **vania-style** layout: `cmd/game` entrypoint + `internal/` with domain-specific packages.

```
/cmd/game/              Entry point and CLI (--seed, --play, --genre flags)
/internal/
  ├── pcg/              Core PCG framework (seed derivation, caching, validation)
  ├── graphics/         Sprite generation (cellular automata), tileset, palette
  ├── audio/            Waveform synthesis, ADSR envelopes, SFX, music generation
  ├── narrative/        Story, theme, character, faction generation
  ├── world/            Graph-based room generation, biome system, platforms
  ├── entity/           Enemy, boss, item, ability generation + AI behaviors
  ├── engine/           Game orchestration, combat system, room transitions
  │   └── ecs/          Entity-Component-System framework (not yet integrated)
  ├── render/           Ebiten rendering, bitmap text, HUD
  ├── input/            Keyboard/gamepad mapping with input buffering
  ├── physics/          Gravity, collision, double-jump, dash, glide, grapple
  ├── camera/           Smooth follow, room-lock, screen-shake
  ├── menu/             Main menu, pause, options, settings screens
  ├── save/             Save slots, checkpoint autosave, persistence
  ├── settings/         Configuration persistence (resolution, volume, keys)
  ├── animation/        Frame-based animation state machine
  ├── particle/         Combat hit sparks, movement dust, effects
  └── achievement/      19 achievements across 6 categories with persistence
/docs/                  System documentation (rendering, combat, AI, etc.)
/.github/               CI/CD workflows, this instructions file
```

---

## ⚠️ CRITICAL: Complete Feature Integration (Zero Dangling Features)

**This is the single most important rule for this codebase.** Every feature, system, component, generator, and integration MUST be fully wired into the runtime. Dangling features are a maintenance burden, a source of frustration, and actively degrade code quality.

### The Dangling Feature Problem

In procedural game codebases like VANIA, it is extremely common for features to be:
1. **Defined but never instantiated** — A system struct exists but is never created in `main()` or game initialization
2. **Instantiated but never integrated** — A system runs but its output is never consumed by other systems
3. **Partially integrated** — A system works for one genre/biome but silently no-ops for others
4. **Tested in isolation but broken in context** — Unit tests pass but the system was never wired into the game loop

**VANIA has documented examples of this problem** (see GAPS.md and AUDIT.md):
- ECS framework exists in `internal/engine/ecs/` but `runner.go` doesn't use it
- `SetGenre()` implemented on render/audio/menu but missing on physics/narrative
- Ranged attack mentioned in README but no projectile system exists
- Status effects (burn, freeze, poison) appear in generation but have no gameplay effect

### Mandatory Checks Before Adding or Modifying Any Feature

**Before writing ANY new code, verify the full integration chain:**

1. **Definition → Instantiation**: Is the struct/system created at runtime? Trace from `main()` through `NewGameApp()` → `startGame()` → `NewGameGenerator()` → `GenerateCompleteGame()`.
2. **Instantiation → Registration**: Is the system registered with the game? Check `GameRunner`, `Game` struct fields.
3. **Registration → Update Loop**: Does the system's `Update()` method get called? Check `GameRunner.Update()` in `runner.go`.
4. **Update → Output**: Does the system produce outputs (components, events, state changes) that other systems consume?
5. **Output → Consumer**: Is there at least one other system that reads this system's output?
6. **Consumer → Player Effect**: Does the chain ultimately produce something visible, audible, or mechanically felt by the player?

If ANY link in this chain is missing, the feature is dangling. **Do not submit dangling features.**

### Specific Anti-Patterns to Reject

```go
// ❌ BAD: Generator exists but is never called in runtime code
type StatusEffectManager struct { ... }
func (s *StatusEffectManager) Apply(effect StatusEffect) { ... }
// Only referenced in status_test.go, never in combat.go or runner.go

// ✅ GOOD: Generator created and integrated into game loop
statusMgr := NewStatusEffectManager()
game.StatusManager = statusMgr
// AND in combat.go:
func (cs *CombatSystem) ApplyDamage(target *Entity, damage int) {
    cs.game.StatusManager.Apply(StatusEffect{Type: "burn", Duration: 3.0})
}
```

```go
// ❌ BAD: SetGenre implemented on some systems, missing on others
// render/renderer.go: func (r *Renderer) SetGenre(genreID string) ✅
// audio/player.go:    func (ap *AudioPlayer) SetGenre(genreID string) ✅
// physics/physics.go: NO SetGenre method ❌ — GAPS.md item

// ✅ GOOD: All systems implement SetGenre per ROADMAP requirements
// When adding any system, ALWAYS add SetGenre(genreID string) method
```

```go
// ❌ BAD: ECS System defined but runner.go ignores it
// engine/ecs/system.go defines System interface
// engine/runner.go has 1224 lines of direct logic, no ECS delegation
// AUDIT.md: "ECS framework exists but GameRunner monolith doesn't use it"

// ✅ GOOD: If using ECS, actually delegate to SystemManager
func (gr *GameRunner) Update() error {
    return gr.systemManager.Update(dt)  // Delegate to ECS
}
```

### Known Gaps (from GAPS.md and AUDIT.md)

| Gap | Priority | Status | Location |
|-----|----------|--------|----------|
| Ranged Combat System | CRITICAL | Not implemented | `internal/engine/combat.go` |
| Status Effect System | CRITICAL | Not implemented | `internal/engine/` |
| ECS Integration | HIGH | Framework built, not used | `internal/engine/ecs/` vs `runner.go` |
| SetGenre on Physics | HIGH | Missing | `internal/physics/physics.go` |
| SetGenre on Narrative | HIGH | Missing | `internal/narrative/story_gen.go` |
| Inventory Screen UI | MEDIUM | Data exists, no UI | `internal/menu/` |
| Headless Test Env | MEDIUM | 7 packages fail without X11 | Multiple packages |

**When working on these gaps, ensure the FULL integration chain is completed.**

### Integration Verification Checklist (run before every PR)

```bash
# Every constructor has at least one non-test caller
grep -rn 'func New' --include='*.go' | grep -v _test.go

# All TODOs are tracked in GAPS.md or ROADMAP.md
grep -rn 'TODO\|FIXME\|HACK\|XXX' --include='*.go'

# No empty method bodies in non-test files
grep -rn '{ *}' --include='*.go' | grep -v _test.go | grep 'func'

# SetGenre implemented on all systems (should match System interface implementers)
grep -rn 'SetGenre' --include='*.go' | grep 'func'

# Verification script passes
./verify.sh
```

---

## Deterministic Procedural Generation

### Core Principle: Same Seed = Same Game

All content generation MUST be deterministic and seed-based. Given the same seed, the game MUST produce identical output across all platforms and runs. This is verified by `verify.sh` in CI.

### Seed Derivation Pattern

VANIA uses SHA-256 hash-based seed derivation to ensure subsystem independence:

```go
// internal/pcg/seed.go — THE authoritative pattern
func HashSeed(masterSeed int64, identifier string) int64 {
    h := sha256.New()
    binary.Write(h, binary.LittleEndian, masterSeed)
    h.Write([]byte(identifier))
    sum := h.Sum(nil)
    return int64(binary.LittleEndian.Uint64(sum[:8]))
}

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

### Deterministic RNG Usage

**ALWAYS create local RNG from seed — NEVER use global rand:**

```go
// ✅ GOOD: Local RNG from seed (23+ instances in codebase)
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rng := rand.New(rand.NewSource(seed))
    value := rng.Intn(100)
    // ...
}

// ✅ GOOD: Using pcg helper
rng := pcg.NewDeterministicRNG(seed)

// ❌ BAD: Global rand (non-deterministic, not thread-safe)
value := rand.Intn(100)

// ❌ BAD: Time-based seeding in generation code
rng := rand.New(rand.NewSource(time.Now().UnixNano()))
```

### Seed Propagation Through Generators

When creating sub-generators, derive seeds to maintain determinism hierarchy:

```go
// ✅ GOOD: Derived seeds for sub-generators
func (gg *GameGenerator) GenerateCompleteGame() (*Game, error) {
    seeds := pcg.DeriveSeeds(gg.MasterSeed)
    
    // Each generator gets its own derived seed
    worldGen := world.NewWorldGenerator(15, 10, 100, 5)
    worldGen.SetSeed(seeds["world"])
    
    entityGen := entity.NewEnemyGenerator(seeds["entity"])
    // ...
}

// ❌ BAD: Shared RNG across generators (breaks independence)
sharedRNG := rand.New(rand.NewSource(masterSeed))
worldGen.UseRNG(sharedRNG)   // World affects entity randomness!
entityGen.UseRNG(sharedRNG)
```

---

## Generator Pattern

All procedural generators in VANIA follow a consistent pattern:

```go
// 1. Define generator struct with configuration
type SpriteGenerator struct {
    Width       int
    Height      int
    Symmetry    SymmetryType
    Constraints SpriteConstraints
}

// 2. Constructor with required params and sensible defaults
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

// 3. Generate method takes seed, returns generated content
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rng := rand.New(rand.NewSource(seed))  // Local RNG!
    // Generation logic...
    return sprite
}

// 4. Optional: SetSeed for re-seeding (used by some generators)
func (sg *SpriteGenerator) SetSeed(seed int64) {
    sg.rng = rand.New(rand.NewSource(seed))
}
```

### Generator Registration Pattern

Generators that support multiple variants should be registered in a dispatch map:

```go
// ✅ GOOD: Registry pattern for genre/biome variants
var terrainGenerators = map[string]TerrainGenerator{
    "fantasy":   &FantasyTerrainGen{},
    "scifi":     &ScifiTerrainGen{},
    "horror":    &HorrorTerrainGen{},
    "cyberpunk": &CyberpunkTerrainGen{},
    "postapoc":  &PostapocTerrainGen{},
}

func GenerateTerrain(genre string, seed int64) *Terrain {
    gen, ok := terrainGenerators[genre]
    if !ok {
        gen = terrainGenerators["fantasy"]  // Fallback
    }
    return gen.Generate(seed)
}
```

---

## GenreSwitcher Interface

Every system that affects player-visible presentation must implement `GenreSwitcher`:

```go
// internal/engine/ecs/system.go
type GenreSwitcher interface {
    SetGenre(genreID string)  // "fantasy" | "scifi" | "horror" | "cyberpunk" | "postapoc"
}
```

### Current Implementation Status

| Package | SetGenre Implemented | Notes |
|---------|---------------------|-------|
| `render/renderer.go` | ✅ Yes | Swaps background colors, clears icon cache |
| `audio/player.go` | ✅ Yes | Swaps instrument packs |
| `menu/menu.go` | ✅ Yes | Swaps UI colors |
| `engine/ecs/system_manager.go` | ✅ Yes | Propagates to all registered systems |
| `physics/physics.go` | ❌ **Missing** | GAPS.md item — need genre-specific hazards |
| `narrative/story_gen.go` | ❌ **Missing** | GAPS.md item — need genre vocabulary tables |

### Adding SetGenre to a New System

```go
// When creating any new system, ALWAYS add SetGenre
type MyNewSystem struct {
    genre           string
    genreParameters map[string]GenreConfig
}

func (s *MyNewSystem) SetGenre(genreID string) {
    s.genre = genreID
    // Load genre-specific parameters
    if config, ok := s.genreParameters[genreID]; ok {
        s.applyConfig(config)
    }
}
```

---

## Code Style Guidelines

### Naming Conventions

- **Packages**: lowercase, single-word (`pcg`, `graphics`, `audio`, `world`, `entity`)
- **Files**: snake_case (`sprite_gen.go`, `music_gen.go`, `enemy_gen.go`)
- **Types**: PascalCase (`SpriteGenerator`, `WorldContext`, `EnemyInstance`)
- **Component types**: PascalCase + "Component" suffix if using ECS
- **System types**: PascalCase + "System" suffix (`CombatSystem`, `RenderSystem`)
- **Interfaces**: PascalCase, often `-er` suffix (`Generator`, `GenreSwitcher`)
- **Constants**: PascalCase for exported, use `iota` for enums
- **Seeds**: Always `int64`, always named `seed` in function parameters

### Error Handling

```go
// ✅ GOOD: Return errors for validation failures
func (gg *GameGenerator) GenerateCompleteGame() (*Game, error) {
    if gg.MasterSeed == 0 {
        return nil, fmt.Errorf("master seed cannot be zero")
    }
    // ...
}

// ✅ GOOD: Implicit error handling for deterministic generation
// (no error return needed — generation always succeeds given valid seed)
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    rng := rand.New(rand.NewSource(seed))
    return sprite
}

// ❌ BAD: Panic in library code
func Generate(seed int64) *Sprite {
    if seed == 0 {
        panic("zero seed")  // Never panic in game logic!
    }
}
```

### Documentation

Every exported type and function must have a godoc comment:

```go
// SpriteGenerator generates procedural pixel art sprites using cellular automata.
// Sprites can have various symmetry types for visual coherence.
type SpriteGenerator struct {
    Width    int
    Height   int
    Symmetry SymmetryType
}

// Generate creates a deterministic sprite from the given seed.
// The same seed will always produce the same sprite.
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
    // ...
}
```

---

## Zero External Assets

The single-binary philosophy is **non-negotiable**:

- **Graphics**: Procedurally generated via cellular automata, symmetry transforms, HSV color theory
- **Audio**: Synthesized via additive synthesis, ADSR envelopes, chord progressions
- **Levels**: Generated via graph theory, ability-gating algorithms, biome assignment
- **Narrative**: Assembled via procedural story generation, constraint-based selection
- **UI**: Built from code using bitmap text rendering, not loaded from image files

**Never add asset files to the repository:**
- No `.png`, `.jpg`, `.svg`, `.gif` images
- No `.mp3`, `.wav`, `.ogg` audio files
- No `.json`, `.yaml` level definition files
- No pre-written dialogue, story scripts, or text assets

If you need test fixtures, generate them in test setup code.

---

## Testing Standards

### Coverage Targets

| Package | Target | Current |
|---------|--------|---------|
| `pcg` | ≥82% | 96.3% ✅ |
| `graphics` | ≥82% | 99.1% ✅ |
| `achievement` | ≥82% | 90.9% ✅ |
| `audio` | ≥82% | 68.1% ⚠️ (GAPS.md item) |
| Display-dependent | ≥30% | 7 packages fail without X11 |

### Table-Driven Tests

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

### Determinism Tests (Critical)

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
            if c1 != c2 {
                t.Errorf("Sprites not deterministic at (%d,%d)", x, y)
            }
        }
    }
}
```

### Headless Testing

7 packages require X11 display. Use `xvfb-run` for CI:

```bash
# Run tests with virtual display
xvfb-run -a go test ./...

# Or use DISPLAY= to identify headless-incompatible tests
DISPLAY= go test ./... 2>&1 | grep -c "^ok"
```

---

## Networking Best Practices (Future v5.0)

VANIA does not currently have networking code, but v5.0 ROADMAP specifies multiplayer co-op. When implementing:

### Interface-Only Network Types (Hard Constraint)

**ALWAYS use interface types for network variables:**

| ❌ Never Use | ✅ Always Use |
|-------------|--------------|
| `*net.UDPAddr` | `net.Addr` |
| `*net.UDPConn` | `net.PacketConn` |
| `*net.TCPConn` | `net.Conn` |
| `*net.TCPListener` | `net.Listener` |

```go
// ✅ GOOD
var addr net.Addr
var conn net.PacketConn

// ❌ BAD
var addr *net.UDPAddr
var conn *net.UDPConn
```

### High-Latency Design (200–5000ms)

Per ROADMAP v5.0, multiplayer must tolerate 200–5000ms latency:

1. **Client-Side Prediction**: Never block game loop waiting for server
2. **State Interpolation**: Interpolate between known server states
3. **Jitter Buffers**: Buffer incoming updates, absorb ±500ms variance
4. **Idempotent Messages**: Safe to process multiple times
5. **No Synchronous RPC**: All network I/O must be async
6. **Timeout Tolerance**: ≥10 second timeouts, ≥3 missed heartbeats for disconnect

---

## Cross-Repository Code Sharing

### Shared Pattern Catalog

When implementing features, follow these patterns for future extraction:

| Pattern | VANIA Package | Sibling Convention |
|---------|---------------|-------------------|
| Seed derivation | `internal/pcg/seed.go` | `pkg/seed/` or inline |
| ECS framework | `internal/engine/ecs/` | `pkg/engine/ecs/` |
| Sprite generation | `internal/graphics/` | `pkg/rendering/` |
| Audio synthesis | `internal/audio/` | `pkg/audio/` |
| Camera system | `internal/camera/` | `pkg/camera/` |
| Input handling | `internal/input/` | `pkg/input/` |
| Particle system | `internal/particle/` | `pkg/particles/` |
| Save/load | `internal/save/` | `pkg/saveload/` |
| Achievement | `internal/achievement/` | `pkg/achievement/` |
| Menu/UI | `internal/menu/` | `pkg/rendering/ui/` |

### Interface Consistency Across Repos

Maintain identical interfaces for future shared packages:

```go
// ECS System interface (must match across all repos)
type System interface {
    Update(dt float64) error
    Draw(screen *ebiten.Image)
    SetGenre(genreID string)
}

// Component identifier (must match across all repos)
type Component interface {
    Type() string
}
```

---

## Performance Requirements

- **Target**: 60 FPS on mid-range hardware
- **Memory budget**: <500MB client
- **Generation time**: ~300ms for complete game (currently achieved)
- **Sprite generation**: ~1ms per sprite
- **Music generation**: ~50ms per track

### Optimization Guidelines

1. **Profile first**: Use `go test -bench=. -benchmem` before optimizing
2. **Cache expensive operations**: Use `AssetCache` in `internal/pcg/cache.go`
3. **Object pooling**: For frequently allocated objects (particles, projectiles)
4. **Spatial partitioning**: Use for entity queries over collections >100
5. **Don't regenerate**: Cache generated sprites, never regenerate same sprite twice

---

## Development Workflow

### Before Submitting Code

```bash
# Format and vet
go fmt ./...
go vet ./...

# Run all tests (with virtual display for CI)
xvfb-run -a go test ./...

# Build
go build -o vania ./cmd/game

# Verify determinism
./vania --seed 42  # Run twice, compare output

# Full verification
./verify.sh
```

### Running the Game

```bash
# Random seed
./vania

# Specific seed
./vania --seed 42

# Play mode (with rendering)
./vania --seed 42 --play

# Specific genre
./vania --seed 42 --play --genre scifi
```

### Adding a New Generator

1. Create file: `internal/<package>/<name>_gen.go`
2. Define result struct and generator struct
3. Implement `New<Generator>()` constructor
4. Implement `Generate(seed int64)` method
5. **Add `SetGenre(genreID string)` if system affects presentation**
6. Create tests: `<name>_gen_test.go`
7. **Integrate into `GameGenerator.GenerateCompleteGame()`**
8. **Verify full integration chain from main() to player effect**
9. Update ROADMAP.md if completing a milestone item

---

## GAPS.md and AUDIT.md Protocol

These files track implementation gaps and audit findings. When you identify a gap:

1. Note it in your response with severity (Critical/High/Medium/Low)
2. Include file path and line number
3. Propose an actionable fix with validation command
4. Suggest adding to GAPS.md if not already tracked

Example GAPS.md entry format:
```markdown
## [Feature Name]

- **Stated Goal**: [What README/ROADMAP claims]
- **Current State**: [What actually exists]
- **Impact**: [Why this matters]
- **Closing the Gap**:
  1. [Step 1]
  2. [Step 2]
  3. Validation: `[command to verify]`
```

---

## Quick Reference

### Common Constructors

```go
// PCG
pcg.NewPCGContext(seed)
pcg.NewDeterministicRNG(seed)
pcg.HashSeed(masterSeed, identifier)
pcg.DeriveSeeds(masterSeed)

// Graphics
graphics.NewSpriteGenerator(width, height, symmetry)
graphics.NewTilesetGenerator(tileSize, biome)
graphics.NewPaletteGenerator(scheme)

// Audio
audio.NewSFXGenerator(sampleRate)
audio.NewMusicGenerator(sampleRate, bpm, rootNote, scale)
audio.NewAudioPlayer()

// World
world.NewWorldGenerator(width, height, roomCount, biomeCount)
world.NewBiomeGenerator()
world.NewPlatformGenerator()

// Entity
entity.NewEnemyGenerator(seed)
entity.NewBossGenerator(seed)
entity.NewItemGenerator(seed)
entity.NewAbilityGenerator(seed)

// Engine
engine.NewGameGenerator(masterSeed)
engine.NewGameGeneratorWithGenre(masterSeed, genre)
engine.NewCombatSystem()
engine.NewRoomTransitionHandler(game)
```

### CLI Flags

```
--seed <int64>   Specify generation seed (default: timestamp)
--play           Launch full game with rendering
--genre <string> Set genre: fantasy|scifi|horror|cyberpunk|postapoc
```

### Genre IDs

```go
const (
    GenreFantasy   = "fantasy"   // Enchanted castle, vine-covered towers
    GenreScifi     = "scifi"     // Derelict space hulk, exposed conduits
    GenreHorror    = "horror"    // Haunted mansion, candlelit halls
    GenreCyberpunk = "cyberpunk" // Megastructure, neon corridors
    GenrePostapoc  = "postapoc"  // Collapsed bunker, concrete debris
)
```

---

**Last Updated**: 2026-03-21
**Maintainer**: opd-ai Project Contributors
**Related Docs**: README.md, GAPS.md, AUDIT.md, ROADMAP.md, PLAN.md
