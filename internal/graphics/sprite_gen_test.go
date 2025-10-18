package graphics

import (
	"testing"
)

func TestSpriteGeneration(t *testing.T) {
	gen := NewSpriteGenerator(16, 16, VerticalSymmetry)
	
	sprite := gen.Generate(12345)
	
	if sprite == nil {
		t.Fatal("Generated sprite is nil")
	}
	
	if sprite.Width != 16 || sprite.Height != 16 {
		t.Errorf("Sprite size mismatch: got %dx%d, want 16x16", sprite.Width, sprite.Height)
	}
	
	if sprite.Image == nil {
		t.Error("Sprite image not initialized")
	}
}

func TestSpriteDeterminism(t *testing.T) {
	gen := NewSpriteGenerator(32, 32, HorizontalSymmetry)
	seed := int64(999)
	
	sprite1 := gen.Generate(seed)
	sprite2 := gen.Generate(seed)
	
	// Compare a few pixels to verify determinism
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			c1 := sprite1.Image.RGBAAt(x, y)
			c2 := sprite2.Image.RGBAAt(x, y)
			
			if c1.R != c2.R || c1.G != c2.G || c1.B != c2.B || c1.A != c2.A {
				t.Errorf("Sprites not deterministic at (%d,%d)", x, y)
				return
			}
		}
	}
}

func TestPaletteGeneration(t *testing.T) {
	gen := NewPaletteGenerator(ComplementaryScheme)
	
	palette := gen.Generate(42, 6)
	
	if len(palette) != 6 {
		t.Errorf("Palette size mismatch: got %d, want 6", len(palette))
	}
	
	// Check that colors are valid
	for i, color := range palette {
		if color.A == 0 {
			t.Errorf("Color %d has zero alpha", i)
		}
	}
}

func TestTilesetGeneration(t *testing.T) {
	gen := NewTilesetGenerator(16, "cave")
	
	tileset := gen.Generate(777)
	
	if tileset == nil {
		t.Fatal("Generated tileset is nil")
	}
	
	if tileset.TileSize != 16 {
		t.Errorf("Tile size mismatch: got %d, want 16", tileset.TileSize)
	}
	
	// Check that all tile types are generated
	expectedTypes := []TileType{SolidTile, PlatformTile, SpikeTile, LiquidTile, BackgroundTile}
	for _, tileType := range expectedTypes {
		if _, ok := tileset.Tiles[tileType]; !ok {
			t.Errorf("Missing tile type: %v", tileType)
		}
	}
}

func TestHSVToRGB(t *testing.T) {
	// Test pure red
	red := hsvToRGB(0, 1.0, 1.0)
	if red.R != 255 || red.G != 0 || red.B != 0 {
		t.Errorf("HSV to RGB conversion failed for red: got RGB(%d,%d,%d)", red.R, red.G, red.B)
	}
	
	// Test white
	white := hsvToRGB(0, 0, 1.0)
	if white.R != 255 || white.G != 255 || white.B != 255 {
		t.Errorf("HSV to RGB conversion failed for white: got RGB(%d,%d,%d)", white.R, white.G, white.B)
	}
	
	// Test black
	black := hsvToRGB(0, 0, 0)
	if black.R != 0 || black.G != 0 || black.B != 0 {
		t.Errorf("HSV to RGB conversion failed for black: got RGB(%d,%d,%d)", black.R, black.G, black.B)
	}
}

func TestSymmetryTypes(t *testing.T) {
	testCases := []struct {
		name     string
		symmetry SymmetryType
	}{
		{"NoSymmetry", NoSymmetry},
		{"HorizontalSymmetry", HorizontalSymmetry},
		{"VerticalSymmetry", VerticalSymmetry},
		{"RadialSymmetry", RadialSymmetry},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gen := NewSpriteGenerator(16, 16, tc.symmetry)
			sprite := gen.Generate(555)
			
			if sprite == nil {
				t.Errorf("Failed to generate sprite with %s", tc.name)
			}
		})
	}
}
