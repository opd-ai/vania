package achievement

import (
	"testing"
	"time"
)

// TestNewAchievementTracker tests creating a new achievement tracker
func TestNewAchievementTracker(t *testing.T) {
	tracker := NewAchievementTracker()
	
	if tracker == nil {
		t.Fatal("Expected tracker to be created")
	}
	
	// Check that default achievements are registered
	achievements := tracker.GetAllAchievements()
	if len(achievements) == 0 {
		t.Error("Expected default achievements to be registered")
	}
	
	// Verify specific achievements exist
	firstBlood := tracker.GetAchievement("first_blood")
	if firstBlood == nil {
		t.Error("Expected 'first_blood' achievement to exist")
	}
	
	if firstBlood != nil && firstBlood.Name != "First Blood" {
		t.Errorf("Expected achievement name 'First Blood', got '%s'", firstBlood.Name)
	}
}

// TestRegisterAchievement tests registering a custom achievement
func TestRegisterAchievement(t *testing.T) {
	tracker := NewAchievementTracker()
	
	customAchievement := &Achievement{
		ID:          "test_achievement",
		Name:        "Test Achievement",
		Description: "A test achievement",
		Category:    CategoryCombat,
		Rarity:      RarityCommon,
		Points:      10,
		RequiresKills: 5,
	}
	
	tracker.RegisterAchievement(customAchievement)
	
	// Verify achievement was registered
	retrieved := tracker.GetAchievement("test_achievement")
	if retrieved == nil {
		t.Fatal("Expected custom achievement to be registered")
	}
	
	if retrieved.Name != "Test Achievement" {
		t.Errorf("Expected name 'Test Achievement', got '%s'", retrieved.Name)
	}
	
	// Verify progress tracking was initialized
	progress := tracker.GetProgress("test_achievement")
	if progress == nil {
		t.Error("Expected progress tracking to be initialized")
	}
	
	if progress != nil && progress.TargetValue != 5 {
		t.Errorf("Expected target value 5, got %d", progress.TargetValue)
	}
}

// TestRecordEnemyKill tests recording enemy kills
func TestRecordEnemyKill(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Record first kill
	tracker.RecordEnemyKill(false)
	
	stats := tracker.GetStatistics()
	if stats.EnemiesDefeated != 1 {
		t.Errorf("Expected 1 enemy defeated, got %d", stats.EnemiesDefeated)
	}
	
	// Check if "first_blood" achievement was unlocked
	if !tracker.IsUnlocked("first_blood") {
		t.Error("Expected 'first_blood' achievement to be unlocked after first kill")
	}
	
	// Record more kills
	for i := 0; i < 49; i++ {
		tracker.RecordEnemyKill(false)
	}
	
	// Check if "slayer" achievement was unlocked
	if !tracker.IsUnlocked("slayer") {
		t.Error("Expected 'slayer' achievement to be unlocked after 50 kills")
	}
}

// TestRecordBossKill tests recording boss defeats
func TestRecordBossKill(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Record first boss kill
	tracker.RecordBossKill(120, false)
	
	stats := tracker.GetStatistics()
	if stats.BossesDefeated != 1 {
		t.Errorf("Expected 1 boss defeated, got %d", stats.BossesDefeated)
	}
	
	// Check if "boss_hunter" achievement was unlocked
	if !tracker.IsUnlocked("boss_hunter") {
		t.Error("Expected 'boss_hunter' achievement to be unlocked after first boss")
	}
	
	// Check fastest boss kill tracking
	if stats.FastestBossKill != 120 {
		t.Errorf("Expected fastest boss kill to be 120, got %d", stats.FastestBossKill)
	}
	
	// Record faster boss kill
	tracker.RecordBossKill(60, false)
	stats = tracker.GetStatistics()
	if stats.FastestBossKill != 60 {
		t.Errorf("Expected fastest boss kill to be updated to 60, got %d", stats.FastestBossKill)
	}
}

// TestRecordRoomVisit tests recording room visits
func TestRecordRoomVisit(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Visit 10 rooms
	for i := 0; i < 10; i++ {
		tracker.RecordRoomVisit(false)
	}
	
	stats := tracker.GetStatistics()
	if stats.RoomsVisited != 10 {
		t.Errorf("Expected 10 rooms visited, got %d", stats.RoomsVisited)
	}
	
	// Check if "explorer" achievement was unlocked
	if !tracker.IsUnlocked("explorer") {
		t.Error("Expected 'explorer' achievement to be unlocked after 10 rooms")
	}
}

// TestRecordItemCollected tests recording item collection
func TestRecordItemCollected(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Collect 10 items
	for i := 0; i < 10; i++ {
		tracker.RecordItemCollected()
	}
	
	stats := tracker.GetStatistics()
	if stats.ItemsCollected != 10 {
		t.Errorf("Expected 10 items collected, got %d", stats.ItemsCollected)
	}
	
	// Check if "treasure_hunter" achievement was unlocked
	if !tracker.IsUnlocked("treasure_hunter") {
		t.Error("Expected 'treasure_hunter' achievement to be unlocked after 10 items")
	}
}

// TestRecordAbilityUnlocked tests recording ability unlocks
func TestRecordAbilityUnlocked(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Unlock 8 abilities
	for i := 0; i < 8; i++ {
		tracker.RecordAbilityUnlocked()
	}
	
	stats := tracker.GetStatistics()
	if stats.AbilitiesUnlocked != 8 {
		t.Errorf("Expected 8 abilities unlocked, got %d", stats.AbilitiesUnlocked)
	}
	
	// Check if "ability_master" achievement was unlocked
	if !tracker.IsUnlocked("ability_master") {
		t.Error("Expected 'ability_master' achievement to be unlocked after 8 abilities")
	}
}

// TestRecordCombo tests recording combo achievements
func TestRecordCombo(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Record a 20-hit combo
	tracker.RecordCombo(20)
	
	stats := tracker.GetStatistics()
	if stats.LongestCombo != 20 {
		t.Errorf("Expected longest combo to be 20, got %d", stats.LongestCombo)
	}
	
	// Check if "combo_master" achievement was unlocked
	if !tracker.IsUnlocked("combo_master") {
		t.Error("Expected 'combo_master' achievement to be unlocked after 20-hit combo")
	}
	
	// Record a higher combo
	tracker.RecordCombo(25)
	stats = tracker.GetStatistics()
	if stats.LongestCombo != 25 {
		t.Errorf("Expected longest combo to be updated to 25, got %d", stats.LongestCombo)
	}
}

// TestPerfectRoomTracking tests perfect room tracking
func TestPerfectRoomTracking(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Visit 10 perfect rooms
	for i := 0; i < 10; i++ {
		tracker.RecordRoomVisit(true)
	}
	
	stats := tracker.GetStatistics()
	if stats.PerfectRooms != 10 {
		t.Errorf("Expected 10 perfect rooms, got %d", stats.PerfectRooms)
	}
	
	// Check if "untouchable" achievement was unlocked
	if !tracker.IsUnlocked("untouchable") {
		t.Error("Expected 'untouchable' achievement to be unlocked after 10 perfect rooms")
	}
}

// TestProgressTracking tests achievement progress tracking
func TestProgressTracking(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Record 25 enemy kills (halfway to "slayer")
	for i := 0; i < 25; i++ {
		tracker.RecordEnemyKill(false)
	}
	
	// Check progress for "slayer" achievement
	progress := tracker.GetProgress("slayer")
	if progress == nil {
		t.Fatal("Expected progress tracking for 'slayer' achievement")
	}
	
	if progress.CurrentValue != 25 {
		t.Errorf("Expected current value 25, got %d", progress.CurrentValue)
	}
	
	if progress.TargetValue != 50 {
		t.Errorf("Expected target value 50, got %d", progress.TargetValue)
	}
	
	expectedProgress := 0.5 // 25/50
	if progress.Progress != expectedProgress {
		t.Errorf("Expected progress %.2f, got %.2f", expectedProgress, progress.Progress)
	}
}

// TestUnlockCallback tests achievement unlock callback
func TestUnlockCallback(t *testing.T) {
	tracker := NewAchievementTracker()
	
	var unlockedAchievement *Achievement
	tracker.SetOnUnlock(func(achievement *Achievement) {
		unlockedAchievement = achievement
	})
	
	// Trigger an achievement
	tracker.RecordEnemyKill(false)
	
	// Check if callback was called
	if unlockedAchievement == nil {
		t.Error("Expected unlock callback to be called")
	}
	
	if unlockedAchievement != nil && unlockedAchievement.ID != "first_blood" {
		t.Errorf("Expected 'first_blood' achievement, got '%s'", unlockedAchievement.ID)
	}
}

// TestCompletionPercentage tests completion percentage calculation
func TestCompletionPercentage(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Initially 0%
	completion := tracker.GetCompletionPercentage()
	if completion != 0.0 {
		t.Errorf("Expected 0%% completion, got %.2f%%", completion)
	}
	
	// Unlock one achievement
	tracker.RecordEnemyKill(false)
	
	completion = tracker.GetCompletionPercentage()
	if completion <= 0.0 || completion >= 100.0 {
		t.Errorf("Expected completion between 0%% and 100%%, got %.2f%%", completion)
	}
}

// TestTotalPoints tests total points calculation
func TestTotalPoints(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Initially 0 points
	points := tracker.GetTotalPoints()
	if points != 0 {
		t.Errorf("Expected 0 points, got %d", points)
	}
	
	// Unlock "first_blood" (10 points)
	tracker.RecordEnemyKill(false)
	
	points = tracker.GetTotalPoints()
	if points < 10 {
		t.Errorf("Expected at least 10 points, got %d", points)
	}
	
	initialPoints := points
	
	// Unlock "boss_hunter" (20 points) and possibly "flash" (60 points)
	tracker.RecordBossKill(60, false)
	
	points = tracker.GetTotalPoints()
	if points <= initialPoints {
		t.Errorf("Expected points to increase from %d, got %d", initialPoints, points)
	}
}

// TestMaxPoints tests maximum points calculation
func TestMaxPoints(t *testing.T) {
	tracker := NewAchievementTracker()
	
	maxPoints := tracker.GetMaxPoints()
	if maxPoints <= 0 {
		t.Error("Expected maximum points to be greater than 0")
	}
	
	// Verify it's the sum of all achievement points
	expectedMax := 0
	for _, achievement := range tracker.GetAllAchievements() {
		expectedMax += achievement.Points
	}
	
	if maxPoints != expectedMax {
		t.Errorf("Expected max points %d, got %d", expectedMax, maxPoints)
	}
}

// TestHiddenAchievements tests hidden achievement behavior
func TestHiddenAchievements(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Check that hidden achievements exist
	secretFinder := tracker.GetAchievement("secret_finder")
	if secretFinder == nil {
		t.Fatal("Expected 'secret_finder' achievement to exist")
	}
	
	if !secretFinder.Hidden {
		t.Error("Expected 'secret_finder' to be a hidden achievement")
	}
}

// TestGetUnlockedAchievements tests retrieving unlocked achievements
func TestGetUnlockedAchievements(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Initially no unlocked achievements
	unlocked := tracker.GetUnlockedAchievements()
	if len(unlocked) != 0 {
		t.Errorf("Expected 0 unlocked achievements, got %d", len(unlocked))
	}
	
	// Unlock some achievements (may unlock more than expected due to multiple criteria)
	tracker.RecordEnemyKill(false)
	tracker.RecordBossKill(60, false)
	
	unlocked = tracker.GetUnlockedAchievements()
	if len(unlocked) < 2 {
		t.Errorf("Expected at least 2 unlocked achievements, got %d", len(unlocked))
	}
	
	// Verify unlocked achievements have timestamps
	for _, u := range unlocked {
		if u.UnlockedAt.IsZero() {
			t.Error("Expected unlocked achievement to have a timestamp")
		}
	}
}

// TestSpecialRequirementNoDamage tests no-damage special requirement
func TestSpecialRequirementNoDamage(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Defeat 10 bosses without taking damage
	for i := 0; i < 10; i++ {
		tracker.RecordBossKill(60, true)
	}
	
	// Ensure no damage was taken
	stats := tracker.GetStatistics()
	stats.DamageTaken = 0
	tracker.UpdateStatistics(stats)
	
	// Check if "survivor" achievement was unlocked
	if !tracker.IsUnlocked("survivor") {
		t.Error("Expected 'survivor' achievement to be unlocked with no deaths and all bosses defeated")
	}
}

// TestDamageTracking tests damage tracking
func TestDamageTracking(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Record damage
	tracker.RecordDamage(100, 50)
	
	stats := tracker.GetStatistics()
	if stats.TotalDamageDealt != 100 {
		t.Errorf("Expected 100 damage dealt, got %d", stats.TotalDamageDealt)
	}
	
	if stats.DamageTaken != 50 {
		t.Errorf("Expected 50 damage taken, got %d", stats.DamageTaken)
	}
}

// TestDeathTracking tests death tracking
func TestDeathTracking(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Record deaths
	tracker.RecordDeath()
	tracker.RecordDeath()
	
	stats := tracker.GetStatistics()
	if stats.DeathCount != 2 {
		t.Errorf("Expected 2 deaths, got %d", stats.DeathCount)
	}
}

// TestUnlockAchievementManually tests manually unlocking an achievement
func TestUnlockAchievementManually(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Manually unlock an achievement
	err := tracker.UnlockAchievement("first_blood")
	if err != nil {
		t.Errorf("Expected no error unlocking achievement, got: %v", err)
	}
	
	// Verify it's unlocked
	if !tracker.IsUnlocked("first_blood") {
		t.Error("Expected achievement to be unlocked")
	}
	
	// Try to unlock again (should not error)
	err = tracker.UnlockAchievement("first_blood")
	if err != nil {
		t.Errorf("Expected no error unlocking already-unlocked achievement, got: %v", err)
	}
	
	// Try to unlock non-existent achievement
	err = tracker.UnlockAchievement("fake_achievement")
	if err == nil {
		t.Error("Expected error unlocking non-existent achievement")
	}
}

// TestUpdateStatistics tests updating statistics in bulk
func TestUpdateStatistics(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Create custom statistics
	stats := Statistics{
		EnemiesDefeated:   100,
		BossesDefeated:    5,
		RoomsVisited:      50,
		ItemsCollected:    25,
		AbilitiesUnlocked: 8,
		StartTime:         time.Now(),
	}
	
	// Update all statistics at once
	tracker.UpdateStatistics(stats)
	
	// Verify statistics were updated
	retrievedStats := tracker.GetStatistics()
	if retrievedStats.EnemiesDefeated != 100 {
		t.Errorf("Expected 100 enemies defeated, got %d", retrievedStats.EnemiesDefeated)
	}
	
	// Verify multiple achievements were unlocked
	expectedUnlocked := []AchievementID{
		"first_blood",
		"slayer",
		"destroyer",
		"boss_hunter",
		"explorer",
		"cartographer",
		"treasure_hunter",
		"hoarder",
		"ability_master",
	}
	
	for _, id := range expectedUnlocked {
		if !tracker.IsUnlocked(id) {
			t.Errorf("Expected achievement '%s' to be unlocked", id)
		}
	}
}

// TestAchievementCategories tests achievement categorization
func TestAchievementCategories(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Count achievements by category
	categoryCounts := make(map[string]int)
	for _, achievement := range tracker.GetAllAchievements() {
		categoryCounts[achievement.Category]++
	}
	
	// Verify all categories have achievements
	expectedCategories := []string{
		CategoryCombat,
		CategoryExploration,
		CategoryCollection,
		CategorySpeed,
		CategoryChallenge,
		CategorySecret,
	}
	
	for _, category := range expectedCategories {
		if categoryCounts[category] == 0 {
			t.Errorf("Expected at least one achievement in category '%s'", category)
		}
	}
}

// TestAchievementRarities tests achievement rarity levels
func TestAchievementRarities(t *testing.T) {
	tracker := NewAchievementTracker()
	
	// Count achievements by rarity
	rarityCounts := make(map[string]int)
	for _, achievement := range tracker.GetAllAchievements() {
		rarityCounts[achievement.Rarity]++
	}
	
	// Verify different rarities exist
	if rarityCounts[RarityCommon] == 0 {
		t.Error("Expected at least one common achievement")
	}
	
	if rarityCounts[RarityLegendary] == 0 {
		t.Error("Expected at least one legendary achievement")
	}
}
