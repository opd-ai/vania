package physics

import (
	"testing"

	"github.com/opd-ai/vania/internal/world"
)

func TestNewBody(t *testing.T) {
	body := NewBody(100, 200, 32, 32)

	if body == nil {
		t.Fatal("NewBody returned nil")
	}

	if body.Position.X != 100 {
		t.Errorf("Expected X=100, got %f", body.Position.X)
	}

	if body.Position.Y != 200 {
		t.Errorf("Expected Y=200, got %f", body.Position.Y)
	}

	if body.Position.Width != 32 {
		t.Errorf("Expected Width=32, got %f", body.Position.Width)
	}

	if body.Position.Height != 32 {
		t.Errorf("Expected Height=32, got %f", body.Position.Height)
	}

	if body.OnGround {
		t.Error("Body should not be on ground initially")
	}
}

func TestApplyGravity(t *testing.T) {
	body := NewBody(100, 100, 32, 32)

	// Not on ground, gravity should apply
	body.ApplyGravity()

	if body.Velocity.Y != Gravity {
		t.Errorf("Expected velocity Y=%f, got %f", Gravity, body.Velocity.Y)
	}

	// Apply gravity multiple times
	for i := 0; i < 30; i++ {
		body.ApplyGravity()
	}

	// Should cap at max fall speed
	if body.Velocity.Y > MaxFallSpeed {
		t.Errorf("Velocity Y exceeded max fall speed: %f > %f", body.Velocity.Y, MaxFallSpeed)
	}

	// On ground, gravity should not apply
	body.OnGround = true
	prevVel := body.Velocity.Y
	body.ApplyGravity()

	if body.Velocity.Y != prevVel {
		t.Error("Gravity should not apply when on ground")
	}
}

func TestUpdate(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	body.Velocity.X = 5
	body.Velocity.Y = 10

	body.Update()

	if body.Position.X != 105 {
		t.Errorf("Expected X=105, got %f", body.Position.X)
	}

	if body.Position.Y != 110 {
		t.Errorf("Expected Y=110, got %f", body.Position.Y)
	}
}

func TestCheckCollision(t *testing.T) {
	a := AABB{X: 0, Y: 0, Width: 32, Height: 32}
	b := AABB{X: 16, Y: 16, Width: 32, Height: 32}

	// These should overlap
	if !CheckCollision(a, b) {
		t.Error("Expected collision between overlapping AABBs")
	}

	// No overlap
	c := AABB{X: 100, Y: 100, Width: 32, Height: 32}
	if CheckCollision(a, c) {
		t.Error("No collision expected for non-overlapping AABBs")
	}
}

func TestMoveHorizontal(t *testing.T) {
	body := NewBody(100, 100, 32, 32)

	// Move right
	body.MoveHorizontal(1.0)
	if body.Velocity.X != PlayerSpeed {
		t.Errorf("Expected velocity X=%f, got %f", PlayerSpeed, body.Velocity.X)
	}

	// Move left
	body.MoveHorizontal(-1.0)
	if body.Velocity.X != -PlayerSpeed {
		t.Errorf("Expected velocity X=%f, got %f", -PlayerSpeed, body.Velocity.X)
	}
}

func TestJump(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	doubleJumpUsed := false

	// Can't jump when not on ground and no double jump and outside coyote window
	body.FramesSinceGrounded = CoyoteFrames + 5 // Well outside coyote window
	if body.Jump(false, &doubleJumpUsed) {
		t.Error("Should not be able to jump when not on ground and outside coyote window")
	}

	// Can jump when on ground
	body.OnGround = true
	body.FramesSinceGrounded = 0
	if !body.Jump(false, &doubleJumpUsed) {
		t.Error("Should be able to jump when on ground")
	}

	if body.Velocity.Y != PlayerJumpSpeed {
		t.Errorf("Expected jump velocity %f, got %f", PlayerJumpSpeed, body.Velocity.Y)
	}

	// Can double jump
	body.OnGround = false
	body.FramesSinceGrounded = CoyoteFrames + 5 // Outside coyote window
	if !body.Jump(true, &doubleJumpUsed) {
		t.Error("Should be able to double jump")
	}

	if !doubleJumpUsed {
		t.Error("Double jump should be marked as used")
	}

	// Can't double jump again
	if body.Jump(true, &doubleJumpUsed) {
		t.Error("Should not be able to double jump twice")
	}
}

func TestDash(t *testing.T) {
	body := NewBody(100, 100, 32, 32)

	// Dash right
	body.Dash(1.0)
	if body.Velocity.X != PlayerDashSpeed {
		t.Errorf("Expected dash velocity %f, got %f", PlayerDashSpeed, body.Velocity.X)
	}

	// Dash left
	body.Dash(-1.0)
	if body.Velocity.X != -PlayerDashSpeed {
		t.Errorf("Expected dash velocity %f, got %f", -PlayerDashSpeed, body.Velocity.X)
	}
}

func TestApplyFriction(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	body.Velocity.X = 10
	body.OnGround = true

	body.ApplyFriction()

	// Velocity should decrease
	if body.Velocity.X >= 10 {
		t.Error("Friction should reduce velocity")
	}

	// Very slow velocity should stop
	body.Velocity.X = 0.05
	body.ApplyFriction()

	if body.Velocity.X != 0 {
		t.Error("Very slow velocity should be stopped by friction")
	}
}

func TestResolveCollisionWithPlatforms(t *testing.T) {
	body := NewBody(100, 460, 32, 32)
	body.Velocity.Y = 10 // Falling

	// Create a platform below
	platforms := []world.Platform{
		{X: 50, Y: 500, Width: 200, Height: 32},
	}

	// Update position (would move down into platform)
	body.Update()

	// Resolve collision
	body.ResolveCollisionWithPlatforms(platforms)

	// Should be on ground now
	if !body.OnGround {
		t.Error("Body should be on ground after landing on platform")
	}

	// Position should be on top of platform
	expectedY := 500.0 - 32.0 // platform Y - body height
	if body.Position.Y != expectedY {
		t.Errorf("Expected Y=%f after landing, got %f", expectedY, body.Position.Y)
	}

	// Velocity Y should be 0
	if body.Velocity.Y != 0 {
		t.Errorf("Velocity Y should be 0 after landing, got %f", body.Velocity.Y)
	}
}

func TestScreenBoundaries(t *testing.T) {
	body := NewBody(100, 100, 32, 32)

	// Move beyond left boundary
	body.Position.X = -10
	body.Velocity.X = -5
	body.ResolveCollisionWithPlatforms([]world.Platform{})

	if body.Position.X < 0 {
		t.Error("Body should be constrained to left boundary")
	}

	// Move beyond right boundary
	body.Position.X = 950
	body.Velocity.X = 5
	body.ResolveCollisionWithPlatforms([]world.Platform{})

	if body.Position.X+body.Position.Width > 960 {
		t.Error("Body should be constrained to right boundary")
	}

	// Move beyond bottom boundary
	body.Position.Y = 650
	body.Velocity.Y = 5
	body.ResolveCollisionWithPlatforms([]world.Platform{})

	if body.Position.Y+body.Position.Height > 640 {
		t.Error("Body should be constrained to bottom boundary")
	}

	if !body.OnGround {
		t.Error("Body should be on ground when at bottom boundary")
	}
}

func TestReleaseJump(t *testing.T) {
	testCases := []struct {
		name             string
		initialVelocity  float64
		expectedVelocity float64
		description      string
	}{
		{
			name:             "ReleaseJumpDuringAscent",
			initialVelocity:  -10.0,
			expectedVelocity: -10.0 * JumpReleaseDamping,
			description:      "Upward velocity should be damped when jump released during ascent",
		},
		{
			name:             "ReleaseJumpAtPeak",
			initialVelocity:  -0.5,
			expectedVelocity: -0.5 * JumpReleaseDamping,
			description:      "Small upward velocity should still be damped",
		},
		{
			name:             "ReleaseJumpWhileFalling",
			initialVelocity:  5.0,
			expectedVelocity: 5.0,
			description:      "Downward velocity should not be affected by jump release",
		},
		{
			name:             "ReleaseJumpWhenStationary",
			initialVelocity:  0.0,
			expectedVelocity: 0.0,
			description:      "Zero velocity should remain zero",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := NewBody(100, 100, 32, 32)
			body.Velocity.Y = tc.initialVelocity

			body.ReleaseJump()

			if body.Velocity.Y != tc.expectedVelocity {
				t.Errorf("%s: Expected velocity Y=%f, got %f", tc.description, tc.expectedVelocity, body.Velocity.Y)
			}
		})
	}
}

func TestVariableHeightJump(t *testing.T) {
	// Test full jump height (hold button entire ascent)
	t.Run("FullJumpHeight", func(t *testing.T) {
		body := NewBody(100, 500, 32, 32)
		body.OnGround = true
		doubleJumpUsed := false

		// Jump
		body.Jump(false, &doubleJumpUsed)

		// Simulate full ascent without releasing
		maxHeight := body.Position.Y
		for i := 0; i < 100; i++ {
			body.ApplyGravity()
			body.Update()
			if body.Velocity.Y > 0 {
				// Started falling, record peak
				maxHeight = body.Position.Y
				break
			}
		}

		fullJumpDistance := 500 - maxHeight

		// Test short jump height (release button immediately)
		body2 := NewBody(100, 500, 32, 32)
		body2.OnGround = true
		doubleJumpUsed2 := false

		// Jump
		body2.Jump(false, &doubleJumpUsed2)

		// Release jump immediately
		body2.ReleaseJump()

		// Simulate ascent
		maxHeight2 := body2.Position.Y
		for i := 0; i < 100; i++ {
			body2.ApplyGravity()
			body2.Update()
			if body2.Velocity.Y > 0 {
				// Started falling, record peak
				maxHeight2 = body2.Position.Y
				break
			}
		}

		shortJumpDistance := 500 - maxHeight2

		// Short jump should be noticeably shorter than full jump
		// With damping of 0.5, short jump should be roughly 25-35% of full jump
		ratio := shortJumpDistance / fullJumpDistance
		if ratio > 0.4 {
			t.Errorf("Short jump too high: %.2f%% of full jump (expected < 40%%)", ratio*100)
		}
		if ratio < 0.2 {
			t.Errorf("Short jump too low: %.2f%% of full jump (expected > 20%%)", ratio*100)
		}
	})

	// Test mid-release jump height
	t.Run("MidReleaseJumpHeight", func(t *testing.T) {
		body := NewBody(100, 500, 32, 32)
		body.OnGround = true
		doubleJumpUsed := false

		// Jump
		body.Jump(false, &doubleJumpUsed)

		// Simulate partial ascent before releasing
		for i := 0; i < 5; i++ {
			body.ApplyGravity()
			body.Update()
		}

		// Release jump mid-ascent
		body.ReleaseJump()

		// Continue to peak
		maxHeight := body.Position.Y
		for i := 0; i < 100; i++ {
			body.ApplyGravity()
			body.Update()
			if body.Velocity.Y > 0 {
				maxHeight = body.Position.Y
				break
			}
		}

		midJumpDistance := 500 - maxHeight

		// Mid-release should produce medium jump height (40-70% of full)
		// This is approximate since physics is discrete
		if midJumpDistance < 20 {
			t.Error("Mid-release jump should reach noticeable height")
		}
	})
}

func TestWallSlide(t *testing.T) {
	testCases := []struct {
		name         string
		onWall       bool
		wallSide     int
		initialVelY  float64
		expectedMaxY float64
		description  string
	}{
		{
			name:         "WallSlideActivates",
			onWall:       true,
			wallSide:     1,
			initialVelY:  0.0,
			expectedMaxY: WallSlideSpeed,
			description:  "Fall speed should be capped at WallSlideSpeed when sliding on wall",
		},
		{
			name:         "NoWallSlideInAir",
			onWall:       false,
			wallSide:     0,
			initialVelY:  0.0,
			expectedMaxY: MaxFallSpeed,
			description:  "Fall speed should reach MaxFallSpeed when not on wall",
		},
		{
			name:         "WallSlideOnLeftWall",
			onWall:       true,
			wallSide:     -1,
			initialVelY:  0.0,
			expectedMaxY: WallSlideSpeed,
			description:  "Wall slide should work on left wall",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := NewBody(100, 100, 32, 32)
			body.OnWall = tc.onWall
			body.WallSide = tc.wallSide
			body.Velocity.Y = tc.initialVelY

			// Apply gravity for many frames to reach terminal velocity
			for i := 0; i < 50; i++ {
				body.ApplyGravity()
			}

			if body.Velocity.Y > tc.expectedMaxY+0.1 {
				t.Errorf("%s: Fall speed %.2f exceeded expected max %.2f",
					tc.description, body.Velocity.Y, tc.expectedMaxY)
			}

			// For wall-slide cases, verify it's significantly slower than max fall speed
			if tc.onWall && body.Velocity.Y > WallSlideSpeed+0.1 {
				t.Errorf("Wall slide not working: velocity %.2f > WallSlideSpeed %.2f",
					body.Velocity.Y, WallSlideSpeed)
			}
		})
	}
}

func TestCoyoteTime(t *testing.T) {
	testCases := []struct {
		name              string
		framesSinceGround int
		shouldAllowJump   bool
		description       string
	}{
		{
			name:              "JumpWithinCoyoteWindow",
			framesSinceGround: 3,
			shouldAllowJump:   true,
			description:       "Should allow jump within coyote-time window",
		},
		{
			name:              "JumpAtCoyoteEdge",
			framesSinceGround: CoyoteFrames,
			shouldAllowJump:   true,
			description:       "Should allow jump exactly at coyote frame limit",
		},
		{
			name:              "JumpAfterCoyoteExpired",
			framesSinceGround: CoyoteFrames + 1,
			shouldAllowJump:   false,
			description:       "Should not allow jump after coyote-time expires",
		},
		{
			name:              "JumpImmediatelyAfterLeaving",
			framesSinceGround: 1,
			shouldAllowJump:   true,
			description:       "Should allow jump on first frame after leaving ground",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := NewBody(100, 100, 32, 32)
			body.OnGround = false
			body.FramesSinceGrounded = tc.framesSinceGround
			doubleJumpUsed := false

			result := body.Jump(false, &doubleJumpUsed)

			if result != tc.shouldAllowJump {
				t.Errorf("%s: Jump returned %v, expected %v", tc.description, result, tc.shouldAllowJump)
			}

			if tc.shouldAllowJump && body.Velocity.Y != PlayerJumpSpeed {
				t.Errorf("Jump velocity incorrect: got %.2f, expected %.2f", body.Velocity.Y, PlayerJumpSpeed)
			}
		})
	}
}

func TestJumpBuffer(t *testing.T) {
	t.Run("BufferedJumpExecutesOnLanding", func(t *testing.T) {
		// Start body above platform
		body := NewBody(100, 460, 32, 32)
		body.OnGround = false
		body.Velocity.Y = 5.0

		// Set buffer timer
		body.JumpBufferTimer = JumpBufferFrames

		// Create a platform for landing
		platforms := []world.Platform{
			{X: 50, Y: 500, Width: 200, Height: 32},
		}

		// Simulate frames until landing
		landed := false
		for i := 0; i < 10; i++ {
			body.ApplyGravity()
			body.Update()
			body.ResolveCollisionWithPlatforms(platforms)

			if body.OnGround {
				landed = true
				break
			}
		}

		if !landed {
			t.Error("Body should have landed on platform")
			return
		}

		// Jump should have executed (negative velocity = upward)
		if body.Velocity.Y >= 0 {
			t.Errorf("Buffered jump should execute on landing, velocity.Y=%.2f", body.Velocity.Y)
		}

		if body.Velocity.Y != PlayerJumpSpeed {
			t.Errorf("Expected jump velocity %.2f, got %.2f", PlayerJumpSpeed, body.Velocity.Y)
		}
	})

	t.Run("BufferedJumpExpiresBeforeLanding", func(t *testing.T) {
		// Position body high in air so it takes many frames to land
		body := NewBody(100, 200, 32, 32)
		body.OnGround = false
		body.Velocity.Y = 1.0 // Start with small downward velocity

		// Set buffer timer to expire quickly
		body.JumpBufferTimer = 2

		// Create a platform far below
		platforms := []world.Platform{
			{X: 50, Y: 500, Width: 200, Height: 32},
		}

		// Run until buffer expires
		for i := 0; i < 3; i++ {
			body.ApplyGravity()
			body.Update()
			body.ResolveCollisionWithPlatforms(platforms)
		}

		// Buffer should be expired now
		if body.JumpBufferTimer != 0 {
			t.Errorf("Buffer should have expired, got timer=%d", body.JumpBufferTimer)
		}

		// Continue until landing
		for i := 0; i < 50; i++ {
			body.ApplyGravity()
			body.Update()
			body.ResolveCollisionWithPlatforms(platforms)

			if body.OnGround {
				break
			}
		}

		// Jump should NOT have executed (velocity should be 0 from landing, not negative from jump)
		if body.Velocity.Y < 0 {
			t.Errorf("Buffered jump should not execute after expiring, velocity.Y=%.2f", body.Velocity.Y)
		}
	})

	t.Run("BufferedJumpAtEdge", func(t *testing.T) {
		// Start body close to platform
		body := NewBody(100, 460, 32, 32)
		body.OnGround = false
		body.Velocity.Y = 5.0

		// Set buffer timer to exactly JumpBufferFrames
		body.JumpBufferTimer = JumpBufferFrames

		// Create a platform for landing
		platforms := []world.Platform{
			{X: 50, Y: 500, Width: 200, Height: 32},
		}

		// Simulate exactly JumpBufferFrames frames
		landed := false
		for i := 0; i < JumpBufferFrames; i++ {
			body.ApplyGravity()
			body.Update()
			body.ResolveCollisionWithPlatforms(platforms)

			if body.OnGround {
				landed = true
				break
			}
		}

		// If landed within window, jump should execute
		if landed && body.Velocity.Y >= 0 {
			t.Error("Buffered jump should execute when landing at buffer edge")
		}
	})
}

func TestBufferJumpMethod(t *testing.T) {
	body := NewBody(100, 100, 32, 32)

	if body.JumpBufferTimer != 0 {
		t.Error("JumpBufferTimer should start at 0")
	}

	body.BufferJump()

	if body.JumpBufferTimer != JumpBufferFrames {
		t.Errorf("BufferJump should set timer to %d, got %d", JumpBufferFrames, body.JumpBufferTimer)
	}
}

func TestBufferDecrement(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	body.JumpBufferTimer = 5

	// Simulate collision resolution which decrements buffer
	platforms := []world.Platform{}
	body.ResolveCollisionWithPlatforms(platforms)

	if body.JumpBufferTimer != 4 {
		t.Errorf("JumpBufferTimer should decrement to 4, got %d", body.JumpBufferTimer)
	}

	// Decrement to 0
	for i := 0; i < 5; i++ {
		body.ResolveCollisionWithPlatforms(platforms)
	}

	if body.JumpBufferTimer != 0 {
		t.Errorf("JumpBufferTimer should not go below 0, got %d", body.JumpBufferTimer)
	}
}

func TestCoyoteTimeTracking(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	platforms := []world.Platform{}

	// Simulate being on ground initially
	body.OnGround = true
	body.FramesSinceGrounded = 0

	// Call resolve to confirm we stay on ground
	// (body.OnGround will be set to false then back to true if on ground)
	// But with no platforms, it will be false
	// So we need a platform or to manually set OnGround

	// Actually, let's use the screen boundary as ground
	body.Position.Y = 640 - body.Position.Height
	body.Velocity.Y = 0
	body.ResolveCollisionWithPlatforms(platforms)

	if body.FramesSinceGrounded != 0 {
		t.Errorf("FramesSinceGrounded should be 0 when on ground, got %d", body.FramesSinceGrounded)
	}

	// Move body up (off ground)
	body.Position.Y = 500
	body.OnGround = false

	// First frame in air - should start counting
	body.ResolveCollisionWithPlatforms(platforms)

	if body.FramesSinceGrounded != 1 {
		t.Errorf("FramesSinceGrounded should be 1 on first frame after leaving ground, got %d", body.FramesSinceGrounded)
	}

	// Continue in air for 5 more frames
	for i := 0; i < 5; i++ {
		body.ResolveCollisionWithPlatforms(platforms)
	}

	if body.FramesSinceGrounded != 6 {
		t.Errorf("FramesSinceGrounded should be 6 after 6 frames in air, got %d", body.FramesSinceGrounded)
	}

	// Land again (use screen boundary)
	body.Position.Y = 640 - body.Position.Height
	body.Velocity.Y = 1 // Falling
	body.ResolveCollisionWithPlatforms(platforms)

	if body.FramesSinceGrounded != 0 {
		t.Errorf("FramesSinceGrounded should reset to 0 on landing, got %d", body.FramesSinceGrounded)
	}
}

func TestWallSlideAndCoyoteTimeCombination(t *testing.T) {
	// Test that wall-slide doesn't interfere with coyote-time tracking
	body := NewBody(100, 100, 32, 32)
	body.OnGround = false
	body.OnWall = true
	body.WallSide = 1
	body.FramesSinceGrounded = 2

	doubleJumpUsed := false

	// Should still allow jump via coyote-time even when on wall
	if !body.Jump(false, &doubleJumpUsed) {
		t.Error("Should allow jump via coyote-time even when on wall")
	}
}

func TestJumpConsumesBuffer(t *testing.T) {
	body := NewBody(100, 100, 32, 32)
	body.OnGround = true
	body.JumpBufferTimer = 5
	doubleJumpUsed := false

	// Execute jump
	body.Jump(false, &doubleJumpUsed)

	if body.JumpBufferTimer != 0 {
		t.Error("Jump should consume buffer timer")
	}
}
