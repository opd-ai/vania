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

// Tests for parry system

func TestPlayerParry(t *testing.T) {
	cs := NewCombatSystem()

	// First parry should succeed
	if !cs.PlayerParry() {
		t.Error("Expected first parry to succeed")
	}

	if !cs.IsPlayerParrying() {
		t.Error("Expected player to be parrying")
	}

	// Second parry should fail while still parrying
	if cs.PlayerParry() {
		t.Error("Expected second parry to fail while still parrying")
	}

	// Complete the parry window
	for i := 0; i <= ParryWindowFrames; i++ {
		cs.Update()
	}

	// Now on cooldown, should still fail
	if cs.PlayerParry() {
		t.Error("Expected parry to fail during cooldown")
	}
}

func TestParryWindowFrames(t *testing.T) {
	cs := NewCombatSystem()

	cs.PlayerParry()

	// Within parry window
	for i := 0; i < ParryWindowFrames; i++ {
		if !cs.IsInParryWindow() {
			t.Errorf("Expected to be in parry window at frame %d", i)
		}
		cs.Update()
	}

	// After parry window
	cs.Update()
	if cs.IsInParryWindow() {
		t.Error("Expected to be outside parry window after duration")
	}

	// Parry should end
	if cs.IsPlayerParrying() {
		t.Error("Expected parry to end after window")
	}
}

func TestParryCooldown(t *testing.T) {
	cs := NewCombatSystem()

	cs.PlayerParry()

	// Complete parry window
	for i := 0; i <= ParryWindowFrames; i++ {
		cs.Update()
	}

	// Should be on cooldown
	if cs.CanParry() {
		t.Error("Expected parry to be on cooldown")
	}

	// Cooldown should expire
	for i := 0; i < ParryCooldownFrames; i++ {
		cs.Update()
	}

	if !cs.CanParry() {
		t.Error("Expected parry cooldown to have expired")
	}
}

func TestParryBlocksAttack(t *testing.T) {
	cs := NewCombatSystem()

	player := &Player{
		Health:    100,
		MaxHealth: 100,
		X:         100,
		Y:         100,
	}

	cs.PlayerParry()

	// Advance to within parry window
	cs.Update()
	cs.Update()

	// Attack during parry window should be blocked
	cs.ApplyDamageToPlayer(player, 25, 150)

	if player.Health != 100 {
		t.Error("Expected parry to block damage")
	}

	if cs.lastParrySucceeded != true {
		t.Error("Expected lastParrySucceeded to be true")
	}

	// Should not be invulnerable (parry doesn't grant invuln)
	if cs.invulnerableFrames > 0 {
		t.Error("Expected no invulnerability from parry")
	}
}

func TestParryOutsideWindowTakesDamage(t *testing.T) {
	cs := NewCombatSystem()

	player := &Player{
		Health:    100,
		MaxHealth: 100,
		X:         100,
		Y:         100,
	}

	cs.PlayerParry()

	// Advance past parry window
	for i := 0; i <= ParryWindowFrames+5; i++ {
		cs.Update()
	}

	// Attack after parry window should deal damage
	cs.ApplyDamageToPlayer(player, 25, 150)

	if player.Health >= 100 {
		t.Error("Expected damage to go through after parry window")
	}
}

func TestCannotParryWhileAttacking(t *testing.T) {
	cs := NewCombatSystem()

	cs.PlayerAttack()

	if cs.CanParry() {
		t.Error("Expected cannot parry while attacking")
	}

	if cs.PlayerParry() {
		t.Error("Expected parry to fail while attacking")
	}
}

func TestCannotParryWhileStaggered(t *testing.T) {
	cs := NewCombatSystem()

	cs.playerStaggered = true
	cs.playerStaggerTime = 10

	if cs.CanParry() {
		t.Error("Expected cannot parry while staggered")
	}

	if cs.PlayerParry() {
		t.Error("Expected parry to fail while staggered")
	}
}

// Tests for stagger system

func TestStaggerPreventsActions(t *testing.T) {
	cs := NewCombatSystem()

	player := &Player{Health: 100, MaxHealth: 100, X: 100, Y: 100}

	// Take damage to trigger stagger
	cs.ApplyDamageToPlayer(player, 10, 150)

	if !cs.IsPlayerStaggered() {
		t.Error("Expected player to be staggered after damage")
	}

	// Cannot attack while staggered
	if cs.CanAttack() {
		t.Error("Expected cannot attack while staggered")
	}

	if cs.PlayerAttack() {
		t.Error("Expected attack to fail while staggered")
	}

	// Cannot parry while staggered
	if cs.CanParry() {
		t.Error("Expected cannot parry while staggered")
	}
}

func TestStaggerDuration(t *testing.T) {
	cs := NewCombatSystem()

	player := &Player{Health: 100, MaxHealth: 100, X: 100, Y: 100}

	cs.ApplyDamageToPlayer(player, 10, 150)

	initialStagger := cs.GetStaggerFrames()
	if initialStagger != StaggerDurationFrames {
		t.Errorf("Expected stagger duration %d, got %d", StaggerDurationFrames, initialStagger)
	}

	// Stagger should decay
	cs.Update()
	if cs.GetStaggerFrames() >= initialStagger {
		t.Error("Expected stagger time to decrease")
	}

	// Stagger should eventually end
	for i := 0; i < StaggerDurationFrames; i++ {
		cs.Update()
	}

	if cs.IsPlayerStaggered() {
		t.Error("Expected stagger to end after duration")
	}

	// Should be able to attack after stagger ends
	if !cs.CanAttack() {
		t.Error("Expected to be able to attack after stagger ends")
	}
}

// Tests for damage numbers

func TestAddDamageNumber(t *testing.T) {
	cs := NewCombatSystem()

	cs.AddDamageNumber(25, 100.0, 100.0, false)

	numbers := cs.GetDamageNumbers()
	if len(numbers) != 1 {
		t.Errorf("Expected 1 damage number, got %d", len(numbers))
	}

	if numbers[0].Value != 25 {
		t.Errorf("Expected damage value 25, got %d", numbers[0].Value)
	}

	if numbers[0].X != 100.0 || numbers[0].Y != 100.0 {
		t.Errorf("Expected position (100, 100), got (%.1f, %.1f)", numbers[0].X, numbers[0].Y)
	}

	if numbers[0].IsCrit {
		t.Error("Expected non-crit damage")
	}
}

func TestDamageNumberLifetime(t *testing.T) {
	cs := NewCombatSystem()

	cs.AddDamageNumber(10, 100.0, 100.0, false)

	initialLifetime := cs.GetDamageNumbers()[0].LifeTime

	// Update should decrease lifetime
	cs.Update()

	numbers := cs.GetDamageNumbers()
	if len(numbers) == 0 {
		t.Fatal("Expected damage number to still exist")
	}

	if numbers[0].LifeTime >= initialLifetime {
		t.Error("Expected lifetime to decrease")
	}
}

func TestDamageNumberExpiration(t *testing.T) {
	cs := NewCombatSystem()

	cs.AddDamageNumber(10, 100.0, 100.0, false)

	// Update for longer than lifetime
	for i := 0; i < 100; i++ {
		cs.Update()
	}

	numbers := cs.GetDamageNumbers()
	if len(numbers) != 0 {
		t.Error("Expected damage number to be removed after expiration")
	}
}

func TestMultipleDamageNumbers(t *testing.T) {
	cs := NewCombatSystem()

	cs.AddDamageNumber(10, 100.0, 100.0, false)
	cs.AddDamageNumber(20, 200.0, 200.0, true)
	cs.AddDamageNumber(30, 300.0, 300.0, false)

	numbers := cs.GetDamageNumbers()
	if len(numbers) != 3 {
		t.Errorf("Expected 3 damage numbers, got %d", len(numbers))
	}

	// Verify each is independent
	if numbers[0].Value == numbers[1].Value {
		t.Error("Expected different damage values")
	}

	if !numbers[1].IsCrit {
		t.Error("Expected second damage to be crit")
	}
}

func TestDamageNumberMovement(t *testing.T) {
	cs := NewCombatSystem()

	cs.AddDamageNumber(10, 100.0, 100.0, false)

	initialY := cs.GetDamageNumbers()[0].Y

	// Update should move damage number
	cs.Update()

	numbers := cs.GetDamageNumbers()
	if len(numbers) == 0 {
		t.Fatal("Expected damage number to still exist")
	}

	// Y should change (moving upward)
	if numbers[0].Y == initialY {
		t.Error("Expected damage number to move")
	}
}

func TestDamageNumberAddedOnEnemyHit(t *testing.T) {
	cs := NewCombatSystem()

	enemy := &entity.Enemy{Health: 50}
	instance := entity.NewEnemyInstance(enemy, 100, 100)

	cs.ApplyDamageToEnemy(instance, 15, 50)

	numbers := cs.GetDamageNumbers()
	if len(numbers) != 1 {
		t.Errorf("Expected 1 damage number from enemy hit, got %d", len(numbers))
	}

	if numbers[0].Value != 15 {
		t.Errorf("Expected damage value 15, got %d", numbers[0].Value)
	}
}

func TestDamageNumberAddedOnPlayerHit(t *testing.T) {
	cs := NewCombatSystem()

	player := &Player{Health: 100, MaxHealth: 100, X: 100, Y: 100}

	cs.ApplyDamageToPlayer(player, 20, 150)

	numbers := cs.GetDamageNumbers()
	if len(numbers) != 1 {
		t.Errorf("Expected 1 damage number from player hit, got %d", len(numbers))
	}

	if numbers[0].Value != 20 {
		t.Errorf("Expected damage value 20, got %d", numbers[0].Value)
	}
}

func TestParryDamageNumber(t *testing.T) {
	cs := NewCombatSystem()

	player := &Player{Health: 100, MaxHealth: 100, X: 100, Y: 100}

	cs.PlayerParry()
	cs.Update()

	// Attack during parry
	cs.ApplyDamageToPlayer(player, 25, 150)

	numbers := cs.GetDamageNumbers()
	if len(numbers) != 1 {
		t.Errorf("Expected 1 damage number for parry feedback, got %d", len(numbers))
	}

	// Parry damage number has value 0 (special marker)
	if numbers[0].Value != 0 {
		t.Errorf("Expected parry damage number value 0, got %d", numbers[0].Value)
	}
}
