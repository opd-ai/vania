package audio

import (
	"testing"

	ebaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

// ── AdaptiveMusicTrack ────────────────────────────────────────────────────────

func TestAdaptiveMusicTrackSetIntensity(t *testing.T) {
	track := NewAdaptiveMusicTrack()
	intensities := []MusicIntensity{IntensityCalm, IntensityTension, IntensityCombat, IntensityBoss}
	for _, intensity := range intensities {
		track.SetIntensity(intensity)
		if track.TargetIntensity != intensity {
			t.Errorf("SetIntensity(%v) not applied; TargetIntensity = %v", intensity, track.TargetIntensity)
		}
	}
}

func TestAdaptiveMusicTrackGetCurrentMixEmpty(t *testing.T) {
	track := NewAdaptiveMusicTrack()
	mix := track.GetCurrentMix()
	if mix == nil {
		t.Error("GetCurrentMix returned nil for empty track")
	}
}

func TestAdaptiveMusicTrackWithLayers(t *testing.T) {
	track := NewAdaptiveMusicTrack()
	layer := &MusicLayer{
		Name:         "test",
		MinIntensity: IntensityCalm,
		Audio:        &AudioSample{Data: []float64{0.1, 0.2}, SampleRate: 22050},
		BaseVolume:   0.8,
	}
	track.AddLayer(layer)
	track.SetIntensity(IntensityCombat)
	track.Update()

	mix := track.GetCurrentMix()
	if _, ok := mix["test"]; !ok {
		t.Error("Expected test layer in current mix in mix after Update")
	}
}

func TestAdaptiveMusicTrackAllIntensityTransitions(t *testing.T) {
	track := NewAdaptiveMusicTrack()
	layer := &MusicLayer{
		Name:         "bass",
		MinIntensity: IntensityCalm,
		Audio:        &AudioSample{Data: []float64{0.5}, SampleRate: 22050},
		BaseVolume:   0.7,
	}
	track.AddLayer(layer)

	transitions := []MusicIntensity{
		IntensityCalm, IntensityTension, IntensityCombat, IntensityBoss,
		IntensityCombat, IntensityTension, IntensityCalm,
	}
	for _, intensity := range transitions {
		track.SetIntensity(intensity)
		track.Update()
		mix := track.GetCurrentMix()
		if mix == nil {
			t.Errorf("nil mix after SetIntensity(%v)", intensity)
		}
	}
}

// ── MusicGenerator ────────────────────────────────────────────────────────────

func TestNewMusicGeneratorDefaults(t *testing.T) {
	mg := NewMusicGenerator(44100, 120, 60, MinorScale)
	if mg == nil {
		t.Fatal("Expected non-nil MusicGenerator")
	}
}

func TestMusicGeneratorGenerateTrack(t *testing.T) {
	mg := NewMusicGenerator(22050, 120, 60, MajorScale)
	track := mg.GenerateTrack(42, 1.0)
	if track == nil {
		t.Fatal("Expected non-nil AudioSample from GenerateTrack")
	}
	if len(track.Data) == 0 {
		t.Error("Expected non-empty samples")
	}
}

func TestMusicGeneratorDeterminism(t *testing.T) {
	mg1 := NewMusicGenerator(22050, 120, 60, MinorScale)
	mg2 := NewMusicGenerator(22050, 120, 60, MinorScale)
	seed := int64(999)

	t1 := mg1.GenerateTrack(seed, 0.5)
	t2 := mg2.GenerateTrack(seed, 0.5)

	if len(t1.Data) != len(t2.Data) {
		t.Errorf("Determinism failed: lengths %d vs %d", len(t1.Data), len(t2.Data))
	}
	if len(t1.Data) > 0 && t1.Data[0] != t2.Data[0] {
		t.Errorf("Determinism failed: first samples differ %f vs %f", t1.Data[0], t2.Data[0])
	}
}

func TestMusicGeneratorAllScales(t *testing.T) {
	scales := []Scale{
		MajorScale, MinorScale, DorianScale, PhrygianScale,
		PentatonicMaj, PentatonicMin,
	}
	for _, scale := range scales {
		mg := NewMusicGenerator(22050, 120, 60, scale)
		track := mg.GenerateTrack(1, 0.5)
		if track == nil || len(track.Data) == 0 {
			t.Errorf("GenerateTrack returned empty for scale %v", scale)
		}
	}
}

func TestMusicGeneratorGenerateAdaptiveTrack(t *testing.T) {
	mg := NewMusicGenerator(22050, 120, 60, MinorScale)
	amt := mg.GenerateAdaptiveMusicTrack(42, 0.5)
	if amt == nil {
		t.Fatal("Expected non-nil AdaptiveMusicTrack")
	}
	if len(amt.Layers) == 0 {
		t.Error("Expected at least one layer in adaptive track")
	}
}

func TestMusicGeneratorVariousBPMs(t *testing.T) {
	bpms := []int{60, 90, 120, 150, 180}
	for _, bpm := range bpms {
		mg := NewMusicGenerator(22050, bpm, 60, MajorScale)
		track := mg.GenerateTrack(1, 0.25)
		if track == nil {
			t.Errorf("GenerateTrack returned nil for BPM=%d", bpm)
		}
	}
}

// ── SFXGenerator extras ───────────────────────────────────────────────────────

func TestSFXGeneratorGenerateExplosion(t *testing.T) {
	gen := NewSFXGenerator(22050)
	sample := gen.GenerateExplosion(42)
	if sample == nil {
		t.Fatal("Expected non-nil explosion sample")
	}
	if len(sample.Data) == 0 {
		t.Error("Expected non-empty explosion samples")
	}
}

func TestSFXGeneratorApplyDistortion(t *testing.T) {
	gen := NewSFXGenerator(22050)
	input := &AudioSample{Data: []float64{0.1, 0.5, -0.5, 1.0, -1.0}}
	gen.ApplyDistortion(input, 2.0)
	for i, s := range input.Data {
		if s > 1.0 || s < -1.0 {
			t.Errorf("Sample[%d] = %f out of [-1,1] after distortion", i, s)
		}
	}
}

// ── AudioPlayer pure-logic tests (no hardware required) ──────────────────────

func makeTestPlayer() *AudioPlayer {
	return &AudioPlayer{
		players:          make(map[string]*ebaudio.Player),
		masterVolume:     0.7,
		sfxVolume:        0.8,
		musicVolume:      0.6,
		currentGenre:     "fantasy",
		genreInstruments: make(map[WaveType]float64),
	}
}

func TestAudioPlayerSetVolumes(t *testing.T) {
	ap := makeTestPlayer()
	ap.SetVolumes(0.5, 0.6, 0.4)
	m, s, mu := ap.GetVolumes()
	if m != 0.5 || s != 0.6 || mu != 0.4 {
		t.Errorf("GetVolumes() = %v %v %v; want 0.5 0.6 0.4", m, s, mu)
	}
}

func TestAudioPlayerSetVolumesClamped(t *testing.T) {
	ap := makeTestPlayer()
	ap.SetVolumes(2.0, -1.0, 0.5)
	m, s, mu := ap.GetVolumes()
	if m != 1.0 || s != 0.0 || mu != 0.5 {
		t.Errorf("Clamping failed: %v %v %v", m, s, mu)
	}
}

func TestAudioPlayerGetMusicIntensityDefault(t *testing.T) {
	ap := makeTestPlayer()
	if ap.GetMusicIntensity() != IntensityCalm {
		t.Error("Expected IntensityCalm when no adaptive track set")
	}
}

func TestAudioPlayerGetMusicIntensityWithTrack(t *testing.T) {
	ap := makeTestPlayer()
	ap.adaptiveTrack = NewAdaptiveMusicTrack()
	ap.adaptiveTrack.SetIntensity(IntensityCombat)
	ap.adaptiveTrack.Update()
	intensity := ap.GetMusicIntensity()
	_ = intensity
}

func TestAudioPlayerIsPlayingMissing(t *testing.T) {
	ap := makeTestPlayer()
	if ap.IsPlaying("nonexistent") {
		t.Error("Expected false for nonexistent sound")
	}
}

func TestAudioPlayerIsMusicPlayingNil(t *testing.T) {
	ap := makeTestPlayer()
	if ap.IsMusicPlaying() {
		t.Error("Expected false when musicPlayer is nil")
	}
}

func TestAudioPlayerSetGenreAllGenres(t *testing.T) {
	genres := []string{"fantasy", "scifi", "horror", "cyberpunk", "postapoc", "unknown"}
	for _, g := range genres {
		ap := makeTestPlayer()
		ap.SetGenre(g)
		if g != "unknown" && ap.currentGenre != g {
			t.Errorf("SetGenre(%q) set currentGenre to %q", g, ap.currentGenre)
		}
		if g != "unknown" && len(ap.GetGenreInstruments()) == 0 {
			t.Errorf("Expected non-empty instruments for genre %q", g)
		}
	}
}

func TestAudioPlayerGetGenreSFXVariation(t *testing.T) {
	ap := makeTestPlayer()
	ap.SetGenre("horror")
	v := ap.GetGenreSFXVariation()
	if v <= 0 {
		t.Errorf("Expected positive SFX variation, got %f", v)
	}
}

func TestAudioPlayerAudioSampleToWAVEmpty(t *testing.T) {
	ap := makeTestPlayer()
	_, err := ap.audioSampleToWAV(nil)
	if err == nil {
		t.Error("Expected error for nil AudioSample")
	}
	_, err = ap.audioSampleToWAV(&AudioSample{Data: []float64{}})
	if err == nil {
		t.Error("Expected error for empty AudioSample")
	}
}

func TestAudioPlayerAudioSampleToWAVValid(t *testing.T) {
	ap := makeTestPlayer()
	sample := &AudioSample{
		Data:       []float64{0.0, 0.5, -0.5, 1.0, -1.0},
		SampleRate: 22050,
		Duration:   0.1,
	}
	wavBytes, err := ap.audioSampleToWAV(sample)
	if err != nil {
		t.Fatalf("audioSampleToWAV failed: %v", err)
	}
	if len(wavBytes) < 44 {
		t.Errorf("Expected at least 44 bytes (WAV header), got %d", len(wavBytes))
	}
	if string(wavBytes[0:4]) != "RIFF" {
		t.Errorf("Expected RIFF header, got %q", wavBytes[0:4])
	}
}

// ── SFX type coverage ─────────────────────────────────────────────────────────

func TestSFXGeneratorAllTypes(t *testing.T) {
	gen := NewSFXGenerator(22050)
	sfxTypes := []SFXType{
		JumpSFX, LandSFX, AttackSFX, HitSFX,
		PickupSFX, DoorSFX, DamageSFX,
	}
	for _, sfxType := range sfxTypes {
		sample := gen.Generate(sfxType, 1)
		if sample == nil || len(sample.Data) == 0 {
			t.Errorf("Generate(%v) returned empty sample", sfxType)
		}
	}
}

// ── Additional coverage tests ─────────────────────────────────────────────────

func TestNewMusicGeneratorInvalidBPM(t *testing.T) {
	// BPM <= 0 should default to 120
	mg := NewMusicGenerator(22050, 0, 60, MajorScale)
	if mg.BPM != 120 {
		t.Errorf("Expected BPM=120 for invalid input, got %d", mg.BPM)
	}
}

func TestNewMusicGeneratorNilScale(t *testing.T) {
	mg := NewMusicGenerator(22050, 120, 60, nil)
	if len(mg.Scale) == 0 {
		t.Error("Expected default scale when nil passed")
	}
}

func TestNewSynthesizerValidation(t *testing.T) {
	s1 := NewSynthesizer(0) // invalid
	if s1.SampleRate <= 0 {
		t.Error("Expected positive sample rate as fallback")
	}
	s2 := NewSynthesizer(44100)
	if s2.SampleRate != 44100 {
		t.Error("Expected sample rate 44100")
	}
}

func TestAdaptiveMusicTrackUpdateNoLayers(t *testing.T) {
	// Update with no layers should not panic
	track := NewAdaptiveMusicTrack()
	track.SetIntensity(IntensityBoss)
	track.Update()
}

func TestAdaptiveMusicTrackGetCurrentMixKeys(t *testing.T) {
	track := NewAdaptiveMusicTrack()
	for _, name := range []string{"a", "b", "c"} {
		track.AddLayer(&MusicLayer{
			Name:         name,
			MinIntensity: IntensityCalm,
			Audio:        &AudioSample{Data: []float64{0.1}, SampleRate: 22050},
			BaseVolume:   0.5,
		})
	}
	mix := track.GetCurrentMix()
	if len(mix) != 3 {
		t.Errorf("Expected 3 layers in mix, got %d", len(mix))
	}
}

func TestAudioPlayerUpdateMusicNilTrack(t *testing.T) {
	ap := makeTestPlayer()
	// Should not panic
	ap.UpdateMusic(&MusicContext{
		InCombat:        false,
		PlayerHealthPct: 1.0,
	})
}

func TestAudioPlayerUpdateMusicWithTrack(t *testing.T) {
	ap := makeTestPlayer()
	ap.adaptiveTrack = NewAdaptiveMusicTrack()
	// musicPlayer is nil, so UpdateMusic returns early
	ap.UpdateMusic(&MusicContext{
		InCombat:        true,
		PlayerHealthPct: 0.5,
	})
}

func TestAudioPlayerCreateInfiniteLoopReaderNil(t *testing.T) {
	ap := makeTestPlayer()
	_, err := ap.CreateInfiniteLoopReader(nil)
	if err == nil {
		t.Error("Expected error for nil sample")
	}
}

func TestAudioPlayerPlaySoundWithPitchMissing(t *testing.T) {
	ap := makeTestPlayer()
	err := ap.PlaySoundWithPitch("nonexistent", 1.5)
	if err == nil {
		t.Error("Expected error for missing sound")
	}
}

func TestSynthesizerFrequencySweepEdge(t *testing.T) {
	s := NewSynthesizer(22050)
	// Test with very short duration
	result := s.FrequencySweep(SineWave, 440.0, 880.0, 0.01)
	if result == nil || len(result.Data) == 0 {
		t.Error("Expected non-empty FrequencySweep result")
	}
}

func TestMusicGeneratorConcatenateSamplesEmpty(t *testing.T) {
	mg := NewMusicGenerator(22050, 120, 60, MajorScale)
	result := mg.concatenateSamples([]*AudioSample{})
	if result == nil {
		t.Error("Expected non-nil result for empty input")
	}
	if len(result.Data) != 0 {
		t.Errorf("Expected empty data, got %d samples", len(result.Data))
	}
}

func TestAudioPlayerCloseEmpty(t *testing.T) {
	// Close with empty players should not panic
	ap := makeTestPlayer()
	ap.Close()
	if len(ap.players) != 0 {
		t.Error("Expected empty players after Close")
	}
	if ap.musicPlayer != nil {
		t.Error("Expected nil musicPlayer after Close")
	}
}

func TestAudioPlayerLoadSoundNil(t *testing.T) {
	ap := makeTestPlayer()
	err := ap.LoadSound("test", nil)
	if err == nil {
		t.Error("Expected error for nil AudioSample in LoadSound")
	}
}

func TestMusicGeneratorTrackDuration(t *testing.T) {
	mg := NewMusicGenerator(22050, 120, 60, MajorScale)
	// Test with different durations
	for _, dur := range []float64{0.1, 0.5, 2.0} {
		track := mg.GenerateTrack(1, dur)
		if track == nil {
			t.Errorf("nil track for duration %v", dur)
		}
	}
}
