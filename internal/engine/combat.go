// Package engine provides combat system implementation for managing
// player attacks, damage calculation, hit detection, and combat state.
package engine

import (
	"math"

	"github.com/opd-ai/vania/internal/entity"
)

const (
	// ProjectileSpeed is the base travel speed for player projectiles (pixels per frame).
	ProjectileSpeed = 10.0

	// ProjectileMaxRange is the maximum travel distance before a projectile expires (pixels).
	ProjectileMaxRange = 400.0

	// ProjectileCooldownFrames is the minimum delay between ranged attacks.
	ProjectileCooldownFrames = 25

	// ParryWindowFrames is the active frame window during which parry can deflect attacks.
	// At 60fps, 8 frames = ~133ms, requiring precise timing from the player.
	ParryWindowFrames = 8

	// ParryCooldownFrames is the cooldown after a parry attempt before another can be initiated.
	ParryCooldownFrames = 30

	// StaggerDurationFrames is how long an entity remains staggered after being hit.
	// During stagger, the entity cannot attack or use abilities.
	StaggerDurationFrames = 20
)

// DamageNumber represents floating damage text
type DamageNumber struct {
	Value    int
	X, Y     float64
	VelY     float64
	LifeTime int
	IsCrit   bool
}

// Projectile represents a ranged attack projectile in flight.
// Damage falls off linearly with distance traveled.
type Projectile struct {
	X, Y         float64
	VelX, VelY   float64
	Damage       int
	DistTraveled float64
	Active       bool
}

// CombatSystem manages all combat interactions
type CombatSystem struct {
	playerAttackCooldown int
	playerAttacking      bool
	playerAttackFrame    int
	knockbackVelX        float64
	knockbackVelY        float64
	invulnerableFrames   int

	// Ranged attack
	rangedCooldown int
	projectiles    []Projectile

	// Parry system
	playerParrying     bool
	parryFrame         int
	parryCooldown      int
	lastParrySucceeded bool

	// Stagger system
	playerStaggered   bool
	playerStaggerTime int

	// Damage numbers for visual feedback
	damageNumbers []DamageNumber
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
		rangedCooldown:       0,
		projectiles:          make([]Projectile, 0),
		playerParrying:       false,
		parryFrame:           0,
		parryCooldown:        0,
		lastParrySucceeded:   false,
		playerStaggered:      false,
		playerStaggerTime:    0,
		damageNumbers:        make([]DamageNumber, 0),
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

	// Update parry state
	if cs.parryCooldown > 0 {
		cs.parryCooldown--
	}

	if cs.playerParrying {
		cs.parryFrame++
		if cs.parryFrame > ParryWindowFrames {
			cs.playerParrying = false
			cs.parryFrame = 0
			cs.parryCooldown = ParryCooldownFrames
		}
	}

	// Update stagger state
	if cs.playerStaggered {
		cs.playerStaggerTime--
		if cs.playerStaggerTime <= 0 {
			cs.playerStaggered = false
		}
	}

	// Update ranged cooldown
	if cs.rangedCooldown > 0 {
		cs.rangedCooldown--
	}

	// Update projectiles
	for i := len(cs.projectiles) - 1; i >= 0; i-- {
		p := &cs.projectiles[i]
		if !p.Active {
			cs.projectiles = append(cs.projectiles[:i], cs.projectiles[i+1:]...)
			continue
		}
		speed := math.Sqrt(p.VelX*p.VelX + p.VelY*p.VelY)
		p.X += p.VelX
		p.Y += p.VelY
		p.DistTraveled += speed
		if p.DistTraveled >= ProjectileMaxRange {
			p.Active = false
		}
	}

	// Update damage numbers
	for i := len(cs.damageNumbers) - 1; i >= 0; i-- {
		cs.damageNumbers[i].Y -= cs.damageNumbers[i].VelY
		cs.damageNumbers[i].VelY *= 0.95 // Decelerate
		cs.damageNumbers[i].LifeTime--

		if cs.damageNumbers[i].LifeTime <= 0 {
			// Remove expired damage number
			cs.damageNumbers = append(cs.damageNumbers[:i], cs.damageNumbers[i+1:]...)
		}
	}
}

// PlayerAttack initiates a player attack
func (cs *CombatSystem) PlayerAttack() bool {
	if cs.playerAttackCooldown <= 0 && !cs.playerStaggered {
		cs.playerAttacking = true
		cs.playerAttackFrame = 0
		cs.playerAttackCooldown = 20 // 20 frames between attacks
		return true
	}
	return false
}

// PlayerParry initiates a parry attempt
func (cs *CombatSystem) PlayerParry() bool {
	if cs.parryCooldown <= 0 && !cs.playerStaggered && !cs.playerAttacking && !cs.playerParrying {
		cs.playerParrying = true
		cs.parryFrame = 0
		cs.lastParrySucceeded = false
		return true
	}
	return false
}

// CanAttack returns true if player can attack (not on cooldown or staggered)
func (cs *CombatSystem) CanAttack() bool {
	return cs.playerAttackCooldown <= 0 && !cs.playerStaggered
}

// CanParry returns true if player can parry (not on cooldown, staggered, attacking, or already parrying)
func (cs *CombatSystem) CanParry() bool {
	return cs.parryCooldown <= 0 && !cs.playerStaggered && !cs.playerAttacking && !cs.playerParrying
}

// IsPlayerParrying returns if player is currently in parry frames
func (cs *CombatSystem) IsPlayerParrying() bool {
	return cs.playerParrying
}

// IsInParryWindow returns if player is within the active parry window
func (cs *CombatSystem) IsInParryWindow() bool {
	return cs.playerParrying && cs.parryFrame <= ParryWindowFrames
}

// IsPlayerStaggered returns if player is currently staggered
func (cs *CombatSystem) IsPlayerStaggered() bool {
	return cs.playerStaggered
}

// GetStaggerFrames returns remaining stagger frames
func (cs *CombatSystem) GetStaggerFrames() int {
	return cs.playerStaggerTime
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

	// Spawn damage number
	cs.AddDamageNumber(damage, enemy.X, enemy.Y-10, false)
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

	// Check for successful parry
	if cs.IsInParryWindow() {
		cs.lastParrySucceeded = true
		cs.playerParrying = false
		cs.parryFrame = 0
		cs.parryCooldown = ParryCooldownFrames / 2 // Shorter cooldown on successful parry

		// Spawn damage number showing "PARRY!"
		cs.AddDamageNumber(0, player.X, player.Y-20, false)
		return
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

	// Apply stagger
	cs.playerStaggered = true
	cs.playerStaggerTime = StaggerDurationFrames

	// Invulnerability frames
	cs.invulnerableFrames = 60 // 1 second of invulnerability

	// Spawn damage number
	cs.AddDamageNumber(damage, player.X, player.Y-10, false)
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

// AddDamageNumber adds a floating damage number for visual feedback
func (cs *CombatSystem) AddDamageNumber(damage int, x, y float64, isCrit bool) {
	cs.damageNumbers = append(cs.damageNumbers, DamageNumber{
		Value:    damage,
		X:        x,
		Y:        y,
		VelY:     2.0, // Initial upward velocity
		LifeTime: 60,  // 1 second at 60fps
		IsCrit:   isCrit,
	})
}

// GetDamageNumbers returns all active damage numbers for rendering
func (cs *CombatSystem) GetDamageNumbers() []DamageNumber {
	return cs.damageNumbers
}

// LastParrySucceeded returns if the last parry successfully deflected an attack
func (cs *CombatSystem) LastParrySucceeded() bool {
	return cs.lastParrySucceeded
}

// ClearLastParrySuccess clears the last parry success flag
func (cs *CombatSystem) ClearLastParrySuccess() {
	cs.lastParrySucceeded = false
}

// PlayerRangedAttack spawns a projectile in the player's facing direction.
// Damage falls off linearly: full damage at range 0, zero damage at ProjectileMaxRange.
// Returns false if the ranged attack is on cooldown or the player is staggered.
func (cs *CombatSystem) PlayerRangedAttack(playerX, playerY, facingDir float64, baseDamage int) bool {
	if cs.rangedCooldown > 0 || cs.playerStaggered {
		return false
	}
	cs.projectiles = append(cs.projectiles, Projectile{
		X:            playerX + 16, // centre of player sprite
		Y:            playerY + 12,
		VelX:         ProjectileSpeed * facingDir,
		VelY:         0,
		Damage:       baseDamage,
		DistTraveled: 0,
		Active:       true,
	})
	cs.rangedCooldown = ProjectileCooldownFrames
	return true
}

// CanRangedAttack returns true when the ranged attack cooldown has elapsed
// and the player is not staggered.
func (cs *CombatSystem) CanRangedAttack() bool {
	return cs.rangedCooldown <= 0 && !cs.playerStaggered
}

// GetProjectiles returns the slice of currently active projectiles for rendering
// and external collision checks.
func (cs *CombatSystem) GetProjectiles() []Projectile {
	return cs.projectiles
}

// CheckProjectileEnemyHit tests every active projectile against the given enemy.
// On first hit the projectile is deactivated and damage (with distance falloff)
// is applied to the enemy.  Returns the damage dealt, or 0 if no hit occurred.
func (cs *CombatSystem) CheckProjectileEnemyHit(enemy *entity.EnemyInstance) int {
	ex, ey, ew, eh := enemy.GetBounds()
	for i := range cs.projectiles {
		p := &cs.projectiles[i]
		if !p.Active {
			continue
		}
		// Simple point-in-AABB check (projectile centre vs enemy bounds)
		if p.X >= ex && p.X <= ex+ew && p.Y >= ey && p.Y <= ey+eh {
			// Linear damage falloff: full damage at distance 0, 0 at max range
			falloff := 1.0 - (p.DistTraveled / ProjectileMaxRange)
			if falloff < 0 {
				falloff = 0
			}
			damage := int(math.Round(float64(p.Damage) * falloff))
			if damage < 1 {
				damage = 1
			}
			p.Active = false
			cs.ApplyDamageToEnemy(enemy, damage, p.X-p.VelX)
			return damage
		}
	}
	return 0
}
