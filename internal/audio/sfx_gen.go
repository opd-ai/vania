package audio

import (
	"math"
	"math/rand"
)

// SFXType defines sound effect categories
type SFXType int

const (
	JumpSFX SFXType = iota
	LandSFX
	AttackSFX
	HitSFX
	PickupSFX
	DoorSFX
	DamageSFX
)

// SFXGenerator generates sound effects
type SFXGenerator struct {
	Synth *Synthesizer
}

// NewSFXGenerator creates a new SFX generator
func NewSFXGenerator(sampleRate int) *SFXGenerator {
	return &SFXGenerator{
		Synth: NewSynthesizer(sampleRate),
	}
}

// Generate creates a sound effect
func (sg *SFXGenerator) Generate(sfxType SFXType, seed int64) *AudioSample {
	rng := rand.New(rand.NewSource(seed))
	
	switch sfxType {
	case JumpSFX:
		return sg.generateJump(rng)
	case LandSFX:
		return sg.generateLand(rng)
	case AttackSFX:
		return sg.generateAttack(rng)
	case HitSFX:
		return sg.generateHit(rng)
	case PickupSFX:
		return sg.generatePickup(rng)
	case DoorSFX:
		return sg.generateDoor(rng)
	case DamageSFX:
		return sg.generateDamage(rng)
	default:
		return sg.generateJump(rng)
	}
}

// generateJump creates jump sound (rising pitch sweep)
func (sg *SFXGenerator) generateJump(rng *rand.Rand) *AudioSample {
	startFreq := 100.0 + rng.Float64()*50.0
	endFreq := 400.0 + rng.Float64()*100.0
	duration := 0.1 + rng.Float64()*0.05
	
	sample := sg.Synth.FrequencySweep(SquareWave, startFreq, endFreq, duration)
	
	envelope := ADSR{
		Attack:  0.01,
		Decay:   0.02,
		Sustain: 0.6,
		Release: 0.05,
	}
	
	return sg.Synth.ApplyEnvelope(sample, envelope)
}

// generateLand creates landing sound (percussive burst)
func (sg *SFXGenerator) generateLand(rng *rand.Rand) *AudioSample {
	duration := 0.1 + rng.Float64()*0.05
	
	sample := sg.Synth.GenerateWave(NoiseWave, 0, duration)
	sample = sg.Synth.ApplyLowPassFilter(sample, 200.0+rng.Float64()*100.0)
	
	envelope := ADSR{
		Attack:  0.005,
		Decay:   0.02,
		Sustain: 0.3,
		Release: 0.05,
	}
	
	return sg.Synth.ApplyEnvelope(sample, envelope)
}

// generateAttack creates attack sound
func (sg *SFXGenerator) generateAttack(rng *rand.Rand) *AudioSample {
	freq := 200.0 + rng.Float64()*100.0
	duration := 0.15 + rng.Float64()*0.05
	
	sample := sg.Synth.GenerateWave(SquareWave, freq, duration)
	
	// Add pitch drop
	sweep := sg.Synth.FrequencySweep(SquareWave, freq, freq*0.7, duration)
	
	mixed := sg.Synth.Mix([]*AudioSample{sample, sweep}, []float64{0.5, 0.5})
	
	envelope := ADSR{
		Attack:  0.01,
		Decay:   0.03,
		Sustain: 0.5,
		Release: 0.08,
	}
	
	return sg.Synth.ApplyEnvelope(mixed, envelope)
}

// generateHit creates hit sound
func (sg *SFXGenerator) generateHit(rng *rand.Rand) *AudioSample {
	freq := 150.0 + rng.Float64()*50.0
	duration := 0.1 + rng.Float64()*0.05
	
	// Mix tone and noise
	tone := sg.Synth.GenerateWave(SineWave, freq, duration)
	noise := sg.Synth.GenerateWave(NoiseWave, 0, duration)
	noise = sg.Synth.ApplyLowPassFilter(noise, 300.0)
	
	mixed := sg.Synth.Mix([]*AudioSample{tone, noise}, []float64{0.6, 0.4})
	
	envelope := ADSR{
		Attack:  0.005,
		Decay:   0.02,
		Sustain: 0.3,
		Release: 0.05,
	}
	
	return sg.Synth.ApplyEnvelope(mixed, envelope)
}

// generatePickup creates pickup sound (ascending arpeggio)
func (sg *SFXGenerator) generatePickup(rng *rand.Rand) *AudioSample {
	baseFreq := 400.0 + rng.Float64()*100.0
	noteDuration := 0.08
	
	// Major chord arpeggio
	notes := []float64{
		baseFreq,           // Root
		baseFreq * 1.25,    // Major third
		baseFreq * 1.5,     // Perfect fifth
	}
	
	samples := make([]*AudioSample, len(notes))
	for i, freq := range notes {
		sample := sg.Synth.GenerateWave(SineWave, freq, noteDuration)
		
		envelope := ADSR{
			Attack:  0.01,
			Decay:   0.02,
			Sustain: 0.7,
			Release: 0.03,
		}
		
		samples[i] = sg.Synth.ApplyEnvelope(sample, envelope)
	}
	
	// Concatenate notes
	totalLen := 0
	for _, s := range samples {
		totalLen += len(s.Data)
	}
	
	result := make([]float64, totalLen)
	offset := 0
	for _, s := range samples {
		copy(result[offset:], s.Data)
		offset += len(s.Data)
	}
	
	return &AudioSample{
		Data:       result,
		SampleRate: sg.Synth.SampleRate,
		Duration:   float64(totalLen) / float64(sg.Synth.SampleRate),
	}
}

// generateDoor creates door sound (mechanical sweep)
func (sg *SFXGenerator) generateDoor(rng *rand.Rand) *AudioSample {
	startFreq := 80.0 + rng.Float64()*20.0
	endFreq := 120.0 + rng.Float64()*30.0
	duration := 0.3 + rng.Float64()*0.1
	
	sample := sg.Synth.FrequencySweep(SawtoothWave, startFreq, endFreq, duration)
	
	envelope := ADSR{
		Attack:  0.05,
		Decay:   0.1,
		Sustain: 0.6,
		Release: 0.1,
	}
	
	return sg.Synth.ApplyEnvelope(sample, envelope)
}

// generateDamage creates damage sound (harsh dissonance)
func (sg *SFXGenerator) generateDamage(rng *rand.Rand) *AudioSample {
	freq1 := 200.0 + rng.Float64()*50.0
	freq2 := freq1 * 1.06 // Slightly dissonant
	duration := 0.2 + rng.Float64()*0.1
	
	sample1 := sg.Synth.GenerateWave(SquareWave, freq1, duration)
	sample2 := sg.Synth.GenerateWave(SquareWave, freq2, duration)
	noise := sg.Synth.GenerateWave(NoiseWave, 0, duration)
	
	mixed := sg.Synth.Mix([]*AudioSample{sample1, sample2, noise}, []float64{0.4, 0.4, 0.2})
	
	envelope := ADSR{
		Attack:  0.005,
		Decay:   0.05,
		Sustain: 0.4,
		Release: 0.1,
	}
	
	return sg.Synth.ApplyEnvelope(mixed, envelope)
}

// GenerateExplosion creates explosion sound
func (sg *SFXGenerator) GenerateExplosion(seed int64) *AudioSample {
	rng := rand.New(rand.NewSource(seed))
	duration := 0.5 + rng.Float64()*0.2
	
	// Low frequency rumble
	rumble := sg.Synth.GenerateWave(SineWave, 40.0+rng.Float64()*20.0, duration)
	
	// Noise burst
	noise := sg.Synth.GenerateWave(NoiseWave, 0, duration)
	noise = sg.Synth.ApplyLowPassFilter(noise, 500.0)
	
	mixed := sg.Synth.Mix([]*AudioSample{rumble, noise}, []float64{0.5, 0.5})
	
	envelope := ADSR{
		Attack:  0.01,
		Decay:   0.1,
		Sustain: 0.3,
		Release: 0.3,
	}
	
	return sg.Synth.ApplyEnvelope(mixed, envelope)
}

// ApplyDistortion adds distortion effect
func (sg *SFXGenerator) ApplyDistortion(sample *AudioSample, amount float64) *AudioSample {
	result := make([]float64, len(sample.Data))
	
	for i, val := range sample.Data {
		// Soft clipping distortion
		result[i] = math.Tanh(val * amount)
	}
	
	return &AudioSample{
		Data:       result,
		SampleRate: sample.SampleRate,
		Duration:   sample.Duration,
	}
}
