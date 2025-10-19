// Package entity provides enemy AI behavior implementations that control
// enemy movement patterns, decision making, and attack logic based on
// procedurally assigned behavior patterns.
package entity

import (
	"math"

	"github.com/opd-ai/vania/internal/animation"
)

// EnemyInstance represents a runtime instance of an enemy with position and state
type EnemyInstance struct {
	Enemy          *Enemy
	X, Y           float64
	VelX, VelY     float64
	CurrentHealth  int
	State          EnemyState
	PatrolMinX     float64
	PatrolMaxX     float64
	PatrolDir      float64
	AttackCooldown int
	OnGround       bool
	AggroRange     float64
	AttackRange    float64
	AnimController *animation.AnimationController
}

// EnemyState represents current enemy state
type EnemyState int

const (
	IdleState EnemyState = iota
	PatrolState
	ChaseState
	AttackState
	FleeState
	DeadState
)

// NewEnemyInstance creates a new enemy runtime instance
func NewEnemyInstance(enemy *Enemy, x, y float64) *EnemyInstance {
	// Set aggro and attack ranges based on enemy size and behavior
	aggroRange := 200.0
	attackRange := 32.0
	
	if enemy.Size == LargeEnemy || enemy.Size == BossEnemy {
		aggroRange = 300.0
		attackRange = 48.0
	}
	
	if enemy.Behavior == FlyingBehavior {
		aggroRange = 250.0
	}
	
	return &EnemyInstance{
		Enemy:          enemy,
		X:              x,
		Y:              y,
		VelX:           0,
		VelY:           0,
		CurrentHealth:  enemy.Health,
		State:          IdleState,
		PatrolMinX:     x - 100,
		PatrolMaxX:     x + 100,
		PatrolDir:      1.0,
		AttackCooldown: 0,
		OnGround:       true,
		AggroRange:     aggroRange,
		AttackRange:    attackRange,
	}
}

// Update updates enemy AI behavior
func (ei *EnemyInstance) Update(playerX, playerY float64) {
	if ei.CurrentHealth <= 0 {
		ei.State = DeadState
		return
	}
	
	// Decrease attack cooldown
	if ei.AttackCooldown > 0 {
		ei.AttackCooldown--
	}
	
	// Calculate distance to player
	dx := playerX - ei.X
	dy := playerY - ei.Y
	distToPlayer := math.Sqrt(dx*dx + dy*dy)
	
	// Update behavior based on pattern
	switch ei.Enemy.Behavior {
	case PatrolBehavior:
		ei.updatePatrolBehavior(distToPlayer, dx, dy)
	case ChaseBehavior:
		ei.updateChaseBehavior(distToPlayer, dx, dy)
	case FleeBehavior:
		ei.updateFleeBehavior(distToPlayer, dx, dy)
	case StationaryBehavior:
		ei.updateStationaryBehavior(distToPlayer, dx, dy)
	case FlyingBehavior:
		ei.updateFlyingBehavior(distToPlayer, dx, dy)
	case JumpingBehavior:
		ei.updateJumpingBehavior(distToPlayer, dx, dy)
	}
	
	// Apply velocity limits
	maxSpeed := ei.Enemy.Speed
	if ei.VelX > maxSpeed {
		ei.VelX = maxSpeed
	} else if ei.VelX < -maxSpeed {
		ei.VelX = -maxSpeed
	}
	
	// Apply gravity for ground-based enemies
	if ei.Enemy.Behavior != FlyingBehavior && !ei.OnGround {
		ei.VelY += 0.5 // Gravity
		if ei.VelY > 10.0 {
			ei.VelY = 10.0
		}
	}
}

// updatePatrolBehavior implements patrol AI
func (ei *EnemyInstance) updatePatrolBehavior(distToPlayer, dx, dy float64) {
	// Check if player is in aggro range
	if distToPlayer < ei.AggroRange {
		ei.State = ChaseState
		ei.chasePlayer(dx, dy)
		return
	}
	
	// Patrol between min and max X
	ei.State = PatrolState
	ei.VelX = ei.Enemy.Speed * ei.PatrolDir
	
	// Reverse direction at boundaries
	if ei.X >= ei.PatrolMaxX {
		ei.PatrolDir = -1.0
	} else if ei.X <= ei.PatrolMinX {
		ei.PatrolDir = 1.0
	}
}

// updateChaseBehavior implements chase AI
func (ei *EnemyInstance) updateChaseBehavior(distToPlayer, dx, dy float64) {
	if distToPlayer < ei.AttackRange && ei.AttackCooldown <= 0 {
		ei.State = AttackState
		ei.VelX = 0
		ei.AttackCooldown = 60 // 1 second cooldown at 60 FPS
		return
	}
	
	ei.State = ChaseState
	ei.chasePlayer(dx, dy)
}

// updateFleeBehavior implements flee AI
func (ei *EnemyInstance) updateFleeBehavior(distToPlayer, dx, dy float64) {
	// Flee if player is too close
	if distToPlayer < ei.AggroRange {
		ei.State = FleeState
		if dx > 0 {
			ei.VelX = -ei.Enemy.Speed
		} else {
			ei.VelX = ei.Enemy.Speed
		}
	} else {
		ei.State = IdleState
		ei.VelX *= 0.8 // Friction
	}
}

// updateStationaryBehavior implements stationary AI
func (ei *EnemyInstance) updateStationaryBehavior(distToPlayer, dx, dy float64) {
	ei.VelX = 0
	
	if distToPlayer < ei.AttackRange && ei.AttackCooldown <= 0 {
		ei.State = AttackState
		ei.AttackCooldown = 90 // Longer cooldown for stationary
	} else {
		ei.State = IdleState
	}
}

// updateFlyingBehavior implements flying AI
func (ei *EnemyInstance) updateFlyingBehavior(distToPlayer, dx, dy float64) {
	if distToPlayer < ei.AttackRange && ei.AttackCooldown <= 0 {
		ei.State = AttackState
		ei.VelX = 0
		ei.VelY = 0
		ei.AttackCooldown = 60
		return
	}
	
	if distToPlayer < ei.AggroRange {
		ei.State = ChaseState
		// Move toward player in both X and Y
		ei.VelX = (dx / distToPlayer) * ei.Enemy.Speed
		ei.VelY = (dy / distToPlayer) * ei.Enemy.Speed
	} else {
		ei.State = PatrolState
		// Hover slowly
		ei.VelX = ei.Enemy.Speed * 0.3 * ei.PatrolDir
		ei.VelY = math.Sin(ei.X*0.1) * 0.5
		
		if ei.X >= ei.PatrolMaxX || ei.X <= ei.PatrolMinX {
			ei.PatrolDir *= -1
		}
	}
}

// updateJumpingBehavior implements jumping AI
func (ei *EnemyInstance) updateJumpingBehavior(distToPlayer, dx, dy float64) {
	if distToPlayer < ei.AttackRange && ei.AttackCooldown <= 0 {
		ei.State = AttackState
		ei.AttackCooldown = 60
		return
	}
	
	if distToPlayer < ei.AggroRange {
		ei.State = ChaseState
		ei.chasePlayer(dx, dy)
		
		// Jump toward player if on ground
		if ei.OnGround && ei.AttackCooldown <= 0 {
			ei.VelY = -8.0
			ei.AttackCooldown = 30
		}
	} else {
		ei.State = PatrolState
		ei.VelX = ei.Enemy.Speed * ei.PatrolDir
		
		if ei.X >= ei.PatrolMaxX || ei.X <= ei.PatrolMinX {
			ei.PatrolDir *= -1
		}
	}
}

// chasePlayer moves enemy toward player
func (ei *EnemyInstance) chasePlayer(dx, dy float64) {
	if dx > 0 {
		ei.VelX = ei.Enemy.Speed
	} else {
		ei.VelX = -ei.Enemy.Speed
	}
}

// TakeDamage applies damage to enemy
func (ei *EnemyInstance) TakeDamage(damage int) {
	ei.CurrentHealth -= damage
	if ei.CurrentHealth < 0 {
		ei.CurrentHealth = 0
	}
}

// IsDead checks if enemy is dead
func (ei *EnemyInstance) IsDead() bool {
	return ei.CurrentHealth <= 0
}

// GetAttackDamage returns damage dealt by enemy attack
func (ei *EnemyInstance) GetAttackDamage() int {
	if ei.State != AttackState {
		return 0
	}
	return ei.Enemy.Damage
}

// GetBounds returns enemy collision bounds
func (ei *EnemyInstance) GetBounds() (x, y, width, height float64) {
	width = 32.0
	height = 32.0
	
	switch ei.Enemy.Size {
	case SmallEnemy:
		width, height = 16.0, 16.0
	case MediumEnemy:
		width, height = 32.0, 32.0
	case LargeEnemy:
		width, height = 64.0, 64.0
	case BossEnemy:
		width, height = 128.0, 128.0
	}
	
	return ei.X, ei.Y, width, height
}
