// Package world generates Metroidvania-style game worlds using graph theory
// to ensure connected, playable level layouts with ability-gated progression,
// biome variety, and procedurally placed rooms, platforms, and hazards.
package world

import (
	"math/rand"
)

// Room represents a single room in the world
type Room struct {
	ID          int
	Type        RoomType
	X, Y        int // Grid position
	Width       int
	Height      int
	Connections []*Room
	Biome       *Biome
	Enemies     []interface{} // Will be populated with enemy data
	Items       []interface{} // Will be populated with item data
	Platforms   []Platform
	Hazards     []Hazard
	Doors       []Door // Exits to other rooms
}

// RoomType defines room archetypes
type RoomType int

const (
	CombatRoom RoomType = iota
	PuzzleRoom
	TreasureRoom
	CorridorRoom
	BossRoom
	StartRoom
	SaveRoom
)

// Platform represents a platform in a room
type Platform struct {
	X, Y   int
	Width  int
	Height int
}

// Hazard represents a dangerous element
type Hazard struct {
	X, Y       int
	Type       string // "spike", "lava", "electric"
	Damage     int
	Width      int
	Height     int
}

// Door represents an exit/entrance to another room
type Door struct {
	X, Y          int    // Position in the room
	Width, Height int    // Size of the door
	Direction     string // "north", "south", "east", "west"
	LeadsTo       *Room  // Connected room
	Locked        bool   // Whether door requires ability/key
}

// World represents the complete game world
type World struct {
	Rooms       []*Room
	StartRoom   *Room
	BossRooms   []*Room
	Biomes      []*Biome
	Width       int // Number of rooms wide
	Height      int // Number of rooms tall
	Graph       *WorldGraph
}

// WorldGraph represents connectivity between rooms
type WorldGraph struct {
	Nodes map[int]*GraphNode
	Edges []GraphEdge
}

// GraphNode represents a room in the graph
type GraphNode struct {
	RoomID   int
	Depth    int // Distance from start
	Required bool // On critical path
}

// GraphEdge represents a connection between rooms
type GraphEdge struct {
	From         int
	To           int
	Requirement  string // "double_jump", "dash", etc.
	IsShortcut   bool
}

// WorldGenerator generates game worlds
type WorldGenerator struct {
	Width       int
	Height      int
	RoomCount   int
	BiomeCount  int
	rng         *rand.Rand
}

// NewWorldGenerator creates a new world generator
func NewWorldGenerator(width, height, roomCount, biomeCount int) *WorldGenerator {
	// Validate and apply defaults
	if width <= 0 {
		width = 15
	}
	if height <= 0 {
		height = 10
	}
	if roomCount <= 0 {
		roomCount = 80
	}
	if biomeCount <= 0 {
		biomeCount = 5
	}
	
	return &WorldGenerator{
		Width:      width,
		Height:     height,
		RoomCount:  roomCount,
		BiomeCount: biomeCount,
	}
}

// Generate creates a complete world
func (wg *WorldGenerator) Generate(seed int64, constraints map[string]interface{}) *World {
	wg.rng = rand.New(rand.NewSource(seed))
	
	world := &World{
		Rooms:  make([]*Room, 0, wg.RoomCount),
		Biomes: make([]*Biome, wg.BiomeCount),
		Width:  wg.Width,
		Height: wg.Height,
		Graph:  &WorldGraph{
			Nodes: make(map[int]*GraphNode),
			Edges: make([]GraphEdge, 0),
		},
	}
	
	// Generate biomes first
	for i := 0; i < wg.BiomeCount; i++ {
		world.Biomes[i] = wg.generateBiome(i)
	}
	
	// Generate world graph structure
	wg.generateGraph(world)
	
	// Create rooms based on graph
	wg.createRooms(world)
	
	// Generate room contents
	for _, room := range world.Rooms {
		wg.populateRoom(room)
	}
	
	// Add shortcuts for backtracking
	wg.addShortcuts(world)
	
	return world
}

// generateGraph creates the world connectivity structure
func (wg *WorldGenerator) generateGraph(world *World) {
	roomID := 0
	
	// Create start room node
	world.Graph.Nodes[roomID] = &GraphNode{
		RoomID:   roomID,
		Depth:    0,
		Required: true,
	}
	roomID++
	
	// Generate critical path
	criticalPathLength := 15 + wg.rng.Intn(10)
	currentNode := 0
	
	for i := 0; i < criticalPathLength; i++ {
		world.Graph.Nodes[roomID] = &GraphNode{
			RoomID:   roomID,
			Depth:    i + 1,
			Required: true,
		}
		
		// Determine if this edge requires an ability
		requirement := ""
		if i > 0 && i%5 == 0 {
			abilities := []string{"double_jump", "dash", "wall_climb", "glide"}
			requirement = abilities[wg.rng.Intn(len(abilities))]
		}
		
		world.Graph.Edges = append(world.Graph.Edges, GraphEdge{
			From:        currentNode,
			To:          roomID,
			Requirement: requirement,
			IsShortcut:  false,
		})
		
		currentNode = roomID
		roomID++
	}
	
	// Add boss room at end of critical path
	bossRoomID := roomID
	world.Graph.Nodes[bossRoomID] = &GraphNode{
		RoomID:   bossRoomID,
		Depth:    criticalPathLength,
		Required: true,
	}
	world.Graph.Edges = append(world.Graph.Edges, GraphEdge{
		From:        currentNode,
		To:          bossRoomID,
		Requirement: "",
		IsShortcut:  false,
	})
	roomID++
	
	// Add side branches for exploration
	branchCount := (wg.RoomCount - roomID) / 3
	for i := 0; i < branchCount; i++ {
		// Pick a random node on critical path
		sourceNode := wg.rng.Intn(criticalPathLength)
		
		// Create 2-3 room branch
		branchLength := 2 + wg.rng.Intn(2)
		currentBranch := sourceNode
		
		for j := 0; j < branchLength; j++ {
			world.Graph.Nodes[roomID] = &GraphNode{
				RoomID:   roomID,
				Depth:    world.Graph.Nodes[currentBranch].Depth + 1,
				Required: false,
			}
			
			world.Graph.Edges = append(world.Graph.Edges, GraphEdge{
				From:        currentBranch,
				To:          roomID,
				Requirement: "",
				IsShortcut:  false,
			})
			
			currentBranch = roomID
			roomID++
		}
	}
}

// createRooms instantiates rooms based on graph
func (wg *WorldGenerator) createRooms(world *World) {
	// Create a map to lookup rooms by ID (map iteration order is random!)
	roomByID := make(map[int]*Room)
	
	for id, node := range world.Graph.Nodes {
		room := &Room{
			ID:          id,
			Type:        wg.determineRoomType(node),
			X:           wg.rng.Intn(world.Width),
			Y:           wg.rng.Intn(world.Height),
			Width:       20 + wg.rng.Intn(10),
			Height:      15 + wg.rng.Intn(5),
			Connections: make([]*Room, 0),
			Biome:       wg.assignBiome(world, node.Depth),
		}
		
		world.Rooms = append(world.Rooms, room)
		roomByID[id] = room // Store for lookup by ID
		
		if room.Type == StartRoom {
			world.StartRoom = room
		} else if room.Type == BossRoom {
			world.BossRooms = append(world.BossRooms, room)
		}
	}
	
	// Connect rooms based on edges using ID lookup
	for _, edge := range world.Graph.Edges {
		fromRoom := roomByID[edge.From]
		toRoom := roomByID[edge.To]
		if fromRoom != nil && toRoom != nil {
			fromRoom.Connections = append(fromRoom.Connections, toRoom)
		}
	}
}

// determineRoomType assigns type based on graph position
func (wg *WorldGenerator) determineRoomType(node *GraphNode) RoomType {
	if node.Depth == 0 {
		return StartRoom
	}
	
	// Check if this should be a boss room (at critical path end)
	if node.Required && node.Depth > 10 {
		maxDepth := 0
		for _, n := range wg.rng.Perm(len([]int{0})) {
			_ = n
			maxDepth++
		}
		// Simplified: assume last node could be boss
		return BossRoom
	}
	
	if !node.Required {
		// Side branch rooms are often treasure or puzzle
		if wg.rng.Float64() < 0.5 {
			return TreasureRoom
		}
		return PuzzleRoom
	}
	
	// Main path rooms
	roll := wg.rng.Float64()
	if roll < 0.5 {
		return CombatRoom
	} else if roll < 0.7 {
		return CorridorRoom
	} else if roll < 0.85 {
		return PuzzleRoom
	} else {
		return SaveRoom
	}
}

// assignBiome assigns a biome to a room based on depth
func (wg *WorldGenerator) assignBiome(world *World, depth int) *Biome {
	// Assign biomes in sequential zones
	biomeIdx := (depth * len(world.Biomes)) / 20
	if biomeIdx >= len(world.Biomes) {
		biomeIdx = len(world.Biomes) - 1
	}
	return world.Biomes[biomeIdx]
}

// populateRoom generates platforms, enemies, and items
func (wg *WorldGenerator) populateRoom(room *Room) {
	// Generate platforms based on room size
	platformCount := 3 + wg.rng.Intn(5)
	room.Platforms = make([]Platform, platformCount)
	
	for i := range room.Platforms {
		room.Platforms[i] = Platform{
			X:      wg.rng.Intn(room.Width - 4),
			Y:      wg.rng.Intn(room.Height - 2),
			Width:  3 + wg.rng.Intn(4),
			Height: 1,
		}
	}
	
	// Add hazards based on room type and biome
	if room.Type == CombatRoom || room.Type == BossRoom {
		hazardCount := 1 + wg.rng.Intn(3)
		room.Hazards = make([]Hazard, hazardCount)
		
		for i := range room.Hazards {
			hazardTypes := []string{"spike", "lava", "electric"}
			room.Hazards[i] = Hazard{
				X:      wg.rng.Intn(room.Width - 2),
				Y:      wg.rng.Intn(room.Height - 2),
				Type:   hazardTypes[wg.rng.Intn(len(hazardTypes))],
				Damage: 1 + wg.rng.Intn(2),
				Width:  2,
				Height: 1,
			}
		}
	}
	
	// Room-specific population
	switch room.Type {
	case CombatRoom:
		// Will be populated with enemies by entity generator
		room.Enemies = make([]interface{}, 2+wg.rng.Intn(3))
	case TreasureRoom:
		// Will be populated with items
		room.Items = make([]interface{}, 1+wg.rng.Intn(2))
	case BossRoom:
		// One boss enemy
		room.Enemies = make([]interface{}, 1)
	}
	
	// Generate doors based on connections
	wg.generateDoors(room)
}

// generateDoors creates doors for each room connection
func (wg *WorldGenerator) generateDoors(room *Room) {
	room.Doors = make([]Door, 0, len(room.Connections))
	
	// Standard room dimensions in pixels (assuming 960x640 screen)
	roomWidthPixels := 960
	roomHeightPixels := 640
	doorWidth := 64
	doorHeight := 96
	
	for i, connectedRoom := range room.Connections {
		if connectedRoom == nil {
			continue
		}
		
		var door Door
		door.LeadsTo = connectedRoom
		door.Width = doorWidth
		door.Height = doorHeight
		door.Locked = false // Will be set based on requirements later
		
		// Determine door direction based on relative position
		// For now, place doors evenly around the room perimeter
		direction := i % 4 // Cycle through: east, west, north, south
		
		switch direction {
		case 0: // East (right side)
			door.Direction = "east"
			door.X = roomWidthPixels - doorWidth - 10
			door.Y = roomHeightPixels/2 - doorHeight/2
			
		case 1: // West (left side)
			door.Direction = "west"
			door.X = 10
			door.Y = roomHeightPixels/2 - doorHeight/2
			
		case 2: // North (top) - for vertical movement
			door.Direction = "north"
			door.X = roomWidthPixels/2 - doorWidth/2
			door.Y = 10
			
		case 3: // South (bottom)
			door.Direction = "south"
			door.X = roomWidthPixels/2 - doorWidth/2
			door.Y = roomHeightPixels - doorHeight - 10
		}
		
		room.Doors = append(room.Doors, door)
	}
}

// addShortcuts creates backtracking shortcuts
func (wg *WorldGenerator) addShortcuts(world *World) {
	shortcutCount := 2 + wg.rng.Intn(2)
	
	for i := 0; i < shortcutCount; i++ {
		// Connect a deep room back to an early room
		if len(world.Rooms) < 10 {
			continue
		}
		
		// Pick random rooms by index
		fromIdx := len(world.Rooms)/2 + wg.rng.Intn(len(world.Rooms)/2)
		toIdx := wg.rng.Intn(len(world.Rooms) / 4)
		
		// Get the actual room IDs (not indices!)
		fromRoomID := world.Rooms[fromIdx].ID
		toRoomID := world.Rooms[toIdx].ID
		
		world.Graph.Edges = append(world.Graph.Edges, GraphEdge{
			From:        fromRoomID,  // Use room ID, not index
			To:          toRoomID,    // Use room ID, not index
			Requirement: "",
			IsShortcut:  true,
		})
		
		// Directly connect the rooms
		world.Rooms[fromIdx].Connections = append(
			world.Rooms[fromIdx].Connections,
			world.Rooms[toIdx],
		)
	}
}

// generateBiome creates a biome
func (wg *WorldGenerator) generateBiome(index int) *Biome {
	biomeTypes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}
	
	return &Biome{
		Name:        biomeTypes[index%len(biomeTypes)],
		Temperature: -10 + wg.rng.Intn(40),
		Moisture:    wg.rng.Intn(100),
		DangerLevel: index + wg.rng.Intn(3),
		Theme:       biomeTypes[index%len(biomeTypes)],
	}
}
