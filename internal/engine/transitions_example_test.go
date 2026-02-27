package engine_test

import (
	"fmt"
	"testing"

	"github.com/opd-ai/vania/internal/engine"
	"github.com/opd-ai/vania/internal/world"
)

// Example showing how to configure transition types
func ExampleRoomTransitionHandler_SetTransitionType() {
	game := &engine.Game{
		CurrentRoom: &world.Room{ID: 1},
	}
	handler := engine.NewRoomTransitionHandler(game)

	// Set fade transition (default)
	handler.SetTransitionType(engine.TransitionFade)
	fmt.Println("Transition type:", handler.GetTransitionType())

	// Set slide transition
	handler.SetTransitionType(engine.TransitionSlide)
	fmt.Println("Transition type:", handler.GetTransitionType())

	// Set iris transition
	handler.SetTransitionType(engine.TransitionIris)
	fmt.Println("Transition type:", handler.GetTransitionType())

	// Output:
	// Transition type: fade
	// Transition type: slide
	// Transition type: iris
}

// Example showing how to configure transition duration
func ExampleRoomTransitionHandler_SetTransitionDuration() {
	game := &engine.Game{
		CurrentRoom: &world.Room{ID: 1},
	}
	handler := engine.NewRoomTransitionHandler(game)

	// Set transition duration to 0.5 seconds
	handler.SetTransitionDuration(0.5)
	fmt.Println("Duration set to 0.5s (30 frames at 60 FPS)")

	// Duration is clamped to 0.3-0.8 seconds
	handler.SetTransitionDuration(0.1) // Will be clamped to 0.3s
	fmt.Println("Duration clamped to minimum 0.3s")

	handler.SetTransitionDuration(1.5) // Will be clamped to 0.8s
	fmt.Println("Duration clamped to maximum 0.8s")

	// Output:
	// Duration set to 0.5s (30 frames at 60 FPS)
	// Duration clamped to minimum 0.3s
	// Duration clamped to maximum 0.8s
}

// Test demonstrating transition type selection based on door properties
func TestTransitionTypeSelection(t *testing.T) {
	room1 := &world.Room{ID: 1, Type: world.CombatRoom}
	room2 := &world.Room{ID: 2, Type: world.TreasureRoom}

	game := &engine.Game{
		CurrentRoom: room1,
		Player:      &engine.Player{X: 100, Y: 200},
	}

	handler := engine.NewRoomTransitionHandler(game)

	// Configure for different types of transitions
	transitionTests := []struct {
		name           string
		transitionType engine.TransitionType
		duration       float64
		doorDirection  string
	}{
		{"Quick fade", engine.TransitionFade, 0.3, "east"},
		{"Slow slide", engine.TransitionSlide, 0.8, "west"},
		{"Medium iris", engine.TransitionIris, 0.5, "north"},
	}

	for _, tt := range transitionTests {
		t.Run(tt.name, func(t *testing.T) {
			handler.SetTransitionType(tt.transitionType)
			handler.SetTransitionDuration(tt.duration)

			door := &world.Door{
				Direction: tt.doorDirection,
				LeadsTo:   room2,
			}

			handler.StartTransition(door)

			if !handler.IsTransitioning() {
				t.Error("Transition should be active after StartTransition")
			}

			if handler.GetTransitionType() != tt.transitionType {
				t.Errorf("Transition type = %v, want %v", handler.GetTransitionType(), tt.transitionType)
			}
		})
	}
}
