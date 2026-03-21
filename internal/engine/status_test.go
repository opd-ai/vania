package engine

import "testing"

func TestNewStatusManager(t *testing.T) {
	sm := NewStatusManager()
	if sm == nil {
		t.Fatal("Expected non-nil StatusManager")
	}
	if len(sm.ActiveEffects()) != 0 {
		t.Error("Expected no active effects on new StatusManager")
	}
}

func TestStatusManagerApply(t *testing.T) {
	sm := NewStatusManager()
	sm.Apply(StatusBurn, 5.0, "player")

	effects := sm.ActiveEffects()
	if len(effects) != 1 {
		t.Fatalf("Expected 1 effect, got %d", len(effects))
	}
	if effects[0].Type != StatusBurn {
		t.Errorf("Expected Burn, got %v", effects[0].Type)
	}
	if effects[0].Stacks != 1 {
		t.Errorf("Expected 1 stack, got %d", effects[0].Stacks)
	}
}

func TestStatusManagerStacking(t *testing.T) {
	sm := NewStatusManager()
	sm.Apply(StatusPoison, 3.0, "player")
	sm.Apply(StatusPoison, 5.0, "player") // Should stack, not duplicate

	effects := sm.ActiveEffects()
	if len(effects) != 1 {
		t.Fatalf("Expected 1 stacked effect, got %d", len(effects))
	}
	if effects[0].Stacks != 2 {
		t.Errorf("Expected 2 stacks, got %d", effects[0].Stacks)
	}
	// Duration should be refreshed to max
	if effects[0].Duration < 4.9 {
		t.Errorf("Expected duration refreshed to ~5.0, got %f", effects[0].Duration)
	}
}

func TestStatusManagerStackCap(t *testing.T) {
	sm := NewStatusManager()
	for i := 0; i < 10; i++ {
		sm.Apply(StatusBurn, 10.0, "player")
	}
	effects := sm.ActiveEffects()
	if effects[0].Stacks > 5 {
		t.Errorf("Stacks should be capped at 5, got %d", effects[0].Stacks)
	}
}

func TestStatusManagerExpiry(t *testing.T) {
	sm := NewStatusManager()
	sm.Apply(StatusSlow, 0.01, "enemy") // Very short duration

	// Advance past expiry
	for i := 0; i < 10; i++ {
		sm.Update(0.01)
	}

	if len(sm.ActiveEffects()) != 0 {
		t.Error("Expected Slow to have expired")
	}
}

func TestStatusManagerSpeedMultiplier(t *testing.T) {
	tests := []struct {
		name    StatusType
		wantLow bool // true = expect mult < 1 (slowing effect)
	}{
		{StatusFreeze, true},
		{StatusSlow, true},
		{StatusHaste, false},
	}
	for _, tc := range tests {
		sm := NewStatusManager()
		sm.Apply(tc.name, 10.0, "test")
		mult := sm.SpeedMultiplier()
		if tc.wantLow && mult >= 1.0 {
			t.Errorf("%v should reduce speed, got multiplier %f", tc.name, mult)
		}
		if !tc.wantLow && mult <= 1.0 {
			t.Errorf("%v should increase speed, got multiplier %f", tc.name, mult)
		}
	}
}

func TestStatusManagerDisplayName(t *testing.T) {
	sm := NewStatusManager()
	sm.SetGenre("scifi")
	name := sm.DisplayName(StatusBurn)
	if name != "Overheating" {
		t.Errorf("Expected 'Overheating' for scifi Burn, got %q", name)
	}

	sm.SetGenre("fantasy")
	name = sm.DisplayName(StatusPoison)
	if name != "Poisoned" {
		t.Errorf("Expected 'Poisoned' for fantasy Poison, got %q", name)
	}
}

func TestStatusManagerClear(t *testing.T) {
	sm := NewStatusManager()
	sm.Apply(StatusBurn, 5.0, "player")
	sm.Apply(StatusPoison, 3.0, "enemy")
	sm.Clear()
	if len(sm.ActiveEffects()) != 0 {
		t.Errorf("Expected empty after Clear, got %d effects", len(sm.ActiveEffects()))
	}
}

func TestStatusManagerPeriodicDamage(t *testing.T) {
	sm := NewStatusManager()
	sm.Apply(StatusBurn, 10.0, "test")
	// Force the tick timer to expire immediately
	sm.effects[0].tickTimer = 1
	dmg := sm.Update(1.0 / 60.0)
	if dmg <= 0 {
		t.Errorf("Expected positive damage from Burn tick, got %d", dmg)
	}
}

func TestStatusManagerSetGenreFallback(t *testing.T) {
	sm := NewStatusManager()
	sm.SetGenre("unknown_genre")
	// Should not panic, should return non-empty string
	name := sm.DisplayName(StatusBurn)
	if name == "" {
		t.Error("Expected non-empty display name for unknown genre")
	}
}
