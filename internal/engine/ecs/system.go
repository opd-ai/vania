// Package ecs provides the core Entity Component System framework for VANIA.
package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// System processes entities with specific component combinations
type System interface {
	// Update processes game logic for the system
	Update(dt float64) error

	// Draw renders visual elements for the system
	Draw(screen *ebiten.Image)

	// SetGenre switches the system's thematic presentation
	SetGenre(genreID string)
}

// GenreSwitcher defines the interface for systems that support genre switching
type GenreSwitcher interface {
	// SetGenre changes the thematic presentation of the system
	// genreID: "fantasy" | "scifi" | "horror" | "cyberpunk" | "postapoc"
	SetGenre(genreID string)
}
