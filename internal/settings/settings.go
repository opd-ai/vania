// Package settings provides comprehensive game settings management withpackage settings

// persistence, input remapping, graphics options, and user preferences.
package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ControlAction represents different game actions that can be mapped to keys
type ControlAction int

const (
	ActionMoveLeft ControlAction = iota
	ActionMoveRight
	ActionJump
	ActionDash
	ActionAttack
	ActionInteract
	ActionPause
	ActionMenu
	ActionInventory
)

// String returns the human-readable name of a control action
func (a ControlAction) String() string {
	switch a {
	case ActionMoveLeft:
		return "Move Left"
	case ActionMoveRight:
		return "Move Right"
	case ActionJump:
		return "Jump"
	case ActionDash:
		return "Dash"
	case ActionAttack:
		return "Attack"
	case ActionInteract:
		return "Interact"
	case ActionPause:
		return "Pause"
	case ActionMenu:
		return "Menu"
	case ActionInventory:
		return "Inventory"
	default:
		return "Unknown"
	}
}

// GraphicsQuality represents different visual quality levels
type GraphicsQuality int

const (
	QualityLow GraphicsQuality = iota
	QualityMedium
	QualityHigh
	QualityUltra
)

func (q GraphicsQuality) String() string {
	switch q {
	case QualityLow:
		return "Low"
	case QualityMedium:
		return "Medium"
	case QualityHigh:
		return "High"
	case QualityUltra:
		return "Ultra"
	default:
		return "Medium"
	}
}

// AudioSettings holds audio-related configuration
type AudioSettings struct {
	MasterVolume float64 `json:"master_volume"`
	SFXVolume    float64 `json:"sfx_volume"`
	MusicVolume  float64 `json:"music_volume"`
	Muted        bool    `json:"muted"`
}

// GraphicsSettings holds graphics-related configuration
type GraphicsSettings struct {
	Quality         GraphicsQuality `json:"quality"`
	Fullscreen      bool            `json:"fullscreen"`
	VSync           bool            `json:"vsync"`
	WindowWidth     int             `json:"window_width"`
	WindowHeight    int             `json:"window_height"`
	ShowFPS         bool            `json:"show_fps"`
	ParticleEffects bool            `json:"particle_effects"`
	ScreenShake     bool            `json:"screen_shake"`
	UIScale         float64         `json:"ui_scale"`
}

// GameplaySettings holds gameplay-related configuration
type GameplaySettings struct {
	Difficulty       int     `json:"difficulty"` // 0=Easy, 1=Normal, 2=Hard, 3=Expert
	AutoSave         bool    `json:"auto_save"`
	ShowHints        bool    `json:"show_hints"`
	InputBuffering   bool    `json:"input_buffering"`
	CameraSmoothing  float64 `json:"camera_smoothing"`
	MouseSensitivity float64 `json:"mouse_sensitivity"`
}

// ControlSettings holds key mapping configuration
type ControlSettings struct {
	KeyBindings    map[ControlAction]ebiten.Key `json:"key_bindings"`
	GamepadEnabled bool                         `json:"gamepad_enabled"`
}

// Settings holds all game configuration
type Settings struct {
	Audio    AudioSettings    `json:"audio"`
	Graphics GraphicsSettings `json:"graphics"`
	Gameplay GameplaySettings `json:"gameplay"`
	Controls ControlSettings  `json:"controls"`
	Version  string           `json:"version"`
}

// SettingsManager manages loading, saving, and validation of game settings
type SettingsManager struct {
	settings     *Settings
	settingsPath string
	callbacks    map[string]func(*Settings)
}

// NewSettingsManager creates a new settings manager with default settings
func NewSettingsManager() *SettingsManager {
	sm := &SettingsManager{
		callbacks: make(map[string]func(*Settings)),
	}

	// Determine settings file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		sm.settingsPath = "vania_settings.json" // Fallback to current directory
	} else {
		configDir := filepath.Join(homeDir, ".config", "vania")
		os.MkdirAll(configDir, 0755)
		sm.settingsPath = filepath.Join(configDir, "settings.json")
	}

	sm.settings = sm.createDefaultSettings()
	sm.LoadSettings()

	return sm
}

// createDefaultSettings creates sensible default settings
func (sm *SettingsManager) createDefaultSettings() *Settings {
	return &Settings{
		Audio: AudioSettings{
			MasterVolume: 0.8,
			SFXVolume:    0.7,
			MusicVolume:  0.6,
			Muted:        false,
		},
		Graphics: GraphicsSettings{
			Quality:         QualityMedium,
			Fullscreen:      false,
			VSync:           true,
			WindowWidth:     1024,
			WindowHeight:    768,
			ShowFPS:         false,
			ParticleEffects: true,
			ScreenShake:     true,
			UIScale:         1.0,
		},
		Gameplay: GameplaySettings{
			Difficulty:       1, // Normal
			AutoSave:         true,
			ShowHints:        true,
			InputBuffering:   true,
			CameraSmoothing:  0.1,
			MouseSensitivity: 1.0,
		},
		Controls: ControlSettings{
			KeyBindings: map[ControlAction]ebiten.Key{
				ActionMoveLeft:  ebiten.KeyA,
				ActionMoveRight: ebiten.KeyD,
				ActionJump:      ebiten.KeySpace,
				ActionDash:      ebiten.KeyShift,
				ActionAttack:    ebiten.KeyJ,
				ActionInteract:  ebiten.KeyF,
				ActionPause:     ebiten.KeyEscape,
				ActionMenu:      ebiten.KeyTab,
				ActionInventory: ebiten.KeyI,
			},
			GamepadEnabled: true,
		},
		Version: "1.0.0",
	}
}

// LoadSettings loads settings from file or creates defaults
func (sm *SettingsManager) LoadSettings() error {
	data, err := ioutil.ReadFile(sm.settingsPath)
	if err != nil {
		// File doesn't exist, use defaults
		return sm.SaveSettings()
	}

	var loadedSettings Settings
	if err := json.Unmarshal(data, &loadedSettings); err != nil {
		return fmt.Errorf("failed to parse settings: %w", err)
	}

	// Validate and merge with defaults to handle version changes
	sm.settings = sm.validateAndMergeSettings(&loadedSettings)

	// Notify callbacks
	sm.notifyCallbacks()

	return nil
}

// SaveSettings saves current settings to file
func (sm *SettingsManager) SaveSettings() error {
	data, err := json.MarshalIndent(sm.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize settings: %w", err)
	}

	return ioutil.WriteFile(sm.settingsPath, data, 0644)
}

// validateAndMergeSettings ensures loaded settings have all required fields
func (sm *SettingsManager) validateAndMergeSettings(loaded *Settings) *Settings {
	defaults := sm.createDefaultSettings()

	// Merge audio settings
	if loaded.Audio.MasterVolume <= 0 {
		loaded.Audio.MasterVolume = defaults.Audio.MasterVolume
	}
	if loaded.Audio.SFXVolume <= 0 {
		loaded.Audio.SFXVolume = defaults.Audio.SFXVolume
	}
	if loaded.Audio.MusicVolume <= 0 {
		loaded.Audio.MusicVolume = defaults.Audio.MusicVolume
	}

	// Merge graphics settings
	if loaded.Graphics.WindowWidth <= 0 {
		loaded.Graphics.WindowWidth = defaults.Graphics.WindowWidth
	}
	if loaded.Graphics.WindowHeight <= 0 {
		loaded.Graphics.WindowHeight = defaults.Graphics.WindowHeight
	}
	if loaded.Graphics.UIScale <= 0 {
		loaded.Graphics.UIScale = defaults.Graphics.UIScale
	}

	// Merge gameplay settings
	if loaded.Gameplay.CameraSmoothing <= 0 {
		loaded.Gameplay.CameraSmoothing = defaults.Gameplay.CameraSmoothing
	}
	if loaded.Gameplay.MouseSensitivity <= 0 {
		loaded.Gameplay.MouseSensitivity = defaults.Gameplay.MouseSensitivity
	}

	// Ensure all key bindings exist
	if loaded.Controls.KeyBindings == nil {
		loaded.Controls.KeyBindings = make(map[ControlAction]ebiten.Key)
	}
	for action, defaultKey := range defaults.Controls.KeyBindings {
		if _, exists := loaded.Controls.KeyBindings[action]; !exists {
			loaded.Controls.KeyBindings[action] = defaultKey
		}
	}

	return loaded
}

// GetSettings returns a copy of the current settings
func (sm *SettingsManager) GetSettings() *Settings {
	settings := *sm.settings
	// Deep copy key bindings
	settings.Controls.KeyBindings = make(map[ControlAction]ebiten.Key)
	for k, v := range sm.settings.Controls.KeyBindings {
		settings.Controls.KeyBindings[k] = v
	}
	return &settings
}

// UpdateAudioSettings updates audio settings and saves
func (sm *SettingsManager) UpdateAudioSettings(audio AudioSettings) error {
	sm.settings.Audio = audio
	sm.notifyCallbacks()
	return sm.SaveSettings()
}

// UpdateGraphicsSettings updates graphics settings and saves
func (sm *SettingsManager) UpdateGraphicsSettings(graphics GraphicsSettings) error {
	sm.settings.Graphics = graphics
	sm.notifyCallbacks()
	return sm.SaveSettings()
}

// UpdateGameplaySettings updates gameplay settings and saves
func (sm *SettingsManager) UpdateGameplaySettings(gameplay GameplaySettings) error {
	sm.settings.Gameplay = gameplay
	sm.notifyCallbacks()
	return sm.SaveSettings()
}

// UpdateControlSettings updates control settings and saves
func (sm *SettingsManager) UpdateControlSettings(controls ControlSettings) error {
	sm.settings.Controls = controls
	sm.notifyCallbacks()
	return sm.SaveSettings()
}

// SetKeyBinding updates a single key binding
func (sm *SettingsManager) SetKeyBinding(action ControlAction, key ebiten.Key) error {
	// Check for conflicts
	for otherAction, otherKey := range sm.settings.Controls.KeyBindings {
		if otherAction != action && otherKey == key {
			return fmt.Errorf("key %v is already bound to %v", key, otherAction)
		}
	}

	sm.settings.Controls.KeyBindings[action] = key
	sm.notifyCallbacks()
	return sm.SaveSettings()
}

// IsActionPressed checks if an action is currently pressed
func (sm *SettingsManager) IsActionPressed(action ControlAction) bool {
	key, exists := sm.settings.Controls.KeyBindings[action]
	if !exists {
		return false
	}
	return ebiten.IsKeyPressed(key)
}

// IsActionJustPressed checks if an action was just pressed this frame
func (sm *SettingsManager) IsActionJustPressed(action ControlAction) bool {
	key, exists := sm.settings.Controls.KeyBindings[action]
	if !exists {
		return false
	}
	return inpututil.IsKeyJustPressed(key)
}

// RegisterCallback registers a callback to be called when settings change
func (sm *SettingsManager) RegisterCallback(name string, callback func(*Settings)) {
	sm.callbacks[name] = callback
}

// UnregisterCallback removes a callback
func (sm *SettingsManager) UnregisterCallback(name string) {
	delete(sm.callbacks, name)
}

// notifyCallbacks calls all registered callbacks with current settings
func (sm *SettingsManager) notifyCallbacks() {
	for _, callback := range sm.callbacks {
		callback(sm.settings)
	}
}

// ApplyGraphicsSettings applies graphics settings to Ebiten
func (sm *SettingsManager) ApplyGraphicsSettings() {
	graphics := &sm.settings.Graphics

	// Apply window size
	ebiten.SetWindowSize(graphics.WindowWidth, graphics.WindowHeight)

	// Apply fullscreen
	ebiten.SetFullscreen(graphics.Fullscreen)

	// Apply VSync
	ebiten.SetVsyncEnabled(graphics.VSync)

	// Window title and other properties would be set by the main application
}

// ResetToDefaults resets all settings to default values
func (sm *SettingsManager) ResetToDefaults() error {
	sm.settings = sm.createDefaultSettings()
	sm.notifyCallbacks()
	return sm.SaveSettings()
}

// GetDifficultyName returns human-readable difficulty name
func (sm *SettingsManager) GetDifficultyName(difficulty int) string {
	switch difficulty {
	case 0:
		return "Easy"
	case 1:
		return "Normal"
	case 2:
		return "Hard"
	case 3:
		return "Expert"
	default:
		return "Normal"
	}
}

// GetQualitySettings returns performance settings for the given quality level
func (sm *SettingsManager) GetQualitySettings(quality GraphicsQuality) (particleCount int, shadowQuality int, textureQuality int) {
	switch quality {
	case QualityLow:
		return 50, 0, 0 // Minimal particles, no shadows, low textures
	case QualityMedium:
		return 150, 1, 1 // Normal particles, basic shadows, medium textures
	case QualityHigh:
		return 300, 2, 2 // Many particles, good shadows, high textures
	case QualityUltra:
		return 500, 3, 3 // Maximum particles, best shadows, ultra textures
	default:
		return 150, 1, 1
	}
}

// ExportSettings exports settings to a JSON string for sharing
func (sm *SettingsManager) ExportSettings() (string, error) {
	data, err := json.MarshalIndent(sm.settings, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to export settings: %w", err)
	}
	return string(data), nil
}

// ImportSettings imports settings from a JSON string
func (sm *SettingsManager) ImportSettings(jsonData string) error {
	var imported Settings
	if err := json.Unmarshal([]byte(jsonData), &imported); err != nil {
		return fmt.Errorf("failed to parse imported settings: %w", err)
	}

	sm.settings = sm.validateAndMergeSettings(&imported)
	sm.notifyCallbacks()
	return sm.SaveSettings()
}
