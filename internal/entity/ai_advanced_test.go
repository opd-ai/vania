package entity

import (
	"math"
	"testing"
)

// TestNewAIMemory verifies AI memory initialization
func TestNewAIMemory(t *testing.T) {
	mem := NewAIMemory()
	
	if mem == nil {
		t.Fatal("NewAIMemory returned nil")
	}
	
	if mem.PlayerSkillEstimate != 0.5 {
		t.Errorf("Expected initial skill estimate 0.5, got %f", mem.PlayerSkillEstimate)
	}
	
	if mem.RetreatThreshold != 0.3 {
		t.Errorf("Expected initial retreat threshold 0.3, got %f", mem.RetreatThreshold)
	}
	
	if mem.LearningRate != 0.05 {
		t.Errorf("Expected learning rate 0.05, got %f", mem.LearningRate)
	}
}

// TestAIMemoryUpdateMemory tests memory updates from player observations
func TestAIMemoryUpdateMemory(t *testing.T) {
	mem := NewAIMemory()
	
	// Initial update
	mem.UpdateMemory(100, 100, false, false, false)
	
	if len(mem.LastPlayerPositions) != 1 {
		t.Errorf("Expected 1 position recorded, got %d", len(mem.LastPlayerPositions))
	}
	
	if !mem.KnowsPlayerPosition {
		t.Error("Memory should know player position after update")
	}
	
	// Update with jump
	for i := 0; i < 10; i++ {
		mem.UpdateMemory(100+float64(i), 100-float64(i)*5, true, false, false)
	}
	
	if mem.JumpFrequency <= 0 {
		t.Error("Jump frequency should increase after jumps")
	}
	
	if mem.TimesSeeingPlayer != 11 {
		t.Errorf("Expected 11 sightings, got %d", mem.TimesSeeingPlayer)
	}
}

// TestAIMemoryRecordCombatEvent tests combat event recording
func TestAIMemoryRecordCombatEvent(t *testing.T) {
	mem := NewAIMemory()
	
	// Record successful hit
	mem.RecordCombatEvent(true, false, 0, 50.0)
	
	if mem.SuccessfulHits != 1 {
		t.Errorf("Expected 1 successful hit, got %d", mem.SuccessfulHits)
	}
	
	// Check preferred distance learned
	if math.Abs(mem.PreferredAttackDistance-50.0) > 5.0 {
		t.Errorf("Expected preferred distance near 50.0, got %f", mem.PreferredAttackDistance)
	}
	
	// Record taking damage
	mem.RecordCombatEvent(false, true, 20, 0)
	
	if mem.DamageReceived != 20 {
		t.Errorf("Expected 20 damage received, got %d", mem.DamageReceived)
	}
}

// TestAIMemoryShouldRetreat tests retreat decision making
func TestAIMemoryShouldRetreat(t *testing.T) {
	mem := NewAIMemory()
	
	// Should not retreat at high health
	if mem.ShouldRetreat(0.8) {
		t.Error("Should not retreat at 80% health")
	}
	
	// Should retreat at low health
	if !mem.ShouldRetreat(0.2) {
		t.Error("Should retreat at 20% health")
	}
	
	// Test threshold adaptation
	mem.RecordCombatEvent(false, true, 60, 0)
	
	// Threshold should increase after taking heavy damage
	if mem.RetreatThreshold <= 0.3 {
		t.Error("Retreat threshold should increase after taking damage")
	}
}

// TestAIMemoryPredictPlayerPosition tests player position prediction
func TestAIMemoryPredictPlayerPosition(t *testing.T) {
	mem := NewAIMemory()
	
	// Need at least 2 positions for prediction
	mem.UpdateMemory(0, 0, false, false, false)
	mem.UpdateMemory(10, 5, false, false, false)
	
	predX, _ := mem.PredictPlayerPosition(1.0)
	
	// Should predict movement in same direction
	if predX < 10 {
		t.Errorf("Expected prediction > 10, got %f", predX)
	}
	
	// With no confidence, prediction should be close to last known
	if math.Abs(predX-10) > 15 {
		t.Errorf("Low confidence prediction should be near last position, got %f", predX)
	}
	
	// Build confidence
	for i := 0; i < 100; i++ {
		mem.UpdateMemory(float64(10+i), float64(5+i/2), false, false, false)
	}
	
	predX, _ = mem.PredictPlayerPosition(1.0)
	
	// With confidence, prediction should extrapolate
	lastPos := mem.LastPlayerPositions[len(mem.LastPlayerPositions)-1]
	if predX <= lastPos.X {
		t.Error("High confidence prediction should extrapolate forward")
	}
}

// TestAIMemoryGetTacticalState tests tactical state determination
func TestAIMemoryGetTacticalState(t *testing.T) {
	tests := []struct {
		name          string
		healthPercent float64
		hasAllies     bool
		distance      float64
		setupMemory   func(*AIMemory)
	}{
		{"Low health solo", 0.2, false, 100, nil},
		{"Low health with allies", 0.2, true, 100, nil},
		{"Winning fight", 0.8, false, 100, func(m *AIMemory) {
			m.SuccessfulHits = 10
			m.DamageReceived = 5
		}},
		{"Defensive", 0.4, false, 100, func(m *AIMemory) {
			m.DamageReceived = 50
		}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewAIMemory()
			
			// Set up memory state
			if tt.setupMemory != nil {
				tt.setupMemory(mem)
			}
			
			state := mem.GetTacticalState(tt.healthPercent, tt.hasAllies, tt.distance)
			
			// Verify the state makes sense for the situation
			switch tt.name {
			case "Low health solo":
				if state != TacticalRetreating {
					t.Errorf("Expected TacticalRetreating, got %v", state)
				}
			case "Low health with allies":
				if state != TacticalRegrouping {
					t.Errorf("Expected TacticalRegrouping, got %v", state)
				}
			case "Winning fight":
				if state != TacticalAggressive && state != TacticalNormal {
					t.Errorf("Expected TacticalAggressive or TacticalNormal, got %v", state)
				}
			case "Defensive":
				if state != TacticalDefensive && state != TacticalNormal {
					t.Errorf("Expected TacticalDefensive or TacticalNormal, got %v", state)
				}
			}
		})
	}
}

// TestAIMemoryRecordEvasion tests evasion recording
func TestAIMemoryRecordEvasion(t *testing.T) {
	mem := NewAIMemory()
	
	for i := 0; i < 5; i++ {
		mem.RecordEvasion()
	}
	
	if mem.AttacksEvaded != 5 {
		t.Errorf("Expected 5 evasions, got %d", mem.AttacksEvaded)
	}
}

// TestNewEnemyGroup tests enemy group creation
func TestNewEnemyGroup(t *testing.T) {
	group := NewEnemyGroup()
	
	if group == nil {
		t.Fatal("NewEnemyGroup returned nil")
	}
	
	if len(group.Members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(group.Members))
	}
	
	if group.Formation != NoFormation {
		t.Errorf("Expected NoFormation, got %v", group.Formation)
	}
}

// TestEnemyGroupAddRemoveMember tests member management
func TestEnemyGroupAddRemoveMember(t *testing.T) {
	group := NewEnemyGroup()
	
	enemy1 := &EnemyInstance{
		Enemy:         &Enemy{Health: 100},
		CurrentHealth: 100,
		State:         IdleState,
	}
	enemy2 := &EnemyInstance{
		Enemy:         &Enemy{Health: 150},
		CurrentHealth: 150,
		State:         IdleState,
	}
	
	// Add members
	group.AddMember(enemy1)
	group.AddMember(enemy2)
	
	if len(group.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(group.Members))
	}
	
	// Leader should be stronger enemy
	if group.Leader != enemy2 {
		t.Error("Leader should be enemy with most health")
	}
	
	// Remove member
	group.RemoveMember(enemy1)
	
	if len(group.Members) != 1 {
		t.Errorf("Expected 1 member after removal, got %d", len(group.Members))
	}
}

// TestEnemyGroupSelectFormation tests formation selection
func TestEnemyGroupSelectFormation(t *testing.T) {
	group := NewEnemyGroup()
	
	// Create enemies
	for i := 0; i < 5; i++ {
		enemy := &EnemyInstance{
			Enemy:         &Enemy{Health: 100},
			CurrentHealth: 100,
			State:         IdleState,
		}
		group.AddMember(enemy)
	}
	
	// Test engaging state with 4+ members
	group.GroupState = GroupEngaging
	group.selectFormation()
	
	if group.Formation != CircleFormation {
		t.Errorf("Expected CircleFormation for 5 engaging enemies, got %v", group.Formation)
	}
	
	// Test patrol state
	group.GroupState = GroupPatrol
	group.selectFormation()
	
	if group.Formation != VFormation {
		t.Errorf("Expected VFormation for 5 patrolling enemies, got %v", group.Formation)
	}
	
	// Test regrouping
	group.GroupState = GroupRegrouping
	group.selectFormation()
	
	if group.Formation != LineFormation {
		t.Errorf("Expected LineFormation for regrouping, got %v", group.Formation)
	}
}

// TestEnemyGroupApplyFormation tests formation positioning
func TestEnemyGroupApplyFormation(t *testing.T) {
	group := NewEnemyGroup()
	
	// Create 4 enemies
	for i := 0; i < 4; i++ {
		enemy := &EnemyInstance{
			Enemy:         &Enemy{Health: 100},
			CurrentHealth: 100,
			X:             0,
			Y:             0,
		}
		group.AddMember(enemy)
	}
	
	playerX, playerY := 100.0, 100.0
	
	// Test circle formation
	group.Formation = CircleFormation
	group.applyFormation(playerX, playerY)
	
	// Check that enemies are positioned around player
	for _, member := range group.Members {
		dx := member.FormationX - playerX
		dy := member.FormationY - playerY
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if math.Abs(dist-120.0) > 5.0 {
			t.Errorf("Expected enemy at ~120 units from player, got %f", dist)
		}
	}
	
	// Test line formation
	group.Formation = LineFormation
	group.applyFormation(playerX, playerY)
	
	// All enemies should have same Y
	firstY := group.Members[0].FormationY
	for i, member := range group.Members {
		if i > 0 && member.FormationY != firstY {
			t.Error("Line formation should have all enemies at same Y")
		}
	}
}

// TestEnemyGroupUpdateGroup tests group coordination update
func TestEnemyGroupUpdateGroup(t *testing.T) {
	group := NewEnemyGroup()
	
	// Create enemies with varying health
	for i := 0; i < 3; i++ {
		health := 100 - i*30 // 100, 70, 40
		enemy := &EnemyInstance{
			Enemy:         &Enemy{Health: 100, Speed: 2.0},
			CurrentHealth: health,
			State:         IdleState,
			X:             float64(i * 50),
			Y:             float64(i * 20),
		}
		group.AddMember(enemy)
	}
	
	// Update with player far away
	group.UpdateGroup(500, 500)
	
	if group.GroupState != GroupPatrol {
		t.Errorf("Expected GroupPatrol when not in combat, got %v", group.GroupState)
	}
	
	// Set one enemy to chase state
	group.Members[0].State = ChaseState
	group.UpdateGroup(500, 500)
	
	if group.GroupState != GroupEngaging {
		t.Errorf("Expected GroupEngaging when chasing, got %v", group.GroupState)
	}
	
	// Set low health on one enemy
	group.Members[2].CurrentHealth = 20
	group.UpdateGroup(500, 500)
	
	if group.GroupState != GroupRegrouping {
		t.Errorf("Expected GroupRegrouping with injured members, got %v", group.GroupState)
	}
}

// TestEnemyGroupRemoveDeadMembers tests dead member cleanup
func TestEnemyGroupRemoveDeadMembers(t *testing.T) {
	group := NewEnemyGroup()
	
	alive := &EnemyInstance{
		Enemy:         &Enemy{Health: 100},
		CurrentHealth: 50,
		State:         IdleState,
	}
	dead := &EnemyInstance{
		Enemy:         &Enemy{Health: 100},
		CurrentHealth: 0,
		State:         DeadState,
	}
	
	group.AddMember(alive)
	group.AddMember(dead)
	
	// Update should remove dead members
	group.UpdateGroup(100, 100)
	
	if len(group.Members) != 1 {
		t.Errorf("Expected 1 alive member, got %d", len(group.Members))
	}
	
	if group.Members[0] != alive {
		t.Error("Wrong member remained in group")
	}
}

// TestGetNearbyAllies tests ally detection
func TestGetNearbyAllies(t *testing.T) {
	enemies := []*EnemyInstance{
		{X: 0, Y: 0, State: IdleState},
		{X: 50, Y: 0, State: IdleState},
		{X: 500, Y: 0, State: IdleState},
		{X: 100, Y: 0, State: DeadState}, // Dead, should be ignored
	}
	
	allies := GetNearbyAllies(enemies, 0, 0, 200)
	
	// Should find enemy at (50,0) but not (500,0) or dead enemy
	if len(allies) != 1 {
		t.Errorf("Expected 1 ally within range, got %d", len(allies))
	}
	
	if allies[0] != enemies[1] {
		t.Error("Wrong ally detected")
	}
}

// TestEnemyInstanceWithAdvancedAI tests enemy with advanced AI fields
func TestEnemyInstanceWithAdvancedAI(t *testing.T) {
	enemy := &Enemy{
		Name:     "Test Enemy",
		Health:   100,
		Damage:   10,
		Speed:    2.0,
		Behavior: ChaseBehavior,
	}
	
	instance := NewEnemyInstance(enemy, 100, 100)
	
	if instance.Memory == nil {
		t.Error("Enemy instance should have AI memory")
	}
	
	if instance.TacticalState != TacticalNormal {
		t.Errorf("Expected TacticalNormal, got %v", instance.TacticalState)
	}
	
	// Update should use memory
	instance.Update(150, 100)
	
	if !instance.Memory.KnowsPlayerPosition {
		t.Error("Memory should know player position after update")
	}
}

// TestEnemyInstanceTakeDamageWithMemory tests damage recording
func TestEnemyInstanceTakeDamageWithMemory(t *testing.T) {
	enemy := &Enemy{Health: 100}
	instance := NewEnemyInstance(enemy, 0, 0)
	
	instance.TakeDamage(20)
	
	if instance.Memory.DamageReceived != 20 {
		t.Errorf("Expected 20 damage in memory, got %d", instance.Memory.DamageReceived)
	}
	
	if instance.CurrentHealth != 80 {
		t.Errorf("Expected 80 health, got %d", instance.CurrentHealth)
	}
}

// TestEnemyInstanceRecordSuccessfulHit tests hit recording
func TestEnemyInstanceRecordSuccessfulHit(t *testing.T) {
	enemy := &Enemy{Health: 100}
	instance := NewEnemyInstance(enemy, 0, 0)
	
	instance.RecordSuccessfulHit(45.0)
	
	if instance.Memory.SuccessfulHits != 1 {
		t.Errorf("Expected 1 successful hit, got %d", instance.Memory.SuccessfulHits)
	}
}

// TestTacticalStateTransitions tests tactical AI state changes
func TestTacticalStateTransitions(t *testing.T) {
	enemy := &Enemy{
		Health:   100,
		Damage:   10,
		Speed:    2.0,
		Behavior: ChaseBehavior,
	}
	
	instance := NewEnemyInstance(enemy, 100, 100)
	
	// Simulate successful combat
	instance.Memory.SuccessfulHits = 5
	instance.Memory.DamageReceived = 10
	
	// Update should recognize advantage
	instance.Update(150, 100)
	
	// With good performance, should be in aggressive or normal state
	if instance.TacticalState != TacticalAggressive && instance.TacticalState != TacticalNormal {
		t.Errorf("Expected aggressive/normal with advantage, got %v", instance.TacticalState)
	}
	
	// Take heavy damage
	instance.TakeDamage(60)
	instance.Update(150, 100)
	
	// Should switch to defensive or retreating
	if instance.TacticalState == TacticalAggressive {
		t.Error("Should not be aggressive after heavy damage")
	}
}

// TestFormationMovement tests formation-based positioning
func TestFormationMovement(t *testing.T) {
	enemy := &Enemy{
		Health:   100,
		Speed:    2.0,
		Behavior: PatrolBehavior,
	}
	
	instance := NewEnemyInstance(enemy, 0, 0)
	
	// Create and join a group
	group := NewEnemyGroup()
	group.AddMember(instance)
	instance.Group = group
	
	// Set formation position far from current position
	instance.FormationX = 100
	instance.FormationY = 50
	
	// Apply formation movement
	instance.applyFormationMovement()
	
	// Velocity should be adjusted toward formation position
	if instance.VelX == 0 {
		t.Error("Formation movement should adjust X velocity")
	}
	
	// Check direction is toward formation
	if instance.VelX < 0 {
		t.Error("Should move in positive X direction toward formation")
	}
}

// TestCoordinatedAttack tests multiple enemies coordinating
func TestCoordinatedAttack(t *testing.T) {
	group := NewEnemyGroup()
	
	// Create 3 enemies at different positions
	for i := 0; i < 3; i++ {
		enemy := &Enemy{
			Health:   100,
			Speed:    2.0,
			Behavior: ChaseBehavior,
		}
		instance := NewEnemyInstance(enemy, float64(i*100), 100)
		instance.State = ChaseState
		group.AddMember(instance)
		instance.Group = group
	}
	
	// Update group to coordinate
	playerX, playerY := 150.0, 100.0
	group.UpdateGroup(playerX, playerY)
	
	// Group should be engaging
	if group.GroupState != GroupEngaging {
		t.Errorf("Expected GroupEngaging, got %v", group.GroupState)
	}
	
	// Should have a formation
	if group.Formation == NoFormation {
		t.Error("Engaging group should have a formation")
	}
	
	// All members should have formation positions assigned
	for i, member := range group.Members {
		if member.FormationX == 0 && member.FormationY == 0 {
			t.Errorf("Member %d has no formation position", i)
		}
	}
}

// TestLearningBehavior tests AI learning over time
func TestLearningBehavior(t *testing.T) {
	mem := NewAIMemory()
	
	initialSkill := mem.PlayerSkillEstimate
	
	// Simulate skilled player behavior
	for i := 0; i < 50; i++ {
		// Frequent attacks and dashes
		mem.UpdateMemory(float64(i*10), 100, i%3 == 0, true, i%2 == 0)
	}
	
	// Skill estimate should increase
	if mem.PlayerSkillEstimate <= initialSkill {
		t.Error("Player skill estimate should increase with skilled play")
	}
	
	// Confidence should increase with more observations
	if mem.ConfidenceLevel <= 0.1 {
		t.Error("Confidence should increase with observations")
	}
}

// TestAdaptiveDifficulty tests AI difficulty adaptation
func TestAdaptiveDifficulty(t *testing.T) {
	enemy := &Enemy{
		Health:   100,
		Damage:   10,
		Speed:    2.0,
		Behavior: ChaseBehavior,
	}
	
	instance := NewEnemyInstance(enemy, 100, 100)
	
	initialThreshold := instance.Memory.RetreatThreshold
	
	// Simulate player dominating - take lots of damage
	for i := 0; i < 10; i++ {
		instance.TakeDamage(10)
	}
	
	// Retreat threshold should adapt to taking heavy damage
	if instance.Memory.RetreatThreshold <= initialThreshold {
		t.Errorf("Retreat threshold should increase after heavy damage: was %f, now %f", 
			initialThreshold, instance.Memory.RetreatThreshold)
	}
}
