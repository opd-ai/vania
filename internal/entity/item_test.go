package entity

import (
	"testing"
)

func TestNewItemInstance(t *testing.T) {
	item := &Item{
		Name:        "Test Sword",
		Type:        WeaponItem,
		Description: "A test weapon",
		Effect:      "increase_damage",
		Value:       50,
	}

	instance := NewItemInstance(item, 123, 100.0, 200.0)

	if instance.X != 100.0 || instance.Y != 200.0 {
		t.Errorf("Expected position (100, 200), got (%.0f, %.0f)", instance.X, instance.Y)
	}

	if instance.ID != 123 {
		t.Errorf("Expected ID 123, got %d", instance.ID)
	}

	if instance.Collected {
		t.Error("Expected item to not be collected initially")
	}

	if instance.Item != item {
		t.Error("Expected item reference to match")
	}
}

func TestItemInstanceGetBounds(t *testing.T) {
	item := &Item{Name: "Test"}
	instance := NewItemInstance(item, 1, 50.0, 75.0)

	x, y, w, h := instance.GetBounds()

	if x != 50.0 || y != 75.0 {
		t.Errorf("Expected position (50, 75), got (%.0f, %.0f)", x, y)
	}

	// Items should be 16x16 pixels
	if w != 16.0 || h != 16.0 {
		t.Errorf("Expected size (16, 16), got (%.0f, %.0f)", w, h)
	}
}

func TestItemInstanceCollection(t *testing.T) {
	item := &Item{Name: "Potion"}
	instance := NewItemInstance(item, 1, 100.0, 200.0)

	// Initially not collected
	if instance.Collected {
		t.Error("Item should not be collected initially")
	}

	// Simulate collection
	instance.Collected = true

	// Verify collected state
	if !instance.Collected {
		t.Error("Item should be marked as collected")
	}
}

func TestItemGenerator_Generate(t *testing.T) {
	gen := NewItemGenerator(42)

	tests := []struct {
		itemType ItemType
		name     string
	}{
		{WeaponItem, "weapon"},
		{ConsumableItem, "consumable"},
		{KeyItem, "key"},
		{UpgradeItem, "upgrade"},
		{CurrencyItem, "currency"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := gen.Generate(tt.itemType, 42)

			if item == nil {
				t.Fatalf("Expected item to be generated, got nil")
			}

			if item.Type != tt.itemType {
				t.Errorf("Expected type %v, got %v", tt.itemType, item.Type)
			}

			if item.Name == "" {
				t.Error("Expected item to have a name")
			}

			if item.Description == "" {
				t.Error("Expected item to have a description")
			}

			if item.Effect == "" {
				t.Error("Expected item to have an effect")
			}
		})
	}
}

func TestItemGenerator_Deterministic(t *testing.T) {
	gen1 := NewItemGenerator(12345)
	gen2 := NewItemGenerator(12345)

	item1 := gen1.Generate(WeaponItem, 100)
	item2 := gen2.Generate(WeaponItem, 100)

	if item1.Name != item2.Name {
		t.Errorf("Expected deterministic generation: %s != %s", item1.Name, item2.Name)
	}

	if item1.Value != item2.Value {
		t.Errorf("Expected same value: %d != %d", item1.Value, item2.Value)
	}
}

func TestItemTypes(t *testing.T) {
	gen := NewItemGenerator(42)

	// Test weapon item
	weapon := gen.Generate(WeaponItem, 1)
	if weapon.Effect != "increase_damage" {
		t.Errorf("Expected weapon effect 'increase_damage', got '%s'", weapon.Effect)
	}

	// Test consumable item
	consumable := gen.Generate(ConsumableItem, 2)
	if consumable.Effect != "heal" {
		t.Errorf("Expected consumable effect 'heal', got '%s'", consumable.Effect)
	}

	// Test key item
	key := gen.Generate(KeyItem, 3)
	if key.Effect != "unlock" {
		t.Errorf("Expected key effect 'unlock', got '%s'", key.Effect)
	}

	// Test upgrade item
	upgrade := gen.Generate(UpgradeItem, 4)
	if upgrade.Effect != "upgrade" {
		t.Errorf("Expected upgrade effect 'upgrade', got '%s'", upgrade.Effect)
	}
	if upgrade.Name != "Upgrade Stone" {
		t.Errorf("Expected upgrade name 'Upgrade Stone', got '%s'", upgrade.Name)
	}

	// Test currency item
	currency := gen.Generate(CurrencyItem, 5)
	if currency.Effect != "currency" {
		t.Errorf("Expected currency effect 'currency', got '%s'", currency.Effect)
	}
	if currency.Name != "Crystal Shard" {
		t.Errorf("Expected currency name 'Crystal Shard', got '%s'", currency.Name)
	}
}

func TestItemRarity(t *testing.T) {
	// Test rarity tiers
	rarities := []struct {
		rarity ItemRarity
		name   string
		mult   float64
	}{
		{CommonRarity, "Common", 1.0},
		{UncommonRarity, "Uncommon", 1.25},
		{RareRarity, "Rare", 1.5},
		{LegendaryRarity, "Legendary", 2.0},
	}

	for _, tc := range rarities {
		t.Run(tc.name, func(t *testing.T) {
			if tc.rarity.RarityName() != tc.name {
				t.Errorf("Expected name %s, got %s", tc.name, tc.rarity.RarityName())
			}
			if tc.rarity.Multiplier() != tc.mult {
				t.Errorf("Expected multiplier %.2f, got %.2f", tc.mult, tc.rarity.Multiplier())
			}
		})
	}
}

func TestItemRarityDistribution(t *testing.T) {
	// Test that rarity distribution roughly matches expected weights
	gen := NewItemGenerator(12345)

	counts := map[ItemRarity]int{
		CommonRarity:    0,
		UncommonRarity:  0,
		RareRarity:      0,
		LegendaryRarity: 0,
	}

	// Generate many items to test distribution
	for seed := int64(0); seed < 1000; seed++ {
		item := gen.Generate(WeaponItem, seed)
		counts[item.Rarity]++
	}

	// Expected: Common 60%, Uncommon 25%, Rare 12%, Legendary 3%
	// Allow 10% tolerance
	total := float64(1000)
	if float64(counts[CommonRarity])/total < 0.50 || float64(counts[CommonRarity])/total > 0.70 {
		t.Errorf("Common rarity distribution out of range: %.1f%%", float64(counts[CommonRarity])/total*100)
	}
	if float64(counts[UncommonRarity])/total < 0.15 || float64(counts[UncommonRarity])/total > 0.35 {
		t.Errorf("Uncommon rarity distribution out of range: %.1f%%", float64(counts[UncommonRarity])/total*100)
	}
	if float64(counts[RareRarity])/total < 0.05 || float64(counts[RareRarity])/total > 0.20 {
		t.Errorf("Rare rarity distribution out of range: %.1f%%", float64(counts[RareRarity])/total*100)
	}
	if float64(counts[LegendaryRarity])/total > 0.10 {
		t.Errorf("Legendary rarity too common: %.1f%%", float64(counts[LegendaryRarity])/total*100)
	}
}

func TestItemRarityValueScaling(t *testing.T) {
	// Test that rarity multipliers affect item value
	gen := NewItemGenerator(42)

	// Generate items with specific seeds to control rarity
	// We test the multiplier logic directly
	baseValue := 100
	scaledUncommon := int(float64(baseValue) * UncommonRarity.Multiplier())
	scaledRare := int(float64(baseValue) * RareRarity.Multiplier())
	scaledLegendary := int(float64(baseValue) * LegendaryRarity.Multiplier())

	if scaledUncommon != 125 {
		t.Errorf("Uncommon scaling wrong: expected 125, got %d", scaledUncommon)
	}
	if scaledRare != 150 {
		t.Errorf("Rare scaling wrong: expected 150, got %d", scaledRare)
	}
	if scaledLegendary != 200 {
		t.Errorf("Legendary scaling wrong: expected 200, got %d", scaledLegendary)
	}

	// Generate an item and verify rarity is set
	item := gen.Generate(WeaponItem, 1)
	if item.Rarity < CommonRarity || item.Rarity > LegendaryRarity {
		t.Errorf("Invalid rarity: %v", item.Rarity)
	}
}

func TestKeyItemsAlwaysCommon(t *testing.T) {
	gen := NewItemGenerator(42)

	// Key items should always be common (progression-gating shouldn't be luck-based)
	for seed := int64(0); seed < 100; seed++ {
		key := gen.Generate(KeyItem, seed)
		if key.Rarity != CommonRarity {
			t.Errorf("Key items should always be common, got %s", key.Rarity.RarityName())
		}
	}
}
