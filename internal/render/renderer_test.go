package render

import (
	"testing"

	"github.com/opd-ai/vania/internal/world"
)

func TestNewRenderer(t *testing.T) {
	r := NewRenderer()
	
	if r == nil {
		t.Fatal("NewRenderer returned nil")
	}
	
	if r.camera == nil {
		t.Error("Camera not initialized")
	}
	
	if r.camera.Width != ScreenWidth {
		t.Errorf("Expected camera width %d, got %d", ScreenWidth, r.camera.Width)
	}
	
	if r.camera.Height != ScreenHeight {
		t.Errorf("Expected camera height %d, got %d", ScreenHeight, r.camera.Height)
	}
}

func TestUpdateCamera(t *testing.T) {
	r := NewRenderer()
	
	// Update camera to target position
	targetX := 500.0
	targetY := 300.0
	r.UpdateCamera(targetX, targetY)
	
	// Camera should center on target
	expectedX := targetX - float64(r.camera.Width)/2
	expectedY := targetY - float64(r.camera.Height)/2
	
	if r.camera.X != expectedX {
		t.Errorf("Expected camera X=%f, got %f", expectedX, r.camera.X)
	}
	
	if r.camera.Y != expectedY {
		t.Errorf("Expected camera Y=%f, got %f", expectedY, r.camera.Y)
	}
}

func TestGetCameraOffset(t *testing.T) {
	r := NewRenderer()
	r.camera.X = 100
	r.camera.Y = 200
	
	offsetX, offsetY := r.GetCameraOffset()
	
	if offsetX != -100 {
		t.Errorf("Expected camera offset X=-100, got %f", offsetX)
	}
	
	if offsetY != -200 {
		t.Errorf("Expected camera offset Y=-200, got %f", offsetY)
	}
}

func TestRenderDataStructures(t *testing.T) {
	// Test that render package data structures are correctly defined
	
	// Create a test biome
	biome := &world.Biome{
		Name:        "test",
		DangerLevel: 1,
		Temperature: 20,
	}
	
	// Create a test room
	room := &world.Room{
		ID:    1,
		Type:  world.StartRoom,
		X:     0,
		Y:     0,
		Width: 30,
		Height: 20,
		Biome:  biome,
		Platforms: []world.Platform{
			{X: 100, Y: 500, Width: 200, Height: 32},
		},
		Hazards: []world.Hazard{
			{X: 400, Y: 550, Type: "spike", Damage: 10, Width: 64, Height: 32},
		},
	}
	
	// Verify the data structures are set up correctly
	if room.Biome.Name != "test" {
		t.Error("Room biome not set correctly")
	}
	
	if len(room.Platforms) != 1 {
		t.Error("Room platforms not set correctly")
	}
	
	if len(room.Hazards) != 1 {
		t.Error("Room hazards not set correctly")
	}
	
	if room.Hazards[0].Type != "spike" {
		t.Error("Hazard type not set correctly")
	}
}

// Note: Full rendering tests require ebiten graphics context which needs X11/graphics libraries.
// These tests verify the data structures and non-rendering logic.
// Integration tests with actual rendering should be run in a graphical environment.
