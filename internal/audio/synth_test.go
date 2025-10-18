package audio

import (
	"math"
	"testing"
)

func TestWaveformGeneration(t *testing.T) {
	synth := NewSynthesizer(44100)
	
	testCases := []struct {
		name     string
		waveType WaveType
	}{
		{"Sine", SineWave},
		{"Square", SquareWave},
		{"Sawtooth", SawtoothWave},
		{"Triangle", TriangleWave},
		{"Noise", NoiseWave},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sample := synth.GenerateWave(tc.waveType, 440.0, 0.1)
			
			if sample == nil {
				t.Fatal("Generated sample is nil")
			}
			
			expectedSamples := int(0.1 * 44100)
			if len(sample.Data) != expectedSamples {
				t.Errorf("Sample length mismatch: got %d, want %d", len(sample.Data), expectedSamples)
			}
			
			// Check that values are in valid range [-1, 1]
			for i, val := range sample.Data {
				if math.Abs(val) > 1.0 {
					t.Errorf("Sample value out of range at index %d: %f", i, val)
				}
			}
		})
	}
}

func TestADSREnvelope(t *testing.T) {
	synth := NewSynthesizer(44100)
	
	sample := synth.GenerateWave(SineWave, 440.0, 1.0)
	
	envelope := ADSR{
		Attack:  0.1,
		Decay:   0.2,
		Sustain: 0.7,
		Release: 0.3,
	}
	
	enveloped := synth.ApplyEnvelope(sample, envelope)
	
	if enveloped == nil {
		t.Fatal("Enveloped sample is nil")
	}
	
	if len(enveloped.Data) != len(sample.Data) {
		t.Error("Envelope changed sample length")
	}
	
	// Check that envelope actually affected the amplitude
	hasLowAmplitude := false
	for _, val := range enveloped.Data[:1000] {
		if math.Abs(val) < 0.5 {
			hasLowAmplitude = true
			break
		}
	}
	
	if !hasLowAmplitude {
		t.Error("Envelope doesn't appear to affect amplitude")
	}
}

func TestSFXGeneration(t *testing.T) {
	gen := NewSFXGenerator(44100)
	
	testCases := []struct {
		name    string
		sfxType SFXType
	}{
		{"Jump", JumpSFX},
		{"Land", LandSFX},
		{"Attack", AttackSFX},
		{"Hit", HitSFX},
		{"Pickup", PickupSFX},
		{"Door", DoorSFX},
		{"Damage", DamageSFX},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sample := gen.Generate(tc.sfxType, 12345)
			
			if sample == nil {
				t.Fatal("Generated SFX is nil")
			}
			
			if len(sample.Data) == 0 {
				t.Error("SFX has no data")
			}
			
			// Verify sample is within bounds
			for _, val := range sample.Data {
				if math.Abs(val) > 1.0 {
					t.Errorf("SFX value out of range: %f", val)
				}
			}
		})
	}
}

func TestMusicGeneration(t *testing.T) {
	gen := NewMusicGenerator(44100, 90, 60, MinorScale)
	
	track := gen.GenerateTrack(42, 5.0)
	
	if track == nil {
		t.Fatal("Generated track is nil")
	}
	
	if len(track.Data) == 0 {
		t.Error("Track has no data")
	}
	
	// Music generation concatenates layers, so final length may be longer
	// Just verify we got a reasonable amount of audio
	minSamples := int(5.0 * 44100)
	
	if len(track.Data) < minSamples {
		t.Errorf("Track too short: got %d samples, want at least %d", len(track.Data), minSamples)
	}
}

func TestAudioMixing(t *testing.T) {
	synth := NewSynthesizer(44100)
	
	sample1 := synth.GenerateWave(SineWave, 440.0, 0.1)
	sample2 := synth.GenerateWave(SineWave, 880.0, 0.1)
	
	mixed := synth.Mix([]*AudioSample{sample1, sample2}, []float64{0.5, 0.5})
	
	if mixed == nil {
		t.Fatal("Mixed sample is nil")
	}
	
	if len(mixed.Data) == 0 {
		t.Error("Mixed sample has no data")
	}
	
	// Check normalization - values should be in [-1, 1]
	for _, val := range mixed.Data {
		if math.Abs(val) > 1.0 {
			t.Errorf("Mixed value out of range: %f", val)
		}
	}
}

func TestLowPassFilter(t *testing.T) {
	synth := NewSynthesizer(44100)
	
	// Generate noise
	noise := synth.GenerateWave(NoiseWave, 0, 0.1)
	
	// Apply low-pass filter
	filtered := synth.ApplyLowPassFilter(noise, 1000.0)
	
	if filtered == nil {
		t.Fatal("Filtered sample is nil")
	}
	
	if len(filtered.Data) != len(noise.Data) {
		t.Error("Filter changed sample length")
	}
	
	// Filtered signal should generally have lower high-frequency content
	// This is a basic check - not a full spectral analysis
	if filtered.Data[0] == noise.Data[0] {
		// First sample might be the same, check a few more
		allSame := true
		for i := 1; i < 10; i++ {
			if filtered.Data[i] != noise.Data[i] {
				allSame = false
				break
			}
		}
		if allSame {
			t.Error("Filter doesn't appear to modify signal")
		}
	}
}

func TestFrequencySweep(t *testing.T) {
	synth := NewSynthesizer(44100)
	
	sweep := synth.FrequencySweep(SineWave, 100.0, 400.0, 0.2)
	
	if sweep == nil {
		t.Fatal("Frequency sweep is nil")
	}
	
	expectedSamples := int(0.2 * 44100)
	if len(sweep.Data) != expectedSamples {
		t.Errorf("Sweep length mismatch: got %d, want %d", len(sweep.Data), expectedSamples)
	}
}
