package engine

import (
	"strings"
	"testing"

	"github.com/opd-ai/vania/internal/narrative"
)

func TestGetRoomDescriptions(t *testing.T) {
	themes := []narrative.StoryTheme{
		narrative.FantasyTheme,
		narrative.SciFiTheme,
		narrative.HorrorTheme,
		narrative.MysticalTheme,
		narrative.PostApocTheme,
	}

	for _, theme := range themes {
		t.Run(string(theme), func(t *testing.T) {
			descs := getRoomDescriptions(theme)
			if len(descs) == 0 {
				t.Errorf("No descriptions for theme %s", theme)
			}
			// Check that corridor type exists (common default)
			if _, ok := descs["corridor"]; !ok {
				t.Errorf("Missing corridor description for theme %s", theme)
			}
		})
	}
}

func TestGetGenericRoomDescription(t *testing.T) {
	testCases := []struct {
		roomType string
		theme    narrative.StoryTheme
	}{
		{"unknown", narrative.FantasyTheme},
		{"test_room", narrative.SciFiTheme},
		{"boss_arena", narrative.HorrorTheme},
	}

	for _, tc := range testCases {
		t.Run(tc.roomType+"_"+string(tc.theme), func(t *testing.T) {
			desc := getGenericRoomDescription(tc.roomType, tc.theme)
			if desc == "" {
				t.Error("Generic description should not be empty")
			}
			if !strings.Contains(desc, tc.roomType) {
				t.Errorf("Description should contain room type %s, got %s", tc.roomType, desc)
			}
		})
	}
}

func TestRoomDescriptionTimerConstants(t *testing.T) {
	// Ensure description shows for reasonable time
	if roomDescriptionDuration < 60 {
		t.Error("Room description duration too short (< 1 second)")
	}
	if roomDescriptionDuration > 600 {
		t.Error("Room description duration too long (> 10 seconds)")
	}
	// Ensure fade is shorter than total duration
	if roomDescriptionFadeDuration >= roomDescriptionDuration/2 {
		t.Error("Fade duration should be less than half of total duration")
	}
}
