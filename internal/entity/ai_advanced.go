// Package entity provides advanced AI capabilities including memory,
// learning behaviors, and coordinated group tactics for enemy entities.
package entity

import (
	"math"
	"time"
)

// AIMemory tracks patterns and learns from player behavior
type AIMemory struct {
	// Player movement patterns
	LastPlayerPositions []Position      // Ring buffer of recent player positions
	PlayerVelocityAvg   float64         // Average player speed
	JumpFrequency       float64         // How often player jumps
	AttackFrequency     float64         // How often player attacks
	DashFrequency       float64         // How often player uses dash
	
	// Combat statistics
	DamageReceived      int             // Total damage taken from this player
	AttacksEvaded       int             // Number of player attacks dodged
	SuccessfulHits      int             // Number of hits landed on player
	LastHitTime         time.Time       // When enemy last hit player
	
	// Learned behaviors
	PreferredAttackDistance float64     // Distance at which enemy has most success
	OptimalApproachAngle    float64     // Best angle to approach player
	RetreatThreshold        float64     // Health % at which to retreat
	
	// Tactical awareness
	KnowsPlayerPosition bool            // Has seen player recently
	LastKnownPlayerX    float64
	LastKnownPlayerY    float64
	TimesSeeingPlayer   int             // Total encounters
	
	// Adaptive difficulty
	PlayerSkillEstimate float64         // 0.0-1.0, higher = more skilled player
	ConfidenceLevel     float64         // 0.0-1.0, how confident in predictions
	
	// Learning parameters
	LastUpdateTime      time.Time
	LearningRate        float64         // How quickly to adapt (0.0-1.0)
}

// Position represents a 2D coordinate
type Position struct {
	X, Y float64
}

// EnemyGroup manages coordinated behavior between multiple enemies
type EnemyGroup struct {
	Members           []*EnemyInstance
	Leader            *EnemyInstance  // Strongest enemy becomes leader
	Formation         FormationType
	TargetX, TargetY  float64        // Shared target position
	GroupState        GroupState
	LastCoordination  time.Time
	CommunicationRange float64       // Max distance for coordination
}

// FormationType defines group tactical formations
type FormationType int

const (
	NoFormation FormationType = iota
	LineFormation      // Enemies in a horizontal line
	CircleFormation    // Surround player
	PincerFormation    // Attack from two sides
	VFormation         // Leader in front, others flanking
	ScatteredFormation // Spread out to avoid AoE
)

// GroupState represents the group's tactical state
type GroupState int

const (
	GroupIdle GroupState = iota
	GroupPatrol
	GroupEngaging
	GroupRetreating
	GroupRegrouping
)

// TacticalState represents advanced tactical decision states
type TacticalState int

const (
	TacticalNormal TacticalState = iota
	TacticalAggressive   // Push advantage
	TacticalDefensive    // Protect self
	TacticalFlanking     // Try to get behind player
	TacticalKiting       // Hit and run
	TacticalRetreating   // Fallback to safety
	TacticalRegrouping   // Wait for allies
)

// NewAIMemory creates a new AI memory system
func NewAIMemory() *AIMemory {
	return &AIMemory{
		LastPlayerPositions:     make([]Position, 0, 20), // Store last 20 positions
		PlayerVelocityAvg:       0.0,
		JumpFrequency:           0.0,
		AttackFrequency:         0.0,
		DashFrequency:           0.0,
		PreferredAttackDistance: 50.0, // Start with default
		OptimalApproachAngle:    0.0,
		RetreatThreshold:        0.3, // Retreat at 30% health
		KnowsPlayerPosition:     false,
		PlayerSkillEstimate:     0.5, // Assume average skill
		ConfidenceLevel:         0.0, // No confidence yet
		LastUpdateTime:          time.Now(),
		LearningRate:            0.05, // Moderate learning speed
	}
}

// UpdateMemory processes new observations and updates learned behaviors
func (mem *AIMemory) UpdateMemory(playerX, playerY float64, playerDidJump, playerDidAttack, playerDidDash bool) {
	now := time.Now()
	deltaTime := now.Sub(mem.LastUpdateTime).Seconds()
	if deltaTime <= 0 {
		deltaTime = 0.016 // Assume ~60 FPS
	}
	
	// Update player position history
	if len(mem.LastPlayerPositions) >= 20 {
		// Shift ring buffer
		mem.LastPlayerPositions = append(mem.LastPlayerPositions[1:], Position{X: playerX, Y: playerY})
	} else {
		mem.LastPlayerPositions = append(mem.LastPlayerPositions, Position{X: playerX, Y: playerY})
	}
	
	// Calculate player velocity if we have enough data
	if len(mem.LastPlayerPositions) >= 2 {
		last := mem.LastPlayerPositions[len(mem.LastPlayerPositions)-1]
		prev := mem.LastPlayerPositions[len(mem.LastPlayerPositions)-2]
		dx := last.X - prev.X
		dy := last.Y - prev.Y
		velocity := math.Sqrt(dx*dx + dy*dy) / deltaTime
		
		// Smooth velocity with exponential moving average
		mem.PlayerVelocityAvg = mem.PlayerVelocityAvg*0.9 + velocity*0.1
	}
	
	// Update action frequencies (exponential moving average)
	if playerDidJump {
		mem.JumpFrequency = mem.JumpFrequency*0.95 + 1.0*0.05
	} else {
		mem.JumpFrequency = mem.JumpFrequency * 0.95
	}
	
	if playerDidAttack {
		mem.AttackFrequency = mem.AttackFrequency*0.95 + 1.0*0.05
	} else {
		mem.AttackFrequency = mem.AttackFrequency * 0.95
	}
	
	if playerDidDash {
		mem.DashFrequency = mem.DashFrequency*0.95 + 1.0*0.05
	} else {
		mem.DashFrequency = mem.DashFrequency * 0.95
	}
	
	// Update tactical awareness
	mem.KnowsPlayerPosition = true
	mem.LastKnownPlayerX = playerX
	mem.LastKnownPlayerY = playerY
	mem.TimesSeeingPlayer++
	
	// Estimate player skill based on action frequency and movement
	// High skill: frequent attacks, good movement, uses dash effectively
	skillIndicator := (mem.AttackFrequency*0.4 + mem.DashFrequency*0.3 + 
	                   math.Min(mem.PlayerVelocityAvg/10.0, 1.0)*0.3)
	mem.PlayerSkillEstimate = mem.PlayerSkillEstimate*(1.0-mem.LearningRate) + 
	                          skillIndicator*mem.LearningRate
	
	// Increase confidence as we gather more data
	if mem.TimesSeeingPlayer > 0 {
		mem.ConfidenceLevel = math.Min(1.0, float64(mem.TimesSeeingPlayer)/100.0)
	}
	
	mem.LastUpdateTime = now
}

// RecordCombatEvent updates memory with combat outcomes
func (mem *AIMemory) RecordCombatEvent(hitPlayer bool, tookDamage bool, damageAmount int, distance float64) {
	if hitPlayer {
		mem.SuccessfulHits++
		mem.LastHitTime = time.Now()
		
		// Learn preferred attack distance
		mem.PreferredAttackDistance = mem.PreferredAttackDistance*0.9 + distance*0.1
	}
	
	if tookDamage {
		mem.DamageReceived += damageAmount
		
		// Adjust retreat threshold based on how much damage we're taking
		if mem.DamageReceived > 50 {
			// Getting hit a lot - be more cautious
			mem.RetreatThreshold += 0.01
			if mem.RetreatThreshold > 0.7 {
				mem.RetreatThreshold = 0.7
			}
		}
	}
}

// RecordEvasion updates memory when successfully evading an attack
func (mem *AIMemory) RecordEvasion() {
	mem.AttacksEvaded++
}

// ShouldRetreat determines if enemy should retreat based on learned patterns
func (mem *AIMemory) ShouldRetreat(currentHealthPercent float64) bool {
	return currentHealthPercent < mem.RetreatThreshold
}

// PredictPlayerPosition attempts to predict where player will be
func (mem *AIMemory) PredictPlayerPosition(deltaTime float64) (float64, float64) {
	if len(mem.LastPlayerPositions) < 2 {
		return mem.LastKnownPlayerX, mem.LastKnownPlayerY
	}
	
	// Use last known position and velocity to predict
	last := mem.LastPlayerPositions[len(mem.LastPlayerPositions)-1]
	prev := mem.LastPlayerPositions[len(mem.LastPlayerPositions)-2]
	
	velX := last.X - prev.X
	velY := last.Y - prev.Y
	
	// Predict based on velocity and confidence
	predX := last.X + velX*deltaTime*mem.ConfidenceLevel
	predY := last.Y + velY*deltaTime*mem.ConfidenceLevel
	
	return predX, predY
}

// GetTacticalState determines tactical state based on memory and current situation
func (mem *AIMemory) GetTacticalState(healthPercent float64, hasAllies bool, distanceToPlayer float64) TacticalState {
	// Retreat if health is low
	if mem.ShouldRetreat(healthPercent) {
		if hasAllies {
			return TacticalRegrouping
		}
		return TacticalRetreating
	}
	
	// Be aggressive if we're winning
	if mem.SuccessfulHits > mem.DamageReceived/10 && healthPercent > 0.6 {
		return TacticalAggressive
	}
	
	// Try flanking if player is skilled
	if mem.PlayerSkillEstimate > 0.7 && distanceToPlayer > 100 {
		return TacticalFlanking
	}
	
	// Kiting behavior for hit-and-run
	if mem.AttacksEvaded > 3 && distanceToPlayer < mem.PreferredAttackDistance {
		return TacticalKiting
	}
	
	// Defensive if taking too much damage
	if mem.DamageReceived > 30 && healthPercent < 0.5 {
		return TacticalDefensive
	}
	
	return TacticalNormal
}

// NewEnemyGroup creates a new enemy group for coordination
func NewEnemyGroup() *EnemyGroup {
	return &EnemyGroup{
		Members:            make([]*EnemyInstance, 0, 5),
		Leader:             nil,
		Formation:          NoFormation,
		GroupState:         GroupIdle,
		LastCoordination:   time.Now(),
		CommunicationRange: 400.0, // Enemies within 400 pixels can coordinate
	}
}

// AddMember adds an enemy to the group
func (g *EnemyGroup) AddMember(enemy *EnemyInstance) {
	g.Members = append(g.Members, enemy)
	g.updateLeader()
}

// RemoveMember removes an enemy from the group
func (g *EnemyGroup) RemoveMember(enemy *EnemyInstance) {
	for i, member := range g.Members {
		if member == enemy {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			break
		}
	}
	g.updateLeader()
}

// updateLeader selects the strongest alive enemy as leader
func (g *EnemyGroup) updateLeader() {
	var strongest *EnemyInstance
	maxHealth := 0
	
	for _, member := range g.Members {
		if member.CurrentHealth > maxHealth && member.State != DeadState {
			maxHealth = member.CurrentHealth
			strongest = member
		}
	}
	
	g.Leader = strongest
}

// UpdateGroup coordinates group behavior
func (g *EnemyGroup) UpdateGroup(playerX, playerY float64) {
	if len(g.Members) == 0 {
		return
	}
	
	// Remove dead members
	aliveMembers := make([]*EnemyInstance, 0, len(g.Members))
	for _, member := range g.Members {
		if member.State != DeadState {
			aliveMembers = append(aliveMembers, member)
		}
	}
	g.Members = aliveMembers
	
	if len(g.Members) == 0 {
		return
	}
	
	// Update target position
	g.TargetX = playerX
	g.TargetY = playerY
	
	// Determine group state
	inCombat := false
	needsRegroup := false
	for _, member := range g.Members {
		if member.State == ChaseState || member.State == AttackState {
			inCombat = true
		}
		if float64(member.CurrentHealth)/float64(member.Enemy.Health) < 0.3 {
			needsRegroup = true
		}
	}
	
	if needsRegroup && len(g.Members) > 1 {
		g.GroupState = GroupRegrouping
	} else if inCombat {
		g.GroupState = GroupEngaging
	} else {
		g.GroupState = GroupPatrol
	}
	
	// Select formation based on group state and size
	g.selectFormation()
	
	// Apply formation positions
	g.applyFormation(playerX, playerY)
	
	g.LastCoordination = time.Now()
}

// selectFormation chooses appropriate formation
func (g *EnemyGroup) selectFormation() {
	memberCount := len(g.Members)
	
	switch g.GroupState {
	case GroupEngaging:
		if memberCount >= 4 {
			g.Formation = CircleFormation // Surround player
		} else if memberCount >= 2 {
			g.Formation = PincerFormation // Pincer attack
		} else {
			g.Formation = NoFormation
		}
	case GroupRegrouping:
		g.Formation = LineFormation // Defensive line
	case GroupPatrol:
		if memberCount >= 3 {
			g.Formation = VFormation // Leader in front
		} else {
			g.Formation = LineFormation
		}
	default:
		g.Formation = NoFormation
	}
}

// applyFormation calculates and applies formation positions to members
func (g *EnemyGroup) applyFormation(playerX, playerY float64) {
	if g.Formation == NoFormation || len(g.Members) == 0 {
		return
	}
	
	switch g.Formation {
	case LineFormation:
		// Arrange in horizontal line
		spacing := 80.0
		startX := g.TargetX - spacing*float64(len(g.Members)-1)/2.0
		for i, member := range g.Members {
			member.FormationX = startX + float64(i)*spacing
			member.FormationY = g.TargetY - 100.0 // Stay back from target
		}
		
	case CircleFormation:
		// Arrange in circle around player
		radius := 120.0
		angleStep := 2.0 * math.Pi / float64(len(g.Members))
		for i, member := range g.Members {
			angle := float64(i) * angleStep
			member.FormationX = playerX + math.Cos(angle)*radius
			member.FormationY = playerY + math.Sin(angle)*radius
		}
		
	case PincerFormation:
		// Split into two groups attacking from sides
		half := len(g.Members) / 2
		for i, member := range g.Members {
			if i < half {
				member.FormationX = playerX - 150.0
			} else {
				member.FormationX = playerX + 150.0
			}
			member.FormationY = playerY + float64(i%half)*60.0
		}
		
	case VFormation:
		// Leader in front, others behind in V shape
		if g.Leader != nil {
			g.Leader.FormationX = playerX
			g.Leader.FormationY = playerY - 100.0
			
			wingIndex := 0
			for _, member := range g.Members {
				if member != g.Leader {
					side := 1.0
					if wingIndex%2 == 0 {
						side = -1.0
					}
					member.FormationX = playerX + side*float64(wingIndex/2+1)*60.0
					member.FormationY = playerY - 150.0 - float64(wingIndex/2)*40.0
					wingIndex++
				}
			}
		}
		
	case ScatteredFormation:
		// Spread out randomly around target
		for i, member := range g.Members {
			angle := float64(i) * 1.618 * math.Pi // Golden angle for good distribution
			distance := 80.0 + float64(i%3)*40.0
			member.FormationX = playerX + math.Cos(angle)*distance
			member.FormationY = playerY + math.Sin(angle)*distance
		}
	}
}

// ShouldCoordinate returns true if members should coordinate
func (g *EnemyGroup) ShouldCoordinate() bool {
	return len(g.Members) >= 2 && time.Since(g.LastCoordination) > time.Millisecond*100
}

// GetNearbyAllies returns allies within communication range of given position
func GetNearbyAllies(enemies []*EnemyInstance, x, y, maxRange float64) []*EnemyInstance {
	allies := make([]*EnemyInstance, 0)
	
	for _, enemy := range enemies {
		if enemy.State == DeadState {
			continue
		}
		
		dx := enemy.X - x
		dy := enemy.Y - y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist < maxRange && dist > 1.0 { // Don't include self
			allies = append(allies, enemy)
		}
	}
	
	return allies
}
