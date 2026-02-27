# Implementation Gaps

## v1.0 Completion — Core Engine + Playable Single-Player

### ECS Architecture Integration Strategy
- **Gap**: The ROADMAP specifies porting venture's `pkg/engine` (ECS core) to `internal/engine/ecs`, but the current codebase uses a monolithic `GameRunner` pattern in `internal/engine/runner.go` (~30KB). There is no specification for how the new ECS framework should coexist with or replace the existing `GameRunner` architecture.
- **Impact**: Without a clear integration strategy, ECS implementation risks either duplicating logic (systems alongside monolith) or requiring a massive rewrite of `runner.go`. Either path has significant regression risk.
- **Resolution needed**: Decide between (a) incremental wrapping — ECS systems delegate to existing `GameRunner` methods, gradually migrating logic into discrete systems, or (b) clean rewrite — build ECS from scratch and port `runner.go` logic into individual systems. Recommend option (a) for v1.0 to minimize risk.

### Grapple Hook Physics Specification
- **Gap**: The grapple hook mechanic is listed in the ROADMAP but has no technical specification for rope physics parameters (rope length, swing angular velocity, anchor-point detection range, launch speed, detach conditions).
- **Impact**: Implementation cannot begin without defined physics parameters; incorrect values will produce unplayable or unfun grapple mechanics that require extensive tuning.
- **Resolution needed**: Define grapple hook physics parameters: max rope length (tiles), swing damping factor, launch velocity, max anchor detection range, detach-on-ground behavior. Reference games like Bionic Commando or Celeste for parameter ranges. Consider prototyping with placeholder values and tuning iteratively.

### GenreSwitcher Runtime vs Startup Behavior
- **Gap**: The ROADMAP defines `SetGenre(genreID string)` on every System but does not specify whether genre can change at runtime (mid-game) or only at game start. Current architecture has no hot-swap mechanism for tilesets, palettes, or audio presets.
- **Impact**: If runtime switching is required, all asset caches must support invalidation and regeneration, which adds significant complexity. If startup-only, the implementation is substantially simpler.
- **Resolution needed**: Confirm scope — recommend startup-only genre selection for v1.0. Document that runtime genre switching is a v2.0+ feature if desired.

### Genre-Specific Tile Vocabularies
- **Gap**: The ROADMAP lists "per-genre room tile vocabulary via `SetGenre()`" but does not define the specific tile types, visual properties, or platforming behaviors per genre. The genre table in the ROADMAP lists high-level concepts (e.g., "vine-covered doorways" for fantasy, "hull-breach bulkheads" for scifi) but not implementable tile specifications.
- **Impact**: Tile vocabulary implementation requires concrete definitions of tile types, collision properties, and visual attributes for each genre. Without these, only the `fantasy` baseline can be fully implemented.
- **Resolution needed**: Define per-genre tile vocabulary tables with: tile name, collision type (solid/platform/hazard/decorative), visual description, and any special behavior. For v1.0, implementing `fantasy` fully and stubbing other genres with palette-swapped versions is acceptable.

### Input Buffering Window Calibration
- **Gap**: Input buffering is listed as a requirement but no specification exists for buffer window durations, which actions should be bufferable, or how buffering interacts with existing input handling in `internal/input/`.
- **Impact**: Incorrect buffer windows produce either unresponsive or overly generous input handling, both of which harm platformer feel.
- **Resolution needed**: Define buffer window in frames (recommend 6 frames at 60fps = 100ms as starting point based on industry standard). Specify bufferable actions: jump, attack, dash. Non-bufferable: movement direction, pause.

### Camera Transition Style Specifications
- **Gap**: The ROADMAP requires "camera transition animations on room change" but does not specify which transition styles are needed, their durations, or how they interact with gameplay (does the player freeze during transition? Do enemies pause?).
- **Impact**: The existing `RoomTransitionHandler` has a basic fade effect. Enhancing it requires knowing the target transition types and gameplay behavior during transitions.
- **Resolution needed**: Define transition types (fade-to-black, directional slide, iris wipe), duration range (0.3–0.8 seconds), and gameplay pause behavior (recommend: freeze all gameplay during transition, resume on completion).

### Backtracking Shortcut Generation Algorithm
- **Gap**: The ROADMAP requires "backtracking shortcuts that unlock as abilities are gained" but the world generation algorithm in `internal/world/` does not have a defined strategy for placing shortcut connections in the room graph.
- **Impact**: Without a defined algorithm, shortcut placement may create sequence-breaking paths or fail to provide meaningful backtracking convenience.
- **Resolution needed**: Define shortcut placement rules: (1) shortcuts connect rooms separated by ≥5 edges on the critical path, (2) shortcuts require an ability gained after the destination room, (3) maximum 3–5 shortcuts per world, (4) shortcuts are one-way until first traversal.
