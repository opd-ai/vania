package entity

import "testing"

func TestBossGeneratorGrantsAbility(t *testing.T) {
	gen := NewBossGenerator(12345)

	// Test without ability
	boss1 := gen.Generate("cave", 42)
	if boss1.GrantsAbility != "" {
		t.Errorf("Expected no ability, got %s", boss1.GrantsAbility)
	}

	// Test with ability
	boss2 := gen.GenerateWithAbility("cave", 42, "double_jump")
	if boss2.GrantsAbility != "double_jump" {
		t.Errorf("Expected double_jump ability, got %s", boss2.GrantsAbility)
	}
}

func TestBossGeneratorDeterminism(t *testing.T) {
	gen1 := NewBossGenerator(42)
	gen2 := NewBossGenerator(42)

	boss1 := gen1.GenerateWithAbility("forest", 100, "dash")
	boss2 := gen2.GenerateWithAbility("forest", 100, "dash")

	if boss1.Name != boss2.Name {
		t.Errorf("Boss names differ: %s vs %s", boss1.Name, boss2.Name)
	}
	if boss1.Health != boss2.Health {
		t.Errorf("Boss health differs: %d vs %d", boss1.Health, boss2.Health)
	}
	if len(boss1.Phases) != len(boss2.Phases) {
		t.Errorf("Boss phase count differs: %d vs %d", len(boss1.Phases), len(boss2.Phases))
	}
	if boss1.GrantsAbility != boss2.GrantsAbility {
		t.Errorf("GrantsAbility differs: %s vs %s", boss1.GrantsAbility, boss2.GrantsAbility)
	}
}
