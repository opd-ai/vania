// Package render provides text rendering abstractions for consistent text display
// throughout the VANIA game, with fallback strategies for different contexts.
package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TextRenderer interface provides abstracted text rendering capabilities
type TextRenderer interface {
	// DrawText renders text at the specified position with the given color
	DrawText(screen *ebiten.Image, text string, x, y int, color color.Color)

	// MeasureText returns the width and height of the text in pixels
	MeasureText(text string) (width, height int)

	// GetLineHeight returns the standard line height for this renderer
	GetLineHeight() int
}

// DebugTextRenderer implements TextRenderer using ebitenutil debug functions
type DebugTextRenderer struct {
	charWidth  int
	charHeight int
}

// NewDebugTextRenderer creates a debug text renderer with standardized font metrics
func NewDebugTextRenderer() *DebugTextRenderer {
	return &DebugTextRenderer{
		charWidth:  8,  // Standardized to match bitmap renderer
		charHeight: 12, // Standardized to match bitmap renderer
	}
}

// DrawText renders text using ebitenutil.DebugPrintAt (color is ignored in debug mode)
func (dtr *DebugTextRenderer) DrawText(screen *ebiten.Image, text string, x, y int, col color.Color) {
	// Note: ebitenutil.DebugPrintAt doesn't support custom colors
	// This is a limitation of the debug text system
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

// MeasureText calculates text dimensions based on character metrics
func (dtr *DebugTextRenderer) MeasureText(text string) (width, height int) {
	lines := 1
	maxLineWidth := 0
	currentLineWidth := 0

	for _, char := range text {
		if char == '\n' {
			lines++
			if currentLineWidth > maxLineWidth {
				maxLineWidth = currentLineWidth
			}
			currentLineWidth = 0
		} else {
			currentLineWidth++
		}
	}

	if currentLineWidth > maxLineWidth {
		maxLineWidth = currentLineWidth
	}

	return maxLineWidth * dtr.charWidth, lines * dtr.charHeight
}

// GetLineHeight returns the line height for debug text
func (dtr *DebugTextRenderer) GetLineHeight() int {
	return dtr.charHeight
}

// BitmapTextRenderer implements TextRenderer using procedural bitmap font
type BitmapTextRenderer struct {
	charWidth  int
	charHeight int
}

// NewBitmapTextRenderer creates a bitmap text renderer
func NewBitmapTextRenderer() *BitmapTextRenderer {
	return &BitmapTextRenderer{
		charWidth:  8,  // Width of our bitmap font characters
		charHeight: 12, // Height of our bitmap font characters
	}
}

// DrawText renders text using procedural bitmap font with full color support
func (btr *BitmapTextRenderer) DrawText(screen *ebiten.Image, text string, x, y int, col color.Color) {
	currentX := x
	currentY := y

	for _, char := range text {
		if char == '\n' {
			currentX = x
			currentY += btr.charHeight
			continue
		}

		btr.drawChar(screen, char, currentX, currentY, btr.charWidth, btr.charHeight, col)
		currentX += btr.charWidth
	}
}

// MeasureText calculates text dimensions for bitmap font
func (btr *BitmapTextRenderer) MeasureText(text string) (width, height int) {
	lines := 1
	maxLineWidth := 0
	currentLineWidth := 0

	for _, char := range text {
		if char == '\n' {
			lines++
			if currentLineWidth > maxLineWidth {
				maxLineWidth = currentLineWidth
			}
			currentLineWidth = 0
		} else {
			currentLineWidth++
		}
	}

	if currentLineWidth > maxLineWidth {
		maxLineWidth = currentLineWidth
	}

	return maxLineWidth * btr.charWidth, lines * btr.charHeight
}

// GetLineHeight returns the line height for bitmap text
func (btr *BitmapTextRenderer) GetLineHeight() int {
	return btr.charHeight
}

// drawChar draws a single character using optimized batch pixel rendering
func (btr *BitmapTextRenderer) drawChar(screen *ebiten.Image, char rune, x, y, width, height int, col color.Color) {
	// Create character image and set all pixels at once for performance
	charImg := ebiten.NewImage(width, height)
	pixels := make([]byte, width*height*4) // RGBA format
	
	pattern := btr.getCharPattern(char, width, height)
	r, g, b, a := col.RGBA()
	
	// Convert to 8-bit values
	r8, g8, b8, a8 := byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8)
	
	// Batch set all pixels for the character
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			offset := (py*width + px) * 4
			if py < len(pattern) && px < len(pattern[py]) && pattern[py][px] {
				// Set pixel to character color
				pixels[offset] = r8     // Red
				pixels[offset+1] = g8   // Green
				pixels[offset+2] = b8   // Blue
				pixels[offset+3] = a8   // Alpha
			} else {
				// Set pixel to transparent
				pixels[offset] = 0
				pixels[offset+1] = 0
				pixels[offset+2] = 0
				pixels[offset+3] = 0
			}
		}
	}
	
	// Apply all pixels in a single operation
	charImg.WritePixels(pixels)
	
	// Single draw call per character
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(charImg, opts)
}

// getCharPattern returns a bitmap pattern for common characters
// This is a simplified version of the menu system's character patterns
func (btr *BitmapTextRenderer) getCharPattern(char rune, width, height int) [][]bool {
	// Initialize empty pattern
	pattern := make([][]bool, height)
	for i := range pattern {
		pattern[i] = make([]bool, width)
	}

	// Simple patterns for most common characters (subset for performance)
	switch char {
	case 'A', 'a':
		pattern[1][3] = true
		pattern[1][4] = true
		pattern[2][2] = true
		pattern[2][5] = true
		pattern[3][2] = true
		pattern[3][5] = true
		pattern[4][1] = true
		pattern[4][6] = true
		pattern[5][1] = true
		pattern[5][6] = true
		pattern[6][1] = true
		pattern[6][2] = true
		pattern[6][5] = true
		pattern[6][6] = true
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][1] = true
		pattern[8][6] = true

	case 'B', 'b':
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		pattern[1][2] = true
		pattern[1][3] = true
		pattern[1][4] = true
		pattern[1][5] = true
		pattern[2][6] = true
		pattern[3][6] = true
		pattern[4][6] = true
		pattern[5][2] = true
		pattern[5][3] = true
		pattern[5][4] = true
		pattern[5][5] = true
		pattern[6][6] = true
		pattern[7][6] = true
		pattern[8][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'C', 'c':
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[3][6] = true
		for y := 4; y < 7; y++ {
			pattern[y][1] = true
		}
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'D', 'd':
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		pattern[1][2] = true
		pattern[1][3] = true
		pattern[1][4] = true
		pattern[1][5] = true
		pattern[2][6] = true
		pattern[3][6] = true
		pattern[4][6] = true
		pattern[5][6] = true
		pattern[6][6] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'E', 'e':
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		for x := 1; x < 7; x++ {
			pattern[1][x] = true
			pattern[5][x] = true
			pattern[8][x] = true
		}

	case 'F', 'f':
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		for x := 1; x < 7; x++ {
			pattern[1][x] = true
			pattern[5][x] = true
		}

	case 'G', 'g':
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[4][1] = true
		pattern[5][1] = true
		pattern[6][1] = true
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true
		pattern[5][4] = true
		pattern[5][5] = true
		pattern[5][6] = true
		pattern[6][6] = true

	case 'H', 'h':
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		for x := 1; x < 7; x++ {
			pattern[5][x] = true
		}

	case 'I', 'i':
		for x := 2; x < 6; x++ {
			pattern[1][x] = true
			pattern[8][x] = true
		}
		for y := 2; y < 8; y++ {
			pattern[y][3] = true
		}

	case ' ':
		// Space - already empty

	case ':':
		pattern[3][3] = true
		pattern[6][3] = true

	case '-':
		pattern[5][2] = true
		pattern[5][3] = true
		pattern[5][4] = true
		pattern[5][5] = true

	case '.':
		pattern[8][3] = true

	case '!':
		for y := 1; y < 7; y++ {
			pattern[y][3] = true
		}
		pattern[8][3] = true

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		switch char {
		case '0':
			pattern[2][2] = true
			pattern[2][3] = true
			pattern[2][4] = true
			pattern[2][5] = true
			for y := 3; y < 8; y++ {
				pattern[y][1] = true
				pattern[y][6] = true
			}
			pattern[8][2] = true
			pattern[8][3] = true
			pattern[8][4] = true
			pattern[8][5] = true
		case '1':
			pattern[2][3] = true
			pattern[3][2] = true
			pattern[3][3] = true
			for y := 4; y < 9; y++ {
				pattern[y][3] = true
			}
			for x := 1; x < 7; x++ {
				pattern[8][x] = true
			}
		case '2':
			pattern[2][2] = true
			pattern[2][3] = true
			pattern[2][4] = true
			pattern[2][5] = true
			pattern[3][6] = true
			pattern[4][6] = true
			pattern[5][5] = true
			pattern[6][4] = true
			pattern[7][3] = true
			pattern[8][2] = true
			for x := 1; x < 7; x++ {
				pattern[8][x] = true
			}
		case '3':
			for x := 2; x < 6; x++ {
				pattern[2][x] = true
				pattern[5][x] = true
				pattern[8][x] = true
			}
			pattern[3][6] = true
			pattern[4][6] = true
			pattern[6][6] = true
			pattern[7][6] = true
		default:
			// Use simple box for other numbers
			for y := 2; y < 9; y++ {
				pattern[y][2] = true
				pattern[y][5] = true
			}
			for x := 2; x < 6; x++ {
				pattern[2][x] = true
				pattern[8][x] = true
			}
		}

	default:
		// Unknown character - draw a simple box
		for y := 2; y < 9; y++ {
			pattern[y][2] = true
			pattern[y][5] = true
		}
		for x := 2; x < 6; x++ {
			pattern[2][x] = true
			pattern[8][x] = true
		}
	}

	return pattern
}

// TextRenderManager manages text rendering strategy with fallback
type TextRenderManager struct {
	primary  TextRenderer
	fallback TextRenderer
	useColor bool
}

// NewTextRenderManager creates a text render manager with consistent font metrics
// Both renderers use standardized 8x12 character dimensions for layout consistency
func NewTextRenderManager(useColor bool) *TextRenderManager {
	primary := NewBitmapTextRenderer()
	fallback := NewDebugTextRenderer()

	return &TextRenderManager{
		primary:  primary,
		fallback: fallback,
		useColor: useColor,
	}
}

// DrawText renders text using the appropriate renderer
func (trm *TextRenderManager) DrawText(screen *ebiten.Image, text string, x, y int, color color.Color) {
	if trm.useColor {
		// Use bitmap renderer for color support
		trm.primary.DrawText(screen, text, x, y, color)
	} else {
		// Use fallback debug renderer
		trm.fallback.DrawText(screen, text, x, y, color)
	}
}

// MeasureText measures text using the current renderer
func (trm *TextRenderManager) MeasureText(text string) (width, height int) {
	if trm.useColor {
		return trm.primary.MeasureText(text)
	}
	return trm.fallback.MeasureText(text)
}

// GetLineHeight returns line height from current renderer
func (trm *TextRenderManager) GetLineHeight() int {
	if trm.useColor {
		return trm.primary.GetLineHeight()
	}
	return trm.fallback.GetLineHeight()
}

// SetColorMode enables or disables color rendering
func (trm *TextRenderManager) SetColorMode(enabled bool) {
	trm.useColor = enabled
}
