// Package engine provides the game runner that integrates Ebiten rendering
// with the procedural generation system, handling the game loop, player
// movement, and visual display.
package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// AudioECSSystem is an ECS-compatible wrapper for the audio subsystem.
// It implements the ecs.System interface to enable genre-aware audio control
// through the SystemManager infrastructure.
type AudioECSSystem struct {
	genre       string
	audioSystem *AudioSystem
}

// NewAudioECSSystem creates a new AudioECSSystem wrapping the given AudioSystem.
func NewAudioECSSystem(audioSys *AudioSystem) *AudioECSSystem {
	return &AudioECSSystem{
		genre:       "fantasy",
		audioSystem: audioSys,
	}
}

// Update ticks the adaptive music tracks each frame.
func (a *AudioECSSystem) Update(dt float64) error {
	if a.audioSystem == nil {
		return nil
	}
	for _, track := range a.audioSystem.AdaptiveTracks {
		track.Update()
	}
	return nil
}

// Draw is a no-op for the audio system (audio produces no visual output).
func (a *AudioECSSystem) Draw(_ *ebiten.Image) {}

// SetGenre updates the genre state for future audio generation parameters.
func (a *AudioECSSystem) SetGenre(genreID string) {
	a.genre = genreID
}
