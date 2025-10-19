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
	
	// Can't jump when not on ground and no double jump
	if body.Jump(false, &doubleJumpUsed) {
		t.Error("Should not be able to jump when not on ground")
	}
	
	// Can jump when on ground
	body.OnGround = true
	if !body.Jump(false, &doubleJumpUsed) {
		t.Error("Should be able to jump when on ground")
	}
	
	if body.Velocity.Y != PlayerJumpSpeed {
		t.Errorf("Expected jump velocity %f, got %f", PlayerJumpSpeed, body.Velocity.Y)
	}
	
	// Can double jump
	body.OnGround = false
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
