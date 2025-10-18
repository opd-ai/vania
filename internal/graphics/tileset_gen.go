package graphics

import (
	"image"
	"image/color"
	"math/rand"
)

// TileType defines different tile categories
type TileType int

const (
	SolidTile TileType = iota
	PlatformTile
	SpikeTile
	LiquidTile
	BackgroundTile
)

// Tileset contains generated tiles
type Tileset struct {
	Tiles    map[TileType]*Sprite
	TileSize int
}

// TilesetGenerator generates tilemap sprites
type TilesetGenerator struct {
	TileSize int
	Biome    string
}

// NewTilesetGenerator creates a new tileset generator
func NewTilesetGenerator(tileSize int, biome string) *TilesetGenerator {
	return &TilesetGenerator{
		TileSize: tileSize,
		Biome:    biome,
	}
}

// Generate creates a complete tileset
func (tg *TilesetGenerator) Generate(seed int64) *Tileset {
	rng := rand.New(rand.NewSource(seed))
	
	tileset := &Tileset{
		Tiles:    make(map[TileType]*Sprite),
		TileSize: tg.TileSize,
	}
	
	// Generate biome-specific palette
	palette := tg.generateBiomePalette(rng)
	
	// Generate each tile type
	tileset.Tiles[SolidTile] = tg.generateSolidTile(rng, palette)
	tileset.Tiles[PlatformTile] = tg.generatePlatformTile(rng, palette)
	tileset.Tiles[SpikeTile] = tg.generateSpikeTile(rng, palette)
	tileset.Tiles[LiquidTile] = tg.generateLiquidTile(rng, palette)
	tileset.Tiles[BackgroundTile] = tg.generateBackgroundTile(rng, palette)
	
	return tileset
}

// generateBiomePalette creates biome-specific colors
func (tg *TilesetGenerator) generateBiomePalette(rng *rand.Rand) []color.RGBA {
	palette := make([]color.RGBA, 4)
	
	// Base colors on biome type
	switch tg.Biome {
	case "cave":
		palette[0] = color.RGBA{40, 40, 50, 255}   // dark stone
		palette[1] = color.RGBA{60, 60, 70, 255}   // lighter stone
		palette[2] = color.RGBA{80, 70, 60, 255}   // brown rock
		palette[3] = color.RGBA{100, 90, 80, 255}  // light rock
	case "forest":
		palette[0] = color.RGBA{34, 80, 34, 255}   // dark green
		palette[1] = color.RGBA{50, 100, 50, 255}  // grass
		palette[2] = color.RGBA{70, 50, 30, 255}   // brown earth
		palette[3] = color.RGBA{90, 70, 50, 255}   // light earth
	case "ruins":
		palette[0] = color.RGBA{100, 90, 80, 255}  // old stone
		palette[1] = color.RGBA{120, 110, 95, 255} // weathered stone
		palette[2] = color.RGBA{80, 85, 90, 255}   // blue-gray
		palette[3] = color.RGBA{140, 130, 115, 255} // light stone
	default:
		// Generic palette
		baseHue := rng.Float64() * 360.0
		for i := range palette {
			palette[i] = hsvToRGB(baseHue+float64(i)*15.0, 0.4, 0.3+float64(i)*0.15)
		}
	}
	
	// Add slight variation
	for i := range palette {
		palette[i].R = uint8(clamp(int(palette[i].R)+rng.Intn(20)-10, 0, 255))
		palette[i].G = uint8(clamp(int(palette[i].G)+rng.Intn(20)-10, 0, 255))
		palette[i].B = uint8(clamp(int(palette[i].B)+rng.Intn(20)-10, 0, 255))
	}
	
	return palette
}

// generateSolidTile creates a solid block tile
func (tg *TilesetGenerator) generateSolidTile(rng *rand.Rand, palette []color.RGBA) *Sprite {
	sprite := &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, tg.TileSize, tg.TileSize)),
		Width:  tg.TileSize,
		Height: tg.TileSize,
	}
	
	// Fill with base color and add texture
	baseColor := palette[0]
	for y := 0; y < tg.TileSize; y++ {
		for x := 0; x < tg.TileSize; x++ {
			// Add noise texture
			if rng.Float64() < 0.2 {
				sprite.Image.Set(x, y, palette[1])
			} else {
				sprite.Image.Set(x, y, baseColor)
			}
		}
	}
	
	// Add edge highlights
	edgeColor := palette[2]
	for x := 0; x < tg.TileSize; x++ {
		sprite.Image.Set(x, 0, edgeColor)
		sprite.Image.Set(x, 1, edgeColor)
	}
	
	return sprite
}

// generatePlatformTile creates a platform tile
func (tg *TilesetGenerator) generatePlatformTile(rng *rand.Rand, palette []color.RGBA) *Sprite {
	sprite := &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, tg.TileSize, tg.TileSize)),
		Width:  tg.TileSize,
		Height: tg.TileSize,
	}
	
	// Transparent background
	transparent := color.RGBA{0, 0, 0, 0}
	for y := 0; y < tg.TileSize; y++ {
		for x := 0; x < tg.TileSize; x++ {
			sprite.Image.Set(x, y, transparent)
		}
	}
	
	// Draw top surface
	platformHeight := 3
	for y := 0; y < platformHeight; y++ {
		for x := 0; x < tg.TileSize; x++ {
			if y == 0 {
				sprite.Image.Set(x, y, palette[2])
			} else {
				sprite.Image.Set(x, y, palette[1])
			}
		}
	}
	
	return sprite
}

// generateSpikeTile creates a spike hazard tile
func (tg *TilesetGenerator) generateSpikeTile(rng *rand.Rand, palette []color.RGBA) *Sprite {
	sprite := &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, tg.TileSize, tg.TileSize)),
		Width:  tg.TileSize,
		Height: tg.TileSize,
	}
	
	// Transparent background
	transparent := color.RGBA{0, 0, 0, 0}
	dangerColor := color.RGBA{180, 50, 50, 255} // Red for danger
	
	for y := 0; y < tg.TileSize; y++ {
		for x := 0; x < tg.TileSize; x++ {
			sprite.Image.Set(x, y, transparent)
		}
	}
	
	// Draw triangular spikes
	spikeCount := 3
	spikeWidth := tg.TileSize / spikeCount
	
	for i := 0; i < spikeCount; i++ {
		startX := i * spikeWidth
		// Draw triangle
		for y := 0; y < tg.TileSize/2; y++ {
			for x := startX; x < startX+spikeWidth; x++ {
				// Check if inside triangle
				relX := x - startX - spikeWidth/2
				if abs(float64(relX)) <= float64(spikeWidth/2-y) {
					sprite.Image.Set(x, tg.TileSize/2+y, dangerColor)
				}
			}
		}
	}
	
	return sprite
}

// generateLiquidTile creates an animated liquid tile
func (tg *TilesetGenerator) generateLiquidTile(rng *rand.Rand, palette []color.RGBA) *Sprite {
	sprite := &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, tg.TileSize, tg.TileSize)),
		Width:  tg.TileSize,
		Height: tg.TileSize,
	}
	
	// Blue/green liquid colors
	liquidColor := color.RGBA{40, 80, 120, 200}
	darkLiquid := color.RGBA{30, 60, 100, 200}
	
	for y := 0; y < tg.TileSize; y++ {
		for x := 0; x < tg.TileSize; x++ {
			// Wave pattern
			if (x+y)%4 < 2 {
				sprite.Image.Set(x, y, liquidColor)
			} else {
				sprite.Image.Set(x, y, darkLiquid)
			}
		}
	}
	
	return sprite
}

// generateBackgroundTile creates a background tile
func (tg *TilesetGenerator) generateBackgroundTile(rng *rand.Rand, palette []color.RGBA) *Sprite {
	sprite := &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, tg.TileSize, tg.TileSize)),
		Width:  tg.TileSize,
		Height: tg.TileSize,
	}
	
	// Darker version of base color
	bgColor := palette[0]
	bgColor.R = bgColor.R / 2
	bgColor.G = bgColor.G / 2
	bgColor.B = bgColor.B / 2
	
	for y := 0; y < tg.TileSize; y++ {
		for x := 0; x < tg.TileSize; x++ {
			sprite.Image.Set(x, y, bgColor)
		}
	}
	
	return sprite
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
