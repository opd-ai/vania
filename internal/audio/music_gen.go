package audio

import (
	"math"
	"math/rand"
)

// Scale represents a musical scale
type Scale []int

var (
	// Common musical scales (semitone intervals from root)
	MajorScale     = Scale{0, 2, 4, 5, 7, 9, 11}
	MinorScale     = Scale{0, 2, 3, 5, 7, 8, 10}
	DorianScale    = Scale{0, 2, 3, 5, 7, 9, 10}
	PhrygianScale  = Scale{0, 1, 3, 5, 7, 8, 10}
	PentatonicMaj  = Scale{0, 2, 4, 7, 9}
	PentatonicMin  = Scale{0, 3, 5, 7, 10}
)

// Chord represents a musical chord
type Chord struct {
	Root     int   // MIDI note
	Intervals []int // Semitone intervals from root
}

// ChordProgression is a sequence of chords
type ChordProgression []Chord

// MusicGenerator generates procedural music
type MusicGenerator struct {
	Synth      *Synthesizer
	BPM        int
	Scale      Scale
	RootNote   int // MIDI note
}

// NewMusicGenerator creates a new music generator
func NewMusicGenerator(sampleRate, bpm, rootNote int, scale Scale) *MusicGenerator {
	return &MusicGenerator{
		Synth:    NewSynthesizer(sampleRate),
		BPM:      bpm,
		Scale:    scale,
		RootNote: rootNote,
	}
}

// GenerateTrack creates a complete music track
func (mg *MusicGenerator) GenerateTrack(seed int64, duration float64) *AudioSample {
	rng := rand.New(rand.NewSource(seed))
	
	// Generate chord progression
	progression := mg.generateProgression(rng, 4)
	
	// Generate musical layers
	bassline := mg.generateBassline(progression, rng)
	melody := mg.generateMelody(progression, rng)
	pads := mg.generatePads(progression, rng)
	drums := mg.generateDrumPattern(rng, duration)
	
	// Mix layers with appropriate volumes
	mixed := mg.Synth.Mix(
		[]*AudioSample{bassline, melody, pads, drums},
		[]float64{0.3, 0.25, 0.2, 0.25},
	)
	
	return mixed
}

// generateProgression creates a chord progression
func (mg *MusicGenerator) generateProgression(rng *rand.Rand, length int) ChordProgression {
	progression := make(ChordProgression, length)
	
	// Common progressions: I-IV-V-I, I-V-vi-IV, etc.
	progressionPatterns := [][]int{
		{0, 3, 4, 0}, // I-IV-V-I
		{0, 4, 5, 3}, // I-V-vi-IV
		{0, 5, 3, 4}, // I-vi-IV-V
		{1, 4, 0, 0}, // ii-V-I-I
	}
	
	pattern := progressionPatterns[rng.Intn(len(progressionPatterns))]
	
	for i := 0; i < length; i++ {
		degree := pattern[i%len(pattern)]
		root := mg.RootNote + mg.Scale[degree%len(mg.Scale)]
		
		// Generate chord type (major or minor)
		if rng.Float64() < 0.7 {
			// Major chord
			progression[i] = Chord{
				Root:     root,
				Intervals: []int{0, 4, 7}, // Root, major third, perfect fifth
			}
		} else {
			// Minor chord
			progression[i] = Chord{
				Root:     root,
				Intervals: []int{0, 3, 7}, // Root, minor third, perfect fifth
			}
		}
	}
	
	return progression
}

// generateBassline creates bass notes following chord roots
func (mg *MusicGenerator) generateBassline(progression ChordProgression, rng *rand.Rand) *AudioSample {
	beatDuration := 60.0 / float64(mg.BPM)
	chordDuration := beatDuration * 4 // 4 beats per chord
	
	var allSamples []*AudioSample
	
	for _, chord := range progression {
		// Bass plays root note, one octave down
		freq := mg.midiToFreq(chord.Root - 12)
		
		sample := mg.Synth.GenerateWave(SawtoothWave, freq, chordDuration)
		
		envelope := ADSR{
			Attack:  0.01,
			Decay:   0.1,
			Sustain: 0.7,
			Release: 0.1,
		}
		
		sample = mg.Synth.ApplyEnvelope(sample, envelope)
		allSamples = append(allSamples, sample)
	}
	
	// Concatenate all bass notes
	return mg.concatenateSamples(allSamples)
}

// generateMelody creates a melody over the chord progression
func (mg *MusicGenerator) generateMelody(progression ChordProgression, rng *rand.Rand) *AudioSample {
	beatDuration := 60.0 / float64(mg.BPM)
	noteDuration := beatDuration / 2 // Eighth notes
	
	var allSamples []*AudioSample
	
	for range progression {
		// Generate 8 melody notes per chord
		for i := 0; i < 8; i++ {
			// Choose note from scale
			scaleIdx := rng.Intn(len(mg.Scale))
			note := mg.RootNote + mg.Scale[scaleIdx] + 12 // One octave up
			
			// Sometimes rest
			if rng.Float64() < 0.2 {
				silence := &AudioSample{
					Data:       make([]float64, int(noteDuration*float64(mg.Synth.SampleRate))),
					SampleRate: mg.Synth.SampleRate,
					Duration:   noteDuration,
				}
				allSamples = append(allSamples, silence)
				continue
			}
			
			freq := mg.midiToFreq(note)
			sample := mg.Synth.GenerateWave(SquareWave, freq, noteDuration)
			
			envelope := ADSR{
				Attack:  0.01,
				Decay:   0.05,
				Sustain: 0.6,
				Release: 0.1,
			}
			
			sample = mg.Synth.ApplyEnvelope(sample, envelope)
			allSamples = append(allSamples, sample)
		}
	}
	
	return mg.concatenateSamples(allSamples)
}

// generatePads creates ambient pad sounds
func (mg *MusicGenerator) generatePads(progression ChordProgression, rng *rand.Rand) *AudioSample {
	beatDuration := 60.0 / float64(mg.BPM)
	chordDuration := beatDuration * 4
	
	var allSamples []*AudioSample
	
	for _, chord := range progression {
		// Play all notes of the chord
		var chordSamples []*AudioSample
		
		for _, interval := range chord.Intervals {
			freq := mg.midiToFreq(chord.Root + interval)
			sample := mg.Synth.GenerateWave(SineWave, freq, chordDuration)
			chordSamples = append(chordSamples, sample)
		}
		
		// Mix chord notes
		mixed := mg.Synth.Mix(chordSamples, nil)
		
		envelope := ADSR{
			Attack:  0.3,
			Decay:   0.2,
			Sustain: 0.6,
			Release: 0.3,
		}
		
		mixed = mg.Synth.ApplyEnvelope(mixed, envelope)
		allSamples = append(allSamples, mixed)
	}
	
	return mg.concatenateSamples(allSamples)
}

// generateDrumPattern creates a drum pattern
func (mg *MusicGenerator) generateDrumPattern(rng *rand.Rand, duration float64) *AudioSample {
	beatDuration := 60.0 / float64(mg.BPM)
	numBeats := int(duration / beatDuration)
	
	if numBeats == 0 {
		numBeats = 16 // Default to 16 beats
	}
	
	var allSamples []*AudioSample
	
	for beat := 0; beat < numBeats; beat++ {
		// Kick on beats 1 and 3
		if beat%4 == 0 || beat%4 == 2 {
			kick := mg.generateKick()
			allSamples = append(allSamples, kick)
		}
		
		// Snare on beats 2 and 4
		if beat%4 == 1 || beat%4 == 3 {
			snare := mg.generateSnare(rng)
			allSamples = append(allSamples, snare)
		}
		
		// Hi-hat every beat
		hihat := mg.generateHiHat(rng)
		allSamples = append(allSamples, hihat)
	}
	
	return mg.concatenateSamples(allSamples)
}

// generateKick creates a kick drum sound
func (mg *MusicGenerator) generateKick() *AudioSample {
	duration := 0.15
	
	// Frequency sweep from 150Hz to 40Hz
	sample := mg.Synth.FrequencySweep(SineWave, 150, 40, duration)
	
	envelope := ADSR{
		Attack:  0.005,
		Decay:   0.05,
		Sustain: 0.3,
		Release: 0.08,
	}
	
	return mg.Synth.ApplyEnvelope(sample, envelope)
}

// generateSnare creates a snare drum sound
func (mg *MusicGenerator) generateSnare(rng *rand.Rand) *AudioSample {
	duration := 0.1
	
	// Mix tone and noise for snare
	tone := mg.Synth.GenerateWave(SineWave, 180, duration)
	noise := mg.Synth.GenerateWave(NoiseWave, 0, duration)
	noise = mg.Synth.ApplyLowPassFilter(noise, 5000)
	
	mixed := mg.Synth.Mix([]*AudioSample{tone, noise}, []float64{0.3, 0.7})
	
	envelope := ADSR{
		Attack:  0.005,
		Decay:   0.03,
		Sustain: 0.2,
		Release: 0.05,
	}
	
	return mg.Synth.ApplyEnvelope(mixed, envelope)
}

// generateHiHat creates a hi-hat sound
func (mg *MusicGenerator) generateHiHat(rng *rand.Rand) *AudioSample {
	duration := 0.05
	
	// High frequency noise
	sample := mg.Synth.GenerateWave(NoiseWave, 0, duration)
	sample = mg.Synth.ApplyLowPassFilter(sample, 8000)
	
	envelope := ADSR{
		Attack:  0.001,
		Decay:   0.01,
		Sustain: 0.1,
		Release: 0.02,
	}
	
	return mg.Synth.ApplyEnvelope(sample, envelope)
}

// midiToFreq converts MIDI note number to frequency
func (mg *MusicGenerator) midiToFreq(midiNote int) float64 {
	return 440.0 * math.Pow(2.0, float64(midiNote-69)/12.0)
}

// concatenateSamples joins audio samples sequentially
func (mg *MusicGenerator) concatenateSamples(samples []*AudioSample) *AudioSample {
	if len(samples) == 0 {
		return &AudioSample{
			Data:       []float64{},
			SampleRate: mg.Synth.SampleRate,
			Duration:   0,
		}
	}
	
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
		SampleRate: mg.Synth.SampleRate,
		Duration:   float64(totalLen) / float64(mg.Synth.SampleRate),
	}
}
