// Package entity provides enemy AI behavior implementations that control
// enemy movement patterns, decision making, and attack logic based on
// procedurally assigned behavior patterns.
package entity

import (
	"math"

	"github.com/opd-ai/vania/internal/animation"
	"github.com/opd-ai/vania/internal/graphics"
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
	
	// Advanced AI fields
	Memory         *AIMemory      // Learning and pattern recognition
	TacticalState  TacticalState  // Current tactical decision state
	Group          *EnemyGroup    // Group coordination (nil if solo)
	FormationX     float64        // Target X position in formation
	FormationY     float64        // Target Y position in formation
	LastPlayerX    float64        // Track player position for learning
	LastPlayerY    float64
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
	
	// Create animation controller if sprite data is available
	var animController *animation.AnimationController
	if sprite, ok := enemy.SpriteData.(*graphics.Sprite); ok && sprite != nil {
		animController = CreateEnemyAnimController(sprite, enemy)
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
		AnimController: animController,
		Memory:         NewAIMemory(),
		TacticalState:  TacticalNormal,
		Group:          nil,
		FormationX:     x,
		FormationY:     y,
		LastPlayerX:    0,
		LastPlayerY:    0,
	}
}

// Update updates enemy AI behavior
func (ei *EnemyInstance) Update(playerX, playerY float64) {
	if ei.CurrentHealth <= 0 {
		ei.State = DeadState
		// Play death animation once
		if ei.AnimController != nil {
			currentAnim := ei.AnimController.GetCurrentAnimation()
			if currentAnim != "death" {
				ei.AnimController.Play("death", true)
			}
			ei.AnimController.Update()
		}
		return
	}
	
	// Update AI memory with player observations
	// Detect if player did actions (simplified detection for now)
	playerDidJump := math.Abs(playerY-ei.LastPlayerY) > 5.0 && playerY < ei.LastPlayerY
	playerDidAttack := false // Would need actual attack detection from game state
	playerDidDash := math.Abs(playerX-ei.LastPlayerX) > 10.0
	
	ei.Memory.UpdateMemory(playerX, playerY, playerDidJump, playerDidAttack, playerDidDash)
	ei.LastPlayerX = playerX
	ei.LastPlayerY = playerY
	
	// Decrease attack cooldown
	if ei.AttackCooldown > 0 {
		ei.AttackCooldown--
	}
	
	// Calculate distance to player
	dx := playerX - ei.X
	dy := playerY - ei.Y
	distToPlayer := math.Sqrt(dx*dx + dy*dy)
	
	// Determine tactical state based on AI memory
	healthPercent := float64(ei.CurrentHealth) / float64(ei.Enemy.Health)
	hasAllies := ei.Group != nil && len(ei.Group.Members) > 1
	ei.TacticalState = ei.Memory.GetTacticalState(healthPercent, hasAllies, distToPlayer)
	
	// Apply tactical state modifications to behavior
	ei.applyTacticalBehavior(distToPlayer, dx, dy, playerX, playerY)
	
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
	
	// Apply formation movement if in a group
	if ei.Group != nil && ei.Group.Formation != NoFormation {
		ei.applyFormationMovement()
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
	
	// Update animation controller
	if ei.AnimController != nil {
		ei.AnimController.Update()
		
		// Set animation based on state
		currentAnim := ei.AnimController.GetCurrentAnimation()
		
		switch ei.State {
		case AttackState:
			if currentAnim != "attack" {
				ei.AnimController.Play("attack", true)
			}
		case PatrolState, ChaseState, FleeState:
			if currentAnim != "patrol" && currentAnim != "attack" {
				ei.AnimController.Play("patrol", false)
			}
		case IdleState:
			if currentAnim != "idle" && currentAnim != "attack" {
				ei.AnimController.Play("idle", false)
			}
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
	
	// Record combat event in memory
	if ei.Memory != nil {
		ei.Memory.RecordCombatEvent(false, true, damage, 0)
	}
	
	// Play hit animation
	if ei.AnimController != nil && ei.CurrentHealth > 0 {
		ei.AnimController.Play("hit", true)
	}
	
	if ei.CurrentHealth < 0 {
		ei.CurrentHealth = 0
	}
}

// applyTacticalBehavior modifies enemy behavior based on tactical state
func (ei *EnemyInstance) applyTacticalBehavior(distToPlayer, dx, dy, playerX, playerY float64) {
	switch ei.TacticalState {
	case TacticalAggressive:
		// Increase aggro range and move faster when aggressive
		ei.AggroRange *= 1.2
		if distToPlayer < ei.AggroRange {
			ei.State = ChaseState
		}
		
	case TacticalDefensive:
		// Increase attack range, reduce aggro range when defensive
		ei.AggroRange *= 0.8
		ei.AttackRange *= 1.3
		
	case TacticalFlanking:
		// Try to get behind or beside player
		if ei.Group != nil && len(ei.Group.Members) > 1 {
			// Let group handle flanking through formation
			return
		}
		// Solo flanking: try to circle around
		angle := math.Atan2(dy, dx) + math.Pi/2 // 90 degrees offset
		targetX := playerX + math.Cos(angle)*100
		targetY := playerY + math.Sin(angle)*100
		
		fdx := targetX - ei.X
		_ = targetY - ei.Y // fdy unused for ground-based flanking
		if math.Abs(fdx) > 10 {
			if fdx > 0 {
				ei.VelX = ei.Enemy.Speed
			} else {
				ei.VelX = -ei.Enemy.Speed
			}
		}
		
	case TacticalKiting:
		// Hit and run: attack then retreat
		if distToPlayer < ei.AttackRange && ei.AttackCooldown <= 0 {
			ei.State = AttackState
			ei.AttackCooldown = 45
		} else if distToPlayer < ei.AttackRange*1.5 {
			// Retreat after attacking
			if dx > 0 {
				ei.VelX = -ei.Enemy.Speed
			} else {
				ei.VelX = ei.Enemy.Speed
			}
		}
		
	case TacticalRetreating:
		// Move away from player
		ei.State = FleeState
		if dx > 0 {
			ei.VelX = -ei.Enemy.Speed * 1.2
		} else {
			ei.VelX = ei.Enemy.Speed * 1.2
		}
		
	case TacticalRegrouping:
		// Move toward group center if we have a group
		if ei.Group != nil && len(ei.Group.Members) > 1 {
			centerX, centerY := 0.0, 0.0
			count := 0
			for _, member := range ei.Group.Members {
				if member.State != DeadState {
					centerX += member.X
					centerY += member.Y
					count++
				}
			}
			if count > 0 {
				centerX /= float64(count)
				centerY /= float64(count)
				
				gdx := centerX - ei.X
				_ = centerY - ei.Y // gdy unused for horizontal regrouping
				
				if math.Abs(gdx) > 10 {
					if gdx > 0 {
						ei.VelX = ei.Enemy.Speed
					} else {
						ei.VelX = -ei.Enemy.Speed
					}
				}
			}
		}
	}
}

// applyFormationMovement moves enemy toward formation position
func (ei *EnemyInstance) applyFormationMovement() {
	// Calculate distance to formation position
	dx := ei.FormationX - ei.X
	dy := ei.FormationY - ei.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	
	// Only apply formation movement if we're far from position
	if dist > 30.0 {
		// Blend formation movement with current velocity
		formationInfluence := 0.3
		
		targetVelX := (dx / dist) * ei.Enemy.Speed
		targetVelY := (dy / dist) * ei.Enemy.Speed
		
		ei.VelX = ei.VelX*(1.0-formationInfluence) + targetVelX*formationInfluence
		
		// Only apply Y velocity for flying enemies
		if ei.Enemy.Behavior == FlyingBehavior {
			ei.VelY = ei.VelY*(1.0-formationInfluence) + targetVelY*formationInfluence
		}
	}
}

// RecordSuccessfulHit records when this enemy successfully hit the player
func (ei *EnemyInstance) RecordSuccessfulHit(distance float64) {
	if ei.Memory != nil {
		ei.Memory.RecordCombatEvent(true, false, 0, distance)
	}
}

// RecordEvasion records when this enemy successfully evaded an attack
func (ei *EnemyInstance) RecordEvasion() {
	if ei.Memory != nil {
		ei.Memory.RecordEvasion()
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

// CreateEnemyAnimController creates an animation controller for an enemy
func CreateEnemyAnimController(baseSprite *graphics.Sprite, enemy *Enemy) *animation.AnimationController {
	// Use enemy's biome and danger level as seed for animation generation
	seed := int64(enemy.DangerLevel * 12345)
	if enemy.BiomeType != "" {
		for _, c := range enemy.BiomeType {
			seed += int64(c)
		}
	}
	
	animGen := animation.NewAnimationGenerator(seed)
	
	// Generate animation frames
	idleFrames := animGen.GenerateEnemyIdleFrames(baseSprite, 4)
	patrolFrames := animGen.GenerateEnemyPatrolFrames(baseSprite, 4)
	attackFrames := animGen.GenerateEnemyAttackFrames(baseSprite, 3)
	deathFrames := animGen.GenerateEnemyDeathFrames(baseSprite, 4)
	hitFrames := animGen.GenerateHitFrames(baseSprite, 2)
	
	// Create animation controller with idle as default
	animController := animation.NewAnimationController("idle")
	animController.AddAnimation(animation.NewAnimation("idle", idleFrames, 15, true))
	animController.AddAnimation(animation.NewAnimation("patrol", patrolFrames, 8, true))
	animController.AddAnimation(animation.NewAnimation("attack", attackFrames, 5, false))
	animController.AddAnimation(animation.NewAnimation("death", deathFrames, 10, false))
	animController.AddAnimation(animation.NewAnimation("hit", hitFrames, 3, false))
	
	return animController
}
