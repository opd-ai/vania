// Package camera provides smooth camera controls with follow mechanics,package camera

// zoom functionality, and screen space conversions for the VANIA game engine.
package camera

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera represents a 2D camera with position, zoom, and target following
type Camera struct {
	// Position (world coordinates that the camera is looking at)
	X, Y float64

	// Target position for smooth following
	targetX, targetY float64

	// Zoom level (1.0 = normal, 2.0 = zoomed in 2x, 0.5 = zoomed out 2x)
	Zoom       float64
	targetZoom float64

	// Follow settings
	followSmoothing float64 // 0.0 = instant, 1.0 = never catches up
	zoomSmoothing   float64 // Smoothing for zoom transitions

	// Screen/viewport dimensions
	screenWidth  int
	screenHeight int

	// Camera bounds (world coordinates)
	minX, minY, maxX, maxY float64
	boundsEnabled          bool

	// Shake effect
	shakeIntensity float64
	shakeDuration  float64
	shakeTime      float64
	shakeOffsetX   float64
	shakeOffsetY   float64

	// Dead zone (area where target can move without camera moving)
	deadZoneWidth  float64
	deadZoneHeight float64

	// Look ahead (camera leads target based on movement direction)
	lookAheadDistance  float64
	lookAheadX         float64
	lookAheadY         float64
	lookAheadSmoothing float64
}

// CameraConfig holds configuration for camera creation
type CameraConfig struct {
	ScreenWidth        int     // Screen/viewport width
	ScreenHeight       int     // Screen/viewport height
	FollowSmoothing    float64 // Camera follow smoothing (0.0 - 1.0)
	ZoomSmoothing      float64 // Zoom transition smoothing (0.0 - 1.0)
	DeadZoneWidth      float64 // Dead zone width in world units
	DeadZoneHeight     float64 // Dead zone height in world units
	LookAheadDistance  float64 // Look ahead distance in world units
	LookAheadSmoothing float64 // Look ahead smoothing (0.0 - 1.0)
}

// NewCamera creates a new camera with the given configuration
func NewCamera(config CameraConfig) *Camera {
	return &Camera{
		X:                  0,
		Y:                  0,
		targetX:            0,
		targetY:            0,
		Zoom:               1.0,
		targetZoom:         1.0,
		followSmoothing:    config.FollowSmoothing,
		zoomSmoothing:      config.ZoomSmoothing,
		screenWidth:        config.ScreenWidth,
		screenHeight:       config.ScreenHeight,
		boundsEnabled:      false,
		shakeIntensity:     0,
		shakeDuration:      0,
		shakeTime:          0,
		deadZoneWidth:      config.DeadZoneWidth,
		deadZoneHeight:     config.DeadZoneHeight,
		lookAheadDistance:  config.LookAheadDistance,
		lookAheadSmoothing: config.LookAheadSmoothing,
	}
}

// NewDefaultCamera creates a camera with sensible defaults
func NewDefaultCamera(screenWidth, screenHeight int) *Camera {
	return NewCamera(CameraConfig{
		ScreenWidth:        screenWidth,
		ScreenHeight:       screenHeight,
		FollowSmoothing:    0.1,
		ZoomSmoothing:      0.15,
		DeadZoneWidth:      64,
		DeadZoneHeight:     32,
		LookAheadDistance:  100,
		LookAheadSmoothing: 0.05,
	})
}

// SetPosition immediately sets camera position (no smoothing)
func (c *Camera) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
	c.targetX = x
	c.targetY = y
}

// SetTarget sets the target position for smooth following
func (c *Camera) SetTarget(x, y float64) {
	c.targetX = x
	c.targetY = y
}

// SetZoom immediately sets zoom level
func (c *Camera) SetZoom(zoom float64) {
	if zoom <= 0 {
		zoom = 0.1
	}
	c.Zoom = zoom
	c.targetZoom = zoom
}

// SetTargetZoom sets target zoom for smooth zooming
func (c *Camera) SetTargetZoom(zoom float64) {
	if zoom <= 0 {
		zoom = 0.1
	}
	c.targetZoom = zoom
}

// SetBounds sets world bounds for the camera
func (c *Camera) SetBounds(minX, minY, maxX, maxY float64) {
	c.minX = minX
	c.minY = minY
	c.maxX = maxX
	c.maxY = maxY
	c.boundsEnabled = true
}

// DisableBounds disables camera bounds checking
func (c *Camera) DisableBounds() {
	c.boundsEnabled = false
}

// StartShake starts a screen shake effect
func (c *Camera) StartShake(intensity, duration float64) {
	c.shakeIntensity = intensity
	c.shakeDuration = duration
	c.shakeTime = 0
}

// FollowTarget smoothly moves camera towards target with dead zone and look ahead
func (c *Camera) FollowTarget(targetX, targetY, velocityX, velocityY float64, deltaTime float64) {
	// Update look ahead based on target velocity
	targetLookAheadX := velocityX * c.lookAheadDistance
	targetLookAheadY := velocityY * c.lookAheadDistance * 0.5 // Less vertical look ahead

	// Smooth look ahead
	c.lookAheadX += (targetLookAheadX - c.lookAheadX) * c.lookAheadSmoothing
	c.lookAheadY += (targetLookAheadY - c.lookAheadY) * c.lookAheadSmoothing

	// Calculate effective target with look ahead
	effectiveTargetX := targetX + c.lookAheadX
	effectiveTargetY := targetY + c.lookAheadY

	// Dead zone calculation
	deadZoneLeft := c.X - c.deadZoneWidth/2
	deadZoneRight := c.X + c.deadZoneWidth/2
	deadZoneTop := c.Y - c.deadZoneHeight/2
	deadZoneBottom := c.Y + c.deadZoneHeight/2

	// Only move camera if target is outside dead zone
	if effectiveTargetX < deadZoneLeft {
		c.targetX = effectiveTargetX + c.deadZoneWidth/2
	} else if effectiveTargetX > deadZoneRight {
		c.targetX = effectiveTargetX - c.deadZoneWidth/2
	}

	if effectiveTargetY < deadZoneTop {
		c.targetY = effectiveTargetY + c.deadZoneHeight/2
	} else if effectiveTargetY > deadZoneBottom {
		c.targetY = effectiveTargetY - c.deadZoneHeight/2
	}
}

// Update updates camera position, zoom, and effects
func (c *Camera) Update(deltaTime float64) {
	// Smooth follow to target
	c.X += (c.targetX - c.X) * c.followSmoothing
	c.Y += (c.targetY - c.Y) * c.followSmoothing

	// Smooth zoom to target
	c.Zoom += (c.targetZoom - c.Zoom) * c.zoomSmoothing

	// Apply bounds if enabled
	if c.boundsEnabled {
		c.applyBounds()
	}

	// Update shake effect
	if c.shakeTime < c.shakeDuration {
		c.shakeTime += deltaTime
		progress := c.shakeTime / c.shakeDuration

		// Shake intensity decreases over time
		currentIntensity := c.shakeIntensity * (1.0 - progress)

		// Random shake offsets
		angle := math.Sin(c.shakeTime*50) * math.Pi * 2
		c.shakeOffsetX = math.Cos(angle) * currentIntensity
		c.shakeOffsetY = math.Sin(angle) * currentIntensity

		if c.shakeTime >= c.shakeDuration {
			c.shakeOffsetX = 0
			c.shakeOffsetY = 0
		}
	}
}

// applyBounds ensures camera doesn't go outside world bounds
func (c *Camera) applyBounds() {
	halfScreenWidth := float64(c.screenWidth) / (2.0 * c.Zoom)
	halfScreenHeight := float64(c.screenHeight) / (2.0 * c.Zoom)

	// Clamp X position
	if c.X-halfScreenWidth < c.minX {
		c.X = c.minX + halfScreenWidth
		c.targetX = c.X
	}
	if c.X+halfScreenWidth > c.maxX {
		c.X = c.maxX - halfScreenWidth
		c.targetX = c.X
	}

	// Clamp Y position
	if c.Y-halfScreenHeight < c.minY {
		c.Y = c.minY + halfScreenHeight
		c.targetY = c.Y
	}
	if c.Y+halfScreenHeight > c.maxY {
		c.Y = c.maxY - halfScreenHeight
		c.targetY = c.Y
	}
}

// GetMatrix returns the camera transformation matrix for rendering
func (c *Camera) GetMatrix() ebiten.GeoM {
	var m ebiten.GeoM

	// Apply zoom around camera center
	m.Scale(c.Zoom, c.Zoom)

	// Translate to center camera on screen
	m.Translate(-c.X*c.Zoom, -c.Y*c.Zoom)
	m.Translate(float64(c.screenWidth)/2, float64(c.screenHeight)/2)

	// Apply shake offset (after zoom and translation)
	m.Translate(c.shakeOffsetX*c.Zoom, c.shakeOffsetY*c.Zoom)

	return m
}

// GetInverseMatrix returns the inverse transformation for converting screen to world coords
func (c *Camera) GetInverseMatrix() ebiten.GeoM {
	var m ebiten.GeoM

	// Reverse the transformations
	m.Translate(-c.shakeOffsetX*c.Zoom, -c.shakeOffsetY*c.Zoom)
	m.Translate(-float64(c.screenWidth)/2, -float64(c.screenHeight)/2)
	m.Translate(c.X*c.Zoom, c.Y*c.Zoom)
	m.Scale(1.0/c.Zoom, 1.0/c.Zoom)

	return m
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY float64) (screenX, screenY float64) {
	// Apply camera transformation
	x := (worldX - c.X) * c.Zoom
	y := (worldY - c.Y) * c.Zoom

	// Translate to screen center
	screenX = x + float64(c.screenWidth)/2
	screenY = y + float64(c.screenHeight)/2

	// Apply shake offset
	screenX += c.shakeOffsetX * c.Zoom
	screenY += c.shakeOffsetY * c.Zoom

	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY float64) (worldX, worldY float64) {
	// Remove shake offset
	screenX -= c.shakeOffsetX * c.Zoom
	screenY -= c.shakeOffsetY * c.Zoom

	// Translate from screen center
	x := screenX - float64(c.screenWidth)/2
	y := screenY - float64(c.screenHeight)/2

	// Apply inverse zoom and camera position
	worldX = x/c.Zoom + c.X
	worldY = y/c.Zoom + c.Y

	return worldX, worldY
}

// IsPointVisible checks if a world point is visible on screen
func (c *Camera) IsPointVisible(worldX, worldY float64) bool {
	screenX, screenY := c.WorldToScreen(worldX, worldY)
	return screenX >= 0 && screenX < float64(c.screenWidth) &&
		screenY >= 0 && screenY < float64(c.screenHeight)
}

// IsRectVisible checks if a world rectangle is visible on screen
func (c *Camera) IsRectVisible(worldX, worldY, width, height float64) bool {
	// Check if any corner of the rectangle is visible
	corners := []struct{ x, y float64 }{
		{worldX, worldY},
		{worldX + width, worldY},
		{worldX, worldY + height},
		{worldX + width, worldY + height},
	}

	for _, corner := range corners {
		if c.IsPointVisible(corner.x, corner.y) {
			return true
		}
	}

	// Also check if camera is inside the rectangle
	if worldX <= c.X && c.X <= worldX+width &&
		worldY <= c.Y && c.Y <= worldY+height {
		return true
	}

	return false
}

// GetVisibleBounds returns the world coordinates of the visible area
func (c *Camera) GetVisibleBounds() (minX, minY, maxX, maxY float64) {
	halfScreenWidth := float64(c.screenWidth) / (2.0 * c.Zoom)
	halfScreenHeight := float64(c.screenHeight) / (2.0 * c.Zoom)

	minX = c.X - halfScreenWidth
	minY = c.Y - halfScreenHeight
	maxX = c.X + halfScreenWidth
	maxY = c.Y + halfScreenHeight

	return minX, minY, maxX, maxY
}

// SetScreenSize updates the camera's screen dimensions
func (c *Camera) SetScreenSize(width, height int) {
	c.screenWidth = width
	c.screenHeight = height
}

// GetPosition returns the current camera position
func (c *Camera) GetPosition() (x, y float64) {
	return c.X, c.Y
}

// GetZoom returns the current zoom level
func (c *Camera) GetZoom() float64 {
	return c.Zoom
}

// ZoomAt zooms in/out while keeping a specific world point at the same screen position
func (c *Camera) ZoomAt(worldX, worldY, newZoom float64) {
	if newZoom <= 0 {
		newZoom = 0.1
	}

	// Calculate screen position of the world point before zoom
	screenX, screenY := c.WorldToScreen(worldX, worldY)

	// Set new zoom
	c.Zoom = newZoom
	c.targetZoom = newZoom

	// Calculate what world position would give us the same screen position
	newWorldX, newWorldY := c.ScreenToWorld(screenX, screenY)

	// Adjust camera position to maintain the screen position
	c.X += worldX - newWorldX
	c.Y += worldY - newWorldY
	c.targetX = c.X
	c.targetY = c.Y
}

// Lerp linearly interpolates between two values
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// SmoothStep provides smooth interpolation with ease in/out
func SmoothStep(t float64) float64 {
	return t * t * (3.0 - 2.0*t)
}
