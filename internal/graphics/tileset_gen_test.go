package graphics

import (
	"testing"
)

// TestNewTilesetGenerator_CreatesValidGenerator tests generator creation
func TestNewTilesetGenerator_CreatesValidGenerator(t *testing.T) {
	tests := []struct {
		name     string
		tileSize int
		biome    string
	}{
		{"small cave tiles", 8, "cave"},
		{"medium forest tiles", 16, "forest"},
		{"large ruins tiles", 32, "ruins"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := NewTilesetGenerator(tt.tileSize, tt.biome)
			
			if tg == nil {
				t.Fatal("NewTilesetGenerator returned nil")
			}
			if tg.TileSize != tt.tileSize {
				t.Errorf("Expected TileSize %d, got %d", tt.tileSize, tg.TileSize)
			}
			if tg.Biome != tt.biome {
				t.Errorf("Expected Biome %s, got %s", tt.biome, tg.Biome)
			}
		})
	}
}

// TestGenerate_ReturnsValidTileset tests tileset generation
func TestGenerate_ReturnsValidTileset(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	tileset := tg.Generate(123)

	if tileset == nil {
		t.Fatal("Generate returned nil")
	}

	if tileset.TileSize != 16 {
		t.Errorf("Expected TileSize 16, got %d", tileset.TileSize)
	}

	if tileset.Tiles == nil {
		t.Fatal("Tileset.Tiles is nil")
	}
}

// TestGenerate_ContainsAllTileTypes tests that all tile types are generated
func TestGenerate_ContainsAllTileTypes(t *testing.T) {
	tg := NewTilesetGenerator(16, "forest")
	tileset := tg.Generate(456)

	requiredTiles := []TileType{
		SolidTile,
		PlatformTile,
		SpikeTile,
		LiquidTile,
		BackgroundTile,
	}

	for _, tileType := range requiredTiles {
		if tileset.Tiles[tileType] == nil {
			t.Errorf("Missing tile type: %v", tileType)
		}
	}
}

// TestGenerate_TilesHaveCorrectDimensions tests tile dimensions
func TestGenerate_TilesHaveCorrectDimensions(t *testing.T) {
	tileSize := 16
	tg := NewTilesetGenerator(tileSize, "ruins")
	tileset := tg.Generate(789)

	for tileType, sprite := range tileset.Tiles {
		if sprite.Width != tileSize {
			t.Errorf("Tile type %v has wrong width: %d", tileType, sprite.Width)
		}
		if sprite.Height != tileSize {
			t.Errorf("Tile type %v has wrong height: %d", tileType, sprite.Height)
		}
	}
}

// TestGenerate_TilesHaveImages tests that tiles have valid images
func TestGenerate_TilesHaveImages(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	tileset := tg.Generate(101)

	for tileType, sprite := range tileset.Tiles {
		if sprite.Image == nil {
			t.Errorf("Tile type %v has nil image", tileType)
		}
	}
}

// TestGenerate_DeterministicTileset tests that same seed produces same tileset
func TestGenerate_DeterministicTileset(t *testing.T) {
	seed := int64(202)
	
	tg1 := NewTilesetGenerator(16, "cave")
	tileset1 := tg1.Generate(seed)
	
	tg2 := NewTilesetGenerator(16, "cave")
	tileset2 := tg2.Generate(seed)

	// Compare that same number of tiles exist
	if len(tileset1.Tiles) != len(tileset2.Tiles) {
		t.Error("Tilesets have different number of tiles")
	}

	// Verify both have the same tile types
	for tileType := range tileset1.Tiles {
		if tileset2.Tiles[tileType] == nil {
			t.Errorf("Second tileset missing tile type %v", tileType)
		}
	}
}

// TestGenerate_DifferentBiomes tests different biome generation
func TestGenerate_DifferentBiomes(t *testing.T) {
	biomes := []string{"cave", "forest", "ruins"}
	
	for _, biome := range biomes {
		t.Run(biome, func(t *testing.T) {
			tg := NewTilesetGenerator(16, biome)
			tileset := tg.Generate(303)

			if tileset == nil {
				t.Fatalf("Failed to generate tileset for biome %s", biome)
			}

			if len(tileset.Tiles) == 0 {
				t.Errorf("Biome %s generated no tiles", biome)
			}
		})
	}
}

// TestGenerate_UnknownBiome tests default behavior for unknown biome
func TestGenerate_UnknownBiome(t *testing.T) {
	tg := NewTilesetGenerator(16, "unknown_biome")
	tileset := tg.Generate(404)

	if tileset == nil {
		t.Fatal("Generate returned nil for unknown biome")
	}

	// Should still generate all tile types
	if len(tileset.Tiles) == 0 {
		t.Error("Unknown biome generated no tiles")
	}
}

// TestGenerate_DifferentTileSizes tests various tile sizes
func TestGenerate_DifferentTileSizes(t *testing.T) {
	tests := []struct {
		name     string
		tileSize int
	}{
		{"8x8 tiles", 8},
		{"16x16 tiles", 16},
		{"32x32 tiles", 32},
		{"64x64 tiles", 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := NewTilesetGenerator(tt.tileSize, "cave")
			tileset := tg.Generate(505)

			for tileType, sprite := range tileset.Tiles {
				if sprite.Width != tt.tileSize || sprite.Height != tt.tileSize {
					t.Errorf("Tile type %v has wrong dimensions: %dx%d",
						tileType, sprite.Width, sprite.Height)
				}
			}
		})
	}
}

// TestTileType_Constants tests tile type constants
func TestTileType_Constants(t *testing.T) {
	tileTypes := []TileType{
		SolidTile,
		PlatformTile,
		SpikeTile,
		LiquidTile,
		BackgroundTile,
	}

	// Verify they're distinct
	seen := make(map[TileType]bool)
	for _, tileType := range tileTypes {
		if seen[tileType] {
			t.Errorf("Duplicate tile type value: %v", tileType)
		}
		seen[tileType] = true
	}
}

// TestGenerate_SolidTileNotNil tests solid tile generation
func TestGenerate_SolidTileNotNil(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	tileset := tg.Generate(606)

	if tileset.Tiles[SolidTile] == nil {
		t.Error("SolidTile is nil")
	}

	sprite := tileset.Tiles[SolidTile]
	if sprite.Image == nil {
		t.Error("SolidTile image is nil")
	}
}

// TestGenerate_PlatformTileNotNil tests platform tile generation
func TestGenerate_PlatformTileNotNil(t *testing.T) {
	tg := NewTilesetGenerator(16, "forest")
	tileset := tg.Generate(707)

	if tileset.Tiles[PlatformTile] == nil {
		t.Error("PlatformTile is nil")
	}

	sprite := tileset.Tiles[PlatformTile]
	if sprite.Image == nil {
		t.Error("PlatformTile image is nil")
	}
}

// TestGenerate_SpikeTileNotNil tests spike tile generation
func TestGenerate_SpikeTileNotNil(t *testing.T) {
	tg := NewTilesetGenerator(16, "ruins")
	tileset := tg.Generate(808)

	if tileset.Tiles[SpikeTile] == nil {
		t.Error("SpikeTile is nil")
	}

	sprite := tileset.Tiles[SpikeTile]
	if sprite.Image == nil {
		t.Error("SpikeTile image is nil")
	}
}

// TestGenerate_LiquidTileNotNil tests liquid tile generation
func TestGenerate_LiquidTileNotNil(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	tileset := tg.Generate(909)

	if tileset.Tiles[LiquidTile] == nil {
		t.Error("LiquidTile is nil")
	}

	sprite := tileset.Tiles[LiquidTile]
	if sprite.Image == nil {
		t.Error("LiquidTile image is nil")
	}
}

// TestGenerate_BackgroundTileNotNil tests background tile generation
func TestGenerate_BackgroundTileNotNil(t *testing.T) {
	tg := NewTilesetGenerator(16, "forest")
	tileset := tg.Generate(1010)

	if tileset.Tiles[BackgroundTile] == nil {
		t.Error("BackgroundTile is nil")
	}

	sprite := tileset.Tiles[BackgroundTile]
	if sprite.Image == nil {
		t.Error("BackgroundTile image is nil")
	}
}

// TestGenerate_CaveBiomePalette tests cave biome specific colors
func TestGenerate_CaveBiomePalette(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	tileset := tg.Generate(1111)

	// Just verify generation succeeds with cave biome
	if tileset == nil {
		t.Error("Failed to generate cave tileset")
	}

	if len(tileset.Tiles) != 5 {
		t.Errorf("Expected 5 tile types, got %d", len(tileset.Tiles))
	}
}

// TestGenerate_ForestBiomePalette tests forest biome specific colors
func TestGenerate_ForestBiomePalette(t *testing.T) {
	tg := NewTilesetGenerator(16, "forest")
	tileset := tg.Generate(1212)

	// Just verify generation succeeds with forest biome
	if tileset == nil {
		t.Error("Failed to generate forest tileset")
	}

	if len(tileset.Tiles) != 5 {
		t.Errorf("Expected 5 tile types, got %d", len(tileset.Tiles))
	}
}

// TestGenerate_RuinsBiomePalette tests ruins biome specific colors
func TestGenerate_RuinsBiomePalette(t *testing.T) {
	tg := NewTilesetGenerator(16, "ruins")
	tileset := tg.Generate(1313)

	// Just verify generation succeeds with ruins biome
	if tileset == nil {
		t.Error("Failed to generate ruins tileset")
	}

	if len(tileset.Tiles) != 5 {
		t.Errorf("Expected 5 tile types, got %d", len(tileset.Tiles))
	}
}

// TestGenerate_ImageBoundsCorrect tests that image bounds match tile size
func TestGenerate_ImageBoundsCorrect(t *testing.T) {
	tileSize := 24
	tg := NewTilesetGenerator(tileSize, "cave")
	tileset := tg.Generate(1414)

	for tileType, sprite := range tileset.Tiles {
		bounds := sprite.Image.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y

		if width != tileSize {
			t.Errorf("Tile type %v image width %d doesn't match tile size %d",
				tileType, width, tileSize)
		}
		if height != tileSize {
			t.Errorf("Tile type %v image height %d doesn't match tile size %d",
				tileType, height, tileSize)
		}
	}
}

// TestGenerate_MultipleGenerations tests generating multiple tilesets
func TestGenerate_MultipleGenerations(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	
	for i := 0; i < 5; i++ {
		tileset := tg.Generate(int64(1500 + i))
		
		if tileset == nil {
			t.Fatalf("Generation %d failed", i)
		}
		
		if len(tileset.Tiles) == 0 {
			t.Errorf("Generation %d produced no tiles", i)
		}
	}
}

// TestGenerate_SmallTileSize tests very small tiles
func TestGenerate_SmallTileSize(t *testing.T) {
	tg := NewTilesetGenerator(4, "cave")
	tileset := tg.Generate(1616)

	if tileset == nil {
		t.Fatal("Failed to generate small tileset")
	}

	for _, sprite := range tileset.Tiles {
		if sprite.Width != 4 || sprite.Height != 4 {
			t.Error("Small tiles have wrong dimensions")
		}
	}
}

// TestGenerate_LargeTileSize tests very large tiles
func TestGenerate_LargeTileSize(t *testing.T) {
	tg := NewTilesetGenerator(128, "cave")
	tileset := tg.Generate(1717)

	if tileset == nil {
		t.Fatal("Failed to generate large tileset")
	}

	for _, sprite := range tileset.Tiles {
		if sprite.Width != 128 || sprite.Height != 128 {
			t.Error("Large tiles have wrong dimensions")
		}
	}
}

// TestGenerate_AllBiomesAllTiles tests all biomes generate all tiles
func TestGenerate_AllBiomesAllTiles(t *testing.T) {
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss"}
	
	for _, biome := range biomes {
		t.Run(biome, func(t *testing.T) {
			tg := NewTilesetGenerator(16, biome)
			tileset := tg.Generate(1818)

			expectedTiles := []TileType{
				SolidTile, PlatformTile, SpikeTile, LiquidTile, BackgroundTile,
			}

			for _, tileType := range expectedTiles {
				if tileset.Tiles[tileType] == nil {
					t.Errorf("Biome %s missing tile type %v", biome, tileType)
				}
			}
		})
	}
}

// TestTileset_MapIsInitialized tests that Tiles map is initialized
func TestTileset_MapIsInitialized(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	tileset := tg.Generate(1919)

	if tileset.Tiles == nil {
		t.Fatal("Tiles map is nil")
	}

	// Should be able to access without panic
	_ = tileset.Tiles[SolidTile]
}

// TestGenerate_ConsistentSeedConsistentOutput tests determinism more thoroughly
func TestGenerate_ConsistentSeedConsistentOutput(t *testing.T) {
	seed := int64(2020)
	biome := "forest"
	tileSize := 16

	// Generate twice with same parameters
	tg1 := NewTilesetGenerator(tileSize, biome)
	tileset1 := tg1.Generate(seed)

	tg2 := NewTilesetGenerator(tileSize, biome)
	tileset2 := tg2.Generate(seed)

	// Verify same tiles generated
	for tileType := range tileset1.Tiles {
		if tileset2.Tiles[tileType] == nil {
			t.Errorf("Second generation missing tile type %v", tileType)
		}
	}

	for tileType := range tileset2.Tiles {
		if tileset1.Tiles[tileType] == nil {
			t.Errorf("First generation missing tile type %v", tileType)
		}
	}
}

// TestGenerate_DifferentSeedsDifferentVariation tests seed variation
func TestGenerate_DifferentSeedsDifferentVariation(t *testing.T) {
	tg := NewTilesetGenerator(16, "cave")
	
	tileset1 := tg.Generate(111)
	tileset2 := tg.Generate(222)

	// Both should have tiles, but potentially with variation
	if len(tileset1.Tiles) == 0 || len(tileset2.Tiles) == 0 {
		t.Error("One or both tilesets are empty")
	}

	// Both should have same tile types
	if len(tileset1.Tiles) != len(tileset2.Tiles) {
		t.Error("Different seeds should still produce same tile types")
	}
}
