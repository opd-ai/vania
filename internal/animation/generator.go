package animation

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/opd-ai/vania/internal/graphics"
)

// AnimationGenerator creates animation frames from base sprites
type AnimationGenerator struct {
	rng *rand.Rand
}

// NewAnimationGenerator creates a new animation generator
func NewAnimationGenerator(seed int64) *AnimationGenerator {
	return &AnimationGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// GenerateWalkFrames creates walking animation frames
func (ag *AnimationGenerator) GenerateWalkFrames(baseSprite *graphics.Sprite, numFrames int) []*graphics.Sprite {
	if baseSprite == nil || numFrames <= 0 {
		return nil
	}
	
	frames := make([]*graphics.Sprite, numFrames)
	
	for i := 0; i < numFrames; i++ {
		frames[i] = ag.createWalkFrame(baseSprite, i, numFrames)
	}
	
	return frames
}

// GenerateAttackFrames creates attack animation frames
func (ag *AnimationGenerator) GenerateAttackFrames(baseSprite *graphics.Sprite, numFrames int) []*graphics.Sprite {
	if baseSprite == nil || numFrames <= 0 {
		return nil
	}
	
	frames := make([]*graphics.Sprite, numFrames)
	
	for i := 0; i < numFrames; i++ {
		frames[i] = ag.createAttackFrame(baseSprite, i, numFrames)
	}
	
	return frames
}

// GenerateJumpFrames creates jump animation frames
func (ag *AnimationGenerator) GenerateJumpFrames(baseSprite *graphics.Sprite, numFrames int) []*graphics.Sprite {
	if baseSprite == nil || numFrames <= 0 {
		return nil
	}
	
	frames := make([]*graphics.Sprite, numFrames)
	
	for i := 0; i < numFrames; i++ {
		frames[i] = ag.createJumpFrame(baseSprite, i, numFrames)
	}
	
	return frames
}

// GenerateIdleFrames creates idle animation frames (subtle breathing effect)
func (ag *AnimationGenerator) GenerateIdleFrames(baseSprite *graphics.Sprite, numFrames int) []*graphics.Sprite {
	if baseSprite == nil || numFrames <= 0 {
		return nil
	}
	
	frames := make([]*graphics.Sprite, numFrames)
	
	for i := 0; i < numFrames; i++ {
		frames[i] = ag.createIdleFrame(baseSprite, i, numFrames)
	}
	
	return frames
}

// createWalkFrame creates a single walking frame with horizontal offset
func (ag *AnimationGenerator) createWalkFrame(baseSprite *graphics.Sprite, frameIndex, totalFrames int) *graphics.Sprite {
	newSprite := ag.copySprite(baseSprite)
	
	// Create bobbing effect (vertical offset)
	progress := float64(frameIndex) / float64(totalFrames)
	// Sin wave for smooth bobbing
	bobOffset := int(2.0 * (1.0 - progress*progress*4.0))
	
	// Shift sprite vertically
	if bobOffset != 0 {
		ag.shiftSpriteVertical(newSprite, bobOffset)
	}
	
	return newSprite
}

// createAttackFrame creates a single attack frame with forward lean
func (ag *AnimationGenerator) createAttackFrame(baseSprite *graphics.Sprite, frameIndex, totalFrames int) *graphics.Sprite {
	newSprite := ag.copySprite(baseSprite)
	
	// Progressive forward shift
	progress := float64(frameIndex) / float64(totalFrames)
	forwardShift := int(progress * 4.0)
	
	if forwardShift > 0 {
		ag.shiftSpriteHorizontal(newSprite, forwardShift)
	}
	
	return newSprite
}

// createJumpFrame creates a single jump frame
func (ag *AnimationGenerator) createJumpFrame(baseSprite *graphics.Sprite, frameIndex, totalFrames int) *graphics.Sprite {
	newSprite := ag.copySprite(baseSprite)
	
	// Compress/stretch effect
	progress := float64(frameIndex) / float64(totalFrames)
	
	// Slight vertical stretch at beginning and end (preparing to jump/landing)
	// Horizontal stretch in middle (mid-air)
	if progress < 0.3 || progress > 0.7 {
		// Crouch/landing - slight vertical compression
		ag.shiftSpriteVertical(newSprite, 1)
	}
	
	return newSprite
}

// createIdleFrame creates a single idle frame with subtle breathing
func (ag *AnimationGenerator) createIdleFrame(baseSprite *graphics.Sprite, frameIndex, totalFrames int) *graphics.Sprite {
	// For idle, we use the base sprite mostly as-is
	// Just create a slight variation for breathing effect
	newSprite := ag.copySprite(baseSprite)
	
	// Very subtle vertical offset (breathing)
	progress := float64(frameIndex) / float64(totalFrames)
	breathOffset := 0
	if progress > 0.25 && progress < 0.75 {
		breathOffset = 1
	}
	
	if breathOffset != 0 {
		ag.shiftSpriteVertical(newSprite, breathOffset)
	}
	
	return newSprite
}

// copySprite creates a deep copy of a sprite
func (ag *AnimationGenerator) copySprite(sprite *graphics.Sprite) *graphics.Sprite {
	if sprite == nil || sprite.Image == nil {
		return &graphics.Sprite{
			Width:  sprite.Width,
			Height: sprite.Height,
		}
	}
	
	// Create new image with same bounds
	bounds := sprite.Image.Bounds()
	newImage := image.NewRGBA(bounds)
	
	// Copy pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImage.Set(x, y, sprite.Image.At(x, y))
		}
	}
	
	return &graphics.Sprite{
		Image:  newImage,
		Width:  sprite.Width,
		Height: sprite.Height,
	}
}

// shiftSpriteVertical shifts sprite content vertically
func (ag *AnimationGenerator) shiftSpriteVertical(sprite *graphics.Sprite, offset int) {
	if sprite == nil || sprite.Image == nil || offset == 0 {
		return
	}
	
	bounds := sprite.Image.Bounds()
	tempImage := image.NewRGBA(bounds)
	
	// Copy with offset
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newY := y + offset
			if newY >= bounds.Min.Y && newY < bounds.Max.Y {
				tempImage.Set(x, newY, sprite.Image.At(x, y))
			}
		}
	}
	
	sprite.Image = tempImage
}

// shiftSpriteHorizontal shifts sprite content horizontally
func (ag *AnimationGenerator) shiftSpriteHorizontal(sprite *graphics.Sprite, offset int) {
	if sprite == nil || sprite.Image == nil || offset == 0 {
		return
	}
	
	bounds := sprite.Image.Bounds()
	tempImage := image.NewRGBA(bounds)
	
	// Copy with offset
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newX := x + offset
			if newX >= bounds.Min.X && newX < bounds.Max.X {
				tempImage.Set(newX, y, sprite.Image.At(x, y))
			}
		}
	}
	
	sprite.Image = tempImage
}

// GenerateHitFrames creates hit/damage animation frames (flash effect)
func (ag *AnimationGenerator) GenerateHitFrames(baseSprite *graphics.Sprite, numFrames int) []*graphics.Sprite {
	if baseSprite == nil || numFrames <= 0 {
		return nil
	}
	
	frames := make([]*graphics.Sprite, numFrames)
	
	for i := 0; i < numFrames; i++ {
		frames[i] = ag.createHitFrame(baseSprite, i, numFrames)
	}
	
	return frames
}

// createHitFrame creates a hit frame with color tinting
func (ag *AnimationGenerator) createHitFrame(baseSprite *graphics.Sprite, frameIndex, totalFrames int) *graphics.Sprite {
	newSprite := ag.copySprite(baseSprite)
	
	if newSprite.Image == nil {
		return newSprite
	}
	
	// Flash red on odd frames
	if frameIndex%2 == 1 {
		ag.tintSprite(newSprite, color.RGBA{255, 100, 100, 255})
	}
	
	return newSprite
}

// tintSprite applies a color tint to sprite
func (ag *AnimationGenerator) tintSprite(sprite *graphics.Sprite, tint color.RGBA) {
	if sprite == nil || sprite.Image == nil {
		return
	}
	
	bounds := sprite.Image.Bounds()
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := sprite.Image.At(x, y)
			r, g, b, a := c.RGBA()
			
			// Skip transparent pixels
			if a == 0 {
				continue
			}
			
			// Blend with tint
			newR := uint8((uint32(r>>8) + uint32(tint.R)) / 2)
			newG := uint8((uint32(g>>8) + uint32(tint.G)) / 2)
			newB := uint8((uint32(b>>8) + uint32(tint.B)) / 2)
			
			sprite.Image.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}
}
