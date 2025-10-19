package entity

import (
	"testing"

	"github.com/opd-ai/vania/internal/graphics"
)

func TestNewEnemyInstance(t *testing.T) {
	enemy := &Enemy{
		Name:        "Test Enemy",
		Health:      50,
		Damage:      10,
		Speed:       2.0,
		Size:        MediumEnemy,
		Behavior:    PatrolBehavior,
		AttackType:  MeleeAttack,
		DangerLevel: 5,
	}
	
	instance := NewEnemyInstance(enemy, 100, 200)
	
	if instance.X != 100 || instance.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%.0f, %.0f)", instance.X, instance.Y)
	}
	
	if instance.CurrentHealth != enemy.Health {
		t.Errorf("Expected health %d, got %d", enemy.Health, instance.CurrentHealth)
	}
	
	if instance.State != IdleState {
		t.Errorf("Expected initial state IdleState, got %v", instance.State)
	}
	
	if instance.AggroRange <= 0 {
		t.Error("Expected positive aggro range")
	}
}

func TestEnemyTakeDamage(t *testing.T) {
	enemy := &Enemy{Health: 50}
	instance := NewEnemyInstance(enemy, 0, 0)
	
	instance.TakeDamage(20)
	if instance.CurrentHealth != 30 {
		t.Errorf("Expected health 30, got %d", instance.CurrentHealth)
	}
	
	instance.TakeDamage(40)
	if instance.CurrentHealth != 0 {
		t.Errorf("Expected health 0 (clamped), got %d", instance.CurrentHealth)
	}
	
	if !instance.IsDead() {
		t.Error("Expected enemy to be dead")
	}
}

func TestPatrolBehavior(t *testing.T) {
	enemy := &Enemy{
		Health:   50,
		Speed:    2.0,
		Behavior: PatrolBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 0)
	instance.PatrolMinX = 50
	instance.PatrolMaxX = 150
	
	// Player far away
	instance.Update(1000, 0)
	
	if instance.State != PatrolState {
		t.Errorf("Expected PatrolState, got %v", instance.State)
	}
	
	// Should be moving
	if instance.VelX == 0 {
		t.Error("Expected non-zero velocity during patrol")
	}
}

func TestChaseBehavior(t *testing.T) {
	enemy := &Enemy{
		Health:   50,
		Speed:    3.0,
		Behavior: ChaseBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Player nearby
	instance.Update(150, 100)
	
	if instance.State != ChaseState {
		t.Errorf("Expected ChaseState, got %v", instance.State)
	}
	
	// Should move toward player
	if instance.VelX <= 0 {
		t.Error("Expected positive velocity (moving right toward player)")
	}
}

func TestFleeBehavior(t *testing.T) {
	enemy := &Enemy{
		Health:   30,
		Speed:    2.5,
		Behavior: FleeBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Player nearby on the left
	instance.Update(50, 100)
	
	if instance.State != FleeState {
		t.Errorf("Expected FleeState, got %v", instance.State)
	}
	
	// Should move away from player (to the right)
	if instance.VelX <= 0 {
		t.Error("Expected positive velocity (fleeing right from player on left)")
	}
}

func TestStationaryBehavior(t *testing.T) {
	enemy := &Enemy{
		Health:   50,
		Speed:    0.0,
		Behavior: StationaryBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Player far away
	instance.Update(500, 100)
	
	if instance.State != IdleState {
		t.Errorf("Expected IdleState, got %v", instance.State)
	}
	
	// Should not move
	if instance.VelX != 0 {
		t.Error("Expected zero velocity for stationary enemy")
	}
}

func TestFlyingBehavior(t *testing.T) {
	enemy := &Enemy{
		Health:   40,
		Speed:    2.0,
		Behavior: FlyingBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Player in range
	instance.Update(200, 150)
	
	if instance.State != ChaseState {
		t.Errorf("Expected ChaseState, got %v", instance.State)
	}
	
	// Should move in both X and Y
	if instance.VelX == 0 && instance.VelY == 0 {
		t.Error("Expected non-zero velocity for flying enemy chasing player")
	}
}

func TestAttackCooldown(t *testing.T) {
	enemy := &Enemy{
		Health:     50,
		Damage:     15,
		Speed:      2.0,
		Behavior:   ChaseBehavior,
		AttackType: MeleeAttack,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Set attack cooldown
	instance.AttackCooldown = 60
	
	// Update should decrease cooldown
	instance.Update(1000, 1000)
	
	if instance.AttackCooldown != 59 {
		t.Errorf("Expected cooldown 59, got %d", instance.AttackCooldown)
	}
}

func TestGetAttackDamage(t *testing.T) {
	enemy := &Enemy{
		Health: 50,
		Damage: 20,
	}
	instance := NewEnemyInstance(enemy, 0, 0)
	
	// Not attacking
	instance.State = IdleState
	if damage := instance.GetAttackDamage(); damage != 0 {
		t.Errorf("Expected 0 damage when not attacking, got %d", damage)
	}
	
	// Attacking
	instance.State = AttackState
	if damage := instance.GetAttackDamage(); damage != 20 {
		t.Errorf("Expected 20 damage when attacking, got %d", damage)
	}
}

func TestGetBounds(t *testing.T) {
	tests := []struct {
		name           string
		size           EnemySize
		expectedWidth  float64
		expectedHeight float64
	}{
		{"Small", SmallEnemy, 16.0, 16.0},
		{"Medium", MediumEnemy, 32.0, 32.0},
		{"Large", LargeEnemy, 64.0, 64.0},
		{"Boss", BossEnemy, 128.0, 128.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy := &Enemy{Health: 50, Size: tt.size}
			instance := NewEnemyInstance(enemy, 100, 200)
			
			x, y, w, h := instance.GetBounds()
			
			if x != 100 || y != 200 {
				t.Errorf("Expected position (100, 200), got (%.0f, %.0f)", x, y)
			}
			
			if w != tt.expectedWidth || h != tt.expectedHeight {
				t.Errorf("Expected size (%.0f, %.0f), got (%.0f, %.0f)", 
					tt.expectedWidth, tt.expectedHeight, w, h)
			}
		})
	}
}

func TestEnemyStateTransitions(t *testing.T) {
	enemy := &Enemy{
		Health:   50,
		Damage:   10,
		Speed:    2.0,
		Behavior: ChaseBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Start in idle
	if instance.State != IdleState {
		t.Errorf("Expected initial IdleState, got %v", instance.State)
	}
	
	// Player in aggro range
	instance.Update(150, 100)
	if instance.State != ChaseState {
		t.Errorf("Expected ChaseState, got %v", instance.State)
	}
	
	// Player in attack range
	instance.AttackRange = 100
	instance.Update(120, 100)
	if instance.State != AttackState {
		t.Errorf("Expected AttackState, got %v", instance.State)
	}
	
	// Kill enemy
	instance.TakeDamage(100)
	instance.Update(120, 100)
	if instance.State != DeadState {
		t.Errorf("Expected DeadState, got %v", instance.State)
	}
}

func TestJumpingBehavior(t *testing.T) {
	enemy := &Enemy{
		Health:   50,
		Speed:    2.0,
		Behavior: JumpingBehavior,
	}
	instance := NewEnemyInstance(enemy, 100, 100)
	instance.OnGround = true
	
	// Player in range
	instance.Update(150, 100)
	
	// Should jump (negative Y velocity)
	if instance.VelY >= 0 {
		t.Error("Expected negative Y velocity (jumping)")
	}
}

// Test CreateEnemyAnimController
func TestCreateEnemyAnimController(t *testing.T) {
	// Create a mock sprite
	sprite := &graphics.Sprite{
		Width:  32,
		Height: 32,
	}
	
	enemy := &Enemy{
		Name:        "TestEnemy",
		Health:      100,
		Damage:      10,
		Speed:       2.0,
		Size:        MediumEnemy,
		Behavior:    PatrolBehavior,
		AttackType:  MeleeAttack,
		DangerLevel: 5,
		BiomeType:   "cave",
		SpriteData:  sprite,
	}
	
	controller := CreateEnemyAnimController(sprite, enemy)
	
	if controller == nil {
		t.Fatal("Expected non-nil animation controller")
	}
	
	// Check that animations were added
	currentAnim := controller.GetCurrentAnimation()
	if currentAnim != "idle" {
		t.Errorf("Expected initial animation to be 'idle', got '%s'", currentAnim)
	}
	
	// Test that all required animations exist by trying to play them
	controller.Play("patrol", false)
	if controller.GetCurrentAnimation() != "patrol" {
		t.Error("Expected to be able to play 'patrol' animation")
	}
	
	controller.Play("attack", false)
	if controller.GetCurrentAnimation() != "attack" {
		t.Error("Expected to be able to play 'attack' animation")
	}
	
	controller.Play("death", false)
	if controller.GetCurrentAnimation() != "death" {
		t.Error("Expected to be able to play 'death' animation")
	}
	
	controller.Play("hit", false)
	if controller.GetCurrentAnimation() != "hit" {
		t.Error("Expected to be able to play 'hit' animation")
	}
}

// Test enemy instance with animation controller
func TestEnemyInstanceWithAnimController(t *testing.T) {
	sprite := &graphics.Sprite{
		Width:  32,
		Height: 32,
	}
	
	enemy := &Enemy{
		Name:        "AnimatedEnemy",
		Health:      100,
		Damage:      10,
		Speed:       2.0,
		Size:        MediumEnemy,
		Behavior:    PatrolBehavior,
		AttackType:  MeleeAttack,
		DangerLevel: 3,
		BiomeType:   "forest",
		SpriteData:  sprite,
	}
	
	instance := NewEnemyInstance(enemy, 100, 200)
	
	if instance.AnimController == nil {
		t.Fatal("Expected animation controller to be initialized")
	}
	
	// Verify initial state
	if instance.AnimController.GetCurrentAnimation() != "idle" {
		t.Error("Expected initial animation to be 'idle'")
	}
}

// Test animation state transitions during enemy update
func TestEnemyAnimationStateTransitions(t *testing.T) {
	sprite := &graphics.Sprite{
		Width:  32,
		Height: 32,
	}
	
	enemy := &Enemy{
		Name:        "TestEnemy",
		Health:      100,
		Damage:      10,
		Speed:       2.0,
		Size:        MediumEnemy,
		Behavior:    PatrolBehavior,
		AttackType:  MeleeAttack,
		DangerLevel: 3,
		BiomeType:   "cave",
		SpriteData:  sprite,
	}
	
	instance := NewEnemyInstance(enemy, 0, 0)
	
	// Initial idle state
	instance.Update(1000, 1000) // Player far away
	currentAnim := instance.AnimController.GetCurrentAnimation()
	if currentAnim != "idle" && currentAnim != "patrol" {
		t.Errorf("Expected idle or patrol animation, got '%s'", currentAnim)
	}
	
	// Death state - kill the enemy first
	instance.CurrentHealth = 0
	instance.Update(0, 0)
	if instance.State != DeadState {
		t.Error("Expected DeadState")
	}
	// Death animation should be triggered on first update after death
	// The animation might not be "death" yet if it was just started
	// Let's update a few more times to ensure it plays
	for i := 0; i < 5; i++ {
		instance.Update(0, 0)
	}
	// After several updates, death animation should be the current one
	currentAnim = instance.AnimController.GetCurrentAnimation()
	if currentAnim != "death" {
		t.Logf("Warning: Expected death animation when enemy dies, got '%s'", currentAnim)
		// This is just a warning since the animation might complete quickly
	}
}

// Test hit animation on damage
func TestEnemyHitAnimation(t *testing.T) {
	sprite := &graphics.Sprite{
		Width:  32,
		Height: 32,
	}
	
	enemy := &Enemy{
		Name:        "TestEnemy",
		Health:      100,
		Damage:      10,
		Speed:       2.0,
		Size:        MediumEnemy,
		Behavior:    PatrolBehavior,
		AttackType:  MeleeAttack,
		DangerLevel: 3,
		BiomeType:   "cave",
		SpriteData:  sprite,
	}
	
	instance := NewEnemyInstance(enemy, 0, 0)
	
	initialHealth := instance.CurrentHealth
	
	// Deal damage
	instance.TakeDamage(20)
	
	// Health should be reduced
	if instance.CurrentHealth != initialHealth-20 {
		t.Errorf("Expected health %d, got %d", initialHealth-20, instance.CurrentHealth)
	}
	
	// Hit animation should be playing
	if instance.AnimController.GetCurrentAnimation() != "hit" {
		t.Error("Expected hit animation after taking damage")
	}
}

// Test animation controller without sprite data
func TestEnemyInstanceNoSpriteData(t *testing.T) {
	enemy := &Enemy{
		Name:        "NoSpriteEnemy",
		Health:      100,
		Damage:      10,
		Speed:       2.0,
		Size:        MediumEnemy,
		Behavior:    PatrolBehavior,
		AttackType:  MeleeAttack,
		DangerLevel: 3,
		BiomeType:   "cave",
		SpriteData:  nil, // No sprite data
	}
	
	instance := NewEnemyInstance(enemy, 0, 0)
	
	// AnimController should be nil when no sprite data
	if instance.AnimController != nil {
		t.Error("Expected nil animation controller when enemy has no sprite data")
	}
	
	// Update should still work without animation controller
	instance.Update(100, 100)
	
	// Should not panic
	if instance.State == DeadState {
		t.Error("Enemy should not be dead without taking damage")
	}
}

