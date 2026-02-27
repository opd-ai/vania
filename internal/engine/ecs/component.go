// Package ecs provides the core Entity Component System framework for VANIA,
// enabling modular game systems with genre-switching capabilities.
package ecs

// ComponentType identifies different component types
type ComponentType int

const (
	ComponentTypeTransform ComponentType = iota
	ComponentTypePhysics
	ComponentTypeSprite
	ComponentTypeHealth
	ComponentTypeAI
	ComponentTypeAbility
	ComponentTypeInventory
	ComponentTypeCombat
	ComponentTypeAnimation
	ComponentTypeAudio
)

// String returns the human-readable name of a component type
func (ct ComponentType) String() string {
	switch ct {
	case ComponentTypeTransform:
		return "Transform"
	case ComponentTypePhysics:
		return "Physics"
	case ComponentTypeSprite:
		return "Sprite"
	case ComponentTypeHealth:
		return "Health"
	case ComponentTypeAI:
		return "AI"
	case ComponentTypeAbility:
		return "Ability"
	case ComponentTypeInventory:
		return "Inventory"
	case ComponentTypeCombat:
		return "Combat"
	case ComponentTypeAnimation:
		return "Animation"
	case ComponentTypeAudio:
		return "Audio"
	default:
		return "Unknown"
	}
}

// Component represents data associated with an entity
type Component interface {
	// Type returns the component type identifier
	Type() ComponentType
}
