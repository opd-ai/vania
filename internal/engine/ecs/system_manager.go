// Package ecs provides the core Entity Component System framework for VANIA.
package ecs

import (
	"fmt"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// SystemEntry represents a system with its priority
type SystemEntry struct {
	System   System
	Priority int
}

// SystemManager manages execution of all game systems
type SystemManager struct {
	systems []SystemEntry
	sorted  bool
}

// NewSystemManager creates a new system manager
func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems: make([]SystemEntry, 0),
		sorted:  true,
	}
}

// Register adds a system to the manager with a priority
// Lower priority values execute first
func (sm *SystemManager) Register(system System, priority int) {
	sm.systems = append(sm.systems, SystemEntry{
		System:   system,
		Priority: priority,
	})
	sm.sorted = false
}

// Unregister removes a system from the manager
func (sm *SystemManager) Unregister(system System) bool {
	for i, entry := range sm.systems {
		if entry.System == system {
			sm.systems = append(sm.systems[:i], sm.systems[i+1:]...)
			return true
		}
	}
	return false
}

// ensureSorted sorts systems by priority if needed
func (sm *SystemManager) ensureSorted() {
	if !sm.sorted {
		sort.Slice(sm.systems, func(i, j int) bool {
			return sm.systems[i].Priority < sm.systems[j].Priority
		})
		sm.sorted = true
	}
}

// Update executes all systems' Update methods in priority order
func (sm *SystemManager) Update(dt float64) error {
	sm.ensureSorted()

	for _, entry := range sm.systems {
		if err := entry.System.Update(dt); err != nil {
			return fmt.Errorf("system update error: %w", err)
		}
	}

	return nil
}

// Draw executes all systems' Draw methods in priority order
func (sm *SystemManager) Draw(screen *ebiten.Image) {
	sm.ensureSorted()

	for _, entry := range sm.systems {
		entry.System.Draw(screen)
	}
}

// SetGenre propagates genre changes to all systems
func (sm *SystemManager) SetGenre(genreID string) {
	for _, entry := range sm.systems {
		entry.System.SetGenre(genreID)
	}
}

// GetSystemCount returns the number of registered systems
func (sm *SystemManager) GetSystemCount() int {
	return len(sm.systems)
}

// Clear removes all systems
func (sm *SystemManager) Clear() {
	sm.systems = make([]SystemEntry, 0)
	sm.sorted = true
}
