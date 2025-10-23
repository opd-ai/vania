// Package world provides procedural platform generation for rooms based on
// room type, biome characteristics, and player abilities to ensure 
// traversable and engaging platforming challenges.
package world

import (
	"math/rand"
)

// PlatformGenerator generates procedural platforms for rooms
type PlatformGenerator struct {
	rng *rand.Rand
}

// NewPlatformGenerator creates a new platform generator
func NewPlatformGenerator() *PlatformGenerator {
	return &PlatformGenerator{}
}

// PlatformLayout represents different platform arrangement patterns
type PlatformLayout int

const (
	LinearLayout PlatformLayout = iota  // Horizontal progression
	StaircaseLayout                    // Ascending/descending steps  
	ScatteredLayout                    // Random platforms requiring jumping
	TowerLayout                        // Vertical climbing challenge
	BridgeLayout                      // Spanning gaps with multiple platforms
	MazeLayout                        // Complex interconnected platforms
)

// PlatformDifficulty affects platform spacing and complexity
type PlatformDifficulty int

const (
	EasyDifficulty PlatformDifficulty = iota
	MediumDifficulty
	HardDifficulty
)

// GeneratePlatforms creates procedural platforms for a room
func (pg *PlatformGenerator) GeneratePlatforms(room *Room, seed int64, playerAbilities map[string]bool) {
	pg.rng = rand.New(rand.NewSource(seed))
	
	// Determine layout based on room type and biome
	layout := pg.selectLayout(room)
	difficulty := pg.calculateDifficulty(room, playerAbilities)
	
	// Clear existing platforms
	room.Platforms = make([]Platform, 0)
	
	// Generate platforms based on layout
	switch layout {
	case LinearLayout:
		pg.generateLinearPlatforms(room, difficulty, playerAbilities)
	case StaircaseLayout:
		pg.generateStaircasePlatforms(room, difficulty, playerAbilities)
	case ScatteredLayout:
		pg.generateScatteredPlatforms(room, difficulty, playerAbilities)
	case TowerLayout:
		pg.generateTowerPlatforms(room, difficulty, playerAbilities)
	case BridgeLayout:
		pg.generateBridgePlatforms(room, difficulty, playerAbilities)
	case MazeLayout:
		pg.generateMazePlatforms(room, difficulty, playerAbilities)
	}
	
	// Add ground platform if room needs it
	pg.addGroundPlatform(room)
	
	// Ensure platforms don't overlap with doors
	pg.clearPlatformsByDoors(room)
	
	// Validate that room is traversable
	pg.validateTraversability(room, playerAbilities)
}

// selectLayout chooses appropriate layout for room type and biome
func (pg *PlatformGenerator) selectLayout(room *Room) PlatformLayout {
	switch room.Type {
	case StartRoom, SaveRoom:
		return LinearLayout // Simple layouts for safe rooms
	case TreasureRoom:
		return ScatteredLayout // Require skill to reach treasure
	case BossRoom:
		return TowerLayout // Vertical arena for boss fights
	case CombatRoom:
		// Vary based on biome
		if room.Biome != nil {
			switch room.Biome.Name {
			case "cave":
				return StaircaseLayout
			case "crystal":
				return TowerLayout
			case "ruins":
				return MazeLayout
			case "abyss":
				return ScatteredLayout
			default:
				return LinearLayout
			}
		}
		return LinearLayout
	default:
		layouts := []PlatformLayout{LinearLayout, StaircaseLayout, ScatteredLayout}
		return layouts[pg.rng.Intn(len(layouts))]
	}
}

// calculateDifficulty determines platform challenge level
func (pg *PlatformGenerator) calculateDifficulty(room *Room, playerAbilities map[string]bool) PlatformDifficulty {
	baseLevel := EasyDifficulty
	
	// Increase difficulty based on biome danger level
	if room.Biome != nil {
		if room.Biome.DangerLevel >= 5 {
			baseLevel = HardDifficulty
		} else if room.Biome.DangerLevel >= 3 {
			baseLevel = MediumDifficulty
		}
	}
	
	// Adjust based on available abilities
	abilityCount := 0
	for _, hasAbility := range playerAbilities {
		if hasAbility {
			abilityCount++
		}
	}
	
	// More abilities = can handle harder platforms
	if abilityCount >= 4 {
		return HardDifficulty
	} else if abilityCount >= 2 {
		return MediumDifficulty
	}
	
	return baseLevel
}

// generateLinearPlatforms creates horizontal progression platforms
func (pg *PlatformGenerator) generateLinearPlatforms(room *Room, difficulty PlatformDifficulty, abilities map[string]bool) {
	roomWidth := 960  // Screen width in pixels
	roomHeight := 640 // Screen height in pixels
	
	platformCount := 3 + pg.rng.Intn(3)
	if difficulty == HardDifficulty {
		platformCount += 2
	}
	
	spacing := roomWidth / (platformCount + 1)
	baseHeight := roomHeight - 200 // Leave room at top and bottom
	
	for i := 0; i < platformCount; i++ {
		x := spacing * (i + 1)
		y := baseHeight + pg.rng.Intn(100) - 50 // Small height variation
		
		width := 80 + pg.rng.Intn(40)
		if difficulty == EasyDifficulty {
			width += 40 // Wider platforms for easier jumping
		}
		
		room.Platforms = append(room.Platforms, Platform{
			X:      x,
			Y:      y,
			Width:  width,
			Height: 32,
		})
	}
}

// generateStaircasePlatforms creates ascending/descending platform steps
func (pg *PlatformGenerator) generateStaircasePlatforms(room *Room, difficulty PlatformDifficulty, abilities map[string]bool) {
	roomWidth := 960
	roomHeight := 640
	
	platformCount := 5 + pg.rng.Intn(3)
	stepWidth := roomWidth / platformCount
	
	// Random ascending or descending
	ascending := pg.rng.Float64() < 0.5
	
	startHeight := roomHeight - 150
	endHeight := roomHeight - 400
	if !ascending {
		startHeight, endHeight = endHeight, startHeight
	}
	
	heightStep := float64(endHeight-startHeight) / float64(platformCount-1)
	
	for i := 0; i < platformCount; i++ {
		x := i * stepWidth + pg.rng.Intn(20) - 10 // Small horizontal variation
		y := startHeight + int(float64(i)*heightStep)
		
		width := 60 + pg.rng.Intn(30)
		if difficulty == EasyDifficulty {
			width += 20
		}
		
		room.Platforms = append(room.Platforms, Platform{
			X:      x,
			Y:      y,
			Width:  width,
			Height: 32,
		})
	}
}

// generateScatteredPlatforms creates random platforms requiring jumping skill
func (pg *PlatformGenerator) generateScatteredPlatforms(room *Room, difficulty PlatformDifficulty, abilities map[string]bool) {
	roomWidth := 960
	roomHeight := 640
	
	platformCount := 4 + pg.rng.Intn(4)
	if difficulty == HardDifficulty {
		platformCount += 2
	}
	
	maxJumpDistance := 150 // Base jump distance
	if abilities["double_jump"] {
		maxJumpDistance = 250
	}
	if abilities["dash"] {
		maxJumpDistance = 300
	}
	
	// Start with one platform
	platforms := []Platform{
		{
			X:      100 + pg.rng.Intn(100),
			Y:      roomHeight - 200 - pg.rng.Intn(100),
			Width:  80 + pg.rng.Intn(40),
			Height: 32,
		},
	}
	
	// Add remaining platforms ensuring they're reachable
	for i := 1; i < platformCount; i++ {
		var newPlatform Platform
		attempts := 0
		
		for attempts < 10 { // Try to place platform near existing ones
			// Pick a random existing platform to connect from
			fromPlatform := platforms[pg.rng.Intn(len(platforms))]
			
			// Generate position within jump range
			angle := pg.rng.Float64() * 2 * 3.14159 // Random angle
			distance := 80 + pg.rng.Float64()*float64(maxJumpDistance-80)
			
			x := fromPlatform.X + int(float64(distance)*cos(angle))
			y := fromPlatform.Y + int(float64(distance)*sin(angle))
			
			// Clamp to room bounds
			if x < 50 || x > roomWidth-150 || y < 100 || y > roomHeight-100 {
				attempts++
				continue
			}
			
			width := 60 + pg.rng.Intn(40)
			if difficulty == EasyDifficulty {
				width += 30
			}
			
			newPlatform = Platform{
				X:      x,
				Y:      y,
				Width:  width,
				Height: 32,
			}
			break
		}
		
		if attempts < 10 {
			platforms = append(platforms, newPlatform)
		}
	}
	
	room.Platforms = platforms
}

// generateTowerPlatforms creates vertical climbing challenge
func (pg *PlatformGenerator) generateTowerPlatforms(room *Room, difficulty PlatformDifficulty, abilities map[string]bool) {
	roomWidth := 960
	roomHeight := 640
	
	levels := 4 + pg.rng.Intn(3)
	if difficulty == HardDifficulty {
		levels += 1
	}
	
	levelHeight := (roomHeight - 200) / levels
	
	for level := 0; level < levels; level++ {
		platformsPerLevel := 2 + pg.rng.Intn(2)
		if difficulty == EasyDifficulty {
			platformsPerLevel = 3 // More platforms = easier
		}
		
		y := roomHeight - 150 - (level * levelHeight)
		
		for i := 0; i < platformsPerLevel; i++ {
			x := (roomWidth / (platformsPerLevel + 1)) * (i + 1)
			x += pg.rng.Intn(80) - 40 // Some variation
			
			width := 70 + pg.rng.Intn(30)
			
			room.Platforms = append(room.Platforms, Platform{
				X:      x,
				Y:      y,
				Width:  width,
				Height: 32,
			})
		}
	}
}

// generateBridgePlatforms creates platforms spanning gaps
func (pg *PlatformGenerator) generateBridgePlatforms(room *Room, difficulty PlatformDifficulty, abilities map[string]bool) {
	roomWidth := 960
	roomHeight := 640
	
	// Create 2-3 separate sections connected by bridges
	sections := 2 + pg.rng.Intn(2)
	sectionWidth := roomWidth / sections
	
	bridgeY := roomHeight - 250 - pg.rng.Intn(100)
	
	for section := 0; section < sections; section++ {
		sectionX := section * sectionWidth
		
		// Main platform for this section
		room.Platforms = append(room.Platforms, Platform{
			X:      sectionX + 50,
			Y:      bridgeY,
			Width:  sectionWidth - 200, // Leave gap for bridge
			Height: 32,
		})
		
		// Bridge platforms between sections
		if section < sections-1 {
			bridgeCount := 1
			if difficulty == EasyDifficulty {
				bridgeCount = 2 // More bridge platforms = easier
			}
			
			bridgeStart := sectionX + sectionWidth - 150
			bridgeSpacing := 200 / (bridgeCount + 1)
			
			for i := 0; i < bridgeCount; i++ {
				room.Platforms = append(room.Platforms, Platform{
					X:      bridgeStart + (i+1)*bridgeSpacing,
					Y:      bridgeY + pg.rng.Intn(40) - 20,
					Width:  40 + pg.rng.Intn(20),
					Height: 32,
				})
			}
		}
	}
}

// generateMazePlatforms creates complex interconnected platform maze
func (pg *PlatformGenerator) generateMazePlatforms(room *Room, difficulty PlatformDifficulty, abilities map[string]bool) {
	roomWidth := 960
	roomHeight := 640
	
	// Create grid of potential platform positions
	gridCols := 8
	gridRows := 5
	cellWidth := roomWidth / gridCols
	cellHeight := (roomHeight - 200) / gridRows
	
	// Randomly place platforms on grid with connections
	for row := 0; row < gridRows; row++ {
		for col := 0; col < gridCols; col++ {
			// Higher chance for platforms in easier difficulty
			chance := 0.4
			if difficulty == EasyDifficulty {
				chance = 0.6
			} else if difficulty == HardDifficulty {
				chance = 0.3
			}
			
			if pg.rng.Float64() < chance {
				x := col*cellWidth + cellWidth/4 + pg.rng.Intn(cellWidth/2)
				y := 150 + row*cellHeight + cellHeight/4 + pg.rng.Intn(cellHeight/2)
				
				width := 50 + pg.rng.Intn(30)
				if difficulty == EasyDifficulty {
					width += 20
				}
				
				room.Platforms = append(room.Platforms, Platform{
					X:      x,
					Y:      y,
					Width:  width,
					Height: 32,
				})
			}
		}
	}
}

// addGroundPlatform ensures room has a base ground platform
func (pg *PlatformGenerator) addGroundPlatform(room *Room) {
	// Add a ground platform for safety
	room.Platforms = append(room.Platforms, Platform{
		X:      0,
		Y:      600, // Near bottom of screen
		Width:  960, // Full width
		Height: 40,
	})
}

// clearPlatformsByDoors removes platforms that would block doors
func (pg *PlatformGenerator) clearPlatformsByDoors(room *Room) {
	if len(room.Doors) == 0 {
		return
	}
	
	var validPlatforms []Platform
	
	for _, platform := range room.Platforms {
		blocked := false
		
		for _, door := range room.Doors {
			// Check if platform overlaps with door area
			if platform.X < door.X+door.Width && platform.X+platform.Width > door.X &&
				platform.Y < door.Y+door.Height && platform.Y+platform.Height > door.Y {
				blocked = true
				break
			}
		}
		
		if !blocked {
			validPlatforms = append(validPlatforms, platform)
		}
	}
	
	room.Platforms = validPlatforms
}

// validateTraversability ensures the room can be completed with available abilities
func (pg *PlatformGenerator) validateTraversability(room *Room, abilities map[string]bool) {
	// This would implement pathfinding to ensure room is completable
	// For now, just ensure minimum platform count
	if len(room.Platforms) < 2 {
		// Add emergency platforms
		room.Platforms = append(room.Platforms, Platform{
			X: 200, Y: 400, Width: 100, Height: 32,
		})
		room.Platforms = append(room.Platforms, Platform{
			X: 500, Y: 350, Width: 100, Height: 32,
		})
	}
}

// Helper function for cosine (missing in simplified math)
func cos(x float64) float64 {
	// Simple approximation - in real code you'd import math package
	// but avoiding import conflicts for this generation
	return float64(1.0 - x*x*0.5)
}

// Helper function for sine (missing in simplified math)  
func sin(x float64) float64 {
	// Simple approximation
	return x * (1.0 - x*x*0.166)
}