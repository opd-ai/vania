package engine

import (
	"testing"
)

func TestNewGameGeneratorWithGenre(t *testing.T) {
	testCases := []struct {
		name    string
		seed    int64
		genreID string
	}{
		{"Fantasy", 42, "fantasy"},
		{"SciFi", 123, "scifi"},
		{"Horror", 456, "horror"},
		{"Cyberpunk", 789, "cyberpunk"},
		{"PostApoc", 999, "postapoc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gen := NewGameGeneratorWithGenre(tc.seed, tc.genreID)

			// Verify generator was created
			if gen == nil {
				t.Fatal("Generator is nil")
			}

			// Verify seed was set
			if gen.MasterSeed != tc.seed {
				t.Errorf("Expected seed %d, got %d", tc.seed, gen.MasterSeed)
			}

			// Verify genre was set
			if gen.Genre != tc.genreID {
				t.Errorf("Expected genre %s, got %s", tc.genreID, gen.Genre)
			}

			// Verify subgenerators were initialized
			if gen.GraphicsGen == nil {
				t.Error("GraphicsGen is nil")
			}
			if gen.AudioGen == nil {
				t.Error("AudioGen is nil")
			}
			if gen.NarrativeGen == nil {
				t.Error("NarrativeGen is nil")
			}
			if gen.WorldGen == nil {
				t.Error("WorldGen is nil")
			}
			if gen.EntityGen == nil {
				t.Error("EntityGen is nil")
			}
			if gen.PCGContext == nil {
				t.Error("PCGContext is nil")
			}
		})
	}
}

func TestNewGameGeneratorDefaultGenre(t *testing.T) {
	seed := int64(42)
	gen := NewGameGenerator(seed)

	// Verify default genre is fantasy
	if gen.Genre != "fantasy" {
		t.Errorf("Expected default genre to be 'fantasy', got '%s'", gen.Genre)
	}
}

func TestGenerateCompleteGameWithGenre(t *testing.T) {
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

	seed := int64(42)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gen := NewGameGeneratorWithGenre(seed, tc.genreID)
			game, err := gen.GenerateCompleteGame()
			// Verify game was generated without error
			if err != nil {
				t.Fatalf("Failed to generate game: %v", err)
			}

			// Verify game structure
			if game == nil {
				t.Fatal("Game is nil")
			}

			// Verify genre was propagated to game
			if game.Genre != tc.genreID {
				t.Errorf("Expected game genre %s, got %s", tc.genreID, game.Genre)
			}

			// Verify game seed matches
			if game.Seed != seed {
				t.Errorf("Expected game seed %d, got %d", seed, game.Seed)
			}

			// Verify all systems were generated
			if game.World == nil {
				t.Error("World is nil")
			}
			if game.Graphics == nil {
				t.Error("Graphics is nil")
			}
			if game.Audio == nil {
				t.Error("Audio is nil")
			}
			if game.Narrative == nil {
				t.Error("Narrative is nil")
			}
			if game.Player == nil {
				t.Error("Player is nil")
			}
		})
	}
}

func TestGenerateCompleteGameDeterminismWithGenre(t *testing.T) {
	seed := int64(123456)
	genre := "scifi"

	gen1 := NewGameGeneratorWithGenre(seed, genre)
	game1, err1 := gen1.GenerateCompleteGame()

	gen2 := NewGameGeneratorWithGenre(seed, genre)
	game2, err2 := gen2.GenerateCompleteGame()

	// Both should generate without error
	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to generate games: %v, %v", err1, err2)
	}

	// Both should have same seed and genre
	if game1.Seed != game2.Seed {
		t.Errorf("Seeds don't match: %d vs %d", game1.Seed, game2.Seed)
	}
	if game1.Genre != game2.Genre {
		t.Errorf("Genres don't match: %s vs %s", game1.Genre, game2.Genre)
	}

	// Both should have same number of rooms (determinism check)
	if len(game1.World.Rooms) != len(game2.World.Rooms) {
		t.Errorf("Room count mismatch: %d vs %d", len(game1.World.Rooms), len(game2.World.Rooms))
	}

	// Both should have same number of entities (determinism check)
	if len(game1.Entities) != len(game2.Entities) {
		t.Errorf("Entity count mismatch: %d vs %d", len(game1.Entities), len(game2.Entities))
	}
}

func TestGenerateCompleteGameGenreIndependence(t *testing.T) {
	seed := int64(42)

	// Generate games with different genres but same seed
	genFantasy := NewGameGeneratorWithGenre(seed, "fantasy")
	gameFantasy, _ := genFantasy.GenerateCompleteGame()

	genSciFi := NewGameGeneratorWithGenre(seed, "scifi")
	gameSciFi, _ := genSciFi.GenerateCompleteGame()

	// Games should have different genres
	if gameFantasy.Genre == gameSciFi.Genre {
		t.Error("Games should have different genres")
	}

	// Games should have same seed (same source of randomness)
	if gameFantasy.Seed != gameSciFi.Seed {
		t.Errorf("Seeds should match: %d vs %d", gameFantasy.Seed, gameSciFi.Seed)
	}

	// World structure should be the same (genre only affects visuals/audio)
	// Room count should be identical since world gen doesn't depend on genre
	if len(gameFantasy.World.Rooms) != len(gameSciFi.World.Rooms) {
		t.Errorf("Room count should match across genres: %d vs %d",
			len(gameFantasy.World.Rooms), len(gameSciFi.World.Rooms))
	}
}
