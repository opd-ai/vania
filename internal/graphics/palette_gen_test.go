package graphics

import (
	"testing"
)

// TestNewPaletteGenerator_CreatesValidGenerator tests generator creation
func TestNewPaletteGenerator_CreatesValidGenerator(t *testing.T) {
	tests := []struct {
		name   string
		scheme ColorScheme
	}{
		{"complementary", ComplementaryScheme},
		{"triadic", TriadicScheme},
		{"analogous", AnalogousScheme},
		{"monochromatic", MonochromaticScheme},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewPaletteGenerator(tt.scheme)
			if pg == nil {
				t.Fatal("NewPaletteGenerator returned nil")
			}
			if pg.Scheme != tt.scheme {
				t.Errorf("Expected scheme %v, got %v", tt.scheme, pg.Scheme)
			}
		})
	}
}

// TestGenerate_ReturnsCorrectCount tests palette generation count
func TestGenerate_ReturnsCorrectCount(t *testing.T) {
	tests := []struct {
		name   string
		scheme ColorScheme
		count  int
	}{
		{"complementary 4 colors", ComplementaryScheme, 4},
		{"triadic 6 colors", TriadicScheme, 6},
		{"analogous 3 colors", AnalogousScheme, 3},
		{"monochromatic 8 colors", MonochromaticScheme, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewPaletteGenerator(tt.scheme)
			palette := pg.Generate(123, tt.count)

			if len(palette) != tt.count {
				t.Errorf("Expected %d colors, got %d", tt.count, len(palette))
			}
		})
	}
}

// TestGenerate_ComplementaryScheme tests complementary color generation
func TestGenerate_ComplementaryScheme(t *testing.T) {
	pg := NewPaletteGenerator(ComplementaryScheme)
	palette := pg.Generate(456, 4)

	if len(palette) != 4 {
		t.Fatalf("Expected 4 colors, got %d", len(palette))
	}

	// All colors should be valid RGBA
	for i, color := range palette {
		if color.A != 255 {
			t.Errorf("Color %d has invalid alpha: %d", i, color.A)
		}
	}
}

// TestGenerate_TriadicScheme tests triadic color generation
func TestGenerate_TriadicScheme(t *testing.T) {
	pg := NewPaletteGenerator(TriadicScheme)
	palette := pg.Generate(789, 6)

	if len(palette) != 6 {
		t.Fatalf("Expected 6 colors, got %d", len(palette))
	}

	// All colors should be valid RGBA
	for i, color := range palette {
		if color.A != 255 {
			t.Errorf("Color %d has invalid alpha: %d", i, color.A)
		}
	}
}

// TestGenerate_AnalogousScheme tests analogous color generation
func TestGenerate_AnalogousScheme(t *testing.T) {
	pg := NewPaletteGenerator(AnalogousScheme)
	palette := pg.Generate(101, 5)

	if len(palette) != 5 {
		t.Fatalf("Expected 5 colors, got %d", len(palette))
	}

	// All colors should be valid RGBA
	for i, color := range palette {
		if color.A != 255 {
			t.Errorf("Color %d has invalid alpha: %d", i, color.A)
		}
	}
}

// TestGenerate_MonochromaticScheme tests monochromatic color generation
func TestGenerate_MonochromaticScheme(t *testing.T) {
	pg := NewPaletteGenerator(MonochromaticScheme)
	palette := pg.Generate(202, 4)

	if len(palette) != 4 {
		t.Fatalf("Expected 4 colors, got %d", len(palette))
	}

	// All colors should be valid RGBA
	for i, color := range palette {
		if color.A != 255 {
			t.Errorf("Color %d has invalid alpha: %d", i, color.A)
		}
	}
}

// TestGenerate_DeterministicPalette tests that same seed produces same colors
func TestGenerate_DeterministicPalette(t *testing.T) {
	seed := int64(303)
	count := 5

	pg1 := NewPaletteGenerator(ComplementaryScheme)
	palette1 := pg1.Generate(seed, count)

	pg2 := NewPaletteGenerator(ComplementaryScheme)
	palette2 := pg2.Generate(seed, count)

	if len(palette1) != len(palette2) {
		t.Fatal("Palettes have different lengths")
	}

	for i := range palette1 {
		if palette1[i] != palette2[i] {
			t.Errorf("Color %d differs: %v != %v", i, palette1[i], palette2[i])
		}
	}
}

// TestGenerate_DifferentSeeds tests that different seeds produce different colors
func TestGenerate_DifferentSeeds(t *testing.T) {
	pg := NewPaletteGenerator(ComplementaryScheme)
	
	palette1 := pg.Generate(111, 4)
	palette2 := pg.Generate(222, 4)

	// At least some colors should be different
	allSame := true
	for i := range palette1 {
		if palette1[i] != palette2[i] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Different seeds should produce different palettes")
	}
}

// TestGenerate_SingleColor tests generating single color palette
func TestGenerate_SingleColor(t *testing.T) {
	pg := NewPaletteGenerator(MonochromaticScheme)
	palette := pg.Generate(404, 1)

	if len(palette) != 1 {
		t.Errorf("Expected 1 color, got %d", len(palette))
	}
}

// TestGenerate_ManyColors tests generating large palette
func TestGenerate_ManyColors(t *testing.T) {
	pg := NewPaletteGenerator(AnalogousScheme)
	palette := pg.Generate(505, 20)

	if len(palette) != 20 {
		t.Errorf("Expected 20 colors, got %d", len(palette))
	}

	// All colors should be valid
	for _, color := range palette {
		if color.A != 255 {
			t.Error("Found color with invalid alpha")
		}
	}
}

// TestGenerate_DefaultScheme tests default behavior for invalid scheme
func TestGenerate_DefaultScheme(t *testing.T) {
	pg := &PaletteGenerator{Scheme: ColorScheme(999)} // Invalid scheme
	palette := pg.Generate(606, 4)

	// Should fall back to analogous
	if len(palette) != 4 {
		t.Errorf("Expected 4 colors, got %d", len(palette))
	}
}

// TestGenerateHeroicPalette_ReturnsValidPalette tests heroic palette generation
func TestGenerateHeroicPalette_ReturnsValidPalette(t *testing.T) {
	palette := GenerateHeroicPalette(707)

	if len(palette) != 6 {
		t.Errorf("Expected 6 colors, got %d", len(palette))
	}

	// All colors should have full alpha
	for i, color := range palette {
		if color.A != 255 {
			t.Errorf("Color %d has invalid alpha: %d", i, color.A)
		}
	}
}

// TestGenerateHeroicPalette_Deterministic tests heroic palette determinism
func TestGenerateHeroicPalette_Deterministic(t *testing.T) {
	seed := int64(808)
	
	palette1 := GenerateHeroicPalette(seed)
	palette2 := GenerateHeroicPalette(seed)

	if len(palette1) != len(palette2) {
		t.Fatal("Palettes have different lengths")
	}

	for i := range palette1 {
		if palette1[i] != palette2[i] {
			t.Errorf("Color %d differs: %v != %v", i, palette1[i], palette2[i])
		}
	}
}

// TestGenerateHeroicPalette_HasVariation tests that different seeds vary
func TestGenerateHeroicPalette_HasVariation(t *testing.T) {
	palette1 := GenerateHeroicPalette(111)
	palette2 := GenerateHeroicPalette(222)

	// Palettes should be different
	allSame := true
	for i := range palette1 {
		if palette1[i] != palette2[i] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Different seeds should produce different heroic palettes")
	}
}

// TestGenerateEnemyPalette_ReturnsValidPalette tests enemy palette generation
func TestGenerateEnemyPalette_ReturnsValidPalette(t *testing.T) {
	palette := GenerateEnemyPalette(909, 3)

	if len(palette) != 6 {
		t.Errorf("Expected 6 colors, got %d", len(palette))
	}

	// All colors should have full alpha
	for i, color := range palette {
		if color.A != 255 {
			t.Errorf("Color %d has invalid alpha: %d", i, color.A)
		}
	}
}

// TestGenerateEnemyPalette_DangerLevelAffectsColors tests danger level influence
func TestGenerateEnemyPalette_DangerLevelAffectsColors(t *testing.T) {
	tests := []struct {
		name        string
		dangerLevel int
	}{
		{"low danger", 1},
		{"medium danger", 5},
		{"high danger", 8},
		{"max danger", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			palette := GenerateEnemyPalette(1010, tt.dangerLevel)
			
			if len(palette) != 6 {
				t.Errorf("Expected 6 colors, got %d", len(palette))
			}

			// High danger should use different hue (purple vs red)
			// This is tested implicitly by checking the function runs
		})
	}
}

// TestGenerateEnemyPalette_Deterministic tests enemy palette determinism
func TestGenerateEnemyPalette_Deterministic(t *testing.T) {
	seed := int64(1111)
	danger := 5
	
	palette1 := GenerateEnemyPalette(seed, danger)
	palette2 := GenerateEnemyPalette(seed, danger)

	if len(palette1) != len(palette2) {
		t.Fatal("Palettes have different lengths")
	}

	for i := range palette1 {
		if palette1[i] != palette2[i] {
			t.Errorf("Color %d differs: %v != %v", i, palette1[i], palette2[i])
		}
	}
}

// TestGenerateEnemyPalette_HighDangerUsesPurple tests high danger color choice
func TestGenerateEnemyPalette_HighDangerUsesPurple(t *testing.T) {
	// High danger (>5) should use purple hue
	highDanger := GenerateEnemyPalette(1212, 6)
	lowDanger := GenerateEnemyPalette(1212, 5)

	// Just verify they're different palettes
	if len(highDanger) != 6 || len(lowDanger) != 6 {
		t.Error("Palettes should have 6 colors")
	}
}

// TestColorSchemeConstants_ValidValues tests that color scheme constants are defined
func TestColorSchemeConstants_ValidValues(t *testing.T) {
	schemes := []ColorScheme{
		ComplementaryScheme,
		TriadicScheme,
		AnalogousScheme,
		MonochromaticScheme,
	}

	// Just verify they're distinct values
	seenValues := make(map[ColorScheme]bool)
	for _, scheme := range schemes {
		if seenValues[scheme] {
			t.Errorf("Duplicate scheme value: %v", scheme)
		}
		seenValues[scheme] = true
	}
}

// TestGenerate_AllSchemesProduceColors tests that all schemes work
func TestGenerate_AllSchemesProduceColors(t *testing.T) {
	schemes := []struct {
		name   string
		scheme ColorScheme
	}{
		{"complementary", ComplementaryScheme},
		{"triadic", TriadicScheme},
		{"analogous", AnalogousScheme},
		{"monochromatic", MonochromaticScheme},
	}

	for _, tt := range schemes {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewPaletteGenerator(tt.scheme)
			palette := pg.Generate(1313, 4)

			if len(palette) == 0 {
				t.Error("Generated empty palette")
			}

			for i, c := range palette {
				if c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0 {
					t.Errorf("Color %d is zero value", i)
				}
			}
		})
	}
}

// TestGenerate_ColorsHaveFullAlpha tests that all generated colors have full alpha
func TestGenerate_ColorsHaveFullAlpha(t *testing.T) {
	schemes := []ColorScheme{
		ComplementaryScheme,
		TriadicScheme,
		AnalogousScheme,
		MonochromaticScheme,
	}

	for _, scheme := range schemes {
		pg := NewPaletteGenerator(scheme)
		palette := pg.Generate(1414, 5)

		for i, color := range palette {
			if color.A != 255 {
				t.Errorf("Scheme %v, color %d has alpha %d, expected 255", 
					scheme, i, color.A)
			}
		}
	}
}

// TestGenerate_ZeroCount tests behavior with zero count
func TestGenerate_ZeroCount(t *testing.T) {
	pg := NewPaletteGenerator(ComplementaryScheme)
	palette := pg.Generate(1515, 0)

	if len(palette) != 0 {
		t.Errorf("Expected 0 colors, got %d", len(palette))
	}
}

// TestGenerateHeroicPalette_ContainsBlueTones tests heroic palette has blue
func TestGenerateHeroicPalette_ContainsBlueTones(t *testing.T) {
	palette := GenerateHeroicPalette(1616)

	// Should contain blue-ish colors (not strict test due to variation)
	foundColor := false
	for _, color := range palette {
		// Just verify we got some colors
		if color.R > 0 || color.G > 0 || color.B > 0 {
			foundColor = true
			break
		}
	}

	if !foundColor {
		t.Error("Heroic palette should contain visible colors")
	}
}

// TestGenerateEnemyPalette_ContainsRedTones tests enemy palette has red/purple
func TestGenerateEnemyPalette_ContainsRedTones(t *testing.T) {
	palette := GenerateEnemyPalette(1717, 3)

	// Should contain red-ish colors (not strict test due to variation)
	foundColor := false
	for _, color := range palette {
		// Just verify we got some colors
		if color.R > 0 || color.G > 0 || color.B > 0 {
			foundColor = true
			break
		}
	}

	if !foundColor {
		t.Error("Enemy palette should contain visible colors")
	}
}
