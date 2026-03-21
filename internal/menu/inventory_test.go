package menu

import (
	"testing"

	"github.com/opd-ai/vania/internal/entity"
)

func TestNewInventoryScreen(t *testing.T) {
	is := NewInventoryScreen()
	if is == nil {
		t.Fatal("Expected non-nil InventoryScreen")
	}
	if is.genre != "fantasy" {
		t.Errorf("Expected default genre 'fantasy', got %q", is.genre)
	}
}

func TestInventoryScreenSetItems(t *testing.T) {
	is := NewInventoryScreen()
	items := []*entity.Item{
		{Name: "Sword", Type: entity.WeaponItem, Description: "A sharp blade"},
		{Name: "Potion", Type: entity.ConsumableItem, Description: "Restores 20 HP", Effect: "+20 HP"},
	}
	is.SetItems(items)
	if len(is.items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(is.items))
	}
}

func TestInventoryScreenSetGenre(t *testing.T) {
	is := NewInventoryScreen()
	is.SetGenre("scifi")
	if is.genre != "scifi" {
		t.Errorf("Expected genre 'scifi', got %q", is.genre)
	}
}

func TestInventoryScreenSelectedItem(t *testing.T) {
	is := NewInventoryScreen()
	items := []*entity.Item{
		{Name: "Key", Type: entity.KeyItem},
	}
	is.SetItems(items)

	item := is.SelectedItem()
	if item == nil {
		t.Fatal("Expected selected item at index 0")
	}
	if item.Name != "Key" {
		t.Errorf("Expected item 'Key', got %q", item.Name)
	}
}

func TestInventoryScreenEmptySlot(t *testing.T) {
	is := NewInventoryScreen()
	is.SetItems([]*entity.Item{})
	if is.SelectedItem() != nil {
		t.Error("Expected nil for empty inventory")
	}
}

func TestInventoryScreenSetAbilities(t *testing.T) {
	is := NewInventoryScreen()
	abilities := []entity.Ability{
		{Name: "Double Jump", Type: entity.MovementAbility},
		{Name: "Dash", Type: entity.MovementAbility},
	}
	is.SetAbilities(abilities)
	if len(is.abilities) != 2 {
		t.Errorf("Expected 2 abilities, got %d", len(is.abilities))
	}
}

func TestInventoryScreenGenreTitle(t *testing.T) {
	tests := []struct {
		genre string
		want  string
	}{
		{"fantasy", "INVENTORY"},
		{"scifi", "CARGO BAY"},
		{"horror", "BELONGINGS"},
		{"cyberpunk", "DATA CACHE"},
		{"postapoc", "SCAVENGED ITEMS"},
		{"unknown", "INVENTORY"}, // fallback
	}
	for _, tc := range tests {
		is := NewInventoryScreen()
		is.SetGenre(tc.genre)
		got := is.genreInventoryTitle()
		if got != tc.want {
			t.Errorf("genre=%q: expected title %q, got %q", tc.genre, tc.want, got)
		}
	}
}

func TestItemTypeName(t *testing.T) {
	tests := []struct {
		t    entity.ItemType
		want string
	}{
		{entity.WeaponItem, "Weapon"},
		{entity.ConsumableItem, "Consumable"},
		{entity.KeyItem, "Key Item"},
		{entity.UpgradeItem, "Upgrade"},
		{entity.CurrencyItem, "Currency"},
	}
	for _, tc := range tests {
		got := itemTypeName(tc.t)
		if got != tc.want {
			t.Errorf("itemTypeName(%v) = %q, want %q", tc.t, got, tc.want)
		}
	}
}

func TestInventoryScreenSelectedIndex(t *testing.T) {
	is := NewInventoryScreen()
	is.selectedRow = 2
	is.selectedCol = 1
	if is.SelectedIndex() != 2*inventoryCols+1 {
		t.Errorf("SelectedIndex wrong: got %d", is.SelectedIndex())
	}
}

func TestInventoryScreenRowClampOnSetItems(t *testing.T) {
	is := NewInventoryScreen()
	is.selectedRow = 5
	// Setting fewer items than current selection row
	is.SetItems([]*entity.Item{
		{Name: "Item1", Type: entity.ConsumableItem},
	})
	if is.selectedRow > 0 {
		t.Errorf("Expected selectedRow clamped to 0, got %d", is.selectedRow)
	}
}
