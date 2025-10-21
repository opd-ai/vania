package pcg

import (
	"sync"
	"testing"
)

// TestNewAssetCache_CreatesValidCache tests cache creation
func TestNewAssetCache_CreatesValidCache(t *testing.T) {
	cache := NewAssetCache()

	if cache == nil {
		t.Fatal("NewAssetCache returned nil")
	}

	if cache.Sprites == nil {
		t.Error("Sprites map is nil")
	}
	if cache.Sounds == nil {
		t.Error("Sounds map is nil")
	}
	if cache.Music == nil {
		t.Error("Music map is nil")
	}
	if cache.Narrative == nil {
		t.Error("Narrative map is nil")
	}
}

// TestNewAssetCache_MapsAreEmpty tests that new cache has empty maps
func TestNewAssetCache_MapsAreEmpty(t *testing.T) {
	cache := NewAssetCache()

	if len(cache.Sprites) != 0 {
		t.Errorf("Expected empty Sprites map, got %d items", len(cache.Sprites))
	}
	if len(cache.Sounds) != 0 {
		t.Errorf("Expected empty Sounds map, got %d items", len(cache.Sounds))
	}
	if len(cache.Music) != 0 {
		t.Errorf("Expected empty Music map, got %d items", len(cache.Music))
	}
	if len(cache.Narrative) != 0 {
		t.Errorf("Expected empty Narrative map, got %d items", len(cache.Narrative))
	}
}

// TestSetSprite_StoresValue tests sprite storage
func TestSetSprite_StoresValue(t *testing.T) {
	cache := NewAssetCache()
	key := "player_sprite"
	value := "sprite_data"

	cache.SetSprite(key, value)

	if len(cache.Sprites) != 1 {
		t.Errorf("Expected 1 sprite, got %d", len(cache.Sprites))
	}
}

// TestGetSprite_RetrievesStoredValue tests sprite retrieval
func TestGetSprite_RetrievesStoredValue(t *testing.T) {
	cache := NewAssetCache()
	key := "enemy_sprite"
	expectedValue := "sprite_data_123"

	cache.SetSprite(key, expectedValue)
	value, ok := cache.GetSprite(key)

	if !ok {
		t.Error("GetSprite returned false for existing key")
	}
	if value != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, value)
	}
}

// TestGetSprite_ReturnsFalseForMissingKey tests missing sprite retrieval
func TestGetSprite_ReturnsFalseForMissingKey(t *testing.T) {
	cache := NewAssetCache()

	value, ok := cache.GetSprite("nonexistent")

	if ok {
		t.Error("GetSprite returned true for nonexistent key")
	}
	if value != nil {
		t.Error("GetSprite returned non-nil value for nonexistent key")
	}
}

// TestSetSound_StoresValue tests sound storage
func TestSetSound_StoresValue(t *testing.T) {
	cache := NewAssetCache()
	key := "jump_sound"
	value := "sound_data"

	cache.SetSound(key, value)

	if len(cache.Sounds) != 1 {
		t.Errorf("Expected 1 sound, got %d", len(cache.Sounds))
	}
}

// TestGetSound_RetrievesStoredValue tests sound retrieval
func TestGetSound_RetrievesStoredValue(t *testing.T) {
	cache := NewAssetCache()
	key := "explosion_sound"
	expectedValue := "sound_data_456"

	cache.SetSound(key, expectedValue)
	value, ok := cache.GetSound(key)

	if !ok {
		t.Error("GetSound returned false for existing key")
	}
	if value != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, value)
	}
}

// TestGetSound_ReturnsFalseForMissingKey tests missing sound retrieval
func TestGetSound_ReturnsFalseForMissingKey(t *testing.T) {
	cache := NewAssetCache()

	value, ok := cache.GetSound("nonexistent")

	if ok {
		t.Error("GetSound returned true for nonexistent key")
	}
	if value != nil {
		t.Error("GetSound returned non-nil value for nonexistent key")
	}
}

// TestSetMusic_StoresValue tests music storage
func TestSetMusic_StoresValue(t *testing.T) {
	cache := NewAssetCache()
	key := "theme_music"
	value := "music_data"

	cache.SetMusic(key, value)

	if len(cache.Music) != 1 {
		t.Errorf("Expected 1 music, got %d", len(cache.Music))
	}
}

// TestGetMusic_RetrievesStoredValue tests music retrieval
func TestGetMusic_RetrievesStoredValue(t *testing.T) {
	cache := NewAssetCache()
	key := "battle_music"
	expectedValue := "music_data_789"

	cache.SetMusic(key, expectedValue)
	value, ok := cache.GetMusic(key)

	if !ok {
		t.Error("GetMusic returned false for existing key")
	}
	if value != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, value)
	}
}

// TestGetMusic_ReturnsFalseForMissingKey tests missing music retrieval
func TestGetMusic_ReturnsFalseForMissingKey(t *testing.T) {
	cache := NewAssetCache()

	value, ok := cache.GetMusic("nonexistent")

	if ok {
		t.Error("GetMusic returned true for nonexistent key")
	}
	if value != nil {
		t.Error("GetMusic returned non-nil value for nonexistent key")
	}
}

// TestSetNarrative_StoresValue tests narrative storage
func TestSetNarrative_StoresValue(t *testing.T) {
	cache := NewAssetCache()
	key := "intro_story"
	value := "narrative_data"

	cache.SetNarrative(key, value)

	if len(cache.Narrative) != 1 {
		t.Errorf("Expected 1 narrative, got %d", len(cache.Narrative))
	}
}

// TestGetNarrative_RetrievesStoredValue tests narrative retrieval
func TestGetNarrative_RetrievesStoredValue(t *testing.T) {
	cache := NewAssetCache()
	key := "boss_story"
	expectedValue := "narrative_data_abc"

	cache.SetNarrative(key, expectedValue)
	value, ok := cache.GetNarrative(key)

	if !ok {
		t.Error("GetNarrative returned false for existing key")
	}
	if value != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, value)
	}
}

// TestGetNarrative_ReturnsFalseForMissingKey tests missing narrative retrieval
func TestGetNarrative_ReturnsFalseForMissingKey(t *testing.T) {
	cache := NewAssetCache()

	value, ok := cache.GetNarrative("nonexistent")

	if ok {
		t.Error("GetNarrative returned true for nonexistent key")
	}
	if value != nil {
		t.Error("GetNarrative returned non-nil value for nonexistent key")
	}
}

// TestClear_RemovesAllCachedAssets tests cache clearing
func TestClear_RemovesAllCachedAssets(t *testing.T) {
	cache := NewAssetCache()

	// Populate cache
	cache.SetSprite("sprite1", "data1")
	cache.SetSound("sound1", "data2")
	cache.SetMusic("music1", "data3")
	cache.SetNarrative("narrative1", "data4")

	// Verify data is stored
	if len(cache.Sprites) == 0 || len(cache.Sounds) == 0 ||
		len(cache.Music) == 0 || len(cache.Narrative) == 0 {
		t.Error("Cache should contain data before Clear")
	}

	// Clear cache
	cache.Clear()

	// Verify all maps are empty
	if len(cache.Sprites) != 0 {
		t.Errorf("Expected empty Sprites after Clear, got %d items", len(cache.Sprites))
	}
	if len(cache.Sounds) != 0 {
		t.Errorf("Expected empty Sounds after Clear, got %d items", len(cache.Sounds))
	}
	if len(cache.Music) != 0 {
		t.Errorf("Expected empty Music after Clear, got %d items", len(cache.Music))
	}
	if len(cache.Narrative) != 0 {
		t.Errorf("Expected empty Narrative after Clear, got %d items", len(cache.Narrative))
	}
}

// TestClear_RetrievalAfterClear tests that retrieval fails after clear
func TestClear_RetrievalAfterClear(t *testing.T) {
	cache := NewAssetCache()
	
	cache.SetSprite("sprite1", "data1")
	cache.Clear()
	
	_, ok := cache.GetSprite("sprite1")
	if ok {
		t.Error("GetSprite should return false after Clear")
	}
}

// TestAssetCache_MultipleTypes tests storing multiple asset types
func TestAssetCache_MultipleTypes(t *testing.T) {
	cache := NewAssetCache()

	// Store different types
	cache.SetSprite("sprite1", "sprite_data")
	cache.SetSound("sound1", 12345)
	cache.SetMusic("music1", []byte{1, 2, 3})
	cache.SetNarrative("narrative1", map[string]string{"key": "value"})

	// Retrieve and verify
	sprite, ok := cache.GetSprite("sprite1")
	if !ok || sprite != "sprite_data" {
		t.Error("Failed to retrieve sprite")
	}

	sound, ok := cache.GetSound("sound1")
	if !ok || sound != 12345 {
		t.Error("Failed to retrieve sound")
	}

	_, ok = cache.GetMusic("music1")
	if !ok {
		t.Error("Failed to retrieve music")
	}

	narrative, ok := cache.GetNarrative("narrative1")
	if !ok {
		t.Error("Failed to retrieve narrative")
	}
	narMap, isMap := narrative.(map[string]string)
	if !isMap || narMap["key"] != "value" {
		t.Error("Failed to retrieve correct narrative data")
	}
}

// TestAssetCache_OverwriteValue tests overwriting cached values
func TestAssetCache_OverwriteValue(t *testing.T) {
	cache := NewAssetCache()
	key := "test_key"

	// Set initial value
	cache.SetSprite(key, "value1")
	value1, _ := cache.GetSprite(key)
	if value1 != "value1" {
		t.Error("Initial value not set correctly")
	}

	// Overwrite with new value
	cache.SetSprite(key, "value2")
	value2, _ := cache.GetSprite(key)
	if value2 != "value2" {
		t.Error("Value not overwritten correctly")
	}

	// Should only have one entry
	if len(cache.Sprites) != 1 {
		t.Errorf("Expected 1 sprite after overwrite, got %d", len(cache.Sprites))
	}
}

// TestAssetCache_ConcurrentAccess tests thread-safe concurrent access
func TestAssetCache_ConcurrentAccess(t *testing.T) {
	cache := NewAssetCache()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := "concurrent_key"
			cache.SetSprite(key, index)
		}(i)
	}

	wg.Wait()

	// Should have one entry (last write wins)
	value, ok := cache.GetSprite("concurrent_key")
	if !ok {
		t.Error("Failed to retrieve concurrently written value")
	}
	if value == nil {
		t.Error("Retrieved nil value")
	}
}

// TestAssetCache_ConcurrentReads tests concurrent read access
func TestAssetCache_ConcurrentReads(t *testing.T) {
	cache := NewAssetCache()
	cache.SetSprite("read_key", "test_value")

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, ok := cache.GetSprite("read_key")
			if !ok {
				errors <- nil
			}
			if value != "test_value" {
				errors <- nil
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for range errors {
		errorCount++
	}
	if errorCount > 0 {
		t.Errorf("Had %d read errors during concurrent access", errorCount)
	}
}

// TestAssetCache_ConcurrentMixedOperations tests mixed concurrent operations
func TestAssetCache_ConcurrentMixedOperations(t *testing.T) {
	cache := NewAssetCache()
	var wg sync.WaitGroup

	// Mix of reads, writes, and clears
	for i := 0; i < 50; i++ {
		wg.Add(3)
		
		go func(index int) {
			defer wg.Done()
			cache.SetSprite("key", index)
		}(i)
		
		go func() {
			defer wg.Done()
			cache.GetSprite("key")
		}()
		
		go func() {
			defer wg.Done()
			if i%10 == 0 {
				cache.Clear()
			}
		}()
	}

	wg.Wait()
	
	// Should not panic and cache should be in valid state
	cache.SetSprite("final", "value")
	_, ok := cache.GetSprite("final")
	if !ok {
		t.Error("Cache in invalid state after concurrent operations")
	}
}

// TestAssetCache_EmptyKeyHandling tests handling of empty keys
func TestAssetCache_EmptyKeyHandling(t *testing.T) {
	cache := NewAssetCache()
	
	// Should be able to use empty string as key
	cache.SetSprite("", "empty_key_value")
	value, ok := cache.GetSprite("")
	
	if !ok {
		t.Error("Should be able to use empty string as key")
	}
	if value != "empty_key_value" {
		t.Error("Failed to retrieve value with empty key")
	}
}

// TestAssetCache_NilValueHandling tests handling of nil values
func TestAssetCache_NilValueHandling(t *testing.T) {
	cache := NewAssetCache()
	
	// Should be able to store nil
	cache.SetSprite("nil_key", nil)
	value, ok := cache.GetSprite("nil_key")
	
	if !ok {
		t.Error("Should return true for stored nil value")
	}
	if value != nil {
		t.Error("Should retrieve nil value")
	}
}

// TestAssetCache_LargeNumberOfEntries tests cache with many entries
func TestAssetCache_LargeNumberOfEntries(t *testing.T) {
	cache := NewAssetCache()
	
	count := 1000
	for i := 0; i < count; i++ {
		cache.SetSprite(string(rune(i)), i)
	}
	
	if len(cache.Sprites) != count {
		t.Errorf("Expected %d sprites, got %d", count, len(cache.Sprites))
	}
	
	// Verify random access works
	value, ok := cache.GetSprite(string(rune(500)))
	if !ok {
		t.Error("Failed to retrieve entry from large cache")
	}
	if value != 500 {
		t.Errorf("Expected value 500, got %v", value)
	}
}
