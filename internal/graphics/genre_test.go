package graphics

import (
	"image/color"
	"testing"
)

func TestGenerateGenrePalette(t *testing.T) {
	testCases := []struct {
		name    string
		genreID string
		seed    int64
		count   int
	}{
		{"Fantasy", "fantasy", 42, 6},
		{"SciFi", "scifi", 42, 6},
		{"Horror", "horror", 42, 6},
		{"Cyberpunk", "cyberpunk", 42, 6},
		{"PostApoc", "postapoc", 42, 6},
		{"Default", "unknown", 42, 6},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			palette := GenerateGenrePalette(tc.genreID, tc.seed, tc.count)

			if len(palette) != tc.count {
				t.Errorf("Expected %d colors, got %d", tc.count, len(palette))
			}

			// Verify all colors are valid (alpha should be 255)
			for i, c := range palette {
				if c.A != 255 {
					t.Errorf("Color %d has alpha %d, expected 255", i, c.A)
				}
			}
		})
	}
}

func TestGenerateGenrePaletteDeterminism(t *testing.T) {
	genres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc"}

	for _, genre := range genres {
		t.Run(genre, func(t *testing.T) {
			seed := int64(999)
			count := 6

			palette1 := GenerateGenrePalette(genre, seed, count)
			palette2 := GenerateGenrePalette(genre, seed, count)

			if len(palette1) != len(palette2) {
				t.Errorf("Palettes have different lengths: %d vs %d", len(palette1), len(palette2))
				return
			}

			for i := range palette1 {
				if palette1[i] != palette2[i] {
					t.Errorf("Color %d differs: %v vs %v", i, palette1[i], palette2[i])
				}
			}
		})
	}
}

func TestGenerateGenrePaletteUniqueness(t *testing.T) {
	// Different genres should produce visually distinct palettes
	seed := int64(42)
	count := 6

	fantasyPalette := GenerateGenrePalette("fantasy", seed, count)
	scifiPalette := GenerateGenrePalette("scifi", seed, count)
	horrorPalette := GenerateGenrePalette("horror", seed, count)

	// Compare average hue ranges (very rough test)
	// This verifies that different genres use different color ranges
	if paletteEquals(fantasyPalette, scifiPalette) {
		t.Error("Fantasy and SciFi palettes should not be identical")
	}
	if paletteEquals(fantasyPalette, horrorPalette) {
		t.Error("Fantasy and Horror palettes should not be identical")
	}
	if paletteEquals(scifiPalette, horrorPalette) {
		t.Error("SciFi and Horror palettes should not be identical")
	}
}

func TestMapGenreToBiome(t *testing.T) {
	testCases := []struct {
		genreID       string
		expectedBiome string
	}{
		{"fantasy", "ruins"},
		{"scifi", "tech"},
		{"horror", "crypt"},
		{"cyberpunk", "neon"},
		{"postapoc", "wasteland"},
		{"unknown", "cave"},
	}

	for _, tc := range testCases {
		t.Run(tc.genreID, func(t *testing.T) {
			biome := MapGenreToBiome(tc.genreID)
			if biome != tc.expectedBiome {
				t.Errorf("Expected biome %s, got %s", tc.expectedBiome, biome)
			}
		})
	}
}

func TestGenerateGenreTileset(t *testing.T) {
	testCases := []struct {
		name     string
		genreID  string
		seed     int64
		tileSize int
	}{
		{"Fantasy", "fantasy", 42, 16},
		{"SciFi", "scifi", 42, 16},
		{"Horror", "horror", 42, 16},
		{"Cyberpunk", "cyberpunk", 42, 16},
		{"PostApoc", "postapoc", 42, 16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tileset := GenerateGenreTileset(tc.genreID, tc.seed, tc.tileSize)

			if tileset == nil {
				t.Fatal("Tileset is nil")
			}

			if tileset.TileSize != tc.tileSize {
				t.Errorf("Expected tile size %d, got %d", tc.tileSize, tileset.TileSize)
			}

			// Verify all tile types are generated
			expectedTiles := []TileType{
				SolidTile,
				PlatformTile,
				SpikeTile,
				LiquidTile,
				BackgroundTile,
			}

			for _, tileType := range expectedTiles {
				tile, ok := tileset.Tiles[tileType]
				if !ok {
					t.Errorf("Missing tile type: %v", tileType)
					continue
				}
				if tile == nil {
					t.Errorf("Tile type %v is nil", tileType)
					continue
				}
				if tile.Image == nil {
					t.Errorf("Tile type %v has nil image", tileType)
				}
			}
		})
	}
}

func TestGenerateGenreTilesetDeterminism(t *testing.T) {
	genreID := "scifi"
	seed := int64(777)
	tileSize := 16

	tileset1 := GenerateGenreTileset(genreID, seed, tileSize)
	tileset2 := GenerateGenreTileset(genreID, seed, tileSize)

	if tileset1.TileSize != tileset2.TileSize {
		t.Errorf("Tile sizes differ: %d vs %d", tileset1.TileSize, tileset2.TileSize)
	}

	// Verify that the same tiles are generated
	for tileType := range tileset1.Tiles {
		tile1 := tileset1.Tiles[tileType]
		tile2 := tileset2.Tiles[tileType]

		if tile1.Width != tile2.Width || tile1.Height != tile2.Height {
			t.Errorf("Tile %v dimensions differ", tileType)
		}

		// Sample a few pixels to verify determinism
		for y := 0; y < 3 && y < tile1.Height; y++ {
			for x := 0; x < 3 && x < tile1.Width; x++ {
				c1 := tile1.Image.RGBAAt(x, y)
				c2 := tile2.Image.RGBAAt(x, y)

				if c1 != c2 {
					t.Errorf("Tile %v pixel (%d,%d) differs: %v vs %v", tileType, x, y, c1, c2)
					return
				}
			}
		}
	}
}

// Helper function to check if two palettes are identical
func paletteEquals(a, b []color.RGBA) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
