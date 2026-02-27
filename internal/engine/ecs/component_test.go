package ecs

import "testing"

func TestComponentType(t *testing.T) {
	testCases := []struct {
		componentType ComponentType
		expectedName  string
	}{
		{ComponentTypeTransform, "Transform"},
		{ComponentTypePhysics, "Physics"},
		{ComponentTypeSprite, "Sprite"},
		{ComponentTypeHealth, "Health"},
		{ComponentTypeAI, "AI"},
		{ComponentTypeAbility, "Ability"},
		{ComponentTypeInventory, "Inventory"},
		{ComponentTypeCombat, "Combat"},
		{ComponentTypeAnimation, "Animation"},
		{ComponentTypeAudio, "Audio"},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedName, func(t *testing.T) {
			if tc.componentType.String() != tc.expectedName {
				t.Errorf("Expected %s, got %s", tc.expectedName, tc.componentType.String())
			}
		})
	}
}

func TestComponentTypeUnknown(t *testing.T) {
	unknownType := ComponentType(999)
	if unknownType.String() != "Unknown" {
		t.Errorf("Expected 'Unknown' for invalid component type, got %s", unknownType.String())
	}
}
