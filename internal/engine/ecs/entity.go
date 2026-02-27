// Package ecs provides the core Entity Component System framework for VANIA.
package ecs

// EntityID is a unique identifier for an entity
type EntityID int

// Entity represents a game object composed of components
type Entity struct {
	ID         EntityID
	Components map[ComponentType]Component
	Active     bool
}

// NewEntity creates a new entity with a unique ID
func NewEntity(id EntityID) *Entity {
	return &Entity{
		ID:         id,
		Components: make(map[ComponentType]Component),
		Active:     true,
	}
}

// AddComponent adds a component to the entity
func (e *Entity) AddComponent(component Component) {
	e.Components[component.Type()] = component
}

// GetComponent retrieves a component by type
func (e *Entity) GetComponent(componentType ComponentType) (Component, bool) {
	comp, exists := e.Components[componentType]
	return comp, exists
}

// RemoveComponent removes a component by type
func (e *Entity) RemoveComponent(componentType ComponentType) {
	delete(e.Components, componentType)
}

// HasComponent checks if entity has a component of the given type
func (e *Entity) HasComponent(componentType ComponentType) bool {
	_, exists := e.Components[componentType]
	return exists
}
