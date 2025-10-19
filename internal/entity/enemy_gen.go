// Package entity generates procedural enemies, bosses, items, and player
// abilities with stats scaled to danger levels, behavior patterns, attack
// types, and Metroidvania-style ability progression systems.
package entity

import (
	"math/rand"
)

// Enemy represents a generated enemy
type Enemy struct {
	Name          string
	Health        int
	Damage        int
	Speed         float64
	Size          EnemySize
	Behavior      BehaviorPattern
	AttackType    AttackType
	SpriteData    interface{} // Will hold generated sprite
	SoundData     interface{} // Will hold generated sounds
	DangerLevel   int
	BiomeType     string
}

// EnemySize defines enemy dimensions
type EnemySize int

const (
	SmallEnemy EnemySize = iota  // 16x16
	MediumEnemy                   // 32x32
	LargeEnemy                    // 64x64
	BossEnemy                     // 128x128
)

// BehaviorPattern defines movement patterns
type BehaviorPattern int

const (
	PatrolBehavior BehaviorPattern = iota
	ChaseBehavior
	FleeBehavior
	StationaryBehavior
	FlyingBehavior
	JumpingBehavior
)

// AttackType defines attack methods
type AttackType int

const (
	MeleeAttack AttackType = iota
	RangedAttack
	AreaAttack
	ContactDamage
)

// Boss represents a boss enemy
type Boss struct {
	Enemy
	Phases        []BossPhase
	UniqueAttacks []string
	ArenaLayout   interface{}
}

// BossPhase represents a phase of a boss fight
type BossPhase struct {
	HealthThreshold float64 // Percentage when phase activates
	Behavior        BehaviorPattern
	AttackPattern   string
	SpeedModifier   float64
}

// Item represents a collectible item
type Item struct {
	Name        string
	Type        ItemType
	Description string
	Effect      string
	Value       int
	SpriteData  interface{}
}

// ItemType defines item categories
type ItemType int

const (
	WeaponItem ItemType = iota
	ConsumableItem
	KeyItem
	UpgradeItem
	CurrencyItem
)

// Ability represents a player ability/upgrade
type Ability struct {
	Name        string
	Type        AbilityType
	Description string
	UnlockOrder int
}

// AbilityType defines ability categories
type AbilityType int

const (
	MovementAbility AbilityType = iota
	CombatAbility
	UtilityAbility
)

// EnemyGenerator generates procedural enemies
type EnemyGenerator struct {
	rng *rand.Rand
}

// NewEnemyGenerator creates a new enemy generator
func NewEnemyGenerator(seed int64) *EnemyGenerator {
	return &EnemyGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Generate creates an enemy for a biome and danger level
func (eg *EnemyGenerator) Generate(biome string, dangerLevel int, seed int64) *Enemy {
	eg.rng = rand.New(rand.NewSource(seed))
	
	enemy := &Enemy{
		Name:        eg.generateName(biome),
		BiomeType:   biome,
		DangerLevel: dangerLevel,
	}
	
	// Scale stats based on danger level
	baseHealth := 10 + dangerLevel*5
	baseDamage := 5 + dangerLevel*2
	
	enemy.Health = baseHealth + eg.rng.Intn(baseHealth/2)
	enemy.Damage = baseDamage + eg.rng.Intn(baseDamage/2)
	enemy.Speed = 1.0 + float64(dangerLevel)*0.1 + eg.rng.Float64()*0.3
	
	// Assign size
	sizeRoll := eg.rng.Float64()
	if sizeRoll < 0.6 {
		enemy.Size = SmallEnemy
	} else if sizeRoll < 0.9 {
		enemy.Size = MediumEnemy
	} else {
		enemy.Size = LargeEnemy
	}
	
	// Assign behavior based on biome
	enemy.Behavior = eg.selectBehavior(biome)
	
	// Assign attack type
	enemy.AttackType = eg.selectAttackType(enemy.Size)
	
	return enemy
}

// generateName creates an enemy name
func (eg *EnemyGenerator) generateName(biome string) string {
	prefixes := map[string][]string{
		"cave":    {"Shadow", "Stone", "Dark", "Deep"},
		"forest":  {"Wild", "Feral", "Primal", "Ancient"},
		"ruins":   {"Cursed", "Haunted", "Lost", "Fallen"},
		"crystal": {"Shattered", "Gleaming", "Frozen", "Radiant"},
		"abyss":   {"Void", "Abyssal", "Corrupted", "Nightmare"},
		"sky":     {"Sky", "Storm", "Cloud", "Wind"},
	}
	
	suffixes := []string{
		"Crawler", "Beast", "Horror", "Fiend", "Wraith",
		"Spawn", "Stalker", "Guardian", "Shade", "Terror",
	}
	
	// Use default if biome not found
	prefixList, ok := prefixes[biome]
	if !ok || len(prefixList) == 0 {
		prefixList = []string{"Unknown", "Strange", "Mysterious", "Enigmatic"}
	}
	
	prefix := prefixList[eg.rng.Intn(len(prefixList))]
	suffix := suffixes[eg.rng.Intn(len(suffixes))]
	
	return prefix + " " + suffix
}

// selectBehavior chooses behavior pattern
func (eg *EnemyGenerator) selectBehavior(biome string) BehaviorPattern {
	behaviors := map[string][]BehaviorPattern{
		"cave":    {PatrolBehavior, ChaseBehavior, StationaryBehavior},
		"forest":  {PatrolBehavior, ChaseBehavior, FleeBehavior},
		"ruins":   {PatrolBehavior, StationaryBehavior},
		"crystal": {FlyingBehavior, ChaseBehavior},
		"abyss":   {ChaseBehavior, FlyingBehavior},
		"sky":     {FlyingBehavior, PatrolBehavior},
	}
	
	// Use default behavior if biome not found
	options, ok := behaviors[biome]
	if !ok || len(options) == 0 {
		options = []BehaviorPattern{PatrolBehavior, ChaseBehavior}
	}
	
	return options[eg.rng.Intn(len(options))]
}

// selectAttackType chooses attack method
func (eg *EnemyGenerator) selectAttackType(size EnemySize) AttackType {
	switch size {
	case SmallEnemy:
		return ContactDamage
	case MediumEnemy:
		if eg.rng.Float64() < 0.5 {
			return MeleeAttack
		}
		return RangedAttack
	case LargeEnemy:
		return AreaAttack
	default:
		return MeleeAttack
	}
}

// BossGenerator generates boss enemies
type BossGenerator struct {
	rng *rand.Rand
}

// NewBossGenerator creates a new boss generator
func NewBossGenerator(seed int64) *BossGenerator {
	return &BossGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Generate creates a boss enemy
func (bg *BossGenerator) Generate(biome string, seed int64) *Boss {
	bg.rng = rand.New(rand.NewSource(seed))
	
	// Create base enemy with boss stats
	baseEnemy := Enemy{
		Name:        bg.generateBossName(biome),
		Health:      200 + bg.rng.Intn(100),
		Damage:      20 + bg.rng.Intn(10),
		Speed:       0.8 + bg.rng.Float64()*0.4,
		Size:        BossEnemy,
		BiomeType:   biome,
		DangerLevel: 10,
	}
	
	boss := &Boss{
		Enemy:  baseEnemy,
		Phases: make([]BossPhase, 2+bg.rng.Intn(2)),
	}
	
	// Generate phases
	for i := range boss.Phases {
		threshold := 1.0 - float64(i+1)/float64(len(boss.Phases)+1)
		boss.Phases[i] = BossPhase{
			HealthThreshold: threshold,
			Behavior:        BehaviorPattern(bg.rng.Intn(int(JumpingBehavior) + 1)),
			AttackPattern:   bg.generateAttackPattern(),
			SpeedModifier:   1.0 + float64(i)*0.2,
		}
	}
	
	// Generate unique attacks
	attackCount := 2 + bg.rng.Intn(2)
	boss.UniqueAttacks = make([]string, attackCount)
	for i := range boss.UniqueAttacks {
		boss.UniqueAttacks[i] = bg.generateUniqueAttack(biome)
	}
	
	return boss
}

// generateBossName creates a boss name
func (bg *BossGenerator) generateBossName(biome string) string {
	titles := []string{"Lord", "Master", "King", "Queen", "Overlord", "Eternal"}
	
	names := map[string][]string{
		"cave":    {"Darkness", "Depths", "Stone", "Shadow"},
		"forest":  {"Thorns", "Wild", "Green", "Nature"},
		"ruins":   {"Ashes", "Ruin", "Lost", "Forgotten"},
		"crystal": {"Frost", "Crystal", "Prism", "Shard"},
		"abyss":   {"Void", "Nightmare", "Terror", "Despair"},
		"sky":     {"Storms", "Winds", "Skies", "Clouds"},
	}
	
	title := titles[bg.rng.Intn(len(titles))]
	
	// Use default if biome not found
	nameList, ok := names[biome]
	if !ok || len(nameList) == 0 {
		nameList = []string{"the Unknown", "Mystery", "Secrets", "the Beyond"}
	}
	
	name := nameList[bg.rng.Intn(len(nameList))]
	
	return title + " of " + name
}

// generateAttackPattern creates an attack pattern description
func (bg *BossGenerator) generateAttackPattern() string {
	patterns := []string{
		"triple_strike",
		"charge_attack",
		"area_blast",
		"summon_minions",
		"projectile_barrage",
		"ground_pound",
	}
	return patterns[bg.rng.Intn(len(patterns))]
}

// generateUniqueAttack creates a special attack
func (bg *BossGenerator) generateUniqueAttack(biome string) string {
	attacks := map[string][]string{
		"cave":    {"rock_fall", "earth_spike", "cave_in"},
		"forest":  {"vine_snare", "poison_cloud", "thorn_burst"},
		"ruins":   {"curse_wave", "soul_drain", "spectral_summon"},
		"crystal": {"ice_prison", "crystal_lance", "shard_storm"},
		"abyss":   {"void_rift", "corruption_beam", "shadow_clone"},
		"sky":     {"lightning_strike", "tornado", "wind_blade"},
	}
	
	// Use default if biome not found
	options, ok := attacks[biome]
	if !ok || len(options) == 0 {
		options = []string{"energy_blast", "power_strike", "devastating_blow"}
	}
	
	return options[bg.rng.Intn(len(options))]
}

// AbilityGenerator generates ability progression
type AbilityGenerator struct {
	rng *rand.Rand
}

// NewAbilityGenerator creates a new ability generator
func NewAbilityGenerator(seed int64) *AbilityGenerator {
	return &AbilityGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// GenerateProgression creates ability unlock order
func (ag *AbilityGenerator) GenerateProgression(seed int64) []Ability {
	ag.rng = rand.New(rand.NewSource(seed))
	
	// Define available abilities
	abilities := []Ability{
		{Name: "Double Jump", Type: MovementAbility, Description: "Jump again in mid-air"},
		{Name: "Dash", Type: MovementAbility, Description: "Quick burst of speed"},
		{Name: "Wall Climb", Type: MovementAbility, Description: "Climb vertical surfaces"},
		{Name: "Glide", Type: MovementAbility, Description: "Slow your fall"},
		{Name: "Swim", Type: MovementAbility, Description: "Move through liquids"},
		{Name: "Charge Attack", Type: CombatAbility, Description: "Powerful charged strike"},
		{Name: "Projectile", Type: CombatAbility, Description: "Ranged attack"},
		{Name: "Shield", Type: UtilityAbility, Description: "Temporary invulnerability"},
	}
	
	// Shuffle for random unlock order
	ag.rng.Shuffle(len(abilities), func(i, j int) {
		abilities[i], abilities[j] = abilities[j], abilities[i]
	})
	
	// Assign unlock order
	for i := range abilities {
		abilities[i].UnlockOrder = i
	}
	
	return abilities
}

// ItemGenerator generates items
type ItemGenerator struct {
	rng *rand.Rand
}

// NewItemGenerator creates a new item generator
func NewItemGenerator(seed int64) *ItemGenerator {
	return &ItemGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Generate creates an item
func (ig *ItemGenerator) Generate(itemType ItemType, seed int64) *Item {
	ig.rng = rand.New(rand.NewSource(seed))
	
	item := &Item{
		Type: itemType,
	}
	
	switch itemType {
	case WeaponItem:
		item.Name = ig.generateWeaponName()
		item.Description = "A powerful weapon"
		item.Effect = "increase_damage"
		item.Value = 50 + ig.rng.Intn(100)
		
	case ConsumableItem:
		item.Name = ig.generatePotionName()
		item.Description = "Restores health"
		item.Effect = "heal"
		item.Value = 10 + ig.rng.Intn(20)
		
	case KeyItem:
		item.Name = ig.generateKeyName()
		item.Description = "Opens new paths"
		item.Effect = "unlock"
		item.Value = 0
		
	case UpgradeItem:
		item.Name = "Upgrade Stone"
		item.Description = "Enhances abilities"
		item.Effect = "upgrade"
		item.Value = 100
		
	case CurrencyItem:
		item.Name = "Crystal Shard"
		item.Description = "Currency"
		item.Effect = "currency"
		item.Value = 1 + ig.rng.Intn(10)
	}
	
	return item
}

func (ig *ItemGenerator) generateWeaponName() string {
	adjectives := []string{"Blazing", "Frozen", "Shadow", "Holy", "Ancient"}
	nouns := []string{"Sword", "Blade", "Axe", "Spear", "Dagger"}
	
	adj := adjectives[ig.rng.Intn(len(adjectives))]
	noun := nouns[ig.rng.Intn(len(nouns))]
	
	return adj + " " + noun
}

func (ig *ItemGenerator) generatePotionName() string {
	colors := []string{"Red", "Blue", "Green", "Purple", "Golden"}
	color := colors[ig.rng.Intn(len(colors))]
	return color + " Potion"
}

func (ig *ItemGenerator) generateKeyName() string {
	materials := []string{"Iron", "Silver", "Gold", "Crystal", "Ancient"}
	material := materials[ig.rng.Intn(len(materials))]
	return material + " Key"
}

// ItemInstance represents a placed item in the game world
type ItemInstance struct {
	Item      *Item
	ID        int     // Unique identifier for this item instance
	X, Y      float64 // Position in the room
	Collected bool    // Whether the item has been collected
}

// NewItemInstance creates a new item instance
func NewItemInstance(item *Item, id int, x, y float64) *ItemInstance {
	return &ItemInstance{
		Item:      item,
		ID:        id,
		X:         x,
		Y:         y,
		Collected: false,
	}
}

// GetBounds returns the bounding box for collision detection
func (ii *ItemInstance) GetBounds() (x, y, width, height float64) {
	// Items are 16x16 pixels
	return ii.X, ii.Y, 16, 16
}
