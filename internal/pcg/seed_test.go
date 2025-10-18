package pcg

import (
	"testing"
)

func TestSeedDeterminism(t *testing.T) {
	seed := int64(12345)
	
	// Generate twice with same seed
	seeds1 := DeriveSeeds(seed)
	seeds2 := DeriveSeeds(seed)
	
	// Should be identical
	for key, val1 := range seeds1 {
		val2, ok := seeds2[key]
		if !ok {
			t.Errorf("Key %s not found in second generation", key)
		}
		if val1 != val2 {
			t.Errorf("Seed mismatch for %s: %d != %d", key, val1, val2)
		}
	}
}

func TestHashSeed(t *testing.T) {
	masterSeed := int64(42)
	
	// Same inputs should produce same output
	result1 := HashSeed(masterSeed, "test")
	result2 := HashSeed(masterSeed, "test")
	
	if result1 != result2 {
		t.Errorf("HashSeed not deterministic: %d != %d", result1, result2)
	}
	
	// Different inputs should produce different outputs
	result3 := HashSeed(masterSeed, "different")
	if result1 == result3 {
		t.Errorf("HashSeed collision: same output for different inputs")
	}
}

func TestPCGContext(t *testing.T) {
	seed := int64(999)
	ctx := NewPCGContext(seed)
	
	if ctx.Seed != seed {
		t.Errorf("Seed mismatch: %d != %d", ctx.Seed, seed)
	}
	
	if ctx.RNG == nil {
		t.Error("RNG not initialized")
	}
	
	if ctx.Cache == nil {
		t.Error("Cache not initialized")
	}
	
	if ctx.Constraints == nil {
		t.Error("Constraints not initialized")
	}
}

func TestAssetCache(t *testing.T) {
	cache := NewAssetCache()
	
	// Test sprite caching
	testSprite := "test_sprite_data"
	cache.SetSprite("player", testSprite)
	
	retrieved, ok := cache.GetSprite("player")
	if !ok {
		t.Error("Failed to retrieve cached sprite")
	}
	
	if retrieved != testSprite {
		t.Error("Retrieved sprite doesn't match original")
	}
	
	// Test non-existent key
	_, ok = cache.GetSprite("nonexistent")
	if ok {
		t.Error("Should not find non-existent sprite")
	}
	
	// Test clear
	cache.Clear()
	_, ok = cache.GetSprite("player")
	if ok {
		t.Error("Cache should be empty after clear")
	}
}
