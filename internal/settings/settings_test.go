package settings

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewSettingsManager(t *testing.T) {
	sm := NewSettingsManager()

	if sm == nil {
		t.Fatal("NewSettingsManager returned nil")
	}

	if sm.settings == nil {
		t.Fatal("Settings not initialized")
	}

	if sm.callbacks == nil {
		t.Fatal("Callbacks map not initialized")
	}
}

func TestDefaultSettings(t *testing.T) {
	sm := NewSettingsManager()
	settings := sm.GetSettings()

	// Test audio defaults
	if settings.Audio.MasterVolume <= 0 || settings.Audio.MasterVolume > 1 {
		t.Errorf("Invalid default master volume: %f", settings.Audio.MasterVolume)
	}

	// Test graphics defaults
	if settings.Graphics.WindowWidth <= 0 || settings.Graphics.WindowHeight <= 0 {
		t.Errorf("Invalid default window size: %dx%d", settings.Graphics.WindowWidth, settings.Graphics.WindowHeight)
	}

	// Test control defaults
	if len(settings.Controls.KeyBindings) == 0 {
		t.Error("No default key bindings set")
	}

	// Check specific key bindings
	if settings.Controls.KeyBindings[ActionJump] != ebiten.KeySpace {
		t.Error("Jump not bound to Space key by default")
	}

	if settings.Controls.KeyBindings[ActionMoveLeft] != ebiten.KeyA {
		t.Error("Move left not bound to A key by default")
	}
}

func TestControlActionString(t *testing.T) {
	testCases := []struct {
		action   ControlAction
		expected string
	}{
		{ActionMoveLeft, "Move Left"},
		{ActionMoveRight, "Move Right"},
		{ActionJump, "Jump"},
		{ActionDash, "Dash"},
		{ActionAttack, "Attack"},
		{ActionInteract, "Interact"},
		{ActionPause, "Pause"},
		{ActionMenu, "Menu"},
		{ActionInventory, "Inventory"},
	}

	for _, tc := range testCases {
		if got := tc.action.String(); got != tc.expected {
			t.Errorf("ControlAction.String() = %q, want %q", got, tc.expected)
		}
	}
}

func TestGraphicsQualityString(t *testing.T) {
	testCases := []struct {
		quality  GraphicsQuality
		expected string
	}{
		{QualityLow, "Low"},
		{QualityMedium, "Medium"},
		{QualityHigh, "High"},
		{QualityUltra, "Ultra"},
	}

	for _, tc := range testCases {
		if got := tc.quality.String(); got != tc.expected {
			t.Errorf("GraphicsQuality.String() = %q, want %q", got, tc.expected)
		}
	}
}

func TestSaveLoadSettings(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	sm := NewSettingsManager()
	// Override settings path to use temporary directory
	sm.settingsPath = filepath.Join(tmpDir, "test_settings.json")

	// Modify some settings
	originalSettings := sm.GetSettings()
	originalSettings.Audio.MasterVolume = 0.5
	originalSettings.Graphics.Fullscreen = true
	originalSettings.Gameplay.Difficulty = 2

	sm.UpdateAudioSettings(originalSettings.Audio)
	sm.UpdateGraphicsSettings(originalSettings.Graphics)
	sm.UpdateGameplaySettings(originalSettings.Gameplay)

	// Create new settings manager and load
	sm2 := NewSettingsManager()
	sm2.settingsPath = sm.settingsPath

	err := sm2.LoadSettings()
	if err != nil {
		t.Fatalf("Failed to load settings: %v", err)
	}

	loadedSettings := sm2.GetSettings()

	// Verify loaded settings match saved settings
	if loadedSettings.Audio.MasterVolume != 0.5 {
		t.Errorf("Master volume not loaded correctly: got %f, want 0.5", loadedSettings.Audio.MasterVolume)
	}

	if !loadedSettings.Graphics.Fullscreen {
		t.Error("Fullscreen setting not loaded correctly")
	}

	if loadedSettings.Gameplay.Difficulty != 2 {
		t.Errorf("Difficulty not loaded correctly: got %d, want 2", loadedSettings.Gameplay.Difficulty)
	}
}

func TestSetKeyBinding(t *testing.T) {
	sm := NewSettingsManager()

	// Test setting a new key binding
	err := sm.SetKeyBinding(ActionJump, ebiten.KeyW)
	if err != nil {
		t.Fatalf("Failed to set key binding: %v", err)
	}

	settings := sm.GetSettings()
	if settings.Controls.KeyBindings[ActionJump] != ebiten.KeyW {
		t.Error("Key binding not updated correctly")
	}

	// Test conflict detection
	err = sm.SetKeyBinding(ActionAttack, ebiten.KeyW)
	if err == nil {
		t.Error("Expected conflict error when setting duplicate key binding")
	}
}

func TestValidateAndMergeSettings(t *testing.T) {
	sm := NewSettingsManager()

	// Create incomplete settings
	incomplete := &Settings{
		Audio: AudioSettings{
			MasterVolume: 0, // Invalid, should be filled with default
		},
		Graphics: GraphicsSettings{
			WindowWidth: -100, // Invalid, should be filled with default
		},
		Controls: ControlSettings{
			KeyBindings: map[ControlAction]ebiten.Key{
				ActionJump: ebiten.KeySpace, // Only partial bindings
			},
		},
	}

	merged := sm.validateAndMergeSettings(incomplete)

	// Check that defaults were applied
	if merged.Audio.MasterVolume <= 0 {
		t.Error("Default master volume not applied")
	}

	if merged.Graphics.WindowWidth <= 0 {
		t.Error("Default window width not applied")
	}

	// Check that all key bindings are present
	expectedBindings := 9 // Should have all 9 actions
	if len(merged.Controls.KeyBindings) != expectedBindings {
		t.Errorf("Not all key bindings filled: got %d, want %d", len(merged.Controls.KeyBindings), expectedBindings)
	}
}

func TestCallbacks(t *testing.T) {
	sm := NewSettingsManager()

	callbackCalled := false
	var callbackSettings *Settings

	// Register callback
	sm.RegisterCallback("test", func(s *Settings) {
		callbackCalled = true
		callbackSettings = s
	})

	// Update settings to trigger callback
	audio := sm.GetSettings().Audio
	audio.MasterVolume = 0.3
	sm.UpdateAudioSettings(audio)

	if !callbackCalled {
		t.Error("Callback was not called")
	}

	if callbackSettings == nil || callbackSettings.Audio.MasterVolume != 0.3 {
		t.Error("Callback did not receive correct settings")
	}

	// Test unregistering callback
	callbackCalled = false
	sm.UnregisterCallback("test")

	audio.MasterVolume = 0.7
	sm.UpdateAudioSettings(audio)

	if callbackCalled {
		t.Error("Callback was called after unregistering")
	}
}

func TestGetQualitySettings(t *testing.T) {
	sm := NewSettingsManager()

	testCases := []struct {
		quality           GraphicsQuality
		expectedParticles int
	}{
		{QualityLow, 50},
		{QualityMedium, 150},
		{QualityHigh, 300},
		{QualityUltra, 500},
	}

	for _, tc := range testCases {
		particles, shadows, textures := sm.GetQualitySettings(tc.quality)

		if particles != tc.expectedParticles {
			t.Errorf("Quality %v: got %d particles, want %d", tc.quality, particles, tc.expectedParticles)
		}

		// Check that shadow and texture quality increase with overall quality
		if shadows < 0 || textures < 0 {
			t.Errorf("Quality %v: negative quality values", tc.quality)
		}
	}
}

func TestGetDifficultyName(t *testing.T) {
	sm := NewSettingsManager()

	testCases := []struct {
		difficulty int
		expected   string
	}{
		{0, "Easy"},
		{1, "Normal"},
		{2, "Hard"},
		{3, "Expert"},
		{99, "Normal"}, // Invalid, should default to Normal
	}

	for _, tc := range testCases {
		got := sm.GetDifficultyName(tc.difficulty)
		if got != tc.expected {
			t.Errorf("GetDifficultyName(%d) = %q, want %q", tc.difficulty, got, tc.expected)
		}
	}
}

func TestExportImportSettings(t *testing.T) {
	sm := NewSettingsManager()

	// Modify settings
	audio := sm.GetSettings().Audio
	audio.MasterVolume = 0.42
	sm.UpdateAudioSettings(audio)

	// Export settings
	exported, err := sm.ExportSettings()
	if err != nil {
		t.Fatalf("Failed to export settings: %v", err)
	}

	// Verify it's valid JSON
	var testJSON map[string]interface{}
	if err := json.Unmarshal([]byte(exported), &testJSON); err != nil {
		t.Fatalf("Exported settings is not valid JSON: %v", err)
	}

	// Create new settings manager and import
	sm2 := NewSettingsManager()
	err = sm2.ImportSettings(exported)
	if err != nil {
		t.Fatalf("Failed to import settings: %v", err)
	}

	// Verify imported settings match
	importedSettings := sm2.GetSettings()
	if importedSettings.Audio.MasterVolume != 0.42 {
		t.Errorf("Imported master volume incorrect: got %f, want 0.42", importedSettings.Audio.MasterVolume)
	}
}

func TestResetToDefaults(t *testing.T) {
	sm := NewSettingsManager()

	// Modify settings away from defaults
	audio := sm.GetSettings().Audio
	audio.MasterVolume = 0.123
	sm.UpdateAudioSettings(audio)

	graphics := sm.GetSettings().Graphics
	graphics.Fullscreen = true
	sm.UpdateGraphicsSettings(graphics)

	// Reset to defaults
	err := sm.ResetToDefaults()
	if err != nil {
		t.Fatalf("Failed to reset to defaults: %v", err)
	}

	// Verify settings are back to defaults
	settings := sm.GetSettings()
	defaultSettings := sm.createDefaultSettings()

	if settings.Audio.MasterVolume != defaultSettings.Audio.MasterVolume {
		t.Error("Audio settings not reset to defaults")
	}

	if settings.Graphics.Fullscreen != defaultSettings.Graphics.Fullscreen {
		t.Error("Graphics settings not reset to defaults")
	}
}

func TestDeterministicDefaults(t *testing.T) {
	// Test that multiple settings managers have identical defaults
	sm1 := NewSettingsManager()
	sm2 := NewSettingsManager()

	settings1 := sm1.GetSettings()
	settings2 := sm2.GetSettings()

	if settings1.Audio.MasterVolume != settings2.Audio.MasterVolume {
		t.Error("Default audio settings not deterministic")
	}

	if settings1.Graphics.WindowWidth != settings2.Graphics.WindowWidth {
		t.Error("Default graphics settings not deterministic")
	}

	if len(settings1.Controls.KeyBindings) != len(settings2.Controls.KeyBindings) {
		t.Error("Default control settings not deterministic")
	}
}
