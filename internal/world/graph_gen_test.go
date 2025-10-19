// Package world provides tests for world generation
package world

import (
	"testing"
)

// TestRoomCreationAndConnections verifies that rooms are properly created
// and connected based on the graph edges, regardless of map iteration order
func TestRoomCreationAndConnections(t *testing.T) {
	wg := NewWorldGenerator(15, 10, 50, 3)
	constraints := make(map[string]interface{})
	
	// Test with multiple seeds to ensure consistency
	seeds := []int64{12345, 67890, 11111, 99999, 42}
	
	for _, seed := range seeds {
		world := wg.Generate(seed, constraints)
		
		// Verify basic structure
		if world.StartRoom == nil {
			t.Errorf("Seed %d: StartRoom is nil", seed)
		}
		
		if len(world.BossRooms) == 0 {
			t.Errorf("Seed %d: No boss rooms generated", seed)
		}
		
		if len(world.Rooms) == 0 {
			t.Errorf("Seed %d: No rooms generated", seed)
		}
		
		// Verify graph node count matches room count
		if len(world.Graph.Nodes) != len(world.Rooms) {
			t.Errorf("Seed %d: Graph nodes (%d) don't match room count (%d)", 
				seed, len(world.Graph.Nodes), len(world.Rooms))
		}
		
		// Build lookup map for verification
		roomByID := make(map[int]*Room)
		for _, room := range world.Rooms {
			if existing, exists := roomByID[room.ID]; exists {
				t.Errorf("Seed %d: Duplicate room ID %d (existing: %v, new: %v)", 
					seed, room.ID, existing, room)
			}
			roomByID[room.ID] = room
		}
		
		// Verify all graph edges correspond to room connections
		for _, edge := range world.Graph.Edges {
			fromRoom, fromExists := roomByID[edge.From]
			toRoom, toExists := roomByID[edge.To]
			
			if !fromExists {
				t.Errorf("Seed %d: Edge references non-existent from room %d", seed, edge.From)
				continue
			}
			
			if !toExists {
				t.Errorf("Seed %d: Edge references non-existent to room %d", seed, edge.To)
				continue
			}
			
			// Check if toRoom is in fromRoom's connections
			found := false
			for _, conn := range fromRoom.Connections {
				if conn == toRoom {
					found = true
					break
				}
			}
			
			if !found {
				t.Errorf("Seed %d: Edge (%d->%d) exists in graph but room %d is not connected to room %d", 
					seed, edge.From, edge.To, edge.From, edge.To)
			}
		}
		
		// Verify each room has proper fields
		for _, room := range world.Rooms {
			if room.Biome == nil {
				t.Errorf("Seed %d: Room %d has nil biome", seed, room.ID)
			}
			
			if room.Width <= 0 || room.Height <= 0 {
				t.Errorf("Seed %d: Room %d has invalid dimensions: %dx%d", 
					seed, room.ID, room.Width, room.Height)
			}
		}
	}
}

// TestWorldGenerationDeterminism verifies that the same seed produces
// identical worlds
func TestWorldGenerationDeterminism(t *testing.T) {
	seed := int64(42)
	constraints := make(map[string]interface{})
	
	wg1 := NewWorldGenerator(15, 10, 50, 3)
	world1 := wg1.Generate(seed, constraints)
	
	wg2 := NewWorldGenerator(15, 10, 50, 3)
	world2 := wg2.Generate(seed, constraints)
	
	// Compare basic properties
	if len(world1.Rooms) != len(world2.Rooms) {
		t.Errorf("Room count mismatch: %d vs %d", len(world1.Rooms), len(world2.Rooms))
	}
	
	if world1.StartRoom == nil || world2.StartRoom == nil {
		t.Errorf("Start room is nil")
	} else if world1.StartRoom.ID != world2.StartRoom.ID {
		t.Errorf("Start room ID mismatch: %d vs %d", 
			world1.StartRoom.ID, world2.StartRoom.ID)
	}
	
	if len(world1.BossRooms) != len(world2.BossRooms) {
		t.Errorf("Boss room count mismatch: %d vs %d", 
			len(world1.BossRooms), len(world2.BossRooms))
	}
	
	// Compare graph structure
	if len(world1.Graph.Nodes) != len(world2.Graph.Nodes) {
		t.Errorf("Graph node count mismatch: %d vs %d", 
			len(world1.Graph.Nodes), len(world2.Graph.Nodes))
	}
	
	if len(world1.Graph.Edges) != len(world2.Graph.Edges) {
		t.Errorf("Graph edge count mismatch: %d vs %d", 
			len(world1.Graph.Edges), len(world2.Graph.Edges))
	}
}

// TestNoRoomIndexOutOfBounds verifies that room connections don't
// cause index out of bounds errors
func TestNoRoomIndexOutOfBounds(t *testing.T) {
	wg := NewWorldGenerator(20, 15, 100, 5)
	constraints := make(map[string]interface{})
	
	// Generate world - this should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("World generation panicked: %v", r)
		}
	}()
	
	world := wg.Generate(12345, constraints)
	
	// Traverse all room connections - this should not panic
	visited := make(map[int]bool)
	var traverse func(*Room)
	traverse = func(room *Room) {
		if visited[room.ID] {
			return
		}
		visited[room.ID] = true
		
		for _, conn := range room.Connections {
			if conn == nil {
				t.Errorf("Room %d has nil connection", room.ID)
				continue
			}
			traverse(conn)
		}
	}
	
	if world.StartRoom != nil {
		traverse(world.StartRoom)
	}
}
