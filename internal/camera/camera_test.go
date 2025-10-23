package camera

import (
	"math"
	"testing"
)

func TestNewCamera(t *testing.T) {
	config := CameraConfig{
		ScreenWidth:        800,
		ScreenHeight:       600,
		FollowSmoothing:    0.1,
		ZoomSmoothing:      0.15,
		DeadZoneWidth:      64,
		DeadZoneHeight:     32,
		LookAheadDistance:  100,
		LookAheadSmoothing: 0.05,
	}

	camera := NewCamera(config)

	if camera == nil {
		t.Fatal("NewCamera returned nil")
	}

	if camera.screenWidth != 800 || camera.screenHeight != 600 {
		t.Errorf("Screen size not set correctly: got %dx%d, want 800x600", camera.screenWidth, camera.screenHeight)
	}

	if camera.Zoom != 1.0 {
		t.Errorf("Initial zoom not 1.0: got %f", camera.Zoom)
	}

	if camera.followSmoothing != 0.1 {
		t.Errorf("Follow smoothing not set correctly: got %f, want 0.1", camera.followSmoothing)
	}
}

func TestNewDefaultCamera(t *testing.T) {
	camera := NewDefaultCamera(1024, 768)

	if camera == nil {
		t.Fatal("NewDefaultCamera returned nil")
	}

	if camera.screenWidth != 1024 || camera.screenHeight != 768 {
		t.Errorf("Screen size not set correctly: got %dx%d, want 1024x768", camera.screenWidth, camera.screenHeight)
	}

	// Check that defaults are reasonable
	if camera.followSmoothing <= 0 || camera.followSmoothing >= 1 {
		t.Errorf("Invalid default follow smoothing: %f", camera.followSmoothing)
	}

	if camera.deadZoneWidth <= 0 || camera.deadZoneHeight <= 0 {
		t.Error("Dead zone should be positive")
	}
}

func TestSetPosition(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	camera.SetPosition(100, 200)

	if camera.X != 100 || camera.Y != 200 {
		t.Errorf("Position not set correctly: got (%f, %f), want (100, 200)", camera.X, camera.Y)
	}

	if camera.targetX != 100 || camera.targetY != 200 {
		t.Errorf("Target position not set correctly: got (%f, %f), want (100, 200)", camera.targetX, camera.targetY)
	}
}

func TestSetZoom(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	camera.SetZoom(2.0)

	if camera.Zoom != 2.0 {
		t.Errorf("Zoom not set correctly: got %f, want 2.0", camera.Zoom)
	}

	if camera.targetZoom != 2.0 {
		t.Errorf("Target zoom not set correctly: got %f, want 2.0", camera.targetZoom)
	}

	// Test invalid zoom
	camera.SetZoom(-1.0)
	if camera.Zoom <= 0 {
		t.Error("Zoom should be clamped to positive values")
	}
}

func TestSetBounds(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	camera.SetBounds(0, 0, 1000, 800)

	if !camera.boundsEnabled {
		t.Error("Bounds should be enabled after SetBounds")
	}

	if camera.minX != 0 || camera.minY != 0 || camera.maxX != 1000 || camera.maxY != 800 {
		t.Errorf("Bounds not set correctly: got (%f, %f, %f, %f), want (0, 0, 1000, 800)",
			camera.minX, camera.minY, camera.maxX, camera.maxY)
	}

	camera.DisableBounds()
	if camera.boundsEnabled {
		t.Error("Bounds should be disabled after DisableBounds")
	}
}

func TestStartShake(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	camera.StartShake(10.0, 0.5)

	if camera.shakeIntensity != 10.0 {
		t.Errorf("Shake intensity not set: got %f, want 10.0", camera.shakeIntensity)
	}

	if camera.shakeDuration != 0.5 {
		t.Errorf("Shake duration not set: got %f, want 0.5", camera.shakeDuration)
	}

	if camera.shakeTime != 0 {
		t.Errorf("Shake time should start at 0: got %f", camera.shakeTime)
	}
}

func TestUpdate(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	// Set target position
	camera.SetTarget(100, 100)

	// Update once
	camera.Update(0.016) // ~60 FPS

	// Camera should move towards target
	if camera.X == 0 && camera.Y == 0 {
		t.Error("Camera should have moved towards target")
	}

	// Should not reach target immediately due to smoothing
	if camera.X == 100 && camera.Y == 100 {
		t.Error("Camera should not reach target immediately with smoothing")
	}
}

func TestSmoothZoom(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	camera.SetTargetZoom(2.0)

	// Update several times
	for i := 0; i < 10; i++ {
		camera.Update(0.016)
	}

	// Zoom should be closer to target but not exactly there due to smoothing
	if camera.Zoom == 1.0 {
		t.Error("Zoom should have changed towards target")
	}

	if camera.Zoom == 2.0 {
		t.Error("Zoom should not reach target immediately with smoothing")
	}

	if camera.Zoom <= 1.0 || camera.Zoom >= 2.0 {
		t.Errorf("Zoom should be between 1.0 and 2.0: got %f", camera.Zoom)
	}
}

func TestShakeEffect(t *testing.T) {
	camera := NewDefaultCamera(800, 600)

	camera.StartShake(5.0, 0.1)

	// Update halfway through shake
	camera.Update(0.05)

	if camera.shakeOffsetX == 0 && camera.shakeOffsetY == 0 {
		t.Error("Shake should produce non-zero offsets")
	}

	// Update past shake duration
	camera.Update(0.1)

	if camera.shakeOffsetX != 0 || camera.shakeOffsetY != 0 {
		t.Error("Shake offsets should be zero after duration ends")
	}
}

func TestWorldToScreen(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)
	camera.SetZoom(1.0)

	// World origin should map to screen center
	screenX, screenY := camera.WorldToScreen(0, 0)
	expectedX := 400.0 // 800 / 2
	expectedY := 300.0 // 600 / 2

	if screenX != expectedX || screenY != expectedY {
		t.Errorf("World origin not at screen center: got (%f, %f), want (%f, %f)",
			screenX, screenY, expectedX, expectedY)
	}

	// Test with camera offset
	camera.SetPosition(100, 50)
	screenX, screenY = camera.WorldToScreen(100, 50)

	if screenX != expectedX || screenY != expectedY {
		t.Errorf("Camera position not centered: got (%f, %f), want (%f, %f)",
			screenX, screenY, expectedX, expectedY)
	}
}

func TestScreenToWorld(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)
	camera.SetZoom(1.0)

	// Screen center should map to camera position
	worldX, worldY := camera.ScreenToWorld(400, 300)

	if worldX != 0 || worldY != 0 {
		t.Errorf("Screen center not at world origin: got (%f, %f), want (0, 0)", worldX, worldY)
	}

	// Test roundtrip conversion
	originalWorldX, originalWorldY := 123.0, 456.0
	screenX, screenY := camera.WorldToScreen(originalWorldX, originalWorldY)
	convertedWorldX, convertedWorldY := camera.ScreenToWorld(screenX, screenY)

	tolerance := 0.001
	if math.Abs(convertedWorldX-originalWorldX) > tolerance ||
		math.Abs(convertedWorldY-originalWorldY) > tolerance {
		t.Errorf("Roundtrip conversion failed: original (%f, %f), converted (%f, %f)",
			originalWorldX, originalWorldY, convertedWorldX, convertedWorldY)
	}
}

func TestZoomEffect(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)

	// Test 2x zoom
	camera.SetZoom(2.0)

	// World point (50, 50) should appear further from center when zoomed
	screenX1, screenY1 := camera.WorldToScreen(50, 50)

	camera.SetZoom(1.0)
	screenX2, screenY2 := camera.WorldToScreen(50, 50)

	// With 2x zoom, the screen distance should be doubled
	expectedX1 := 400.0 + 50.0*2.0 // center + world_offset * zoom
	expectedY1 := 300.0 + 50.0*2.0
	expectedX2 := 400.0 + 50.0*1.0 // center + world_offset * zoom
	expectedY2 := 300.0 + 50.0*1.0

	if screenX1 != expectedX1 || screenY1 != expectedY1 {
		t.Errorf("2x zoom incorrect: got (%f, %f), want (%f, %f)",
			screenX1, screenY1, expectedX1, expectedY1)
	}

	if screenX2 != expectedX2 || screenY2 != expectedY2 {
		t.Errorf("1x zoom incorrect: got (%f, %f), want (%f, %f)",
			screenX2, screenY2, expectedX2, expectedY2)
	}
}

func TestIsPointVisible(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)
	camera.SetZoom(1.0)

	// Point at camera center should be visible
	if !camera.IsPointVisible(0, 0) {
		t.Error("Point at camera center should be visible")
	}

	// Point way off screen should not be visible
	if camera.IsPointVisible(10000, 10000) {
		t.Error("Point far off screen should not be visible")
	}

	// Point near edge should be visible
	if !camera.IsPointVisible(350, 250) {
		t.Error("Point near edge should be visible")
	}
}

func TestGetVisibleBounds(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)
	camera.SetZoom(1.0)

	minX, minY, maxX, maxY := camera.GetVisibleBounds()

	// At 1x zoom, visible area should be screen size
	expectedMinX := -400.0 // -screenWidth/2
	expectedMinY := -300.0 // -screenHeight/2
	expectedMaxX := 400.0  // screenWidth/2
	expectedMaxY := 300.0  // screenHeight/2

	if minX != expectedMinX || minY != expectedMinY ||
		maxX != expectedMaxX || maxY != expectedMaxY {
		t.Errorf("Visible bounds incorrect: got (%f, %f, %f, %f), want (%f, %f, %f, %f)",
			minX, minY, maxX, maxY, expectedMinX, expectedMinY, expectedMaxX, expectedMaxY)
	}

	// Test with 2x zoom (should see half as much)
	camera.SetZoom(2.0)
	minX, minY, maxX, maxY = camera.GetVisibleBounds()

	expectedMinX = -200.0
	expectedMinY = -150.0
	expectedMaxX = 200.0
	expectedMaxY = 150.0

	if minX != expectedMinX || minY != expectedMinY ||
		maxX != expectedMaxX || maxY != expectedMaxY {
		t.Errorf("2x zoom bounds incorrect: got (%f, %f, %f, %f), want (%f, %f, %f, %f)",
			minX, minY, maxX, maxY, expectedMinX, expectedMinY, expectedMaxX, expectedMaxY)
	}
}

func TestFollowTarget(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)

	// Target moving right
	camera.FollowTarget(100, 0, 50, 0, 0.016)
	camera.Update(0.016)

	// Camera should move towards target
	if camera.X <= 0 {
		t.Error("Camera should move towards target")
	}

	// Look ahead should be applied (camera should be ahead of target position)
	if camera.targetX <= 100 {
		t.Error("Camera should have look ahead when target is moving")
	}
}

func TestZoomAt(t *testing.T) {
	camera := NewDefaultCamera(800, 600)
	camera.SetPosition(0, 0)
	camera.SetZoom(1.0)

	// Zoom at world point (100, 100)
	worldX, worldY := 100.0, 100.0
	screenX, screenY := camera.WorldToScreen(worldX, worldY)

	camera.ZoomAt(worldX, worldY, 2.0)

	// The same world point should still be at the same screen position
	newScreenX, newScreenY := camera.WorldToScreen(worldX, worldY)

	tolerance := 1.0
	if math.Abs(newScreenX-screenX) > tolerance || math.Abs(newScreenY-screenY) > tolerance {
		t.Errorf("ZoomAt didn't maintain screen position: before (%f, %f), after (%f, %f)",
			screenX, screenY, newScreenX, newScreenY)
	}

	if camera.Zoom != 2.0 {
		t.Errorf("Zoom not set correctly: got %f, want 2.0", camera.Zoom)
	}
}

func TestBoundsConstraint(t *testing.T) {
	camera := NewDefaultCamera(400, 300) // Smaller screen for easier testing
	camera.SetBounds(0, 0, 800, 600)
	camera.SetZoom(1.0)

	// Try to move camera outside bounds
	camera.SetPosition(-1000, -1000)
	camera.Update(0.016)

	// Camera should be constrained within bounds
	minX, minY, maxX, maxY := camera.GetVisibleBounds()

	if minX < 0 || minY < 0 {
		t.Errorf("Camera moved outside minimum bounds: visible area (%f, %f, %f, %f)",
			minX, minY, maxX, maxY)
	}

	// Try to move camera to other extreme
	camera.SetPosition(2000, 2000)
	camera.Update(0.016)

	minX, minY, maxX, maxY = camera.GetVisibleBounds()

	if maxX > 800 || maxY > 600 {
		t.Errorf("Camera moved outside maximum bounds: visible area (%f, %f, %f, %f)",
			minX, minY, maxX, maxY)
	}
}

func TestLerp(t *testing.T) {
	result := Lerp(0, 100, 0.5)
	if result != 50 {
		t.Errorf("Lerp(0, 100, 0.5) = %f, want 50", result)
	}

	result = Lerp(10, 20, 0.0)
	if result != 10 {
		t.Errorf("Lerp(10, 20, 0.0) = %f, want 10", result)
	}

	result = Lerp(10, 20, 1.0)
	if result != 20 {
		t.Errorf("Lerp(10, 20, 1.0) = %f, want 20", result)
	}
}

func TestSmoothStep(t *testing.T) {
	result := SmoothStep(0.0)
	if result != 0.0 {
		t.Errorf("SmoothStep(0.0) = %f, want 0.0", result)
	}

	result = SmoothStep(1.0)
	if result != 1.0 {
		t.Errorf("SmoothStep(1.0) = %f, want 1.0", result)
	}

	result = SmoothStep(0.5)
	if result != 0.5 {
		t.Errorf("SmoothStep(0.5) = %f, want 0.5", result)
	}

	// Test that SmoothStep provides smooth interpolation
	// (derivative should be 0 at both ends)
	result1 := SmoothStep(0.1)
	result2 := SmoothStep(0.9)

	if result1 >= 0.1 || result2 <= 0.9 {
		t.Error("SmoothStep should provide ease in/out behavior")
	}
}
