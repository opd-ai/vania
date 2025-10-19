// Package particle provides preset particle effect generators for common
// game events like combat hits, movement effects, and environmental particles.
package particle

import (
	"image/color"
	"math"
)

// ParticlePresets provides factory methods for common particle effects
type ParticlePresets struct{}

// CreateHitEffect creates a hit spark effect at the specified position
func (pp *ParticlePresets) CreateHitEffect(x, y float64, direction float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, HitSpark)
	emitter.EmitRate = 20
	emitter.Spread = math.Pi / 3 // 60 degrees
	emitter.Speed = 4.0
	emitter.SpeedVariance = 2.0
	emitter.Life = 15 // 0.25 seconds
	emitter.LifeVariance = 5
	emitter.Size = 3.0
	emitter.SizeVariance = 1.0
	emitter.Gravity = 0.1
	emitter.Color = color.RGBA{255, 200, 100, 255} // Orange-yellow
	emitter.OneShot = true
	
	// Adjust direction based on hit direction
	if direction < 0 {
		emitter.Spread = -emitter.Spread
	}
	
	return emitter
}

// CreateDashTrail creates a dash trail effect
func (pp *ParticlePresets) CreateDashTrail(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, DashTrail)
	emitter.EmitRate = 10
	emitter.Spread = math.Pi // 180 degrees
	emitter.Speed = 1.0
	emitter.SpeedVariance = 0.5
	emitter.Life = 20 // 0.33 seconds
	emitter.LifeVariance = 5
	emitter.Size = 4.0
	emitter.SizeVariance = 1.0
	emitter.Gravity = 0.05
	emitter.Color = color.RGBA{100, 150, 255, 200} // Light blue, semi-transparent
	emitter.OneShot = false // Continuous while dashing
	
	return emitter
}

// CreateJumpDust creates dust particles when jumping
func (pp *ParticlePresets) CreateJumpDust(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, JumpDust)
	emitter.EmitRate = 15
	emitter.Spread = math.Pi / 4 // 45 degrees up
	emitter.Speed = 2.0
	emitter.SpeedVariance = 1.0
	emitter.Life = 25 // 0.4 seconds
	emitter.LifeVariance = 8
	emitter.Size = 3.0
	emitter.SizeVariance = 1.0
	emitter.Gravity = 0.15
	emitter.Color = color.RGBA{150, 130, 100, 180} // Brownish dust
	emitter.OneShot = true
	
	return emitter
}

// CreateLandDust creates dust particles when landing
func (pp *ParticlePresets) CreateLandDust(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, LandDust)
	emitter.EmitRate = 20
	emitter.Spread = math.Pi // Spread outward
	emitter.Speed = 3.0
	emitter.SpeedVariance = 1.5
	emitter.Life = 20 // 0.33 seconds
	emitter.LifeVariance = 8
	emitter.Size = 4.0
	emitter.SizeVariance = 1.5
	emitter.Gravity = 0.2
	emitter.Color = color.RGBA{140, 120, 90, 200} // Brownish dust
	emitter.OneShot = true
	
	return emitter
}

// CreateWalkDust creates subtle dust for walking
func (pp *ParticlePresets) CreateWalkDust(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, WalkDust)
	emitter.EmitRate = 5
	emitter.Spread = math.Pi / 6 // 30 degrees
	emitter.Speed = 0.5
	emitter.SpeedVariance = 0.3
	emitter.Life = 15 // 0.25 seconds
	emitter.LifeVariance = 5
	emitter.Size = 2.0
	emitter.SizeVariance = 0.5
	emitter.Gravity = 0.05
	emitter.Color = color.RGBA{120, 110, 90, 100} // Light dust, very transparent
	emitter.OneShot = true
	
	return emitter
}

// CreateBloodSplatter creates blood particles when enemy is hit
func (pp *ParticlePresets) CreateBloodSplatter(x, y float64, direction float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, BloodSplatter)
	emitter.EmitRate = 15
	emitter.Spread = math.Pi / 2 // 90 degrees
	emitter.Speed = 3.5
	emitter.SpeedVariance = 1.5
	emitter.Life = 30 // 0.5 seconds
	emitter.LifeVariance = 10
	emitter.Size = 2.5
	emitter.SizeVariance = 1.0
	emitter.Gravity = 0.3
	emitter.Color = color.RGBA{200, 50, 50, 255} // Red blood
	emitter.OneShot = true
	
	return emitter
}

// CreateExplosion creates an explosion effect
func (pp *ParticlePresets) CreateExplosion(x, y float64, size float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Explosion)
	emitter.EmitRate = 30
	emitter.Spread = math.Pi * 2 // 360 degrees
	emitter.Speed = 5.0 * size
	emitter.SpeedVariance = 2.0 * size
	emitter.Life = int(40 * size) // Scales with size
	emitter.LifeVariance = int(15 * size)
	emitter.Size = 5.0 * size
	emitter.SizeVariance = 2.0 * size
	emitter.Gravity = 0.1
	emitter.Color = color.RGBA{255, 150, 50, 255} // Orange fire
	emitter.OneShot = true
	
	return emitter
}

// CreateSmoke creates smoke particles
func (pp *ParticlePresets) CreateSmoke(x, y float64, continuous bool) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Smoke)
	emitter.EmitRate = 3
	emitter.Spread = math.Pi / 4 // 45 degrees upward
	emitter.Speed = 1.0
	emitter.SpeedVariance = 0.5
	emitter.Life = 60 // 1 second
	emitter.LifeVariance = 20
	emitter.Size = 6.0
	emitter.SizeVariance = 2.0
	emitter.Gravity = -0.05 // Rise upward
	emitter.Color = color.RGBA{100, 100, 100, 150} // Gray smoke, semi-transparent
	emitter.OneShot = !continuous
	
	return emitter
}

// CreateRain creates rain particles
func (pp *ParticlePresets) CreateRain(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Rain)
	emitter.EmitRate = 10
	emitter.Spread = 0.1 // Nearly vertical
	emitter.Speed = 8.0
	emitter.SpeedVariance = 2.0
	emitter.Life = 120 // 2 seconds
	emitter.LifeVariance = 30
	emitter.Size = 2.0
	emitter.SizeVariance = 0.5
	emitter.Gravity = 0.3
	emitter.Color = color.RGBA{150, 180, 220, 180} // Light blue water
	emitter.OneShot = false // Continuous
	
	return emitter
}

// CreateSnow creates snow particles
func (pp *ParticlePresets) CreateSnow(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Snow)
	emitter.EmitRate = 5
	emitter.Spread = 0.3 // Slight spread
	emitter.Speed = 1.0
	emitter.SpeedVariance = 0.5
	emitter.Life = 180 // 3 seconds - slow falling
	emitter.LifeVariance = 60
	emitter.Size = 3.0
	emitter.SizeVariance = 1.0
	emitter.Gravity = 0.02 // Very light
	emitter.Color = color.RGBA{240, 240, 255, 220} // White snow
	emitter.OneShot = false // Continuous
	
	return emitter
}

// CreateEmbers creates ember particles
func (pp *ParticlePresets) CreateEmbers(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Embers)
	emitter.EmitRate = 3
	emitter.Spread = math.Pi / 3 // 60 degrees upward
	emitter.Speed = 1.5
	emitter.SpeedVariance = 0.8
	emitter.Life = 90 // 1.5 seconds
	emitter.LifeVariance = 30
	emitter.Size = 2.5
	emitter.SizeVariance = 1.0
	emitter.Gravity = -0.1 // Rise upward
	emitter.Color = color.RGBA{255, 100, 50, 200} // Orange-red embers
	emitter.OneShot = false // Continuous
	
	return emitter
}

// CreateSparkles creates sparkle particles
func (pp *ParticlePresets) CreateSparkles(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Sparkles)
	emitter.EmitRate = 5
	emitter.Spread = math.Pi * 2 // 360 degrees
	emitter.Speed = 0.5
	emitter.SpeedVariance = 0.3
	emitter.Life = 40 // 0.67 seconds
	emitter.LifeVariance = 15
	emitter.Size = 2.0
	emitter.SizeVariance = 0.5
	emitter.Gravity = 0.0 // Float
	emitter.Color = color.RGBA{255, 255, 100, 255} // Bright yellow
	emitter.OneShot = false // Continuous
	
	return emitter
}

// CreateBubbles creates bubble particles
func (pp *ParticlePresets) CreateBubbles(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Bubbles)
	emitter.EmitRate = 2
	emitter.Spread = math.Pi / 6 // 30 degrees upward
	emitter.Speed = 1.0
	emitter.SpeedVariance = 0.5
	emitter.Life = 90 // 1.5 seconds
	emitter.LifeVariance = 30
	emitter.Size = 4.0
	emitter.SizeVariance = 2.0
	emitter.Gravity = -0.08 // Rise upward (buoyancy)
	emitter.Color = color.RGBA{150, 200, 255, 100} // Light blue bubbles, transparent
	emitter.OneShot = false // Continuous
	
	return emitter
}

// CreateLightning creates lightning flash particles
func (pp *ParticlePresets) CreateLightning(x, y float64) *ParticleEmitter {
	emitter := NewParticleEmitter(x, y, Lightning)
	emitter.EmitRate = 50
	emitter.Spread = math.Pi * 2 // 360 degrees
	emitter.Speed = 6.0
	emitter.SpeedVariance = 3.0
	emitter.Life = 8 // Very brief flash
	emitter.LifeVariance = 3
	emitter.Size = 4.0
	emitter.SizeVariance = 2.0
	emitter.Gravity = 0.0
	emitter.Color = color.RGBA{200, 220, 255, 255} // Electric blue-white
	emitter.OneShot = true
	
	return emitter
}

// CreateDamageNumber creates a floating damage number particle
func (pp *ParticlePresets) CreateDamageNumber(x, y float64, damage int) *Particle {
	particle := NewParticle(x, y, 0, -1.5, 60, 1.0, color.RGBA{255, 255, 255, 255}, DamageNumber)
	particle.Data = damage // Store damage value
	particle.AccelY = 0.0 // No gravity, just float up
	return particle
}
