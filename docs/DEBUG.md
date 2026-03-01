# Generic UI/UX Debugging Prompt

Use this prompt (adapted to your current codebase state) when asking an AI coding agent to audit and fix UI/UX issues in the VANIA game engine. Copy the template below, fill in the `[bracketed]` sections, and submit.

---

## Prompt Template

```
# VANIA UI/UX Debug Pass

## Objective
Audit and fix all UI/UX issues in the VANIA game engine that affect player experience. Apply fixes directly to the codebase. Do not ask for approval before each change.

## Constraints (do not violate)
- Preserve deterministic generation — do not modify seed derivation or RNG usage.
- Maintain zero-external-assets philosophy — no font files, image files, etc.
- Keep all existing tests passing (`go test ./...`).
- Follow project conventions: local RNG from seeds, `New*` constructors, no global state mutation.
- Target 960×640 fixed resolution (ScreenWidth/ScreenHeight constants).

## Audit Checklist

Work through every category below. For each item, verify correctness by reading the relevant source, identify any mismatch or bug, and fix it in-place. Skip items that are already correct.

### 1. Coordinate System Consistency
- [ ] Camera offset applied with correct sign in every `Render*` function (world-to-screen = subtract camera position).
- [ ] Camera clamped to room bounds so single-screen rooms stay at origin (0,0).
- [ ] All world-space positions (player, enemies, items, doors, platforms, particles, attack hitboxes) use the same coordinate convention.
- [ ] `GetCameraOffset()` return values match the sign convention used by callers.

### 2. Platform & Tile Rendering
- [ ] Platform `Width`/`Height` fields are in pixels. Rendering loop divides by `TileSize` to get tile count (not iterating Width×Height times).
- [ ] Tile draw positions align to the grid: `platform.X + col*TileSize`, `platform.Y + row*TileSize`.
- [ ] Background tiles fill the full room area without gaps or overflow.
- [ ] Hazard and door rectangles match their collision boundaries exactly.

### 3. Collision & Physics
- [ ] `PlayerWidth`/`PlayerHeight` in physics.go matches the rendered player sprite dimensions.
- [ ] `ResolveCollisionWithPlatforms` uses pixel-based platform AABB, not tile counts.
- [ ] Screen-boundary collision (floor, walls) uses `ScreenWidth`/`ScreenHeight` constants, not magic numbers.
- [ ] Enemy-platform collision resolution places enemies on the platform surface, not inside it.
- [ ] Item pickup collision zone matches rendered item size (`GetBounds()` returns correct pixel dimensions).
- [ ] Door interaction zone matches rendered door rectangle.
- [ ] Attack hitbox aligns with the player sprite's facing side and vertical center.

### 4. Spawn Placement
- [ ] Player spawn position lands on a valid platform (not floating or inside geometry).
- [ ] Enemy spawn Y is calculated from the ground platform surface minus enemy height.
- [ ] Item spawn Y is calculated from the ground platform surface minus item height.
- [ ] After room transitions, `SpawnEnemiesForRoom` and `SpawnItemsForRoom` use platform-aware placement.
- [ ] Boss spawn position fits within the room (boss sprites can be 128×128).

### 5. HUD Layout
- [ ] Health bar position uses layout constants (`HealthBarX`, `HealthBarY`), not magic numbers.
- [ ] Ability icons are positioned below the health bar using `AbilityIconY` constant.
- [ ] Ability icons are cached and only regenerated when ability state changes (not every frame).
- [ ] HUD elements do not overlap each other at 960×640.
- [ ] Debug overlay (F3) is positioned below all HUD elements using calculated Y from layout constants.

### 6. Text Rendering
- [ ] `BitmapTextRenderer` and `DebugTextRenderer` use identical character metrics (8×12).
- [ ] `MeasureText` fallback returns `len(text)*8, 12` (not 6×16 or other mismatched values).
- [ ] Menu title and instruction text are horizontally centered: `(ScreenWidth - textWidth) / 2`.
- [ ] `CharWidth`/`CharHeight` constants in the menu package match the render package metrics.
- [ ] Per-character rendering uses `WritePixels` batch operation, not individual 1×1 `DrawImage` calls.

### 7. Message & Notification Display
- [ ] Locked-door message is centered on screen: `(ScreenWidth - MessageWidth) / 2`.
- [ ] Item-collection message is positioned below the HUD area (not overlapping ability icons).
- [ ] Message background, progress bar, and text are drawn in correct z-order: border → background → progress → text.
- [ ] Controls hint in top-right corner accounts for actual text pixel width.

### 8. Menu System
- [ ] Menu items use consistent `CharWidth`/`CharHeight` for positioning.
- [ ] Selection indicator (">>") aligns vertically with the selected item text.
- [ ] Scaled text for selected items centers the scale transform around the text midpoint.
- [ ] Instructions line at bottom is centered using measured text width.

### 9. Transitions
- [ ] After a room transition completes, the camera resets to the new room's origin.
- [ ] Player position is set to a valid entry point in the new room (near the source door).
- [ ] Enemies and items from the previous room are fully replaced (no stale instances).
- [ ] Transition overlay (fade/slide/iris) covers the full 960×640 screen.

### 10. Visual Feedback
- [ ] Damage numbers render at the correct world-to-screen position (with camera offset).
- [ ] Particle effects use consistent camera transform (subtract, not add).
- [ ] Enemy health bars are positioned above the enemy sprite (screenY - offset), not overlapping it.
- [ ] Invulnerability flicker/transparency is visible (alpha < 255 during i-frames).

## Procedural Checks (things that break silently)
- [ ] `math.Sin`/`math.Cos` are used for trigonometry (not hand-rolled Taylor approximations).
- [ ] No division that can produce zero-size images (`ebiten.NewImage(0, 0)` panics).
- [ ] Progress bar width calculation doesn't go negative or exceed container width.
- [ ] All `map` fields are initialized with `make()` before first write.

## Output Format
For each fix applied, report:
1. **File and function** modified
2. **Problem** (one line)
3. **Fix** (one line)

After all fixes, run:
    go build ./cmd/game && go test ./...
and confirm both succeed with zero errors.

## Known Issues (fill in current problems)
[Paste any specific bugs, screenshots, or error messages here.
 If starting from a clean audit, delete this section.]
```

---

## How to Use

1. **Copy** the template above (inside the code fence).
2. **Fill in** the "Known Issues" section with any specific bugs you've observed, or delete it for a full audit.
3. **Submit** to your AI coding agent.
4. **Review** the reported fixes and verify `go build` + `go test` pass.

## Checklist Coverage

The audit covers these source files and their typical failure modes:

| File | Common Issues |
|------|--------------|
| `internal/render/renderer.go` | Camera sign, tile loop math, HUD layout, ability icon caching |
| `internal/render/text.go` | Mismatched font metrics between renderers |
| `internal/menu/menu.go` | Text centering, selection indicator alignment |
| `internal/physics/physics.go` | Hitbox size vs sprite size, screen boundary constants |
| `internal/world/platform_gen.go` | Trig functions, platform pixel dimensions |
| `internal/engine/runner.go` | Spawn positions, message z-order, debug overlay overlap |
| `internal/engine/transitions.go` | Post-transition spawn placement |
| `internal/engine/combat.go` | Attack hitbox alignment with facing direction |
| `internal/camera/camera.go` | Bounds clamping, viewport calculation |

## Revision History

| Date | Change |
|------|--------|
| 2026-03-01 | Initial version based on playability fix audit |
