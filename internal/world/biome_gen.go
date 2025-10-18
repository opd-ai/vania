package world

// Biome represents an environmental zone
type Biome struct {
	Name        string
	Temperature int // -20 to 40 degrees
	Moisture    int // 0-100 percent
	DangerLevel int // 1-10
	Theme       string
	ColorScheme []string
	EnemyTypes  []string
	Hazards     []string
}

// BiomeGenerator creates biome definitions
type BiomeGenerator struct {
}

// NewBiomeGenerator creates a new biome generator
func NewBiomeGenerator() *BiomeGenerator {
	return &BiomeGenerator{}
}

// Generate creates a biome with the given parameters
func (bg *BiomeGenerator) Generate(name string, seed int64) *Biome {
	biome := &Biome{
		Name:        name,
		Theme:       name,
		ColorScheme: []string{},
		EnemyTypes:  []string{},
		Hazards:     []string{},
	}
	
	// Set biome-specific properties
	switch name {
	case "cave":
		biome.Temperature = 10
		biome.Moisture = 80
		biome.DangerLevel = 3
		biome.ColorScheme = []string{"#2a2a3a", "#3a3a4a", "#4a4550"}
		biome.EnemyTypes = []string{"bat", "slime", "spider"}
		biome.Hazards = []string{"spike", "pit"}
		
	case "forest":
		biome.Temperature = 20
		biome.Moisture = 70
		biome.DangerLevel = 2
		biome.ColorScheme = []string{"#2a5a2a", "#3a6a3a", "#4a7a4a"}
		biome.EnemyTypes = []string{"wolf", "plant", "insect"}
		biome.Hazards = []string{"thorns", "poison"}
		
	case "ruins":
		biome.Temperature = 15
		biome.Moisture = 40
		biome.DangerLevel = 5
		biome.ColorScheme = []string{"#5a5a50", "#6a6a60", "#7a7a70"}
		biome.EnemyTypes = []string{"golem", "ghost", "construct"}
		biome.Hazards = []string{"trap", "curse"}
		
	case "crystal":
		biome.Temperature = 5
		biome.Moisture = 30
		biome.DangerLevel = 6
		biome.ColorScheme = []string{"#4a5a7a", "#5a6a8a", "#6a7a9a"}
		biome.EnemyTypes = []string{"elemental", "crystal_beast", "wisp"}
		biome.Hazards = []string{"ice", "energy"}
		
	case "abyss":
		biome.Temperature = -10
		biome.Moisture = 20
		biome.DangerLevel = 8
		biome.ColorScheme = []string{"#1a1a2a", "#2a2a3a", "#3a3a4a"}
		biome.EnemyTypes = []string{"shadow", "demon", "horror"}
		biome.Hazards = []string{"void", "corruption"}
		
	case "sky":
		biome.Temperature = 10
		biome.Moisture = 60
		biome.DangerLevel = 4
		biome.ColorScheme = []string{"#6a8aaa", "#7a9aba", "#8aaaca"}
		biome.EnemyTypes = []string{"bird", "cloud_beast", "aerial"}
		biome.Hazards = []string{"wind", "lightning"}
		
	default:
		// Generic biome
		biome.Temperature = 15
		biome.Moisture = 50
		biome.DangerLevel = 3
		biome.ColorScheme = []string{"#4a4a4a", "#5a5a5a", "#6a6a6a"}
		biome.EnemyTypes = []string{"basic_enemy"}
		biome.Hazards = []string{"generic"}
	}
	
	return biome
}

// GetEnvironmentalEffect returns effect based on biome
func (b *Biome) GetEnvironmentalEffect() string {
	if b.Temperature < 0 {
		return "freezing" // Player moves slower
	} else if b.Temperature > 30 {
		return "scorching" // Player takes periodic damage
	} else if b.Moisture > 80 {
		return "slippery" // Reduced friction
	}
	return "normal"
}

// GetMusicMood returns appropriate music mood for biome
func (b *Biome) GetMusicMood() string {
	switch b.Theme {
	case "cave":
		return "dark_ambient"
	case "forest":
		return "peaceful"
	case "ruins":
		return "mysterious"
	case "crystal":
		return "ethereal"
	case "abyss":
		return "horror"
	case "sky":
		return "uplifting"
	default:
		return "neutral"
	}
}
