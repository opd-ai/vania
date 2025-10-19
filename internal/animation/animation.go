// Package animation provides sprite animation system for the game,
// supporting frame-based animations with configurable timing and looping.
package animation

import (
	"github.com/opd-ai/vania/internal/graphics"
)

// Animation represents a sequence of sprite frames
type Animation struct {
	Name       string
	Frames     []*graphics.Sprite
	FrameTime  int // frames per sprite frame (at 60 FPS)
	Loop       bool
	currentFrame int
	timer      int
}

// AnimationController manages animations for an entity
type AnimationController struct {
	animations      map[string]*Animation
	currentAnim     string
	defaultAnim     string
	playing         bool
}

// NewAnimation creates a new animation
func NewAnimation(name string, frames []*graphics.Sprite, frameTime int, loop bool) *Animation {
	if frameTime <= 0 {
		frameTime = 10 // Default to 10 frames per sprite frame
	}
	
	return &Animation{
		Name:       name,
		Frames:     frames,
		FrameTime:  frameTime,
		Loop:       loop,
		currentFrame: 0,
		timer:      0,
	}
}

// NewAnimationController creates a new animation controller
func NewAnimationController(defaultAnimation string) *AnimationController {
	return &AnimationController{
		animations:  make(map[string]*Animation),
		currentAnim: defaultAnimation,
		defaultAnim: defaultAnimation,
		playing:     true,
	}
}

// AddAnimation adds an animation to the controller
func (ac *AnimationController) AddAnimation(anim *Animation) {
	ac.animations[anim.Name] = anim
}

// Play starts or switches to an animation
func (ac *AnimationController) Play(name string, restart bool) {
	if anim, exists := ac.animations[name]; exists {
		// If switching animations or restart requested, reset the animation
		if ac.currentAnim != name || restart {
			anim.Reset()
			ac.currentAnim = name
		}
		ac.playing = true
	}
}

// Stop stops the current animation
func (ac *AnimationController) Stop() {
	ac.playing = false
}

// Update updates the current animation (call each frame)
func (ac *AnimationController) Update() {
	if !ac.playing {
		return
	}
	
	anim, exists := ac.animations[ac.currentAnim]
	if !exists {
		return
	}
	
	anim.Update()
	
	// If animation finished and not looping, stop or return to default
	if anim.IsFinished() && !anim.Loop {
		if ac.currentAnim != ac.defaultAnim {
			ac.Play(ac.defaultAnim, true)
		} else {
			ac.playing = false
		}
	}
}

// GetCurrentFrame returns the current frame sprite
func (ac *AnimationController) GetCurrentFrame() *graphics.Sprite {
	anim, exists := ac.animations[ac.currentAnim]
	if !exists {
		return nil
	}
	return anim.GetCurrentFrame()
}

// GetCurrentAnimation returns the name of current animation
func (ac *AnimationController) GetCurrentAnimation() string {
	return ac.currentAnim
}

// IsPlaying returns whether an animation is currently playing
func (ac *AnimationController) IsPlaying() bool {
	return ac.playing
}

// Update updates the animation state (advances frame if needed)
func (a *Animation) Update() {
	if len(a.Frames) == 0 {
		return
	}
	
	a.timer++
	
	// Time to advance to next frame?
	if a.timer >= a.FrameTime {
		a.timer = 0
		a.currentFrame++
		
		// Handle looping or clamping
		if a.currentFrame >= len(a.Frames) {
			if a.Loop {
				a.currentFrame = 0
			} else {
				a.currentFrame = len(a.Frames) - 1
			}
		}
	}
}

// GetCurrentFrame returns the current frame sprite
func (a *Animation) GetCurrentFrame() *graphics.Sprite {
	if len(a.Frames) == 0 {
		return nil
	}
	
	// Ensure currentFrame is valid
	if a.currentFrame < 0 || a.currentFrame >= len(a.Frames) {
		a.currentFrame = 0
	}
	
	return a.Frames[a.currentFrame]
}

// Reset resets the animation to the first frame
func (a *Animation) Reset() {
	a.currentFrame = 0
	a.timer = 0
}

// IsFinished returns true if the animation has completed (non-looping only)
func (a *Animation) IsFinished() bool {
	if a.Loop {
		return false
	}
	return a.currentFrame >= len(a.Frames)-1 && a.timer >= a.FrameTime-1
}

// GetProgress returns animation progress from 0.0 to 1.0
func (a *Animation) GetProgress() float64 {
	if len(a.Frames) == 0 {
		return 1.0
	}
	
	totalFrames := len(a.Frames) * a.FrameTime
	currentPosition := a.currentFrame*a.FrameTime + a.timer
	
	return float64(currentPosition) / float64(totalFrames)
}

// Clone creates a copy of the animation with independent state
func (a *Animation) Clone() *Animation {
	return &Animation{
		Name:       a.Name,
		Frames:     a.Frames, // Share frame data, don't deep copy
		FrameTime:  a.FrameTime,
		Loop:       a.Loop,
		currentFrame: 0,
		timer:      0,
	}
}
