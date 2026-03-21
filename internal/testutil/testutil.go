// Package testutil provides helpers for testing VANIA packages that depend
// on Ebiten rendering infrastructure. These utilities allow tests to skip
// display-dependent scenarios when running in a headless CI environment.
//
// Tests requiring an X11 display should use SkipIfNoDisplay to avoid panics.
// The verify.sh script uses xvfb-run to provide a virtual display for the
// full test suite; this package enables graceful degradation when xvfb is
// unavailable.
package testutil

import (
	"image"
	"image/color"
	"os"
	"testing"
)

// SkipIfNoDisplay calls t.Skip when the DISPLAY environment variable is unset
// or empty. Call this at the start of any test that instantiates Ebiten
// windows, renderers, or other display-dependent objects.
func SkipIfNoDisplay(t *testing.T) {
	t.Helper()
	if os.Getenv("DISPLAY") == "" {
		t.Skip("DISPLAY not set: skipping display-dependent test (run with xvfb-run -a go test ./...)")
	}
}

// MockImage is a minimal image.RGBA-backed stand-in for *ebiten.Image in unit
// tests that only need to verify pixel data without a real graphics context.
type MockImage struct {
	*image.RGBA
}

// NewMockImage creates a blank MockImage of the given dimensions.
func NewMockImage(width, height int) *MockImage {
	return &MockImage{RGBA: image.NewRGBA(image.Rect(0, 0, width, height))}
}

// SetPixel sets a single pixel to the given colour.
func (m *MockImage) SetPixel(x, y int, c color.RGBA) {
	m.Set(x, y, c)
}

// PixelAt returns the colour of a single pixel.
func (m *MockImage) PixelAt(x, y int) color.RGBA {
	r, g, b, a := m.RGBAAt(x, y).RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
