package world

import (
	"testing"
)

// TestNewBiomeGenerator_CreatesValidGenerator tests generator creation
func TestNewBiomeGenerator_CreatesValidGenerator(t *testing.T) {
	bg := NewBiomeGenerator()

	if bg == nil {
		t.Fatal("NewBiomeGenerator returned nil")
	}
}

// TestGenerate_CaveBiome tests cave biome generation
func TestGenerate_CaveBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("cave", 123)

	if biome == nil {
		t.Fatal("Generate returned nil")
	}

	if biome.Name != "cave" {
		t.Errorf("Expected name 'cave', got '%s'", biome.Name)
	}

	if biome.Theme != "cave" {
		t.Errorf("Expected theme 'cave', got '%s'", biome.Theme)
	}

	if biome.Temperature != 10 {
		t.Errorf("Expected temperature 10, got %d", biome.Temperature)
	}

	if biome.Moisture != 80 {
		t.Errorf("Expected moisture 80, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 3 {
		t.Errorf("Expected danger level 3, got %d", biome.DangerLevel)
	}

	if len(biome.ColorScheme) == 0 {
		t.Error("ColorScheme should not be empty")
	}

	if len(biome.EnemyTypes) == 0 {
		t.Error("EnemyTypes should not be empty")
	}

	if len(biome.Hazards) == 0 {
		t.Error("Hazards should not be empty")
	}
}

// TestGenerate_ForestBiome tests forest biome generation
func TestGenerate_ForestBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("forest", 456)

	if biome.Name != "forest" {
		t.Errorf("Expected name 'forest', got '%s'", biome.Name)
	}

	if biome.Temperature != 20 {
		t.Errorf("Expected temperature 20, got %d", biome.Temperature)
	}

	if biome.Moisture != 70 {
		t.Errorf("Expected moisture 70, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 2 {
		t.Errorf("Expected danger level 2, got %d", biome.DangerLevel)
	}
}

// TestGenerate_RuinsBiome tests ruins biome generation
func TestGenerate_RuinsBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("ruins", 789)

	if biome.Name != "ruins" {
		t.Errorf("Expected name 'ruins', got '%s'", biome.Name)
	}

	if biome.Temperature != 15 {
		t.Errorf("Expected temperature 15, got %d", biome.Temperature)
	}

	if biome.Moisture != 40 {
		t.Errorf("Expected moisture 40, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 5 {
		t.Errorf("Expected danger level 5, got %d", biome.DangerLevel)
	}
}

// TestGenerate_CrystalBiome tests crystal biome generation
func TestGenerate_CrystalBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("crystal", 101)

	if biome.Name != "crystal" {
		t.Errorf("Expected name 'crystal', got '%s'", biome.Name)
	}

	if biome.Temperature != 5 {
		t.Errorf("Expected temperature 5, got %d", biome.Temperature)
	}

	if biome.Moisture != 30 {
		t.Errorf("Expected moisture 30, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 6 {
		t.Errorf("Expected danger level 6, got %d", biome.DangerLevel)
	}
}

// TestGenerate_AbyssBiome tests abyss biome generation
func TestGenerate_AbyssBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("abyss", 202)

	if biome.Name != "abyss" {
		t.Errorf("Expected name 'abyss', got '%s'", biome.Name)
	}

	if biome.Temperature != -10 {
		t.Errorf("Expected temperature -10, got %d", biome.Temperature)
	}

	if biome.Moisture != 20 {
		t.Errorf("Expected moisture 20, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 8 {
		t.Errorf("Expected danger level 8, got %d", biome.DangerLevel)
	}
}

// TestGenerate_SkyBiome tests sky biome generation
func TestGenerate_SkyBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("sky", 303)

	if biome.Name != "sky" {
		t.Errorf("Expected name 'sky', got '%s'", biome.Name)
	}

	if biome.Temperature != 10 {
		t.Errorf("Expected temperature 10, got %d", biome.Temperature)
	}

	if biome.Moisture != 60 {
		t.Errorf("Expected moisture 60, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 4 {
		t.Errorf("Expected danger level 4, got %d", biome.DangerLevel)
	}
}

// TestGenerate_DefaultBiome tests default biome for unknown names
func TestGenerate_DefaultBiome(t *testing.T) {
	bg := NewBiomeGenerator()
	biome := bg.Generate("unknown", 404)

	if biome.Name != "unknown" {
		t.Errorf("Expected name 'unknown', got '%s'", biome.Name)
	}

	if biome.Temperature != 15 {
		t.Errorf("Expected temperature 15, got %d", biome.Temperature)
	}

	if biome.Moisture != 50 {
		t.Errorf("Expected moisture 50, got %d", biome.Moisture)
	}

	if biome.DangerLevel != 3 {
		t.Errorf("Expected danger level 3, got %d", biome.DangerLevel)
	}

	if len(biome.ColorScheme) == 0 {
		t.Error("Default biome should have color scheme")
	}
}

// TestGenerate_AllBiomes_TableDriven tests all biome types
func TestGenerate_AllBiomes_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		biomeName    string
		expectedTemp int
		expectedMois int
		expectedDang int
	}{
		{"cave biome", "cave", 10, 80, 3},
		{"forest biome", "forest", 20, 70, 2},
		{"ruins biome", "ruins", 15, 40, 5},
		{"crystal biome", "crystal", 5, 30, 6},
		{"abyss biome", "abyss", -10, 20, 8},
		{"sky biome", "sky", 10, 60, 4},
		{"default biome", "unknown_type", 15, 50, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := NewBiomeGenerator()
			biome := bg.Generate(tt.biomeName, 1234)

			if biome.Temperature != tt.expectedTemp {
				t.Errorf("Expected temperature %d, got %d", tt.expectedTemp, biome.Temperature)
			}

			if biome.Moisture != tt.expectedMois {
				t.Errorf("Expected moisture %d, got %d", tt.expectedMois, biome.Moisture)
			}

			if biome.DangerLevel != tt.expectedDang {
				t.Errorf("Expected danger level %d, got %d", tt.expectedDang, biome.DangerLevel)
			}
		})
	}
}

// TestGenerate_ConsistentOutput tests that same input produces same output
func TestGenerate_ConsistentOutput(t *testing.T) {
	bg := NewBiomeGenerator()
	
	biome1 := bg.Generate("cave", 555)
	biome2 := bg.Generate("cave", 555)

	if biome1.Temperature != biome2.Temperature {
		t.Error("Temperature should be consistent")
	}

	if biome1.Moisture != biome2.Moisture {
		t.Error("Moisture should be consistent")
	}

	if biome1.DangerLevel != biome2.DangerLevel {
		t.Error("DangerLevel should be consistent")
	}
}

// TestGetEnvironmentalEffect_Freezing tests freezing effect
func TestGetEnvironmentalEffect_Freezing(t *testing.T) {
	biome := &Biome{Temperature: -5}
	effect := biome.GetEnvironmentalEffect()

	if effect != "freezing" {
		t.Errorf("Expected 'freezing', got '%s'", effect)
	}
}

// TestGetEnvironmentalEffect_Scorching tests scorching effect
func TestGetEnvironmentalEffect_Scorching(t *testing.T) {
	biome := &Biome{Temperature: 35}
	effect := biome.GetEnvironmentalEffect()

	if effect != "scorching" {
		t.Errorf("Expected 'scorching', got '%s'", effect)
	}
}

// TestGetEnvironmentalEffect_Slippery tests slippery effect
func TestGetEnvironmentalEffect_Slippery(t *testing.T) {
	biome := &Biome{Temperature: 15, Moisture: 85}
	effect := biome.GetEnvironmentalEffect()

	if effect != "slippery" {
		t.Errorf("Expected 'slippery', got '%s'", effect)
	}
}

// TestGetEnvironmentalEffect_Normal tests normal conditions
func TestGetEnvironmentalEffect_Normal(t *testing.T) {
	biome := &Biome{Temperature: 20, Moisture: 50}
	effect := biome.GetEnvironmentalEffect()

	if effect != "normal" {
		t.Errorf("Expected 'normal', got '%s'", effect)
	}
}

// TestGetEnvironmentalEffect_BoundaryConditions tests boundary values
func TestGetEnvironmentalEffect_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name        string
		temperature int
		moisture    int
		expected    string
	}{
		{"exactly zero temp", 0, 50, "normal"},
		{"exactly 30 temp", 30, 50, "normal"},
		{"just above 30", 31, 50, "scorching"},
		{"just below zero", -1, 50, "freezing"},
		{"exactly 80 moisture", 15, 80, "normal"},
		{"just above 80 moisture", 15, 81, "slippery"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biome := &Biome{Temperature: tt.temperature, Moisture: tt.moisture}
			effect := biome.GetEnvironmentalEffect()

			if effect != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, effect)
			}
		})
	}
}

// TestGetMusicMood_Cave tests cave music mood
func TestGetMusicMood_Cave(t *testing.T) {
	biome := &Biome{Theme: "cave"}
	mood := biome.GetMusicMood()

	if mood != "dark_ambient" {
		t.Errorf("Expected 'dark_ambient', got '%s'", mood)
	}
}

// TestGetMusicMood_Forest tests forest music mood
func TestGetMusicMood_Forest(t *testing.T) {
	biome := &Biome{Theme: "forest"}
	mood := biome.GetMusicMood()

	if mood != "peaceful" {
		t.Errorf("Expected 'peaceful', got '%s'", mood)
	}
}

// TestGetMusicMood_Ruins tests ruins music mood
func TestGetMusicMood_Ruins(t *testing.T) {
	biome := &Biome{Theme: "ruins"}
	mood := biome.GetMusicMood()

	if mood != "mysterious" {
		t.Errorf("Expected 'mysterious', got '%s'", mood)
	}
}

// TestGetMusicMood_Crystal tests crystal music mood
func TestGetMusicMood_Crystal(t *testing.T) {
	biome := &Biome{Theme: "crystal"}
	mood := biome.GetMusicMood()

	if mood != "ethereal" {
		t.Errorf("Expected 'ethereal', got '%s'", mood)
	}
}

// TestGetMusicMood_Abyss tests abyss music mood
func TestGetMusicMood_Abyss(t *testing.T) {
	biome := &Biome{Theme: "abyss"}
	mood := biome.GetMusicMood()

	if mood != "horror" {
		t.Errorf("Expected 'horror', got '%s'", mood)
	}
}

// TestGetMusicMood_Sky tests sky music mood
func TestGetMusicMood_Sky(t *testing.T) {
	biome := &Biome{Theme: "sky"}
	mood := biome.GetMusicMood()

	if mood != "uplifting" {
		t.Errorf("Expected 'uplifting', got '%s'", mood)
	}
}

// TestGetMusicMood_Default tests default music mood
func TestGetMusicMood_Default(t *testing.T) {
	biome := &Biome{Theme: "unknown"}
	mood := biome.GetMusicMood()

	if mood != "neutral" {
		t.Errorf("Expected 'neutral', got '%s'", mood)
	}
}

// TestGetMusicMood_AllThemes tests all theme music moods
func TestGetMusicMood_AllThemes(t *testing.T) {
	tests := []struct {
		theme    string
		expected string
	}{
		{"cave", "dark_ambient"},
		{"forest", "peaceful"},
		{"ruins", "mysterious"},
		{"crystal", "ethereal"},
		{"abyss", "horror"},
		{"sky", "uplifting"},
		{"other", "neutral"},
	}

	for _, tt := range tests {
		t.Run(tt.theme, func(t *testing.T) {
			biome := &Biome{Theme: tt.theme}
			mood := biome.GetMusicMood()

			if mood != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, mood)
			}
		})
	}
}

// TestBiome_ColorSchemeNotEmpty tests that biomes have color schemes
func TestBiome_ColorSchemeNotEmpty(t *testing.T) {
	bg := NewBiomeGenerator()
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}

	for _, biomeName := range biomes {
		t.Run(biomeName, func(t *testing.T) {
			biome := bg.Generate(biomeName, 999)
			if len(biome.ColorScheme) == 0 {
				t.Errorf("Biome '%s' has empty color scheme", biomeName)
			}
		})
	}
}

// TestBiome_EnemyTypesNotEmpty tests that biomes have enemy types
func TestBiome_EnemyTypesNotEmpty(t *testing.T) {
	bg := NewBiomeGenerator()
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}

	for _, biomeName := range biomes {
		t.Run(biomeName, func(t *testing.T) {
			biome := bg.Generate(biomeName, 999)
			if len(biome.EnemyTypes) == 0 {
				t.Errorf("Biome '%s' has empty enemy types", biomeName)
			}
		})
	}
}

// TestBiome_HazardsNotEmpty tests that biomes have hazards
func TestBiome_HazardsNotEmpty(t *testing.T) {
	bg := NewBiomeGenerator()
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}

	for _, biomeName := range biomes {
		t.Run(biomeName, func(t *testing.T) {
			biome := bg.Generate(biomeName, 999)
			if len(biome.Hazards) == 0 {
				t.Errorf("Biome '%s' has empty hazards", biomeName)
			}
		})
	}
}

// TestBiome_DangerLevelRange tests that danger levels are reasonable
func TestBiome_DangerLevelRange(t *testing.T) {
	bg := NewBiomeGenerator()
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}

	for _, biomeName := range biomes {
		t.Run(biomeName, func(t *testing.T) {
			biome := bg.Generate(biomeName, 999)
			if biome.DangerLevel < 1 || biome.DangerLevel > 10 {
				t.Errorf("Biome '%s' has danger level %d outside range [1,10]",
					biomeName, biome.DangerLevel)
			}
		})
	}
}

// TestBiome_TemperatureRange tests that temperatures are reasonable
func TestBiome_TemperatureRange(t *testing.T) {
	bg := NewBiomeGenerator()
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}

	for _, biomeName := range biomes {
		t.Run(biomeName, func(t *testing.T) {
			biome := bg.Generate(biomeName, 999)
			if biome.Temperature < -20 || biome.Temperature > 40 {
				t.Errorf("Biome '%s' has temperature %d outside range [-20,40]",
					biomeName, biome.Temperature)
			}
		})
	}
}

// TestBiome_MoistureRange tests that moisture levels are reasonable
func TestBiome_MoistureRange(t *testing.T) {
	bg := NewBiomeGenerator()
	biomes := []string{"cave", "forest", "ruins", "crystal", "abyss", "sky"}

	for _, biomeName := range biomes {
		t.Run(biomeName, func(t *testing.T) {
			biome := bg.Generate(biomeName, 999)
			if biome.Moisture < 0 || biome.Moisture > 100 {
				t.Errorf("Biome '%s' has moisture %d outside range [0,100]",
					biomeName, biome.Moisture)
			}
		})
	}
}

// TestBiome_SpecificEnemies tests that biomes have appropriate enemies
func TestBiome_SpecificEnemies(t *testing.T) {
	bg := NewBiomeGenerator()
	
	cave := bg.Generate("cave", 1)
	if !contains(cave.EnemyTypes, "bat") {
		t.Error("Cave should have bat enemy")
	}

	forest := bg.Generate("forest", 1)
	if !contains(forest.EnemyTypes, "wolf") {
		t.Error("Forest should have wolf enemy")
	}

	abyss := bg.Generate("abyss", 1)
	if !contains(abyss.EnemyTypes, "shadow") {
		t.Error("Abyss should have shadow enemy")
	}
}

// TestBiome_SpecificHazards tests that biomes have appropriate hazards
func TestBiome_SpecificHazards(t *testing.T) {
	bg := NewBiomeGenerator()
	
	cave := bg.Generate("cave", 1)
	if !contains(cave.Hazards, "spike") {
		t.Error("Cave should have spike hazard")
	}

	crystal := bg.Generate("crystal", 1)
	if !contains(crystal.Hazards, "ice") {
		t.Error("Crystal should have ice hazard")
	}

	sky := bg.Generate("sky", 1)
	if !contains(sky.Hazards, "lightning") {
		t.Error("Sky should have lightning hazard")
	}
}

// TestGetEnvironmentalEffect_PrioritizesFreezing tests freezing takes priority
func TestGetEnvironmentalEffect_PrioritizesFreezing(t *testing.T) {
	// Even with high moisture, freezing should be returned
	biome := &Biome{Temperature: -5, Moisture: 85}
	effect := biome.GetEnvironmentalEffect()

	if effect != "freezing" {
		t.Errorf("Expected 'freezing', got '%s'", effect)
	}
}

// TestGetEnvironmentalEffect_PrioritizesScorching tests scorching takes priority over slippery
func TestGetEnvironmentalEffect_PrioritizesScorching(t *testing.T) {
	// High temp with high moisture should return scorching
	biome := &Biome{Temperature: 35, Moisture: 85}
	effect := biome.GetEnvironmentalEffect()

	if effect != "scorching" {
		t.Errorf("Expected 'scorching', got '%s'", effect)
	}
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
