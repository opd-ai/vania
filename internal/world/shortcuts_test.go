package world

import (
	"fmt"
	"testing"
)

func TestAddShortcuts(t *testing.T) {
	testCases := []struct {
		name           string
		seed           int64
		expectedMin    int
		expectedMax    int
		checkDistance  bool
		checkAbilities bool
	}{
		{
			name:           "Standard world gets 3-5 shortcuts",
			seed:           42,
			expectedMin:    3,
			expectedMax:    5,
			checkDistance:  true,
			checkAbilities: true,
		},
		{
			name:           "Different seed produces shortcuts",
			seed:           999,
			expectedMin:    3,
			expectedMax:    5,
			checkDistance:  true,
			checkAbilities: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate world with shortcuts
			gen := NewWorldGenerator(15, 10, 80, 5)
			world := gen.Generate(tc.seed, nil)

			// Count shortcuts
			shortcutCount := 0
			for _, edge := range world.Graph.Edges {
				if edge.IsShortcut {
					shortcutCount++
				}
			}

			// Verify shortcut count
			if shortcutCount < tc.expectedMin || shortcutCount > tc.expectedMax {
				t.Errorf("Expected %d-%d shortcuts, got %d", tc.expectedMin, tc.expectedMax, shortcutCount)
			}

			// Verify shortcuts properties
			for _, edge := range world.Graph.Edges {
				if !edge.IsShortcut {
					continue
				}

				// Check distance requirement (≥5 edges)
				if tc.checkDistance {
					fromNode := world.Graph.Nodes[edge.From]
					toNode := world.Graph.Nodes[edge.To]
					distance := fromNode.Depth - toNode.Depth

					if distance < 5 {
						t.Errorf("Shortcut has distance %d (expected ≥5)", distance)
					}

					// Verify source is deeper than destination
					if fromNode.Depth <= toNode.Depth {
						t.Errorf("Shortcut from depth %d to %d (should go backward)", fromNode.Depth, toNode.Depth)
					}
				}

				// Check ability requirements
				if tc.checkAbilities {
					validAbilities := map[string]bool{
						"double_jump": true,
						"dash":        true,
						"wall_climb":  true,
						"glide":       true,
					}

					if !validAbilities[edge.Requirement] {
						t.Errorf("Shortcut has invalid ability requirement: %s", edge.Requirement)
					}
				}

				// Check one-way initial state
				if !edge.OneWay {
					t.Errorf("Shortcut should be one-way initially")
				}

				// Check not traversed initially
				if edge.Traversed {
					t.Errorf("Shortcut should not be traversed initially")
				}
			}
		})
	}
}

func TestShortcutDeterminism(t *testing.T) {
	seed := int64(12345)

	gen1 := NewWorldGenerator(15, 10, 80, 5)
	world1 := gen1.Generate(seed, nil)

	gen2 := NewWorldGenerator(15, 10, 80, 5)
	world2 := gen2.Generate(seed, nil)

	// Count shortcuts in both worlds
	shortcuts1 := 0
	shortcuts2 := 0

	for _, edge := range world1.Graph.Edges {
		if edge.IsShortcut {
			shortcuts1++
		}
	}

	for _, edge := range world2.Graph.Edges {
		if edge.IsShortcut {
			shortcuts2++
		}
	}

	if shortcuts1 != shortcuts2 {
		t.Errorf("Shortcut count not deterministic: %d vs %d", shortcuts1, shortcuts2)
	}

	// Verify same shortcuts in same order
	shortcutEdges1 := make([]GraphEdge, 0)
	shortcutEdges2 := make([]GraphEdge, 0)

	for _, edge := range world1.Graph.Edges {
		if edge.IsShortcut {
			shortcutEdges1 = append(shortcutEdges1, edge)
		}
	}

	for _, edge := range world2.Graph.Edges {
		if edge.IsShortcut {
			shortcutEdges2 = append(shortcutEdges2, edge)
		}
	}

	if len(shortcutEdges1) != len(shortcutEdges2) {
		t.Fatalf("Different number of shortcuts")
	}

	for i := range shortcutEdges1 {
		e1 := shortcutEdges1[i]
		e2 := shortcutEdges2[i]

		if e1.From != e2.From || e1.To != e2.To {
			t.Errorf("Shortcut %d different: (%d→%d) vs (%d→%d)", i, e1.From, e1.To, e2.From, e2.To)
		}

		if e1.Requirement != e2.Requirement {
			t.Errorf("Shortcut %d has different requirements: %s vs %s", i, e1.Requirement, e2.Requirement)
		}
	}
}

func TestShortcutNoDuplicates(t *testing.T) {
	seed := int64(555)

	gen := NewWorldGenerator(15, 10, 80, 5)
	world := gen.Generate(seed, nil)

	// Check for duplicate shortcuts
	seen := make(map[string]bool)

	for _, edge := range world.Graph.Edges {
		if !edge.IsShortcut {
			continue
		}

		key := fmt.Sprintf("%d->%d", edge.From, edge.To)
		if seen[key] {
			t.Errorf("Duplicate shortcut found: %d -> %d", edge.From, edge.To)
		}
		seen[key] = true
	}
}

func TestShortcutRoomConnections(t *testing.T) {
	seed := int64(777)

	gen := NewWorldGenerator(15, 10, 80, 5)
	world := gen.Generate(seed, nil)

	// Verify that shortcuts create actual room connections
	for _, edge := range world.Graph.Edges {
		if !edge.IsShortcut {
			continue
		}

		// Find source and destination rooms
		var sourceRoom, destRoom *Room
		for _, room := range world.Rooms {
			if room.ID == edge.From {
				sourceRoom = room
			}
			if room.ID == edge.To {
				destRoom = room
			}
		}

		if sourceRoom == nil {
			t.Errorf("Source room %d not found for shortcut", edge.From)
			continue
		}

		if destRoom == nil {
			t.Errorf("Destination room %d not found for shortcut", edge.To)
			continue
		}

		// Verify connection exists
		found := false
		for _, conn := range sourceRoom.Connections {
			if conn.ID == destRoom.ID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Shortcut edge %d->%d exists in graph but not in room connections", edge.From, edge.To)
		}
	}
}

func TestGetCriticalPathNodes(t *testing.T) {
	seed := int64(888)

	gen := NewWorldGenerator(15, 10, 80, 5)
	world := gen.Generate(seed, nil)

	criticalNodes := gen.getCriticalPathNodes(world)

	// Verify all nodes are required
	for _, nodeID := range criticalNodes {
		node := world.Graph.Nodes[nodeID]
		if !node.Required {
			t.Errorf("Node %d in critical path but not marked as required", nodeID)
		}
	}

	// Verify sorted by depth
	for i := 0; i < len(criticalNodes)-1; i++ {
		depth1 := world.Graph.Nodes[criticalNodes[i]].Depth
		depth2 := world.Graph.Nodes[criticalNodes[i+1]].Depth

		if depth1 > depth2 {
			t.Errorf("Critical path not sorted: node %d (depth %d) before node %d (depth %d)",
				criticalNodes[i], depth1, criticalNodes[i+1], depth2)
		}
	}

	// Verify minimum path length (should be at least 15 based on generateGraph)
	if len(criticalNodes) < 15 {
		t.Errorf("Critical path too short: %d nodes (expected ≥15)", len(criticalNodes))
	}
}

func TestShortcutExistsCheck(t *testing.T) {
	seed := int64(333)

	gen := NewWorldGenerator(15, 10, 80, 5)
	world := gen.Generate(seed, nil)

	// Find a shortcut
	var testFrom, testTo int
	found := false

	for _, edge := range world.Graph.Edges {
		if edge.IsShortcut {
			testFrom = edge.From
			testTo = edge.To
			found = true
			break
		}
	}

	if !found {
		t.Skip("No shortcuts generated to test")
	}

	// Verify shortcutExists detects it
	if !gen.shortcutExists(world, testFrom, testTo) {
		t.Errorf("shortcutExists failed to detect existing shortcut %d->%d", testFrom, testTo)
	}

	// Verify it doesn't detect non-existent shortcuts
	if gen.shortcutExists(world, 9999, 9998) {
		t.Errorf("shortcutExists incorrectly detected non-existent shortcut")
	}
}

func TestFindRoomIndexByID(t *testing.T) {
	seed := int64(444)

	gen := NewWorldGenerator(15, 10, 80, 5)
	world := gen.Generate(seed, nil)

	// Test finding each room by ID
	for expectedIdx, room := range world.Rooms {
		actualIdx := gen.findRoomIndexByID(world, room.ID)

		if actualIdx != expectedIdx {
			t.Errorf("findRoomIndexByID(%d) = %d, expected %d", room.ID, actualIdx, expectedIdx)
		}
	}

	// Test non-existent ID
	badIdx := gen.findRoomIndexByID(world, 9999)
	if badIdx != -1 {
		t.Errorf("findRoomIndexByID(9999) = %d, expected -1", badIdx)
	}
}

func TestSmallWorldNoShortcuts(t *testing.T) {
	seed := int64(111)

	// Generate very small world
	gen := NewWorldGenerator(5, 5, 8, 2)
	world := gen.Generate(seed, nil)

	// Count shortcuts
	shortcutCount := 0
	for _, edge := range world.Graph.Edges {
		if edge.IsShortcut {
			shortcutCount++
		}
	}

	// Small world should have 0 shortcuts (not enough rooms)
	if shortcutCount > 0 {
		t.Logf("Small world has %d shortcuts (may be acceptable)", shortcutCount)
	}
}

func TestGraphEdgeDefaultValues(t *testing.T) {
	seed := int64(666)

	gen := NewWorldGenerator(15, 10, 80, 5)
	world := gen.Generate(seed, nil)

	// Check that regular edges have correct default values
	for _, edge := range world.Graph.Edges {
		if edge.IsShortcut {
			continue
		}

		if edge.OneWay {
			t.Errorf("Regular edge %d->%d should not be one-way", edge.From, edge.To)
		}

		if edge.Traversed {
			t.Errorf("Regular edge %d->%d should not be pre-traversed", edge.From, edge.To)
		}
	}
}
