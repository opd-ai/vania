// Package engine provides room transition functionality for moving between
// connected rooms in the game world.
package engine

import (
	"fmt"

	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/world"
)

// TransitionType represents the visual style of room transitions
type TransitionType string

const (
	// TransitionFade fades to black between rooms
	TransitionFade TransitionType = "fade"
	// TransitionSlide slides camera from old room to new room
	TransitionSlide TransitionType = "slide"
	// TransitionIris uses an iris wipe effect
	TransitionIris TransitionType = "iris"
)

// DefaultTransitionDuration is the default transition time in frames at 60 FPS (0.5 seconds)
const DefaultTransitionDuration = 30

// RoomTransitionHandler manages room transitions
type RoomTransitionHandler struct {
	game              *Game
	transitionActive  bool
	transitionTimer   int
	maxTransitionTime int
	targetRoom        *world.Room
	sourceRoom        *world.Room
	transitionType    TransitionType
	slideDirection    string // Direction for slide transitions: "left", "right", "up", "down"
}

// NewRoomTransitionHandler creates a new room transition handler
func NewRoomTransitionHandler(game *Game) *RoomTransitionHandler {
	return &RoomTransitionHandler{
		game:              game,
		transitionActive:  false,
		transitionTimer:   0,
		maxTransitionTime: DefaultTransitionDuration,
		transitionType:    TransitionFade,
		slideDirection:    "right",
	}
}

// SetTransitionType sets the type of transition to use
func (rth *RoomTransitionHandler) SetTransitionType(transitionType TransitionType) {
	rth.transitionType = transitionType
}

// SetTransitionDuration sets the transition duration in seconds
func (rth *RoomTransitionHandler) SetTransitionDuration(seconds float64) {
	// Clamp to 0.3-0.8 seconds as specified
	if seconds < 0.3 {
		seconds = 0.3
	}
	if seconds > 0.8 {
		seconds = 0.8
	}
	rth.maxTransitionTime = int(seconds * 60.0) // Convert to frames at 60 FPS
}

// CheckDoorCollision checks if player is touching a door
func (rth *RoomTransitionHandler) CheckDoorCollision(playerX, playerY, playerW, playerH float64, unlockedDoors map[string]bool) *world.Door {
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

			// Check if door is locked and not unlocked yet
			doorKey := rth.GetDoorKey(door)
			if door.Locked && !unlockedDoors[doorKey] {
				// Door is locked - caller will handle UI message
				return nil
			}

			return door
		}
	}

	return nil
}

// GetDoorKey generates a unique key for a door based on its position and room
func (rth *RoomTransitionHandler) GetDoorKey(door *world.Door) string {
	if rth.game.CurrentRoom == nil {
		return ""
	}
	// Create unique door identifier using room ID and door properties
	return fmt.Sprintf("room_%d_door_%d_%d_%s",
		rth.game.CurrentRoom.ID, door.X, door.Y, door.Direction)
}

// CanUnlockDoor checks if player has the required ability/key to unlock a door
func (rth *RoomTransitionHandler) CanUnlockDoor(door *world.Door, playerAbilities map[string]bool, collectedItems map[int]bool) bool {
	if door == nil || !door.Locked {
		return true
	}

	// Check if player has required ability
	// Door requirements are stored in world graph edges
	if door.LeadsTo != nil {
		// Look up edge requirement for this connection
		requirement := rth.findEdgeRequirement(rth.game.CurrentRoom.ID, door.LeadsTo.ID)
		if requirement != "" {
			// Check if player has the required ability
			if !playerAbilities[requirement] {
				return false
			}
		}
	}

	return true
}

// findEdgeRequirement finds the ability requirement for transitioning between rooms
func (rth *RoomTransitionHandler) findEdgeRequirement(fromRoomID, toRoomID int) string {
	if rth.game.World == nil || rth.game.World.Graph == nil {
		return ""
	}

	// Search through graph edges for matching connection
	for _, edge := range rth.game.World.Graph.Edges {
		if edge.From == fromRoomID && edge.To == toRoomID {
			return edge.Requirement
		}
		// Check reverse direction as well
		if edge.From == toRoomID && edge.To == fromRoomID {
			return edge.Requirement
		}
	}

	return ""
}

// StartTransition initiates a room transition
func (rth *RoomTransitionHandler) StartTransition(door *world.Door) {
	if door == nil || door.LeadsTo == nil {
		return
	}

	rth.transitionActive = true
	rth.transitionTimer = rth.maxTransitionTime
	rth.sourceRoom = rth.game.CurrentRoom
	rth.targetRoom = door.LeadsTo

	// Determine slide direction based on door direction
	rth.slideDirection = door.Direction
	if rth.slideDirection == "" {
		rth.slideDirection = "right" // Default direction
	}
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

	progress := 1.0 - (float64(rth.transitionTimer) / float64(rth.maxTransitionTime))
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	return progress
}

// GetTransitionType returns the current transition type
func (rth *RoomTransitionHandler) GetTransitionType() TransitionType {
	return rth.transitionType
}

// GetSlideDirection returns the direction for slide transitions
func (rth *RoomTransitionHandler) GetSlideDirection() string {
	return rth.slideDirection
}

// GetSourceRoom returns the room we're transitioning from
func (rth *RoomTransitionHandler) GetSourceRoom() *world.Room {
	return rth.sourceRoom
}

// GetTargetRoom returns the room we're transitioning to
func (rth *RoomTransitionHandler) GetTargetRoom() *world.Room {
	return rth.targetRoom
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

// SpawnItemsForRoom creates item instances for the current room
func (rth *RoomTransitionHandler) SpawnItemsForRoom(room *world.Room) []*entity.ItemInstance {
	return createItemInstancesForRoom(room, rth.game.Items)
}
