package menu

import (
	"strings"
	"testing"

	"github.com/opd-ai/vania/internal/save"
)

func TestSaveLoadMenuDisplay(t *testing.T) {
	mm := NewMenuManager()

	// Create save manager with test data
	saveManager, err := save.NewSaveManager("")
	if err != nil {
		t.Fatalf("Failed to create save manager: %v", err)
	}
	mm.saveManager = saveManager

	// Create a test save in slot 1
	testSave := &save.SaveData{
		Seed:     12345,
		PlayTime: 3665, // 1 hour, 1 minute, 5 seconds
		SlotID:   1,
		PlayerX:  100,
		PlayerY:  200,
	}

	err = mm.saveManager.SaveGame(testSave, 1)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}
	defer mm.saveManager.DeleteSave(1) // Clean up

	// Show save/load menu
	mm.ShowSaveLoadMenu()

	// Verify menu type
	if mm.currentMenu != SaveLoadMenu {
		t.Errorf("Expected SaveLoadMenu, got %v", mm.currentMenu)
	}

	// Verify items were created (5 slots + back button)
	expectedItems := 6
	if len(mm.items) != expectedItems {
		t.Errorf("Expected %d items, got %d", expectedItems, len(mm.items))
	}

	// Verify slot 1 shows metadata
	slot1Text := mm.items[1].Text
	if !strings.Contains(slot1Text, "1h 1m") {
		t.Errorf("Slot 1 should show play time '1h 1m', got: %s", slot1Text)
	}
	if !strings.Contains(slot1Text, "12345") {
		t.Errorf("Slot 1 should show seed '12345', got: %s", slot1Text)
	}

	// Verify empty slots show "Empty"
	emptySlots := []int{0, 2, 3, 4}
	for _, idx := range emptySlots {
		slotText := mm.items[idx].Text
		if !strings.Contains(slotText, "Empty") {
			t.Errorf("Slot %d should show 'Empty', got: %s", idx, slotText)
		}
	}

	// Verify all slots are enabled
	for i := 0; i < 5; i++ {
		if !mm.items[i].Enabled {
			t.Errorf("Slot %d should be enabled", i)
		}
	}

	// Verify back button exists and is enabled
	backIdx := len(mm.items) - 1
	if !mm.items[backIdx].Enabled {
		t.Errorf("Back button should be enabled")
	}
	if !strings.Contains(mm.items[backIdx].Text, "Back") {
		t.Errorf("Last item should be Back button, got: %s", mm.items[backIdx].Text)
	}
}

func TestSaveLoadMenuEmptySlots(t *testing.T) {
	mm := NewMenuManager()

	// Create save manager without any saves
	saveManager, err := save.NewSaveManager("")
	if err != nil {
		t.Fatalf("Failed to create save manager: %v", err)
	}
	mm.saveManager = saveManager

	mm.ShowSaveLoadMenu()

	// Verify all slots show "Empty"
	for i := 0; i < 5; i++ {
		slotText := mm.items[i].Text
		if !strings.Contains(slotText, "Empty") {
			t.Errorf("Slot %d should show 'Empty', got: %s", i, slotText)
		}
	}
}

func TestSaveLoadMenuMultipleSaves(t *testing.T) {
	mm := NewMenuManager()

	saveManager, err := save.NewSaveManager("")
	if err != nil {
		t.Fatalf("Failed to create save manager: %v", err)
	}
	mm.saveManager = saveManager

	// Create multiple test saves with different data
	testSaves := []struct {
		slot     int
		seed     int64
		playTime int64
	}{
		{slot: 0, seed: 111, playTime: 60},   // 1 minute
		{slot: 2, seed: 222, playTime: 3600}, // 1 hour
		{slot: 4, seed: 333, playTime: 7265}, // 2 hours, 1 minute, 5 seconds
	}

	for _, ts := range testSaves {
		testSave := &save.SaveData{
			Seed:     ts.seed,
			PlayTime: ts.playTime,
			SlotID:   ts.slot,
		}
		err = mm.saveManager.SaveGame(testSave, ts.slot)
		if err != nil {
			t.Fatalf("Failed to save test data for slot %d: %v", ts.slot, err)
		}
		defer mm.saveManager.DeleteSave(ts.slot)
	}

	mm.ShowSaveLoadMenu()

	// Verify saved slots show metadata
	if !strings.Contains(mm.items[0].Text, "111") {
		t.Errorf("Slot 0 should show seed 111, got: %s", mm.items[0].Text)
	}
	if !strings.Contains(mm.items[2].Text, "222") {
		t.Errorf("Slot 2 should show seed 222, got: %s", mm.items[2].Text)
	}
	if !strings.Contains(mm.items[4].Text, "333") {
		t.Errorf("Slot 4 should show seed 333, got: %s", mm.items[4].Text)
	}

	// Verify empty slots
	emptySlots := []int{1, 3}
	for _, idx := range emptySlots {
		if !strings.Contains(mm.items[idx].Text, "Empty") {
			t.Errorf("Slot %d should show 'Empty', got: %s", idx, mm.items[idx].Text)
		}
	}
}

func TestSaveLoadMenuLoadAction(t *testing.T) {
	mm := NewMenuManager()

	saveManager, err := save.NewSaveManager("")
	if err != nil {
		t.Fatalf("Failed to create save manager: %v", err)
	}
	mm.saveManager = saveManager

	// Track which slot was requested to load
	loadedSlot := -1
	mm.SetCallbacks(
		nil,
		func(slot int) error {
			loadedSlot = slot
			return nil
		},
		nil,
		nil,
		nil,
	)

	mm.ShowSaveLoadMenu()

	// Simulate selecting slot 2
	mm.selectedIndex = 2
	err = mm.selectCurrentItem()
	if err != nil {
		t.Errorf("Failed to select item: %v", err)
	}

	if loadedSlot != 2 {
		t.Errorf("Expected to load slot 2, got slot %d", loadedSlot)
	}
}

func TestSaveLoadMenuBackNavigation(t *testing.T) {
	mm := NewMenuManager()
	mm.ShowSaveLoadMenu()

	// Select back button (last item)
	mm.selectedIndex = len(mm.items) - 1
	err := mm.selectCurrentItem()
	if err != nil {
		t.Errorf("Back action failed: %v", err)
	}

	// Should return to main menu
	if mm.currentMenu != MainMenu {
		t.Errorf("Expected MainMenu after back, got %v", mm.currentMenu)
	}
}

func TestSaveLoadMenuTimeFormatting(t *testing.T) {
	testCases := []struct {
		name     string
		playTime int64
		expected string
	}{
		{
			name:     "Less than 1 minute",
			playTime: 45,
			expected: "0h 0m",
		},
		{
			name:     "Exactly 1 minute",
			playTime: 60,
			expected: "0h 1m",
		},
		{
			name:     "Multiple minutes",
			playTime: 300,
			expected: "0h 5m",
		},
		{
			name:     "Exactly 1 hour",
			playTime: 3600,
			expected: "1h 0m",
		},
		{
			name:     "Hours and minutes",
			playTime: 3665,
			expected: "1h 1m",
		},
		{
			name:     "Multiple hours",
			playTime: 7325,
			expected: "2h 2m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mm := NewMenuManager()

			saveManager, err := save.NewSaveManager("")
			if err != nil {
				t.Fatalf("Failed to create save manager: %v", err)
			}
			mm.saveManager = saveManager

			testSave := &save.SaveData{
				Seed:     999,
				PlayTime: tc.playTime,
				SlotID:   0,
			}

			err = mm.saveManager.SaveGame(testSave, 0)
			if err != nil {
				t.Fatalf("Failed to save test data: %v", err)
			}
			defer mm.saveManager.DeleteSave(0)

			mm.ShowSaveLoadMenu()

			slotText := mm.items[0].Text
			if !strings.Contains(slotText, tc.expected) {
				t.Errorf("Expected time format '%s' in slot text, got: %s", tc.expected, slotText)
			}
		})
	}
}
