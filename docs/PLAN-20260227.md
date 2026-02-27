# UI Fix Plan for VANIA - Procedural Metroidvania Game

## Executive Summary

- **Total Issues Found**: 18
- **Critical**: 2 | **High**: 6 | **Medium**: 7 | **Low**: 3
- **Estimated Total Effort**: 16-24 hours
- **Phases**: 4

This plan addresses UI positioning, layout, and usability issues in the VANIA ebitengine-based game. The game is functional but suffers from hardcoded positions, non-responsive layouts, and inconsistent coordinate systems.

## Table of Contents

1. [Codebase Analysis](#1-codebase-analysis)
2. [Issues Identified](#2-issues-identified)
3. [Root Cause Analysis](#3-root-cause-analysis)
4. [Solution Architecture](#4-solution-architecture)
5. [Implementation Roadmap](#5-implementation-roadmap)
6. [Risk Assessment](#6-risk-assessment)
7. [Testing & Validation](#7-testing--validation)
8. [Success Metrics](#8-success-metrics)
9. [Appendices](#appendices)

---

## 1. Codebase Analysis

### 1.1 UI Architecture Overview

The VANIA game uses **Ebiten v2.6.3** as its game engine. The UI system consists of:

- **Entry Point**: `cmd/game/main.go` - Implements ebiten.Game interface
- **Game Loop**: `internal/engine/runner.go` - Update() and Draw() methods
- **Rendering**: `internal/render/renderer.go` - Core rendering primitives
- **Menus**: `internal/menu/menu.go` - Menu system implementation

**Current Approach**:
- Fixed 960x640 resolution despite window resizing enabled
- Direct positioning using pixel coordinates
- No layout manager or coordinate utilities
- Mix of absolute and relative positioning

### 1.2 UI Element Inventory

| Element | Type | File | Function | Line | Current Position | Dimensions | Purpose | Visibility |
|---------|------|------|----------|------|------------------|------------|---------|------------|
| **Health Bar** | HUD | renderer.go | RenderUI() | 234-258 | (10, 10) | 200x20 | Display player health | Always |
| **Health Bar Background** | HUD | renderer.go | RenderUI() | 242-246 | (10, 10) | 200x20 | Health bar container | Always |
| **Health Bar Fill** | HUD | renderer.go | RenderUI() | 249-257 | (10, 10) | Dynamic width | Visual health indicator | Always |
| **Ability Icons** | HUD | renderer.go | RenderUI() | 260-284 | (10, 40) | 30x30 each | Show unlocked abilities | Always |
| **Debug Info** | Debug | runner.go | Draw() | 561-590 | (0, 0) | Full text | FPS, position, controls | Always |
| **Locked Door Message** | Notification | runner.go | Draw() | 527-541 | (380, 300) | 200x40 | Door unlock requirement | Conditional |
| **Item Collection Message** | Notification | runner.go | Draw() | 544-558 | (380, 80) | 200x40 | Item pickup notification | Conditional |
| **Menu Title** | Menu | menu.go | Draw() | 222-226 | `480 - len*4, 100` | Variable | Menu screen title | Menu active |
| **Menu Items** | Menu | menu.go | Draw() | 230-248 | (220, 200+i*40) | Variable | Selectable options | Menu active |
| **Menu Selection Indicator** | Menu | menu.go | Draw() | 242-244 | (200, varies) | Text | Shows selected item | Menu active |
| **Menu Instructions** | Menu | menu.go | Draw() | 251-253 | `480 - len*3, 500` | Variable | Control instructions | Menu active |
| **Enemy Health Bars** | HUD | renderer.go | RenderEnemy() | 341-362 | Above enemy | Enemy width x 4 | Enemy health display | Per enemy |
| **Transition Fade** | Effect | renderer.go | RenderTransitionEffect() | 385-396 | (0, 0) | 960x640 | Room transition effect | During transition |
| **Particles** | Effect | renderer.go | RenderParticles() | 399-442 | World space | 1-10px | Visual effects | Always |
| **Attack Effect** | Effect | renderer.go | RenderAttackEffect() | 365-382 | Attack hitbox | Variable | Player attack visual | During attack |
| **Item Glow** | Effect | renderer.go | RenderItem() | 462-469 | Item position | 1.5x item size | Collectible indicator | Uncollected items |

### 1.3 Coordinate Systems in Use

**Three different coordinate systems identified:**

1. **Screen Space** (renderer.go, menu.go)
   - Origin: Top-left corner (0, 0)
   - Range: 0-960 width, 0-640 height
   - Used for: UI elements, menus

2. **World Space** (runner.go, renderer.go)
   - Origin: Variable based on camera position
   - Range: Unlimited
   - Used for: Game objects, enemies, items
   - Transformed to screen space via camera offset

3. **Approximation-based** (menu.go)
   - Uses text length * character width estimate
   - Example: `titleX := 480 - len(title)*4`
   - Imprecise and unreliable

---

## 2. Issues Identified

### 2.1 Critical Issues

#### ISSUE-01: Non-Responsive Layout Despite Window Resizing
**Severity**: Critical  
**Category**: Layout  
**Affected Element**: All UI elements  
**Location**: cmd/game/main.go:91, runner.go:594-596  
**Current State**: Layout() returns fixed 960x640 despite `WindowResizingModeEnabled` set to true  
**Impact**: Window can resize but UI elements remain at fixed positions, causing cropping or excessive margins  
**Screenshot/Visual**: At 1920x1080, health bar stays at top-left with large unused space; at 800x600, elements may be cut off  

#### ISSUE-02: Debug Text Overlaps Game UI
**Severity**: Critical  
**Category**: Usability/Readability  
**Affected Element**: Debug Info  
**Location**: runner.go:590  
**Current State**: Debug info renders at (0, 0) overlapping health bar and ability icons  
**Impact**: Makes health bar difficult to read; cluttered visual experience  
**Screenshot/Visual**: Multi-line debug text covers first 150px of top-left corner where UI elements are positioned  

### 2.2 High Priority Issues

#### ISSUE-03: Approximate Text Centering
**Severity**: High  
**Category**: Positioning  
**Affected Element**: Menu Title, Menu Instructions  
**Location**: menu.go:224, 253  
**Current State**: Uses `len(title)*4` and `len(instructions)*3` to estimate center position  
**Impact**: Text is rarely centered correctly; varies by font and string content  
**Screenshot/Visual**: Titles appear off-center by 10-50 pixels depending on content  

#### ISSUE-04: Hardcoded UI Positions
**Severity**: High  
**Category**: Positioning  
**Affected Element**: Health Bar, Ability Icons  
**Location**: renderer.go:238-239, 261-264  
**Current State**: Uses `barX := 10`, `barY := 10`, `abilityX := barX`  
**Impact**: No constants defined; changing UI layout requires searching all files  
**Screenshot/Visual**: Health bar always at exact top-left regardless of screen size  

#### ISSUE-05: Message Centering Based on Fixed Width
**Severity**: High  
**Category**: Positioning  
**Affected Element**: Locked Door Message, Item Collection Message  
**Location**: runner.go:529, 546  
**Current State**: Uses `ScreenWidth/2 - 100` assuming 200px message width  
**Impact**: Messages not truly centered; breaks if message dimensions change  
**Screenshot/Visual**: 200px wide messages centered, but if width changes, alignment breaks  

#### ISSUE-06: No Named Constants for UI Dimensions
**Severity**: High  
**Category**: Layout  
**Affected Element**: All UI elements  
**Location**: Throughout renderer.go, runner.go, menu.go  
**Current State**: Magic numbers like 200, 20, 30, 100, 40, 10 scattered in code  
**Impact**: Difficult to maintain consistent spacing; changes require multiple file edits  
**Screenshot/Visual**: Health bar 200px, height 20px, abilities 30px, margins 10px - all hardcoded  

#### ISSUE-07: Inconsistent Spacing Between UI Elements
**Severity**: High  
**Category**: Layout  
**Affected Element**: Health Bar + Ability Icons  
**Location**: renderer.go:261  
**Current State**: `abilityY := barY + barHeight + 10` - spacing is arbitrary  
**Impact**: No spacing scale or system; visual inconsistency  
**Screenshot/Visual**: 10px gap between health and abilities, but other elements use different gaps  

#### ISSUE-08: Menu Item Spacing Hardcoded
**Severity**: High  
**Category**: Layout  
**Affected Element**: Menu Items  
**Location**: menu.go:231  
**Current State**: `y := startY + i*40` - assumes 40px per item  
**Impact**: Cannot adjust for longer text or accessibility needs  
**Screenshot/Visual**: Menu items evenly spaced at 40px intervals regardless of content  

### 2.3 Medium Priority Issues

#### ISSUE-09: Camera Offset Not Applied to UI Elements
**Severity**: Medium  
**Category**: Positioning  
**Affected Element**: Enemy Health Bars  
**Location**: renderer.go:302-310, 344  
**Current State**: Enemy rendering applies camera offset, but health bars use `screenY - 8`  
**Impact**: Health bars positioned correctly but logic is duplicated and fragile  
**Screenshot/Visual**: Health bars float above enemies but recalculate offset separately  

#### ISSUE-10: Fixed Message Box Dimensions
**Severity**: Medium  
**Category**: Layout  
**Affected Element**: Locked Door Message, Item Message  
**Location**: runner.go:533-534, 550-551  
**Current State**: Creates `NewImage(200, 40)` boxes regardless of text length  
**Impact**: Long messages truncated; short messages have excess space  
**Screenshot/Visual**: All messages fit in 200x40 box, but text may overflow or be too small  

#### ISSUE-11: Ability Icon Labels Missing
**Severity**: Medium  
**Category**: Usability  
**Affected Element**: Ability Icons  
**Location**: renderer.go:260-284  
**Current State**: Shows colored boxes but no text labels  
**Impact**: Players don't know which ability each icon represents  
**Screenshot/Visual**: Four colored squares without indication of double_jump, dash, wall_jump, glide  

#### ISSUE-12: No Visual Feedback for Selected Menu Item
**Severity**: Medium  
**Category**: Usability  
**Affected Element**: Menu Selection Indicator  
**Location**: menu.go:242-244  
**Current State**: Only shows `">"` character; color change doesn't render due to DebugPrint limitations  
**Impact**: Selection indicator is subtle; easy to miss current selection  
**Screenshot/Visual**: Small ">" to left of selected item, same white color as all text  

#### ISSUE-13: Menu Title Position Varies by Text Length
**Severity**: Medium  
**Category**: Positioning  
**Affected Element**: Menu Title  
**Location**: menu.go:264-278  
**Current State**: Different menu titles have different lengths, causing position shift  
**Impact**: Inconsistent visual experience; titles jump around  
**Screenshot/Visual**: "VANIA - Procedural Metroidvania" vs "Settings" cause 100+ pixel position difference  

#### ISSUE-14: Enemy Health Bars Overlap When Enemies Stack
**Severity**: Medium  
**Category**: Layout  
**Affected Element**: Enemy Health Bars  
**Location**: renderer.go:341-362  
**Current State**: Each enemy renders its health bar 8px above sprite  
**Impact**: When multiple enemies are close together, health bars overlap unreadably  
**Screenshot/Visual**: Two enemies stacked vertically show overlapping health bars  

#### ISSUE-15: Particle Rendering Order Inconsistent
**Severity**: Medium  
**Category**: Positioning  
**Affected Element**: Particles  
**Location**: runner.go:499  
**Current State**: Particles rendered before player but should be layered  
**Impact**: Some effects (like dust) should be behind player; others (hits) should be in front  
**Screenshot/Visual**: Jump dust particles appear in front of player sprite  

### 2.4 Low Priority Issues

#### ISSUE-16: No Minimap or Room Indicator
**Severity**: Low  
**Category**: Usability  
**Affected Element**: HUD (missing element)  
**Location**: N/A  
**Current State**: Players have no visual indication of room layout or explored areas  
**Impact**: Navigation difficulty in procedurally generated worlds  
**Screenshot/Visual**: No map visible; only room name in debug text  

#### ISSUE-17: Item Glow Effect Not Animated
**Severity**: Low  
**Category**: Usability  
**Affected Element**: Item Glow  
**Location**: renderer.go:462-469  
**Current State**: Static glow; no pulsing or animation  
**Impact**: Items less noticeable; missed collectibles  
**Screenshot/Visual**: Gold glow around items is static  

#### ISSUE-18: Settings Values Not Persisted
**Severity**: Low  
**Category**: Usability  
**Affected Element**: Settings Menu  
**Location**: menu.go:396-465  
**Current State**: Volume and display settings reset on restart  
**Impact**: Players must reconfigure settings each session  
**Screenshot/Visual**: Settings menu functional but changes don't save  

---

## 3. Root Cause Analysis

### ROOT CAUSE 1: No Centralized UI Layout System
**Evidence**: ISSUE-01, ISSUE-04, ISSUE-06, ISSUE-07, ISSUE-08  
**Scope**: All UI elements (18 affected)  
**Systemic Impact**: Every UI element uses ad-hoc positioning logic. There's no `UILayout` struct, no layout manager, no positioning helpers. Each developer hardcodes coordinates independently.

### ROOT CAUSE 2: Hardcoded Screen Dimensions
**Evidence**: ISSUE-01, ISSUE-05, ISSUE-13  
**Scope**: 15+ locations across 3 files  
**Systemic Impact**: `960` and `640` appear throughout code. Window resizing is enabled but unused. Layout() returns fixed size, making responsive design impossible.

### ROOT CAUSE 3: Ebiten's Limited Text Rendering
**Evidence**: ISSUE-03, ISSUE-12, ISSUE-13  
**Scope**: All text elements (menu titles, menu items, debug text)  
**Systemic Impact**: `ebitenutil.DebugPrint` doesn't support color, font sizing, or text measurement. This forces developers to use length-based approximation for centering.

### ROOT CAUSE 4: Missing UI Constants
**Evidence**: ISSUE-06, ISSUE-07, ISSUE-10  
**Scope**: 20+ magic numbers  
**Systemic Impact**: Values like `10`, `20`, `30`, `40`, `100`, `200` are scattered without meaning. No `const UIMargin = 10` or `const HealthBarWidth = 200`. This makes consistent spacing impossible.

### ROOT CAUSE 5: Mixed Coordinate Systems Without Utilities
**Evidence**: ISSUE-09, ISSUE-15  
**Scope**: World space vs screen space elements  
**Systemic Impact**: No helper functions for coordinate conversion. Each renderer calculates camera offset independently. Error-prone and duplicated logic.

---

## 4. Solution Architecture

### 4.1 Overall Strategy

**Phase 1: Foundation** - Create UI constants and coordinate utilities  
**Phase 2: Refactor Positioning** - Replace hardcoded values with responsive calculations  
**Phase 3: Layout Manager** - Implement centralized layout system  
**Phase 4: Polish** - Add missing features and animations  

**Key Principles:**
1. **Responsive by default** - All UI scales with window size
2. **Centralized constants** - Single source of truth for dimensions
3. **Coordinate utilities** - Helper functions for common calculations
4. **Backward compatible** - Don't break gameplay mechanics

### 4.2 Individual Solutions

---

#### SOLUTION-01: Create Responsive Layout System
**For**: ISSUE-01  
**Approach**: Implement window size tracking and scale UI proportionally  
**Implementation Type**: Add new UILayout struct  
**Affected Files**: 
- `internal/render/layout.go` (new file)
- `internal/render/renderer.go` (modify)
- `cmd/game/main.go` (modify)

**Estimated Complexity**: Complex  
**Prerequisites**: None  
**Risk Level**: Low - Additive change, doesn't modify existing logic initially  

**STEP-BY-STEP IMPLEMENTATION:**

1. Create `internal/render/layout.go` with UILayout struct
2. Add methods: `GetScreenWidth()`, `GetScreenHeight()`, `ScaleX(val)`, `ScaleY(val)`
3. Modify `cmd/game/main.go` Layout() to return actual window dimensions
4. Pass UILayout instance to Renderer
5. Update Renderer to use scaled coordinates

**CODE CHANGES REQUIRED:**

File: `internal/render/layout.go` (new)
```go
package render

type UILayout struct {
    baseWidth  int // 960
    baseHeight int // 640
    actualWidth  int
    actualHeight int
    scaleX float64
    scaleY float64
}

func NewUILayout(baseWidth, baseHeight int) *UILayout {
    return &UILayout{
        baseWidth:  baseWidth,
        baseHeight: baseHeight,
        actualWidth:  baseWidth,
        actualHeight: baseHeight,
        scaleX: 1.0,
        scaleY: 1.0,
    }
}

func (ui *UILayout) Update(width, height int) {
    ui.actualWidth = width
    ui.actualHeight = height
    ui.scaleX = float64(width) / float64(ui.baseWidth)
    ui.scaleY = float64(height) / float64(ui.baseHeight)
}

func (ui *UILayout) ScaleX(x float64) float64 {
    return x * ui.scaleX
}

func (ui *UILayout) ScaleY(y float64) float64 {
    return y * ui.scaleY
}

func (ui *UILayout) GetScreenWidth() int {
    return ui.actualWidth
}

func (ui *UILayout) GetScreenHeight() int {
    return ui.actualHeight
}

func (ui *UILayout) CenterX(width float64) float64 {
    return (float64(ui.actualWidth) - width) / 2.0
}

func (ui *UILayout) CenterY(height float64) float64 {
    return (float64(ui.actualHeight) - height) / 2.0
}
```

File: `cmd/game/main.go`
- Change 1: Store window dimensions
  Current: `func (app *GameApp) Layout(outsideWidth, outsideHeight int) (int, int) { return 960, 640 }`
  New: Store dimensions and return them
  ```go
  func (app *GameApp) Layout(outsideWidth, outsideHeight int) (int, int) {
      // Update layout system with actual dimensions
      if app.gameRunner != nil {
          app.gameRunner.UpdateLayout(outsideWidth, outsideHeight)
      }
      return outsideWidth, outsideHeight
  }
  ```

**CONSTANTS TO DEFINE:**
```go
const (
    BaseScreenWidth  = 960
    BaseScreenHeight = 640
)
```

**TESTING CRITERIA:**
- [ ] Game renders at 960x640 unchanged
- [ ] Game renders correctly at 1920x1080
- [ ] Game renders correctly at 800x600
- [ ] UI elements scale proportionally
- [ ] No gameplay mechanics affected

**ROLLBACK PLAN:**
If issues arise: Revert Layout() to return `960, 640` and remove UILayout usage

---

#### SOLUTION-02: Separate Debug UI Layer
**For**: ISSUE-02  
**Approach**: Add debug UI toggle and render on separate layer  
**Implementation Type**: Refactor + Add feature  
**Affected Files**: 
- `internal/engine/runner.go`
- `internal/render/renderer.go`

**Estimated Complexity**: Simple  
**Prerequisites**: None  
**Risk Level**: Low  

**STEP-BY-STEP IMPLEMENTATION:**

1. Add `showDebugUI` bool field to GameRunner
2. Check for F3 key press to toggle debug UI
3. Move debug text rendering to separate function
4. Render debug UI last (on top) with semi-transparent background
5. Position debug text at bottom-left instead of top-left

**CODE CHANGES REQUIRED:**

File: `internal/engine/runner.go`
- Change 1: Add debug UI toggle field
  Current: No field for debug UI state
  New: Add to GameRunner struct
  ```go
  type GameRunner struct {
      // ... existing fields ...
      showDebugUI bool
  }
  ```

- Change 2: Toggle debug UI in Update()
  Current: Debug always shown
  New: Add toggle check
  ```go
  // In Update() method
  if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
      gr.showDebugUI = !gr.showDebugUI
  }
  ```

- Change 3: Conditional debug rendering
  Current: `ebitenutil.DebugPrint(screen, debugInfo)` at line 590
  New: Wrap in condition and move to bottom
  ```go
  if gr.showDebugUI {
      gr.renderDebugUI(screen, debugInfo)
  }
  ```

- Change 4: Create renderDebugUI() method
  New method:
  ```go
  func (gr *GameRunner) renderDebugUI(screen *ebiten.Image, info string) {
      // Draw semi-transparent background
      bgHeight := 100
      bgImg := ebiten.NewImage(ScreenWidth, bgHeight)
      bgImg.Fill(color.RGBA{0, 0, 0, 180})
      opts := &ebiten.DrawImageOptions{}
      opts.GeoM.Translate(0, float64(ScreenHeight-bgHeight))
      screen.DrawImage(bgImg, opts)
      
      // Draw debug text at bottom-left
      ebitenutil.DebugPrintAt(screen, info, 10, ScreenHeight-90)
  }
  ```

**TESTING CRITERIA:**
- [ ] F3 toggles debug UI on/off
- [ ] Debug UI starts hidden by default
- [ ] Debug UI renders at bottom with background
- [ ] Health bar and abilities visible without overlap
- [ ] Debug info remains readable

**ROLLBACK PLAN:**
Remove `showDebugUI` field and conditional; restore original DebugPrint call

---

#### SOLUTION-03: Create UI Constants File
**For**: ISSUE-04, ISSUE-06, ISSUE-07, ISSUE-08, ISSUE-10  
**Approach**: Centralize all UI dimension constants  
**Implementation Type**: Add constants file + refactor  
**Affected Files**:
- `internal/render/constants.go` (new)
- `internal/render/renderer.go` (modify)
- `internal/menu/menu.go` (modify)
- `internal/engine/runner.go` (modify)

**Estimated Complexity**: Moderate  
**Prerequisites**: None  
**Risk Level**: Low  

**STEP-BY-STEP IMPLEMENTATION:**

1. Create `internal/render/constants.go`
2. Define all UI dimension constants
3. Replace magic numbers in renderer.go
4. Replace magic numbers in menu.go
5. Replace magic numbers in runner.go

**CODE CHANGES REQUIRED:**

File: `internal/render/constants.go` (new)
```go
package render

const (
    // Screen dimensions (base)
    BaseScreenWidth  = 960
    BaseScreenHeight = 640
    
    // UI margins and padding
    UIMargin        = 10
    UISpacing       = 10
    UILargeSpacing  = 20
    
    // Health bar
    HealthBarWidth  = 200
    HealthBarHeight = 20
    HealthBarX      = UIMargin
    HealthBarY      = UIMargin
    
    // Ability icons
    AbilityIconSize    = 30
    AbilityIconSpacing = 5
    AbilityIconY       = HealthBarY + HealthBarHeight + UISpacing
    
    // Messages
    MessageBoxWidth  = 200
    MessageBoxHeight = 40
    MessagePadding   = 10
    
    // Menu
    MenuTitleY       = 100
    MenuItemStartY   = 200
    MenuItemSpacing  = 40
    MenuItemX        = 220
    MenuSelectorX    = 200
    MenuInstructionY = 500
    
    // Enemy health bars
    EnemyHealthBarHeight = 4
    EnemyHealthBarOffset = 8
    
    // Tile rendering
    TileSize = 32
)
```

File: `internal/render/renderer.go`
- Change 1: Replace health bar position
  Current: `barX := 10` and `barY := 10` (lines 238-239)
  New: `barX := HealthBarX` and `barY := HealthBarY`

- Change 2: Replace health bar dimensions
  Current: `barWidth := 200` and `barHeight := 20` (lines 236-237)
  New: `barWidth := HealthBarWidth` and `barHeight := HealthBarHeight`

- Change 3: Replace ability icon constants
  Current: Multiple magic numbers (lines 261-282)
  New: Use AbilityIconSize, AbilityIconSpacing, AbilityIconY

File: `internal/menu/menu.go`
- Change 1: Replace menu item spacing
  Current: `y := startY + i*40` (line 231)
  New: `y := render.MenuItemStartY + i*render.MenuItemSpacing`

File: `internal/engine/runner.go`
- Change 1: Replace message dimensions
  Current: `NewImage(200, 40)` (lines 533, 550)
  New: `NewImage(render.MessageBoxWidth, render.MessageBoxHeight)`

**TESTING CRITERIA:**
- [ ] All UI elements render in same positions as before
- [ ] No visual regressions
- [ ] Constants are accessible from all files
- [ ] Code is more readable

**ROLLBACK PLAN:**
Delete constants.go and restore hardcoded values

---

#### SOLUTION-04: Implement Proper Text Centering
**For**: ISSUE-03, ISSUE-13  
**Approach**: Use Ebiten text bounds or add custom font rendering  
**Implementation Type**: Add text utility functions  
**Affected Files**:
- `internal/render/text.go` (new)
- `internal/menu/menu.go` (modify)

**Estimated Complexity**: Moderate  
**Prerequisites**: Consider adding `golang.org/x/image/font` dependency  
**Risk Level**: Medium - May require external font library  

**STEP-BY-STEP IMPLEMENTATION:**

1. Research Ebiten text measurement options
2. Create text utility functions for centering
3. Replace approximation-based centering in menu
4. Add font loading if needed

**CODE CHANGES REQUIRED:**

File: `internal/render/text.go` (new)
```go
package render

import (
    "image"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
)

// EstimateTextWidth estimates text width using basicfont
func EstimateTextWidth(text string) int {
    // Using basicfont.Face7x13 as default
    bounds, _ := font.BoundString(basicfont.Face7x13, text)
    return (bounds.Max.X - bounds.Min.X).Ceil()
}

// CenterTextX calculates X position to center text
func CenterTextX(text string, screenWidth int) int {
    textWidth := EstimateTextWidth(text)
    return (screenWidth - textWidth) / 2
}
```

File: `internal/menu/menu.go`
- Change 1: Replace title centering
  Current: `titleX := 480 - len(title)*4` (line 224)
  New: `titleX := render.CenterTextX(title, 960)`

- Change 2: Replace instruction centering
  Current: `480-len(instructions)*3` (line 253)
  New: `render.CenterTextX(instructions, 960)`

**TESTING CRITERIA:**
- [ ] Menu titles appear centered
- [ ] Instructions appear centered
- [ ] Works for all menu types (main, pause, settings)
- [ ] No visual shift from previous approximate centering

**ROLLBACK PLAN:**
Remove text.go and restore length-based approximation

---

#### SOLUTION-05: Fix Message Box Centering
**For**: ISSUE-05, ISSUE-10  
**Approach**: Calculate center based on actual message dimensions  
**Implementation Type**: Refactor positioning logic  
**Affected Files**:
- `internal/engine/runner.go`

**Estimated Complexity**: Simple  
**Prerequisites**: SOLUTION-03 (constants)  
**Risk Level**: Low  

**STEP-BY-STEP IMPLEMENTATION:**

1. Use UILayout.CenterX() for message positioning
2. Make message box size dynamic based on text length
3. Add padding calculations

**CODE CHANGES REQUIRED:**

File: `internal/engine/runner.go`
- Change 1: Fix locked door message centering
  Current: `messageX := render.ScreenWidth/2 - 100` (line 529)
  New: 
  ```go
  messageWidth := render.MessageBoxWidth
  messageX := gr.renderer.layout.CenterX(float64(messageWidth))
  ```

- Change 2: Fix item message centering
  Current: `messageX := render.ScreenWidth/2 - 100` (line 546)
  New: Use same CenterX() approach

- Change 3: Dynamic message height (optional enhancement)
  Current: Fixed 40px height
  New: Calculate based on text length if multi-line

**TESTING CRITERIA:**
- [ ] Messages appear centered at all resolutions
- [ ] Messages remain centered when window resizes
- [ ] Text doesn't overflow message box
- [ ] Background properly encompasses text

**ROLLBACK PLAN:**
Restore `ScreenWidth/2 - 100` calculation

---

#### SOLUTION-06: Add Ability Icon Labels
**For**: ISSUE-11  
**Approach**: Add text labels below ability icons  
**Implementation Type**: Add feature  
**Affected Files**:
- `internal/render/renderer.go`

**Estimated Complexity**: Simple  
**Prerequisites**: SOLUTION-04 (text utilities)  
**Risk Level**: Low  

**STEP-BY-STEP IMPLEMENTATION:**

1. Define short ability names (J, D, W, G for Jump, Dash, Wall, Glide)
2. Render text below each icon
3. Adjust icon positioning to make room for labels

**CODE CHANGES REQUIRED:**

File: `internal/render/renderer.go`
- Change 1: Add ability name labels
  Location: After drawing ability icons (around line 283)
  New code:
  ```go
  // After drawing ability icons
  abilityLabels := []string{"DJ", "DS", "WJ", "GL"}
  for i, label := range abilityLabels {
      if i < len(abilityNames) {
          labelX := abilityX + i*(abilitySize+abilitySpacing) + abilitySize/2 - 6
          labelY := abilityY + abilitySize + 2
          ebitenutil.DebugPrintAt(screen, label, labelX, labelY)
      }
  }
  ```

**TESTING CRITERIA:**
- [ ] Labels appear below icons
- [ ] Labels aligned with icon centers
- [ ] Labels readable against background
- [ ] Locked abilities have dimmed labels

**ROLLBACK PLAN:**
Remove label rendering code

---

#### SOLUTION-07: Improve Menu Selection Indicator
**For**: ISSUE-12  
**Approach**: Add background highlight to selected item  
**Implementation Type**: Add visual enhancement  
**Affected Files**:
- `internal/menu/menu.go`

**Estimated Complexity**: Simple  
**Prerequisites**: None  
**Risk Level**: Low  

**STEP-BY-STEP IMPLEMENTATION:**

1. Draw colored background rectangle for selected menu item
2. Keep ">" indicator for additional clarity
3. Use contrasting color

**CODE CHANGES REQUIRED:**

File: `internal/menu/menu.go`
- Change 1: Add background to selected item
  Location: Before drawing menu items (around line 238)
  New code:
  ```go
  // Draw selection background
  if i == mm.selectedIndex {
      highlightWidth := 400
      highlightHeight := 35
      highlightImg := ebiten.NewImage(highlightWidth, highlightHeight)
      highlightImg.Fill(color.RGBA{50, 50, 100, 255}) // Blue highlight
      
      opts := &ebiten.DrawImageOptions{}
      opts.GeoM.Translate(float64(200), float64(y-5))
      screen.DrawImage(highlightImg, opts)
  }
  ```

**TESTING CRITERIA:**
- [ ] Selected item has visible background
- [ ] Background color contrasts with menu background
- [ ] ">" indicator still visible
- [ ] Navigation feels responsive

**ROLLBACK PLAN:**
Remove background drawing code

---

#### SOLUTION-08: Add Particle Layer Management
**For**: ISSUE-15  
**Approach**: Create particle layers (background, foreground)  
**Implementation Type**: Refactor particle system  
**Affected Files**:
- `internal/particle/particle.go`
- `internal/engine/runner.go`

**Estimated Complexity**: Moderate  
**Prerequisites**: None  
**Risk Level**: Low  

**STEP-BY-STEP IMPLEMENTATION:**

1. Add `Layer` field to Particle struct
2. Separate particles by layer in rendering
3. Render background particles before player
4. Render foreground particles after player

**CODE CHANGES REQUIRED:**

File: `internal/particle/particle.go`
- Change 1: Add Layer field
  ```go
  type ParticleLayer int
  
  const (
      LayerBackground ParticleLayer = iota
      LayerForeground
  )
  
  type Particle struct {
      // ... existing fields ...
      Layer ParticleLayer
  }
  ```

File: `internal/engine/runner.go`
- Change 1: Render particles in two passes
  Current: Single render call at line 499
  New: Two render calls
  ```go
  // Render background particles (dust, etc.)
  backgroundParticles := gr.particleSystem.GetParticlesByLayer(particle.LayerBackground)
  gr.renderer.RenderParticles(screen, backgroundParticles)
  
  // Render player
  gr.renderer.RenderPlayer(...)
  
  // Render foreground particles (hits, sparkles)
  foregroundParticles := gr.particleSystem.GetParticlesByLayer(particle.LayerForeground)
  gr.renderer.RenderParticles(screen, foregroundParticles)
  ```

**TESTING CRITERIA:**
- [ ] Jump dust appears behind player
- [ ] Hit effects appear in front of player
- [ ] Blood splatter renders correctly
- [ ] Sparkles visible on top

**ROLLBACK PLAN:**
Remove Layer field and render all particles in single pass

---

### 4.3 Dependency Graph

```
SOLUTION-01 (Responsive Layout)
    ↓
SOLUTION-03 (UI Constants) → SOLUTION-05 (Message Centering)
    ↓                           ↓
SOLUTION-02 (Debug Layer)   SOLUTION-06 (Ability Labels)
    ↓                           ↓
SOLUTION-04 (Text Centering)   SOLUTION-07 (Menu Selection)
    ↓
SOLUTION-08 (Particle Layers)
```

**BLOCKING SOLUTIONS (must be done first):**
- SOLUTION-01: Responsive layout system - Foundation for all responsive changes
- SOLUTION-03: UI constants - Required by multiple other solutions

**INDEPENDENT SOLUTIONS (can be done in parallel after foundation):**
- SOLUTION-02: Debug layer
- SOLUTION-06: Ability labels  
- SOLUTION-07: Menu selection highlight
- SOLUTION-08: Particle layers

---

## 5. Implementation Roadmap

### Phase 1: Foundation (Blocking changes, low risk)

**Goal**: Establish responsive layout system and UI constants  
**Duration Estimate**: 4-6 hours  
**Solutions**: SOLUTION-01, SOLUTION-03  
**Deliverable**: Layout system and constants in place, no visual changes yet  
**Verification**:
- [ ] UILayout struct exists and compiles
- [ ] Constants file exists with all values
- [ ] Game runs without visual regressions
- [ ] Window can be resized (even if UI doesn't scale yet)

**Detailed Steps:**
1. Create `internal/render/layout.go` with UILayout implementation
2. Create `internal/render/constants.go` with all UI constants
3. Modify `cmd/game/main.go` to track window dimensions
4. Add UILayout to Renderer struct
5. Run game and verify no crashes
6. Test at multiple resolutions

### Phase 2: Core Positioning Fixes (High priority issues)

**Goal**: Fix hardcoded positions and enable responsive UI  
**Duration Estimate**: 6-8 hours  
**Solutions**: SOLUTION-02, SOLUTION-04, SOLUTION-05  
**Deliverable**: UI elements scale with window; text properly centered; debug UI separated  
**Verification**:
- [ ] Health bar scales with window size
- [ ] Menu text centered accurately
- [ ] Messages centered at all resolutions
- [ ] F3 toggles debug UI
- [ ] Debug text doesn't overlap health bar

**Detailed Steps:**
1. Implement SOLUTION-02 (Debug UI layer)
2. Implement SOLUTION-04 (Text centering utilities)
3. Implement SOLUTION-05 (Message box centering)
4. Replace all hardcoded positions with scaled values
5. Test at 960x640, 1920x1080, 800x600
6. Verify gameplay mechanics unchanged

### Phase 3: Visual Polish (Medium priority improvements)

**Goal**: Improve visual clarity and usability  
**Duration Estimate**: 4-5 hours  
**Solutions**: SOLUTION-06, SOLUTION-07, SOLUTION-08  
**Deliverable**: Better visual feedback, clearer UI indicators  
**Verification**:
- [ ] Ability icons have labels
- [ ] Menu selection clearly visible
- [ ] Particles render in correct layers
- [ ] Overall visual consistency improved

**Detailed Steps:**
1. Implement SOLUTION-06 (Ability labels)
2. Implement SOLUTION-07 (Menu selection highlight)
3. Implement SOLUTION-08 (Particle layers)
4. Visual QA pass on all screens
5. Adjust colors/spacing as needed

### Phase 4: Future Enhancements (Low priority, optional)

**Goal**: Add nice-to-have features  
**Duration Estimate**: 6-8 hours  
**Solutions**: Address ISSUE-16, ISSUE-17, ISSUE-18  
**Deliverable**: Minimap, animated item glow, persistent settings  
**Verification**:
- [ ] Minimap shows current room and explored areas
- [ ] Item glow pulses attractively
- [ ] Settings persist across sessions

**Note**: This phase is optional and can be done after main UI issues are resolved.

---

## 6. Risk Assessment

### High-Risk Changes

**SOLUTION-01: Responsive Layout System**
- **Risk**: Could break existing positioning logic
- **Mitigation**: Implement with scaling factor of 1.0 initially; add gradual scaling
- **Testing**: Extensive testing at base resolution before enabling scaling
- **Rollback**: Single flag to disable responsive mode

**SOLUTION-04: Text Centering with Font Library**
- **Risk**: External dependency; may affect build process
- **Mitigation**: Use standard library if possible; fallback to approximation
- **Testing**: Test on multiple platforms (Linux, Windows, macOS)
- **Rollback**: Remove font import and use previous approximation

### Medium-Risk Changes

**SOLUTION-08: Particle Layer Management**
- **Risk**: Could affect particle performance
- **Mitigation**: Keep layer logic simple; no deep sorting
- **Testing**: Monitor FPS with many particles
- **Rollback**: Single-layer rendering easy to restore

**SOLUTION-05: Message Box Centering**
- **Risk**: Messages might appear in wrong locations
- **Mitigation**: Test all message types (door, item, etc.)
- **Testing**: Trigger all message scenarios
- **Rollback**: Restore hardcoded offsets

### Low-Risk Changes

**SOLUTION-02, 03, 06, 07**: All additive changes with minimal existing code modification
- **Risk**: Minimal; mostly adding new features
- **Mitigation**: Not needed
- **Testing**: Standard functional testing
- **Rollback**: Simple deletion of new code

### Mitigation Strategies

1. **Feature Flags**: Add flags to enable/disable new systems during testing
2. **Gradual Rollout**: Implement responsive layout per-element, not all at once
3. **Comprehensive Testing**: Test at multiple resolutions (960x640, 1280x720, 1920x1080, 800x600)
4. **Visual Comparison**: Take screenshots before/after each phase
5. **Performance Monitoring**: Track FPS during implementation
6. **Git Branching**: Use feature branches for risky changes
7. **Backup Constants**: Keep original magic numbers in comments initially

---

## 7. Testing & Validation

### Per-Phase Testing

**Phase 1 Testing:**
```
Resolution Tests:
  [ ] 960x640 (base) - UI renders identically to before
  [ ] 1920x1080 - Window opens, game runs, no crashes
  [ ] 800x600 - Window opens, game runs, no crashes
  [ ] Window resize - Can resize without crashes

Constant Tests:
  [ ] Health bar at (10, 10) with constants
  [ ] Ability icons properly spaced with constants
  [ ] Menu items evenly spaced with constants
  [ ] All magic numbers replaced

Regression Tests:
  [ ] Player movement works
  [ ] Combat system works
  [ ] Enemy AI works
  [ ] Item collection works
  [ ] Room transitions work
```

**Phase 2 Testing:**
```
Responsive UI Tests:
  [ ] Health bar scales proportionally at 1920x1080
  [ ] Ability icons scale proportionally at 1920x1080
  [ ] Health bar visible at 800x600 (no cutoff)
  [ ] Messages centered at all resolutions
  [ ] Menu text centered at all resolutions

Debug UI Tests:
  [ ] F3 toggles debug UI
  [ ] Debug UI hidden by default (or shown, as configured)
  [ ] Health bar fully visible with debug off
  [ ] Debug info readable at bottom of screen
  [ ] Debug background semi-transparent

Text Centering Tests:
  [ ] "VANIA - Procedural Metroidvania" centered
  [ ] "Settings" centered
  [ ] "Game Paused" centered
  [ ] Instructions text centered
  [ ] All menu titles aligned consistently
```

**Phase 3 Testing:**
```
Visual Polish Tests:
  [ ] Ability labels visible and aligned
  [ ] Ability labels update when abilities unlocked
  [ ] Menu selection has visible highlight
  [ ] Selection navigation smooth
  [ ] Dust particles behind player
  [ ] Hit effects in front of player
  [ ] Blood splatter renders correctly

Color/Contrast Tests:
  [ ] All text readable against backgrounds
  [ ] Selection highlight contrasts with menu
  [ ] Ability labels readable
  [ ] Message boxes have sufficient contrast
```

### Final Validation Criteria

- [ ] All UI elements visible and properly positioned at 960x640
- [ ] All UI elements visible and properly scaled at 1920x1080
- [ ] All UI elements visible and properly scaled at 800x600
- [ ] No overlapping elements (except intentional layering)
- [ ] Consistent spacing and alignment across all screens
- [ ] Responsive to screen size changes (scales appropriately)
- [ ] All text readable with proper contrast
- [ ] Visual hierarchy clear and logical (health bar prominent, debug info subtle)
- [ ] No gameplay mechanics affected
- [ ] FPS remains stable (no performance regression)
- [ ] Menu navigation intuitive
- [ ] All messages properly centered
- [ ] Ability icons clearly labeled
- [ ] Debug UI toggleable and non-intrusive

### Automated Testing (Future)

```go
// Example test structure
func TestUILayout_Responsive(t *testing.T) {
    layout := NewUILayout(960, 640)
    
    // Test scaling at 1920x1080
    layout.Update(1920, 1080)
    assert.Equal(t, 2.0, layout.scaleX)
    assert.Equal(t, 1.6875, layout.scaleY)
    
    // Test position scaling
    scaledX := layout.ScaleX(10) // Health bar X
    assert.Equal(t, 20.0, scaledX)
}

func TestTextCentering(t *testing.T) {
    text := "Test Title"
    centerX := CenterTextX(text, 960)
    
    // Should be reasonably centered
    assert.Greater(t, centerX, 400)
    assert.Less(t, centerX, 520)
}
```

---

## 8. Success Metrics

### Quantitative Metrics

1. **Code Quality**
   - Magic numbers reduced from 20+ to 0
   - UI constants centralized in 1 file
   - Lines of duplicated positioning code: 0

2. **Resolution Support**
   - Tested resolutions: 5+ (960x640, 1280x720, 1920x1080, 800x600, 2560x1440)
   - UI elements visible at all tested resolutions: 100%
   - Scaling accuracy: Within 2px of expected position

3. **Performance**
   - FPS at 960x640: No change from baseline
   - FPS at 1920x1080: Within 5% of baseline
   - Memory usage: No significant increase

4. **User Experience**
   - Menu text centering accuracy: Within 5px of perfect center
   - Debug UI toggle response time: <100ms
   - Selection indicator visibility improvement: Measured by user feedback

### Qualitative Metrics

1. **Visual Consistency**
   - All spacing follows defined constants
   - Color scheme consistent across UI elements
   - Typography consistent (limited by Ebiten's default font)

2. **Maintainability**
   - New UI elements can be added without hardcoding positions
   - UI adjustments require changing constants only
   - Code reviewers find positioning logic clear

3. **Player Feedback**
   - Players report UI clarity improved
   - No confusion about ability icons
   - Menu navigation feels responsive
   - Debug info no longer intrusive

### Acceptance Criteria

✅ **Must Have (All must pass):**
- All critical issues (ISSUE-01, ISSUE-02) resolved
- All high priority issues (ISSUE-03 through ISSUE-08) resolved
- No gameplay regressions
- Performance within acceptable range
- Code passes existing tests

✅ **Should Have (80% must pass):**
- Medium priority issues resolved
- Automated tests for UI layout
- Documentation updated
- Visual consistency metrics met

✅ **Nice to Have (50% pass acceptable):**
- Low priority issues resolved
- Additional features (minimap, animations)
- Platform-specific testing complete

---

## Appendices

### A. Code Location Reference

| Issue ID | File | Line(s) | Element |
|----------|------|---------|---------|
| ISSUE-01 | cmd/game/main.go | 91, runner.go:594-596 | Layout return value |
| ISSUE-02 | runner.go | 590 | Debug text |
| ISSUE-03 | menu.go | 224, 253 | Title/instruction centering |
| ISSUE-04 | renderer.go | 238-239 | Health bar position |
| ISSUE-05 | runner.go | 529, 546 | Message centering |
| ISSUE-06 | renderer.go, menu.go, runner.go | Multiple | Magic numbers |
| ISSUE-07 | renderer.go | 261 | Ability icon spacing |
| ISSUE-08 | menu.go | 231 | Menu item spacing |
| ISSUE-09 | renderer.go | 302-310, 344 | Enemy health bars |
| ISSUE-10 | runner.go | 533-534, 550-551 | Message box size |
| ISSUE-11 | renderer.go | 260-284 | Ability icons |
| ISSUE-12 | menu.go | 242-244 | Selection indicator |
| ISSUE-13 | menu.go | 264-278 | Menu title position |
| ISSUE-14 | renderer.go | 341-362 | Enemy health overlap |
| ISSUE-15 | runner.go | 499 | Particle render order |
| ISSUE-16 | N/A | N/A | Missing minimap |
| ISSUE-17 | renderer.go | 462-469 | Item glow |
| ISSUE-18 | menu.go | 396-465 | Settings persistence |

### B. Glossary

**Ebiten**: Go game engine used by VANIA. Version 2.6.3.

**Layout()**: Ebiten interface method that defines logical screen dimensions. Currently returns fixed 960x640.

**Draw()**: Ebiten interface method called every frame to render graphics.

**Update()**: Ebiten interface method called every tick (60 TPS) for game logic.

**DebugPrint**: Ebiten utility function for rendering text. Limited to white text, no sizing.

**Screen Space**: Coordinate system where (0,0) is top-left of screen, used for UI.

**World Space**: Coordinate system for game objects, transformed by camera offset.

**Magic Number**: Hardcoded numeric value without named constant, making code hard to maintain.

**Responsive Layout**: UI that adapts to different screen sizes and resolutions.

**UILayout**: Proposed struct to manage responsive positioning and scaling.

**TPS**: Ticks Per Second. Ebiten runs at 60 TPS for game logic.

**FPS**: Frames Per Second. Target is 60 FPS for rendering.

**HUD**: Heads-Up Display. On-screen UI elements during gameplay (health, abilities).

**Render Order**: Sequence in which elements are drawn. Later draws appear on top.

**Camera Offset**: Translation applied to world space coordinates for viewport scrolling.

**Particle Layer**: Background or foreground classification for particle effects.

---

## Implementation Notes

### Development Environment Setup

```bash
# Clone repository
git clone https://github.com/opd-ai/vania.git
cd vania

# Install dependencies
go mod download

# Build game
go build -o vania ./cmd/game

# Run game
./vania --play

# Run with specific seed
./vania --play --seed 42

# Run in stats-only mode (no rendering)
./vania --stats-only --seed 42
```

### Code Style Guidelines

- Follow existing Go conventions in the codebase
- Use `const` blocks for related constants
- Add comments for non-obvious calculations
- Keep functions focused and under 50 lines when possible
- Use descriptive variable names (avoid single letters except loops)

### Git Workflow

```bash
# Create feature branch for each phase
git checkout -b feature/ui-foundation
git checkout -b feature/ui-positioning
git checkout -b feature/ui-polish

# Commit after each solution
git commit -m "feat(ui): Add responsive layout system (SOLUTION-01)"

# Tag each phase completion
git tag phase-1-complete
```

### Testing Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/render

# Run with verbose output
go test -v ./internal/render

# Build without running
go build ./cmd/game
```

---

**Document Version**: 1.0  
**Last Updated**: 2025-01-26  
**Author**: Code Analysis System  
**Status**: Ready for Implementation  
**Estimated Total Time**: 16-24 hours across 4 phases

