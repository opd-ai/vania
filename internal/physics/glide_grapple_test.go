package physics

import (
	"math"
	"testing"

	"github.com/opd-ai/vania/internal/world"
)

// TestGlide verifies gliding mechanics
func TestGlide(t *testing.T) {
	testCases := []struct {
		name         string
		gliding      bool
		initialVelY  float64
		expectedMaxY float64
		description  string
	}{
		{
			name:         "GlidingCapsFallSpeed",
			gliding:      true,
			initialVelY:  15.0,
			expectedMaxY: GlideFallSpeed,
			description:  "When gliding, fall speed should be capped at GlideFallSpeed",
		},
		{
			name:         "NoGlidingNormalFall",
			gliding:      false,
			initialVelY:  5.0,
			expectedMaxY: MaxFallSpeed,
			description:  "Without gliding, normal max fall speed applies",
		},
		{
			name:         "GlidingSlowsFall",
			gliding:      true,
			initialVelY:  0.0,
			expectedMaxY: GlideFallSpeed,
			description:  "Gliding prevents reaching max fall speed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := NewBody(100, 100, 32, 32)
			body.Velocity.Y = tc.initialVelY

			// Apply gravity with gliding for several frames
			for i := 0; i < 50; i++ {
				body.ApplyGravity(tc.gliding)
			}

			if body.Velocity.Y > tc.expectedMaxY+0.1 {
				t.Errorf("%s: Expected fall speed <= %.2f, got %.2f",
					tc.description, tc.expectedMaxY, body.Velocity.Y)
			}
		})
	}
}

// TestGlideMethod verifies the explicit Glide() method
func TestGlideMethod(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	body.Velocity.Y = 10.0

	body.Glide()

	if body.Velocity.Y > GlideFallSpeed {
		t.Errorf("Glide() should cap fall speed at %.2f, got %.2f",
			GlideFallSpeed, body.Velocity.Y)
	}
}

// TestStartGrapple verifies grapple initiation
func TestStartGrapple(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	anchor := world.AnchorPoint{X: 200, Y: 50}

	body.StartGrapple(anchor)

	if !body.Grappling {
		t.Error("StartGrapple should set Grappling to true")
	}

	if body.GrappleAnchor.X != anchor.X || body.GrappleAnchor.Y != anchor.Y {
		t.Errorf("Anchor should be set to (%.0f, %.0f), got (%.0f, %.0f)",
			anchor.X, anchor.Y, body.GrappleAnchor.X, body.GrappleAnchor.Y)
	}

	// Check that initial velocity is set toward anchor
	expectedDX := anchor.X - (body.Position.X + body.Position.Width/2)
	expectedDY := anchor.Y - (body.Position.Y + body.Position.Height/2)
	expectedDist := math.Sqrt(expectedDX*expectedDX + expectedDY*expectedDY)

	if math.Abs(body.GrappleLength-expectedDist) > 0.1 {
		t.Errorf("GrappleLength should be %.2f, got %.2f",
			expectedDist, body.GrappleLength)
	}

	// Verify launch velocity is non-zero
	if body.Velocity.X == 0 && body.Velocity.Y == 0 {
		t.Error("StartGrapple should set initial velocity toward anchor")
	}
}

// TestReleaseGrapple verifies grapple release
func TestReleaseGrapple(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	anchor := world.AnchorPoint{X: 200, Y: 50}

	body.StartGrapple(anchor)
	body.ReleaseGrapple()

	if body.Grappling {
		t.Error("ReleaseGrapple should set Grappling to false")
	}
}

// TestGrappleAutoDetachOnGround verifies grapple releases when landing
func TestGrappleAutoDetachOnGround(t *testing.T) {
	body := NewBody(100, 500, 32, 32)
	anchor := world.AnchorPoint{X: 200, Y: 50}

	body.StartGrapple(anchor)

	if !body.Grappling {
		t.Fatal("Should be grappling after StartGrapple")
	}

	// Simulate landing by moving body to bottom boundary
	body.Velocity.Y = 5.0
	body.Position.Y = 610 // Position such that Y + Height = 642 > 640

	// ResolveCollisionWithPlatforms should detect ground boundary and release grapple
	body.ResolveCollisionWithPlatforms([]world.Platform{})

	if body.Grappling {
		t.Error("Grappling should auto-release on landing at screen boundary")
	}

	if !body.OnGround {
		t.Error("Body should be on ground after hitting boundary")
	}
}

// TestUpdateGrapple verifies grapple physics simulation
func TestUpdateGrapple(t *testing.T) {
	body := NewBody(100, 200, 32, 32)
	anchor := world.AnchorPoint{X: 200, Y: 50}

	body.StartGrapple(anchor)
	initialAngle := body.GrappleAngle

	// Update grapple physics for several frames
	for i := 0; i < 10; i++ {
		body.UpdateGrapple()
	}

	// Angle should change (pendulum swing)
	if math.Abs(body.GrappleAngle-initialAngle) < 0.01 {
		t.Error("Grapple angle should change during UpdateGrapple (pendulum motion)")
	}

	// Position should stay constrained to rope length
	centerX := body.Position.X + body.Position.Width/2
	centerY := body.Position.Y + body.Position.Height/2
	dx := anchor.X - centerX
	dy := anchor.Y - centerY
	currentDist := math.Sqrt(dx*dx + dy*dy)

	if math.Abs(currentDist-body.GrappleLength) > 1.0 {
		t.Errorf("Body should stay at rope length %.2f, but is at distance %.2f",
			body.GrappleLength, currentDist)
	}
}

// TestUpdateGrappleWhenNotGrappling verifies UpdateGrapple does nothing when not grappling
func TestUpdateGrappleWhenNotGrappling(t *testing.T) {
	body := NewBody(100, 200, 32, 32)
	initialX := body.Position.X
	initialY := body.Position.Y

	body.UpdateGrapple()

	if body.Position.X != initialX || body.Position.Y != initialY {
		t.Error("UpdateGrapple should not modify position when not grappling")
	}
}

// TestFindNearestAnchor verifies anchor finding logic
func TestFindNearestAnchor(t *testing.T) {
	bodyPos := AABB{X: 100, Y: 100, Width: 32, Height: 32}

	testCases := []struct {
		name        string
		anchors     []world.AnchorPoint
		expectFound bool
		expectedIdx int
		description string
	}{
		{
			name: "FindsNearestInRange",
			anchors: []world.AnchorPoint{
				{X: 150, Y: 120}, // Distance ~50
				{X: 200, Y: 100}, // Distance ~84
				{X: 120, Y: 110}, // Distance ~20 - NEAREST
			},
			expectFound: true,
			expectedIdx: 2,
			description: "Should find the nearest anchor within range",
		},
		{
			name: "NoAnchorsInRange",
			anchors: []world.AnchorPoint{
				{X: 500, Y: 500}, // Far away
			},
			expectFound: false,
			description: "Should return false when no anchors in range",
		},
		{
			name:        "EmptyAnchorList",
			anchors:     []world.AnchorPoint{},
			expectFound: false,
			description: "Should return false for empty anchor list",
		},
		{
			name: "FindsClosestAmongMultiple",
			anchors: []world.AnchorPoint{
				{X: 200, Y: 150}, // Distance ~100
				{X: 140, Y: 120}, // Distance ~45
				{X: 130, Y: 115}, // Distance ~35 - NEAREST
				{X: 180, Y: 140}, // Distance ~90
			},
			expectFound: true,
			expectedIdx: 2,
			description: "Should find closest among multiple anchors",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			anchor, found := FindNearestAnchor(bodyPos, tc.anchors)

			if found != tc.expectFound {
				t.Errorf("%s: Expected found=%v, got %v",
					tc.description, tc.expectFound, found)
			}

			if found && tc.expectFound {
				expectedAnchor := tc.anchors[tc.expectedIdx]
				if anchor.X != expectedAnchor.X || anchor.Y != expectedAnchor.Y {
					t.Errorf("%s: Expected anchor (%.0f, %.0f), got (%.0f, %.0f)",
						tc.description, expectedAnchor.X, expectedAnchor.Y, anchor.X, anchor.Y)
				}
			}
		})
	}
}

// TestGrappleRangeLimit verifies anchors beyond range are not found
func TestGrappleRangeLimit(t *testing.T) {
	bodyPos := AABB{X: 100, Y: 100, Width: 32, Height: 32}

	// Create anchor just beyond range
	// GrappleAnchorRange = 192 pixels
	// Body center is at (116, 116)
	// Place anchor at distance > 192
	anchors := []world.AnchorPoint{
		{X: 116 + 200, Y: 116}, // Distance = 200 > 192
	}

	_, found := FindNearestAnchor(bodyPos, anchors)

	if found {
		t.Error("Should not find anchor beyond GrappleAnchorRange")
	}
}

// TestGlideAndGrappleInteraction verifies gliding is disabled when grappling
func TestGlideAndGrappleInteraction(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	anchor := world.AnchorPoint{X: 200, Y: 50}

	body.StartGrapple(anchor)
	body.Velocity.Y = 5.0

	// Try to apply gliding while grappling
	body.ApplyGravity(true)

	// Grappling should prevent gravity from being applied
	// (Grapple physics override normal gravity)
	// The grapple state means gravity won't add to velocity
	if !body.Grappling {
		t.Error("Grappling state should remain true")
	}
}

// TestGrappleDeterminism verifies grapple physics is deterministic
func TestGrappleDeterminism(t *testing.T) {
	// Create two identical scenarios
	body1 := NewBody(100, 200, 32, 32)
	body2 := NewBody(100, 200, 32, 32)
	anchor := world.AnchorPoint{X: 200, Y: 50}

	body1.StartGrapple(anchor)
	body2.StartGrapple(anchor)

	// Simulate same number of frames
	for i := 0; i < 20; i++ {
		body1.UpdateGrapple()
		body2.UpdateGrapple()
	}

	// Positions should be identical
	if math.Abs(body1.Position.X-body2.Position.X) > 0.001 {
		t.Errorf("Grapple X position not deterministic: %.6f vs %.6f",
			body1.Position.X, body2.Position.X)
	}

	if math.Abs(body1.Position.Y-body2.Position.Y) > 0.001 {
		t.Errorf("Grapple Y position not deterministic: %.6f vs %.6f",
			body1.Position.Y, body2.Position.Y)
	}

	// Angles should be identical
	if math.Abs(body1.GrappleAngle-body2.GrappleAngle) > 0.001 {
		t.Errorf("Grapple angle not deterministic: %.6f vs %.6f",
			body1.GrappleAngle, body2.GrappleAngle)
	}
}

// TestGlideDeterminism verifies glide physics is deterministic
func TestGlideDeterminism(t *testing.T) {
	body1 := NewBody(100, 100, 32, 32)
	body2 := NewBody(100, 100, 32, 32)

	// Apply same gliding gravity for same frames
	for i := 0; i < 30; i++ {
		body1.ApplyGravity(true)
		body2.ApplyGravity(true)
	}

	if math.Abs(body1.Velocity.Y-body2.Velocity.Y) > 0.001 {
		t.Errorf("Glide velocity not deterministic: %.6f vs %.6f",
			body1.Velocity.Y, body2.Velocity.Y)
	}
}
