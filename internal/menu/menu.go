// Package menu provides the main menu, pause menu, and settings menu systems
// for the VANIA game, handling navigation and user interface for game state
// management outside of gameplay.
package menu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/opd-ai/vania/internal/input"
	"github.com/opd-ai/vania/internal/save"
	settingspkg "github.com/opd-ai/vania/internal/settings"
)

const (
	// Menu Layout Constants
	ScreenWidth  = 960
	ScreenHeight = 640
	
	// Font constants
	CharWidth  = 8
	CharHeight = 12
	
	// Menu positioning
	MenuTitleY      = 100
	MenuStartY      = 200
	MenuItemSpacing = 40
	InstructionsY   = 500
)

// MenuType represents different menu types
type MenuType int

const (
	MainMenu MenuType = iota
	PauseMenu
	SettingsMenu
	SaveLoadMenu
	GameOverMenu
)

// MenuState represents current menu state
type MenuState int

const (
	MenuStateActive MenuState = iota
	MenuStateInactive
	MenuStateTransitioning
)

// MenuItem represents a selectable menu item
type MenuItem struct {
	Text     string
	Action   func() error
	Enabled  bool
	Selected bool
}

// MenuManager handles all menu operations
type MenuManager struct {
	currentMenu   MenuType
	state         MenuState
	items         []*MenuItem
	selectedIndex int
	inputHandler  *input.InputHandler
	saveManager   *save.SaveManager

	// Callbacks
	onNewGame    func(seed int64) error
	onLoadGame   func(slot int) error
	onSettings   func() error
	onQuitGame   func() error
	onResumeGame func() error

	// Settings
	settings        *GameSettings
	settingsManager *settingspkg.SettingsManager

	// Visual properties
	backgroundColor color.Color
	textColor       color.Color
	selectedColor   color.Color
	disabledColor   color.Color
}

// GameSettings holds user configurable settings
type GameSettings struct {
	MasterVolume float64
	SFXVolume    float64
	MusicVolume  float64
	FullScreen   bool
	VSync        bool
	ShowFPS      bool
	KeyBindings  map[string][]ebiten.Key
}

// NewMenuManager creates a new menu manager
func NewMenuManager() *MenuManager {
	settings := &GameSettings{
		MasterVolume: 0.7,
		SFXVolume:    0.8,
		MusicVolume:  0.6,
		FullScreen:   false,
		VSync:        true,
		ShowFPS:      false,
		KeyBindings: map[string][]ebiten.Key{
			"move_left":  {ebiten.KeyA, ebiten.KeyArrowLeft},
			"move_right": {ebiten.KeyD, ebiten.KeyArrowRight},
			"jump":       {ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp},
			"attack":     {ebiten.KeyJ, ebiten.KeyZ},
			"dash":       {ebiten.KeyK, ebiten.KeyX},
			"pause":      {ebiten.KeyP, ebiten.KeyEscape},
		},
	}

	saveManager, _ := save.NewSaveManager("")
	settingsManager := settingspkg.NewSettingsManager()

	return &MenuManager{
		currentMenu:     MainMenu,
		state:           MenuStateActive,
		items:           make([]*MenuItem, 0),
		selectedIndex:   0,
		inputHandler:    input.NewInputHandler(),
		saveManager:     saveManager,
		settings:        settings,
		settingsManager: settingsManager,
		backgroundColor: color.RGBA{20, 20, 30, 255},
		textColor:       color.RGBA{200, 200, 200, 255},
		selectedColor:   color.RGBA{255, 255, 100, 255},
		disabledColor:   color.RGBA{100, 100, 100, 255},
	}
}

// SetCallbacks sets callback functions for menu actions
func (mm *MenuManager) SetCallbacks(
	onNewGame func(seed int64) error,
	onLoadGame func(slot int) error,
	onSettings func() error,
	onQuitGame func() error,
	onResumeGame func() error,
) {
	mm.onNewGame = onNewGame
	mm.onLoadGame = onLoadGame
	mm.onSettings = onSettings
	mm.onQuitGame = onQuitGame
	mm.onResumeGame = onResumeGame
}

// ShowMainMenu displays the main menu
func (mm *MenuManager) ShowMainMenu() {
	mm.currentMenu = MainMenu
	mm.state = MenuStateActive
	mm.selectedIndex = 0
	mm.buildMainMenuItems()
}

// ShowPauseMenu displays the pause menu
func (mm *MenuManager) ShowPauseMenu() {
	mm.currentMenu = PauseMenu
	mm.state = MenuStateActive
	mm.selectedIndex = 0
	mm.buildPauseMenuItems()
}

// ShowSettingsMenu displays the settings menu
func (mm *MenuManager) ShowSettingsMenu() {
	mm.currentMenu = SettingsMenu
	mm.state = MenuStateActive
	mm.selectedIndex = 0
	mm.buildSettingsMenuItems()
}

// ShowGameOverMenu displays the game over menu
func (mm *MenuManager) ShowGameOverMenu() {
	mm.currentMenu = GameOverMenu
	mm.state = MenuStateActive
	mm.selectedIndex = 0
	mm.buildGameOverMenuItems()
}

// Hide hides the current menu
func (mm *MenuManager) Hide() {
	mm.state = MenuStateInactive
}

// IsActive returns true if menu is currently active
func (mm *MenuManager) IsActive() bool {
	return mm.state == MenuStateActive
}

// GetCurrentMenu returns the current menu type
func (mm *MenuManager) GetCurrentMenu() MenuType {
	return mm.currentMenu
}

// Update handles menu input and updates
func (mm *MenuManager) Update() error {
	if mm.state != MenuStateActive {
		return nil
	}

	inputState := mm.inputHandler.Update()

	// Handle navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		mm.navigateUp()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		mm.navigateDown()
	}

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return mm.selectCurrentItem()
	}

	// Handle back/escape
	if inputState.PausePress {
		return mm.handleBack()
	}

	return nil
}

// Draw renders the current menu
func (mm *MenuManager) Draw(screen *ebiten.Image) {
	if mm.state != MenuStateActive {
		return
	}

	// Clear screen with background
	screen.Fill(mm.backgroundColor)

	// Draw menu title with proper centering
	title := mm.getMenuTitle()
	titleWidth := len(title) * CharWidth
	titleX := (ScreenWidth - titleWidth) / 2
	titleY := MenuTitleY
	mm.drawColoredText(screen, title, titleX, titleY, mm.textColor)

	// Draw menu items
	startY := MenuStartY
	for i, item := range mm.items {
		y := startY + i*MenuItemSpacing

		// Choose color based on state
		itemColor := mm.textColor
		if !item.Enabled {
			itemColor = mm.disabledColor
		} else if i == mm.selectedIndex {
			itemColor = mm.selectedColor
		}

		// Draw selection indicator
		if i == mm.selectedIndex {
			ebitenutil.DebugPrintAt(screen, ">", 200, y)
		}

		// Draw item text
		mm.drawColoredText(screen, item.Text, 220, y, itemColor)
	}

	// Draw instructions with proper centering
	instructions := "Use W/S or Arrow Keys to navigate, Enter to select, Esc to back"
	instructWidth := len(instructions) * CharWidth
	instructX := (ScreenWidth - instructWidth) / 2
	instructY := InstructionsY
	mm.drawColoredText(screen, instructions, instructX, instructY, mm.disabledColor)
}

// drawColoredText draws text with specified color using procedural bitmap font
func (mm *MenuManager) drawColoredText(screen *ebiten.Image, text string, x, y int, col color.Color) {
	for i, char := range text {
		charX := x + i*CharWidth
		mm.drawChar(screen, char, charX, y, CharWidth, CharHeight, col)
	}
}

// drawChar draws a single character using procedural bitmap rendering
func (mm *MenuManager) drawChar(screen *ebiten.Image, char rune, x, y, width, height int, col color.Color) {
	// Create character bitmap based on simple patterns
	charImg := ebiten.NewImage(width, height)
	pixels := mm.getCharPattern(char, width, height)
	
	// Draw character pixels
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			if py < len(pixels) && px < len(pixels[py]) && pixels[py][px] {
				pixelImg := ebiten.NewImage(1, 1)
				pixelImg.Fill(col)
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(px), float64(py))
				charImg.DrawImage(pixelImg, opts)
			}
		}
	}
	
	// Draw character to screen
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(charImg, opts)
}

// getCharPattern returns a bitmap pattern for common characters
func (mm *MenuManager) getCharPattern(char rune, width, height int) [][]bool {
	// Initialize empty pattern
	pattern := make([][]bool, height)
	for i := range pattern {
		pattern[i] = make([]bool, width)
	}
	
	// Simple 8x12 bitmap patterns for basic ASCII characters
	switch char {
	case 'A', 'a':
		// A pattern
		pattern[1][3] = true; pattern[1][4] = true
		pattern[2][2] = true; pattern[2][5] = true
		pattern[3][2] = true; pattern[3][5] = true
		pattern[4][1] = true; pattern[4][6] = true
		pattern[5][1] = true; pattern[5][6] = true
		pattern[6][1] = true; pattern[6][2] = true; pattern[6][5] = true; pattern[6][6] = true
		pattern[7][1] = true; pattern[7][6] = true
		pattern[8][1] = true; pattern[8][6] = true
		
	case 'B', 'b':
		// B pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true }
		pattern[1][2] = true; pattern[1][3] = true; pattern[1][4] = true; pattern[1][5] = true
		pattern[2][6] = true; pattern[3][6] = true; pattern[4][6] = true
		pattern[5][2] = true; pattern[5][3] = true; pattern[5][4] = true; pattern[5][5] = true
		pattern[6][6] = true; pattern[7][6] = true; pattern[8][6] = true
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		
	case 'C', 'c':
		// C pattern
		pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
		pattern[3][1] = true; pattern[3][6] = true
		for y := 4; y < 7; y++ { pattern[y][1] = true }
		pattern[7][1] = true; pattern[7][6] = true
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		
	case 'E', 'e':
		// E pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true }
		for x := 1; x < 7; x++ { pattern[1][x] = true; pattern[5][x] = true; pattern[8][x] = true }
		
	case 'G', 'g':
		// G pattern
		pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
		pattern[3][1] = true; pattern[4][1] = true; pattern[5][1] = true; pattern[6][1] = true
		pattern[7][1] = true; pattern[7][6] = true
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		pattern[5][4] = true; pattern[5][5] = true; pattern[5][6] = true
		pattern[6][6] = true
		
	case 'L', 'l':
		// L pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true }
		for x := 1; x < 7; x++ { pattern[8][x] = true }
		
	case 'M', 'm':
		// M pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[2][2] = true; pattern[2][5] = true
		pattern[3][3] = true; pattern[3][4] = true
		
	case 'N', 'n':
		// N pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[2][2] = true; pattern[3][2] = true; pattern[4][3] = true
		pattern[5][4] = true; pattern[6][5] = true; pattern[7][5] = true
		
	case 'O', 'o':
		// O pattern
		pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
		pattern[3][1] = true; pattern[3][6] = true
		for y := 4; y < 7; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[7][1] = true; pattern[7][6] = true
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		
	case 'P', 'p':
		// P pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true }
		for x := 2; x < 6; x++ { pattern[1][x] = true; pattern[5][x] = true }
		pattern[2][6] = true; pattern[3][6] = true; pattern[4][6] = true
		
	case 'Q', 'q':
		// Q pattern (O with tail)
		pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
		pattern[3][1] = true; pattern[3][6] = true
		for y := 4; y < 7; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[6][4] = true // Inner diagonal
		pattern[7][1] = true; pattern[7][5] = true; pattern[7][6] = true
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][6] = true
		
	case 'R', 'r':
		// R pattern
		for y := 1; y < 9; y++ { pattern[y][1] = true }
		for x := 2; x < 6; x++ { pattern[1][x] = true; pattern[5][x] = true }
		pattern[2][6] = true; pattern[3][6] = true; pattern[4][6] = true
		pattern[6][4] = true; pattern[7][5] = true; pattern[8][6] = true
		
	case 'S', 's':
		// S pattern
		pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
		pattern[3][1] = true; pattern[4][1] = true
		pattern[5][2] = true; pattern[5][3] = true; pattern[5][4] = true; pattern[5][5] = true
		pattern[6][6] = true; pattern[7][6] = true
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		
	case 'T', 't':
		// T pattern
		for x := 1; x < 7; x++ { pattern[1][x] = true }
		for y := 1; y < 9; y++ { pattern[y][3] = true }
		
	case 'U', 'u':
		// U pattern
		for y := 1; y < 8; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		
	case 'V', 'v':
		// V pattern
		for y := 1; y < 7; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[7][2] = true; pattern[7][5] = true
		pattern[8][3] = true; pattern[8][4] = true
		
	case 'W', 'w':
		// W pattern
		for y := 1; y < 8; y++ { pattern[y][1] = true; pattern[y][6] = true }
		pattern[6][3] = true; pattern[6][4] = true
		pattern[7][2] = true; pattern[7][3] = true; pattern[7][4] = true; pattern[7][5] = true
		
	case 'X', 'x':
		// X pattern
		pattern[1][1] = true; pattern[1][6] = true
		pattern[2][2] = true; pattern[2][5] = true
		pattern[3][3] = true; pattern[3][4] = true
		pattern[4][3] = true; pattern[4][4] = true
		pattern[5][2] = true; pattern[5][5] = true
		pattern[6][1] = true; pattern[6][6] = true
		
	case 'Y', 'y':
		// Y pattern
		pattern[1][1] = true; pattern[1][6] = true
		pattern[2][2] = true; pattern[2][5] = true
		pattern[3][3] = true; pattern[3][4] = true
		for y := 4; y < 9; y++ { pattern[y][3] = true }
		
	case 'Z', 'z':
		// Z pattern
		for x := 1; x < 7; x++ { pattern[1][x] = true; pattern[8][x] = true }
		pattern[2][6] = true; pattern[3][5] = true; pattern[4][4] = true
		pattern[5][3] = true; pattern[6][2] = true; pattern[7][1] = true
		
	case ' ':
		// Space - already empty
		
	case ':':
		// Colon
		pattern[3][3] = true; pattern[6][3] = true
		
	case '-':
		// Hyphen
		pattern[5][2] = true; pattern[5][3] = true; pattern[5][4] = true; pattern[5][5] = true
		
	case '.':
		// Period
		pattern[8][3] = true
		
	case '!':
		// Exclamation
		for y := 1; y < 7; y++ { pattern[y][3] = true }
		pattern[8][3] = true
		
	case '?':
		// Question mark
		pattern[1][2] = true; pattern[1][3] = true; pattern[1][4] = true; pattern[1][5] = true
		pattern[2][6] = true; pattern[3][6] = true
		pattern[4][4] = true; pattern[4][5] = true
		pattern[5][3] = true
		pattern[8][3] = true
		
	case '(':
		// Left parenthesis
		pattern[2][4] = true; pattern[3][3] = true
		for y := 4; y < 7; y++ { pattern[y][2] = true }
		pattern[7][3] = true; pattern[8][4] = true
		
	case ')':
		// Right parenthesis  
		pattern[2][3] = true; pattern[3][4] = true
		for y := 4; y < 7; y++ { pattern[y][5] = true }
		pattern[7][4] = true; pattern[8][3] = true
		
	case '/':
		// Forward slash
		pattern[1][6] = true; pattern[2][5] = true; pattern[3][5] = true
		pattern[4][4] = true; pattern[5][3] = true; pattern[6][3] = true
		pattern[7][2] = true; pattern[8][1] = true
		
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Numbers
		switch char {
		case '0':
			pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
			for y := 3; y < 8; y++ { pattern[y][1] = true; pattern[y][6] = true }
			pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		case '1':
			pattern[2][3] = true; pattern[3][2] = true; pattern[3][3] = true
			for y := 4; y < 9; y++ { pattern[y][3] = true }
			for x := 1; x < 7; x++ { pattern[8][x] = true }
		case '2':
			pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
			pattern[3][6] = true; pattern[4][6] = true; pattern[5][5] = true
			pattern[6][4] = true; pattern[7][3] = true; pattern[8][2] = true
			for x := 1; x < 7; x++ { pattern[8][x] = true }
		case '3':
			for x := 2; x < 6; x++ { pattern[2][x] = true; pattern[5][x] = true; pattern[8][x] = true }
			pattern[3][6] = true; pattern[4][6] = true
			pattern[6][6] = true; pattern[7][6] = true
		case '4':
			for y := 2; y < 6; y++ { pattern[y][1] = true; pattern[y][5] = true }
			for x := 1; x < 7; x++ { pattern[5][x] = true }
			for y := 6; y < 9; y++ { pattern[y][5] = true }
		case '5':
			for x := 1; x < 7; x++ { pattern[2][x] = true; pattern[5][x] = true }
			for y := 3; y < 5; y++ { pattern[y][1] = true }
			pattern[6][6] = true; pattern[7][6] = true
			for x := 2; x < 6; x++ { pattern[8][x] = true }
		case '6':
			pattern[2][3] = true; pattern[2][4] = true; pattern[3][2] = true
			for y := 4; y < 8; y++ { pattern[y][1] = true }
			pattern[5][2] = true; pattern[5][3] = true; pattern[5][4] = true; pattern[5][5] = true
			pattern[6][6] = true; pattern[7][6] = true
			pattern[8][2] = true; pattern[8][3] = true; pattern[8][4] = true; pattern[8][5] = true
		case '7':
			for x := 1; x < 7; x++ { pattern[2][x] = true }
			pattern[3][6] = true; pattern[4][5] = true; pattern[5][4] = true
			pattern[6][3] = true; pattern[7][3] = true; pattern[8][3] = true
		case '8':
			for x := 2; x < 6; x++ { pattern[2][x] = true; pattern[5][x] = true; pattern[8][x] = true }
			pattern[3][1] = true; pattern[3][6] = true; pattern[4][1] = true; pattern[4][6] = true
			pattern[6][1] = true; pattern[6][6] = true; pattern[7][1] = true; pattern[7][6] = true
		case '9':
			pattern[2][2] = true; pattern[2][3] = true; pattern[2][4] = true; pattern[2][5] = true
			pattern[3][1] = true; pattern[3][6] = true; pattern[4][1] = true; pattern[4][6] = true
			pattern[5][2] = true; pattern[5][3] = true; pattern[5][4] = true; pattern[5][6] = true
			pattern[6][5] = true; pattern[7][4] = true
			pattern[8][3] = true; pattern[8][4] = true
		}
		
	default:
		// Unknown character - draw a simple box
		for y := 2; y < 9; y++ {
			pattern[y][2] = true; pattern[y][5] = true
		}
		for x := 2; x < 6; x++ {
			pattern[2][x] = true; pattern[8][x] = true
		}
	}
	
	return pattern
}

// getMenuTitle returns the title for the current menu
func (mm *MenuManager) getMenuTitle() string {
	switch mm.currentMenu {
	case MainMenu:
		return "VANIA - Procedural Metroidvania"
	case PauseMenu:
		return "Game Paused"
	case SettingsMenu:
		return "Settings"
	case SaveLoadMenu:
		return "Save / Load"
	case GameOverMenu:
		return "Game Over"
	default:
		return "Menu"
	}
}

// buildMainMenuItems creates main menu items
func (mm *MenuManager) buildMainMenuItems() {
	mm.items = []*MenuItem{
		{
			Text:    "New Game (Random Seed)",
			Enabled: true,
			Action: func() error {
				if mm.onNewGame != nil {
					return mm.onNewGame(0) // 0 = random seed
				}
				return nil
			},
		},
		{
			Text:    "New Game (Seed: 42)",
			Enabled: true,
			Action: func() error {
				if mm.onNewGame != nil {
					return mm.onNewGame(42)
				}
				return nil
			},
		},
		{
			Text:    "Load Game",
			Enabled: mm.saveManager != nil && mm.hasSaveFiles(),
			Action: func() error {
				mm.ShowSaveLoadMenu()
				return nil
			},
		},
		{
			Text:    "Settings",
			Enabled: true,
			Action: func() error {
				mm.ShowSettingsMenu()
				return nil
			},
		},
		{
			Text:    "Quit",
			Enabled: true,
			Action: func() error {
				if mm.onQuitGame != nil {
					return mm.onQuitGame()
				}
				return ebiten.Termination
			},
		},
	}
}

// buildPauseMenuItems creates pause menu items
func (mm *MenuManager) buildPauseMenuItems() {
	mm.items = []*MenuItem{
		{
			Text:    "Resume",
			Enabled: true,
			Action: func() error {
				if mm.onResumeGame != nil {
					return mm.onResumeGame()
				}
				mm.Hide()
				return nil
			},
		},
		{
			Text:    "Save Game",
			Enabled: mm.saveManager != nil,
			Action: func() error {
				// Default to slot 0 for quick save
				if mm.saveManager != nil {
					// Would need game state - for now just hide menu
					mm.Hide()
				}
				return nil
			},
		},
		{
			Text:    "Load Game",
			Enabled: mm.saveManager != nil && mm.hasSaveFiles(),
			Action: func() error {
				mm.ShowSaveLoadMenu()
				return nil
			},
		},
		{
			Text:    "Settings",
			Enabled: true,
			Action: func() error {
				mm.ShowSettingsMenu()
				return nil
			},
		},
		{
			Text:    "Main Menu",
			Enabled: true,
			Action: func() error {
				mm.ShowMainMenu()
				return nil
			},
		},
		{
			Text:    "Quit Game",
			Enabled: true,
			Action: func() error {
				if mm.onQuitGame != nil {
					return mm.onQuitGame()
				}
				return ebiten.Termination
			},
		},
	}
}

// buildSettingsMenuItems creates settings menu items
func (mm *MenuManager) buildSettingsMenuItems() {
	mm.items = []*MenuItem{
		{
			Text:    fmt.Sprintf("Master Volume: %.0f%%", mm.settings.MasterVolume*100),
			Enabled: true,
			Action: func() error {
				mm.settings.MasterVolume += 0.1
				if mm.settings.MasterVolume > 1.0 {
					mm.settings.MasterVolume = 0.0
				}
				mm.buildSettingsMenuItems() // Rebuild to update display
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("SFX Volume: %.0f%%", mm.settings.SFXVolume*100),
			Enabled: true,
			Action: func() error {
				mm.settings.SFXVolume += 0.1
				if mm.settings.SFXVolume > 1.0 {
					mm.settings.SFXVolume = 0.0
				}
				mm.buildSettingsMenuItems() // Rebuild to update display
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Music Volume: %.0f%%", mm.settings.MusicVolume*100),
			Enabled: true,
			Action: func() error {
				mm.settings.MusicVolume += 0.1
				if mm.settings.MusicVolume > 1.0 {
					mm.settings.MusicVolume = 0.0
				}
				mm.buildSettingsMenuItems() // Rebuild to update display
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Fullscreen: %v", mm.settings.FullScreen),
			Enabled: true,
			Action: func() error {
				mm.settings.FullScreen = !mm.settings.FullScreen
				if mm.settings.FullScreen {
					ebiten.SetFullscreen(true)
				} else {
					ebiten.SetFullscreen(false)
				}
				mm.buildSettingsMenuItems() // Rebuild to update display
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Show FPS: %v", mm.settings.ShowFPS),
			Enabled: true,
			Action: func() error {
				mm.settings.ShowFPS = !mm.settings.ShowFPS
				mm.buildSettingsMenuItems() // Rebuild to update display
				return nil
			},
		},
		{
			Text:    "Back",
			Enabled: true,
			Action: func() error {
				return mm.handleBack()
			},
		},
	}
}

// buildGameOverMenuItems creates game over menu items
func (mm *MenuManager) buildGameOverMenuItems() {
	mm.items = []*MenuItem{
		{
			Text:    "Try Again",
			Enabled: true,
			Action: func() error {
				if mm.onNewGame != nil {
					return mm.onNewGame(0) // New random game
				}
				return nil
			},
		},
		{
			Text:    "Load Game",
			Enabled: mm.saveManager != nil && mm.hasSaveFiles(),
			Action: func() error {
				mm.ShowSaveLoadMenu()
				return nil
			},
		},
		{
			Text:    "Main Menu",
			Enabled: true,
			Action: func() error {
				mm.ShowMainMenu()
				return nil
			},
		},
		{
			Text:    "Quit",
			Enabled: true,
			Action: func() error {
				if mm.onQuitGame != nil {
					return mm.onQuitGame()
				}
				return ebiten.Termination
			},
		},
	}
}

// ShowSaveLoadMenu displays save/load menu
func (mm *MenuManager) ShowSaveLoadMenu() {
	mm.currentMenu = SaveLoadMenu
	mm.state = MenuStateActive
	mm.selectedIndex = 0
	mm.buildSaveLoadMenuItems()
}

// buildSaveLoadMenuItems creates save/load menu items
func (mm *MenuManager) buildSaveLoadMenuItems() {
	mm.items = make([]*MenuItem, 0)

	// Add save slots
	for i := 0; i < 5; i++ {
		slotText := fmt.Sprintf("Slot %d", i+1)
		if mm.saveManager != nil {
			if saveData, err := mm.saveManager.LoadGame(i); err == nil {
				// Format play time
				hours := saveData.PlayTime / 3600
				minutes := (saveData.PlayTime % 3600) / 60
				slotText = fmt.Sprintf("Slot %d - %dh %dm (Seed: %d)", i+1, hours, minutes, saveData.Seed)
			} else {
				slotText = fmt.Sprintf("Slot %d - Empty", i+1)
			}
		}

		slot := i // Capture for closure
		mm.items = append(mm.items, &MenuItem{
			Text:    slotText,
			Enabled: true,
			Action: func() error {
				if mm.onLoadGame != nil {
					return mm.onLoadGame(slot)
				}
				return nil
			},
		})
	}

	// Add back option
	mm.items = append(mm.items, &MenuItem{
		Text:    "Back",
		Enabled: true,
		Action: func() error {
			return mm.handleBack()
		},
	})
}

// navigateUp moves selection up
func (mm *MenuManager) navigateUp() {
	if len(mm.items) == 0 {
		return
	}

	mm.selectedIndex--
	if mm.selectedIndex < 0 {
		mm.selectedIndex = len(mm.items) - 1
	}

	// Skip disabled items
	if !mm.items[mm.selectedIndex].Enabled {
		mm.navigateUp()
	}
}

// navigateDown moves selection down
func (mm *MenuManager) navigateDown() {
	if len(mm.items) == 0 {
		return
	}

	mm.selectedIndex++
	if mm.selectedIndex >= len(mm.items) {
		mm.selectedIndex = 0
	}

	// Skip disabled items
	if !mm.items[mm.selectedIndex].Enabled {
		mm.navigateDown()
	}
}

// selectCurrentItem executes the action of the selected item
func (mm *MenuManager) selectCurrentItem() error {
	if mm.selectedIndex >= 0 && mm.selectedIndex < len(mm.items) {
		item := mm.items[mm.selectedIndex]
		if item.Enabled && item.Action != nil {
			return item.Action()
		}
	}
	return nil
}

// handleBack handles the back/escape action
func (mm *MenuManager) handleBack() error {
	switch mm.currentMenu {
	case MainMenu:
		// Quit from main menu
		if mm.onQuitGame != nil {
			return mm.onQuitGame()
		}
		return ebiten.Termination
	case PauseMenu:
		// Resume game from pause menu
		if mm.onResumeGame != nil {
			return mm.onResumeGame()
		}
		mm.Hide()
	case SettingsMenu, SaveLoadMenu:
		// Go back to previous menu
		mm.ShowMainMenu()
	case GameOverMenu:
		// Go to main menu from game over
		mm.ShowMainMenu()
	}
	return nil
}

// hasSaveFiles checks if any save files exist
func (mm *MenuManager) hasSaveFiles() bool {
	if mm.saveManager == nil {
		return false
	}

	for i := 0; i < 5; i++ {
		if _, err := mm.saveManager.LoadGame(i); err == nil {
			return true
		}
	}
	return false
}

// GetSettings returns the current settings
func (mm *MenuManager) GetSettings() *GameSettings {
	return mm.settings
}

// SetSettings updates the settings
func (mm *MenuManager) SetSettings(settings *GameSettings) {
	mm.settings = settings
}
