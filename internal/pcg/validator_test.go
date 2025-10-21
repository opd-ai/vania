package pcg

import (
	"testing"
)

// TestNewValidator_CreatesValidValidator tests validator creation
func TestNewValidator_CreatesValidValidator(t *testing.T) {
	minScore := 7.5
	validator := NewValidator(minScore)

	if validator == nil {
		t.Fatal("NewValidator returned nil")
	}

	if validator.minQualityScore != minScore {
		t.Errorf("Expected minQualityScore %f, got %f", minScore, validator.minQualityScore)
	}
}

// TestNewValidator_DifferentScores tests validator creation with various scores
func TestNewValidator_DifferentScores(t *testing.T) {
	tests := []struct {
		name     string
		minScore float64
	}{
		{"zero score", 0.0},
		{"low score", 3.5},
		{"mid score", 5.0},
		{"high score", 8.5},
		{"max score", 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.minScore)
			if validator.minQualityScore != tt.minScore {
				t.Errorf("Expected minQualityScore %f, got %f", tt.minScore, validator.minQualityScore)
			}
		})
	}
}

// TestValidateSprite_ReturnsTrue tests sprite validation
func TestValidateSprite_ReturnsTrue(t *testing.T) {
	validator := NewValidator(5.0)
	sprite := "test_sprite_data"

	result := validator.ValidateSprite(sprite)

	if !result {
		t.Error("ValidateSprite should return true")
	}
}

// TestValidateSprite_VariousInputs tests sprite validation with different inputs
func TestValidateSprite_VariousInputs(t *testing.T) {
	validator := NewValidator(5.0)

	tests := []struct {
		name   string
		sprite interface{}
	}{
		{"string sprite", "sprite_data"},
		{"int sprite", 12345},
		{"nil sprite", nil},
		{"map sprite", map[string]interface{}{"data": "value"}},
		{"slice sprite", []byte{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateSprite(tt.sprite)
			if !result {
				t.Errorf("ValidateSprite returned false for %s", tt.name)
			}
		})
	}
}

// TestValidateAudio_ReturnsTrue tests audio validation
func TestValidateAudio_ReturnsTrue(t *testing.T) {
	validator := NewValidator(5.0)
	audio := "test_audio_data"

	result := validator.ValidateAudio(audio)

	if !result {
		t.Error("ValidateAudio should return true")
	}
}

// TestValidateAudio_VariousInputs tests audio validation with different inputs
func TestValidateAudio_VariousInputs(t *testing.T) {
	validator := NewValidator(5.0)

	tests := []struct {
		name  string
		audio interface{}
	}{
		{"string audio", "audio_data"},
		{"int audio", 98765},
		{"nil audio", nil},
		{"struct audio", struct{ data string }{data: "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateAudio(tt.audio)
			if !result {
				t.Errorf("ValidateAudio returned false for %s", tt.name)
			}
		})
	}
}

// TestValidateNarrative_ReturnsTrue tests narrative validation
func TestValidateNarrative_ReturnsTrue(t *testing.T) {
	validator := NewValidator(5.0)
	narrative := "test_narrative_data"

	result := validator.ValidateNarrative(narrative)

	if !result {
		t.Error("ValidateNarrative should return true")
	}
}

// TestValidateNarrative_VariousInputs tests narrative validation with different inputs
func TestValidateNarrative_VariousInputs(t *testing.T) {
	validator := NewValidator(5.0)

	tests := []struct {
		name      string
		narrative interface{}
	}{
		{"string narrative", "narrative_data"},
		{"complex narrative", map[string]string{"story": "A tale of heroes"}},
		{"nil narrative", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateNarrative(tt.narrative)
			if !result {
				t.Errorf("ValidateNarrative returned false for %s", tt.name)
			}
		})
	}
}

// TestValidateWorld_ReturnsTrue tests world validation
func TestValidateWorld_ReturnsTrue(t *testing.T) {
	validator := NewValidator(5.0)
	world := "test_world_data"

	result := validator.ValidateWorld(world)

	if !result {
		t.Error("ValidateWorld should return true")
	}
}

// TestValidateWorld_VariousInputs tests world validation with different inputs
func TestValidateWorld_VariousInputs(t *testing.T) {
	validator := NewValidator(5.0)

	tests := []struct {
		name  string
		world interface{}
	}{
		{"string world", "world_data"},
		{"map world", map[string]interface{}{"rooms": 10}},
		{"nil world", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateWorld(tt.world)
			if !result {
				t.Errorf("ValidateWorld returned false for %s", tt.name)
			}
		})
	}
}

// TestCalculateQualityScore_BasicMetrics tests quality score calculation
func TestCalculateQualityScore_BasicMetrics(t *testing.T) {
	validator := NewValidator(5.0)

	metrics := &QualityMetrics{
		Completability:   0.8,  // 80% beatable
		GenerationTime:   100,
		VisualCoherence:  8.0,
		AudioHarmony:     7.0,
		NarrativeScore:   6.0,
		ContentDiversity: 0.9,
		PerformanceFPS:   60,
	}

	score := validator.CalculateQualityScore(metrics)

	// Expected: (8.0 * 0.25) + (7.0 * 0.20) + (6.0 * 0.20) + (0.8 * 10.0 * 0.35)
	// = 2.0 + 1.4 + 1.2 + 2.8 = 7.4
	expected := 7.4
	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}
}

// TestCalculateQualityScore_PerfectMetrics tests score with perfect metrics
func TestCalculateQualityScore_PerfectMetrics(t *testing.T) {
	validator := NewValidator(5.0)

	metrics := &QualityMetrics{
		Completability:   1.0,   // 100% beatable
		GenerationTime:   50,
		VisualCoherence:  10.0,
		AudioHarmony:     10.0,
		NarrativeScore:   10.0,
		ContentDiversity: 1.0,
		PerformanceFPS:   120,
	}

	score := validator.CalculateQualityScore(metrics)

	// Expected: (10.0 * 0.25) + (10.0 * 0.20) + (10.0 * 0.20) + (1.0 * 10.0 * 0.35)
	// = 2.5 + 2.0 + 2.0 + 3.5 = 10.0
	expected := 10.0
	if score != expected {
		t.Errorf("Expected perfect score %f, got %f", expected, score)
	}
}

// TestCalculateQualityScore_ZeroMetrics tests score with zero metrics
func TestCalculateQualityScore_ZeroMetrics(t *testing.T) {
	validator := NewValidator(5.0)

	metrics := &QualityMetrics{
		Completability:   0.0,
		GenerationTime:   0,
		VisualCoherence:  0.0,
		AudioHarmony:     0.0,
		NarrativeScore:   0.0,
		ContentDiversity: 0.0,
		PerformanceFPS:   0,
	}

	score := validator.CalculateQualityScore(metrics)

	expected := 0.0
	if score != expected {
		t.Errorf("Expected zero score %f, got %f", expected, score)
	}
}

// TestCalculateQualityScore_TableDriven tests various metric combinations
func TestCalculateQualityScore_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		metrics  *QualityMetrics
		expected float64
	}{
		{
			name: "low visual high completability",
			metrics: &QualityMetrics{
				Completability:  1.0,
				VisualCoherence: 2.0,
				AudioHarmony:    5.0,
				NarrativeScore:  5.0,
			},
			expected: (2.0 * 0.25) + (5.0 * 0.20) + (5.0 * 0.20) + (1.0 * 10.0 * 0.35),
		},
		{
			name: "high visual low completability",
			metrics: &QualityMetrics{
				Completability:  0.2,
				VisualCoherence: 10.0,
				AudioHarmony:    10.0,
				NarrativeScore:  10.0,
			},
			expected: (10.0 * 0.25) + (10.0 * 0.20) + (10.0 * 0.20) + (0.2 * 10.0 * 0.35),
		},
		{
			name: "mixed metrics",
			metrics: &QualityMetrics{
				Completability:  0.5,
				VisualCoherence: 6.0,
				AudioHarmony:    7.5,
				NarrativeScore:  8.0,
			},
			expected: (6.0 * 0.25) + (7.5 * 0.20) + (8.0 * 0.20) + (0.5 * 10.0 * 0.35),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(5.0)
			score := validator.CalculateQualityScore(tt.metrics)

			if score != tt.expected {
				t.Errorf("Expected score %f, got %f", tt.expected, score)
			}
		})
	}
}

// TestMeetsThreshold_PassingScore tests threshold passing
func TestMeetsThreshold_PassingScore(t *testing.T) {
	validator := NewValidator(5.0)

	metrics := &QualityMetrics{
		Completability:   0.8,
		VisualCoherence:  8.0,
		AudioHarmony:     7.0,
		NarrativeScore:   6.0,
	}

	result := validator.MeetsThreshold(metrics)

	if !result {
		t.Error("MeetsThreshold should return true for passing score")
	}
}

// TestMeetsThreshold_FailingScore tests threshold failing
func TestMeetsThreshold_FailingScore(t *testing.T) {
	validator := NewValidator(8.0)

	metrics := &QualityMetrics{
		Completability:   0.2,
		VisualCoherence:  3.0,
		AudioHarmony:     2.0,
		NarrativeScore:   1.0,
	}

	result := validator.MeetsThreshold(metrics)

	if result {
		t.Error("MeetsThreshold should return false for failing score")
	}
}

// TestMeetsThreshold_ExactThreshold tests exact threshold match
func TestMeetsThreshold_ExactThreshold(t *testing.T) {
	threshold := 7.0
	validator := NewValidator(threshold)

	// Create metrics that give exactly 7.0
	// Score = (8 * 0.25) + (6 * 0.20) + (5 * 0.20) + (0.5 * 10 * 0.35)
	//       = 2.0 + 1.2 + 1.0 + 1.75 = 5.95
	// Let's recalculate for exactly 7.0
	// We need: visual*0.25 + audio*0.20 + narrative*0.20 + completability*10*0.35 = 7.0
	// Using: visual=10, audio=10, narrative=10, completability=0.4
	// = 2.5 + 2.0 + 2.0 + 1.4 = 7.9 (too high)
	// Using: visual=8, audio=8, narrative=8, completability=0.4
	// = 2.0 + 1.6 + 1.6 + 1.4 = 6.6 (too low)
	// Using: visual=9, audio=8, narrative=7, completability=0.4571
	// = 2.25 + 1.6 + 1.4 + 1.6 = 6.85 (close)
	// Using: visual=10, audio=8, narrative=8, completability=0.4
	// = 2.5 + 1.6 + 1.6 + 1.4 = 7.1 (close enough)
	// Let's use simple calc: if all are 7.0 and completability is 0.7
	// = 7*0.25 + 7*0.20 + 7*0.20 + 0.7*10*0.35 = 1.75 + 1.4 + 1.4 + 2.45 = 7.0

	metrics := &QualityMetrics{
		Completability:  0.7,
		VisualCoherence: 7.0,
		AudioHarmony:    7.0,
		NarrativeScore:  7.0,
	}

	result := validator.MeetsThreshold(metrics)

	if !result {
		t.Error("MeetsThreshold should return true for score equal to threshold")
	}
}

// TestMeetsThreshold_VariousThresholds tests different threshold values
func TestMeetsThreshold_VariousThresholds(t *testing.T) {
	tests := []struct {
		name           string
		threshold      float64
		metrics        *QualityMetrics
		shouldPass     bool
	}{
		{
			name:      "low threshold passes",
			threshold: 2.0,
			metrics: &QualityMetrics{
				Completability:  0.3,
				VisualCoherence: 5.0,
				AudioHarmony:    4.0,
				NarrativeScore:  3.0,
			},
			shouldPass: true,
		},
		{
			name:      "high threshold fails",
			threshold: 9.0,
			metrics: &QualityMetrics{
				Completability:  0.7,
				VisualCoherence: 8.0,
				AudioHarmony:    8.0,
				NarrativeScore:  8.0,
			},
			shouldPass: false,
		},
		{
			name:      "medium threshold passes",
			threshold: 6.0,
			metrics: &QualityMetrics{
				Completability:  0.8,
				VisualCoherence: 7.0,
				AudioHarmony:    7.0,
				NarrativeScore:  7.0,
			},
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.threshold)
			result := validator.MeetsThreshold(tt.metrics)

			if result != tt.shouldPass {
				score := validator.CalculateQualityScore(tt.metrics)
				t.Errorf("Expected %v, got %v (score: %f, threshold: %f)",
					tt.shouldPass, result, score, tt.threshold)
			}
		})
	}
}

// TestQualityMetrics_AllFields tests QualityMetrics structure
func TestQualityMetrics_AllFields(t *testing.T) {
	metrics := &QualityMetrics{
		Completability:   0.95,
		GenerationTime:   150,
		VisualCoherence:  8.5,
		AudioHarmony:     9.0,
		NarrativeScore:   7.5,
		ContentDiversity: 0.85,
		PerformanceFPS:   60,
	}

	if metrics.Completability != 0.95 {
		t.Errorf("Completability mismatch")
	}
	if metrics.GenerationTime != 150 {
		t.Errorf("GenerationTime mismatch")
	}
	if metrics.VisualCoherence != 8.5 {
		t.Errorf("VisualCoherence mismatch")
	}
	if metrics.AudioHarmony != 9.0 {
		t.Errorf("AudioHarmony mismatch")
	}
	if metrics.NarrativeScore != 7.5 {
		t.Errorf("NarrativeScore mismatch")
	}
	if metrics.ContentDiversity != 0.85 {
		t.Errorf("ContentDiversity mismatch")
	}
	if metrics.PerformanceFPS != 60 {
		t.Errorf("PerformanceFPS mismatch")
	}
}

// TestValidator_NegativeMinScore tests validator with negative threshold
func TestValidator_NegativeMinScore(t *testing.T) {
	validator := NewValidator(-5.0)

	metrics := &QualityMetrics{
		Completability:  0.0,
		VisualCoherence: 0.0,
		AudioHarmony:    0.0,
		NarrativeScore:  0.0,
	}

	// Even zero score should pass negative threshold
	result := validator.MeetsThreshold(metrics)
	if !result {
		t.Error("Negative threshold should allow zero score to pass")
	}
}

// TestValidator_VeryHighMinScore tests validator with very high threshold
func TestValidator_VeryHighMinScore(t *testing.T) {
	validator := NewValidator(15.0)

	metrics := &QualityMetrics{
		Completability:  1.0,
		VisualCoherence: 10.0,
		AudioHarmony:    10.0,
		NarrativeScore:  10.0,
	}

	// Even perfect score (10.0) should fail threshold of 15.0
	result := validator.MeetsThreshold(metrics)
	if result {
		t.Error("Perfect score should fail threshold above 10.0")
	}
}

// TestCalculateQualityScore_WeightedCorrectly tests that weights are applied
func TestCalculateQualityScore_WeightedCorrectly(t *testing.T) {
	validator := NewValidator(5.0)

	// Test that visual weight (0.25) is applied
	visualOnlyMetrics := &QualityMetrics{
		VisualCoherence: 10.0,
		AudioHarmony:    0.0,
		NarrativeScore:  0.0,
		Completability:  0.0,
	}
	visualScore := validator.CalculateQualityScore(visualOnlyMetrics)
	expectedVisual := 10.0 * 0.25
	if visualScore != expectedVisual {
		t.Errorf("Visual weight incorrect: expected %f, got %f", expectedVisual, visualScore)
	}

	// Test that audio weight (0.20) is applied
	audioOnlyMetrics := &QualityMetrics{
		VisualCoherence: 0.0,
		AudioHarmony:    10.0,
		NarrativeScore:  0.0,
		Completability:  0.0,
	}
	audioScore := validator.CalculateQualityScore(audioOnlyMetrics)
	expectedAudio := 10.0 * 0.20
	if audioScore != expectedAudio {
		t.Errorf("Audio weight incorrect: expected %f, got %f", expectedAudio, audioScore)
	}

	// Test that narrative weight (0.20) is applied
	narrativeOnlyMetrics := &QualityMetrics{
		VisualCoherence: 0.0,
		AudioHarmony:    0.0,
		NarrativeScore:  10.0,
		Completability:  0.0,
	}
	narrativeScore := validator.CalculateQualityScore(narrativeOnlyMetrics)
	expectedNarrative := 10.0 * 0.20
	if narrativeScore != expectedNarrative {
		t.Errorf("Narrative weight incorrect: expected %f, got %f", expectedNarrative, narrativeScore)
	}

	// Test that playable weight (0.35) is applied
	playableOnlyMetrics := &QualityMetrics{
		VisualCoherence: 0.0,
		AudioHarmony:    0.0,
		NarrativeScore:  0.0,
		Completability:  1.0,
	}
	playableScore := validator.CalculateQualityScore(playableOnlyMetrics)
	expectedPlayable := 1.0 * 10.0 * 0.35
	if playableScore != expectedPlayable {
		t.Errorf("Playable weight incorrect: expected %f, got %f", expectedPlayable, playableScore)
	}
}

// TestValidator_MultipleValidations tests validator reuse
func TestValidator_MultipleValidations(t *testing.T) {
	validator := NewValidator(6.0)

	metrics1 := &QualityMetrics{
		Completability:  0.9,
		VisualCoherence: 8.0,
		AudioHarmony:    8.0,
		NarrativeScore:  8.0,
	}

	metrics2 := &QualityMetrics{
		Completability:  0.3,
		VisualCoherence: 3.0,
		AudioHarmony:    3.0,
		NarrativeScore:  3.0,
	}

	// First validation should pass
	if !validator.MeetsThreshold(metrics1) {
		t.Error("First validation should pass")
	}

	// Second validation should fail
	if validator.MeetsThreshold(metrics2) {
		t.Error("Second validation should fail")
	}

	// Validator should remain in valid state
	if !validator.MeetsThreshold(metrics1) {
		t.Error("Validator state corrupted after multiple uses")
	}
}
