// Package particle provides a particle effects system for visual feedback
// in combat, movement, and environmental effects.
package particle

import (
	"image/color"
	"math"
	"math/rand"
)

// ParticleType defines the type of particle
type ParticleType int

const (
	// Combat particles
	HitSpark ParticleType = iota
	DamageNumber
	BloodSplatter
	
	// Movement particles
	DashTrail
	JumpDust
	LandDust
	WalkDust
	
	// Environmental particles
	Rain
	Snow
	Embers
	Sparkles
	Bubbles
	
	// Effect particles
	Explosion
	Smoke
	Lightning
)

// Particle represents a single particle
type Particle struct {
	X, Y           float64
	VelX, VelY     float64
	AccelX, AccelY float64
	Life           int // frames remaining
	MaxLife        int
	Size           float64
	Color          color.RGBA
	Alpha          uint8
	Type           ParticleType
	Rotation       float64
	RotationSpeed  float64
	Data           interface{} // Custom data (e.g., damage number text)
}

// ParticleEmitter generates and manages particles
type ParticleEmitter struct {
	X, Y           float64
	Active         bool
	EmitRate       int     // particles per frame
	EmitTimer      int
	Spread         float64 // angle spread in radians
	Speed          float64 // initial velocity
	SpeedVariance  float64
	Life           int     // particle lifetime in frames
	LifeVariance   int
	Size           float64
	SizeVariance   float64
	Gravity        float64
	Type           ParticleType
	Color          color.RGBA
	OneShot        bool // emit once then deactivate
	Particles      []*Particle
}

// ParticleSystem manages all particle emitters and particles
type ParticleSystem struct {
	emitters []*ParticleEmitter
	particles []*Particle
	maxParticles int
}

// NewParticle creates a new particle
func NewParticle(x, y, velX, velY float64, life int, size float64, col color.RGBA, ptype ParticleType) *Particle {
	return &Particle{
		X:       x,
		Y:       y,
		VelX:    velX,
		VelY:    velY,
		AccelX:  0,
		AccelY:  0,
		Life:    life,
		MaxLife: life,
		Size:    size,
		Color:   col,
		Alpha:   255,
		Type:    ptype,
		Rotation: 0,
		RotationSpeed: 0,
	}
}

// NewParticleEmitter creates a new particle emitter
func NewParticleEmitter(x, y float64, ptype ParticleType) *ParticleEmitter {
	return &ParticleEmitter{
		X:            x,
		Y:            y,
		Active:       false,
		EmitRate:     1,
		EmitTimer:    0,
		Spread:       math.Pi / 2, // 90 degrees
		Speed:        2.0,
		SpeedVariance: 0.5,
		Life:         30, // 0.5 seconds at 60 FPS
		LifeVariance: 10,
		Size:         2.0,
		SizeVariance: 0.5,
		Gravity:      0.2,
		Type:         ptype,
		Color:        color.RGBA{255, 255, 255, 255},
		OneShot:      false,
		Particles:    make([]*Particle, 0),
	}
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(maxParticles int) *ParticleSystem {
	return &ParticleSystem{
		emitters:     make([]*ParticleEmitter, 0),
		particles:    make([]*Particle, 0),
		maxParticles: maxParticles,
	}
}

// Update updates all particles and emitters
func (ps *ParticleSystem) Update() {
	// Update emitters
	for i := len(ps.emitters) - 1; i >= 0; i-- {
		emitter := ps.emitters[i]
		emitter.Update()
		
		// Remove inactive one-shot emitters
		if emitter.OneShot && !emitter.Active && len(emitter.Particles) == 0 {
			ps.emitters = append(ps.emitters[:i], ps.emitters[i+1:]...)
		}
	}
	
	// Update standalone particles
	for i := len(ps.particles) - 1; i >= 0; i-- {
		particle := ps.particles[i]
		particle.Update()
		
		// Remove dead particles
		if particle.Life <= 0 {
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
		}
	}
}

// AddEmitter adds an emitter to the system
func (ps *ParticleSystem) AddEmitter(emitter *ParticleEmitter) {
	if len(ps.emitters) < ps.maxParticles/10 { // Limit emitters to 1/10 of max particles
		ps.emitters = append(ps.emitters, emitter)
	}
}

// AddParticle adds a single particle to the system
func (ps *ParticleSystem) AddParticle(particle *Particle) {
	if len(ps.particles) < ps.maxParticles {
		ps.particles = append(ps.particles, particle)
	}
}

// GetAllParticles returns all active particles (from emitters and standalone)
func (ps *ParticleSystem) GetAllParticles() []*Particle {
	allParticles := make([]*Particle, 0, len(ps.particles))
	
	// Add standalone particles
	allParticles = append(allParticles, ps.particles...)
	
	// Add particles from emitters
	for _, emitter := range ps.emitters {
		allParticles = append(allParticles, emitter.Particles...)
	}
	
	return allParticles
}

// Clear removes all particles and emitters
func (ps *ParticleSystem) Clear() {
	ps.emitters = make([]*ParticleEmitter, 0)
	ps.particles = make([]*Particle, 0)
}

// GetParticleCount returns total number of active particles
func (ps *ParticleSystem) GetParticleCount() int {
	count := len(ps.particles)
	for _, emitter := range ps.emitters {
		count += len(emitter.Particles)
	}
	return count
}

// Update updates the particle state
func (p *Particle) Update() {
	// Apply acceleration
	p.VelX += p.AccelX
	p.VelY += p.AccelY
	
	// Update position
	p.X += p.VelX
	p.Y += p.VelY
	
	// Update rotation
	p.Rotation += p.RotationSpeed
	
	// Decrease life
	p.Life--
	
	// Fade out as life decreases
	if p.Life < p.MaxLife/3 {
		p.Alpha = uint8(float64(p.Alpha) * 0.95)
	}
}

// IsAlive returns true if the particle is still active
func (p *Particle) IsAlive() bool {
	return p.Life > 0
}

// Update updates the emitter and emits particles
func (e *ParticleEmitter) Update() {
	// Always update existing particles, even if emitter is inactive
	for i := len(e.Particles) - 1; i >= 0; i-- {
		particle := e.Particles[i]
		particle.Update()
		
		// Remove dead particles
		if particle.Life <= 0 {
			e.Particles = append(e.Particles[:i], e.Particles[i+1:]...)
		}
	}
	
	// Only emit new particles if active
	if !e.Active {
		return
	}
	
	e.EmitTimer++
	
	// Emit particles based on rate
	if e.EmitTimer >= 60/e.EmitRate {
		e.EmitTimer = 0
		e.EmitParticles(1)
	}
}

// EmitParticles emits a number of particles
func (e *ParticleEmitter) EmitParticles(count int) {
	for i := 0; i < count; i++ {
		// Random angle within spread
		angle := (rand.Float64() - 0.5) * e.Spread
		
		// Random speed with variance
		speed := e.Speed + (rand.Float64()-0.5)*e.SpeedVariance
		
		// Calculate velocity from angle and speed
		velX := math.Cos(angle) * speed
		velY := math.Sin(angle) * speed
		
		// Random life with variance
		life := e.Life + rand.Intn(e.LifeVariance*2) - e.LifeVariance
		if life < 1 {
			life = 1
		}
		
		// Random size with variance
		size := e.Size + (rand.Float64()-0.5)*e.SizeVariance
		if size < 0.5 {
			size = 0.5
		}
		
		// Create particle
		particle := NewParticle(e.X, e.Y, velX, velY, life, size, e.Color, e.Type)
		particle.AccelY = e.Gravity
		
		e.Particles = append(e.Particles, particle)
	}
	
	// Deactivate if one-shot
	if e.OneShot {
		e.Active = false
	}
}

// Start activates the emitter
func (e *ParticleEmitter) Start() {
	e.Active = true
}

// Stop deactivates the emitter
func (e *ParticleEmitter) Stop() {
	e.Active = false
}

// SetPosition updates the emitter position
func (e *ParticleEmitter) SetPosition(x, y float64) {
	e.X = x
	e.Y = y
}

// Burst emits a burst of particles and deactivates
func (e *ParticleEmitter) Burst(count int) {
	e.Active = true
	e.OneShot = true
	e.EmitParticles(count)
	e.Active = false
}
