package particle

import (
	"image/color"
	"math"
	"testing"
)

// TestParticlePresets_CreateHitEffect tests hit effect creation
func TestParticlePresets_CreateHitEffect(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 100.0, 200.0
	direction := 1.0

	emitter := pp.CreateHitEffect(x, y, direction)

	if emitter == nil {
		t.Fatal("CreateHitEffect returned nil")
	}

	if emitter.X != x || emitter.Y != y {
		t.Errorf("Expected position (%.1f, %.1f), got (%.1f, %.1f)", x, y, emitter.X, emitter.Y)
	}

	if emitter.Type != HitSpark {
		t.Errorf("Expected type HitSpark, got %d", emitter.Type)
	}

	if !emitter.OneShot {
		t.Error("Hit effect should be one-shot")
	}

	if emitter.EmitRate != 20 {
		t.Errorf("Expected EmitRate 20, got %d", emitter.EmitRate)
	}
}

// TestParticlePresets_CreateHitEffect_Direction tests hit effect with different directions
func TestParticlePresets_CreateHitEffect_Direction(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 100.0, 200.0

	tests := []struct {
		name      string
		direction float64
	}{
		{"positive direction", 1.0},
		{"negative direction", -1.0},
		{"zero direction", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emitter := pp.CreateHitEffect(x, y, tt.direction)
			if emitter == nil {
				t.Error("CreateHitEffect returned nil")
			}
			// Negative direction should invert spread
			if tt.direction < 0 && emitter.Spread >= 0 {
				t.Error("Expected negative spread for negative direction")
			}
		})
	}
}

// TestParticlePresets_CreateDashTrail tests dash trail creation
func TestParticlePresets_CreateDashTrail(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 150.0, 250.0

	emitter := pp.CreateDashTrail(x, y)

	if emitter == nil {
		t.Fatal("CreateDashTrail returned nil")
	}

	if emitter.Type != DashTrail {
		t.Errorf("Expected type DashTrail, got %d", emitter.Type)
	}

	if emitter.OneShot {
		t.Error("Dash trail should be continuous")
	}

	if emitter.Spread != math.Pi {
		t.Errorf("Expected spread PI, got %f", emitter.Spread)
	}
}

// TestParticlePresets_CreateJumpDust tests jump dust creation
func TestParticlePresets_CreateJumpDust(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 120.0, 180.0

	emitter := pp.CreateJumpDust(x, y)

	if emitter == nil {
		t.Fatal("CreateJumpDust returned nil")
	}

	if emitter.Type != JumpDust {
		t.Errorf("Expected type JumpDust, got %d", emitter.Type)
	}

	if !emitter.OneShot {
		t.Error("Jump dust should be one-shot")
	}

	if emitter.EmitRate != 15 {
		t.Errorf("Expected EmitRate 15, got %d", emitter.EmitRate)
	}
}

// TestParticlePresets_CreateLandDust tests land dust creation
func TestParticlePresets_CreateLandDust(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 130.0, 190.0

	emitter := pp.CreateLandDust(x, y)

	if emitter == nil {
		t.Fatal("CreateLandDust returned nil")
	}

	if emitter.Type != LandDust {
		t.Errorf("Expected type LandDust, got %d", emitter.Type)
	}

	if !emitter.OneShot {
		t.Error("Land dust should be one-shot")
	}

	if emitter.EmitRate != 20 {
		t.Errorf("Expected EmitRate 20, got %d", emitter.EmitRate)
	}

	if emitter.Spread != math.Pi {
		t.Errorf("Expected spread PI, got %f", emitter.Spread)
	}
}

// TestParticlePresets_CreateWalkDust tests walk dust creation
func TestParticlePresets_CreateWalkDust(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 110.0, 170.0

	emitter := pp.CreateWalkDust(x, y)

	if emitter == nil {
		t.Fatal("CreateWalkDust returned nil")
	}

	if emitter.Type != WalkDust {
		t.Errorf("Expected type WalkDust, got %d", emitter.Type)
	}

	if !emitter.OneShot {
		t.Error("Walk dust should be one-shot")
	}

	if emitter.EmitRate != 5 {
		t.Errorf("Expected EmitRate 5, got %d", emitter.EmitRate)
	}
}

// TestParticlePresets_CreateBloodSplatter tests blood splatter creation
func TestParticlePresets_CreateBloodSplatter(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 140.0, 200.0
	direction := 1.0

	emitter := pp.CreateBloodSplatter(x, y, direction)

	if emitter == nil {
		t.Fatal("CreateBloodSplatter returned nil")
	}

	if emitter.Type != BloodSplatter {
		t.Errorf("Expected type BloodSplatter, got %d", emitter.Type)
	}

	if !emitter.OneShot {
		t.Error("Blood splatter should be one-shot")
	}

	// Check for red color
	if emitter.Color.R < 150 {
		t.Error("Blood splatter should have high red component")
	}
}

// TestParticlePresets_CreateExplosion tests explosion creation
func TestParticlePresets_CreateExplosion(t *testing.T) {
	tests := []struct {
		name string
		size float64
	}{
		{"small explosion", 0.5},
		{"normal explosion", 1.0},
		{"large explosion", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp := &ParticlePresets{}
			x, y := 160.0, 220.0

			emitter := pp.CreateExplosion(x, y, tt.size)

			if emitter == nil {
				t.Fatal("CreateExplosion returned nil")
			}

			if emitter.Type != Explosion {
				t.Errorf("Expected type Explosion, got %d", emitter.Type)
			}

			if !emitter.OneShot {
				t.Error("Explosion should be one-shot")
			}

			// Verify size scaling
			expectedSpeed := 5.0 * tt.size
			if emitter.Speed != expectedSpeed {
				t.Errorf("Expected speed %f, got %f", expectedSpeed, emitter.Speed)
			}

			if emitter.Spread != math.Pi*2 {
				t.Errorf("Expected 360 degree spread, got %f", emitter.Spread)
			}
		})
	}
}

// TestParticlePresets_CreateSmoke tests smoke creation
func TestParticlePresets_CreateSmoke(t *testing.T) {
	tests := []struct {
		name       string
		continuous bool
	}{
		{"one-shot smoke", false},
		{"continuous smoke", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp := &ParticlePresets{}
			x, y := 170.0, 230.0

			emitter := pp.CreateSmoke(x, y, tt.continuous)

			if emitter == nil {
				t.Fatal("CreateSmoke returned nil")
			}

			if emitter.Type != Smoke {
				t.Errorf("Expected type Smoke, got %d", emitter.Type)
			}

			if emitter.OneShot == tt.continuous {
				t.Errorf("Expected OneShot=%v for continuous=%v", !tt.continuous, tt.continuous)
			}

			// Smoke should rise
			if emitter.Gravity >= 0 {
				t.Error("Smoke should have negative gravity (rise upward)")
			}
		})
	}
}

// TestParticlePresets_CreateRain tests rain creation
func TestParticlePresets_CreateRain(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 180.0, 240.0

	emitter := pp.CreateRain(x, y)

	if emitter == nil {
		t.Fatal("CreateRain returned nil")
	}

	if emitter.Type != Rain {
		t.Errorf("Expected type Rain, got %d", emitter.Type)
	}

	if emitter.OneShot {
		t.Error("Rain should be continuous")
	}

	// Rain should be nearly vertical
	if emitter.Spread > 0.2 {
		t.Errorf("Expected small spread for rain, got %f", emitter.Spread)
	}

	// Rain should fall fast
	if emitter.Speed < 5.0 {
		t.Error("Rain should have high speed")
	}
}

// TestParticlePresets_CreateSnow tests snow creation
func TestParticlePresets_CreateSnow(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 190.0, 250.0

	emitter := pp.CreateSnow(x, y)

	if emitter == nil {
		t.Fatal("CreateSnow returned nil")
	}

	if emitter.Type != Snow {
		t.Errorf("Expected type Snow, got %d", emitter.Type)
	}

	if emitter.OneShot {
		t.Error("Snow should be continuous")
	}

	// Snow should fall slowly
	if emitter.Speed > 2.0 {
		t.Error("Snow should have slow speed")
	}

	// Snow should have very light gravity
	if emitter.Gravity > 0.05 {
		t.Error("Snow should have very light gravity")
	}
}

// TestParticlePresets_CreateEmbers tests ember creation
func TestParticlePresets_CreateEmbers(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 200.0, 260.0

	emitter := pp.CreateEmbers(x, y)

	if emitter == nil {
		t.Fatal("CreateEmbers returned nil")
	}

	if emitter.Type != Embers {
		t.Errorf("Expected type Embers, got %d", emitter.Type)
	}

	if emitter.OneShot {
		t.Error("Embers should be continuous")
	}

	// Embers should rise
	if emitter.Gravity >= 0 {
		t.Error("Embers should have negative gravity (rise upward)")
	}
}

// TestParticlePresets_CreateSparkles tests sparkle creation
func TestParticlePresets_CreateSparkles(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 210.0, 270.0

	emitter := pp.CreateSparkles(x, y)

	if emitter == nil {
		t.Fatal("CreateSparkles returned nil")
	}

	if emitter.Type != Sparkles {
		t.Errorf("Expected type Sparkles, got %d", emitter.Type)
	}

	if emitter.OneShot {
		t.Error("Sparkles should be continuous")
	}

	// Sparkles should spread in all directions
	if emitter.Spread != math.Pi*2 {
		t.Errorf("Expected 360 degree spread, got %f", emitter.Spread)
	}

	// Sparkles should float (no gravity)
	if emitter.Gravity != 0.0 {
		t.Errorf("Expected zero gravity, got %f", emitter.Gravity)
	}
}

// TestParticlePresets_CreateBubbles tests bubble creation
func TestParticlePresets_CreateBubbles(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 220.0, 280.0

	emitter := pp.CreateBubbles(x, y)

	if emitter == nil {
		t.Fatal("CreateBubbles returned nil")
	}

	if emitter.Type != Bubbles {
		t.Errorf("Expected type Bubbles, got %d", emitter.Type)
	}

	if emitter.OneShot {
		t.Error("Bubbles should be continuous")
	}

	// Bubbles should rise (negative gravity for buoyancy)
	if emitter.Gravity >= 0 {
		t.Error("Bubbles should have negative gravity (buoyancy)")
	}
}

// TestParticlePresets_CreateLightning tests lightning creation
func TestParticlePresets_CreateLightning(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 230.0, 290.0

	emitter := pp.CreateLightning(x, y)

	if emitter == nil {
		t.Fatal("CreateLightning returned nil")
	}

	if emitter.Type != Lightning {
		t.Errorf("Expected type Lightning, got %d", emitter.Type)
	}

	if !emitter.OneShot {
		t.Error("Lightning should be one-shot")
	}

	// Lightning should have high emit rate
	if emitter.EmitRate < 30 {
		t.Error("Lightning should have high emit rate")
	}

	// Lightning should have very short life
	if emitter.Life > 10 {
		t.Error("Lightning should have very short life")
	}

	// Lightning should have no gravity
	if emitter.Gravity != 0.0 {
		t.Errorf("Expected zero gravity, got %f", emitter.Gravity)
	}
}

// TestParticlePresets_CreateDamageNumber tests damage number creation
func TestParticlePresets_CreateDamageNumber(t *testing.T) {
	pp := &ParticlePresets{}
	x, y := 240.0, 300.0
	damage := 42

	particle := pp.CreateDamageNumber(x, y, damage)

	if particle == nil {
		t.Fatal("CreateDamageNumber returned nil")
	}

	if particle.X != x || particle.Y != y {
		t.Errorf("Expected position (%.1f, %.1f), got (%.1f, %.1f)", x, y, particle.X, particle.Y)
	}

	if particle.Type != DamageNumber {
		t.Errorf("Expected type DamageNumber, got %d", particle.Type)
	}

	// Damage value should be stored
	if particle.Data != damage {
		t.Errorf("Expected damage %d, got %v", damage, particle.Data)
	}

	// Should float upward
	if particle.VelY >= 0 {
		t.Error("Damage number should float upward (negative VelY)")
	}

	// Should have no gravity
	if particle.AccelY != 0.0 {
		t.Errorf("Expected zero AccelY, got %f", particle.AccelY)
	}
}

// TestParticlePresets_CreateDamageNumber_VariousDamages tests different damage values
func TestParticlePresets_CreateDamageNumber_VariousDamages(t *testing.T) {
	tests := []struct {
		name   string
		damage int
	}{
		{"zero damage", 0},
		{"small damage", 5},
		{"medium damage", 50},
		{"large damage", 500},
		{"negative damage", -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp := &ParticlePresets{}
			particle := pp.CreateDamageNumber(100.0, 100.0, tt.damage)

			if particle.Data != tt.damage {
				t.Errorf("Expected damage %d, got %v", tt.damage, particle.Data)
			}
		})
	}
}

// TestParticlePresets_AllEffectsHaveValidColors tests all effects have colors set
func TestParticlePresets_AllEffectsHaveValidColors(t *testing.T) {
	pp := &ParticlePresets{}
	
	emitters := []*ParticleEmitter{
		pp.CreateHitEffect(0, 0, 1),
		pp.CreateDashTrail(0, 0),
		pp.CreateJumpDust(0, 0),
		pp.CreateLandDust(0, 0),
		pp.CreateWalkDust(0, 0),
		pp.CreateBloodSplatter(0, 0, 1),
		pp.CreateExplosion(0, 0, 1),
		pp.CreateSmoke(0, 0, false),
		pp.CreateRain(0, 0),
		pp.CreateSnow(0, 0),
		pp.CreateEmbers(0, 0),
		pp.CreateSparkles(0, 0),
		pp.CreateBubbles(0, 0),
		pp.CreateLightning(0, 0),
	}

	for i, emitter := range emitters {
		// Check that color is not default zero value
		if emitter.Color == (color.RGBA{}) {
			t.Errorf("Emitter %d has zero color value", i)
		}
	}
}

// TestParticlePresets_AllEffectsHaveValidSizes tests all effects have reasonable sizes
func TestParticlePresets_AllEffectsHaveValidSizes(t *testing.T) {
	pp := &ParticlePresets{}
	
	emitters := []*ParticleEmitter{
		pp.CreateHitEffect(0, 0, 1),
		pp.CreateDashTrail(0, 0),
		pp.CreateJumpDust(0, 0),
		pp.CreateLandDust(0, 0),
		pp.CreateWalkDust(0, 0),
		pp.CreateBloodSplatter(0, 0, 1),
		pp.CreateExplosion(0, 0, 1),
		pp.CreateSmoke(0, 0, false),
		pp.CreateRain(0, 0),
		pp.CreateSnow(0, 0),
		pp.CreateEmbers(0, 0),
		pp.CreateSparkles(0, 0),
		pp.CreateBubbles(0, 0),
		pp.CreateLightning(0, 0),
	}

	for i, emitter := range emitters {
		if emitter.Size <= 0 {
			t.Errorf("Emitter %d has invalid size: %f", i, emitter.Size)
		}
	}
}

// TestParticlePresets_AllEffectsHaveValidLife tests all effects have positive life
func TestParticlePresets_AllEffectsHaveValidLife(t *testing.T) {
	pp := &ParticlePresets{}
	
	emitters := []*ParticleEmitter{
		pp.CreateHitEffect(0, 0, 1),
		pp.CreateDashTrail(0, 0),
		pp.CreateJumpDust(0, 0),
		pp.CreateLandDust(0, 0),
		pp.CreateWalkDust(0, 0),
		pp.CreateBloodSplatter(0, 0, 1),
		pp.CreateExplosion(0, 0, 1),
		pp.CreateSmoke(0, 0, false),
		pp.CreateRain(0, 0),
		pp.CreateSnow(0, 0),
		pp.CreateEmbers(0, 0),
		pp.CreateSparkles(0, 0),
		pp.CreateBubbles(0, 0),
		pp.CreateLightning(0, 0),
	}

	for i, emitter := range emitters {
		if emitter.Life <= 0 {
			t.Errorf("Emitter %d has invalid life: %d", i, emitter.Life)
		}
	}
}

// TestParticlePresets_PositionPreservation tests that positions are preserved
func TestParticlePresets_PositionPreservation(t *testing.T) {
	pp := &ParticlePresets{}
	testCases := []struct {
		name string
		x, y float64
		fn   func(float64, float64) *ParticleEmitter
	}{
		{"hit effect", 123.45, 678.90, func(x, y float64) *ParticleEmitter { return pp.CreateHitEffect(x, y, 1) }},
		{"dash trail", 111.11, 222.22, pp.CreateDashTrail},
		{"jump dust", 333.33, 444.44, pp.CreateJumpDust},
		{"rain", 555.55, 666.66, pp.CreateRain},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emitter := tc.fn(tc.x, tc.y)
			if emitter.X != tc.x || emitter.Y != tc.y {
				t.Errorf("Expected position (%.2f, %.2f), got (%.2f, %.2f)", 
					tc.x, tc.y, emitter.X, emitter.Y)
			}
		})
	}
}

// TestParticlePresets_OneShotVsContinuous tests one-shot vs continuous classification
func TestParticlePresets_OneShotVsContinuous(t *testing.T) {
	pp := &ParticlePresets{}
	
	oneShotEffects := []struct {
		name    string
		emitter *ParticleEmitter
	}{
		{"hit effect", pp.CreateHitEffect(0, 0, 1)},
		{"jump dust", pp.CreateJumpDust(0, 0)},
		{"land dust", pp.CreateLandDust(0, 0)},
		{"walk dust", pp.CreateWalkDust(0, 0)},
		{"blood splatter", pp.CreateBloodSplatter(0, 0, 1)},
		{"explosion", pp.CreateExplosion(0, 0, 1)},
		{"lightning", pp.CreateLightning(0, 0)},
	}

	for _, effect := range oneShotEffects {
		if !effect.emitter.OneShot {
			t.Errorf("%s should be one-shot", effect.name)
		}
	}

	continuousEffects := []struct {
		name    string
		emitter *ParticleEmitter
	}{
		{"dash trail", pp.CreateDashTrail(0, 0)},
		{"rain", pp.CreateRain(0, 0)},
		{"snow", pp.CreateSnow(0, 0)},
		{"embers", pp.CreateEmbers(0, 0)},
		{"sparkles", pp.CreateSparkles(0, 0)},
		{"bubbles", pp.CreateBubbles(0, 0)},
	}

	for _, effect := range continuousEffects {
		if effect.emitter.OneShot {
			t.Errorf("%s should be continuous", effect.name)
		}
	}
}

// TestParticlePresets_ExplosionScaling tests explosion size scaling
func TestParticlePresets_ExplosionScaling(t *testing.T) {
	pp := &ParticlePresets{}
	
	size1 := 1.0
	size2 := 2.0
	
	exp1 := pp.CreateExplosion(0, 0, size1)
	exp2 := pp.CreateExplosion(0, 0, size2)
	
	// Speed should scale
	if exp2.Speed <= exp1.Speed {
		t.Error("Larger explosion should have higher speed")
	}
	
	// Size should scale
	if exp2.Size <= exp1.Size {
		t.Error("Larger explosion should have bigger particle size")
	}
	
	// Life should scale
	if exp2.Life <= exp1.Life {
		t.Error("Larger explosion should have longer life")
	}
}

// TestParticlePresets_SmokeMode tests smoke continuous mode
func TestParticlePresets_SmokeMode(t *testing.T) {
	pp := &ParticlePresets{}
	
	oneShotSmoke := pp.CreateSmoke(0, 0, false)
	continuousSmoke := pp.CreateSmoke(0, 0, true)
	
	if !oneShotSmoke.OneShot {
		t.Error("One-shot smoke should have OneShot=true")
	}
	
	if continuousSmoke.OneShot {
		t.Error("Continuous smoke should have OneShot=false")
	}
}

// TestParticlePresets_DamageNumberProperties tests damage number particle properties
func TestParticlePresets_DamageNumberProperties(t *testing.T) {
	pp := &ParticlePresets{}
	particle := pp.CreateDamageNumber(50.0, 75.0, 123)
	
	// Should have positive life
	if particle.Life <= 0 {
		t.Error("Damage number should have positive life")
	}
	
	// Should have positive size
	if particle.Size <= 0 {
		t.Error("Damage number should have positive size")
	}
	
	// Should have white color (for visibility)
	if particle.Color.R != 255 || particle.Color.G != 255 || particle.Color.B != 255 {
		t.Error("Damage number should have white color")
	}
}
