# AUDIT — 2026-03-21

## Project Goals

VANIA is a procedural Metroidvania game engine that claims to generate **ALL** game content (graphics, audio, narrative, levels) algorithmically at runtime from a single seed value. The project promises:

1. **Zero External Assets**: No pre-rendered images, bundled audio files, or static narrative content
2. **Deterministic Generation**: Same seed produces identical results
3. **Procedural Graphics**: Pixel art sprites via cellular automata with symmetry transforms
4. **Procedural Audio**: Waveform synthesis with ADSR envelopes, chord progressions, layered music
5. **Procedural Narrative**: Story themes, character generation, item descriptions, world lore
6. **World Generation**: Graph-based 80-150 rooms, 4-6 biomes, ability-gated progression
7. **Combat System**: Melee, ranged, block/parry, knockback, damage numbers
8. **AI System**: Patrol/chase/flee/flying patterns, group coordination (5 formations)
9. **Platforming Physics**: Double-jump, wall-jump, dash, glide, grapple hook, coyote-time
10. **Genre System**: 5 genres (fantasy, scifi, horror, cyberpunk, postapoc) with `SetGenre()` on all systems
11. **Achievement System**: 19 achievements across 6 categories
12. **Adaptive Music**: Dynamic multi-layer music responding to gameplay
13. **Save System**: Multiple save slots with checkpoint autosave

**Target Audience**: Developers interested in PCG techniques; players seeking infinite unique Metroidvania experiences.

---

## Goal-Achievement Summary

| Goal | Status | Evidence |
|------|--------|----------|
| Procedural Sprite Generation (CA) | ✅ Achieved | `internal/graphics/sprite_gen.go:107-128` |
| Audio Waveform Synthesis | ✅ Achieved | `internal/audio/synth.go:54-86` (5 waveforms) |
| ADSR Envelopes | ✅ Achieved | `internal/audio/synth.go:88-129` |
| Music Chord Progressions | ✅ Achieved | `internal/audio/music_gen.go:80-115` |
| Layered Music Composition | ✅ Achieved | `internal/audio/music_gen.go:58-78` (4 layers) |
| Narrative Generation | ✅ Achieved | `internal/narrative/story_gen.go` |
| Graph-Based World | ✅ Achieved | `internal/world/graph_gen.go:81-103` |
| Ability-Gating | ✅ Achieved | `internal/world/graph_gen.go:198-202` |
| Melee Combat | ✅ Achieved | `internal/engine/combat.go:128-137` |
| Ranged Combat | ❌ Missing | No projectile system in `combat.go` |
| Block/Parry | ✅ Achieved | `internal/engine/combat.go:139-168` |
| Knockback | ✅ Achieved | `internal/engine/combat.go:295-312` |
| Damage Numbers | ✅ Achieved | `internal/engine/combat.go:324-339` |
| AI Patrol/Chase/Flee/Flying | ✅ Achieved | `internal/entity/ai.go:202-316` |
| AI Group Coordination (5 formations) | ✅ Achieved | `internal/entity/ai_advanced.go:66-70` |
| Double-Jump | ✅ Achieved | `internal/physics/physics.go:237-257` |
| Wall-Jump | ✅ Achieved | `internal/physics/physics.go:249-254` |
| Dash | ✅ Achieved | `internal/physics/physics.go:276-281` |
| Glide | ✅ Achieved | `internal/physics/physics.go:300-304` |
| Grapple Hook | ✅ Achieved | `internal/physics/physics.go:306-382` |
| Coyote Time | ✅ Achieved | `internal/physics/physics.go:32` (6 frames) |
| GenreSwitcher Interface | ✅ Achieved | `internal/engine/ecs/system.go:21-25` |
| SetGenre on Render | ✅ Achieved | `internal/render/renderer.go:96` |
| SetGenre on Audio | ✅ Achieved | `internal/audio/player.go:57` |
| SetGenre on Menu | ✅ Achieved | `internal/menu/menu.go:157` |
| SetGenre on Physics | ❌ Missing | Not implemented |
| SetGenre on Narrative | ❌ Missing | Not implemented |
| Status Effects (burn/freeze/poison) | ❌ Missing | Not implemented |
| ECS Framework | ⚠️ Partial | Built in `internal/engine/ecs/` but not integrated with `runner.go` |
| Achievement System (19 achievements) | ✅ Achieved | `internal/achievement/achievement.go:132-347` |
| Adaptive Music | ✅ Achieved | `internal/audio/adaptive.go` (4 intensities, 5 layers) |
| Save System (5 slots) | ✅ Achieved | `internal/save/save_manager.go:72` |
| Checkpoint Autosave | ✅ Achieved | `internal/save/checkpoint.go:35-54` |
| Multi-Phase Bosses | ✅ Achieved | `internal/entity/enemy_gen.go:237-276` |
| Inventory Screen UI | ❌ Missing | Only data structure exists (`game.go:47`), no UI |
| Stat-Upgrade Nodes | ❌ Missing | Not implemented |
| Ability Tree UI | ❌ Missing | Not implemented |
| Determinism Testing | ⚠️ Partial | Seed tests exist but no full game reproducibility test |

---

## Findings

### CRITICAL

- [ ] **Ranged Attack System Not Implemented** — `internal/engine/combat.go` — The README claims "Ranged attack (projectile with falloff)" but no ranged attack or projectile system exists in the combat code. Only melee attacks are implemented. — **Remediation:** Implement `PlayerRangedAttack()` in `combat.go` with projectile spawning, trajectory, and falloff damage calculation. Add projectile entity type to `internal/entity/`. Validation: `grep -n "Projectile\|ranged" internal/engine/combat.go` should return results.

- [ ] **Status Effect System Not Implemented** — `internal/engine/` — ROADMAP v2.0 specifies "Burn, freeze, shock, poison, bleed, slow, haste" status effects but none exist. The string labels "poison" appear in enemy generation but no gameplay mechanics are implemented. — **Remediation:** Create `internal/engine/status.go` with `StatusEffect` struct (type, duration, stacks, source), `StatusManager` for application/tick/expiration. Integrate with combat damage flow. Validation: `go test ./internal/engine -v -run TestStatusEffects`.

### HIGH

- [ ] **ECS Framework Not Integrated** — `internal/engine/runner.go:1-1224` — The ECS framework in `internal/engine/ecs/` is complete but the 1224-line `GameRunner` monolith does not use it. The ROADMAP states ECS is required, but `runner.go` directly manages all systems without delegation. — **Remediation:** Incrementally migrate `runner.go` logic into ECS systems: create `RenderSystem`, `PhysicsSystem`, `CombatSystem`, `AudioSystem` that implement `ecs.System`. Have `GameRunner` delegate to `ecs.SystemManager`. Validation: `grep -c "SystemManager" internal/engine/runner.go` should be ≥5.

- [ ] **SetGenre Missing on Physics Package** — `internal/physics/physics.go` — ROADMAP requires `SetGenre()` on all systems for genre-specific hazards (magic barriers, airlock vents, radiation clouds), but physics has no genre awareness. — **Remediation:** Add `SetGenre(genreID string)` method to `physics.go` that configures genre-specific hazard parameters. Add hazard collision types per genre. Validation: `grep -n "SetGenre" internal/physics/physics.go`.

- [ ] **SetGenre Missing on Narrative Package** — `internal/narrative/story_gen.go` — ROADMAP requires `SetGenre()` for re-skinning narrative vocabulary, but narrative generators are stateless functions with no genre state. — **Remediation:** Add `NarrativeGenerator` struct with `SetGenre(genreID string)` method that configures vocabulary tables per genre. Validation: `grep -n "SetGenre" internal/narrative/*.go`.

- [ ] **Tests Fail Without X11 Display** — `internal/camera`, `internal/engine`, `internal/input`, `internal/menu`, `internal/render`, `internal/settings` — 7 packages fail tests with "X11: The DISPLAY environment variable is missing" panic. CI/CD cannot run full test suite. — **Remediation:** Create `internal/testutil/ebiten_mock.go` with mock image/screen types. Refactor tests to use mocks when `DISPLAY` is unset. Validation: `DISPLAY= go test ./... 2>&1 | grep -c "^ok"` should equal package count.

### MEDIUM

- [ ] **Duplicate getCharPattern Implementation** — `internal/render/text.go:179` and `internal/menu/menu.go:409` — Identical 44-complexity character bitmap functions exist in two files. Menu version has cyclomatic complexity 44, render version 25 (different character coverage). — **Remediation:** Extract shared implementation to `internal/render/font.go`, export as `GetCharPattern()`, import in both packages. Validation: `grep -r "getCharPattern" internal/ | wc -l` should show 1 definition.

- [ ] **Inventory Screen UI Not Implemented** — `internal/menu/` — ROADMAP v2.0 specifies "Inventory screen (grid layout)" but only data structure exists in `game.go:47`. No inventory UI renders or handles navigation. — **Remediation:** Create `internal/menu/inventory.go` with grid-based item display, tooltips, and keyboard/gamepad navigation. Validation: `ls internal/menu/inventory.go` should exist.

- [ ] **runner.Update() Extreme Complexity** — `internal/engine/runner.go:142` — Cyclomatic complexity 63, far exceeding the 15-threshold for maintainability. Single function handles input, physics, combat, AI, transitions, achievements, music, and rendering decisions. — **Remediation:** Extract state-specific update logic into dedicated methods: `updatePlaying()`, `updatePaused()`, `updateTransition()`, `updateGameOver()`. Delegate to subsystem managers. Validation: `go-stats-generator analyze . --format json | jq '.functions[] | select(.name == "Update" and .file | contains("runner")) | .complexity.cyclomatic'` should be ≤25.

- [ ] **Audio Package Test Coverage Below Target** — `internal/audio/` — Current coverage 68.1%, target is 82% per ROADMAP. Music generation and adaptive music have limited test coverage. — **Remediation:** Add tests for `music_gen.go` chord progressions, `adaptive.go` intensity transitions. Validation: `go test -cover ./internal/audio` should show ≥82%.

### LOW

- [ ] **Undocumented Exported Function** — `internal/settings/settings.go:68` — `String()` method on settings type lacks doc comment. 1 of 626 functions undocumented (99.8% coverage). — **Remediation:** Add `// String returns a human-readable representation of the settings.` before line 68. Validation: `go doc github.com/opd-ai/vania/internal/settings Settings.String`.

- [ ] **Duplication Ratio Near Threshold** — Project-wide — 10% duplication ratio (1610 lines of 9125), slightly below the 10% warning threshold. 38 clone pairs detected. — **Remediation:** Review top clone pairs in `runner.go:996-1007` and `cmd/game/main.go:361-373`. Extract shared display formatting logic. Validation: `go-stats-generator analyze . --format json | jq '.duplication.duplication_ratio'` should be <0.08.

- [ ] **No Comprehensive Determinism Test** — `internal/pcg/seed_test.go` — Seed derivation determinism is tested, but no test verifies full game generation produces identical output for same seed across runs. — **Remediation:** Create `TestFullGameDeterminism` that generates two games with same seed, compares world graph, entity counts, and asset checksums. Validation: `go test ./internal/pcg -v -run TestFullGameDeterminism`.

---

## Metrics Snapshot

| Metric | Value |
|--------|-------|
| Total Lines of Code | 9,125 |
| Total Functions | 626 (87 standalone, 539 methods) |
| Total Structs | 114 |
| Total Interfaces | 4 |
| Total Packages | 19 |
| Total Files | 46 |
| Average Function Length | 18.3 lines |
| Average Cyclomatic Complexity | 2.87 |
| Functions with Complexity ≥15 | 7 |
| Documentation Coverage | 98.2% (615/626) |
| Duplication Ratio | 9.97% (1,610 lines) |
| Test Coverage (pcg) | 96.3% |
| Test Coverage (graphics) | 99.1% |
| Test Coverage (achievement) | 90.9% |
| Test Coverage (audio) | 68.1% |
| Packages Failing Tests (X11) | 7 |

---

## Dependency Assessment

| Dependency | Version | Status | Notes |
|------------|---------|--------|-------|
| `github.com/hajimehoshi/ebiten/v2` | v2.6.3 | ⚠️ Outdated | v2.8+ available; no security issues |
| `github.com/ebitengine/oto/v3` | v3.1.0 | ✅ Stable | Audio backend |
| `golang.org/x/image` | v0.12.0 | ✅ Stable | Image processing |

**Recommendation:** Upgrade Ebiten to v2.8+ after v2.0 features complete for bug fixes and performance improvements.
