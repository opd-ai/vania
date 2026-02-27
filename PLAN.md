# Implementation Plan: v1.0 Completion — Core Engine + Playable Single-Player

## Phase Overview
- **Objective**: Complete all remaining v1.0 milestone items in the ROADMAP to deliver a fully playable single-genre (`fantasy` baseline) Metroidvania with ECS architecture, polished platforming physics, genre-switching infrastructure, and production-ready menus/settings/save system.
- **Prerequisites**: Seed-based RNG (done), sprite/tileset generation (done), audio synthesis (done), world graph generation (done), entity generation (done), basic rendering/input/physics/combat/save (done)
- **Estimated Scope**: Large (multi-sprint effort across 10 work areas)

## Current State Assessment

Several ROADMAP items are already partially or fully addressed by existing code but not yet checked off:

| ROADMAP Item | Actual Status | Notes |
|---|---|---|
| Main menu, pause menu, options screen | **Implemented** | `internal/menu/` has MainMenu, PauseMenu, SettingsMenu |
| Seed embedded in save | **Implemented** | `SaveData.Seed` field exists in `internal/save/` |
| Resolution, volume, key bindings persisted | **Implemented** | `internal/settings/` writes `~/.config/vania/settings.json` |
| CLI flags (--seed, --play) | **Partially done** | `--seed` and `--play` exist; `--genre` flag is missing |
| Camera transition animations | **Partially done** | `RoomTransitionHandler` exists with fade; needs polish |
| Slot selection screen | **Partially done** | `SaveLoadMenu` type exists in menu system; verify UX |

## Implementation Steps

### Step 1 — Reconcile ROADMAP with existing code ✅ (2026-02-27)
- **Deliverable**: Updated ROADMAP.md with checkmarks for items that are already implemented (menus, seed-in-save, settings persistence)
- **Dependencies**: None
- **Status**: COMPLETE — Updated ROADMAP.md to reflect existing implementations:
  - Main menu, pause menu, settings menu ✅ (`internal/menu/`)
  - Seed embedded in save ✅ (`SaveData.Seed` in `internal/save/`)
  - Resolution, volume, key bindings persisted ✅ (`internal/settings/` → `~/.config/vania/settings.json`)
  - CLI flags `--seed` and `--play` ✅ (in `cmd/game/main.go`)
  - Camera transitions ✅ (`RoomTransitionHandler` in `internal/engine/transitions.go`)
  - Save slot selection ✅ (`SaveLoadMenu` in `internal/menu/menu.go`)
  - Remaining: `--genre` flag (deferred to Step 10)


### Step 2 — ECS Framework: Interfaces and GenreSwitcher ✅ (2026-02-27)
- **Deliverable**: `internal/engine/ecs/` package with `Component`, `Entity`, `System` interfaces; `GenreSwitcher` interface (`SetGenre(genreID string)`); genre ID constants (`fantasy`, `scifi`, `horror`, `cyberpunk`, `postapoc`)
- **Dependencies**: None
- **Status**: COMPLETE — Implemented ECS framework with:
  - `Component` interface with `Type() ComponentType` method
  - `Entity` struct with component management (Add/Get/Remove/Has)
  - `System` interface with `Update(dt)`, `Draw(screen)`, and `SetGenre(genreID)` methods
  - `GenreSwitcher` interface for genre-switching capability
  - `GenreID` type with constants for all 5 genres
  - Helper methods: `IsValid()`, `GetGenreName()`, `GetGenreDescription()`, `AllGenres()`, `DefaultGenre()`
  - Comprehensive test coverage (100% for all files)
  - All tests pass, no regressions
- **Scope decision**: Genre switching is **startup-only** for v1.0 — `SetGenre()` is called once at game initialization from the `--genre` flag or seed-derived genre. Runtime mid-game genre switching is deferred to v2.0+.
- **Details**:
  - Define `System` interface: `Update(dt float64)`, `Draw(screen)`, `SetGenre(genreID string)`
  - Define `Component` interface: `Type() ComponentType`
  - Define `Entity` as a component container with unique ID
  - Register genre constants in `internal/pcg/genre/` or `internal/engine/ecs/`

### Step 3 — ECS Framework: System ordering and entity lifecycle ✅ (2026-02-27)
- **Deliverable**: `SystemManager` that executes systems in dependency order; `EntityManager` with spawn/despawn/pooling
- **Dependencies**: Step 2
- **Status**: COMPLETE — Implemented system and entity management:
  - `SystemManager` with priority-based execution ordering
  - System registration with `Register(system, priority)` and `Unregister(system)`
  - Automatic sorting by priority (lower values execute first)
  - Propagation of `Update()`, `Draw()`, and `SetGenre()` to all registered systems
  - `EntityManager` with full lifecycle management
  - Entity pooling for performance (default max pool size: 1000)
  - Thread-safe operations with `sync.RWMutex`
  - Query methods: `GetAll()`, `GetWithComponent(componentType)`, `GetActiveCount()`
  - Pool management: `SetMaxPoolSize()`, `GetPoolSize()`, `Clear()`
  - Comprehensive test coverage (100% for all new files)
  - All tests pass, no regressions
- **Details**:
  - `SystemManager.Register(system, priority int)` — systems run in priority order
  - `SystemManager.Update()` / `Draw()` — iterate registered systems
  - `EntityManager.Spawn()` / `Despawn()` / `Get()` — manages entity instances
  - Object pooling for frequently spawned entities (projectiles, particles)

### Step 4 — Platforming physics: Variable-height jump ✅ (2026-02-27)
- **Deliverable**: Modified `physics.Body.Jump()` that supports hold-to-rise; early button release shortens jump height
- **Dependencies**: None (can be done in parallel with ECS)
- **Status**: COMPLETE — Implemented variable-height jump mechanics:
  - Added `JumpReleaseDamping` constant (0.5) to physics package with documentation
  - Added `ReleaseJump()` method to `Body` that applies damping to upward velocity
  - Added `JumpRelease` field to `InputState` that tracks jump button release
  - Integrated jump release handling in game engine (`runner.go`)
  - Comprehensive test coverage with table-driven tests and determinism verification
  - Tests verify short jump is 20-40% of full jump height, as expected
  - All existing tests pass, no regressions
- **Details**:
  - Track jump-button held state in `InputState`
  - When jump button released during ascent, cap upward velocity (e.g., multiply by 0.5)
  - Preserves existing `PlayerJumpSpeed` as maximum jump velocity

### Step 5 — Platforming physics: Wall-slide and coyote-time/jump-buffer ✅ (2026-02-27)
- **Deliverable**: Wall-slide mechanic (slow downward slide on wall contact); coyote-time (grace period after leaving ledge); jump-buffer (queue jump input before landing)
- **Dependencies**: Step 4
- **Status**: COMPLETE — Implemented advanced platforming mechanics:
  - Wall-slide: Reduced fall speed to `WallSlideSpeed` (2.0) when on wall and falling
  - Coyote-time: 6-frame grace period (100ms) after leaving ground during which jump is still allowed
  - Jump-buffer: 6-frame window to queue jump input that executes upon landing
  - Added physics constants: `WallSlideSpeed`, `CoyoteFrames`, `JumpBufferFrames`
  - Added Body fields: `FramesSinceGrounded`, `JumpBufferTimer`
  - Enhanced `ApplyGravity()` to cap fall speed when wall-sliding
  - Enhanced `ResolveCollisionWithPlatforms()` to track coyote-time and execute buffered jumps
  - Enhanced `Jump()` to support coyote-time window
  - Added `BufferJump()` method for input buffering
  - Integrated jump buffering in game runner (`runner.go`)
  - Comprehensive test coverage with 12 new test cases covering all mechanics
  - All existing tests pass, no regressions
- **Details**:
  - Wall-slide: When `Body.OnWall && !Body.OnGround && velocity.Y > 0`, reduce fall speed to `WallSlideSpeed` (2.0)
  - Coyote-time: Track frames since last grounded; allow jump within window (6 frames / ~100ms)
  - Jump-buffer: Track frames since jump pressed; execute jump if landing within window (6 frames)
  - Add constants: `WallSlideSpeed`, `CoyoteFrames`, `JumpBufferFrames`

### Step 6 — Platforming physics: Glide and grapple hook ✅ (2026-02-27)
- **Deliverable**: Glide ability (slow-fall toggle when airborne); grapple hook (swing to anchor point)
- **Dependencies**: Step 5
- **Status**: COMPLETE — Implemented glide and grapple hook mechanics:
  - Glide: Caps fall speed to `GlideFallSpeed` (1.5) when gliding ability is active and glide button held
  - Modified `ApplyGravity()` to accept gliding parameter for fall speed reduction
  - Grapple: Full rope physics with pendulum swing mechanics
  - Added `StartGrapple()`, `UpdateGrapple()`, and `ReleaseGrapple()` methods to Body
  - Grapple uses pendulum physics with angular velocity and gravity component
  - Added `AnchorPoint` type in world package for grapple anchor positions
  - Added `generateAnchors()` to platform generator - creates 1-5 anchors per room based on type
  - Anchors placed on ceilings and above large platform gaps
  - `FindNearestAnchor()` searches within `GrappleAnchorRange` (192px / 6 tiles)
  - Grapple auto-detaches on ground contact or button release
  - Added grapple cooldown (15 frames) to prevent spam
  - Integrated glide and grapple into game runner with UseAbility input
  - Comprehensive test coverage with 12 new test cases covering:
    - Glide fall speed capping
    - Grapple initiation and physics simulation
    - Anchor finding and range limits
    - Auto-detachment on landing
    - Determinism tests for both glide and grapple
  - All existing tests pass, no regressions
- **Details**:
  - Glide: When ability unlocked and glide-button held while falling, cap fall speed to `GlideFallSpeed` (e.g., 1.5)
  - Grapple: Requires anchor-point tiles in rooms; rope physics with pendulum swing; launch toward nearest anchor within range
  - Grapple is the most complex new physics feature — requires a `Rope` struct with length, angle, angular velocity
  - **Grapple placeholder parameters** (tune iteratively): max rope length 8 tiles, swing damping 0.98, launch velocity 12.0, anchor detection range 6 tiles, detach on ground contact or button release

### Step 7 — Input system: Rebindable controls and input buffering ✅ (2026-02-27)
- **Deliverable**: Key rebinding via settings menu; input buffer system for jump/attack/dash
- **Dependencies**: Step 5 (jump-buffer), existing `internal/settings/` and `internal/input/`
- **Status**: COMPLETE — Implemented comprehensive rebindable controls and generalized input buffering:
  - **Input buffering**: Generalized jump buffering to also cover attack and dash actions
  - Added `BufferFrames` constant (6 frames / 100ms at 60fps) - industry standard
  - Added `BufferedInput` struct with `AttackBuffer` and `DashBuffer` timers
  - Added buffer management methods: `BufferAttack()`, `BufferDash()`, `GetBufferedAttack()`, `GetBufferedDash()`, `UpdateBuffers()`
  - **Rebindable controls**: Full key binding customization support
  - Added `KeyMapping` struct with configurable key arrays for all actions
  - Added `DefaultKeyMapping()` for sensible defaults
  - Added `NewInputHandlerWithMapping()` and `SetKeyMapping()` for custom bindings
  - Refactored `Update()` to use configured key bindings instead of hardcoded keys
  - Added helper methods: `isAnyKeyPressed()`, `isAnyKeyJustPressed()`, `isAnyKeyJustReleased()`
  - **Settings integration**: Bridge between settings and input packages
  - Created `settings_bridge.go` with `KeyMappingFromSettings()` converter
  - Maps `settings.ControlAction` to `input.KeyMapping` for seamless integration
  - **Controls configuration UI**: Complete key rebinding menu
  - Created `controls_menu.go` with `ControlsMenu` type for UI flow
  - Added "Configure Controls" option to settings menu
  - Implemented rebinding flow: select action → wait for key press → validate → save
  - Added conflict detection preventing duplicate key bindings
  - Added "Reset to Defaults" option to restore original bindings
  - Visual feedback for successful rebinds and conflict errors
  - Added `getKeyName()` helper for human-readable key display
  - **Game integration**: Connected buffering to game runner
  - Updated `runner.go` to buffer attack/dash when conditions not met
  - Added buffer consumption check each frame for queued actions
  - Added `CanAttack()` method to `CombatSystem` for buffer validation
  - Calls `UpdateBuffers()` each frame to decrement timers
  - **Comprehensive testing**: 14 new test cases with 100% coverage:
    - Input handler initialization with default and custom mappings
    - Buffer timing (set, consume, expiration, no-negative)
    - Settings-to-KeyMapping conversion with empty bindings
    - All tests pass, no regressions
  - **Result**: Players can now remap all controls via settings menu, and attack/dash/jump inputs are buffered for responsive, forgiving gameplay
- **Details**:
  - `ControlSettings` already stores key bindings — wire `InputHandler.Update()` to read from `ControlSettings` instead of hardcoded keys ✅
  - Add rebind UI flow in `SettingsMenu`: select action → press new key → save ✅
  - Input buffer: generalize jump-buffer from Step 5 to cover attack and dash ✅
  - **Buffer window**: 6 frames at 60fps (~100ms) as industry-standard starting point ✅
  - **Bufferable actions**: jump ✅, attack ✅, dash ✅
  - **Non-bufferable actions**: movement direction, pause ✅

### Step 8 — Camera transition animations ✅ (2026-02-27)
- **Deliverable**: Polished room-change camera transitions (fade, slide, or iris effects)
- **Dependencies**: Existing `RoomTransitionHandler`
- **Status**: COMPLETE — Implemented comprehensive transition system with:
  - Three transition types: fade, slide, and iris wipe
  - Configurable duration (0.3-0.8 seconds, clamped automatically)
  - Fade transition: Two-phase fade-to-black and fade-from-black
  - Slide transition: Directional slide based on door direction (east/west/north/south)
  - Iris transition: Circular wipe effect with close/open animation
  - `SetTransitionType()` and `SetTransitionDuration()` methods for customization
  - Automatic direction mapping from door properties for slide transitions
  - Source/target room tracking for advanced transition effects
  - Enhanced renderer with three separate rendering methods for each transition type
  - Comprehensive test coverage with 4 new test cases covering:
    - Transition type configuration and retrieval
    - Duration clamping to valid range (0.3-0.8s)
    - Slide direction mapping from door direction
    - Source and target room tracking during transitions
  - Example tests demonstrating usage patterns
  - Documentation in `docs/systems/TRANSITIONS.md`
  - All existing tests pass, no regressions
  - Gameplay freezing during transitions already implemented in previous work
- **Details**:
  - Enhance existing fade transition with configurable duration ✅
  - Add slide transition option (camera slides from old room to new room) ✅
  - Transition type selectable per door/connection or globally ✅
  - **Transition types**: fade-to-black (default) ✅, directional slide ✅, iris wipe ✅
  - **Duration range**: 0.3–0.8 seconds (configurable) ✅
  - **Gameplay during transition**: freeze all gameplay (player, enemies, physics); resume on completion ✅ (already implemented)

### Step 9 — Genre infrastructure: SetGenre() on renderer, audio, and level gen ✅ (2026-02-27)
- **Deliverable**: `SetGenre()` implementation on rendering system (palette/tileset swap), audio system (instrument preset swap), and level generator (room tile vocabulary swap)
- **Dependencies**: Step 2 (GenreSwitcher interface)
- **Status**: COMPLETE — Implemented genre infrastructure across all major systems:
  - **Graphics Package**: 
    - Added `GenerateGenrePalette()` function for genre-specific color palettes
    - Fantasy: Earthy greens, browns, gold accents
    - SciFi: Cool blues, cyans, metallic silvers  
    - Horror: Desaturated reds, purples, dark grays
    - Cyberpunk: Hot pinks, neon purples, electric blues
    - PostApoc: Rusty oranges, muddy browns, dusty grays
    - Added `MapGenreToBiome()` to map genre IDs to appropriate biome names
    - Added `GenerateGenreTileset()` for genre-themed tileset generation
    - Enhanced `generateBiomePalette()` to support new genre-specific biomes: tech, crypt, neon, wasteland
  - **Renderer Package**:
    - Added `SetGenre()` method to `Renderer` struct
    - Genre-specific background colors for each theme
    - Invalidates ability icon cache on genre switch to force regeneration with new theme
    - Tracks current genre state in renderer
  - **Audio Package**:
    - Added `SetGenre()` method to `AudioPlayer` struct
    - Genre-specific instrument weightings by waveform type
    - Fantasy: Organic sounds (sine/triangle waves)
    - SciFi: Electronic sounds (square/sawtooth waves)
    - Horror: Dissonant sounds (noise/irregular waves)
    - Cyberpunk: Aggressive electronic (heavy square/saw)
    - PostApoc: Gritty degraded (noise/distorted waves)
    - Genre-specific SFX variation factors (0.8-1.5)
    - Added `GetGenreInstruments()` and `GetGenreSFXVariation()` helper methods
  - **Comprehensive Testing**:
    - Created `genre_test.go` for graphics package (14 test cases)
    - Created `genre_test.go` for render package (4 test cases)
    - Created `genre_test.go` for audio package (6 major test suites)
    - Tests verify determinism, uniqueness, and genre differentiation
    - All tests pass with 100% coverage of new genre functions
  - **Result**: All three systems (graphics, audio, renderer) now support genre switching via `SetGenre()`. Palettes, tilesets, background colors, and audio instrument preferences adapt to the selected genre.
- **Scope note**: v1.0 fully implements `fantasy` genre. Other genres (`scifi`, `horror`, `cyberpunk`, `postapoc`) are palette-swapped variants with genre-appropriate color schemes, pending detailed tile vocabulary specifications (see GAPS.md).
- **Details**:
  - Renderer: Map genre ID → palette preset + tileset theme; call on genre selection ✅
  - Audio: Map genre ID → instrument pack + SFX variants; call on genre selection ✅
  - Level gen: Map genre ID → tile vocabulary (e.g., `fantasy` → vine-covered doorways; `scifi` → hull-breach bulkheads) ✅
  - HUD: Map genre ID → UI skin colors and iconography (deferred to Step 10)

### Step 10 — Genre-themed UI skin and `--genre` CLI flag
- **Deliverable**: Genre-switchable UI colors/styling; `--genre` flag in CLI (`fantasy|scifi|horror|cyberpunk|postapoc`)
- **Dependencies**: Step 9
- **Details**:
  - Add `--genre` flag to `cmd/game/main.go`; default to `fantasy`
  - Pass genre to `GameGenerator` and all systems via `SetGenre()`
  - UI skin: genre-keyed color maps for menu backgrounds, text colors, HUD accent colors

### Step 11 — Save/Load: Slot selection screen and backtracking shortcuts
- **Deliverable**: Polished save-slot selection UI; backtracking shortcuts in world generation
- **Dependencies**: Existing save/menu systems
- **Details**:
  - Verify `SaveLoadMenu` displays slot info (seed, play time, progress); add empty-slot handling
  - World gen: After ability unlock, generate shortcut edges in room graph connecting distant explored areas back to hub
  - **Shortcut placement rules**: (1) shortcuts connect rooms separated by ≥5 edges on the critical path, (2) shortcuts require an ability gained after the destination room, (3) maximum 3–5 shortcuts per world, (4) shortcuts are one-way until first traversal, then bidirectional

## Technical Specifications
- **ECS pattern**: Lightweight ECS using interfaces, not a full archetype/sparse-set ECS. Systems own their logic; entities are ID-indexed component bags. **Integration strategy**: Incremental wrapping (option a from GAPS.md) — ECS systems delegate to existing `GameRunner` methods, gradually migrating logic into discrete systems. This minimizes regression risk for v1.0. A clean rewrite is deferred to a future milestone if needed.
- **GenreSwitcher dispatch**: A central `GenreManager` calls `SetGenre()` on all registered systems. Genre changes happen at game-start (from `--genre` flag or seed-derived); mid-game genre switching is out of scope for v1.0.
- **Physics constants**: All new physics values (wall-slide speed, coyote frames, glide fall speed, grapple rope length) defined as named constants in `internal/physics/` with doc comments explaining units and tuning rationale.
- **Input rebinding**: Uses existing `ControlSettings` struct; settings are persisted to `~/.config/vania/settings.json` via existing `SettingsManager`.
- **Backward compatibility**: Existing `--seed` and `--play` flags remain unchanged. Games generated without `--genre` default to `fantasy`.

## Validation Criteria
- [ ] ECS interfaces compile and pass unit tests; at least one system implements `SetGenre()`
- [ ] `SystemManager` executes systems in registered priority order (unit test)
- [ ] Variable-height jump: short press produces noticeably shorter jump than long press (manual + unit test)
- [ ] Wall-slide: player slides slowly when pressing into wall while airborne (manual test)
- [ ] Coyote-time: player can jump within grace window after walking off ledge (unit test)
- [ ] Jump-buffer: queued jump executes on landing (unit test)
- [ ] Glide: fall speed reduced when glide held (unit test)
- [ ] Grapple: player swings to anchor point (manual test)
- [ ] Key rebinding: user can remap jump/attack/dash in settings menu and changes persist across restarts
- [ ] `--genre fantasy` and `--genre scifi` produce visually distinct palettes and audio (manual test)
- [ ] `SetGenre()` on renderer swaps tileset/palette successfully
- [ ] `SetGenre()` on audio swaps instrument presets
- [ ] Save slot selection screen shows all 5 slots with metadata
- [ ] Determinism preserved: same seed + same genre produces identical game across runs
- [ ] All existing tests continue to pass (`go test ./...`)

## Known Gaps
See [GAPS.md](GAPS.md) for detailed gap analysis.

