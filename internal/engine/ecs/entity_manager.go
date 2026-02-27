// Package ecs provides the core Entity Component System framework for VANIA.
package ecs

import (
	"fmt"
	"sync"
)

// EntityManager manages entity lifecycle (spawn, despawn, pooling)
type EntityManager struct {
	entities       map[EntityID]*Entity
	nextID         EntityID
	entityPool     []*Entity
	maxPoolSize    int
	activeEntities int
	mu             sync.RWMutex
}

// NewEntityManager creates a new entity manager
func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities:    make(map[EntityID]*Entity),
		nextID:      1,
		entityPool:  make([]*Entity, 0),
		maxPoolSize: 1000,
	}
}

// Spawn creates a new entity or reuses one from the pool
func (em *EntityManager) Spawn() *Entity {
	em.mu.Lock()
	defer em.mu.Unlock()

	var entity *Entity

	// Try to reuse from pool
	if len(em.entityPool) > 0 {
		entity = em.entityPool[len(em.entityPool)-1]
		em.entityPool = em.entityPool[:len(em.entityPool)-1]

		// Reset entity
		entity.Active = true
		entity.Components = make(map[ComponentType]Component)
	} else {
		// Create new entity
		entity = NewEntity(em.nextID)
		em.nextID++
	}

	em.entities[entity.ID] = entity
	em.activeEntities++

	return entity
}

// Despawn deactivates an entity and returns it to the pool
func (em *EntityManager) Despawn(entityID EntityID) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	entity, exists := em.entities[entityID]
	if !exists {
		return fmt.Errorf("entity %d does not exist", entityID)
	}

	if !entity.Active {
		return fmt.Errorf("entity %d is already inactive", entityID)
	}

	entity.Active = false
	em.activeEntities--

	// Return to pool if not full
	if len(em.entityPool) < em.maxPoolSize {
		em.entityPool = append(em.entityPool, entity)
	}

	// Remove from active entities
	delete(em.entities, entityID)

	return nil
}

// Get retrieves an entity by ID
func (em *EntityManager) Get(entityID EntityID) (*Entity, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entity, exists := em.entities[entityID]
	if !exists {
		return nil, fmt.Errorf("entity %d does not exist", entityID)
	}

	return entity, nil
}

// GetAll returns all active entities
func (em *EntityManager) GetAll() []*Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entities := make([]*Entity, 0, len(em.entities))
	for _, entity := range em.entities {
		if entity.Active {
			entities = append(entities, entity)
		}
	}

	return entities
}

// GetWithComponent returns all entities that have a specific component type
func (em *EntityManager) GetWithComponent(componentType ComponentType) []*Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entities := make([]*Entity, 0)
	for _, entity := range em.entities {
		if entity.Active && entity.HasComponent(componentType) {
			entities = append(entities, entity)
		}
	}

	return entities
}

// GetActiveCount returns the number of active entities
func (em *EntityManager) GetActiveCount() int {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.activeEntities
}

// GetPoolSize returns the current pool size
func (em *EntityManager) GetPoolSize() int {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return len(em.entityPool)
}

// SetMaxPoolSize sets the maximum pool size
func (em *EntityManager) SetMaxPoolSize(size int) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.maxPoolSize = size

	// Trim pool if necessary
	if len(em.entityPool) > size {
		em.entityPool = em.entityPool[:size]
	}
}

// Clear removes all entities and clears the pool
func (em *EntityManager) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.entities = make(map[EntityID]*Entity)
	em.entityPool = make([]*Entity, 0)
	em.activeEntities = 0
}
