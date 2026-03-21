// Package engine provides the game runner that integrates Ebiten rendering
// with the procedural generation system, handling the game loop, player
// movement, and visual display.
package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/vania/internal/particle"
	"github.com/opd-ai/vania/internal/render"
)

// ParticleECSSystem is an ECS-compatible wrapper that manages particle
// simulation updates and rendering through the SystemManager infrastructure.
type ParticleECSSystem struct {
	system   *particle.ParticleSystem
	renderer *render.Renderer
}

// NewParticleECSSystem creates a new ParticleECSSystem wrapping the given
// particle system and renderer.
func NewParticleECSSystem(ps *particle.ParticleSystem, r *render.Renderer) *ParticleECSSystem {
	return &ParticleECSSystem{system: ps, renderer: r}
}

// Update advances all active particle emitters by one simulation step.
func (p *ParticleECSSystem) Update(_ float64) error {
	if p.system != nil {
		p.system.Update()
	}
	return nil
}

// Draw renders all active particles to the screen.
func (p *ParticleECSSystem) Draw(screen *ebiten.Image) {
	if p.system != nil && p.renderer != nil {
		p.renderer.RenderParticles(screen, p.system.GetAllParticles())
	}
}

// SetGenre is a no-op for the particle system (particles are genre-neutral).
func (p *ParticleECSSystem) SetGenre(_ string) {}
