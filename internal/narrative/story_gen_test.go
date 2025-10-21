package narrative

import (
	"testing"
)

// TestNewNarrativeGenerator_CreatesValidGenerator tests generator creation
func TestNewNarrativeGenerator_CreatesValidGenerator(t *testing.T) {
	seed := int64(12345)
	ng := NewNarrativeGenerator(seed)

	if ng == nil {
		t.Fatal("NewNarrativeGenerator returned nil")
	}
	if ng.rng == nil {
		t.Error("NarrativeGenerator.rng is nil")
	}
}

// TestGenerate_CreatesCompleteWorldContext tests full world generation
func TestGenerate_CreatesCompleteWorldContext(t *testing.T) {
	tests := []struct {
		name string
		seed int64
	}{
		{"seed 1", 1},
		{"seed 42", 42},
		{"seed 100", 100},
		{"seed 999", 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ng := NewNarrativeGenerator(tt.seed)
			ctx := ng.Generate(tt.seed)

			if ctx == nil {
				t.Fatal("Generate returned nil WorldContext")
			}

			// Verify theme is set
			if ctx.Theme == "" {
				t.Error("WorldContext.Theme is empty")
			}

			// Verify mood is set
			if ctx.Mood == "" {
				t.Error("WorldContext.Mood is empty")
			}

			// Verify civilization type is set
			if ctx.CivilizationType == "" {
				t.Error("WorldContext.CivilizationType is empty")
			}

			// Verify catastrophe is set
			if ctx.Catastrophe == "" {
				t.Error("WorldContext.Catastrophe is empty")
			}

			// Verify player motivation is set
			if ctx.PlayerMotivation == "" {
				t.Error("WorldContext.PlayerMotivation is empty")
			}

			// Verify factions are generated
			if len(ctx.Factions) != 3 {
				t.Errorf("Expected 3 factions, got %d", len(ctx.Factions))
			}

			// Verify world constraints are set
			if ctx.WorldConstraints == nil {
				t.Error("WorldConstraints is nil")
			}

			// Check specific constraints exist
			if _, ok := ctx.WorldConstraints["dangerLevel"]; !ok {
				t.Error("WorldConstraints missing dangerLevel")
			}
			if _, ok := ctx.WorldConstraints["mysteryLevel"]; !ok {
				t.Error("WorldConstraints missing mysteryLevel")
			}
			if _, ok := ctx.WorldConstraints["techLevel"]; !ok {
				t.Error("WorldConstraints missing techLevel")
			}
		})
	}
}

// TestGenerate_DeterministicOutput tests that same seed produces same output
func TestGenerate_DeterministicOutput(t *testing.T) {
	seed := int64(54321)
	
	ng1 := NewNarrativeGenerator(seed)
	ctx1 := ng1.Generate(seed)
	
	ng2 := NewNarrativeGenerator(seed)
	ctx2 := ng2.Generate(seed)

	if ctx1.Theme != ctx2.Theme {
		t.Errorf("Theme mismatch: %s != %s", ctx1.Theme, ctx2.Theme)
	}
	if ctx1.Mood != ctx2.Mood {
		t.Errorf("Mood mismatch: %s != %s", ctx1.Mood, ctx2.Mood)
	}
	if ctx1.CivilizationType != ctx2.CivilizationType {
		t.Errorf("CivilizationType mismatch: %s != %s", ctx1.CivilizationType, ctx2.CivilizationType)
	}
	if ctx1.Catastrophe != ctx2.Catastrophe {
		t.Errorf("Catastrophe mismatch: %s != %s", ctx1.Catastrophe, ctx2.Catastrophe)
	}
	if ctx1.PlayerMotivation != ctx2.PlayerMotivation {
		t.Errorf("PlayerMotivation mismatch: %s != %s", ctx1.PlayerMotivation, ctx2.PlayerMotivation)
	}
}

// TestSelectTheme_ReturnsValidTheme tests theme selection
func TestSelectTheme_ReturnsValidTheme(t *testing.T) {
	ng := NewNarrativeGenerator(123)
	validThemes := map[StoryTheme]bool{
		FantasyTheme:  true,
		SciFiTheme:    true,
		HorrorTheme:   true,
		MysticalTheme: true,
		PostApocTheme: true,
	}

	// Test multiple selections to ensure variety
	for i := 0; i < 20; i++ {
		theme := ng.selectTheme()
		if !validThemes[theme] {
			t.Errorf("selectTheme returned invalid theme: %s", theme)
		}
	}
}

// TestSelectMood_ReturnsValidMood tests mood selection
func TestSelectMood_ReturnsValidMood(t *testing.T) {
	ng := NewNarrativeGenerator(456)
	validMoods := map[Mood]bool{
		DarkMood:       true,
		HopefulMood:    true,
		MysteriousMood: true,
		EpicMood:       true,
	}

	// Test multiple selections to ensure variety
	for i := 0; i < 20; i++ {
		mood := ng.selectMood()
		if !validMoods[mood] {
			t.Errorf("selectMood returned invalid mood: %s", mood)
		}
	}
}

// TestGenerateCivilizationType_AllThemes tests civilization generation for all themes
func TestGenerateCivilizationType_AllThemes(t *testing.T) {
	tests := []struct {
		theme StoryTheme
	}{
		{FantasyTheme},
		{SciFiTheme},
		{HorrorTheme},
		{MysticalTheme},
		{PostApocTheme},
	}

	for _, tt := range tests {
		t.Run(string(tt.theme), func(t *testing.T) {
			ng := NewNarrativeGenerator(789)
			civ := ng.generateCivilizationType(tt.theme)

			if civ == "" {
				t.Errorf("generateCivilizationType returned empty string for theme %s", tt.theme)
			}
		})
	}
}

// TestGenerateCatastrophe_AllThemes tests catastrophe generation for all themes
func TestGenerateCatastrophe_AllThemes(t *testing.T) {
	tests := []struct {
		theme StoryTheme
	}{
		{FantasyTheme},
		{SciFiTheme},
		{HorrorTheme},
		{MysticalTheme},
		{PostApocTheme},
	}

	for _, tt := range tests {
		t.Run(string(tt.theme), func(t *testing.T) {
			ng := NewNarrativeGenerator(101)
			catastrophe := ng.generateCatastrophe(tt.theme)

			if catastrophe == "" {
				t.Errorf("generateCatastrophe returned empty string for theme %s", tt.theme)
			}
		})
	}
}

// TestGenerateFactions_CreatesCorrectCount tests faction generation
func TestGenerateFactions_CreatesCorrectCount(t *testing.T) {
	tests := []struct {
		name  string
		theme StoryTheme
		count int
	}{
		{"fantasy 3", FantasyTheme, 3},
		{"scifi 5", SciFiTheme, 5},
		{"horror 2", HorrorTheme, 2},
		{"mystical 4", MysticalTheme, 4},
		{"postapoc 1", PostApocTheme, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ng := NewNarrativeGenerator(202)
			factions := ng.generateFactions(tt.theme, tt.count)

			if len(factions) != tt.count {
				t.Errorf("Expected %d factions, got %d", tt.count, len(factions))
			}

			// Verify each faction has required fields
			for i, faction := range factions {
				if faction.Name == "" {
					t.Errorf("Faction %d has empty name", i)
				}
				if faction.Description == "" {
					t.Errorf("Faction %d has empty description", i)
				}
				if faction.Relationship == "" {
					t.Errorf("Faction %d has empty relationship", i)
				}

				// Verify relationship is valid
				validRels := map[string]bool{"ally": true, "enemy": true, "neutral": true}
				if !validRels[faction.Relationship] {
					t.Errorf("Faction %d has invalid relationship: %s", i, faction.Relationship)
				}
			}
		})
	}
}

// TestGenerateFactionDescription_ReturnsNonEmpty tests faction description generation
func TestGenerateFactionDescription_ReturnsNonEmpty(t *testing.T) {
	ng := NewNarrativeGenerator(303)
	
	for i := 0; i < 10; i++ {
		desc := ng.generateFactionDescription(FantasyTheme)
		if desc == "" {
			t.Error("generateFactionDescription returned empty string")
		}
	}
}

// TestGeneratePlayerMotivation_AllThemes tests player motivation generation
func TestGeneratePlayerMotivation_AllThemes(t *testing.T) {
	tests := []struct {
		theme StoryTheme
	}{
		{FantasyTheme},
		{SciFiTheme},
		{HorrorTheme},
		{MysticalTheme},
		{PostApocTheme},
	}

	for _, tt := range tests {
		t.Run(string(tt.theme), func(t *testing.T) {
			ng := NewNarrativeGenerator(404)
			motivation := ng.generatePlayerMotivation(tt.theme)

			if motivation == "" {
				t.Errorf("generatePlayerMotivation returned empty string for theme %s", tt.theme)
			}
		})
	}
}

// TestGetTechLevel_AllThemes tests tech level assignment
func TestGetTechLevel_AllThemes(t *testing.T) {
	tests := []struct {
		theme         StoryTheme
		expectedLevel int
	}{
		{FantasyTheme, 2},
		{SciFiTheme, 9},
		{HorrorTheme, 4},
		{MysticalTheme, 3},
		{PostApocTheme, 5},
	}

	for _, tt := range tests {
		t.Run(string(tt.theme), func(t *testing.T) {
			ng := NewNarrativeGenerator(505)
			level := ng.getTechLevel(tt.theme)

			if level != tt.expectedLevel {
				t.Errorf("Expected tech level %d for theme %s, got %d",
					tt.expectedLevel, tt.theme, level)
			}
		})
	}
}

// TestGenerateCharacter_CreatesValidCharacter tests character generation
func TestGenerateCharacter_CreatesValidCharacter(t *testing.T) {
	tests := []struct {
		name string
		role string
	}{
		{"merchant role", "merchant"},
		{"guard role", "guard"},
		{"wizard role", "wizard"},
		{"companion role", "companion"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ng := NewNarrativeGenerator(606)
			char := ng.GenerateCharacter(tt.role)

			if char == nil {
				t.Fatal("GenerateCharacter returned nil")
			}

			if char.Name == "" {
				t.Error("Character has empty name")
			}

			if char.Role != tt.role {
				t.Errorf("Expected role %s, got %s", tt.role, char.Role)
			}

			if len(char.Traits) < 2 || len(char.Traits) > 3 {
				t.Errorf("Expected 2-3 traits, got %d", len(char.Traits))
			}

			if char.Motivation == "" {
				t.Error("Character has empty motivation")
			}
		})
	}
}

// TestGenerateCharacter_DeterministicOutput tests character generation determinism
func TestGenerateCharacter_DeterministicOutput(t *testing.T) {
	seed := int64(707)
	role := "hero"
	
	ng1 := NewNarrativeGenerator(seed)
	char1 := ng1.GenerateCharacter(role)
	
	ng2 := NewNarrativeGenerator(seed)
	char2 := ng2.GenerateCharacter(role)

	if char1.Name != char2.Name {
		t.Errorf("Character names differ: %s != %s", char1.Name, char2.Name)
	}
	if char1.Motivation != char2.Motivation {
		t.Errorf("Character motivations differ: %s != %s", char1.Motivation, char2.Motivation)
	}
}

// TestGenerateItemDescription_AllTypes tests item description generation
func TestGenerateItemDescription_AllTypes(t *testing.T) {
	tests := []struct {
		itemType string
		theme    StoryTheme
	}{
		{"weapon", FantasyTheme},
		{"weapon", SciFiTheme},
		{"key_item", HorrorTheme},
		{"key_item", MysticalTheme},
		{"consumable", PostApocTheme},
		{"consumable", FantasyTheme},
	}

	for _, tt := range tests {
		t.Run(tt.itemType+"_"+string(tt.theme), func(t *testing.T) {
			ng := NewNarrativeGenerator(808)
			desc := ng.GenerateItemDescription(tt.itemType, tt.theme)

			if desc == "" {
				t.Errorf("GenerateItemDescription returned empty string for %s/%s",
					tt.itemType, tt.theme)
			}
		})
	}
}

// TestGenerateItemDescription_UnknownType tests handling of unknown item types
func TestGenerateItemDescription_UnknownType(t *testing.T) {
	ng := NewNarrativeGenerator(909)
	desc := ng.GenerateItemDescription("unknown_type", FantasyTheme)

	if desc == "" {
		t.Error("GenerateItemDescription returned empty string for unknown type")
	}

	expected := "A remarkable item of unknown origin."
	if desc != expected {
		t.Errorf("Expected default description, got: %s", desc)
	}
}

// TestGenerateItemDescription_UnknownTheme tests handling of unknown themes
func TestGenerateItemDescription_UnknownTheme(t *testing.T) {
	ng := NewNarrativeGenerator(1010)
	// Use a theme that doesn't exist in the adjectives map
	unknownTheme := StoryTheme("unknown")
	desc := ng.GenerateItemDescription("weapon", unknownTheme)

	if desc == "" {
		t.Error("GenerateItemDescription returned empty string for unknown theme")
	}
}

// TestGenerateRoomDescription_AllTypes tests room description generation
func TestGenerateRoomDescription_AllTypes(t *testing.T) {
	tests := []struct {
		roomType string
		theme    StoryTheme
	}{
		{"combat", FantasyTheme},
		{"treasure", SciFiTheme},
		{"puzzle", HorrorTheme},
	}

	for _, tt := range tests {
		t.Run(tt.roomType+"_"+string(tt.theme), func(t *testing.T) {
			ng := NewNarrativeGenerator(1111)
			desc := ng.GenerateRoomDescription(tt.roomType, tt.theme)

			if desc == "" {
				t.Errorf("GenerateRoomDescription returned empty string for %s/%s",
					tt.roomType, tt.theme)
			}
		})
	}
}

// TestGenerateRoomDescription_UnknownType tests handling of unknown room types
func TestGenerateRoomDescription_UnknownType(t *testing.T) {
	ng := NewNarrativeGenerator(1212)
	desc := ng.GenerateRoomDescription("unknown_room", FantasyTheme)

	expected := "A mysterious chamber awaits exploration."
	if desc != expected {
		t.Errorf("Expected default description, got: %s", desc)
	}
}

// TestWorldContext_ConstraintRanges tests that world constraints are within expected ranges
func TestWorldContext_ConstraintRanges(t *testing.T) {
	ng := NewNarrativeGenerator(1313)
	
	// Generate multiple contexts to check ranges
	for i := 0; i < 20; i++ {
		ctx := ng.Generate(int64(1313 + i))

		dangerLevel, ok := ctx.WorldConstraints["dangerLevel"].(int)
		if !ok {
			t.Error("dangerLevel is not an int")
		}
		if dangerLevel < 3 || dangerLevel > 7 {
			t.Errorf("dangerLevel %d out of expected range [3,7]", dangerLevel)
		}

		mysteryLevel, ok := ctx.WorldConstraints["mysteryLevel"].(int)
		if !ok {
			t.Error("mysteryLevel is not an int")
		}
		if mysteryLevel < 3 || mysteryLevel > 7 {
			t.Errorf("mysteryLevel %d out of expected range [3,7]", mysteryLevel)
		}

		techLevel, ok := ctx.WorldConstraints["techLevel"].(int)
		if !ok {
			t.Error("techLevel is not an int")
		}
		if techLevel < 2 || techLevel > 9 {
			t.Errorf("techLevel %d out of expected range [2,9]", techLevel)
		}
	}
}

// TestFactionVariety_MultipleCalls tests that faction generation produces variety
func TestFactionVariety_MultipleCalls(t *testing.T) {
	ng := NewNarrativeGenerator(1414)
	
	// Generate multiple faction sets
	factionNames := make(map[string]bool)
	for i := 0; i < 10; i++ {
		factions := ng.generateFactions(FantasyTheme, 3)
		for _, faction := range factions {
			factionNames[faction.Name] = true
		}
	}

	// We should see more than 3 unique names across 10 generations
	if len(factionNames) <= 3 {
		t.Errorf("Expected more variety in faction names, got only %d unique names", len(factionNames))
	}
}

// TestCharacterTraitsVariety tests that character traits show variety
func TestCharacterTraitsVariety(t *testing.T) {
	ng := NewNarrativeGenerator(1515)
	
	allTraits := make(map[string]bool)
	for i := 0; i < 10; i++ {
		char := ng.GenerateCharacter("test")
		for _, trait := range char.Traits {
			allTraits[trait] = true
		}
	}

	// Should have variety in traits
	if len(allTraits) < 3 {
		t.Errorf("Expected more variety in character traits, got only %d unique traits", len(allTraits))
	}
}

// TestStoryThemeConstants_ValidValues tests that theme constants are properly defined
func TestStoryThemeConstants_ValidValues(t *testing.T) {
	themes := []StoryTheme{FantasyTheme, SciFiTheme, HorrorTheme, MysticalTheme, PostApocTheme}
	
	for _, theme := range themes {
		if theme == "" {
			t.Error("Found empty theme constant")
		}
	}
}

// TestMoodConstants_ValidValues tests that mood constants are properly defined
func TestMoodConstants_ValidValues(t *testing.T) {
	moods := []Mood{DarkMood, HopefulMood, MysteriousMood, EpicMood}
	
	for _, mood := range moods {
		if mood == "" {
			t.Error("Found empty mood constant")
		}
	}
}

// TestGenerateSeedConsistency_MultipleGenerations tests seed consistency
func TestGenerateSeedConsistency_MultipleGenerations(t *testing.T) {
	seed := int64(1616)
	ng := NewNarrativeGenerator(seed)
	
	// Generate multiple contexts with same seed - should be identical
	ctx1 := ng.Generate(seed)
	ctx2 := ng.Generate(seed)
	
	if ctx1.Theme != ctx2.Theme {
		t.Error("Theme inconsistent with same seed")
	}
	if ctx1.Mood != ctx2.Mood {
		t.Error("Mood inconsistent with same seed")
	}
}
