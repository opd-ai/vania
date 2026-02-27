// Package physics provides collision detection, movement physics, and
// gravity simulation for game entities, platforms, and the player character.
package physics

import (
	"github.com/opd-ai/vania/internal/world"
)

const (
	// Physics constants
	Gravity         = 0.5
	MaxFallSpeed    = 10.0
	PlayerSpeed     = 4.0
	PlayerJumpSpeed = -12.0
	PlayerDashSpeed = 8.0
	PlayerWidth     = 32
	PlayerHeight    = 32

	// JumpReleaseDamping is applied to upward velocity when jump button is released.
	// Multiplies current Y velocity to reduce jump height on early release.
	// Value of 0.5 cuts jump height approximately in half for instant release.
	JumpReleaseDamping = 0.5

	// WallSlideSpeed is the maximum fall speed when sliding down a wall (units/frame).
	// Reduced from MaxFallSpeed to create a slow, controlled descent.
	WallSlideSpeed = 2.0

	// CoyoteFrames is the grace period (in frames at 60fps) after leaving a ledge
	// during which the player can still jump. Industry standard: ~100ms = 6 frames.
	CoyoteFrames = 6

	// JumpBufferFrames is the window (in frames at 60fps) during which a jump
	// input is buffered and will execute upon landing. Industry standard: ~100ms = 6 frames.
	JumpBufferFrames = 6
)

// AABB represents an axis-aligned bounding box
type AABB struct {
	X, Y          float64
	Width, Height float64
}

// Body represents a physics body
type Body struct {
	Position            AABB
	Velocity            Vector2D
	OnGround            bool
	OnWall              bool
	WallSide            int // -1 for left, 1 for right, 0 for none
	FramesSinceGrounded int // Used for coyote-time
	JumpBufferTimer     int // Countdown for buffered jump input
}

// Vector2D represents a 2D vector
type Vector2D struct {
	X, Y float64
}

// NewBody creates a new physics body
func NewBody(x, y, width, height float64) *Body {
	return &Body{
		Position: AABB{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Velocity:            Vector2D{X: 0, Y: 0},
		OnGround:            false,
		OnWall:              false,
		WallSide:            0,
		FramesSinceGrounded: 0,
		JumpBufferTimer:     0,
	}
}

// ApplyGravity applies gravity to the body with wall-slide support
func (b *Body) ApplyGravity() {
	if !b.OnGround {
		b.Velocity.Y += Gravity

		// Wall-slide: slow fall speed when sliding down a wall
		if b.OnWall && b.Velocity.Y > WallSlideSpeed {
			b.Velocity.Y = WallSlideSpeed
		} else if b.Velocity.Y > MaxFallSpeed {
			b.Velocity.Y = MaxFallSpeed
		}
	}
}

// Update updates the body position
func (b *Body) Update() {
	b.Position.X += b.Velocity.X
	b.Position.Y += b.Velocity.Y
}

// CheckCollision checks if two AABBs collide
func CheckCollision(a, b AABB) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// ResolveCollisionWithPlatforms checks and resolves collisions with platforms
func (b *Body) ResolveCollisionWithPlatforms(platforms []world.Platform) {
	// Store previous state for coyote-time tracking
	wasOnGround := b.OnGround

	// Reset ground and wall state
	b.OnGround = false
	b.OnWall = false
	b.WallSide = 0

	for _, platform := range platforms {
		platformAABB := AABB{
			X:      float64(platform.X),
			Y:      float64(platform.Y),
			Width:  float64(platform.Width),
			Height: float64(platform.Height),
		}

		if CheckCollision(b.Position, platformAABB) {
			// Determine collision direction
			// Check if player was above the platform (landed on top)
			if b.Velocity.Y > 0 && b.Position.Y+b.Position.Height-b.Velocity.Y <= platformAABB.Y {
				// Collision from top
				b.Position.Y = platformAABB.Y - b.Position.Height
				b.Velocity.Y = 0
				b.OnGround = true
			} else if b.Velocity.Y < 0 && b.Position.Y-b.Velocity.Y >= platformAABB.Y+platformAABB.Height {
				// Collision from bottom
				b.Position.Y = platformAABB.Y + platformAABB.Height
				b.Velocity.Y = 0
			} else if b.Velocity.X > 0 {
				// Collision from left
				b.Position.X = platformAABB.X - b.Position.Width
				b.Velocity.X = 0
				b.OnWall = true
				b.WallSide = 1
			} else if b.Velocity.X < 0 {
				// Collision from right
				b.Position.X = platformAABB.X + platformAABB.Width
				b.Velocity.X = 0
				b.OnWall = true
				b.WallSide = -1
			}
		}
	}

	// Check screen boundaries (floor at bottom)
	if b.Position.Y+b.Position.Height >= 640 {
		b.Position.Y = 640 - b.Position.Height
		b.Velocity.Y = 0
		b.OnGround = true
	}

	// Keep player on screen (left/right boundaries)
	if b.Position.X < 0 {
		b.Position.X = 0
		b.Velocity.X = 0
	}
	if b.Position.X+b.Position.Width > 960 {
		b.Position.X = 960 - b.Position.Width
		b.Velocity.X = 0
	}

	// Update coyote-time tracking
	if b.OnGround {
		b.FramesSinceGrounded = 0
		// Execute buffered jump if any
		if b.JumpBufferTimer > 0 {
			b.Velocity.Y = PlayerJumpSpeed
			b.JumpBufferTimer = 0
		}
	} else if wasOnGround {
		// Just left ground, start counting
		b.FramesSinceGrounded = 1
	} else {
		// In air, increment counter
		b.FramesSinceGrounded++
	}

	// Decrement jump buffer timer
	if b.JumpBufferTimer > 0 {
		b.JumpBufferTimer--
	}

	// Landing detection (just landed)
	if !wasOnGround && b.OnGround {
		// Could trigger landing sound/animation here
	}
}

// MoveHorizontal applies horizontal movement
func (b *Body) MoveHorizontal(direction float64) {
	b.Velocity.X = direction * PlayerSpeed
}

// Jump makes the body jump if on ground, in coyote-time window, or wall.
// Returns true if jump was executed.
func (b *Body) Jump(hasDoubleJump bool, doubleJumpUsed *bool) bool {
	// Ground jump or coyote-time jump
	if b.OnGround || b.FramesSinceGrounded <= CoyoteFrames {
		b.Velocity.Y = PlayerJumpSpeed
		*doubleJumpUsed = false
		b.JumpBufferTimer = 0 // Consume buffered jump
		return true
	} else if hasDoubleJump && !*doubleJumpUsed {
		b.Velocity.Y = PlayerJumpSpeed
		*doubleJumpUsed = true
		b.JumpBufferTimer = 0
		return true
	} else if b.OnWall {
		// Wall jump
		b.Velocity.Y = PlayerJumpSpeed
		b.Velocity.X = float64(-b.WallSide) * PlayerSpeed * 1.5
		b.JumpBufferTimer = 0
		return true
	}
	return false
}

// BufferJump stores a jump input for execution upon landing.
// Should be called when jump is pressed but jump conditions aren't met.
func (b *Body) BufferJump() {
	b.JumpBufferTimer = JumpBufferFrames
}

// ReleaseJump applies variable-height jump mechanics.
// When called during upward movement (negative Y velocity), it reduces
// the jump height by damping the velocity. This allows for short-hop jumps
// when the jump button is released early.
func (b *Body) ReleaseJump() {
	// Only apply damping when moving upward (negative Y velocity)
	if b.Velocity.Y < 0 {
		b.Velocity.Y *= JumpReleaseDamping
	}
}

// Dash performs a dash move
func (b *Body) Dash(direction float64) {
	if direction != 0 {
		b.Velocity.X = direction * PlayerDashSpeed
	}
}

// ApplyFriction applies friction to horizontal movement
func (b *Body) ApplyFriction() {
	if b.OnGround {
		b.Velocity.X *= 0.8
		// Stop if moving very slowly
		if b.Velocity.X > -0.1 && b.Velocity.X < 0.1 {
			b.Velocity.X = 0
		}
	} else {
		// Air resistance
		b.Velocity.X *= 0.95
	}
}
