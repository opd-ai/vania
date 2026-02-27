// Package menu provides controls configuration UI for key rebinding
package menu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/opd-ai/vania/internal/settings"
)

// ControlsMenuState represents the state of the controls configuration menu
type ControlsMenuState int

const (
	ControlsMenuBrowsing ControlsMenuState = iota
	ControlsMenuWaitingForKey
)

// ControlsMenu handles key rebinding UI
type ControlsMenu struct {
	menuManager      *MenuManager
	state            ControlsMenuState
	selectedAction   settings.ControlAction
	actions          []settings.ControlAction
	waitingFrames    int // Frames to wait before accepting key input
	conflictDetected bool
	conflictMessage  string
	successMessage   string
	messageTimer     int
}

// NewControlsMenu creates a new controls configuration menu
func NewControlsMenu(mm *MenuManager) *ControlsMenu {
	return &ControlsMenu{
		menuManager: mm,
		state:       ControlsMenuBrowsing,
		actions: []settings.ControlAction{
			settings.ActionMoveLeft,
			settings.ActionMoveRight,
			settings.ActionJump,
			settings.ActionAttack,
			settings.ActionDash,
			settings.ActionInteract,
			settings.ActionPause,
		},
		waitingFrames:    0,
		conflictDetected: false,
		messageTimer:     0,
	}
}

// Show displays the controls menu
func (cm *ControlsMenu) Show() {
	cm.state = ControlsMenuBrowsing
	cm.conflictDetected = false
	cm.successMessage = ""
	cm.messageTimer = 0
	cm.menuManager.buildControlsMenuItems()
}

// Update handles controls menu input
func (cm *ControlsMenu) Update() error {
	if cm.messageTimer > 0 {
		cm.messageTimer--
		if cm.messageTimer == 0 {
			cm.successMessage = ""
		}
	}

	if cm.state == ControlsMenuWaitingForKey {
		return cm.waitForKeyPress()
	}

	// Normal browsing mode
	inputState := cm.menuManager.inputHandler.Update()

	// Handle back/escape
	if inputState.PausePress {
		cm.menuManager.ShowSettingsMenu()
		return nil
	}

	return nil
}

// waitForKeyPress waits for user to press a key for rebinding
func (cm *ControlsMenu) waitForKeyPress() error {
	// Wait a few frames to avoid capturing the Enter key that selected the action
	if cm.waitingFrames < 10 {
		cm.waitingFrames++
		return nil
	}

	// Check for ESC to cancel
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		cm.state = ControlsMenuBrowsing
		cm.waitingFrames = 0
		cm.menuManager.buildControlsMenuItems()
		return nil
	}

	// Check all keys for input
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if inpututil.IsKeyJustPressed(key) {
			// Don't allow rebinding Escape (reserved for back/menu)
			if key == ebiten.KeyEscape {
				continue
			}

			// Try to set the key binding
			err := cm.menuManager.settingsManager.SetKeyBinding(cm.selectedAction, key)
			if err != nil {
				// Conflict detected
				cm.conflictDetected = true
				cm.conflictMessage = err.Error()
				cm.messageTimer = 180 // Show message for 3 seconds
			} else {
				// Success
				cm.successMessage = fmt.Sprintf("Key bound successfully!")
				cm.messageTimer = 120 // Show message for 2 seconds
				cm.conflictDetected = false
			}

			// Return to browsing mode
			cm.state = ControlsMenuBrowsing
			cm.waitingFrames = 0
			cm.menuManager.buildControlsMenuItems()
			return nil
		}
	}

	return nil
}

// StartRebind initiates key rebinding for the given action
func (cm *ControlsMenu) StartRebind(action settings.ControlAction) {
	cm.selectedAction = action
	cm.state = ControlsMenuWaitingForKey
	cm.waitingFrames = 0
	cm.conflictDetected = false
	cm.successMessage = ""
	cm.messageTimer = 0
}

// BuildControlsMenuItems creates menu items for controls configuration
func (mm *MenuManager) buildControlsMenuItems() {
	if mm.controlsMenu == nil {
		mm.controlsMenu = NewControlsMenu(mm)
	}

	controls := mm.settingsManager.GetSettings().Controls

	mm.items = []*MenuItem{
		{
			Text:    fmt.Sprintf("Move Left: %s", getKeyName(controls.KeyBindings[settings.ActionMoveLeft])),
			Enabled: true,
			Action: func() error {
				mm.controlsMenu.StartRebind(settings.ActionMoveLeft)
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Move Right: %s", getKeyName(controls.KeyBindings[settings.ActionMoveRight])),
			Enabled: true,
			Action: func() error {
				mm.controlsMenu.StartRebind(settings.ActionMoveRight)
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Jump: %s", getKeyName(controls.KeyBindings[settings.ActionJump])),
			Enabled: true,
			Action: func() error {
				mm.controlsMenu.StartRebind(settings.ActionJump)
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Attack: %s", getKeyName(controls.KeyBindings[settings.ActionAttack])),
			Enabled: true,
			Action: func() error {
				mm.controlsMenu.StartRebind(settings.ActionAttack)
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Dash: %s", getKeyName(controls.KeyBindings[settings.ActionDash])),
			Enabled: true,
			Action: func() error {
				mm.controlsMenu.StartRebind(settings.ActionDash)
				return nil
			},
		},
		{
			Text:    fmt.Sprintf("Use Ability: %s", getKeyName(controls.KeyBindings[settings.ActionInteract])),
			Enabled: true,
			Action: func() error {
				mm.controlsMenu.StartRebind(settings.ActionInteract)
				return nil
			},
		},
		{
			Text:    "Reset to Defaults",
			Enabled: true,
			Action: func() error {
				// Reset only control settings
				defaults := mm.settingsManager.GetSettings()
				defaultControls := settings.ControlSettings{
					KeyBindings: map[settings.ControlAction]ebiten.Key{
						settings.ActionMoveLeft:  ebiten.KeyA,
						settings.ActionMoveRight: ebiten.KeyD,
						settings.ActionJump:      ebiten.KeySpace,
						settings.ActionDash:      ebiten.KeyShift,
						settings.ActionAttack:    ebiten.KeyJ,
						settings.ActionInteract:  ebiten.KeyF,
						settings.ActionPause:     ebiten.KeyEscape,
					},
					GamepadEnabled: true,
				}
				defaults.Controls = defaultControls
				mm.settingsManager.UpdateControlSettings(defaultControls)
				mm.buildControlsMenuItems()
				return nil
			},
		},
		{
			Text:    "Back",
			Enabled: true,
			Action: func() error {
				mm.ShowSettingsMenu()
				return nil
			},
		},
	}
}

// getKeyName returns a human-readable name for a key
func getKeyName(key ebiten.Key) string {
	switch key {
	case ebiten.KeyA:
		return "A"
	case ebiten.KeyB:
		return "B"
	case ebiten.KeyC:
		return "C"
	case ebiten.KeyD:
		return "D"
	case ebiten.KeyE:
		return "E"
	case ebiten.KeyF:
		return "F"
	case ebiten.KeyG:
		return "G"
	case ebiten.KeyH:
		return "H"
	case ebiten.KeyI:
		return "I"
	case ebiten.KeyJ:
		return "J"
	case ebiten.KeyK:
		return "K"
	case ebiten.KeyL:
		return "L"
	case ebiten.KeyM:
		return "M"
	case ebiten.KeyN:
		return "N"
	case ebiten.KeyO:
		return "O"
	case ebiten.KeyP:
		return "P"
	case ebiten.KeyQ:
		return "Q"
	case ebiten.KeyR:
		return "R"
	case ebiten.KeyS:
		return "S"
	case ebiten.KeyT:
		return "T"
	case ebiten.KeyU:
		return "U"
	case ebiten.KeyV:
		return "V"
	case ebiten.KeyW:
		return "W"
	case ebiten.KeyX:
		return "X"
	case ebiten.KeyY:
		return "Y"
	case ebiten.KeyZ:
		return "Z"
	case ebiten.KeySpace:
		return "SPACE"
	case ebiten.KeyShift:
		return "SHIFT"
	case ebiten.KeyControl:
		return "CTRL"
	case ebiten.KeyAlt:
		return "ALT"
	case ebiten.KeyEnter:
		return "ENTER"
	case ebiten.KeyEscape:
		return "ESC"
	case ebiten.KeyArrowLeft:
		return "LEFT"
	case ebiten.KeyArrowRight:
		return "RIGHT"
	case ebiten.KeyArrowUp:
		return "UP"
	case ebiten.KeyArrowDown:
		return "DOWN"
	default:
		return fmt.Sprintf("Key%d", key)
	}
}
