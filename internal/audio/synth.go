// Package audio implements procedural audio synthesis and music generation
// using waveform generators, ADSR envelopes, and music theory to create
// dynamic sound effects and biome-specific musical tracks at runtime.
package audio

import (
	"math"
	"math/rand"
)

// WaveType defines waveform types
type WaveType int

const (
	SineWave WaveType = iota
	SquareWave
	SawtoothWave
	TriangleWave
	NoiseWave
)

// ADSR envelope parameters
type ADSR struct {
	Attack  float64 // seconds
	Decay   float64 // seconds
	Sustain float64 // amplitude (0-1)
	Release float64 // seconds
}

// AudioSample represents generated audio data
type AudioSample struct {
	Data       []float64
	SampleRate int
	Duration   float64
}

// Synthesizer generates waveforms
type Synthesizer struct {
	SampleRate int
}

// NewSynthesizer creates a new synthesizer
func NewSynthesizer(sampleRate int) *Synthesizer {
	return &Synthesizer{
		SampleRate: sampleRate,
	}
}

// GenerateWave creates a waveform
func (s *Synthesizer) GenerateWave(waveType WaveType, frequency, duration float64) *AudioSample {
	numSamples := int(duration * float64(s.SampleRate))
	data := make([]float64, numSamples)
	
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(s.SampleRate)
		phase := 2.0 * math.Pi * frequency * t
		
		switch waveType {
		case SineWave:
			data[i] = math.Sin(phase)
		case SquareWave:
			if math.Sin(phase) >= 0 {
				data[i] = 1.0
			} else {
				data[i] = -1.0
			}
		case SawtoothWave:
			data[i] = 2.0*(phase/(2.0*math.Pi)-math.Floor(phase/(2.0*math.Pi)+0.5))
		case TriangleWave:
			data[i] = 2.0*math.Abs(2.0*(phase/(2.0*math.Pi)-math.Floor(phase/(2.0*math.Pi)+0.5))) - 1.0
		case NoiseWave:
			data[i] = rand.Float64()*2.0 - 1.0
		}
	}
	
	return &AudioSample{
		Data:       data,
		SampleRate: s.SampleRate,
		Duration:   duration,
	}
}

// ApplyEnvelope applies ADSR envelope to audio
func (s *Synthesizer) ApplyEnvelope(sample *AudioSample, envelope ADSR) *AudioSample {
	numSamples := len(sample.Data)
	result := make([]float64, numSamples)
	
	attackSamples := int(envelope.Attack * float64(s.SampleRate))
	decaySamples := int(envelope.Decay * float64(s.SampleRate))
	releaseSamples := int(envelope.Release * float64(s.SampleRate))
	sustainSamples := numSamples - attackSamples - decaySamples - releaseSamples
	
	if sustainSamples < 0 {
		sustainSamples = 0
	}
	
	for i := 0; i < numSamples; i++ {
		var amplitude float64
		
		if i < attackSamples {
			// Attack phase - linear ramp up
			amplitude = float64(i) / float64(attackSamples)
		} else if i < attackSamples+decaySamples {
			// Decay phase - linear ramp down to sustain
			t := float64(i-attackSamples) / float64(decaySamples)
			amplitude = 1.0 + t*(envelope.Sustain-1.0)
		} else if i < attackSamples+decaySamples+sustainSamples {
			// Sustain phase - constant
			amplitude = envelope.Sustain
		} else {
			// Release phase - linear ramp down
			t := float64(i-attackSamples-decaySamples-sustainSamples) / float64(releaseSamples)
			amplitude = envelope.Sustain * (1.0 - t)
		}
		
		result[i] = sample.Data[i] * amplitude
	}
	
	return &AudioSample{
		Data:       result,
		SampleRate: s.SampleRate,
		Duration:   sample.Duration,
	}
}

// ApplyLowPassFilter applies simple low-pass filter
func (s *Synthesizer) ApplyLowPassFilter(sample *AudioSample, cutoff float64) *AudioSample {
	result := make([]float64, len(sample.Data))
	
	// Simple RC low-pass filter
	rc := 1.0 / (2.0 * math.Pi * cutoff)
	dt := 1.0 / float64(s.SampleRate)
	alpha := dt / (rc + dt)
	
	result[0] = sample.Data[0]
	for i := 1; i < len(sample.Data); i++ {
		result[i] = result[i-1] + alpha*(sample.Data[i]-result[i-1])
	}
	
	return &AudioSample{
		Data:       result,
		SampleRate: s.SampleRate,
		Duration:   sample.Duration,
	}
}

// Mix combines multiple audio samples
func (s *Synthesizer) Mix(samples []*AudioSample, volumes []float64) *AudioSample {
	if len(samples) == 0 {
		return &AudioSample{
			Data:       []float64{},
			SampleRate: s.SampleRate,
			Duration:   0,
		}
	}
	
	// Find longest sample
	maxLen := 0
	for _, sample := range samples {
		if len(sample.Data) > maxLen {
			maxLen = len(sample.Data)
		}
	}
	
	result := make([]float64, maxLen)
	
	for i, sample := range samples {
		volume := 1.0
		if i < len(volumes) {
			volume = volumes[i]
		}
		
		for j := 0; j < len(sample.Data) && j < maxLen; j++ {
			result[j] += sample.Data[j] * volume
		}
	}
	
	// Normalize to prevent clipping
	maxAmp := 0.0
	for _, val := range result {
		if math.Abs(val) > maxAmp {
			maxAmp = math.Abs(val)
		}
	}
	
	if maxAmp > 1.0 {
		for i := range result {
			result[i] /= maxAmp
		}
	}
	
	return &AudioSample{
		Data:       result,
		SampleRate: s.SampleRate,
		Duration:   float64(maxLen) / float64(s.SampleRate),
	}
}

// FrequencySweep generates a frequency sweep
func (s *Synthesizer) FrequencySweep(waveType WaveType, startFreq, endFreq, duration float64) *AudioSample {
	numSamples := int(duration * float64(s.SampleRate))
	data := make([]float64, numSamples)
	
	phase := 0.0
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := startFreq + (endFreq-startFreq)*t
		
		dt := 1.0 / float64(s.SampleRate)
		phase += 2.0 * math.Pi * freq * dt
		
		switch waveType {
		case SineWave:
			data[i] = math.Sin(phase)
		case SquareWave:
			if math.Sin(phase) >= 0 {
				data[i] = 1.0
			} else {
				data[i] = -1.0
			}
		case SawtoothWave:
			data[i] = 2.0*(phase/(2.0*math.Pi)-math.Floor(phase/(2.0*math.Pi)+0.5))
		case TriangleWave:
			data[i] = 2.0*math.Abs(2.0*(phase/(2.0*math.Pi)-math.Floor(phase/(2.0*math.Pi)+0.5))) - 1.0
		}
	}
	
	return &AudioSample{
		Data:       data,
		SampleRate: s.SampleRate,
		Duration:   duration,
	}
}
