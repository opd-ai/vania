package render

import (
	"testing"
)

func TestRendererSetGenre(t *testing.T) {
	testCases := []struct {
		name    string
		genreID string
	}{
		{"Fantasy", "fantasy"},
		{"SciFi", "scifi"},
		{"Horror", "horror"},
		{"Cyberpunk", "cyberpunk"},
		{"PostApoc", "postapoc"},
		{"Unknown", "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			renderer := NewRenderer()

			// Set genre
			renderer.SetGenre(tc.genreID)

			// Verify genre was set
			if renderer.currentGenre != tc.genreID {
				t.Errorf("Expected genre %s, got %s", tc.genreID, renderer.currentGenre)
			}

			// Verify background color was updated
			if renderer.bgColor == nil {
				t.Error("Background color is nil after SetGenre")
			}

			// Verify icon cache was cleared
			if len(renderer.abilityIconCache) != 0 {
				t.Error("Ability icon cache should be cleared after SetGenre")
			}
			if len(renderer.lastAbilities) != 0 {
				t.Error("Last abilities cache should be cleared after SetGenre")
			}
		})
	}
}

func TestRendererSetGenreBackgroundColors(t *testing.T) {
	renderer := NewRenderer()

	genres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}

	// Store background colors for each genre
	bgColors := make(map[string]interface{})

	for _, genre := range genres {
		renderer.SetGenre(genre)
		bgColors[genre] = renderer.bgColor
	}

	// Verify that different genres produce different background colors
	// At least some genres should have distinct colors
	fantasyBg := bgColors["fantasy"]
	scifiBg := bgColors["scifi"]

	if fantasyBg == scifiBg {
		// Colors might be the same by chance, but let's check if they're all identical
		allSame := true
		for _, bg := range bgColors {
			if bg != fantasyBg {
				allSame = false
				break
			}
		}
		if allSame {
			t.Error("All genre background colors are identical - genres should have distinct visual themes")
		}
	}
}

func TestRendererSetGenreInvalidatesCache(t *testing.T) {
	renderer := NewRenderer()

	// Populate caches
	renderer.abilityIconCache["test"] = nil
	renderer.lastAbilities["test"] = true

	if len(renderer.abilityIconCache) == 0 || len(renderer.lastAbilities) == 0 {
		t.Fatal("Failed to populate test caches")
	}

	// Set genre should clear caches
	renderer.SetGenre("scifi")

	if len(renderer.abilityIconCache) != 0 {
		t.Error("SetGenre should clear ability icon cache")
	}
	if len(renderer.lastAbilities) != 0 {
		t.Error("SetGenre should clear last abilities cache")
	}
}

func TestRendererDefaultGenre(t *testing.T) {
	renderer := NewRenderer()

	if renderer.currentGenre != "fantasy" {
		t.Errorf("Expected default genre to be 'fantasy', got '%s'", renderer.currentGenre)
	}
}
