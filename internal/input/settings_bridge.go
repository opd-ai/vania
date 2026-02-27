// Package input provides helpers to integrate with the settings package
package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/vania/internal/settings"
)

// KeyMappingFromSettings converts settings.ControlSettings to input.KeyMapping.
// This bridges the two packages, allowing the input handler to use
// key bindings configured in the settings system.
func KeyMappingFromSettings(controls *settings.ControlSettings) *KeyMapping {
	mapping := &KeyMapping{
		MoveLeft:   []ebiten.Key{},
		MoveRight:  []ebiten.Key{},
		Jump:       []ebiten.Key{},
		Attack:     []ebiten.Key{},
		Dash:       []ebiten.Key{},
		UseAbility: []ebiten.Key{},
		Pause:      []ebiten.Key{},
	}

	// Map each ControlAction to the corresponding KeyMapping field
	for action, key := range controls.KeyBindings {
		switch action {
		case settings.ActionMoveLeft:
			mapping.MoveLeft = append(mapping.MoveLeft, key)
		case settings.ActionMoveRight:
			mapping.MoveRight = append(mapping.MoveRight, key)
		case settings.ActionJump:
			mapping.Jump = append(mapping.Jump, key)
		case settings.ActionAttack:
			mapping.Attack = append(mapping.Attack, key)
		case settings.ActionDash:
			mapping.Dash = append(mapping.Dash, key)
		case settings.ActionInteract:
			// Map Interact to UseAbility
			mapping.UseAbility = append(mapping.UseAbility, key)
		case settings.ActionPause:
			mapping.Pause = append(mapping.Pause, key)
		case settings.ActionMenu:
			// Also map Menu to Pause for menu navigation
			mapping.Pause = append(mapping.Pause, key)
		}
	}

	return mapping
}
