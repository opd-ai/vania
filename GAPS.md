# Implementation Gaps — 2026-03-21

## Ranged Combat System

- **Stated Goal**: README claims "Ranged attack (projectile with falloff)" as a core combat feature
- **Current State**: Only melee combat exists in `internal/engine/combat.go`. No projectile entity, trajectory calculation, or ranged damage logic implemented.
- **Impact**: Players cannot engage enemies at range, limiting combat tactical options and making certain enemy types (ranged enemies mentioned in AI) asymmetric — enemies can shoot but player cannot.
- **Closing the Gap**: 
  1. Create `Projectile` struct in `internal/entity/projectile.go` with position, velocity, damage, falloff parameters
  2. Add `PlayerRangedAttack()` to `combat.go` that spawns projectiles
  3. Add projectile update/collision logic to `runner.go` game loop
  4. Implement damage falloff based on travel distance
  5. Test with `go test ./internal/engine -v -run TestRangedAttack`

---

## Status Effect System

- **Stated Goal**: ROADMAP v2.0 specifies "Burn, freeze, shock, poison, bleed, slow, haste" with stack/duration management and genre-mapped variants
- **Current State**: No status effect system exists. String labels like "poison" appear in procedural generation but provide no gameplay effect.
- **Impact**: Combat lacks depth — no damage-over-time, no crowd control, no buff/debuff strategy. Genre theming cannot extend to combat effects.
- **Closing the Gap**:
  1. Create `internal/engine/status.go` with:
     - `StatusEffect` struct (type, stacks, duration, source, tickDamage)
     - `StatusManager` with `Apply()`, `Tick()`, `Remove()`, `HasEffect()` methods
  2. Define 7 effect types: Burn, Freeze, Shock, Poison, Bleed, Slow, Haste
  3. Integrate with combat damage flow — effects applied on hit, ticked in update loop
  4. Add `SetGenre()` for effect name/color mapping (e.g., "Irradiate" for postapoc)
  5. Render status icons in HUD
  6. Test with `go test ./internal/engine -v -run TestStatusEffects`

---

## ECS Integration with GameRunner

- **Stated Goal**: ROADMAP v1.0 specifies "Component / Entity / System interfaces" with ECS framework powering the engine
- **Current State**: ECS framework exists in `internal/engine/ecs/` with complete System, Entity, Component interfaces. However, `internal/engine/runner.go` is a 1224-line monolith that directly manages all game logic without using ECS.
- **Impact**: Adding new features requires modifying the monolithic runner. Systems cannot be composed, reordered, or disabled independently. Genre switching cannot propagate through a unified system manager.
- **Closing the Gap**:
  1. Create ECS-compatible system wrappers: `RenderSystem`, `PhysicsSystem`, `CombatSystem`, `AudioSystem`, `AISystem`
  2. Each system implements `ecs.System` interface with `Update(dt float64)` and `SetGenre(genreID string)`
  3. Modify `GameRunner` to hold `*ecs.SystemManager` and delegate updates
  4. Migrate logic incrementally — start with `AudioSystem` (least coupled)
  5. Validate with `go test ./internal/engine/ecs -v`

---

## GenreSwitcher on Physics System

- **Stated Goal**: ROADMAP requires `SetGenre()` on all systems including physics for genre-specific hazards: "magic barriers, airlock vents, darkness, voltage floors, radiation clouds"
- **Current State**: `internal/physics/physics.go` has no `SetGenre()` method and no genre-aware hazard system.
- **Impact**: Physics hazards are not genre-themed. Radiation clouds in postapoc biome behave identically to magic barriers in fantasy — no visual or mechanical differentiation.
- **Closing the Gap**:
  1. Add `currentGenre string` field to physics package
  2. Implement `SetGenre(genreID string)` that configures:
     - Hazard damage types and amounts
     - Hazard visual parameters (passed to renderer)
     - Genre-specific collision behaviors
  3. Create hazard configuration per genre in `internal/physics/hazards.go`
  4. Test with `go test ./internal/physics -v -run TestGenreHazards`

---

## GenreSwitcher on Narrative System

- **Stated Goal**: ROADMAP requires `SetGenre()` for "re-skinning narrative vocabulary" — same story structure with genre-appropriate language
- **Current State**: `internal/narrative/` contains stateless generator functions. No `SetGenre()` method, no vocabulary tables per genre.
- **Impact**: Generated narrative text uses fixed vocabulary regardless of genre. Fantasy game may mention "data terminals" or scifi game may reference "enchanted swords" due to vocabulary bleed.
- **Closing the Gap**:
  1. Create `NarrativeGenerator` struct with `genre string` field
  2. Build vocabulary tables per genre for: locations, characters, items, verbs, adjectives
  3. Implement `SetGenre(genreID string)` that selects vocabulary table
  4. Modify all generation functions to use genre-appropriate vocabulary
  5. Test with `go test ./internal/narrative -v -run TestGenreVocabulary`

---

## Inventory Screen UI

- **Stated Goal**: ROADMAP v2.0 specifies "Inventory screen (grid layout)" with consumables, key items, and equipment slots
- **Current State**: Data structure exists (`game.go:47` has `Inventory []*entity.Item`) but no UI renders it. No grid layout, no item selection, no consumable usage.
- **Impact**: Players collect items but cannot view, organize, or use them. Item collection is effectively non-functional from user perspective.
- **Closing the Gap**:
  1. Create `internal/menu/inventory.go` with:
     - `InventoryScreen` struct with grid dimensions (4×8 recommended)
     - Item rendering with icons and stack counts
     - Tooltip display on hover/selection
     - Consumable use with confirmation
     - Equipment slot display (weapon, charm, armor)
  2. Add keyboard/gamepad navigation (arrow keys, confirm, cancel)
  3. Integrate with `MenuManager` for pause menu access
  4. Test with manual verification: `./vania --seed 42 --play` then open inventory

---

## Headless Test Environment

- **Stated Goal**: ROADMAP specifies "Test Coverage ≥ 82%" across all packages
- **Current State**: 7 packages (`camera`, `engine`, `engine/ecs`, `input`, `menu`, `render`, `settings`) fail tests with "X11: The DISPLAY environment variable is missing" panic. CI cannot run full test suite.
- **Impact**: ~40% of packages cannot be tested in headless CI environments. Regressions in these packages go undetected until manual testing.
- **Closing the Gap**:
  1. Create `internal/testutil/ebiten_mock.go` with mock `*ebiten.Image` factory
  2. Add build tag `//go:build !headless` to tests requiring display
  3. Create parallel `_headless_test.go` files using mocks
  4. Alternatively: use Xvfb in CI (`xvfb-run go test ./...`)
  5. Validate with `DISPLAY= go test ./... 2>&1 | grep -c "^ok"` = 19

---

## Stat-Upgrade and Ability Tree System

- **Stated Goal**: ROADMAP v2.0 specifies "Stat-upgrade nodes (max HP, attack, defense) interspersed with movement unlocks" and "Ability tree UI with lock/unlock animations"
- **Current State**: Abilities are unlocked linearly via boss defeats. No stat upgrade nodes, no ability tree structure, no upgrade UI.
- **Impact**: Progression is linear — no player choice in character build. Reduces replay value and strategic depth.
- **Closing the Gap**:
  1. Create `internal/entity/ability_tree.go` with:
     - `AbilityNode` struct (id, name, requirements, statBoosts, unlocked)
     - Tree structure connecting movement abilities with stat nodes
  2. Create `internal/menu/ability_tree.go` for visual tree display
  3. Integrate upgrade points earned from bosses/exploration
  4. Add `SetGenre()` for ability name flavoring
  5. Test with `go test ./internal/entity -v -run TestAbilityTree`

---

## Full Game Determinism Verification

- **Stated Goal**: README states "deterministic: same seed = same game, identical output given identical inputs"
- **Current State**: `internal/pcg/seed_test.go` tests seed derivation. `internal/physics/glide_grapple_test.go` tests grapple physics determinism. No test verifies full game generation produces identical complete games.
- **Impact**: Determinism claim is not fully validated. Subtle non-determinism (map ordering, goroutine scheduling) could produce different games for same seed.
- **Closing the Gap**:
  1. Create `internal/pcg/determinism_test.go` with `TestFullGameDeterminism`:
     - Generate two complete games with seed 42
     - Compare world graph structure (room count, edge list)
     - Compare entity counts and positions
     - Compare generated asset checksums (sprite pixels, audio samples)
  2. Run in CI to catch determinism regressions
  3. Validate with `go test ./internal/pcg -v -run TestFullGameDeterminism -count=10`

---

## Audio Package Test Coverage

- **Stated Goal**: ROADMAP specifies "Test Coverage ≥ 82%" target
- **Current State**: `internal/audio` has 68.1% coverage — 14 percentage points below target.
- **Impact**: Music generation and adaptive music logic are undertested. Bugs in chord progression or layer mixing may go undetected.
- **Closing the Gap**:
  1. Add tests for `music_gen.go`:
     - `TestChordProgressionGeneration` — verify valid chord sequences
     - `TestLayerMixing` — verify 4 layers combine correctly
     - `TestBiomeTrackDeterminism` — same seed = same track
  2. Add tests for `adaptive.go`:
     - `TestIntensityTransitions` — verify smooth crossfades
     - `TestLayerActivation` — correct layers at each intensity
  3. Validate with `go test -cover ./internal/audio` showing ≥82%

---

## Summary Table

| Gap | Priority | Effort | ROADMAP Phase |
|-----|----------|--------|---------------|
| Ranged Combat | CRITICAL | Medium | v2.0 |
| Status Effects | CRITICAL | Medium | v2.0 |
| ECS Integration | HIGH | Large | v1.0 (incomplete) |
| Physics SetGenre | HIGH | Small | v1.0 (incomplete) |
| Narrative SetGenre | HIGH | Small | v1.0 (incomplete) |
| Inventory UI | MEDIUM | Medium | v2.0 |
| Headless Tests | MEDIUM | Medium | v5.0 |
| Ability Tree | MEDIUM | Large | v2.0 |
| Determinism Test | LOW | Small | v1.0 |
| Audio Coverage | LOW | Small | v5.0 |
