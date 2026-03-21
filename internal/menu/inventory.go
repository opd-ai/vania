// Package menu provides inventory UI for displaying, navigating, and using
// items collected by the player.
package menu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/render"
)

const (
	// Inventory grid dimensions
	inventoryCols = 4
	inventoryRows = 8

	// Cell dimensions in pixels
	inventoryCellSize = 40
	inventoryPadding  = 8

	// Tooltip display duration in frames at 60 fps
	tooltipDuration = 180
)

// inventoryPanel positions
const (
	inventoryPanelX = 80
	inventoryPanelY = 60
)

// InventoryScreen renders and handles input for the player inventory.
// It displays a 4×8 item grid with item tooltips and supports consumable use.
type InventoryScreen struct {
	items         []*entity.Item
	abilities     []entity.Ability
	selectedRow   int
	selectedCol   int
	tooltipTimer  int
	useConfirm    bool // waiting for confirmation to use a consumable
	renderer      *render.BitmapTextRenderer
	genre         string
}

// NewInventoryScreen creates an InventoryScreen for rendering.
func NewInventoryScreen() *InventoryScreen {
	return &InventoryScreen{
		renderer: render.NewBitmapTextRenderer(),
		genre:    "fantasy",
	}
}

// SetGenre updates the UI vocabulary for the active genre.
func (is *InventoryScreen) SetGenre(genreID string) {
	is.genre = genreID
}

// SetItems replaces the item list with the player's current inventory.
func (is *InventoryScreen) SetItems(items []*entity.Item) {
	is.items = items
	// Clamp selection to new list length
	maxRow := (len(is.items) - 1) / inventoryCols
	if is.selectedRow > maxRow {
		is.selectedRow = maxRow
	}
}

// SetAbilities replaces the abilities list for the equipment-slots section.
func (is *InventoryScreen) SetAbilities(abilities []entity.Ability) {
	is.abilities = abilities
}

// SelectedIndex returns the flat index of the currently highlighted cell.
func (is *InventoryScreen) SelectedIndex() int {
	return is.selectedRow*inventoryCols + is.selectedCol
}

// SelectedItem returns the item at the cursor, or nil if the slot is empty.
func (is *InventoryScreen) SelectedItem() *entity.Item {
	idx := is.SelectedIndex()
	if idx < 0 || idx >= len(is.items) {
		return nil
	}
	return is.items[idx]
}

// Update processes one frame of input.
// Returns an ItemUseEvent when the player confirms consuming an item.
func (is *InventoryScreen) Update() *ItemUseEvent {
	if is.tooltipTimer > 0 {
		is.tooltipTimer--
	}

	// Navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if is.selectedRow > 0 {
			is.selectedRow--
		}
		is.useConfirm = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if is.selectedRow < inventoryRows-1 {
			is.selectedRow++
		}
		is.useConfirm = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if is.selectedCol > 0 {
			is.selectedCol--
		}
		is.useConfirm = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if is.selectedCol < inventoryCols-1 {
			is.selectedCol++
		}
		is.useConfirm = false
	}

	// Show tooltip on hover (reset timer when cursor moves)
	is.tooltipTimer = tooltipDuration

	item := is.SelectedItem()
	if item == nil {
		return nil
	}

	// Confirm use for consumable items
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		if item.Type == entity.ConsumableItem {
			if is.useConfirm {
				// Second press — execute use
				is.useConfirm = false
				return &ItemUseEvent{Item: item, Index: is.SelectedIndex()}
			}
			is.useConfirm = true
		}
	}
	// Cancel confirm
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		is.useConfirm = false
	}
	return nil
}

// ItemUseEvent is returned by Update when the player has confirmed using an item.
type ItemUseEvent struct {
	Item  *entity.Item
	Index int
}

// Draw renders the inventory grid, tooltips, and abilities section.
func (is *InventoryScreen) Draw(screen *ebiten.Image) {
	is.drawBackground(screen)
	is.drawTitle(screen)
	is.drawGrid(screen)
	is.drawTooltip(screen)
	is.drawAbilitiesSection(screen)
}

// drawBackground draws a semi-transparent panel behind the inventory.
func (is *InventoryScreen) drawBackground(screen *ebiten.Image) {
	panelW := inventoryCols*inventoryCellSize + (inventoryCols+1)*inventoryPadding + 200
	panelH := inventoryRows*inventoryCellSize + (inventoryRows+1)*inventoryPadding + 80
	ebitenutil.DrawRect(screen,
		float64(inventoryPanelX-inventoryPadding),
		float64(inventoryPanelY-inventoryPadding),
		float64(panelW),
		float64(panelH),
		color.RGBA{0, 0, 0, 200})
}

// drawTitle renders the section heading.
func (is *InventoryScreen) drawTitle(screen *ebiten.Image) {
	title := is.genreInventoryTitle()
	titleX := inventoryPanelX
	titleY := inventoryPanelY - 2
	is.drawText(screen, title, titleX, titleY, color.RGBA{255, 220, 100, 255})
}

// genreInventoryTitle returns a genre-appropriate heading.
func (is *InventoryScreen) genreInventoryTitle() string {
	titles := map[string]string{
		"fantasy":   "INVENTORY",
		"scifi":     "CARGO BAY",
		"horror":    "BELONGINGS",
		"cyberpunk": "DATA CACHE",
		"postapoc":  "SCAVENGED ITEMS",
	}
	if t, ok := titles[is.genre]; ok {
		return t
	}
	return "INVENTORY"
}

// drawGrid renders each inventory cell with its item (if any) and the cursor.
func (is *InventoryScreen) drawGrid(screen *ebiten.Image) {
	for row := 0; row < inventoryRows; row++ {
		for col := 0; col < inventoryCols; col++ {
			cellX := inventoryPanelX + col*(inventoryCellSize+inventoryPadding)
			cellY := inventoryPanelY + 16 + row*(inventoryCellSize+inventoryPadding)

			// Cell background
			cellColor := color.RGBA{40, 40, 60, 220}
			if row == is.selectedRow && col == is.selectedCol {
				cellColor = color.RGBA{80, 80, 140, 240}
			}
			ebitenutil.DrawRect(screen,
				float64(cellX), float64(cellY),
				float64(inventoryCellSize), float64(inventoryCellSize),
				cellColor)

			// Cell border
			borderColor := color.RGBA{100, 100, 150, 200}
			if row == is.selectedRow && col == is.selectedCol {
				borderColor = color.RGBA{200, 200, 255, 255}
			}
			is.drawBorder(screen, cellX, cellY, inventoryCellSize, inventoryCellSize, borderColor)

			// Item indicator
			idx := row*inventoryCols + col
			if idx < len(is.items) {
				item := is.items[idx]
				dotColor := itemTypeColor(item.Type)
				ebitenutil.DrawRect(screen,
					float64(cellX+2), float64(cellY+2),
					float64(inventoryCellSize-4), float64(inventoryCellSize-4),
					dotColor)
				// Item initial letter
				if len(item.Name) > 0 {
					is.drawText(screen, string([]rune(item.Name)[:1]),
						cellX+inventoryCellSize/2-4,
						cellY+inventoryCellSize/2-6,
						color.RGBA{255, 255, 255, 255})
				}
			}
		}
	}
}

// drawTooltip shows the selected item's name, type, and description.
func (is *InventoryScreen) drawTooltip(screen *ebiten.Image) {
	item := is.SelectedItem()
	if item == nil || is.tooltipTimer <= 0 {
		return
	}

	tipX := inventoryPanelX + inventoryCols*(inventoryCellSize+inventoryPadding) + 12
	tipY := inventoryPanelY + 20

	ebitenutil.DrawRect(screen, float64(tipX-4), float64(tipY-4), 188, 120, color.RGBA{20, 20, 40, 230})

	is.drawText(screen, item.Name, tipX, tipY, color.RGBA{255, 220, 100, 255})
	is.drawText(screen, fmt.Sprintf("[%s]", itemTypeName(item.Type)), tipX, tipY+14, color.RGBA{180, 180, 220, 255})
	is.drawWrappedText(screen, item.Description, tipX, tipY+30, 180, color.RGBA{220, 220, 220, 255})

	if item.Effect != "" {
		is.drawText(screen, "Effect: "+item.Effect, tipX, tipY+70, color.RGBA{100, 255, 100, 255})
	}

	if item.Type == entity.ConsumableItem {
		prompt := "Press Z/J to use"
		if is.useConfirm {
			prompt = "Confirm? Press again"
		}
		is.drawText(screen, prompt, tipX, tipY+88, color.RGBA{255, 180, 60, 255})
	}
}

// drawAbilitiesSection renders unlocked abilities below the grid.
func (is *InventoryScreen) drawAbilitiesSection(screen *ebiten.Image) {
	sectionX := inventoryPanelX
	sectionY := inventoryPanelY + 16 + inventoryRows*(inventoryCellSize+inventoryPadding) + 8

	is.drawText(screen, "ABILITIES", sectionX, sectionY, color.RGBA{180, 220, 255, 255})
	if len(is.abilities) == 0 {
		is.drawText(screen, "None unlocked", sectionX, sectionY+14, color.RGBA{140, 140, 160, 255})
		return
	}
	for i, ab := range is.abilities {
		col := i % inventoryCols
		row := i / inventoryCols
		ax := sectionX + col*100
		ay := sectionY + 14 + row*14
		is.drawText(screen, "• "+ab.Name, ax, ay, color.RGBA{200, 255, 200, 255})
	}
}

// drawText is a convenience wrapper around the bitmap renderer.
func (is *InventoryScreen) drawText(screen *ebiten.Image, text string, x, y int, col color.Color) {
	is.renderer.DrawText(screen, text, x, y, col)
}

// drawWrappedText renders text wrapping at maxWidth pixels.
func (is *InventoryScreen) drawWrappedText(screen *ebiten.Image, text string, x, y, maxWidth int, col color.Color) {
	charsPerLine := maxWidth / CharWidth
	if charsPerLine < 1 {
		charsPerLine = 1
	}
	line := ""
	lineY := y
	for _, ch := range text {
		line += string(ch)
		if len(line) >= charsPerLine {
			is.drawText(screen, line, x, lineY, col)
			line = ""
			lineY += CharHeight + 2
		}
	}
	if line != "" {
		is.drawText(screen, line, x, lineY, col)
	}
}

// drawBorder draws a 1-pixel rectangular outline.
func (is *InventoryScreen) drawBorder(screen *ebiten.Image, x, y, w, h int, col color.Color) {
	// Top
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), 1, col)
	// Bottom
	ebitenutil.DrawRect(screen, float64(x), float64(y+h-1), float64(w), 1, col)
	// Left
	ebitenutil.DrawRect(screen, float64(x), float64(y), 1, float64(h), col)
	// Right
	ebitenutil.DrawRect(screen, float64(x+w-1), float64(y), 1, float64(h), col)
}

// itemTypeColor returns a colour hint for each item category.
func itemTypeColor(t entity.ItemType) color.Color {
	switch t {
	case entity.WeaponItem:
		return color.RGBA{180, 60, 60, 180}
	case entity.ConsumableItem:
		return color.RGBA{60, 180, 60, 180}
	case entity.KeyItem:
		return color.RGBA{60, 60, 180, 180}
	case entity.UpgradeItem:
		return color.RGBA{180, 140, 60, 180}
	case entity.CurrencyItem:
		return color.RGBA{220, 200, 60, 180}
	default:
		return color.RGBA{100, 100, 100, 180}
	}
}

// itemTypeName returns a display label for each item category.
func itemTypeName(t entity.ItemType) string {
	switch t {
	case entity.WeaponItem:
		return "Weapon"
	case entity.ConsumableItem:
		return "Consumable"
	case entity.KeyItem:
		return "Key Item"
	case entity.UpgradeItem:
		return "Upgrade"
	case entity.CurrencyItem:
		return "Currency"
	default:
		return "Item"
	}
}
