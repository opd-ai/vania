// Package graphics provides procedural pixel art generation including sprites,
// tilesets, and color palettes using cellular automata, symmetry transforms,
// and HSV color theory to create coherent retro-style visual assets.
package graphics

import (
	"image"
	"image/color"
	"math/rand"
)

// SymmetryType defines sprite symmetry
type SymmetryType int

const (
	NoSymmetry SymmetryType = iota
	HorizontalSymmetry
	VerticalSymmetry
	RadialSymmetry
)

// SpriteConstraints defines rules for sprite generation
type SpriteConstraints struct {
	MinDensity     float64
	MaxDensity     float64
	RequireOutline bool
	ColorCount     int
}

// Sprite represents a generated sprite
type Sprite struct {
	Image  *image.RGBA
	Width  int
	Height int
}

// SpriteGenerator generates procedural pixel art sprites
type SpriteGenerator struct {
	Width       int
	Height      int
	Symmetry    SymmetryType
	Constraints SpriteConstraints
}

// NewSpriteGenerator creates a new sprite generator
func NewSpriteGenerator(width, height int, symmetry SymmetryType) *SpriteGenerator {
	// Validate inputs
	if width <= 0 {
		width = 32 // Default width
	}
	if height <= 0 {
		height = 32 // Default height
	}
	
	return &SpriteGenerator{
		Width:    width,
		Height:   height,
		Symmetry: symmetry,
		Constraints: SpriteConstraints{
			MinDensity:     0.3,
			MaxDensity:     0.7,
			RequireOutline: true,
			ColorCount:     6,
		},
	}
}

// Generate creates a sprite from a seed
func (sg *SpriteGenerator) Generate(seed int64) *Sprite {
	rng := rand.New(rand.NewSource(seed))
	
	// Create base shape using cellular automata
	grid := sg.generateBaseShape(rng)
	
	// Apply symmetry
	grid = sg.applySymmetry(grid)
	
	// Generate color palette
	palette := sg.generatePalette(rng, sg.Constraints.ColorCount)
	
	// Flood fill regions with colors
	colored := sg.floodFillColors(grid, palette, rng)
	
	// Add shading
	shaded := sg.applyShading(colored, palette)
	
	// Add outline if required
	if sg.Constraints.RequireOutline {
		shaded = sg.addOutline(shaded)
	}
	
	return shaded
}

// generateBaseShape creates initial pixel layout using cellular automata
func (sg *SpriteGenerator) generateBaseShape(rng *rand.Rand) [][]bool {
	grid := make([][]bool, sg.Height)
	for y := range grid {
		grid[y] = make([]bool, sg.Width)
		for x := range grid[y] {
			grid[y][x] = rng.Float64() < 0.5
		}
	}
	
	// Run cellular automata iterations
	for i := 0; i < 3; i++ {
		grid = sg.cellularAutomataStep(grid)
	}
	
	return grid
}

// cellularAutomataStep performs one CA iteration
func (sg *SpriteGenerator) cellularAutomataStep(grid [][]bool) [][]bool {
	newGrid := make([][]bool, len(grid))
	for y := range newGrid {
		newGrid[y] = make([]bool, len(grid[y]))
		for x := range newGrid[y] {
			neighbors := sg.countNeighbors(grid, x, y)
			if grid[y][x] {
				newGrid[y][x] = neighbors >= 4
			} else {
				newGrid[y][x] = neighbors >= 5
			}
		}
	}
	return newGrid
}

// countNeighbors counts active neighbors in grid
func (sg *SpriteGenerator) countNeighbors(grid [][]bool, x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < sg.Width && ny >= 0 && ny < sg.Height {
				if grid[ny][nx] {
					count++
				}
			}
		}
	}
	return count
}

// applySymmetry applies symmetry transformation to grid
func (sg *SpriteGenerator) applySymmetry(grid [][]bool) [][]bool {
	switch sg.Symmetry {
	case HorizontalSymmetry:
		return sg.applyHorizontalSymmetry(grid)
	case VerticalSymmetry:
		return sg.applyVerticalSymmetry(grid)
	case RadialSymmetry:
		return sg.applyRadialSymmetry(grid)
	default:
		return grid
	}
}

// applyHorizontalSymmetry mirrors grid horizontally
func (sg *SpriteGenerator) applyHorizontalSymmetry(grid [][]bool) [][]bool {
	for y := range grid {
		for x := 0; x < sg.Width/2; x++ {
			grid[y][sg.Width-1-x] = grid[y][x]
		}
	}
	return grid
}

// applyVerticalSymmetry mirrors grid vertically
func (sg *SpriteGenerator) applyVerticalSymmetry(grid [][]bool) [][]bool {
	for y := 0; y < sg.Height/2; y++ {
		for x := range grid[y] {
			grid[sg.Height-1-y][x] = grid[y][x]
		}
	}
	return grid
}

// applyRadialSymmetry applies 4-way rotational symmetry
func (sg *SpriteGenerator) applyRadialSymmetry(grid [][]bool) [][]bool {
	// Simplified radial symmetry - copy quadrants
	midX, midY := sg.Width/2, sg.Height/2
	for y := 0; y < midY; y++ {
		for x := 0; x < midX; x++ {
			val := grid[y][x]
			grid[y][sg.Width-1-x] = val
			grid[sg.Height-1-y][x] = val
			grid[sg.Height-1-y][sg.Width-1-x] = val
		}
	}
	return grid
}

// generatePalette creates a color palette
func (sg *SpriteGenerator) generatePalette(rng *rand.Rand, count int) []color.RGBA {
	palette := make([]color.RGBA, count)
	
	// Generate base hue
	baseHue := rng.Float64() * 360.0
	
	for i := range palette {
		// Create complementary and analogous colors
		hue := baseHue + float64(i)*30.0
		for hue >= 360.0 {
			hue -= 360.0
		}
		
		saturation := 0.6 + rng.Float64()*0.3
		value := 0.4 + float64(i)*0.1
		
		palette[i] = hsvToRGB(hue, saturation, value)
	}
	
	return palette
}

// hsvToRGB converts HSV to RGB color
func hsvToRGB(h, s, v float64) color.RGBA {
	c := v * s
	x := c * (1.0 - abs((mod(h/60.0, 2.0) - 1.0)))
	m := v - c
	
	var r, g, b float64
	
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	
	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func mod(x, y float64) float64 {
	return x - y*float64(int(x/y))
}

// floodFillColors assigns colors to regions
func (sg *SpriteGenerator) floodFillColors(grid [][]bool, palette []color.RGBA, rng *rand.Rand) *Sprite {
	sprite := &Sprite{
		Image:  image.NewRGBA(image.Rect(0, 0, sg.Width, sg.Height)),
		Width:  sg.Width,
		Height: sg.Height,
	}
	
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] {
				colorIdx := rng.Intn(len(palette))
				sprite.Image.Set(x, y, palette[colorIdx])
			} else {
				sprite.Image.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	
	return sprite
}

// applyShading adds depth to sprite
func (sg *SpriteGenerator) applyShading(sprite *Sprite, palette []color.RGBA) *Sprite {
	// Simple shading - darken bottom pixels
	for y := sprite.Height / 2; y < sprite.Height; y++ {
		for x := 0; x < sprite.Width; x++ {
			c := sprite.Image.RGBAAt(x, y)
			if c.A > 0 {
				factor := 0.8
				c.R = uint8(float64(c.R) * factor)
				c.G = uint8(float64(c.G) * factor)
				c.B = uint8(float64(c.B) * factor)
				sprite.Image.Set(x, y, c)
			}
		}
	}
	return sprite
}

// addOutline adds outline to sprite for clarity
func (sg *SpriteGenerator) addOutline(sprite *Sprite) *Sprite {
	outlined := image.NewRGBA(sprite.Image.Bounds())
	
	// Copy original
	for y := 0; y < sprite.Height; y++ {
		for x := 0; x < sprite.Width; x++ {
			outlined.Set(x, y, sprite.Image.At(x, y))
		}
	}
	
	// Add outline where transparent meets opaque
	black := color.RGBA{0, 0, 0, 255}
	for y := 0; y < sprite.Height; y++ {
		for x := 0; x < sprite.Width; x++ {
			c := sprite.Image.RGBAAt(x, y)
			if c.A == 0 {
				// Check if adjacent to opaque pixel
				if sg.hasOpaqueNeighbor(sprite, x, y) {
					outlined.Set(x, y, black)
				}
			}
		}
	}
	
	sprite.Image = outlined
	return sprite
}

// hasOpaqueNeighbor checks if pixel has opaque neighbors
func (sg *SpriteGenerator) hasOpaqueNeighbor(sprite *Sprite, x, y int) bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < sprite.Width && ny >= 0 && ny < sprite.Height {
				if sprite.Image.RGBAAt(nx, ny).A > 0 {
					return true
				}
			}
		}
	}
	return false
}
