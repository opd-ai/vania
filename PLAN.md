# Implementation Plan: v2.0 Status Effects & Inventory Systems

## Project Context
- **What it does**: A procedural Metroidvania engine generating all assets (graphics, audio, narrative, levels) algorithmically from a single seed
- **Current goal**: Complete v2.0 features — Status Effects, full Inventory system, and remaining AI/Combat enhancements
- **Estimated Scope**: Medium (5–15 items requiring implementation)

## Goal-Achievement Status

| Stated Goal | Current Status | This Plan Addresses |
|-------------|---------------|---------------------|
| v1.0 Core Engine + Playable Single-Player | ✅ Complete | No |
| v2.0 Combat System (Melee, Ranged, Block/Parry) | ✅ Complete | No |
| v2.0 Status Effects (Burn, Freeze, Poison, etc.) | ❌ Not Started | **Yes** |
| v2.0 AI Platformer-Aware Pathfinding | ⚠️ Partial | **Yes** |
| v2.0 Inventory Screen & Equipment Slots | ❌ Not Started | **Yes** |
| v2.0 Drop Rarity Tiers | ❌ Not Started | **Yes** |
| v2.0 All 5 Genres Full Integration | ❌ Not Started | No (v2.0 late-stage) |
| Test Coverage ≥ 82% | ⚠️ Partial (~60%) | **Yes** |
| UI Performance Optimization | ⚠️ Identified Issues | **Yes** |

## Metrics Summary
- **Complexity hotspots on goal-critical paths**: 6 functions with cyclomatic complexity ≥10
  - `runner.Update()`: 63 (game loop orchestrator — complexity justified)
  - `menu.getCharPattern()`: 44 (font rendering)
  - `render.getCharPattern()`: 25 (font rendering)
  - `runner.Draw()`: 25 (rendering orchestrator)
  - `physics.ResolveCollisionWithPlatforms()`: 17
  - `entity.Update()` (AI): 16
- **Duplication ratio**: Low (<3%) — `getCharPattern` exists in both `menu` and `render` packages
- **Doc coverage**: 98.2% (615/626 functions documented)
- **Test failures**: 8 packages fail due to missing X11/DISPLAY (Ebiten requires display for tests)
- **Package coupling**: `engine` package (2279 LOC) is the integration hub; appropriately sized for orchestration role

## Dependency Analysis

### External Dependencies
| Dependency | Version | Status | Notes |
|------------|---------|--------|-------|
| `github.com/hajimehoshi/ebiten/v2` | v2.6.3 | ⚠️ Outdated | v2.9.x available; consider upgrade for bug fixes |
| `github.com/ebitengine/oto/v3` | v3.1.0 | ✅ Stable | Audio library |
| `golang.org/x/image` | v0.12.0 | ✅ Stable | |

### Recommended Dependency Action
Upgrade Ebiten to v2.8+ after v2.0 features complete. Breaking changes are minimal within v2.x.

---

## Implementation Steps

### Step 1: Fix Text Rendering Duplication
- **Deliverable**: Consolidate `getCharPattern()` from `internal/menu/menu.go` and `internal/render/text.go` into single shared implementation in `internal/render/text.go`
- **Dependencies**: None
- **Goal Impact**: UI Audit Issue #3 (inconsistent text measurement); prepares for UI performance work
- **Acceptance**: 
  - Single `getCharPattern` implementation
  - `menu` package imports from `render`
  - All existing tests pass
- **Validation**: 
  ```bash
  grep -r "getCharPattern" internal/ | wc -l  # Should show 1 definition
  go test ./internal/menu ./internal/render
  ```

### Step 2: Implement Text Rendering Performance Optimization
- **Deliverable**: Batch pixel rendering in `internal/render/text.go` using `WritePixels()` instead of per-pixel draw calls
- **Dependencies**: Step 1
- **Goal Impact**: UI Audit Issue #1 (performance); ROADMAP Performance Optimisation
- **Acceptance**:
  - `BitmapTextRenderer.drawChar()` uses single `WritePixels()` call per character
  - Character images cached and reused
  - Frame time for text-heavy screens reduced by ≥50%
- **Validation**: 
  ```bash
  go test -bench=. ./internal/render  # Add benchmark comparing old vs new
  ```

### Step 3: Implement Ability Icon Caching
- **Deliverable**: Cache ability icons in `internal/render/renderer.go` with invalidation on ability state change
- **Dependencies**: None
- **Goal Impact**: UI Audit Issue #2 (performance)
- **Acceptance**:
  - `Renderer` struct has `abilityIconCache map[string]*ebiten.Image`
  - Icons regenerated only when `abilities` map changes
  - No per-frame icon generation
- **Validation**: 
  ```bash
  go test ./internal/render -v -run TestAbilityIconCaching  # Add test
  ```

### Step 4: Implement Status Effects System
- **Deliverable**: New `internal/engine/status.go` with status effect types, stack/duration management, and genre-mapped variants
- **Dependencies**: None
- **Goal Impact**: v2.0 Status Effects (Burn, Freeze, Shock, Poison, Bleed, Slow, Haste)
- **Acceptance**:
  - `StatusEffect` struct with type, stacks, duration, source
  - `StatusManager` handles application, stacking, tick, and expiration
  - Genre-specific naming via `SetGenre()` (e.g., "irradiate" for `postapoc`)
  - Visual indicators render in HUD
  - Status effects modify combat/movement per type
- **Validation**: 
  ```bash
  go test ./internal/engine -v -run TestStatusEffects
  go-stats-generator analyze ./internal/engine --sections functions | jq '.functions[] | select(.name | contains("Status"))'
  ```

### Step 5: Implement Inventory Screen UI
- **Deliverable**: New `internal/menu/inventory.go` with grid layout inventory screen
- **Dependencies**: Step 1 (text rendering consolidation)
- **Goal Impact**: v2.0 Inventory & Items — grid layout, consumables, key items
- **Acceptance**:
  - Grid-based inventory display (4×8 slots)
  - Item tooltips with procedurally generated descriptions
  - Consumable use with confirmation
  - Key items display with lore text
  - Equipment slots section (weapon, charm, armour)
  - Keyboard/gamepad navigation
- **Validation**: 
  ```bash
  go test ./internal/menu -v -run TestInventory
  ./vania --seed 42 --play  # Manual verification of inventory UI
  ```

### Step 6: Implement Drop Rarity System
- **Deliverable**: Extend `internal/entity/enemy_gen.go` with rarity tiers and drop VFX
- **Dependencies**: Step 4 (status effects for rare item buffs)
- **Goal Impact**: v2.0 Loot/Drops — rarity tiers, drop VFX, audio feedback
- **Acceptance**:
  - `ItemRarity` enum: Common, Uncommon, Rare, Legendary
  - Drop tables weighted by rarity
  - Rarity affects item stats (scaling multipliers)
  - Particle VFX on drop (color-coded by rarity)
  - Procedural audio feedback on pickup
- **Validation**: 
  ```bash
  go test ./internal/entity -v -run TestDropRarity
  go-stats-generator analyze ./internal/entity --sections functions | jq '.functions[] | select(.name | contains("Drop"))'
  ```

### Step 7: Implement Platformer-Aware AI Pathfinding
- **Deliverable**: Extend `internal/entity/ai.go` with ledge detection, wall awareness, and preferred attack range positioning
- **Dependencies**: None
- **Goal Impact**: v2.0 AI Behavior Trees — platformer-aware pathfinding
- **Acceptance**:
  - Enemies detect platform edges and avoid falling
  - Ranged enemies maintain preferred attack distance
  - Flying enemies manage altitude relative to player
  - Wall detection prevents enemies from walking into walls
- **Validation**: 
  ```bash
  go test ./internal/entity -v -run TestPlatformerAI
  go-stats-generator analyze ./internal/entity --sections functions | jq '[.functions[] | select(.complexity.cyclomatic >= 10)]'  # Should not increase complexity significantly
  ```

### Step 8: Add Headless Test Support for Ebiten Packages
- **Deliverable**: Create test helpers that mock Ebiten dependencies for CI environments without display
- **Dependencies**: None
- **Goal Impact**: ROADMAP Test Coverage ≥ 82%
- **Acceptance**:
  - `internal/testutil/ebiten_mock.go` provides mock image/screen types
  - Tests in `camera`, `engine`, `input`, `menu`, `render`, `settings` pass without DISPLAY
  - CI workflow runs all tests successfully
- **Validation**: 
  ```bash
  DISPLAY= go test ./... 2>&1 | grep -c "^ok"  # All packages should pass
  ```

### Step 9: Boss Ability Gate Integration
- **Deliverable**: Connect boss defeats to ability unlocks in `internal/engine/combat.go`
- **Dependencies**: Step 4 (status effects for boss-granted buffs)
- **Goal Impact**: v2.0 Boss Gatekeeper Encounters — boss guards ability, door unlocks on kill
- **Acceptance**:
  - Boss defeat triggers specific ability unlock
  - Ability-gate door behind boss auto-unlocks
  - Genre-themed boss skins via existing `SetGenre()` infrastructure
- **Validation**: 
  ```bash
  go test ./internal/engine -v -run TestBossAbilityGate
  ```

### Step 10: Room Description HUD Integration
- **Deliverable**: Display procedurally generated room descriptions on room entry in HUD
- **Dependencies**: None
- **Goal Impact**: v2.0 Narrative Generation — room descriptions surfaced in HUD
- **Acceptance**:
  - Room entry triggers description display (fade-in/out)
  - Descriptions generated from `narrative` package context
  - Genre-appropriate vocabulary via `SetGenre()`
  - Non-intrusive positioning (bottom of screen)
- **Validation**: 
  ```bash
  go test ./internal/engine -v -run TestRoomDescription
  ```

---

## Scope Assessment Calibration

| Metric | Measured Value | Scope Category |
|--------|---------------|----------------|
| Functions with complexity ≥ 9 | 20 | Medium |
| Functions with complexity ≥ 15 | 6 | Small |
| Undocumented exported functions | 1 | Small |
| Test failure packages (env issue) | 8 | Medium |
| Duplicated implementations | 1 (`getCharPattern`) | Small |

**Overall Scope**: **Medium** — 10 implementation steps, each independently testable with clear metrics.

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Ebiten tests require display | High | Medium | Step 8 adds headless mocks |
| Status effects increase combat complexity | Medium | Medium | Keep `StatusManager` isolated; delegate to existing combat flow |
| Inventory UI complexity | Medium | Low | Reuse existing menu patterns from `internal/menu` |
| runner.Update() complexity growth | Low | High | Status/inventory delegate to sub-managers; don't add to main loop |

---

## Validation Commands Summary

```bash
# After all steps complete:
go test ./...                                    # All tests pass
go vet ./...                                     # No new warnings
go-stats-generator analyze . --sections functions | jq '[.functions[] | select(.complexity.cyclomatic >= 15)] | length'  # Should remain ≤ 6
./vania --seed 42 --play                         # Manual gameplay verification
```

---

## Notes

### Complexity Hotspots (Acknowledged)
The following high-complexity functions are structural necessities:
- `runner.Update()` (63): Main game loop orchestrator — complexity is inherent to its coordination role
- `getCharPattern()` (44/25): Procedural font bitmap data — large switch is appropriate for static data

### GAPS.md Items Deferred
Per GAPS.md, the following require design decisions before implementation:
- ECS Architecture Integration Strategy (v1.0 item marked complete but integration unclear)
- Grapple Hook Physics Specification (marked complete but parameters unspecified)
- GenreSwitcher Runtime vs Startup Behavior (recommend startup-only for v2.0)

These items should be addressed in a separate design document before v2.0 completion.
