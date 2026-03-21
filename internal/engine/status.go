// Package engine provides status effect management for combat gameplay.
// Status effects apply periodic damage or movement modifiers with optional
// stacking and genre-specific naming.
package engine

import "math"

// StatusType identifies the mechanical class of a status effect.
type StatusType int

const (
	// StatusBurn deals periodic fire damage (or equivalent per genre).
	StatusBurn StatusType = iota
	// StatusFreeze halves movement speed.
	StatusFreeze
	// StatusShock deals small periodic damage and briefly stuns.
	StatusShock
	// StatusPoison deals periodic damage that ignores armour.
	StatusPoison
	// StatusBleed deals periodic damage scaling with movement speed.
	StatusBleed
	// StatusSlow reduces movement speed by 40%.
	StatusSlow
	// StatusHaste increases movement speed by 40%.
	StatusHaste
)

// statusTickRate is the number of game-frames between damage ticks.
const statusTickRate = 60 // once per second at 60 fps

// genreNames maps each genre string to a per-StatusType display name.
var genreNames = map[string][7]string{
	// [Burn, Freeze, Shock, Poison, Bleed, Slow, Haste]
	"fantasy":   {"Burning", "Frozen", "Shocked", "Poisoned", "Bleeding", "Slowed", "Hasted"},
	"scifi":     {"Overheating", "Cryogenic", "Electrified", "Irradiated", "Haemorrhaging", "Dampened", "Overclocked"},
	"horror":    {"Aflame", "Petrified", "Cursed", "Infected", "Exsanguinating", "Terrified", "Frenzied"},
	"cyberpunk": {"Melting", "Cryo-locked", "Voltage", "Nanobot Swarm", "Hemorrhage", "Lag", "Overclock"},
	"postapoc":  {"Irradiated", "Frost-bitten", "Shocked", "Contaminated", "Wounded", "Crippled", "Stimmed"},
}

// StatusEffect represents a single active status applied to an entity.
type StatusEffect struct {
	Type     StatusType
	Stacks   int     // 1–5; higher stacks increase damage/duration
	Duration float64 // remaining seconds
	Source   string  // e.g. "player", "enemy_poison_dart"
	tickTimer int    // counts frames until next damage tick
}

// StatusManager tracks and resolves all active status effects on a single entity.
type StatusManager struct {
	effects []StatusEffect
	genre   string
}

// NewStatusManager returns an initialised StatusManager.
func NewStatusManager() *StatusManager {
	return &StatusManager{
		effects: make([]StatusEffect, 0, 4),
		genre:   "fantasy",
	}
}

// SetGenre configures the display vocabulary for status effect names.
// Accepted genre IDs: "fantasy", "scifi", "horror", "cyberpunk", "postapoc".
func (sm *StatusManager) SetGenre(genreID string) {
	sm.genre = genreID
}

// Apply adds a status effect, stacking with any existing effect of the same type
// (capped at 5 stacks).  Duration is refreshed to the maximum of the current and
// the incoming duration.
func (sm *StatusManager) Apply(effectType StatusType, duration float64, source string) {
	for i := range sm.effects {
		if sm.effects[i].Type == effectType {
			if sm.effects[i].Stacks < 5 {
				sm.effects[i].Stacks++
			}
			if duration > sm.effects[i].Duration {
				sm.effects[i].Duration = duration
			}
			return
		}
	}
	sm.effects = append(sm.effects, StatusEffect{
		Type:      effectType,
		Stacks:    1,
		Duration:  duration,
		Source:    source,
		tickTimer: statusTickRate,
	})
}

// Update advances all effects by one frame, ticking damage and expiring finished
// effects.  Returns the total periodic damage dealt this frame.
func (sm *StatusManager) Update(dt float64) int {
	totalDmg := 0
	for i := len(sm.effects) - 1; i >= 0; i-- {
		e := &sm.effects[i]
		e.Duration -= dt
		if e.Duration <= 0 {
			sm.effects = append(sm.effects[:i], sm.effects[i+1:]...)
			continue
		}
		e.tickTimer--
		if e.tickTimer <= 0 {
			e.tickTimer = statusTickRate
			totalDmg += sm.tickDamage(e)
		}
	}
	return totalDmg
}

// tickDamage returns the periodic damage for one tick of the given effect.
func (sm *StatusManager) tickDamage(e *StatusEffect) int {
	switch e.Type {
	case StatusBurn:
		return int(math.Round(float64(3*e.Stacks) * 1.0))
	case StatusShock:
		return int(math.Round(float64(2*e.Stacks) * 1.0))
	case StatusPoison:
		return int(math.Round(float64(2*e.Stacks) * 1.2)) // piercing
	case StatusBleed:
		return int(math.Round(float64(e.Stacks) * 1.5))
	default:
		// Freeze, Slow, Haste – no direct damage
		return 0
	}
}

// SpeedMultiplier returns the movement speed modifier imposed by active effects
// (multiplicative; values <1.0 slow, >1.0 haste).
func (sm *StatusManager) SpeedMultiplier() float64 {
	mult := 1.0
	for _, e := range sm.effects {
		switch e.Type {
		case StatusFreeze:
			mult *= 0.5
		case StatusSlow:
			mult *= 0.6
		case StatusHaste:
			mult *= 1.4
		}
	}
	return mult
}

// IsStunned returns true when a Shock effect is active (brief stun on tick).
func (sm *StatusManager) IsStunned() bool {
	for _, e := range sm.effects {
		if e.Type == StatusShock && e.tickTimer == statusTickRate {
			return true
		}
	}
	return false
}

// ActiveEffects returns a copy of the current effect slice for HUD rendering.
func (sm *StatusManager) ActiveEffects() []StatusEffect {
	out := make([]StatusEffect, len(sm.effects))
	copy(out, sm.effects)
	return out
}

// DisplayName returns the genre-appropriate name for the given status type.
func (sm *StatusManager) DisplayName(t StatusType) string {
	names, ok := genreNames[sm.genre]
	if !ok {
		names = genreNames["fantasy"]
	}
	if int(t) < len(names) {
		return names[t]
	}
	return "Unknown"
}

// Clear removes all active status effects (e.g., on room change or respawn).
func (sm *StatusManager) Clear() {
	sm.effects = sm.effects[:0]
}
