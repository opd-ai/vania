package audio

import (
	"testing"
)

// TestAudioPlayerGenre tests all genre-related functionality with a single audio context
func TestAudioPlayerGenre(t *testing.T) {
	// Create player once for all subtests
	player, err := NewAudioPlayer()
	if err != nil {
		t.Fatalf("Failed to create audio player: %v", err)
	}
	defer player.Close()

	t.Run("SetGenre", func(t *testing.T) {
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
				player.SetGenre(tc.genreID)

				expectedGenre := tc.genreID
				if tc.genreID == "unknown" {
					expectedGenre = "fantasy"
				}

				if player.currentGenre != expectedGenre {
					t.Errorf("Expected genre %s, got %s", expectedGenre, player.currentGenre)
				}

				if len(player.genreInstruments) == 0 {
					t.Error("Genre instruments map is empty after SetGenre")
				}

				if player.genreSFXVariation <= 0 {
					t.Error("Genre SFX variation should be positive")
				}
			})
		}
	})

	t.Run("InstrumentWeights", func(t *testing.T) {
		genres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}

		for _, genre := range genres {
			t.Run(genre, func(t *testing.T) {
				player.SetGenre(genre)

				instruments := player.GetGenreInstruments()
				if len(instruments) == 0 {
					t.Error("No instruments configured for genre")
					return
				}

				totalWeight := 0.0
				for _, weight := range instruments {
					if weight < 0.0 {
						t.Errorf("Negative weight found: %f", weight)
					}
					totalWeight += weight
				}

				if totalWeight < 0.9 || totalWeight > 1.1 {
					t.Errorf("Total instrument weights (%f) should be close to 1.0", totalWeight)
				}
			})
		}
	})

	t.Run("InstrumentDifferences", func(t *testing.T) {
		player.SetGenre("fantasy")
		fantasyInstruments := copyInstrumentMap(player.GetGenreInstruments())

		player.SetGenre("scifi")
		scifiInstruments := copyInstrumentMap(player.GetGenreInstruments())

		player.SetGenre("horror")
		horrorInstruments := copyInstrumentMap(player.GetGenreInstruments())

		if instrumentMapsEqual(fantasyInstruments, scifiInstruments) {
			t.Error("Fantasy and SciFi should have different instrument preferences")
		}
		if instrumentMapsEqual(fantasyInstruments, horrorInstruments) {
			t.Error("Fantasy and Horror should have different instrument preferences")
		}
		if instrumentMapsEqual(scifiInstruments, horrorInstruments) {
			t.Error("SciFi and Horror should have different instrument preferences")
		}
	})

	t.Run("SFXVariation", func(t *testing.T) {
		testCases := []struct {
			genre         string
			expectHighVar bool
		}{
			{"fantasy", false},
			{"scifi", true},
			{"horror", true},
			{"cyberpunk", true},
			{"postapoc", true},
		}

		for _, tc := range testCases {
			t.Run(tc.genre, func(t *testing.T) {
				player.SetGenre(tc.genre)
				variation := player.GetGenreSFXVariation()

				if variation <= 0 {
					t.Error("SFX variation should be positive")
				}

				if tc.expectHighVar && variation <= 1.0 {
					t.Errorf("Expected high variation (>1.0) for %s, got %f", tc.genre, variation)
				}
			})
		}
	})

	t.Run("DefaultGenre", func(t *testing.T) {
		player.SetGenre("fantasy")

		if player.currentGenre != "fantasy" {
			t.Errorf("Expected default genre to be 'fantasy', got '%s'", player.currentGenre)
		}

		if len(player.genreInstruments) == 0 {
			t.Error("Default genre instruments should be set")
		}
	})

	t.Run("MultipleGenreSwitches", func(t *testing.T) {
		player.SetGenre("scifi")
		if player.currentGenre != "scifi" {
			t.Error("Failed to set genre to scifi")
		}

		player.SetGenre("horror")
		if player.currentGenre != "horror" {
			t.Error("Failed to set genre to horror")
		}

		player.SetGenre("fantasy")
		if player.currentGenre != "fantasy" {
			t.Error("Failed to set genre back to fantasy")
		}

		if len(player.genreInstruments) == 0 {
			t.Error("Instruments should be set after genre switch")
		}
	})
}

func copyInstrumentMap(m map[WaveType]float64) map[WaveType]float64 {
	result := make(map[WaveType]float64)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func instrumentMapsEqual(a, b map[WaveType]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
