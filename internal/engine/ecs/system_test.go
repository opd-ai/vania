package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// MockSystem is a test system that implements the System interface
type MockSystem struct {
	updateCalled bool
	drawCalled   bool
	currentGenre string
	updateError  error
}

func (ms *MockSystem) Update(dt float64) error {
	ms.updateCalled = true
	return ms.updateError
}

func (ms *MockSystem) Draw(screen *ebiten.Image) {
	ms.drawCalled = true
}

func (ms *MockSystem) SetGenre(genreID string) {
	ms.currentGenre = genreID
}

func TestMockSystemImplementsSystem(t *testing.T) {
	var _ System = (*MockSystem)(nil)
}

func TestMockSystemImplementsGenreSwitcher(t *testing.T) {
	var _ GenreSwitcher = (*MockSystem)(nil)
}

func TestSystemUpdate(t *testing.T) {
	system := &MockSystem{}

	if system.updateCalled {
		t.Error("Expected update not to be called yet")
	}

	err := system.Update(0.016)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !system.updateCalled {
		t.Error("Expected update to be called")
	}
}

func TestSystemDraw(t *testing.T) {
	system := &MockSystem{}

	if system.drawCalled {
		t.Error("Expected draw not to be called yet")
	}

	system.Draw(nil)

	if !system.drawCalled {
		t.Error("Expected draw to be called")
	}
}

func TestSystemSetGenre(t *testing.T) {
	system := &MockSystem{}

	if system.currentGenre != "" {
		t.Errorf("Expected empty genre, got %s", system.currentGenre)
	}

	system.SetGenre("fantasy")

	if system.currentGenre != "fantasy" {
		t.Errorf("Expected genre 'fantasy', got %s", system.currentGenre)
	}
}

func TestGenreSwitcherInterface(t *testing.T) {
	system := &MockSystem{}

	var switcher GenreSwitcher = system

	switcher.SetGenre("scifi")

	if system.currentGenre != "scifi" {
		t.Errorf("Expected genre 'scifi', got %s", system.currentGenre)
	}
}
