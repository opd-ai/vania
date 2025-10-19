// Package achievement provides a comprehensive achievement tracking system
// that monitors player progress, combat performance, exploration, and special
// accomplishments throughout the game.
package achievement

import (
	"fmt"
	"time"
)

// AchievementID uniquely identifies an achievement
type AchievementID string

// Achievement categories
const (
	CategoryCombat      = "Combat"
	CategoryExploration = "Exploration"
	CategoryCollection  = "Collection"
	CategorySpeed       = "Speed"
	CategoryChallenge   = "Challenge"
	CategorySecret      = "Secret"
)

// Achievement rarity levels
const (
	RarityCommon    = "Common"
	RarityUncommon  = "Uncommon"
	RarityRare      = "Rare"
	RarityEpic      = "Epic"
	RarityLegendary = "Legendary"
)

// Achievement represents a single achievement definition
type Achievement struct {
	ID          AchievementID
	Name        string
	Description string
	Category    string
	Rarity      string
	Hidden      bool   // Hidden until unlocked
	Points      int    // Achievement points value
	IconIndex   int    // Index for procedurally generated icon
	
	// Requirements (criteria for unlocking)
	RequiresBosses     int     // Number of bosses to defeat
	RequiresRooms      int     // Number of rooms to visit
	RequiresItems      int     // Number of items to collect
	RequiresAbilities  int     // Number of abilities to unlock
	RequiresKills      int     // Number of enemies to kill
	RequiresDamage     int     // Total damage to deal
	RequiresNoDamage   bool    // Complete without taking damage
	RequiresTimeLimit  int64   // Complete within time limit (seconds)
	RequiresSpecial    string  // Special condition (custom logic)
}

// UnlockedAchievement represents an achievement that has been unlocked
type UnlockedAchievement struct {
	AchievementID AchievementID
	UnlockedAt    time.Time
	Progress      float64 // 0.0 to 1.0
}

// AchievementProgress tracks progress toward an achievement
type AchievementProgress struct {
	AchievementID   AchievementID
	CurrentValue    int     // Current progress value
	TargetValue     int     // Target value needed
	Progress        float64 // 0.0 to 1.0
	LastUpdated     time.Time
}

// AchievementTracker manages achievement tracking and unlocking
type AchievementTracker struct {
	achievements map[AchievementID]*Achievement
	unlocked     map[AchievementID]*UnlockedAchievement
	progress     map[AchievementID]*AchievementProgress
	
	// Runtime statistics
	stats Statistics
	
	// Event handlers
	onUnlock func(achievement *Achievement)
}

// Statistics tracks player statistics for achievement calculations
type Statistics struct {
	// Combat
	EnemiesDefeated    int
	BossesDefeated     int
	TotalDamageDealt   int
	DamageTaken        int
	PerfectKills       int  // Enemies killed without taking damage
	
	// Exploration
	RoomsVisited       int
	BiomesExplored     int
	SecretsFound       int
	
	// Collection
	ItemsCollected     int
	AbilitiesUnlocked  int
	
	// Speed/Time
	StartTime          time.Time
	PlayTime           int64  // seconds
	FastestBossKill    int64  // seconds
	
	// Special
	DeathCount         int
	PerfectRooms       int  // Rooms cleared without damage
	ConsecutiveKills   int
	LongestCombo       int
}

// NewAchievementTracker creates a new achievement tracker
func NewAchievementTracker() *AchievementTracker {
	tracker := &AchievementTracker{
		achievements: make(map[AchievementID]*Achievement),
		unlocked:     make(map[AchievementID]*UnlockedAchievement),
		progress:     make(map[AchievementID]*AchievementProgress),
		stats:        Statistics{
			StartTime: time.Now(),
		},
	}
	
	// Register all achievements
	tracker.registerDefaultAchievements()
	
	return tracker
}

// registerDefaultAchievements registers all built-in achievements
func (at *AchievementTracker) registerDefaultAchievements() {
	achievements := []*Achievement{
		// Combat Achievements
		{
			ID:          "first_blood",
			Name:        "First Blood",
			Description: "Defeat your first enemy",
			Category:    CategoryCombat,
			Rarity:      RarityCommon,
			Points:      10,
			IconIndex:   0,
			RequiresKills: 1,
		},
		{
			ID:          "slayer",
			Name:        "Slayer",
			Description: "Defeat 50 enemies",
			Category:    CategoryCombat,
			Rarity:      RarityUncommon,
			Points:      25,
			IconIndex:   1,
			RequiresKills: 50,
		},
		{
			ID:          "destroyer",
			Name:        "Destroyer",
			Description: "Defeat 100 enemies",
			Category:    CategoryCombat,
			Rarity:      RarityRare,
			Points:      50,
			IconIndex:   2,
			RequiresKills: 100,
		},
		{
			ID:          "boss_hunter",
			Name:        "Boss Hunter",
			Description: "Defeat your first boss",
			Category:    CategoryCombat,
			Rarity:      RarityCommon,
			Points:      20,
			IconIndex:   3,
			RequiresBosses: 1,
		},
		{
			ID:          "boss_slayer",
			Name:        "Boss Slayer",
			Description: "Defeat all bosses",
			Category:    CategoryCombat,
			Rarity:      RarityEpic,
			Points:      100,
			IconIndex:   4,
			RequiresBosses: 10, // Will be adjusted based on actual boss count
		},
		{
			ID:          "perfectionist",
			Name:        "Perfectionist",
			Description: "Defeat a boss without taking damage",
			Category:    CategoryChallenge,
			Rarity:      RarityRare,
			Points:      75,
			IconIndex:   5,
			RequiresNoDamage: true,
			RequiresBosses: 1,
		},
		
		// Exploration Achievements
		{
			ID:          "explorer",
			Name:        "Explorer",
			Description: "Visit 10 different rooms",
			Category:    CategoryExploration,
			Rarity:      RarityCommon,
			Points:      10,
			IconIndex:   6,
			RequiresRooms: 10,
		},
		{
			ID:          "cartographer",
			Name:        "Cartographer",
			Description: "Visit 50 different rooms",
			Category:    CategoryExploration,
			Rarity:      RarityUncommon,
			Points:      30,
			IconIndex:   7,
			RequiresRooms: 50,
		},
		{
			ID:          "master_explorer",
			Name:        "Master Explorer",
			Description: "Visit all rooms in the world",
			Category:    CategoryExploration,
			Rarity:      RarityEpic,
			Points:      100,
			IconIndex:   8,
			RequiresRooms: 100, // Will be adjusted based on world size
		},
		
		// Collection Achievements
		{
			ID:          "treasure_hunter",
			Name:        "Treasure Hunter",
			Description: "Collect 10 items",
			Category:    CategoryCollection,
			Rarity:      RarityCommon,
			Points:      15,
			IconIndex:   9,
			RequiresItems: 10,
		},
		{
			ID:          "hoarder",
			Name:        "Hoarder",
			Description: "Collect 25 items",
			Category:    CategoryCollection,
			Rarity:      RarityUncommon,
			Points:      35,
			IconIndex:   10,
			RequiresItems: 25,
		},
		{
			ID:          "ability_master",
			Name:        "Ability Master",
			Description: "Unlock all abilities",
			Category:    CategoryCollection,
			Rarity:      RarityRare,
			Points:      80,
			IconIndex:   11,
			RequiresAbilities: 8, // Will be adjusted based on ability count
		},
		
		// Speed Achievements
		{
			ID:          "speedrunner",
			Name:        "Speedrunner",
			Description: "Complete the game in under 30 minutes",
			Category:    CategorySpeed,
			Rarity:      RarityEpic,
			Points:      150,
			IconIndex:   12,
			RequiresTimeLimit: 1800, // 30 minutes
			RequiresBosses: 10,
		},
		{
			ID:          "flash",
			Name:        "Flash",
			Description: "Defeat a boss in under 60 seconds",
			Category:    CategorySpeed,
			Rarity:      RarityRare,
			Points:      60,
			IconIndex:   13,
			RequiresTimeLimit: 60,
			RequiresBosses: 1,
		},
		
		// Challenge Achievements
		{
			ID:          "untouchable",
			Name:        "Untouchable",
			Description: "Clear 10 rooms without taking damage",
			Category:    CategoryChallenge,
			Rarity:      RarityRare,
			Points:      90,
			IconIndex:   14,
			RequiresSpecial: "perfect_rooms_10",
		},
		{
			ID:          "combo_master",
			Name:        "Combo Master",
			Description: "Achieve a 20-hit combo",
			Category:    CategoryChallenge,
			Rarity:      RarityUncommon,
			Points:      40,
			IconIndex:   15,
			RequiresSpecial: "combo_20",
		},
		{
			ID:          "survivor",
			Name:        "Survivor",
			Description: "Complete the game without dying",
			Category:    CategoryChallenge,
			Rarity:      RarityLegendary,
			Points:      200,
			IconIndex:   16,
			RequiresSpecial: "no_deaths",
			RequiresBosses: 10,
		},
		
		// Secret Achievements
		{
			ID:          "secret_finder",
			Name:        "Secret Finder",
			Description: "Discover a hidden secret",
			Category:    CategorySecret,
			Rarity:      RarityUncommon,
			Points:      30,
			IconIndex:   17,
			Hidden:      true,
			RequiresSpecial: "secret_1",
		},
		{
			ID:          "completionist",
			Name:        "Completionist",
			Description: "Unlock all achievements",
			Category:    CategorySecret,
			Rarity:      RarityLegendary,
			Points:      250,
			IconIndex:   18,
			Hidden:      true,
			RequiresSpecial: "all_achievements",
		},
	}
	
	for _, achievement := range achievements {
		at.RegisterAchievement(achievement)
	}
}

// RegisterAchievement registers a new achievement
func (at *AchievementTracker) RegisterAchievement(achievement *Achievement) {
	at.achievements[achievement.ID] = achievement
	
	// Initialize progress tracking
	at.progress[achievement.ID] = &AchievementProgress{
		AchievementID: achievement.ID,
		CurrentValue:  0,
		TargetValue:   at.getTargetValue(achievement),
		Progress:      0.0,
		LastUpdated:   time.Now(),
	}
}

// getTargetValue calculates the target value for an achievement
func (at *AchievementTracker) getTargetValue(achievement *Achievement) int {
	if achievement.RequiresKills > 0 {
		return achievement.RequiresKills
	}
	if achievement.RequiresBosses > 0 {
		return achievement.RequiresBosses
	}
	if achievement.RequiresRooms > 0 {
		return achievement.RequiresRooms
	}
	if achievement.RequiresItems > 0 {
		return achievement.RequiresItems
	}
	if achievement.RequiresAbilities > 0 {
		return achievement.RequiresAbilities
	}
	return 1 // Default for special achievements
}

// UpdateStatistics updates statistics and checks for achievement unlocks
func (at *AchievementTracker) UpdateStatistics(stats Statistics) {
	at.stats = stats
	at.checkAchievements()
}

// RecordEnemyKill records an enemy kill
func (at *AchievementTracker) RecordEnemyKill(wasPerfect bool) {
	at.stats.EnemiesDefeated++
	if wasPerfect {
		at.stats.PerfectKills++
	}
	at.checkAchievements()
}

// RecordBossKill records a boss defeat
func (at *AchievementTracker) RecordBossKill(timeTaken int64, wasPerfect bool) {
	at.stats.BossesDefeated++
	if wasPerfect {
		at.stats.PerfectKills++
	}
	
	// Track fastest boss kill
	if at.stats.FastestBossKill == 0 || timeTaken < at.stats.FastestBossKill {
		at.stats.FastestBossKill = timeTaken
	}
	
	at.checkAchievements()
}

// RecordRoomVisit records visiting a room
func (at *AchievementTracker) RecordRoomVisit(isPerfect bool) {
	at.stats.RoomsVisited++
	if isPerfect {
		at.stats.PerfectRooms++
	}
	at.checkAchievements()
}

// RecordItemCollected records collecting an item
func (at *AchievementTracker) RecordItemCollected() {
	at.stats.ItemsCollected++
	at.checkAchievements()
}

// RecordAbilityUnlocked records unlocking an ability
func (at *AchievementTracker) RecordAbilityUnlocked() {
	at.stats.AbilitiesUnlocked++
	at.checkAchievements()
}

// RecordDamage records damage dealt or taken
func (at *AchievementTracker) RecordDamage(dealt int, taken int) {
	at.stats.TotalDamageDealt += dealt
	at.stats.DamageTaken += taken
	at.checkAchievements()
}

// RecordCombo records a combo achievement
func (at *AchievementTracker) RecordCombo(comboCount int) {
	at.stats.ConsecutiveKills = comboCount
	if comboCount > at.stats.LongestCombo {
		at.stats.LongestCombo = comboCount
	}
	at.checkAchievements()
}

// RecordDeath records a player death
func (at *AchievementTracker) RecordDeath() {
	at.stats.DeathCount++
}

// checkAchievements checks all achievements and unlocks completed ones
func (at *AchievementTracker) checkAchievements() {
	for id, achievement := range at.achievements {
		// Skip already unlocked achievements
		if at.IsUnlocked(id) {
			continue
		}
		
		// Check if requirements are met
		if at.checkRequirements(achievement) {
			at.UnlockAchievement(id)
		} else {
			// Update progress
			at.updateProgress(achievement)
		}
	}
}

// checkRequirements checks if an achievement's requirements are met
func (at *AchievementTracker) checkRequirements(achievement *Achievement) bool {
	// Check kill requirements
	if achievement.RequiresKills > 0 && at.stats.EnemiesDefeated < achievement.RequiresKills {
		return false
	}
	
	// Check boss requirements
	if achievement.RequiresBosses > 0 && at.stats.BossesDefeated < achievement.RequiresBosses {
		return false
	}
	
	// Check room requirements
	if achievement.RequiresRooms > 0 && at.stats.RoomsVisited < achievement.RequiresRooms {
		return false
	}
	
	// Check item requirements
	if achievement.RequiresItems > 0 && at.stats.ItemsCollected < achievement.RequiresItems {
		return false
	}
	
	// Check ability requirements
	if achievement.RequiresAbilities > 0 && at.stats.AbilitiesUnlocked < achievement.RequiresAbilities {
		return false
	}
	
	// Check time limit
	if achievement.RequiresTimeLimit > 0 {
		playTime := at.stats.PlayTime
		if playTime == 0 {
			playTime = int64(time.Since(at.stats.StartTime).Seconds())
		}
		if playTime > achievement.RequiresTimeLimit {
			return false
		}
	}
	
	// Check special requirements
	if achievement.RequiresSpecial != "" {
		return at.checkSpecialRequirement(achievement)
	}
	
	return true
}

// checkSpecialRequirement checks custom achievement requirements
func (at *AchievementTracker) checkSpecialRequirement(achievement *Achievement) bool {
	switch achievement.RequiresSpecial {
	case "perfect_rooms_10":
		return at.stats.PerfectRooms >= 10
	case "combo_20":
		return at.stats.LongestCombo >= 20
	case "no_deaths":
		return at.stats.DeathCount == 0 && at.stats.BossesDefeated >= 10
	case "secret_1":
		return at.stats.SecretsFound >= 1
	case "all_achievements":
		// Check if all other achievements are unlocked
		totalAchievements := len(at.achievements) - 1 // Exclude this one
		return len(at.unlocked) >= totalAchievements
	default:
		return false
	}
}

// updateProgress updates progress tracking for an achievement
func (at *AchievementTracker) updateProgress(achievement *Achievement) {
	progress := at.progress[achievement.ID]
	if progress == nil {
		return
	}
	
	// Determine current value based on achievement type
	var currentValue int
	if achievement.RequiresKills > 0 {
		currentValue = at.stats.EnemiesDefeated
	} else if achievement.RequiresBosses > 0 {
		currentValue = at.stats.BossesDefeated
	} else if achievement.RequiresRooms > 0 {
		currentValue = at.stats.RoomsVisited
	} else if achievement.RequiresItems > 0 {
		currentValue = at.stats.ItemsCollected
	} else if achievement.RequiresAbilities > 0 {
		currentValue = at.stats.AbilitiesUnlocked
	}
	
	progress.CurrentValue = currentValue
	progress.Progress = float64(currentValue) / float64(progress.TargetValue)
	if progress.Progress > 1.0 {
		progress.Progress = 1.0
	}
	progress.LastUpdated = time.Now()
}

// UnlockAchievement unlocks an achievement
func (at *AchievementTracker) UnlockAchievement(id AchievementID) error {
	achievement, exists := at.achievements[id]
	if !exists {
		return fmt.Errorf("achievement not found: %s", id)
	}
	
	// Check if already unlocked
	if at.IsUnlocked(id) {
		return nil
	}
	
	// Create unlocked achievement record
	unlocked := &UnlockedAchievement{
		AchievementID: id,
		UnlockedAt:    time.Now(),
		Progress:      1.0,
	}
	
	at.unlocked[id] = unlocked
	
	// Update progress to complete
	if progress := at.progress[id]; progress != nil {
		progress.Progress = 1.0
		progress.CurrentValue = progress.TargetValue
	}
	
	// Trigger unlock callback
	if at.onUnlock != nil {
		at.onUnlock(achievement)
	}
	
	return nil
}

// IsUnlocked checks if an achievement is unlocked
func (at *AchievementTracker) IsUnlocked(id AchievementID) bool {
	_, unlocked := at.unlocked[id]
	return unlocked
}

// GetProgress returns progress for an achievement
func (at *AchievementTracker) GetProgress(id AchievementID) *AchievementProgress {
	return at.progress[id]
}

// GetAchievement returns an achievement by ID
func (at *AchievementTracker) GetAchievement(id AchievementID) *Achievement {
	return at.achievements[id]
}

// GetAllAchievements returns all registered achievements
func (at *AchievementTracker) GetAllAchievements() []*Achievement {
	achievements := make([]*Achievement, 0, len(at.achievements))
	for _, achievement := range at.achievements {
		achievements = append(achievements, achievement)
	}
	return achievements
}

// GetUnlockedAchievements returns all unlocked achievements
func (at *AchievementTracker) GetUnlockedAchievements() []*UnlockedAchievement {
	unlocked := make([]*UnlockedAchievement, 0, len(at.unlocked))
	for _, u := range at.unlocked {
		unlocked = append(unlocked, u)
	}
	return unlocked
}

// GetStatistics returns current statistics
func (at *AchievementTracker) GetStatistics() Statistics {
	return at.stats
}

// SetOnUnlock sets the callback for when an achievement is unlocked
func (at *AchievementTracker) SetOnUnlock(callback func(achievement *Achievement)) {
	at.onUnlock = callback
}

// GetCompletionPercentage returns overall achievement completion percentage
func (at *AchievementTracker) GetCompletionPercentage() float64 {
	if len(at.achievements) == 0 {
		return 0.0
	}
	return float64(len(at.unlocked)) / float64(len(at.achievements)) * 100.0
}

// GetTotalPoints returns total achievement points earned
func (at *AchievementTracker) GetTotalPoints() int {
	total := 0
	for _, unlocked := range at.unlocked {
		if achievement := at.achievements[unlocked.AchievementID]; achievement != nil {
			total += achievement.Points
		}
	}
	return total
}

// GetMaxPoints returns maximum possible achievement points
func (at *AchievementTracker) GetMaxPoints() int {
	total := 0
	for _, achievement := range at.achievements {
		total += achievement.Points
	}
	return total
}
