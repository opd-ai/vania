package graphics

import (
	"image/color"
	"math/rand"
)

// ColorScheme represents a color palette scheme
type ColorScheme int

const (
	ComplementaryScheme ColorScheme = iota
	TriadicScheme
	AnalogousScheme
	MonochromaticScheme
)

// PaletteGenerator generates color palettes
type PaletteGenerator struct {
	Scheme ColorScheme
}

// NewPaletteGenerator creates a new palette generator
func NewPaletteGenerator(scheme ColorScheme) *PaletteGenerator {
	return &PaletteGenerator{
		Scheme: scheme,
	}
}

// Generate creates a color palette
func (pg *PaletteGenerator) Generate(seed int64, count int) []color.RGBA {
	rng := rand.New(rand.NewSource(seed))

	// Generate base hue
	baseHue := rng.Float64() * 360.0

	switch pg.Scheme {
	case ComplementaryScheme:
		return pg.generateComplementary(baseHue, count)
	case TriadicScheme:
		return pg.generateTriadic(baseHue, count)
	case AnalogousScheme:
		return pg.generateAnalogous(baseHue, count)
	case MonochromaticScheme:
		return pg.generateMonochromatic(baseHue, count)
	default:
		return pg.generateAnalogous(baseHue, count)
	}
}

// generateComplementary creates complementary color scheme
func (pg *PaletteGenerator) generateComplementary(baseHue float64, count int) []color.RGBA {
	palette := make([]color.RGBA, count)

	for i := 0; i < count; i++ {
		var hue float64
		if i%2 == 0 {
			hue = baseHue
		} else {
			hue = mod(baseHue+180.0, 360.0)
		}

		saturation := 0.5 + float64(i%3)*0.15
		value := 0.4 + float64(i%4)*0.15

		palette[i] = hsvToRGB(hue, saturation, value)
	}

	return palette
}

// generateTriadic creates triadic color scheme
func (pg *PaletteGenerator) generateTriadic(baseHue float64, count int) []color.RGBA {
	palette := make([]color.RGBA, count)

	for i := 0; i < count; i++ {
		hue := mod(baseHue+float64(i%3)*120.0, 360.0)
		saturation := 0.6 + float64(i/3)*0.1
		value := 0.5 + float64(i%2)*0.2

		palette[i] = hsvToRGB(hue, saturation, value)
	}

	return palette
}

// generateAnalogous creates analogous color scheme
func (pg *PaletteGenerator) generateAnalogous(baseHue float64, count int) []color.RGBA {
	palette := make([]color.RGBA, count)

	for i := 0; i < count; i++ {
		hue := mod(baseHue+float64(i)*30.0-45.0, 360.0)
		saturation := 0.5 + float64(i%3)*0.15
		value := 0.4 + float64(i)*0.1

		palette[i] = hsvToRGB(hue, saturation, value)
	}

	return palette
}

// generateMonochromatic creates monochromatic color scheme
func (pg *PaletteGenerator) generateMonochromatic(baseHue float64, count int) []color.RGBA {
	palette := make([]color.RGBA, count)

	for i := 0; i < count; i++ {
		saturation := 0.3 + float64(i)*0.1
		value := 0.2 + float64(i)*0.15

		palette[i] = hsvToRGB(baseHue, saturation, value)
	}

	return palette
}

// GenerateHeroicPalette creates a palette suitable for hero characters
func GenerateHeroicPalette(seed int64) []color.RGBA {
	rng := rand.New(rand.NewSource(seed))

	// Blues and golds typical of heroes
	palette := []color.RGBA{
		{30, 60, 120, 255},   // Deep blue
		{50, 90, 180, 255},   // Bright blue
		{180, 140, 60, 255},  // Gold
		{220, 180, 80, 255},  // Bright gold
		{200, 200, 200, 255}, // Silver/white
		{80, 80, 90, 255},    // Dark accent
	}

	// Add variation
	for i := range palette {
		palette[i].R = uint8(clamp(int(palette[i].R)+rng.Intn(30)-15, 0, 255))
		palette[i].G = uint8(clamp(int(palette[i].G)+rng.Intn(30)-15, 0, 255))
		palette[i].B = uint8(clamp(int(palette[i].B)+rng.Intn(30)-15, 0, 255))
	}

	return palette
}

// GenerateEnemyPalette creates a palette suitable for enemies
func GenerateEnemyPalette(seed int64, dangerLevel int) []color.RGBA {
	rng := rand.New(rand.NewSource(seed))

	// Reds, purples, and dark colors for enemies
	baseHue := 0.0 // Red
	if dangerLevel > 5 {
		baseHue = 280.0 // Purple for high danger
	}

	palette := make([]color.RGBA, 6)
	for i := range palette {
		hue := mod(baseHue+float64(i)*15.0, 360.0)
		saturation := 0.6 + float64(dangerLevel)*0.05
		value := 0.3 + float64(i)*0.1

		palette[i] = hsvToRGB(hue, saturation, value)
	}

	// Add variation
	for i := range palette {
		palette[i].R = uint8(clamp(int(palette[i].R)+rng.Intn(20)-10, 0, 255))
		palette[i].G = uint8(clamp(int(palette[i].G)+rng.Intn(20)-10, 0, 255))
		palette[i].B = uint8(clamp(int(palette[i].B)+rng.Intn(20)-10, 0, 255))
	}

	return palette
}

// GenerateGenrePalette creates a genre-specific color palette
func GenerateGenrePalette(genreID string, seed int64, count int) []color.RGBA {
	rng := rand.New(rand.NewSource(seed))
	palette := make([]color.RGBA, count)

	switch genreID {
	case "fantasy":
		// Earthy greens, browns, gold accents
		baseHue := 100.0 + rng.Float64()*40.0 // Green-yellow range
		for i := 0; i < count; i++ {
			hue := mod(baseHue+float64(i)*25.0, 360.0)
			saturation := 0.45 + float64(i%3)*0.15
			value := 0.35 + float64(i)*0.12
			palette[i] = hsvToRGB(hue, saturation, value)
		}

	case "scifi":
		// Cool blues, cyans, metallic silvers
		baseHue := 180.0 + rng.Float64()*60.0 // Cyan-blue range
		for i := 0; i < count; i++ {
			hue := mod(baseHue+float64(i)*20.0, 360.0)
			saturation := 0.5 + float64(i%2)*0.2
			value := 0.4 + float64(i)*0.1
			palette[i] = hsvToRGB(hue, saturation, value)
		}

	case "horror":
		// Desaturated reds, purples, dark grays
		baseHue := 0.0 + rng.Float64()*30.0 // Red range
		for i := 0; i < count; i++ {
			hue := mod(baseHue+float64(i)*15.0, 360.0)
			saturation := 0.25 + float64(i%3)*0.1 // Low saturation
			value := 0.2 + float64(i)*0.08        // Dark values
			palette[i] = hsvToRGB(hue, saturation, value)
		}

	case "cyberpunk":
		// Hot pinks, neon purples, electric blues
		baseHue := 290.0 + rng.Float64()*50.0 // Magenta-purple range
		for i := 0; i < count; i++ {
			hue := mod(baseHue+float64(i)*30.0, 360.0)
			saturation := 0.7 + float64(i%2)*0.15 // High saturation
			value := 0.5 + float64(i)*0.1
			palette[i] = hsvToRGB(hue, saturation, value)
		}

	case "postapoc":
		// Rusty oranges, muddy browns, dusty grays
		baseHue := 20.0 + rng.Float64()*30.0 // Orange-brown range
		for i := 0; i < count; i++ {
			hue := mod(baseHue+float64(i)*20.0, 360.0)
			saturation := 0.35 + float64(i%3)*0.12 // Muted
			value := 0.3 + float64(i)*0.1
			palette[i] = hsvToRGB(hue, saturation, value)
		}

	default:
		// Fallback to fantasy
		return GenerateGenrePalette("fantasy", seed, count)
	}

	return palette
}
