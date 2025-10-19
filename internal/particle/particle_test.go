package particle

import (
	"image/color"
	"testing"
)

func TestNewParticle(t *testing.T) {
	p := NewParticle(100, 200, 1.5, -2.0, 60, 3.0, color.RGBA{255, 0, 0, 255}, HitSpark)
	
	if p.X != 100 || p.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", p.X, p.Y)
	}
	if p.VelX != 1.5 || p.VelY != -2.0 {
		t.Errorf("Expected velocity (1.5, -2.0), got (%f, %f)", p.VelX, p.VelY)
	}
	if p.Life != 60 || p.MaxLife != 60 {
		t.Errorf("Expected life 60, got Life=%d, MaxLife=%d", p.Life, p.MaxLife)
	}
	if p.Size != 3.0 {
		t.Errorf("Expected size 3.0, got %f", p.Size)
	}
	if p.Type != HitSpark {
		t.Errorf("Expected type HitSpark, got %v", p.Type)
	}
}

func TestParticleUpdate(t *testing.T) {
	p := NewParticle(0, 0, 2.0, 3.0, 60, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
	p.AccelX = 0.1
	p.AccelY = 0.2
	
	initialLife := p.Life
	p.Update()
	
	// Acceleration is applied first: VelX = 2.0 + 0.1 = 2.1, VelY = 3.0 + 0.2 = 3.2
	// Then position updated: X = 0 + 2.1 = 2.1, Y = 0 + 3.2 = 3.2
	if p.X != 2.1 || p.Y != 3.2 {
		t.Errorf("Expected position (2.1, 3.2) after update, got (%f, %f)", p.X, p.Y)
	}
	
	// Check velocity updated with acceleration
	if p.VelX != 2.1 || p.VelY != 3.2 {
		t.Errorf("Expected velocity (2.1, 3.2) after update, got (%f, %f)", p.VelX, p.VelY)
	}
	
	// Check life decreased
	if p.Life != initialLife-1 {
		t.Errorf("Expected life to decrease by 1, got %d", p.Life)
	}
}

func TestParticleIsAlive(t *testing.T) {
	p := NewParticle(0, 0, 0, 0, 5, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
	
	if !p.IsAlive() {
		t.Error("Particle should be alive with life > 0")
	}
	
	// Update until dead
	for i := 0; i < 10; i++ {
		p.Update()
	}
	
	if p.IsAlive() {
		t.Error("Particle should be dead with life <= 0")
	}
}

func TestNewParticleEmitter(t *testing.T) {
	e := NewParticleEmitter(150, 250, DashTrail)
	
	if e.X != 150 || e.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", e.X, e.Y)
	}
	if e.Active {
		t.Error("Emitter should not be active by default")
	}
	if e.Type != DashTrail {
		t.Errorf("Expected type DashTrail, got %v", e.Type)
	}
	if len(e.Particles) != 0 {
		t.Errorf("Expected 0 particles initially, got %d", len(e.Particles))
	}
}

func TestEmitterStartStop(t *testing.T) {
	e := NewParticleEmitter(0, 0, HitSpark)
	
	if e.Active {
		t.Error("Emitter should start inactive")
	}
	
	e.Start()
	if !e.Active {
		t.Error("Emitter should be active after Start()")
	}
	
	e.Stop()
	if e.Active {
		t.Error("Emitter should be inactive after Stop()")
	}
}

func TestEmitterSetPosition(t *testing.T) {
	e := NewParticleEmitter(0, 0, HitSpark)
	e.SetPosition(100, 200)
	
	if e.X != 100 || e.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", e.X, e.Y)
	}
}

func TestEmitterEmitParticles(t *testing.T) {
	e := NewParticleEmitter(0, 0, HitSpark)
	e.EmitParticles(10)
	
	if len(e.Particles) != 10 {
		t.Errorf("Expected 10 particles, got %d", len(e.Particles))
	}
	
	// Check that particles are created at emitter position
	for _, p := range e.Particles {
		if p.X != e.X || p.Y != e.Y {
			t.Errorf("Particle not created at emitter position: (%f, %f) vs (%f, %f)", 
				p.X, p.Y, e.X, e.Y)
		}
	}
}

func TestEmitterBurst(t *testing.T) {
	e := NewParticleEmitter(0, 0, HitSpark)
	e.Burst(15)
	
	if len(e.Particles) != 15 {
		t.Errorf("Expected 15 particles from burst, got %d", len(e.Particles))
	}
	
	if e.Active {
		t.Error("Emitter should be inactive after burst")
	}
	
	if !e.OneShot {
		t.Error("Emitter should be marked as OneShot after burst")
	}
}

func TestEmitterUpdate(t *testing.T) {
	e := NewParticleEmitter(0, 0, HitSpark)
	e.Life = 10 // Short life for testing
	e.Start()
	e.EmitParticles(5)
	
	initialCount := len(e.Particles)
	
	// Update particles multiple times until they die
	for i := 0; i < 20; i++ {
		e.Update()
	}
	
	// All particles should have died
	if len(e.Particles) != 0 {
		t.Errorf("Expected all particles to be removed after their lifetime, got %d remaining", len(e.Particles))
	}
	
	// Verify particles were alive initially
	if initialCount != 5 {
		t.Errorf("Expected 5 initial particles, got %d", initialCount)
	}
}

func TestNewParticleSystem(t *testing.T) {
	ps := NewParticleSystem(1000)
	
	if ps.maxParticles != 1000 {
		t.Errorf("Expected maxParticles 1000, got %d", ps.maxParticles)
	}
	
	if len(ps.emitters) != 0 || len(ps.particles) != 0 {
		t.Error("ParticleSystem should start empty")
	}
}

func TestParticleSystemAddEmitter(t *testing.T) {
	ps := NewParticleSystem(100)
	e := NewParticleEmitter(0, 0, HitSpark)
	
	ps.AddEmitter(e)
	
	if len(ps.emitters) != 1 {
		t.Errorf("Expected 1 emitter, got %d", len(ps.emitters))
	}
}

func TestParticleSystemAddParticle(t *testing.T) {
	ps := NewParticleSystem(100)
	p := NewParticle(0, 0, 0, 0, 60, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
	
	ps.AddParticle(p)
	
	if len(ps.particles) != 1 {
		t.Errorf("Expected 1 particle, got %d", len(ps.particles))
	}
}

func TestParticleSystemMaxParticles(t *testing.T) {
	ps := NewParticleSystem(5) // Very small limit
	
	// Try to add more particles than the limit
	for i := 0; i < 10; i++ {
		p := NewParticle(0, 0, 0, 0, 60, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
		ps.AddParticle(p)
	}
	
	if len(ps.particles) > 5 {
		t.Errorf("Expected max 5 particles, got %d", len(ps.particles))
	}
}

func TestParticleSystemUpdate(t *testing.T) {
	ps := NewParticleSystem(100)
	
	// Add a particle with short life
	p := NewParticle(0, 0, 0, 0, 5, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
	ps.AddParticle(p)
	
	// Update until particle dies
	for i := 0; i < 10; i++ {
		ps.Update()
	}
	
	if len(ps.particles) != 0 {
		t.Errorf("Expected dead particles to be removed, got %d particles", len(ps.particles))
	}
}

func TestParticleSystemGetAllParticles(t *testing.T) {
	ps := NewParticleSystem(100)
	
	// Add standalone particles
	for i := 0; i < 3; i++ {
		p := NewParticle(float64(i), 0, 0, 0, 60, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
		ps.AddParticle(p)
	}
	
	// Add emitter with particles
	e := NewParticleEmitter(0, 0, HitSpark)
	e.EmitParticles(2)
	ps.AddEmitter(e)
	
	allParticles := ps.GetAllParticles()
	
	if len(allParticles) != 5 {
		t.Errorf("Expected 5 total particles (3 standalone + 2 from emitter), got %d", len(allParticles))
	}
}

func TestParticleSystemGetParticleCount(t *testing.T) {
	ps := NewParticleSystem(100)
	
	// Add standalone particles
	for i := 0; i < 4; i++ {
		p := NewParticle(float64(i), 0, 0, 0, 60, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
		ps.AddParticle(p)
	}
	
	// Add emitter with particles
	e := NewParticleEmitter(0, 0, HitSpark)
	e.EmitParticles(3)
	ps.AddEmitter(e)
	
	count := ps.GetParticleCount()
	
	if count != 7 {
		t.Errorf("Expected 7 total particles, got %d", count)
	}
}

func TestParticleSystemClear(t *testing.T) {
	ps := NewParticleSystem(100)
	
	// Add some content
	ps.AddParticle(NewParticle(0, 0, 0, 0, 60, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark))
	e := NewParticleEmitter(0, 0, HitSpark)
	e.EmitParticles(5)
	ps.AddEmitter(e)
	
	ps.Clear()
	
	if len(ps.particles) != 0 || len(ps.emitters) != 0 {
		t.Error("ParticleSystem should be empty after Clear()")
	}
}

func TestParticleSystemRemoveOneShotEmitters(t *testing.T) {
	ps := NewParticleSystem(100)
	
	// Add a one-shot emitter with short-lived particles
	e := NewParticleEmitter(0, 0, HitSpark)
	e.Life = 2 // Very short life
	e.Burst(3)
	ps.AddEmitter(e)
	
	// Update until particles die and emitter is removed
	for i := 0; i < 10; i++ {
		ps.Update()
	}
	
	// One-shot emitter with no particles should be removed
	if len(ps.emitters) != 0 {
		t.Errorf("Expected one-shot emitter to be removed, got %d emitters", len(ps.emitters))
	}
}

func TestParticleAlphaFade(t *testing.T) {
	p := NewParticle(0, 0, 0, 0, 30, 1.0, color.RGBA{255, 255, 255, 255}, HitSpark)
	
	initialAlpha := p.Alpha
	
	// Update until particle is in fade-out phase (last 1/3 of life)
	for i := 0; i < 25; i++ {
		p.Update()
	}
	
	// Alpha should have decreased
	if p.Alpha >= initialAlpha {
		t.Errorf("Expected alpha to fade, got %d (initial: %d)", p.Alpha, initialAlpha)
	}
}
