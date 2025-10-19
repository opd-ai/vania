package audio

import (
	"testing"
)

func TestMusicIntensityLevels(t *testing.T) {
	tests := []struct {
		name     string
		context  *MusicContext
		expected MusicIntensity
	}{
		{
			name: "Calm exploration",
			context: &MusicContext{
				InCombat:         false,
				IsBossFight:      false,
				NearbyEnemyCount: 0,
				PlayerHealthPct:  1.0,
				RoomDangerLevel:  1,
			},
			expected: IntensityCalm,
		},
		{
			name: "Tension from nearby enemies",
			context: &MusicContext{
				InCombat:         false,
				IsBossFight:      false,
				NearbyEnemyCount: 2,
				PlayerHealthPct:  1.0,
				RoomDangerLevel:  5,
			},
			expected: IntensityTension,
		},
		{
			name: "Tension from low health",
			context: &MusicContext{
				InCombat:         false,
				IsBossFight:      false,
				NearbyEnemyCount: 0,
				PlayerHealthPct:  0.2,
				RoomDangerLevel:  3,
			},
			expected: IntensityTension,
		},
		{
			name: "Tension from high danger room",
			context: &MusicContext{
				InCombat:         false,
				IsBossFight:      false,
				NearbyEnemyCount: 0,
				PlayerHealthPct:  1.0,
				RoomDangerLevel:  8,
			},
			expected: IntensityTension,
		},
		{
			name: "Combat intensity",
			context: &MusicContext{
				InCombat:         true,
				IsBossFight:      false,
				NearbyEnemyCount: 3,
				PlayerHealthPct:  0.8,
				RoomDangerLevel:  5,
			},
			expected: IntensityCombat,
		},
		{
			name: "Boss fight intensity",
			context: &MusicContext{
				InCombat:         true,
				IsBossFight:      true,
				NearbyEnemyCount: 1,
				PlayerHealthPct:  0.5,
				RoomDangerLevel:  10,
			},
			expected: IntensityBoss,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.context.CalculateIntensity()
			if result != tt.expected {
				t.Errorf("Expected intensity %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestAdaptiveMusicTrack(t *testing.T) {
	track := NewAdaptiveMusicTrack()

	if track == nil {
		t.Fatal("Expected non-nil track")
	}

	if track.CurrentIntensity != IntensityCalm {
		t.Errorf("Expected initial intensity to be Calm, got %d", track.CurrentIntensity)
	}

	if track.CurrentMix == nil {
		t.Error("Expected non-nil CurrentMix")
	}
}

func TestMusicLayerVolume(t *testing.T) {
	track := NewAdaptiveMusicTrack()

	// Add a layer that requires tension
	layer := &MusicLayer{
		Name:         "test",
		Audio:        &AudioSample{},
		BaseVolume:   0.5,
		MinIntensity: IntensityTension,
	}
	track.AddLayer(layer)

	// At calm intensity, layer should be silent
	track.CurrentIntensity = IntensityCalm
	volume := track.calculateLayerVolume(layer)
	if volume != 0.0 {
		t.Errorf("Expected volume 0.0 for layer below min intensity, got %f", volume)
	}

	// At tension intensity, layer should have some volume
	track.CurrentIntensity = IntensityTension
	volume = track.calculateLayerVolume(layer)
	if volume <= 0.0 {
		t.Errorf("Expected volume > 0 for layer at min intensity, got %f", volume)
	}

	// At boss intensity, layer should be at base volume
	track.CurrentIntensity = IntensityBoss
	volume = track.calculateLayerVolume(layer)
	if volume != layer.BaseVolume {
		t.Errorf("Expected volume %f at max intensity, got %f", layer.BaseVolume, volume)
	}
}

func TestMusicIntensityTransition(t *testing.T) {
	track := NewAdaptiveMusicTrack()
	track.CurrentIntensity = IntensityCalm
	track.TargetIntensity = IntensityCombat

	// Should increase by 1 each update
	track.Update()
	if track.CurrentIntensity != IntensityTension {
		t.Errorf("Expected intensity to increase to Tension, got %d", track.CurrentIntensity)
	}

	track.Update()
	if track.CurrentIntensity != IntensityCombat {
		t.Errorf("Expected intensity to increase to Combat, got %d", track.CurrentIntensity)
	}

	// Test decreasing intensity
	track.TargetIntensity = IntensityCalm
	track.Update()
	if track.CurrentIntensity != IntensityTension {
		t.Errorf("Expected intensity to decrease to Tension, got %d", track.CurrentIntensity)
	}
}

func TestGenerateAdaptiveMusicTrack(t *testing.T) {
	generator := NewMusicGenerator(44100, 120, 60, MinorScale)
	seed := int64(12345)
	duration := 4.0

	track := generator.GenerateAdaptiveMusicTrack(seed, duration)

	if track == nil {
		t.Fatal("Expected non-nil adaptive track")
	}

	// Should have 5 layers: pads, melody, bass, drums, lead
	if len(track.Layers) != 5 {
		t.Errorf("Expected 5 layers, got %d", len(track.Layers))
	}

	// Check that all layers have audio
	for _, layer := range track.Layers {
		if layer.Audio == nil {
			t.Errorf("Layer %s has nil audio", layer.Name)
		}
		if layer.Audio.Data == nil || len(layer.Audio.Data) == 0 {
			t.Errorf("Layer %s has no audio data", layer.Name)
		}
	}

	// Check layer names
	expectedLayers := map[string]bool{
		"pads":   false,
		"melody": false,
		"bass":   false,
		"drums":  false,
		"lead":   false,
	}

	for _, layer := range track.Layers {
		if _, exists := expectedLayers[layer.Name]; !exists {
			t.Errorf("Unexpected layer name: %s", layer.Name)
		}
		expectedLayers[layer.Name] = true
	}

	for name, found := range expectedLayers {
		if !found {
			t.Errorf("Missing expected layer: %s", name)
		}
	}
}

func TestLayerMinIntensity(t *testing.T) {
	generator := NewMusicGenerator(44100, 120, 60, MinorScale)
	track := generator.GenerateAdaptiveMusicTrack(12345, 4.0)

	// Verify each layer has correct minimum intensity
	layerIntensities := map[string]MusicIntensity{
		"pads":   IntensityCalm,
		"melody": IntensityCalm,
		"bass":   IntensityTension,
		"drums":  IntensityCombat,
		"lead":   IntensityBoss,
	}

	for _, layer := range track.Layers {
		expected, ok := layerIntensities[layer.Name]
		if !ok {
			t.Errorf("Unknown layer: %s", layer.Name)
			continue
		}
		if layer.MinIntensity != expected {
			t.Errorf("Layer %s: expected min intensity %d, got %d",
				layer.Name, expected, layer.MinIntensity)
		}
	}
}

func TestSmoothVolumeTransition(t *testing.T) {
	track := NewAdaptiveMusicTrack()

	layer := &MusicLayer{
		Name:         "test",
		Audio:        &AudioSample{},
		BaseVolume:   1.0,
		MinIntensity: IntensityCalm,
	}
	track.AddLayer(layer)

	track.CurrentIntensity = IntensityCalm
	track.TargetIntensity = IntensityCalm

	// Initial volume should be 0
	if track.CurrentMix["test"] != 0.0 {
		t.Errorf("Expected initial volume 0, got %f", track.CurrentMix["test"])
	}

	// After several updates, volume should gradually increase
	previousVolume := 0.0
	for i := 0; i < 100; i++ {
		track.Update()
		currentVolume := track.CurrentMix["test"]
		if currentVolume < previousVolume {
			t.Errorf("Volume decreased at update %d: %f -> %f", i, previousVolume, currentVolume)
		}
		previousVolume = currentVolume
	}

	// Eventually should reach near target volume (30% for calm at min intensity)
	expectedVolume := 0.3 // 30% of base volume at minimum intensity
	if track.CurrentMix["test"] < expectedVolume*0.8 {
		t.Errorf("Expected volume to reach near %f, got %f after 100 updates", expectedVolume, track.CurrentMix["test"])
	}
}

func TestDeterministicGeneration(t *testing.T) {
	seed := int64(99999)
	generator := NewMusicGenerator(44100, 120, 60, MinorScale)

	track1 := generator.GenerateAdaptiveMusicTrack(seed, 4.0)
	track2 := generator.GenerateAdaptiveMusicTrack(seed, 4.0)

	if len(track1.Layers) != len(track2.Layers) {
		t.Errorf("Layer counts differ: %d vs %d", len(track1.Layers), len(track2.Layers))
	}

	// Check that layers have same audio data length (deterministic generation)
	for i := 0; i < len(track1.Layers) && i < len(track2.Layers); i++ {
		layer1 := track1.Layers[i]
		layer2 := track2.Layers[i]

		if len(layer1.Audio.Data) != len(layer2.Audio.Data) {
			t.Errorf("Layer %d audio length differs: %d vs %d",
				i, len(layer1.Audio.Data), len(layer2.Audio.Data))
		}
	}
}

func TestMusicContextUpdate(t *testing.T) {
	context := NewMusicContext()

	if context == nil {
		t.Fatal("Expected non-nil context")
	}

	// Default should be calm
	if context.CalculateIntensity() != IntensityCalm {
		t.Error("Expected default intensity to be calm")
	}

	// Update context for combat
	context.InCombat = true
	if context.CalculateIntensity() != IntensityCombat {
		t.Error("Expected combat intensity when InCombat is true")
	}

	// Boss fight should override everything
	context.IsBossFight = true
	if context.CalculateIntensity() != IntensityBoss {
		t.Error("Expected boss intensity when IsBossFight is true")
	}
}
