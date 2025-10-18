// Package pcg provides core procedural content generation framework including
// seed management, caching, validation, and quality metrics for deterministic
// generation across all game subsystems.
package pcg

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

// PCGContext holds the seed and random number generator for procedural generation
type PCGContext struct {
	Seed        int64
	RNG         *rand.Rand
	Cache       *AssetCache
	Constraints *GenerationRules
}

// GenerationRules defines constraints for PCG quality
type GenerationRules struct {
	MinQualityScore     float64
	MaxGenerationTime   int64 // milliseconds
	EnableValidation    bool
	DifficultyPreset    string
	ThemeBias           string
	ArtisticStyle       string
	MusicalGenre        string
}

// NewPCGContext creates a new PCG context with the given seed
func NewPCGContext(seed int64) *PCGContext {
	return &PCGContext{
		Seed:  seed,
		RNG:   rand.New(rand.NewSource(seed)),
		Cache: NewAssetCache(),
		Constraints: &GenerationRules{
			MinQualityScore:   7.0,
			MaxGenerationTime: 30000,
			EnableValidation:  true,
			DifficultyPreset:  "normal",
			ThemeBias:         "fantasy",
			ArtisticStyle:     "retro",
			MusicalGenre:      "chiptune",
		},
	}
}

// HashSeed derives a subsystem seed from a master seed and identifier
func HashSeed(masterSeed int64, identifier string) int64 {
	h := sha256.New()
	// binary.Write to hash.Hash never returns an error, but we check for robustness
	if err := binary.Write(h, binary.LittleEndian, masterSeed); err != nil {
		// This should never happen with hash.Hash, but handle gracefully
		panic("binary.Write to hash failed: " + err.Error())
	}
	h.Write([]byte(identifier))
	sum := h.Sum(nil)
	return int64(binary.LittleEndian.Uint64(sum[:8]))
}

// DeriveSeeds generates all subsystem seeds from a master seed
func DeriveSeeds(masterSeed int64) map[string]int64 {
	return map[string]int64{
		"graphics":  HashSeed(masterSeed, "graphics"),
		"audio":     HashSeed(masterSeed, "audio"),
		"narrative": HashSeed(masterSeed, "narrative"),
		"world":     HashSeed(masterSeed, "world"),
		"entity":    HashSeed(masterSeed, "entity"),
	}
}

// NewDeterministicRNG creates a new deterministic RNG from a seed
func NewDeterministicRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}
