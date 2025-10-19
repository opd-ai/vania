// Package engine provides room transition functionality for moving between
// connected rooms in the game world.
package engine

import (
	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/world"
)

// RoomTransitionHandler manages room transitions
type RoomTransitionHandler struct {
	game              *Game
	transitionActive  bool
	transitionTimer   int
	targetRoom        *world.Room
	transitionEffect  string // "fade", "slide", etc.
}

// NewRoomTransitionHandler creates a new room transition handler
func NewRoomTransitionHandler(game *Game) *RoomTransitionHandler {
	return &RoomTransitionHandler{
		game:              game,
		transitionActive:  false,
		transitionTimer:   0,
		transitionEffect:  "fade",
	}
}

// CheckDoorCollision checks if player is touching a door
func (rth *RoomTransitionHandler) CheckDoorCollision(playerX, playerY, playerW, playerH float64) *world.Door {
	if rth.game.CurrentRoom == nil {
		return nil
	}
	
	// Check collision with each door
	for i := range rth.game.CurrentRoom.Doors {
		door := &rth.game.CurrentRoom.Doors[i]
		
		// Simple AABB collision
		doorX := float64(door.X)
		doorY := float64(door.Y)
		doorW := float64(door.Width)
		doorH := float64(door.Height)
		
		if playerX < doorX+doorW &&
			playerX+playerW > doorX &&
			playerY < doorY+doorH &&
			playerY+playerH > doorY {
			
			// Check if door is locked
			if door.Locked {
				// TODO: Show "locked" message or play sound
				continue
			}
			
			return door
		}
	}
	
	return nil
}

// StartTransition initiates a room transition
func (rth *RoomTransitionHandler) StartTransition(door *world.Door) {
	if door == nil || door.LeadsTo == nil {
		return
	}
	
	rth.transitionActive = true
	rth.transitionTimer = 30 // 0.5 seconds at 60 FPS
	rth.targetRoom = door.LeadsTo
}

// Update updates the transition state
func (rth *RoomTransitionHandler) Update() bool {
	if !rth.transitionActive {
		return false
	}
	
	rth.transitionTimer--
	
	if rth.transitionTimer <= 0 {
		// Complete the transition
		rth.CompleteTransition()
		rth.transitionActive = false
		return true
	}
	
	return false
}

// CompleteTransition finishes the room transition
func (rth *RoomTransitionHandler) CompleteTransition() {
	if rth.targetRoom == nil {
		return
	}
	
	// Switch to new room
	rth.game.CurrentRoom = rth.targetRoom
	
	// Position player based on entry direction
	// For now, place player near the opposite side of where they entered
	rth.game.Player.X = 100.0 // Default position
	rth.game.Player.Y = 500.0
	
	// Reset player velocity
	rth.game.Player.VelX = 0
	rth.game.Player.VelY = 0
}

// IsTransitioning returns if a transition is in progress
func (rth *RoomTransitionHandler) IsTransitioning() bool {
	return rth.transitionActive
}

// GetTransitionProgress returns progress from 0.0 to 1.0
func (rth *RoomTransitionHandler) GetTransitionProgress() float64 {
	if !rth.transitionActive {
		return 0.0
	}
	
	maxTime := 30.0
	progress := 1.0 - (float64(rth.transitionTimer) / maxTime)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	return progress
}

// SpawnEnemiesForRoom creates enemy instances for the current room
func (rth *RoomTransitionHandler) SpawnEnemiesForRoom(room *world.Room) []*entity.EnemyInstance {
	var enemyInstances []*entity.EnemyInstance
	
	if room == nil || len(rth.game.Entities) == 0 {
		return enemyInstances
	}
	
	// Determine enemy count based on room type
	enemyCount := 0
	switch room.Type {
	case world.CombatRoom:
		enemyCount = 3 + (len(room.Enemies) % 3) // 3-5 enemies
	case world.BossRoom:
		enemyCount = 1 // One boss
	case world.TreasureRoom:
		enemyCount = 1 + (len(room.Enemies) % 2) // 1-2 guards
	default:
		enemyCount = 0
	}
	
	// Spawn enemies
	for i := 0; i < enemyCount && i < len(rth.game.Entities); i++ {
		enemy := rth.game.Entities[i]
		
		// Position enemies across the room
		enemyX := 300.0 + float64(i*150)
		enemyY := 500.0
		
		enemyInstances = append(enemyInstances, entity.NewEnemyInstance(enemy, enemyX, enemyY))
	}
	
	return enemyInstances
}
