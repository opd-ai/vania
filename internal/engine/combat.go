// Package engine provides combat system implementation for managing
// player attacks, damage calculation, hit detection, and combat state.
package engine

import (
	"math"
	
	"github.com/opd-ai/vania/internal/entity"
)

// CombatSystem manages all combat interactions
type CombatSystem struct {
	playerAttackCooldown int
	playerAttacking      bool
	playerAttackFrame    int
	knockbackVelX        float64
	knockbackVelY        float64
	invulnerableFrames   int
}

// NewCombatSystem creates a new combat system
func NewCombatSystem() *CombatSystem {
	return &CombatSystem{
		playerAttackCooldown: 0,
		playerAttacking:      false,
		playerAttackFrame:    0,
		knockbackVelX:        0,
		knockbackVelY:        0,
		invulnerableFrames:   0,
	}
}

// Update updates combat system state
func (cs *CombatSystem) Update() {
	if cs.playerAttackCooldown > 0 {
		cs.playerAttackCooldown--
	}
	
	if cs.playerAttacking {
		cs.playerAttackFrame++
		if cs.playerAttackFrame > 15 { // Attack lasts 15 frames
			cs.playerAttacking = false
			cs.playerAttackFrame = 0
		}
	}
	
	if cs.invulnerableFrames > 0 {
		cs.invulnerableFrames--
	}
}

// PlayerAttack initiates a player attack
func (cs *CombatSystem) PlayerAttack() bool {
	if cs.playerAttackCooldown <= 0 {
		cs.playerAttacking = true
		cs.playerAttackFrame = 0
		cs.playerAttackCooldown = 20 // 20 frames between attacks
		return true
	}
	return false
}

// IsPlayerAttacking returns if player is currently attacking
func (cs *CombatSystem) IsPlayerAttacking() bool {
	return cs.playerAttacking
}

// GetAttackHitbox returns player attack hitbox during attack frames
func (cs *CombatSystem) GetAttackHitbox(playerX, playerY, facingDir float64) (x, y, width, height float64) {
	if !cs.playerAttacking || cs.playerAttackFrame < 3 || cs.playerAttackFrame > 10 {
		return 0, 0, 0, 0 // No hitbox outside active frames
	}
	
	// Attack hitbox in front of player
	width = 40.0
	height = 32.0
	x = playerX
	y = playerY
	
	if facingDir >= 0 {
		x = playerX + 32 // Right side
	} else {
		x = playerX - 40 // Left side
	}
	
	return x, y, width, height
}

// CheckEnemyHit checks if attack hit an enemy
func (cs *CombatSystem) CheckEnemyHit(attackX, attackY, attackW, attackH float64, enemy *entity.EnemyInstance) bool {
	if attackW <= 0 || attackH <= 0 {
		return false
	}
	
	ex, ey, ew, eh := enemy.GetBounds()
	
	// AABB collision check
	return attackX < ex+ew &&
		attackX+attackW > ex &&
		attackY < ey+eh &&
		attackY+attackH > ey
}

// ApplyDamageToEnemy applies damage and knockback to enemy
func (cs *CombatSystem) ApplyDamageToEnemy(enemy *entity.EnemyInstance, damage int, playerX float64) {
	enemy.TakeDamage(damage)
	
	// Apply knockback
	knockbackDir := 1.0
	if enemy.X < playerX {
		knockbackDir = -1.0
	}
	
	enemy.VelX = knockbackDir * 5.0
	enemy.VelY = -3.0
}

// CheckPlayerEnemyCollision checks if player touched enemy
func (cs *CombatSystem) CheckPlayerEnemyCollision(playerX, playerY, playerW, playerH float64, enemy *entity.EnemyInstance) bool {
	if cs.invulnerableFrames > 0 {
		return false // Player is invulnerable
	}
	
	ex, ey, ew, eh := enemy.GetBounds()
	
	return playerX < ex+ew &&
		playerX+playerW > ex &&
		playerY < ey+eh &&
		playerY+playerH > ey
}

// ApplyDamageToPlayer applies damage and knockback to player
func (cs *CombatSystem) ApplyDamageToPlayer(player *Player, damage int, enemyX float64) {
	if cs.invulnerableFrames > 0 {
		return // Player is invulnerable
	}
	
	player.Health -= damage
	if player.Health < 0 {
		player.Health = 0
	}
	
	// Apply knockback
	knockbackDir := 1.0
	if player.X < enemyX {
		knockbackDir = -1.0
	}
	
	cs.knockbackVelX = knockbackDir * 8.0
	cs.knockbackVelY = -5.0
	
	// Invulnerability frames
	cs.invulnerableFrames = 60 // 1 second of invulnerability
}

// GetKnockback returns current knockback velocity
func (cs *CombatSystem) GetKnockback() (float64, float64) {
	vx := cs.knockbackVelX
	vy := cs.knockbackVelY
	
	// Decay knockback
	cs.knockbackVelX *= 0.8
	cs.knockbackVelY *= 0.8
	
	if math.Abs(cs.knockbackVelX) < 0.1 {
		cs.knockbackVelX = 0
	}
	if math.Abs(cs.knockbackVelY) < 0.1 {
		cs.knockbackVelY = 0
	}
	
	return vx, vy
}

// IsInvulnerable returns if player is invulnerable
func (cs *CombatSystem) IsInvulnerable() bool {
	return cs.invulnerableFrames > 0
}

// GetInvulnerableFrames returns remaining invulnerable frames
func (cs *CombatSystem) GetInvulnerableFrames() int {
	return cs.invulnerableFrames
}
