package ecs

import "testing"

// MockComponent is a test component
type MockComponent struct {
	componentType ComponentType
	data          int
}

func (mc *MockComponent) Type() ComponentType {
	return mc.componentType
}

func TestNewEntity(t *testing.T) {
	entity := NewEntity(42)

	if entity.ID != 42 {
		t.Errorf("Expected entity ID 42, got %d", entity.ID)
	}

	if !entity.Active {
		t.Error("Expected entity to be active by default")
	}

	if entity.Components == nil {
		t.Error("Expected components map to be initialized")
	}

	if len(entity.Components) != 0 {
		t.Errorf("Expected empty components map, got %d components", len(entity.Components))
	}
}

func TestEntityAddComponent(t *testing.T) {
	entity := NewEntity(1)
	component := &MockComponent{
		componentType: ComponentTypeTransform,
		data:          100,
	}

	entity.AddComponent(component)

	if len(entity.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(entity.Components))
	}

	if !entity.HasComponent(ComponentTypeTransform) {
		t.Error("Expected entity to have Transform component")
	}
}

func TestEntityGetComponent(t *testing.T) {
	entity := NewEntity(1)
	component := &MockComponent{
		componentType: ComponentTypeHealth,
		data:          100,
	}

	entity.AddComponent(component)

	retrieved, exists := entity.GetComponent(ComponentTypeHealth)
	if !exists {
		t.Error("Expected component to exist")
	}

	mockComp, ok := retrieved.(*MockComponent)
	if !ok {
		t.Error("Expected component to be MockComponent type")
	}

	if mockComp.data != 100 {
		t.Errorf("Expected component data to be 100, got %d", mockComp.data)
	}
}

func TestEntityGetComponentNotExists(t *testing.T) {
	entity := NewEntity(1)

	_, exists := entity.GetComponent(ComponentTypePhysics)
	if exists {
		t.Error("Expected component to not exist")
	}
}

func TestEntityRemoveComponent(t *testing.T) {
	entity := NewEntity(1)
	component := &MockComponent{
		componentType: ComponentTypeSprite,
		data:          50,
	}

	entity.AddComponent(component)

	if !entity.HasComponent(ComponentTypeSprite) {
		t.Error("Expected entity to have Sprite component")
	}

	entity.RemoveComponent(ComponentTypeSprite)

	if entity.HasComponent(ComponentTypeSprite) {
		t.Error("Expected entity to not have Sprite component after removal")
	}

	if len(entity.Components) != 0 {
		t.Errorf("Expected 0 components after removal, got %d", len(entity.Components))
	}
}

func TestEntityMultipleComponents(t *testing.T) {
	entity := NewEntity(1)

	transformComp := &MockComponent{componentType: ComponentTypeTransform, data: 1}
	physicsComp := &MockComponent{componentType: ComponentTypePhysics, data: 2}
	spriteComp := &MockComponent{componentType: ComponentTypeSprite, data: 3}

	entity.AddComponent(transformComp)
	entity.AddComponent(physicsComp)
	entity.AddComponent(spriteComp)

	if len(entity.Components) != 3 {
		t.Errorf("Expected 3 components, got %d", len(entity.Components))
	}

	if !entity.HasComponent(ComponentTypeTransform) {
		t.Error("Expected entity to have Transform component")
	}
	if !entity.HasComponent(ComponentTypePhysics) {
		t.Error("Expected entity to have Physics component")
	}
	if !entity.HasComponent(ComponentTypeSprite) {
		t.Error("Expected entity to have Sprite component")
	}
}

func TestEntityComponentReplacement(t *testing.T) {
	entity := NewEntity(1)

	comp1 := &MockComponent{componentType: ComponentTypeHealth, data: 100}
	entity.AddComponent(comp1)

	comp2 := &MockComponent{componentType: ComponentTypeHealth, data: 200}
	entity.AddComponent(comp2)

	if len(entity.Components) != 1 {
		t.Errorf("Expected 1 component (replaced), got %d", len(entity.Components))
	}

	retrieved, _ := entity.GetComponent(ComponentTypeHealth)
	mockComp := retrieved.(*MockComponent)

	if mockComp.data != 200 {
		t.Errorf("Expected component data to be 200 (replaced), got %d", mockComp.data)
	}
}
