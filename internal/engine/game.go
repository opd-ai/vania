package engine

import (
	"github.com/opd-ai/vania/internal/audio"
	"github.com/opd-ai/vania/internal/entity"
	"github.com/opd-ai/vania/internal/graphics"
	"github.com/opd-ai/vania/internal/narrative"
	"github.com/opd-ai/vania/internal/pcg"
	"github.com/opd-ai/vania/internal/world"
	"time"
)

// Game represents the main game state
type Game struct {
	World          *world.World
	Entities       []*entity.Enemy
	Bosses         []*entity.Boss
	Items          []*entity.Item
	Abilities      []entity.Ability
	Graphics       *GraphicsSystem
	Audio          *AudioSystem
	Narrative      *narrative.WorldContext
	Player         *Player
	CurrentRoom    *world.Room
	Running        bool
	Seed           int64
	PCGContext     *pcg.PCGContext
}

// Player represents the player character
type Player struct {
	X, Y           float64
	VelX, VelY     float64
	Health         int
	MaxHealth      int
	Damage         int
	Speed          float64
	Abilities      map[string]bool
	Inventory      []*entity.Item
	Sprite         *graphics.Sprite
}

// GraphicsSystem manages all graphics
type GraphicsSystem struct {
	SpriteGen  *graphics.SpriteGenerator
	TilesetGen *graphics.TilesetGenerator
	PaletteGen *graphics.PaletteGenerator
	Tilesets   map[string]*graphics.Tileset
	Sprites    map[string]*graphics.Sprite
}

// AudioSystem manages all audio
type AudioSystem struct {
	SFXGen   *audio.SFXGenerator
	MusicGen *audio.MusicGenerator
	Sounds   map[string]*audio.AudioSample
	Music    map[string]*audio.AudioSample
}

// GameGenerator orchestrates all generation
type GameGenerator struct {
	MasterSeed   int64
	GraphicsGen  *GraphicsGenerator
	AudioGen     *AudioGenerator
	NarrativeGen *narrative.NarrativeGenerator
	WorldGen     *world.WorldGenerator
	EntityGen    *EntityGenerator
	PCGContext   *pcg.PCGContext
}

// GraphicsGenerator manages graphics generation
type GraphicsGenerator struct {
	Seed int64
}

// AudioGenerator manages audio generation
type AudioGenerator struct {
	Seed int64
}

// EntityGenerator manages entity generation
type EntityGenerator struct {
	Seed int64
}

// NewGameGenerator creates a new game generator
func NewGameGenerator(masterSeed int64) *GameGenerator {
	seeds := pcg.DeriveSeeds(masterSeed)
	
	return &GameGenerator{
		MasterSeed:   masterSeed,
		GraphicsGen:  &GraphicsGenerator{Seed: seeds["graphics"]},
		AudioGen:     &AudioGenerator{Seed: seeds["audio"]},
		NarrativeGen: narrative.NewNarrativeGenerator(seeds["narrative"]),
		WorldGen:     world.NewWorldGenerator(15, 10, 100, 5),
		EntityGen:    &EntityGenerator{Seed: seeds["entity"]},
		PCGContext:   pcg.NewPCGContext(masterSeed),
	}
}

// GenerateCompleteGame creates a full game from seed
func (gg *GameGenerator) GenerateCompleteGame() (*Game, error) {
	startTime := time.Now()
	
	// Generate narrative context first (influences other systems)
	narrative := gg.NarrativeGen.Generate(pcg.HashSeed(gg.MasterSeed, "narrative"))
	
	// Generate visual style based on narrative theme
	graphicsSystem := gg.generateGraphics(narrative)
	
	// Generate world using narrative constraints
	worldData := gg.WorldGen.Generate(
		pcg.HashSeed(gg.MasterSeed, "world"),
		narrative.WorldConstraints,
	)
	
	// Generate entities that fit world biomes
	entities, bosses, items, abilities := gg.generateEntities(worldData, narrative)
	
	// Generate audio matching narrative tone
	audioSystem := gg.generateAudio(narrative, worldData)
	
	// Create player
	player := gg.createPlayer(graphicsSystem)
	
	// Validate generation
	if !gg.validate(worldData, entities, narrative) {
		// Could regenerate or return error
	}
	
	game := &Game{
		World:       worldData,
		Entities:    entities,
		Bosses:      bosses,
		Items:       items,
		Abilities:   abilities,
		Graphics:    graphicsSystem,
		Audio:       audioSystem,
		Narrative:   narrative,
		Player:      player,
		CurrentRoom: worldData.StartRoom,
		Running:     true,
		Seed:        gg.MasterSeed,
		PCGContext:  gg.PCGContext,
	}
	
	generationTime := time.Since(startTime)
	println("Game generated in", generationTime.Seconds(), "seconds")
	
	return game, nil
}

// generateGraphics creates all graphics
func (gg *GameGenerator) generateGraphics(narrative *narrative.WorldContext) *GraphicsSystem {
	system := &GraphicsSystem{
		SpriteGen:  graphics.NewSpriteGenerator(32, 32, graphics.VerticalSymmetry),
		TilesetGen: graphics.NewTilesetGenerator(16, string(narrative.Theme)),
		PaletteGen: graphics.NewPaletteGenerator(graphics.AnalogousScheme),
		Tilesets:   make(map[string]*graphics.Tileset),
		Sprites:    make(map[string]*graphics.Sprite),
	}
	
	// Generate player sprite
	playerSpriteGen := graphics.NewSpriteGenerator(16, 16, graphics.VerticalSymmetry)
	system.Sprites["player"] = playerSpriteGen.Generate(gg.GraphicsGen.Seed)
	
	// Generate tilesets for each biome
	biomeTypes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}
	for i, biome := range biomeTypes {
		tilesetGen := graphics.NewTilesetGenerator(16, biome)
		system.Tilesets[biome] = tilesetGen.Generate(gg.GraphicsGen.Seed + int64(i))
	}
	
	return system
}

// generateAudio creates all audio
func (gg *GameGenerator) generateAudio(narrative *narrative.WorldContext, worldData *world.World) *AudioSystem {
	system := &AudioSystem{
		SFXGen:   audio.NewSFXGenerator(44100),
		MusicGen: audio.NewMusicGenerator(44100, 90, 60, audio.MinorScale),
		Sounds:   make(map[string]*audio.AudioSample),
		Music:    make(map[string]*audio.AudioSample),
	}
	
	// Generate sound effects
	sfxTypes := []audio.SFXType{
		audio.JumpSFX,
		audio.LandSFX,
		audio.AttackSFX,
		audio.HitSFX,
		audio.PickupSFX,
		audio.DoorSFX,
		audio.DamageSFX,
	}
	
	for i, sfxType := range sfxTypes {
		key := []string{"jump", "land", "attack", "hit", "pickup", "door", "damage"}[i]
		system.Sounds[key] = system.SFXGen.Generate(sfxType, gg.AudioGen.Seed+int64(i))
	}
	
	// Generate music for each biome
	for i, biome := range worldData.Biomes {
		musicGen := gg.selectMusicGenerator(biome)
		system.Music[biome.Name] = musicGen.GenerateTrack(
			gg.AudioGen.Seed+int64(i*100),
			60.0, // 60 seconds
		)
	}
	
	return system
}

// selectMusicGenerator chooses appropriate music generator for biome
func (gg *GameGenerator) selectMusicGenerator(biome *world.Biome) *audio.MusicGenerator {
	mood := biome.GetMusicMood()
	
	var scale audio.Scale
	var bpm int
	
	switch mood {
	case "dark_ambient":
		scale = audio.MinorScale
		bpm = 70
	case "peaceful":
		scale = audio.MajorScale
		bpm = 80
	case "mysterious":
		scale = audio.DorianScale
		bpm = 75
	case "ethereal":
		scale = audio.PentatonicMaj
		bpm = 85
	case "horror":
		scale = audio.PhrygianScale
		bpm = 60
	default:
		scale = audio.MinorScale
		bpm = 90
	}
	
	return audio.NewMusicGenerator(44100, bpm, 60, scale)
}

// generateEntities creates all enemies, bosses, items, and abilities
func (gg *GameGenerator) generateEntities(worldData *world.World, narrative *narrative.WorldContext) ([]*entity.Enemy, []*entity.Boss, []*entity.Item, []entity.Ability) {
	enemyGen := entity.NewEnemyGenerator(gg.EntityGen.Seed)
	bossGen := entity.NewBossGenerator(gg.EntityGen.Seed + 1000)
	itemGen := entity.NewItemGenerator(gg.EntityGen.Seed + 2000)
	abilityGen := entity.NewAbilityGenerator(gg.EntityGen.Seed + 3000)
	
	var enemies []*entity.Enemy
	var bosses []*entity.Boss
	var items []*entity.Item
	
	// Generate enemies for each room
	for i, room := range worldData.Rooms {
		if room.Type == world.CombatRoom {
			for j := 0; j < len(room.Enemies); j++ {
				enemy := enemyGen.Generate(
					room.Biome.Name,
					room.Biome.DangerLevel,
					gg.EntityGen.Seed+int64(i*1000+j),
				)
				enemies = append(enemies, enemy)
			}
		} else if room.Type == world.BossRoom {
			boss := bossGen.Generate(
				room.Biome.Name,
				gg.EntityGen.Seed+int64(i*1000),
			)
			bosses = append(bosses, boss)
		}
		
		// Generate items for treasure rooms
		if room.Type == world.TreasureRoom {
			for j := 0; j < len(room.Items); j++ {
				itemType := entity.ItemType(j % 3)
				item := itemGen.Generate(
					itemType,
					gg.EntityGen.Seed+int64(i*100+j),
				)
				items = append(items, item)
			}
		}
	}
	
	// Generate ability progression
	abilities := abilityGen.GenerateProgression(gg.EntityGen.Seed)
	
	return enemies, bosses, items, abilities
}

// createPlayer creates the player character
func (gg *GameGenerator) createPlayer(gfx *GraphicsSystem) *Player {
	return &Player{
		X:         100,
		Y:         100,
		Health:    100,
		MaxHealth: 100,
		Damage:    10,
		Speed:     5.0,
		Abilities: make(map[string]bool),
		Inventory: make([]*entity.Item, 0),
		Sprite:    gfx.Sprites["player"],
	}
}

// validate checks if generation is valid
func (gg *GameGenerator) validate(worldData *world.World, entities []*entity.Enemy, narrative *narrative.WorldContext) bool {
	// Check world has start room
	if worldData.StartRoom == nil {
		return false
	}
	
	// Check world has boss rooms
	if len(worldData.BossRooms) == 0 {
		return false
	}
	
	// Check minimum room count
	if len(worldData.Rooms) < 50 {
		return false
	}
	
	// More validation could be added
	return true
}

// Run starts the game loop (stub for now)
func (g *Game) Run() {
	println("Game starting with seed:", g.Seed)
	println("Theme:", g.Narrative.Theme)
	println("Player motivation:", g.Narrative.PlayerMotivation)
	println("World rooms:", len(g.World.Rooms))
	println("Enemies:", len(g.Entities))
	println("Bosses:", len(g.Bosses))
	println("Items:", len(g.Items))
	println("Abilities:", len(g.Abilities))
	
	// Game loop would go here
	// For now, just print stats
}
