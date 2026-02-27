package menu

import (
	"image/color"
	"testing"
)

func TestMenuManagerSetGenre(t *testing.T) {
	testCases := []struct {
		name    string
		genreID string
	}{
		{"Fantasy", "fantasy"},
		{"SciFi", "scifi"},
		{"Horror", "horror"},
		{"Cyberpunk", "cyberpunk"},
		{"PostApoc", "postapoc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mm := NewMenuManager()

			// Set genre
			mm.SetGenre(tc.genreID)

			// Verify genre was set
			if mm.currentGenre != tc.genreID {
				t.Errorf("Expected genre %s, got %s", tc.genreID, mm.currentGenre)
			}

			// Verify colors are not nil
			if mm.backgroundColor == nil {
				t.Error("Background color is nil after SetGenre")
			}
			if mm.textColor == nil {
				t.Error("Text color is nil after SetGenre")
			}
			if mm.selectedColor == nil {
				t.Error("Selected color is nil after SetGenre")
			}
			if mm.disabledColor == nil {
				t.Error("Disabled color is nil after SetGenre")
			}
		})
	}
}

func TestMenuManagerSetGenreUniqueness(t *testing.T) {
	mm := NewMenuManager()

	genres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}
	bgColors := make(map[string]color.Color)

	// Store background colors for each genre
	for _, genre := range genres {
		mm.SetGenre(genre)
		bgColors[genre] = mm.backgroundColor
	}

	// Verify that different genres produce different background colors
	// At least some genres should have distinct colors
	uniqueColors := make(map[color.Color]bool)
	for _, c := range bgColors {
		uniqueColors[c] = true
	}

	// We expect at least 3 unique colors across all genres
	if len(uniqueColors) < 3 {
		t.Errorf("Expected at least 3 unique genre background colors, got %d", len(uniqueColors))
	}
}

func TestMenuManagerSetGenreInvalidDefault(t *testing.T) {
	mm := NewMenuManager()

	// Set invalid genre
	mm.SetGenre("invalid_genre")

	// Should default to fantasy theme colors
	expectedBg := color.RGBA{20, 25, 35, 255}
	if mm.backgroundColor != expectedBg {
		t.Errorf("Expected default fantasy background color %v, got %v", expectedBg, mm.backgroundColor)
	}
}

func TestMenuManagerDefaultGenre(t *testing.T) {
	mm := NewMenuManager()

	// Verify default genre is fantasy
	if mm.currentGenre != "fantasy" {
		t.Errorf("Expected default genre to be 'fantasy', got '%s'", mm.currentGenre)
	}
}

func TestMenuManagerGenreColorsDistinct(t *testing.T) {
	mm := NewMenuManager()

	testCases := []struct {
		genre string
	}{
		{"fantasy"},
		{"scifi"},
		{"horror"},
		{"cyberpunk"},
		{"postapoc"},
	}

	for _, tc := range testCases {
		t.Run(tc.genre, func(t *testing.T) {
			mm.SetGenre(tc.genre)

			// Verify that all four color types are different from each other
			colors := []color.Color{
				mm.backgroundColor,
				mm.textColor,
				mm.selectedColor,
				mm.disabledColor,
			}

			for i := 0; i < len(colors); i++ {
				for j := i + 1; j < len(colors); j++ {
					if colors[i] == colors[j] {
						t.Errorf("Genre %s has duplicate colors: index %d and %d are identical", tc.genre, i, j)
					}
				}
			}
		})
	}
}
