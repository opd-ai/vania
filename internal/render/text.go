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

	pattern := GetCharPattern(char, width, height)
	r, g, b, a := col.RGBA()

	// Convert to 8-bit values
	r8, g8, b8, a8 := byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8)

	// Batch set all pixels for the character
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			offset := (py*width + px) * 4
			if py < len(pattern) && px < len(pattern[py]) && pattern[py][px] {
				// Set pixel to character color
				pixels[offset] = r8   // Red
				pixels[offset+1] = g8 // Green
				pixels[offset+2] = b8 // Blue
				pixels[offset+3] = a8 // Alpha
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
