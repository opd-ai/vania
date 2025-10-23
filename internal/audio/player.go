// Package audio provides Ebiten-compatible audio playback integration
// for the procedurally generated audio samples and adaptive music system.
package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// AudioPlayer manages audio playback using Ebiten's audio system
type AudioPlayer struct {
	audioContext  *audio.Context
	players       map[string]*audio.Player
	musicPlayer   *audio.Player
	currentTrack  string
	masterVolume  float64
	sfxVolume     float64
	musicVolume   float64
	adaptiveTrack *AdaptiveMusicTrack
	musicContext  *MusicContext
}

// NewAudioPlayer creates a new audio player
func NewAudioPlayer() (*AudioPlayer, error) {
	// Create audio context with standard sample rate
	audioContext := audio.NewContext(44100)

	return &AudioPlayer{
		audioContext: audioContext,
		players:      make(map[string]*audio.Player),
		masterVolume: 0.7,
		sfxVolume:    0.8,
		musicVolume:  0.6,
		musicContext: NewMusicContext(),
	}, nil
}

// LoadSound converts AudioSample to Ebiten-compatible format and loads it
func (ap *AudioPlayer) LoadSound(name string, sample *AudioSample) error {
	if sample == nil {
		return fmt.Errorf("audio sample is nil")
	}

	// Convert float64 samples to 16-bit PCM
	wavData, err := ap.audioSampleToWAV(sample)
	if err != nil {
		return fmt.Errorf("failed to convert audio sample: %v", err)
	}

	// Create WAV reader
	wavReader := bytes.NewReader(wavData)

	// Decode WAV data
	stream, err := wav.DecodeWithoutResampling(wavReader)
	if err != nil {
		return fmt.Errorf("failed to decode WAV: %v", err)
	}

	// Create audio player
	player, err := ap.audioContext.NewPlayer(stream)
	if err != nil {
		return fmt.Errorf("failed to create audio player: %v", err)
	}

	// Store player
	ap.players[name] = player

	return nil
}

// PlaySound plays a loaded sound effect
func (ap *AudioPlayer) PlaySound(name string) error {
	player, exists := ap.players[name]
	if !exists {
		return fmt.Errorf("sound '%s' not loaded", name)
	}

	// Reset and play from beginning
	if err := player.Rewind(); err != nil {
		return fmt.Errorf("failed to rewind sound: %v", err)
	}

	// Set volume (SFX volume * master volume)
	player.SetVolume(ap.sfxVolume * ap.masterVolume)

	player.Play()
	return nil
}

// LoadMusic loads an adaptive music track
func (ap *AudioPlayer) LoadMusic(track *AdaptiveMusicTrack) error {
	if track == nil {
		return fmt.Errorf("adaptive track is nil")
	}

	ap.adaptiveTrack = track
	return nil
}

// PlayMusic starts playing the loaded music
func (ap *AudioPlayer) PlayMusic() error {
	if ap.adaptiveTrack == nil {
		return fmt.Errorf("no music loaded")
	}

	// For now, just play the first available layer
	// In a full implementation, you'd mix layers dynamically
	if len(ap.adaptiveTrack.Layers) == 0 {
		return fmt.Errorf("no music layers available")
	}

	// Get the base layer (usually pads or melody)
	baseLayer := ap.adaptiveTrack.Layers[0]
	if baseLayer.Audio == nil {
		return fmt.Errorf("base layer has no audio data")
	}

	// Convert to WAV and create player
	wavData, err := ap.audioSampleToWAV(baseLayer.Audio)
	if err != nil {
		return fmt.Errorf("failed to convert music to WAV: %v", err)
	}

	wavReader := bytes.NewReader(wavData)
	stream, err := wav.DecodeWithoutResampling(wavReader)
	if err != nil {
		return fmt.Errorf("failed to decode music WAV: %v", err)
	}

	// Create looping stream
	loopStream := audio.NewInfiniteLoop(stream, int64(baseLayer.Audio.Duration*float64(baseLayer.Audio.SampleRate)))

	player, err := ap.audioContext.NewPlayer(loopStream)
	if err != nil {
		return fmt.Errorf("failed to create music player: %v", err)
	}

	// Stop current music if playing
	if ap.musicPlayer != nil {
		ap.musicPlayer.Close()
	}

	ap.musicPlayer = player
	ap.musicPlayer.SetVolume(ap.musicVolume * ap.masterVolume)
	ap.musicPlayer.Play()

	return nil
}

// StopMusic stops the current music
func (ap *AudioPlayer) StopMusic() {
	if ap.musicPlayer != nil {
		ap.musicPlayer.Pause()
	}
}

// UpdateMusic updates adaptive music based on game context
func (ap *AudioPlayer) UpdateMusic(context *MusicContext) {
	if ap.adaptiveTrack == nil || ap.musicPlayer == nil {
		return
	}

	// Update music context
	ap.musicContext = context

	// Calculate intensity
	intensity := context.CalculateIntensity()
	ap.adaptiveTrack.SetIntensity(intensity)
	ap.adaptiveTrack.Update()

	// Adjust volume based on intensity and health
	baseVolume := ap.musicVolume * ap.masterVolume

	// Increase volume during combat
	intensityMultiplier := 1.0
	switch intensity {
	case IntensityTension:
		intensityMultiplier = 1.1
	case IntensityCombat:
		intensityMultiplier = 1.2
	case IntensityBoss:
		intensityMultiplier = 1.3
	}

	// Decrease volume when player has low health (adds tension)
	healthMultiplier := 0.7 + 0.3*context.PlayerHealthPct

	finalVolume := baseVolume * intensityMultiplier * healthMultiplier
	if finalVolume > 1.0 {
		finalVolume = 1.0
	}

	ap.musicPlayer.SetVolume(finalVolume)
}

// SetVolumes sets the volume levels
func (ap *AudioPlayer) SetVolumes(master, sfx, music float64) {
	ap.masterVolume = math.Max(0.0, math.Min(1.0, master))
	ap.sfxVolume = math.Max(0.0, math.Min(1.0, sfx))
	ap.musicVolume = math.Max(0.0, math.Min(1.0, music))

	// Update current music volume
	if ap.musicPlayer != nil {
		ap.musicPlayer.SetVolume(ap.musicVolume * ap.masterVolume)
	}

	// Update SFX volumes (they'll be applied on next play)
}

// GetVolumes returns current volume levels
func (ap *AudioPlayer) GetVolumes() (master, sfx, music float64) {
	return ap.masterVolume, ap.sfxVolume, ap.musicVolume
}

// IsPlaying checks if a sound is currently playing
func (ap *AudioPlayer) IsPlaying(name string) bool {
	player, exists := ap.players[name]
	if !exists {
		return false
	}
	return player.IsPlaying()
}

// IsMusicPlaying checks if music is currently playing
func (ap *AudioPlayer) IsMusicPlaying() bool {
	return ap.musicPlayer != nil && ap.musicPlayer.IsPlaying()
}

// Close releases all audio resources
func (ap *AudioPlayer) Close() {
	// Close all sound players
	for _, player := range ap.players {
		player.Close()
	}

	// Close music player
	if ap.musicPlayer != nil {
		ap.musicPlayer.Close()
	}

	ap.players = make(map[string]*audio.Player)
	ap.musicPlayer = nil
}

// audioSampleToWAV converts AudioSample to WAV format bytes
func (ap *AudioPlayer) audioSampleToWAV(sample *AudioSample) ([]byte, error) {
	if sample == nil || len(sample.Data) == 0 {
		return nil, fmt.Errorf("empty audio sample")
	}

	// Create buffer for WAV data
	buf := new(bytes.Buffer)

	// WAV header
	sampleRate := int32(sample.SampleRate)
	numSamples := int32(len(sample.Data))
	bitsPerSample := int16(16)
	numChannels := int16(1) // Mono
	byteRate := sampleRate * int32(numChannels) * int32(bitsPerSample) / 8
	blockAlign := numChannels * bitsPerSample / 8
	dataSize := numSamples * int32(blockAlign)

	// RIFF header
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, int32(36+dataSize))
	buf.WriteString("WAVE")

	// Format chunk
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, int32(16)) // Chunk size
	binary.Write(buf, binary.LittleEndian, int16(1))  // Audio format (PCM)
	binary.Write(buf, binary.LittleEndian, numChannels)
	binary.Write(buf, binary.LittleEndian, sampleRate)
	binary.Write(buf, binary.LittleEndian, byteRate)
	binary.Write(buf, binary.LittleEndian, blockAlign)
	binary.Write(buf, binary.LittleEndian, bitsPerSample)

	// Data chunk
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, dataSize)

	// Convert float64 samples to 16-bit PCM
	for _, sample := range sample.Data {
		// Clamp to [-1.0, 1.0] and convert to 16-bit
		clamped := math.Max(-1.0, math.Min(1.0, sample))
		pcmValue := int16(clamped * 32767.0)
		binary.Write(buf, binary.LittleEndian, pcmValue)
	}

	return buf.Bytes(), nil
}

// LoadSoundFromGenerator generates and loads a sound effect
func (ap *AudioPlayer) LoadSoundFromGenerator(name string, gen *SFXGenerator, soundType SFXType) error {
	// Generate the sound
	sample := gen.Generate(soundType, 42) // Use fixed seed for consistent sounds
	if sample == nil {
		return fmt.Errorf("failed to generate sound '%s'", name)
	}

	// Load into player
	return ap.LoadSound(name, sample)
}

// PreloadGameSounds loads commonly used game sounds
func (ap *AudioPlayer) PreloadGameSounds(sfxGen *SFXGenerator) error {
	soundTypes := map[string]SFXType{
		"jump":   JumpSFX,
		"attack": AttackSFX,
		"hit":    HitSFX,
		"pickup": PickupSFX,
		"door":   DoorSFX,
		"damage": DamageSFX,
		"land":   LandSFX,
	}

	for name, soundType := range soundTypes {
		if err := ap.LoadSoundFromGenerator(name, sfxGen, soundType); err != nil {
			return fmt.Errorf("failed to preload sound '%s': %v", name, err)
		}
	}

	return nil
}

// CreateInfiniteLoopReader creates an infinite loop reader for music
func (ap *AudioPlayer) CreateInfiniteLoopReader(sample *AudioSample) (io.ReadSeeker, error) {
	wavData, err := ap.audioSampleToWAV(sample)
	if err != nil {
		return nil, err
	}

	wavReader := bytes.NewReader(wavData)
	stream, err := wav.DecodeWithoutResampling(wavReader)
	if err != nil {
		return nil, err
	}

	return audio.NewInfiniteLoop(stream, int64(sample.Duration*float64(sample.SampleRate))), nil
}

// PlaySoundWithPitch plays a sound with pitch modification
func (ap *AudioPlayer) PlaySoundWithPitch(name string, pitchFactor float64) error {
	_, exists := ap.players[name]
	if !exists {
		return fmt.Errorf("sound '%s' not loaded", name)
	}

	// Note: Ebiten doesn't support pitch shifting directly
	// This would require resampling the audio data
	// For now, just play at normal pitch
	_ = pitchFactor // Suppress unused warning
	return ap.PlaySound(name)
}

// GetMusicIntensity returns current music intensity
func (ap *AudioPlayer) GetMusicIntensity() MusicIntensity {
	if ap.adaptiveTrack == nil {
		return IntensityCalm
	}
	return ap.adaptiveTrack.CurrentIntensity
}
