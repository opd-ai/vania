// Package narrative generates procedural story elements including themes,
// world lore, factions, character backgrounds, and item descriptions to
// create cohesive narrative contexts that influence other generation systems.
package narrative

import (
	"fmt"
	"math/rand"
	"strings"
)

// StoryTheme defines the narrative theme
type StoryTheme string

const (
	FantasyTheme    StoryTheme = "fantasy"
	SciFiTheme      StoryTheme = "scifi"
	HorrorTheme     StoryTheme = "horror"
	MysticalTheme   StoryTheme = "mystical"
	PostApocTheme   StoryTheme = "postapoc"
)

// Mood represents emotional tone
type Mood string

const (
	DarkMood      Mood = "dark"
	HopefulMood   Mood = "hopeful"
	MysteriousMood Mood = "mysterious"
	EpicMood      Mood = "epic"
)

// WorldContext represents the generated world setting
type WorldContext struct {
	Theme            StoryTheme
	Mood             Mood
	CivilizationType string
	Catastrophe      string
	Factions         []Faction
	PlayerMotivation string
	WorldConstraints map[string]interface{}
}

// Faction represents a group in the world
type Faction struct {
	Name         string
	Description  string
	Relationship string // "ally", "enemy", "neutral"
}

// Character represents a generated character
type Character struct {
	Name       string
	Traits     []string
	Motivation string
	Role       string
}

// StoryElement represents a piece of narrative
type StoryElement struct {
	Type        string
	Content     string
	Context     map[string]string
}

// NarrativeGenerator generates procedural stories
type NarrativeGenerator struct {
	rng *rand.Rand
}

// NewNarrativeGenerator creates a new narrative generator
func NewNarrativeGenerator(seed int64) *NarrativeGenerator {
	return &NarrativeGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Generate creates a complete narrative context
func (ng *NarrativeGenerator) Generate(seed int64) *WorldContext {
	ng.rng = rand.New(rand.NewSource(seed))
	
	theme := ng.selectTheme()
	mood := ng.selectMood()
	
	ctx := &WorldContext{
		Theme:            theme,
		Mood:             mood,
		CivilizationType: ng.generateCivilizationType(theme),
		Catastrophe:      ng.generateCatastrophe(theme),
		Factions:         ng.generateFactions(theme, 3),
		PlayerMotivation: ng.generatePlayerMotivation(theme),
		WorldConstraints: make(map[string]interface{}),
	}
	
	// Set constraints based on theme
	ctx.WorldConstraints["dangerLevel"] = ng.rng.Intn(5) + 3
	ctx.WorldConstraints["mysteryLevel"] = ng.rng.Intn(5) + 3
	ctx.WorldConstraints["techLevel"] = ng.getTechLevel(theme)
	
	return ctx
}

// selectTheme randomly chooses a story theme
func (ng *NarrativeGenerator) selectTheme() StoryTheme {
	themes := []StoryTheme{
		FantasyTheme,
		SciFiTheme,
		HorrorTheme,
		MysticalTheme,
		PostApocTheme,
	}
	return themes[ng.rng.Intn(len(themes))]
}

// selectMood randomly chooses a mood
func (ng *NarrativeGenerator) selectMood() Mood {
	moods := []Mood{
		DarkMood,
		HopefulMood,
		MysteriousMood,
		EpicMood,
	}
	return moods[ng.rng.Intn(len(moods))]
}

// generateCivilizationType creates a civilization description
func (ng *NarrativeGenerator) generateCivilizationType(theme StoryTheme) string {
	civilizations := map[StoryTheme][]string{
		FantasyTheme: {
			"ancient elven kingdom",
			"dwarven underground empire",
			"magical city-states",
			"tribal confederation",
		},
		SciFiTheme: {
			"advanced spacefaring society",
			"cybernetic collective",
			"colony of explorers",
			"artificial intelligence network",
		},
		HorrorTheme: {
			"cursed settlement",
			"corrupted monastery",
			"twisted research facility",
			"haunted asylum",
		},
		MysticalTheme: {
			"order of mystics",
			"ethereal plane dwellers",
			"dreamwalker sanctuary",
			"crystal city",
		},
		PostApocTheme: {
			"survivor encampment",
			"wasteland traders",
			"military remnants",
			"mutant tribes",
		},
	}
	
	options := civilizations[theme]
	return options[ng.rng.Intn(len(options))]
}

// generateCatastrophe creates a disaster event
func (ng *NarrativeGenerator) generateCatastrophe(theme StoryTheme) string {
	catastrophes := map[StoryTheme][]string{
		FantasyTheme: {
			"the Dark Lord's curse shattered the realm",
			"a forbidden ritual tore reality asunder",
			"the Great Dragons vanished, leaving chaos",
			"an ancient evil awakened from its slumber",
		},
		SciFiTheme: {
			"the AI rebellion destroyed civilization",
			"a dimensional rift collapsed the colonies",
			"the alien invasion decimated humanity",
			"quantum experiments fractured spacetime",
		},
		HorrorTheme: {
			"madness spread from the ancient ruins",
			"the dead rose to consume the living",
			"eldritch beings breached our world",
			"a plague transformed people into monsters",
		},
		MysticalTheme: {
			"the Spirit Realm merged with reality",
			"cosmic alignment shattered the barriers",
			"a prophet's vision came to pass",
			"the Void consumed the light",
		},
		PostApocTheme: {
			"nuclear war left only ruins",
			"biological weapons wiped out billions",
			"environmental collapse ended civilization",
			"the machines turned against their makers",
		},
	}
	
	options := catastrophes[theme]
	return options[ng.rng.Intn(len(options))]
}

// generateFactions creates factions for the world
func (ng *NarrativeGenerator) generateFactions(theme StoryTheme, count int) []Faction {
	factions := make([]Faction, count)
	
	nameTemplates := map[StoryTheme][]string{
		FantasyTheme:  {"Order of", "Brotherhood of", "Circle of", "Guild of"},
		SciFiTheme:    {"The", "Sector", "Division", "Protocol"},
		HorrorTheme:   {"Cult of", "Children of", "Followers of", "Sect of"},
		MysticalTheme: {"Seekers of", "Keepers of", "Guardians of", "Watchers of"},
		PostApocTheme: {"The", "New", "Free", "United"},
	}
	
	nouns := map[StoryTheme][]string{
		FantasyTheme:  {"the Phoenix", "the Silver Moon", "the Iron Crown", "the Ancient Oak"},
		SciFiTheme:    {"Nexus", "Genesis", "Vanguard", "Horizon"},
		HorrorTheme:   {"the Crimson Eye", "Eternal Night", "the Whispering Dark", "the Void"},
		MysticalTheme: {"True Sight", "the Eternal Flame", "Cosmic Balance", "Hidden Knowledge"},
		PostApocTheme: {"Survivors", "Resistance", "Haven", "Outcasts"},
	}
	
	relationships := []string{"ally", "enemy", "neutral"}
	
	templates := nameTemplates[theme]
	nounList := nouns[theme]
	
	for i := 0; i < count; i++ {
		template := templates[ng.rng.Intn(len(templates))]
		noun := nounList[ng.rng.Intn(len(nounList))]
		
		factions[i] = Faction{
			Name:         fmt.Sprintf("%s %s", template, noun),
			Description:  ng.generateFactionDescription(theme),
			Relationship: relationships[ng.rng.Intn(len(relationships))],
		}
	}
	
	return factions
}

// generateFactionDescription creates faction description
func (ng *NarrativeGenerator) generateFactionDescription(theme StoryTheme) string {
	descriptions := []string{
		"seeks to restore the old ways",
		"controls vital resources",
		"guards ancient secrets",
		"wages war against the darkness",
		"manipulates events from the shadows",
	}
	return descriptions[ng.rng.Intn(len(descriptions))]
}

// generatePlayerMotivation creates player goal
func (ng *NarrativeGenerator) generatePlayerMotivation(theme StoryTheme) string {
	motivations := map[StoryTheme][]string{
		FantasyTheme: {
			"break the curse on your homeland",
			"recover the stolen artifacts of power",
			"defeat the tyrant who destroyed your village",
			"find the legendary weapon to save the kingdom",
		},
		SciFiTheme: {
			"escape the dying space station",
			"uncover the conspiracy behind the colony collapse",
			"rescue your crew from alien captivity",
			"prevent the AI from destroying humanity",
		},
		HorrorTheme: {
			"escape the nightmare realm alive",
			"break the curse binding you to this place",
			"discover what happened to the missing researchers",
			"stop the ritual before it's too late",
		},
		MysticalTheme: {
			"restore balance to the Spirit Realm",
			"awaken your true mystical potential",
			"seal the breach between worlds",
			"fulfill the ancient prophecy",
		},
		PostApocTheme: {
			"find a safe haven for your people",
			"locate the rumored untouched vault",
			"avenge your fallen community",
			"discover the truth about the apocalypse",
		},
	}
	
	options := motivations[theme]
	return options[ng.rng.Intn(len(options))]
}

// getTechLevel returns technology level for theme
func (ng *NarrativeGenerator) getTechLevel(theme StoryTheme) int {
	techLevels := map[StoryTheme]int{
		FantasyTheme:  2, // Medieval
		SciFiTheme:    9, // Advanced
		HorrorTheme:   4, // Modern
		MysticalTheme: 3, // Ancient magic
		PostApocTheme: 5, // Scavenged tech
	}
	return techLevels[theme]
}

// GenerateCharacter creates a character
func (ng *NarrativeGenerator) GenerateCharacter(role string) *Character {
	firstNames := []string{
		"Aria", "Kael", "Zara", "Theron", "Lyra",
		"Drake", "Nova", "Cipher", "Echo", "Raven",
	}
	
	lastNames := []string{
		"Shadowbane", "Ironheart", "Stormwind", "Nightshade",
		"Brightblade", "Darkwater", "Swiftfoot", "Firebrand",
	}
	
	traits := []string{
		"brave", "cunning", "wise", "fierce", "mysterious",
		"loyal", "ambitious", "skilled", "ruthless", "compassionate",
	}
	
	motivations := []string{
		"seeking redemption",
		"protecting loved ones",
		"pursuing knowledge",
		"gaining power",
		"escaping the past",
	}
	
	firstName := firstNames[ng.rng.Intn(len(firstNames))]
	lastName := lastNames[ng.rng.Intn(len(lastNames))]
	
	numTraits := 2 + ng.rng.Intn(2)
	selectedTraits := make([]string, numTraits)
	for i := 0; i < numTraits; i++ {
		selectedTraits[i] = traits[ng.rng.Intn(len(traits))]
	}
	
	return &Character{
		Name:       fmt.Sprintf("%s %s", firstName, lastName),
		Traits:     selectedTraits,
		Motivation: motivations[ng.rng.Intn(len(motivations))],
		Role:       role,
	}
}

// GenerateItemDescription generates item lore
func (ng *NarrativeGenerator) GenerateItemDescription(itemType string, theme StoryTheme) string {
	adjectives := map[StoryTheme][]string{
		FantasyTheme:  {"enchanted", "ancient", "blessed", "legendary"},
		SciFiTheme:    {"advanced", "prototype", "quantum", "neural"},
		HorrorTheme:   {"cursed", "twisted", "forbidden", "eldritch"},
		MysticalTheme: {"ethereal", "transcendent", "sacred", "cosmic"},
		PostApocTheme: {"salvaged", "modified", "reinforced", "makeshift"},
	}
	
	templates := map[string][]string{
		"weapon": {
			"A %s blade forged in the fires of the %s.",
			"This %s weapon has seen countless battles.",
			"The %s craftsmanship is evident in every detail.",
		},
		"key_item": {
			"A %s artifact of immense power.",
			"This %s object holds the key to secrets long forgotten.",
			"The %s nature of this item is unmistakable.",
		},
		"consumable": {
			"A %s potion that glows with inner light.",
			"This %s elixir was crafted by master alchemists.",
			"The %s properties make it invaluable.",
		},
	}
	
	// Use default if theme not found
	adjList, ok := adjectives[theme]
	if !ok || len(adjList) == 0 {
		adjList = []string{"mysterious", "powerful", "rare", "valuable"}
	}
	
	// Use default if item type not found
	tmplList, ok := templates[itemType]
	if !ok || len(tmplList) == 0 {
		return "A remarkable item of unknown origin."
	}
	
	adj := adjList[ng.rng.Intn(len(adjList))]
	tmpl := tmplList[ng.rng.Intn(len(tmplList))]
	
	if strings.Contains(tmpl, "fires of the %s") {
		locations := []string{"ancients", "fallen kingdom", "first age", "old world"}
		return fmt.Sprintf(tmpl, adj, locations[ng.rng.Intn(len(locations))])
	}
	
	return fmt.Sprintf(tmpl, adj)
}

// GenerateRoomDescription generates room description
func (ng *NarrativeGenerator) GenerateRoomDescription(roomType string, theme StoryTheme) string {
	descriptions := map[string][]string{
		"combat": {
			"The chamber echoes with the sounds of battle.",
			"Danger lurks in every shadow of this arena.",
			"Ancient weapons line the walls of this proving ground.",
		},
		"treasure": {
			"Glittering prizes await those brave enough to claim them.",
			"The air shimmers with the promise of riches.",
			"Valuable artifacts rest on ornate pedestals.",
		},
		"puzzle": {
			"Strange mechanisms hint at hidden solutions.",
			"Cryptic symbols cover every surface.",
			"The room holds secrets waiting to be unraveled.",
		},
	}
	
	if descs, ok := descriptions[roomType]; ok {
		return descs[ng.rng.Intn(len(descs))]
	}
	
	return "A mysterious chamber awaits exploration."
}
