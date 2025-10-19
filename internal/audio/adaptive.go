package audio

import (
	"math"
	"math/rand"
)

// MusicIntensity represents the current music intensity level
type MusicIntensity int

const (
	// IntensityCalm represents peaceful exploration
	IntensityCalm MusicIntensity = iota
	// IntensityTension represents light danger or nearby enemies
	IntensityTension
	// IntensityCombat represents active combat
	IntensityCombat
	// IntensityBoss represents boss battle
	IntensityBoss
)

// MusicLayer represents a single musical layer (drums, bass, melody, etc.)
type MusicLayer struct {
	Name       string
	Audio      *AudioSample
	BaseVolume float64 // Base volume (0-1)
	MinIntensity MusicIntensity // Minimum intensity for this layer
}

// AdaptiveMusicTrack contains multiple layers that can be mixed dynamically
type AdaptiveMusicTrack struct {
	Layers           []*MusicLayer
	CurrentIntensity MusicIntensity
	TargetIntensity  MusicIntensity
	TransitionSpeed  float64 // How fast to transition (0-1 per update)
	CurrentMix       map[string]float64 // Current volume per layer
}

// NewAdaptiveMusicTrack creates a new adaptive music track
func NewAdaptiveMusicTrack() *AdaptiveMusicTrack {
	return &AdaptiveMusicTrack{
		Layers:           make([]*MusicLayer, 0),
		CurrentIntensity: IntensityCalm,
		TargetIntensity:  IntensityCalm,
		TransitionSpeed:  0.05, // Smooth transitions
		CurrentMix:       make(map[string]float64),
	}
}

// AddLayer adds a music layer to the track
func (amt *AdaptiveMusicTrack) AddLayer(layer *MusicLayer) {
	amt.Layers = append(amt.Layers, layer)
	amt.CurrentMix[layer.Name] = 0.0
}

// SetIntensity updates the target intensity level
func (amt *AdaptiveMusicTrack) SetIntensity(intensity MusicIntensity) {
	amt.TargetIntensity = intensity
}

// Update updates the music mix based on current/target intensity
func (amt *AdaptiveMusicTrack) Update() {
	// Smoothly transition current intensity toward target
	if amt.CurrentIntensity < amt.TargetIntensity {
		amt.CurrentIntensity++
	} else if amt.CurrentIntensity > amt.TargetIntensity {
		amt.CurrentIntensity--
	}
	
	// Update layer volumes based on current intensity
	for _, layer := range amt.Layers {
		targetVolume := amt.calculateLayerVolume(layer)
		currentVolume := amt.CurrentMix[layer.Name]
		
		// Smooth crossfade toward target volume
		diff := targetVolume - currentVolume
		amt.CurrentMix[layer.Name] = currentVolume + diff*amt.TransitionSpeed
	}
}

// calculateLayerVolume determines the appropriate volume for a layer
func (amt *AdaptiveMusicTrack) calculateLayerVolume(layer *MusicLayer) float64 {
	// Layer is silent if current intensity is below its minimum
	if amt.CurrentIntensity < layer.MinIntensity {
		return 0.0
	}
	
	// If we're at minimum intensity, use a base volume level (30%)
	if amt.CurrentIntensity == layer.MinIntensity {
		return layer.BaseVolume * 0.3
	}
	
	// Calculate volume based on how far above minimum intensity we are
	intensityDiff := float64(amt.CurrentIntensity - layer.MinIntensity)
	maxIntensityDiff := float64(IntensityBoss - layer.MinIntensity)
	
	if maxIntensityDiff <= 0 {
		return layer.BaseVolume
	}
	
	// Smooth volume curve from 30% to 100%
	volumeFactor := 0.3 + (0.7 * math.Min(1.0, intensityDiff/maxIntensityDiff))
	return layer.BaseVolume * volumeFactor
}

// GetCurrentMix returns the current volume mix for all layers
func (amt *AdaptiveMusicTrack) GetCurrentMix() map[string]float64 {
	return amt.CurrentMix
}

// MusicContext tracks the game state for adaptive music
type MusicContext struct {
	InCombat         bool
	IsBossFight      bool
	NearbyEnemyCount int
	PlayerHealthPct  float64
	RoomDangerLevel  int
}

// NewMusicContext creates a new music context
func NewMusicContext() *MusicContext {
	return &MusicContext{
		InCombat:         false,
		IsBossFight:      false,
		NearbyEnemyCount: 0,
		PlayerHealthPct:  1.0,
		RoomDangerLevel:  0,
	}
}

// CalculateIntensity determines the appropriate music intensity
func (mc *MusicContext) CalculateIntensity() MusicIntensity {
	// Boss fight always has highest priority
	if mc.IsBossFight {
		return IntensityBoss
	}
	
	// Active combat
	if mc.InCombat {
		return IntensityCombat
	}
	
	// Tension from nearby enemies or low health
	if mc.NearbyEnemyCount > 0 || mc.PlayerHealthPct < 0.3 || mc.RoomDangerLevel >= 7 {
		return IntensityTension
	}
	
	// Peaceful exploration
	return IntensityCalm
}

// GenerateAdaptiveMusicTrack creates an adaptive track with multiple layers
func (mg *MusicGenerator) GenerateAdaptiveMusicTrack(seed int64, duration float64) *AdaptiveMusicTrack {
	track := NewAdaptiveMusicTrack()
	
	// Generate chord progression (shared by all layers)
	rng := rand.New(rand.NewSource(seed))
	progression := mg.generateProgression(rng, 4)
	
	// Layer 1: Ambient pads (always present at low intensity)
	pads := mg.generatePads(progression, rng)
	track.AddLayer(&MusicLayer{
		Name:         "pads",
		Audio:        pads,
		BaseVolume:   0.15,
		MinIntensity: IntensityCalm,
	})
	
	// Layer 2: Light melody (exploration)
	melody := mg.generateMelody(progression, rng)
	track.AddLayer(&MusicLayer{
		Name:         "melody",
		Audio:        melody,
		BaseVolume:   0.20,
		MinIntensity: IntensityCalm,
	})
	
	// Layer 3: Bassline (tension and up)
	bassline := mg.generateBassline(progression, rng)
	track.AddLayer(&MusicLayer{
		Name:         "bass",
		Audio:        bassline,
		BaseVolume:   0.25,
		MinIntensity: IntensityTension,
	})
	
	// Layer 4: Drums (combat intensity)
	drums := mg.generateDrumPattern(rng, duration)
	track.AddLayer(&MusicLayer{
		Name:         "drums",
		Audio:        drums,
		BaseVolume:   0.30,
		MinIntensity: IntensityCombat,
	})
	
	// Layer 5: Intense lead (boss fights)
	intenseLead := mg.generateIntenseLead(progression, rng)
	track.AddLayer(&MusicLayer{
		Name:         "lead",
		Audio:        intenseLead,
		BaseVolume:   0.35,
		MinIntensity: IntensityBoss,
	})
	
	return track
}

// generateIntenseLead creates an intense lead melody for boss fights
func (mg *MusicGenerator) generateIntenseLead(progression ChordProgression, rng *rand.Rand) *AudioSample {
	beatDuration := 60.0 / float64(mg.BPM)
	
	totalSamples := 0
	for range progression {
		// 2 beats per chord for intense lead
		totalSamples += int(beatDuration * 2.0 * float64(mg.Synth.SampleRate))
	}
	
	data := make([]float64, totalSamples)
	sampleIndex := 0
	
	for _, chord := range progression {
		chordDuration := beatDuration * 2.0
		
		// Fast arpeggio for intensity
		numNotes := 8
		noteDuration := chordDuration / float64(numNotes)
		
		for i := 0; i < numNotes; i++ {
			noteIndex := i % len(chord.Intervals)
			frequency := mg.midiToFreq(chord.Root + chord.Intervals[noteIndex] + 12) // One octave up
			
			// Use square wave for intensity
			note := mg.Synth.GenerateWave(SquareWave, frequency, noteDuration)
			envelope := ADSR{Attack: 0.01, Decay: 0.05, Sustain: 0.7, Release: 0.1}
			note = mg.Synth.ApplyEnvelope(note, envelope)
			
			// Copy to output
			for j := 0; j < len(note.Data) && sampleIndex < len(data); j++ {
				data[sampleIndex] = note.Data[j] * 0.8 // Slightly louder
				sampleIndex++
			}
		}
	}
	
	return &AudioSample{
		Data:       data,
		SampleRate: mg.Synth.SampleRate,
		Duration:   float64(len(data)) / float64(mg.Synth.SampleRate),
	}
}
