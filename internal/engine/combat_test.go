package engine

import (
	"testing"
	
	"github.com/opd-ai/vania/internal/entity"
)

func TestNewCombatSystem(t *testing.T) {
	cs := NewCombatSystem()
	
	if cs == nil {
		t.Fatal("Expected non-nil combat system")
	}
	
	if cs.playerAttackCooldown != 0 {
		t.Errorf("Expected initial cooldown 0, got %d", cs.playerAttackCooldown)
	}
	
	if cs.playerAttacking {
		t.Error("Expected player not attacking initially")
	}
}

func TestPlayerAttack(t *testing.T) {
	cs := NewCombatSystem()
	
	// First attack should succeed
	if !cs.PlayerAttack() {
		t.Error("Expected first attack to succeed")
	}
	
	if !cs.IsPlayerAttacking() {
		t.Error("Expected player to be attacking")
	}
	
	// Second immediate attack should fail (cooldown)
	if cs.PlayerAttack() {
		t.Error("Expected second attack to fail due to cooldown")
	}
}

func TestAttackCooldown(t *testing.T) {
	cs := NewCombatSystem()
	
	cs.PlayerAttack()
	initialCooldown := cs.playerAttackCooldown
	
	if initialCooldown <= 0 {
		t.Error("Expected positive cooldown after attack")
	}
	
	// Update should decrease cooldown
	cs.Update()
	
	if cs.playerAttackCooldown >= initialCooldown {
		t.Error("Expected cooldown to decrease after update")
	}
}

func TestAttackFrames(t *testing.T) {
	cs := NewCombatSystem()
	
	cs.PlayerAttack()
	
	// Attack should last multiple frames
	for i := 0; i < 15; i++ {
		if !cs.IsPlayerAttacking() {
			t.Errorf("Expected to be attacking at frame %d", i)
		}
		cs.Update()
	}
	
	// After 15+ frames, attack should end
	cs.Update()
	if cs.IsPlayerAttacking() {
		t.Error("Expected attack to end after duration")
	}
}

func TestGetAttackHitbox(t *testing.T) {
	cs := NewCombatSystem()
	
	// No hitbox when not attacking
	x, y, w, h := cs.GetAttackHitbox(100, 100, 1.0)
	if w != 0 || h != 0 {
		t.Error("Expected no hitbox when not attacking")
	}
	
	// Hitbox appears during active frames
	cs.PlayerAttack()
	for i := 0; i < 3; i++ {
		cs.Update()
	}
	
	x, y, w, h = cs.GetAttackHitbox(100, 100, 1.0)
	if w <= 0 || h <= 0 {
		t.Error("Expected valid hitbox during active attack frames")
	}
	
	if x == 0 || y == 0 {
		t.Error("Expected non-zero hitbox position")
	}
	
	// Hitbox position depends on facing direction
	xRight, _, _, _ := cs.GetAttackHitbox(100, 100, 1.0)
	
	cs2 := NewCombatSystem()
	cs2.PlayerAttack()
	for i := 0; i < 3; i++ {
		cs2.Update()
	}
	xLeft, _, _, _ := cs2.GetAttackHitbox(100, 100, -1.0)
	
	if xRight == xLeft {
		t.Error("Expected different hitbox positions for different facing directions")
	}
}

func TestCheckEnemyHit(t *testing.T) {
	cs := NewCombatSystem()
	
	enemy := &entity.Enemy{
		Health: 50,
		Size:   entity.MediumEnemy,
	}
	instance := entity.NewEnemyInstance(enemy, 150, 100)
	
	// Hit detection - overlapping
	if !cs.CheckEnemyHit(140, 90, 40, 32, instance) {
		t.Error("Expected hit when hitboxes overlap")
	}
	
	// No hit - not overlapping
	if cs.CheckEnemyHit(10, 10, 40, 32, instance) {
		t.Error("Expected no hit when hitboxes don't overlap")
	}
	
	// No hitbox
	if cs.CheckEnemyHit(150, 100, 0, 0, instance) {
		t.Error("Expected no hit with zero-size hitbox")
	}
}

func TestApplyDamageToEnemy(t *testing.T) {
	cs := NewCombatSystem()
	
	enemy := &entity.Enemy{
		Health: 50,
	}
	instance := entity.NewEnemyInstance(enemy, 150, 100)
	
	initialHealth := instance.CurrentHealth
	
	cs.ApplyDamageToEnemy(instance, 20, 100)
	
	if instance.CurrentHealth >= initialHealth {
		t.Error("Expected health to decrease after damage")
	}
	
	// Check knockback applied
	if instance.VelX == 0 && instance.VelY == 0 {
		t.Error("Expected knockback velocity after damage")
	}
}

func TestCheckPlayerEnemyCollision(t *testing.T) {
	cs := NewCombatSystem()
	
	enemy := &entity.Enemy{
		Health: 50,
		Size:   entity.MediumEnemy,
	}
	instance := entity.NewEnemyInstance(enemy, 110, 100)
	
	// Collision detected
	if !cs.CheckPlayerEnemyCollision(100, 100, 32, 32, instance) {
		t.Error("Expected collision when player and enemy overlap")
	}
	
	// No collision
	if cs.CheckPlayerEnemyCollision(10, 10, 32, 32, instance) {
		t.Error("Expected no collision when far apart")
	}
	
	// No collision during invulnerability
	cs.invulnerableFrames = 60
	if cs.CheckPlayerEnemyCollision(110, 100, 32, 32, instance) {
		t.Error("Expected no collision during invulnerability")
	}
}

func TestApplyDamageToPlayer(t *testing.T) {
	cs := NewCombatSystem()
	
	player := &Player{
		Health:    100,
		MaxHealth: 100,
		X:         100,
		Y:         100,
	}
	
	cs.ApplyDamageToPlayer(player, 25, 150)
	
	if player.Health >= 100 {
		t.Error("Expected player health to decrease")
	}
	
	if player.Health != 75 {
		t.Errorf("Expected health 75, got %d", player.Health)
	}
	
	if !cs.IsInvulnerable() {
		t.Error("Expected player to be invulnerable after taking damage")
	}
	
	// Damage during invulnerability should be ignored
	cs.ApplyDamageToPlayer(player, 25, 150)
	if player.Health != 75 {
		t.Error("Expected no additional damage during invulnerability")
	}
}

func TestInvulnerabilityFrames(t *testing.T) {
	cs := NewCombatSystem()
	
	if cs.IsInvulnerable() {
		t.Error("Expected not invulnerable initially")
	}
	
	player := &Player{Health: 100, MaxHealth: 100}
	cs.ApplyDamageToPlayer(player, 10, 150)
	
	initialFrames := cs.GetInvulnerableFrames()
	if initialFrames <= 0 {
		t.Error("Expected positive invulnerable frames after damage")
	}
	
	// Update should decrease frames
	cs.Update()
	if cs.GetInvulnerableFrames() >= initialFrames {
		t.Error("Expected invulnerable frames to decrease")
	}
}

func TestKnockback(t *testing.T) {
	cs := NewCombatSystem()
	
	// Initial knockback should be zero
	vx, vy := cs.GetKnockback()
	if vx != 0 || vy != 0 {
		t.Error("Expected zero initial knockback")
	}
	
	// Apply damage creates knockback
	player := &Player{Health: 100, MaxHealth: 100, X: 100}
	cs.ApplyDamageToPlayer(player, 10, 150)
	
	vx, vy = cs.GetKnockback()
	if vx == 0 && vy == 0 {
		t.Error("Expected non-zero knockback after damage")
	}
	
	// Knockback eventually decays to zero
	for i := 0; i < 100; i++ {
		vx, vy = cs.GetKnockback()
		if vx == 0 && vy == 0 {
			return // Successfully decayed
		}
	}
	t.Error("Expected knockback to eventually decay to zero")
}

func TestPlayerHealthClamp(t *testing.T) {
	cs := NewCombatSystem()
	
	player := &Player{Health: 20, MaxHealth: 100}
	
	// Health shouldn't go below zero
	cs.ApplyDamageToPlayer(player, 50, 150)
	
	if player.Health < 0 {
		t.Errorf("Expected health clamped to 0, got %d", player.Health)
	}
}

func TestCombatSystemUpdate(t *testing.T) {
	cs := NewCombatSystem()
	
	// Set up various state
	cs.playerAttackCooldown = 10
	cs.playerAttacking = true
	cs.playerAttackFrame = 0
	cs.invulnerableFrames = 30
	
	cs.Update()
	
	// All counters should decrease
	if cs.playerAttackCooldown != 9 {
		t.Errorf("Expected cooldown 9, got %d", cs.playerAttackCooldown)
	}
	
	if cs.playerAttackFrame != 1 {
		t.Errorf("Expected attack frame 1, got %d", cs.playerAttackFrame)
	}
	
	if cs.invulnerableFrames != 29 {
		t.Errorf("Expected invuln frames 29, got %d", cs.invulnerableFrames)
	}
}
