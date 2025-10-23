package menu

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewMenuManager(t *testing.T) {
	mm := NewMenuManager()
	
	if mm == nil {
		t.Fatal("NewMenuManager returned nil")
	}
	
	if mm.currentMenu != MainMenu {
		t.Errorf("Expected current menu to be MainMenu, got %v", mm.currentMenu)
	}
	
	if mm.state != MenuStateActive {
		t.Errorf("Expected state to be MenuStateActive, got %v", mm.state)
	}
	
	if mm.settings == nil {
		t.Error("Settings should not be nil")
	}
	
	// Check default settings
	if mm.settings.MasterVolume != 0.7 {
		t.Errorf("Expected MasterVolume to be 0.7, got %f", mm.settings.MasterVolume)
	}
}

func TestMenuNavigation(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowMainMenu()
	
	if len(mm.items) == 0 {
		t.Fatal("Main menu should have items")
	}
	
	initialIndex := mm.selectedIndex
	
	// Test navigation down
	mm.navigateDown()
	if mm.selectedIndex <= initialIndex {
		t.Error("Navigation down should increase selected index")
	}
	
	// Test navigation up
	mm.navigateUp()
	if mm.selectedIndex != initialIndex {
		t.Error("Navigation up should return to initial index")
	}
}

func TestMenuStates(t *testing.T) {
	mm := NewMenuManager()
	
	// Test showing different menus
	mm.ShowMainMenu()
	if mm.currentMenu != MainMenu || !mm.IsActive() {
		t.Error("ShowMainMenu should set MainMenu and active state")
	}
	
	mm.ShowPauseMenu()
	if mm.currentMenu != PauseMenu || !mm.IsActive() {
		t.Error("ShowPauseMenu should set PauseMenu and active state")
	}
	
	mm.ShowSettingsMenu()
	if mm.currentMenu != SettingsMenu || !mm.IsActive() {
		t.Error("ShowSettingsMenu should set SettingsMenu and active state")
	}
	
	mm.Hide()
	if mm.IsActive() {
		t.Error("Hide should deactivate menu")
	}
}

func TestSettingsToggle(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowSettingsMenu()
	
	// Test fullscreen toggle
	initialFullscreen := mm.settings.FullScreen
	
	// Find fullscreen item and trigger it
	for i, item := range mm.items {
		if item.Text == "Fullscreen: false" || item.Text == "Fullscreen: true" {
			mm.selectedIndex = i
			err := mm.selectCurrentItem()
			if err != nil {
				t.Errorf("Error selecting fullscreen item: %v", err)
			}
			break
		}
	}
	
	if mm.settings.FullScreen == initialFullscreen {
		t.Error("Fullscreen setting should have toggled")
	}
}

func TestVolumeAdjustment(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowSettingsMenu()
	
	initialVolume := mm.settings.MasterVolume
	
	// Find master volume item and trigger it
	for i, item := range mm.items {
		if len(item.Text) > 13 && item.Text[:13] == "Master Volume" {
			mm.selectedIndex = i
			err := mm.selectCurrentItem()
			if err != nil {
				t.Errorf("Error selecting volume item: %v", err)
			}
			break
		}
	}
	
	if mm.settings.MasterVolume == initialVolume {
		t.Error("Master volume should have changed")
	}
}

func TestCallbacksSetup(t *testing.T) {
	mm := NewMenuManager()
	
	called := false
	mm.SetCallbacks(
		func(seed int64) error { called = true; return nil }, // onNewGame
		nil, // onLoadGame
		nil, // onSettings
		nil, // onQuitGame
		nil, // onResumeGame
	)
	
	if mm.onNewGame == nil {
		t.Error("onNewGame callback should be set")
	}
	
	// Test callback is called
	mm.ShowMainMenu()
	
	// Find new game item and trigger it
	for i, item := range mm.items {
		if item.Text == "New Game (Random Seed)" {
			mm.selectedIndex = i
			err := mm.selectCurrentItem()
			if err != nil {
				t.Errorf("Error selecting new game: %v", err)
			}
			break
		}
	}
	
	if !called {
		t.Error("New game callback should have been called")
	}
}

func TestMenuItemEnabling(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowMainMenu()
	
	// Check that items are properly enabled/disabled
	foundEnabledItem := false
	for _, item := range mm.items {
		if item.Enabled {
			foundEnabledItem = true
		}
		// All main menu items should be enabled except possibly Load Game
		if item.Text != "Load Game" && !item.Enabled {
			t.Errorf("Menu item '%s' should be enabled", item.Text)
		}
	}
	
	if !foundEnabledItem {
		t.Error("At least one menu item should be enabled")
	}
}

func TestGameOverMenu(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowGameOverMenu()
	
	if mm.currentMenu != GameOverMenu {
		t.Error("Should be showing game over menu")
	}
	
	if len(mm.items) == 0 {
		t.Error("Game over menu should have items")
	}
	
	// Check that essential items exist
	hasRetry := false
	hasMainMenu := false
	
	for _, item := range mm.items {
		if item.Text == "Try Again" {
			hasRetry = true
		}
		if item.Text == "Main Menu" {
			hasMainMenu = true
		}
	}
	
	if !hasRetry {
		t.Error("Game over menu should have 'Try Again' option")
	}
	
	if !hasMainMenu {
		t.Error("Game over menu should have 'Main Menu' option")
	}
}

func TestMenuDrawing(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowMainMenu()
	
	// Create a test screen
	screen := ebiten.NewImage(960, 640)
	
	// Drawing should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Menu drawing panicked: %v", r)
		}
	}()
	
	mm.Draw(screen)
	
	// When inactive, should not draw
	mm.Hide()
	mm.Draw(screen) // Should not panic
}

func TestKeyBindings(t *testing.T) {
	mm := NewMenuManager()
	
	if mm.settings.KeyBindings == nil {
		t.Error("Key bindings should be initialized")
	}
	
	// Check essential key bindings exist
	requiredBindings := []string{"move_left", "move_right", "jump", "attack", "dash", "pause"}
	
	for _, binding := range requiredBindings {
		if _, exists := mm.settings.KeyBindings[binding]; !exists {
			t.Errorf("Key binding '%s' should exist", binding)
		}
	}
	
	// Check that bindings have actual keys
	if len(mm.settings.KeyBindings["jump"]) == 0 {
		t.Error("Jump binding should have at least one key")
	}
}