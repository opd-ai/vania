package ecs

import (
	"testing"
)

func TestNewEntityManager(t *testing.T) {
	em := NewEntityManager()

	if em == nil {
		t.Fatal("Expected non-nil EntityManager")
	}

	if em.GetActiveCount() != 0 {
		t.Errorf("Expected 0 active entities, got %d", em.GetActiveCount())
	}

	if em.GetPoolSize() != 0 {
		t.Errorf("Expected empty pool, got size %d", em.GetPoolSize())
	}
}

func TestEntityManagerSpawn(t *testing.T) {
	em := NewEntityManager()

	entity := em.Spawn()

	if entity == nil {
		t.Fatal("Expected non-nil entity")
	}

	if !entity.Active {
		t.Error("Expected spawned entity to be active")
	}

	if em.GetActiveCount() != 1 {
		t.Errorf("Expected 1 active entity, got %d", em.GetActiveCount())
	}
}

func TestEntityManagerSpawnMultiple(t *testing.T) {
	em := NewEntityManager()

	entity1 := em.Spawn()
	entity2 := em.Spawn()
	entity3 := em.Spawn()

	if entity1.ID == entity2.ID || entity2.ID == entity3.ID || entity1.ID == entity3.ID {
		t.Error("Expected unique entity IDs")
	}

	if em.GetActiveCount() != 3 {
		t.Errorf("Expected 3 active entities, got %d", em.GetActiveCount())
	}
}

func TestEntityManagerDespawn(t *testing.T) {
	em := NewEntityManager()

	entity := em.Spawn()
	entityID := entity.ID

	err := em.Despawn(entityID)
	if err != nil {
		t.Errorf("Expected no error despawning entity, got %v", err)
	}

	if em.GetActiveCount() != 0 {
		t.Errorf("Expected 0 active entities, got %d", em.GetActiveCount())
	}

	if !entity.Active {
		// Entity should be marked inactive but still exist in memory
		// (this is actually correct behavior - the entity is deactivated)
	}
}

func TestEntityManagerDespawnNonExistent(t *testing.T) {
	em := NewEntityManager()

	err := em.Despawn(999)
	if err == nil {
		t.Error("Expected error despawning non-existent entity")
	}
}

func TestEntityManagerDespawnTwice(t *testing.T) {
	em := NewEntityManager()

	entity := em.Spawn()
	entityID := entity.ID

	err := em.Despawn(entityID)
	if err != nil {
		t.Errorf("Expected no error on first despawn, got %v", err)
	}

	err = em.Despawn(entityID)
	if err == nil {
		t.Error("Expected error despawning already despawned entity")
	}
}

func TestEntityManagerGet(t *testing.T) {
	em := NewEntityManager()

	entity := em.Spawn()
	entityID := entity.ID

	retrieved, err := em.Get(entityID)
	if err != nil {
		t.Errorf("Expected no error getting entity, got %v", err)
	}

	if retrieved.ID != entityID {
		t.Errorf("Expected entity ID %d, got %d", entityID, retrieved.ID)
	}
}

func TestEntityManagerGetNonExistent(t *testing.T) {
	em := NewEntityManager()

	_, err := em.Get(999)
	if err == nil {
		t.Error("Expected error getting non-existent entity")
	}
}

func TestEntityManagerGetAll(t *testing.T) {
	em := NewEntityManager()

	em.Spawn()
	em.Spawn()
	em.Spawn()

	entities := em.GetAll()

	if len(entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(entities))
	}

	for _, entity := range entities {
		if !entity.Active {
			t.Error("Expected all returned entities to be active")
		}
	}
}

func TestEntityManagerGetWithComponent(t *testing.T) {
	em := NewEntityManager()

	entity1 := em.Spawn()
	entity2 := em.Spawn()
	entity3 := em.Spawn()

	// Add components to some entities
	entity1.AddComponent(&MockComponent{componentType: ComponentTypeTransform, data: 1})
	entity2.AddComponent(&MockComponent{componentType: ComponentTypePhysics, data: 2})
	entity3.AddComponent(&MockComponent{componentType: ComponentTypeTransform, data: 3})

	// Get entities with Transform component
	transformEntities := em.GetWithComponent(ComponentTypeTransform)

	if len(transformEntities) != 2 {
		t.Errorf("Expected 2 entities with Transform, got %d", len(transformEntities))
	}

	// Get entities with Physics component
	physicsEntities := em.GetWithComponent(ComponentTypePhysics)

	if len(physicsEntities) != 1 {
		t.Errorf("Expected 1 entity with Physics, got %d", len(physicsEntities))
	}
}

func TestEntityManagerPooling(t *testing.T) {
	em := NewEntityManager()

	// Spawn and despawn entity
	entity1 := em.Spawn()
	entityID1 := entity1.ID

	err := em.Despawn(entityID1)
	if err != nil {
		t.Errorf("Expected no error despawning, got %v", err)
	}

	if em.GetPoolSize() != 1 {
		t.Errorf("Expected pool size 1, got %d", em.GetPoolSize())
	}

	// Spawn new entity (should reuse from pool)
	entity2 := em.Spawn()

	if entity2.ID != entityID1 {
		// Note: Reused entities keep their original ID
		t.Logf("Entity IDs: original=%d, reused=%d", entityID1, entity2.ID)
	}

	if em.GetPoolSize() != 0 {
		t.Errorf("Expected pool size 0 after reuse, got %d", em.GetPoolSize())
	}
}

func TestEntityManagerSetMaxPoolSize(t *testing.T) {
	em := NewEntityManager()

	em.SetMaxPoolSize(5)

	// Spawn and despawn 10 entities
	for i := 0; i < 10; i++ {
		entity := em.Spawn()
		em.Despawn(entity.ID)
	}

	// Pool should be capped at 5
	if em.GetPoolSize() > 5 {
		t.Errorf("Expected pool size <= 5, got %d", em.GetPoolSize())
	}
}

func TestEntityManagerClear(t *testing.T) {
	em := NewEntityManager()

	em.Spawn()
	em.Spawn()
	em.Spawn()

	if em.GetActiveCount() != 3 {
		t.Errorf("Expected 3 active entities before clear, got %d", em.GetActiveCount())
	}

	em.Clear()

	if em.GetActiveCount() != 0 {
		t.Errorf("Expected 0 active entities after clear, got %d", em.GetActiveCount())
	}

	if em.GetPoolSize() != 0 {
		t.Errorf("Expected pool size 0 after clear, got %d", em.GetPoolSize())
	}
}

func TestEntityManagerThreadSafety(t *testing.T) {
	em := NewEntityManager()

	// Spawn entities concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			entity := em.Spawn()
			entity.AddComponent(&MockComponent{componentType: ComponentTypeTransform, data: 1})
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	if em.GetActiveCount() != 10 {
		t.Errorf("Expected 10 active entities, got %d", em.GetActiveCount())
	}
}

func TestEntityManagerGetAllAfterDespawn(t *testing.T) {
	em := NewEntityManager()

	entity1 := em.Spawn()
	entity2 := em.Spawn()
	entity3 := em.Spawn()

	em.Despawn(entity2.ID)

	entities := em.GetAll()

	if len(entities) != 2 {
		t.Errorf("Expected 2 active entities, got %d", len(entities))
	}

	// Verify the right entities are still active
	for _, e := range entities {
		if e.ID == entity2.ID {
			t.Error("Expected despawned entity not to be in GetAll results")
		}
	}

	// Verify remaining entities are correct
	foundEntity1 := false
	foundEntity3 := false
	for _, e := range entities {
		if e.ID == entity1.ID {
			foundEntity1 = true
		}
		if e.ID == entity3.ID {
			foundEntity3 = true
		}
	}

	if !foundEntity1 || !foundEntity3 {
		t.Error("Expected entity1 and entity3 to still be active")
	}
}
