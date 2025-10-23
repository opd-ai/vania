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
	currentMenu     MenuType
	state          MenuState
	items          []*MenuItem
	selectedIndex  int
	inputHandler   *input.InputHandler
	saveManager    *save.SaveManager
	
	// Callbacks
	onNewGame    func(seed int64) error
	onLoadGame   func(slot int) error
	onSettings   func() error
	onQuitGame   func() error
	onResumeGame func() error
	
	// Settings
	settings *GameSettings
	
	// Visual properties
	backgroundColor color.Color
	textColor       color.Color
	selectedColor   color.Color
	disabledColor   color.Color
}

// GameSettings holds user configurable settings
type GameSettings struct {
	MasterVolume   float64
	SFXVolume      float64
	MusicVolume    float64
	FullScreen     bool
	VSync          bool
	ShowFPS        bool
	KeyBindings    map[string][]ebiten.Key
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

	return &MenuManager{
		currentMenu:     MainMenu,
		state:          MenuStateActive,
		items:          make([]*MenuItem, 0),
		selectedIndex:  0,
		inputHandler:   input.NewInputHandler(),
		saveManager:    saveManager,
		settings:       settings,
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

	// Draw menu title
	title := mm.getMenuTitle()
	titleX := 480 - len(title)*4 // Approximate centering
	titleY := 100
	ebitenutil.DebugPrintAt(screen, title, titleX, titleY)

	// Draw menu items
	startY := 200
	for i, item := range mm.items {
		y := startY + i*40

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

	// Draw instructions
	instructions := "Use W/S or Arrow Keys to navigate, Enter to select, Esc to back"
	instructY := 500
	ebitenutil.DebugPrintAt(screen, instructions, 480-len(instructions)*3, instructY)
}

// drawColoredText draws text with specified color (simplified implementation)
func (mm *MenuManager) drawColoredText(screen *ebiten.Image, text string, x, y int, col color.Color) {
	// For now, use basic debug print (Ebiten doesn't have easy colored text)
	// In a full implementation, you'd use a font rendering system
	ebitenutil.DebugPrintAt(screen, text, x, y)
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