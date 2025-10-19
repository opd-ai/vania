package engine

import (
	"testing"

	"github.com/opd-ai/vania/internal/world"
)

func TestRoomTransitionHandler_CheckDoorCollision(t *testing.T) {
	// Create a simple game with a room
	game := &Game{
		CurrentRoom: &world.Room{
			ID:   1,
			Type: world.CombatRoom,
			Doors: []world.Door{
				{
					X:      100,
					Y:      200,
					Width:  64,
					Height: 96,
					Direction: "east",
					Locked: false,
				},
			},
		},
	}

	handler := NewRoomTransitionHandler(game)

	tests := []struct {
		name     string
		playerX  float64
		playerY  float64
		playerW  float64
		playerH  float64
		wantDoor bool
	}{
		{
			name:     "Player at door",
			playerX:  110,
			playerY:  210,
			playerW:  32,
			playerH:  32,
			wantDoor: true,
		},
		{
			name:     "Player far from door",
			playerX:  500,
			playerY:  500,
			playerW:  32,
			playerH:  32,
			wantDoor: false,
		},
		{
			name:     "Player touching edge of door",
			playerX:  100,
			playerY:  200,
			playerW:  32,
			playerH:  32,
			wantDoor: true,
		},
		{
			name:     "Player just past door",
			playerX:  165,
			playerY:  297,
			playerW:  32,
			playerH:  32,
			wantDoor: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			door := handler.CheckDoorCollision(tt.playerX, tt.playerY, tt.playerW, tt.playerH)
			if (door != nil) != tt.wantDoor {
				t.Errorf("CheckDoorCollision() returned door = %v, want door = %v", door != nil, tt.wantDoor)
			}
		})
	}
}

func TestRoomTransitionHandler_StartTransition(t *testing.T) {
	// Create game with two rooms
	room1 := &world.Room{ID: 1, Type: world.CombatRoom}
	room2 := &world.Room{ID: 2, Type: world.TreasureRoom}

	game := &Game{
		CurrentRoom: room1,
		Player: &Player{
			X: 100,
			Y: 200,
		},
	}

	handler := NewRoomTransitionHandler(game)

	door := &world.Door{
		X:         100,
		Y:         200,
		Width:     64,
		Height:    96,
		Direction: "east",
		LeadsTo:   room2,
		Locked:    false,
	}

	// Start transition
	handler.StartTransition(door)

	if !handler.IsTransitioning() {
		t.Error("Expected transition to be active after StartTransition")
	}

	if handler.targetRoom != room2 {
		t.Errorf("Expected target room to be room2, got %v", handler.targetRoom)
	}
}

func TestRoomTransitionHandler_Update(t *testing.T) {
	room1 := &world.Room{ID: 1, Type: world.CombatRoom}
	room2 := &world.Room{ID: 2, Type: world.TreasureRoom}

	game := &Game{
		CurrentRoom: room1,
		Player: &Player{
			X:     100,
			Y:     200,
			VelX:  5.0,
			VelY:  3.0,
		},
	}

	handler := NewRoomTransitionHandler(game)

	door := &world.Door{
		LeadsTo: room2,
	}

	handler.StartTransition(door)

	// Update until transition completes
	completed := false
	maxIterations := 100
	iterations := 0

	for iterations < maxIterations {
		if handler.Update() {
			completed = true
			break
		}
		iterations++
	}

	if !completed {
		t.Error("Transition did not complete within expected time")
	}

	if game.CurrentRoom != room2 {
		t.Errorf("Expected current room to be room2 after transition, got %v", game.CurrentRoom)
	}

	if game.Player.VelX != 0 || game.Player.VelY != 0 {
		t.Errorf("Expected player velocity to be reset, got VelX=%v, VelY=%v", game.Player.VelX, game.Player.VelY)
	}
}

func TestRoomTransitionHandler_LockedDoor(t *testing.T) {
	room1 := &world.Room{
		ID:   1,
		Type: world.CombatRoom,
		Doors: []world.Door{
			{
				X:         100,
				Y:         200,
				Width:     64,
				Height:    96,
				Direction: "east",
				Locked:    true, // Door is locked
			},
		},
	}

	game := &Game{
		CurrentRoom: room1,
	}

	handler := NewRoomTransitionHandler(game)

	// Try to collide with locked door
	door := handler.CheckDoorCollision(110, 210, 32, 32)

	// Should return nil because door is locked
	if door != nil {
		t.Error("Expected nil door when colliding with locked door")
	}
}

func TestRoomTransitionHandler_GetTransitionProgress(t *testing.T) {
	game := &Game{
		CurrentRoom: &world.Room{ID: 1},
		Player:      &Player{},
	}

	handler := NewRoomTransitionHandler(game)

	// No transition
	progress := handler.GetTransitionProgress()
	if progress != 0.0 {
		t.Errorf("Expected progress 0.0 when not transitioning, got %v", progress)
	}

	// Start transition
	door := &world.Door{
		LeadsTo: &world.Room{ID: 2},
	}
	handler.StartTransition(door)

	// Progress should increase as we update
	firstProgress := handler.GetTransitionProgress()
	handler.Update()
	secondProgress := handler.GetTransitionProgress()

	if secondProgress <= firstProgress {
		t.Errorf("Expected progress to increase, got first=%v, second=%v", firstProgress, secondProgress)
	}
}

func TestRoomTransitionHandler_SpawnEnemiesForRoom(t *testing.T) {
	// Create a simple game (SpawnEnemiesForRoom returns empty if no entities)
	game := &Game{
		Entities: nil, // Will check logic even with nil entities
	}

	handler := NewRoomTransitionHandler(game)

	tests := []struct {
		name       string
		roomType   world.RoomType
		shouldSpawn bool
	}{
		{
			name:       "Combat room logic",
			roomType:   world.CombatRoom,
			shouldSpawn: true,
		},
		{
			name:       "Boss room logic",
			roomType:   world.BossRoom,
			shouldSpawn: true,
		},
		{
			name:       "Corridor room spawns nothing",
			roomType:   world.CorridorRoom,
			shouldSpawn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &world.Room{
				Type:    tt.roomType,
				Enemies: make([]interface{}, 3),
			}

			enemies := handler.SpawnEnemiesForRoom(room)
			
			// With no entities in game, should return empty list
			if len(enemies) != 0 {
				t.Errorf("Expected 0 enemies with nil entity list, got %d", len(enemies))
			}
		})
	}
}

