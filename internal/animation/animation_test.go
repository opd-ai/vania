package animation

import (
	"testing"

	"github.com/opd-ai/vania/internal/graphics"
)

// Helper function to create test sprite
func createTestSprite(id int) *graphics.Sprite {
	return &graphics.Sprite{
		Width:  32,
		Height: 32,
	}
}

// Test NewAnimation
func TestNewAnimation(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
		createTestSprite(3),
	}
	
	anim := NewAnimation("test", frames, 10, true)
	
	if anim.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", anim.Name)
	}
	if len(anim.Frames) != 3 {
		t.Errorf("Expected 3 frames, got %d", len(anim.Frames))
	}
	if anim.FrameTime != 10 {
		t.Errorf("Expected FrameTime 10, got %d", anim.FrameTime)
	}
	if !anim.Loop {
		t.Error("Expected Loop to be true")
	}
	if anim.currentFrame != 0 {
		t.Errorf("Expected currentFrame 0, got %d", anim.currentFrame)
	}
	if anim.timer != 0 {
		t.Errorf("Expected timer 0, got %d", anim.timer)
	}
}

// Test NewAnimation with invalid frameTime
func TestNewAnimationDefaultFrameTime(t *testing.T) {
	frames := []*graphics.Sprite{createTestSprite(1)}
	
	anim := NewAnimation("test", frames, 0, false)
	
	if anim.FrameTime != 10 {
		t.Errorf("Expected default FrameTime 10, got %d", anim.FrameTime)
	}
}

// Test Animation.Update
func TestAnimationUpdate(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
		createTestSprite(3),
	}
	
	anim := NewAnimation("test", frames, 5, true)
	
	// Update 4 times (not enough to advance frame)
	for i := 0; i < 4; i++ {
		anim.Update()
	}
	if anim.currentFrame != 0 {
		t.Errorf("Expected frame 0 after 4 updates, got %d", anim.currentFrame)
	}
	
	// Update once more (should advance to frame 1)
	anim.Update()
	if anim.currentFrame != 1 {
		t.Errorf("Expected frame 1 after 5 updates, got %d", anim.currentFrame)
	}
	
	// Update 5 more times (should advance to frame 2)
	for i := 0; i < 5; i++ {
		anim.Update()
	}
	if anim.currentFrame != 2 {
		t.Errorf("Expected frame 2 after 10 updates, got %d", anim.currentFrame)
	}
	
	// Update 5 more times (should loop back to frame 0)
	for i := 0; i < 5; i++ {
		anim.Update()
	}
	if anim.currentFrame != 0 {
		t.Errorf("Expected frame 0 after looping, got %d", anim.currentFrame)
	}
}

// Test non-looping animation
func TestAnimationNoLoop(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	
	anim := NewAnimation("test", frames, 5, false)
	
	// Update enough to reach last frame
	for i := 0; i < 10; i++ {
		anim.Update()
	}
	
	// Should be on last frame
	if anim.currentFrame != 1 {
		t.Errorf("Expected frame 1, got %d", anim.currentFrame)
	}
	
	// More updates shouldn't advance beyond last frame
	for i := 0; i < 10; i++ {
		anim.Update()
	}
	if anim.currentFrame != 1 {
		t.Errorf("Expected to stay on frame 1, got %d", anim.currentFrame)
	}
}

// Test GetCurrentFrame
func TestGetCurrentFrame(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	
	anim := NewAnimation("test", frames, 5, true)
	
	frame := anim.GetCurrentFrame()
	if frame != frames[0] {
		t.Error("Expected first frame")
	}
	
	// Advance to next frame
	for i := 0; i < 5; i++ {
		anim.Update()
	}
	
	frame = anim.GetCurrentFrame()
	if frame != frames[1] {
		t.Error("Expected second frame")
	}
}

// Test Reset
func TestAnimationReset(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
		createTestSprite(3),
	}
	
	anim := NewAnimation("test", frames, 5, true)
	
	// Advance animation
	for i := 0; i < 10; i++ {
		anim.Update()
	}
	
	// Reset
	anim.Reset()
	
	if anim.currentFrame != 0 {
		t.Errorf("Expected currentFrame 0 after reset, got %d", anim.currentFrame)
	}
	if anim.timer != 0 {
		t.Errorf("Expected timer 0 after reset, got %d", anim.timer)
	}
}

// Test IsFinished
func TestIsFinished(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	
	// Looping animation never finishes
	loopAnim := NewAnimation("loop", frames, 5, true)
	if loopAnim.IsFinished() {
		t.Error("Looping animation should never be finished")
	}
	for i := 0; i < 20; i++ {
		loopAnim.Update()
	}
	if loopAnim.IsFinished() {
		t.Error("Looping animation should never be finished")
	}
	
	// Non-looping animation finishes
	noLoopAnim := NewAnimation("noloop", frames, 5, false)
	if noLoopAnim.IsFinished() {
		t.Error("Animation should not be finished at start")
	}
	
	// Update to last frame
	for i := 0; i < 9; i++ {
		noLoopAnim.Update()
		if noLoopAnim.IsFinished() && i < 8 {
			t.Errorf("Animation finished too early at update %d", i)
		}
	}
	
	if !noLoopAnim.IsFinished() {
		t.Error("Animation should be finished")
	}
}

// Test GetProgress
func TestGetProgress(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	
	anim := NewAnimation("test", frames, 5, true)
	
	progress := anim.GetProgress()
	if progress != 0.0 {
		t.Errorf("Expected progress 0.0 at start, got %f", progress)
	}
	
	// Update halfway through first frame
	for i := 0; i < 2; i++ {
		anim.Update()
	}
	progress = anim.GetProgress()
	if progress < 0.15 || progress > 0.25 {
		t.Errorf("Expected progress around 0.2, got %f", progress)
	}
	
	// Update to end (looping animation loops back, so check before loop)
	for i := 0; i < 7; i++ {
		anim.Update()
	}
	progress = anim.GetProgress()
	// Should be near 90% progress (9 out of 10 total frames)
	if progress < 0.85 || progress > 0.95 {
		t.Errorf("Expected progress around 0.9, got %f", progress)
	}
}

// Test Clone
func TestClone(t *testing.T) {
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	
	original := NewAnimation("test", frames, 5, true)
	original.Update()
	original.Update()
	
	clone := original.Clone()
	
	if clone.Name != original.Name {
		t.Error("Clone should have same name")
	}
	if clone.FrameTime != original.FrameTime {
		t.Error("Clone should have same FrameTime")
	}
	if clone.Loop != original.Loop {
		t.Error("Clone should have same Loop")
	}
	if clone.currentFrame != 0 {
		t.Errorf("Clone should reset currentFrame to 0, got %d", clone.currentFrame)
	}
	if clone.timer != 0 {
		t.Errorf("Clone should reset timer to 0, got %d", clone.timer)
	}
	
	// Frames should be shared (same reference)
	if len(clone.Frames) != len(original.Frames) {
		t.Error("Clone should have same frames")
	}
}

// Test NewAnimationController
func TestNewAnimationController(t *testing.T) {
	controller := NewAnimationController("idle")
	
	if controller.currentAnim != "idle" {
		t.Errorf("Expected currentAnim 'idle', got '%s'", controller.currentAnim)
	}
	if controller.defaultAnim != "idle" {
		t.Errorf("Expected defaultAnim 'idle', got '%s'", controller.defaultAnim)
	}
	if !controller.playing {
		t.Error("Expected playing to be true")
	}
	if len(controller.animations) != 0 {
		t.Errorf("Expected empty animations map, got %d entries", len(controller.animations))
	}
}

// Test AddAnimation
func TestAddAnimation(t *testing.T) {
	controller := NewAnimationController("idle")
	
	frames := []*graphics.Sprite{createTestSprite(1)}
	anim := NewAnimation("walk", frames, 5, true)
	
	controller.AddAnimation(anim)
	
	if len(controller.animations) != 1 {
		t.Errorf("Expected 1 animation, got %d", len(controller.animations))
	}
	
	if controller.animations["walk"] != anim {
		t.Error("Animation not properly added")
	}
}

// Test Play
func TestPlay(t *testing.T) {
	controller := NewAnimationController("idle")
	
	idleFrames := []*graphics.Sprite{createTestSprite(1)}
	walkFrames := []*graphics.Sprite{createTestSprite(2), createTestSprite(3)}
	
	idleAnim := NewAnimation("idle", idleFrames, 10, true)
	walkAnim := NewAnimation("walk", walkFrames, 5, true)
	
	controller.AddAnimation(idleAnim)
	controller.AddAnimation(walkAnim)
	
	// Play walk animation
	controller.Play("walk", false)
	
	if controller.currentAnim != "walk" {
		t.Errorf("Expected currentAnim 'walk', got '%s'", controller.currentAnim)
	}
	
	// Advance walk animation
	for i := 0; i < 3; i++ {
		controller.Update()
	}
	
	// Play walk again without restart (should not reset)
	controller.Play("walk", false)
	if walkAnim.timer == 0 {
		t.Error("Animation should not have been reset")
	}
	
	// Play walk with restart
	controller.Play("walk", true)
	if walkAnim.timer != 0 || walkAnim.currentFrame != 0 {
		t.Error("Animation should have been reset")
	}
}

// Test Stop
func TestStop(t *testing.T) {
	controller := NewAnimationController("idle")
	
	frames := []*graphics.Sprite{createTestSprite(1)}
	anim := NewAnimation("idle", frames, 5, true)
	controller.AddAnimation(anim)
	
	if !controller.playing {
		t.Error("Controller should be playing initially")
	}
	
	controller.Stop()
	
	if controller.playing {
		t.Error("Controller should not be playing after Stop()")
	}
}

// Test Update
func TestControllerUpdate(t *testing.T) {
	controller := NewAnimationController("idle")
	
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	anim := NewAnimation("idle", frames, 5, true)
	controller.AddAnimation(anim)
	
	// Update several times (10 frames at FrameTime=5 means we need 5 updates to change frame)
	for i := 0; i < 6; i++ {
		controller.Update()
	}
	
	// Should have advanced to second frame (after 5 updates)
	if anim.currentFrame != 1 {
		t.Errorf("Expected currentFrame 1 after 6 updates, got %d", anim.currentFrame)
	}
	
	// Stop and update (should not advance)
	controller.Stop()
	prevFrame := anim.currentFrame
	prevTimer := anim.timer
	
	controller.Update()
	
	if anim.currentFrame != prevFrame || anim.timer != prevTimer {
		t.Error("Animation should not advance when stopped")
	}
}

// Test GetCurrentFrame
func TestControllerGetCurrentFrame(t *testing.T) {
	controller := NewAnimationController("idle")
	
	frames := []*graphics.Sprite{
		createTestSprite(1),
		createTestSprite(2),
	}
	anim := NewAnimation("idle", frames, 5, true)
	controller.AddAnimation(anim)
	
	frame := controller.GetCurrentFrame()
	if frame != frames[0] {
		t.Error("Expected first frame")
	}
	
	// Advance
	for i := 0; i < 5; i++ {
		controller.Update()
	}
	
	frame = controller.GetCurrentFrame()
	if frame != frames[1] {
		t.Error("Expected second frame")
	}
}

// Test animation completion behavior
func TestAnimationCompletion(t *testing.T) {
	controller := NewAnimationController("idle")
	
	idleFrames := []*graphics.Sprite{createTestSprite(1)}
	attackFrames := []*graphics.Sprite{createTestSprite(2), createTestSprite(3)}
	
	idleAnim := NewAnimation("idle", idleFrames, 10, true)
	attackAnim := NewAnimation("attack", attackFrames, 5, false)
	
	controller.AddAnimation(idleAnim)
	controller.AddAnimation(attackAnim)
	
	// Play attack animation
	controller.Play("attack", true)
	
	if controller.currentAnim != "attack" {
		t.Error("Should be playing attack")
	}
	
	// Update until attack completes
	for i := 0; i < 15; i++ {
		controller.Update()
	}
	
	// Should return to idle (default)
	if controller.currentAnim != "idle" {
		t.Errorf("Expected to return to 'idle', got '%s'", controller.currentAnim)
	}
}

// Test GetCurrentAnimation
func TestGetCurrentAnimation(t *testing.T) {
	controller := NewAnimationController("idle")
	
	frames := []*graphics.Sprite{createTestSprite(1)}
	idleAnim := NewAnimation("idle", frames, 10, true)
	walkAnim := NewAnimation("walk", frames, 5, true)
	
	controller.AddAnimation(idleAnim)
	controller.AddAnimation(walkAnim)
	
	if controller.GetCurrentAnimation() != "idle" {
		t.Error("Expected idle animation")
	}
	
	controller.Play("walk", false)
	
	if controller.GetCurrentAnimation() != "walk" {
		t.Error("Expected walk animation")
	}
}

// Test IsPlaying
func TestIsPlaying(t *testing.T) {
	controller := NewAnimationController("idle")
	
	if !controller.IsPlaying() {
		t.Error("Controller should be playing initially")
	}
	
	controller.Stop()
	
	if controller.IsPlaying() {
		t.Error("Controller should not be playing after stop")
	}
	
	frames := []*graphics.Sprite{createTestSprite(1)}
	anim := NewAnimation("idle", frames, 10, true)
	controller.AddAnimation(anim)
	
	controller.Play("idle", false)
	
	if !controller.IsPlaying() {
		t.Error("Controller should be playing after play")
	}
}

// Test GenerateEnemyIdleFrames
func TestGenerateEnemyIdleFrames(t *testing.T) {
	gen := NewAnimationGenerator(12345)
	baseSprite := createTestSprite(1)
	
	frames := gen.GenerateEnemyIdleFrames(baseSprite, 4)
	
	if len(frames) != 4 {
		t.Errorf("Expected 4 frames, got %d", len(frames))
	}
	
	for i, frame := range frames {
		if frame == nil {
			t.Errorf("Frame %d is nil", i)
		}
		if frame.Width != baseSprite.Width {
			t.Errorf("Frame %d width mismatch: expected %d, got %d", i, baseSprite.Width, frame.Width)
		}
		if frame.Height != baseSprite.Height {
			t.Errorf("Frame %d height mismatch: expected %d, got %d", i, baseSprite.Height, frame.Height)
		}
	}
}

// Test GenerateEnemyPatrolFrames
func TestGenerateEnemyPatrolFrames(t *testing.T) {
	gen := NewAnimationGenerator(12345)
	baseSprite := createTestSprite(1)
	
	frames := gen.GenerateEnemyPatrolFrames(baseSprite, 4)
	
	if len(frames) != 4 {
		t.Errorf("Expected 4 frames, got %d", len(frames))
	}
	
	for i, frame := range frames {
		if frame == nil {
			t.Errorf("Frame %d is nil", i)
		}
	}
}

// Test GenerateEnemyAttackFrames
func TestGenerateEnemyAttackFrames(t *testing.T) {
	gen := NewAnimationGenerator(12345)
	baseSprite := createTestSprite(1)
	
	frames := gen.GenerateEnemyAttackFrames(baseSprite, 3)
	
	if len(frames) != 3 {
		t.Errorf("Expected 3 frames, got %d", len(frames))
	}
	
	for i, frame := range frames {
		if frame == nil {
			t.Errorf("Frame %d is nil", i)
		}
	}
}

// Test GenerateEnemyDeathFrames
func TestGenerateEnemyDeathFrames(t *testing.T) {
	gen := NewAnimationGenerator(12345)
	baseSprite := createTestSprite(1)
	
	frames := gen.GenerateEnemyDeathFrames(baseSprite, 4)
	
	if len(frames) != 4 {
		t.Errorf("Expected 4 frames, got %d", len(frames))
	}
	
	for i, frame := range frames {
		if frame == nil {
			t.Errorf("Frame %d is nil", i)
		}
	}
}

// Test enemy animation with nil sprite
func TestEnemyAnimationNilSprite(t *testing.T) {
	gen := NewAnimationGenerator(12345)
	
	idleFrames := gen.GenerateEnemyIdleFrames(nil, 4)
	if idleFrames != nil {
		t.Error("Expected nil frames for nil sprite")
	}
	
	patrolFrames := gen.GenerateEnemyPatrolFrames(nil, 4)
	if patrolFrames != nil {
		t.Error("Expected nil frames for nil sprite")
	}
	
	attackFrames := gen.GenerateEnemyAttackFrames(nil, 3)
	if attackFrames != nil {
		t.Error("Expected nil frames for nil sprite")
	}
	
	deathFrames := gen.GenerateEnemyDeathFrames(nil, 4)
	if deathFrames != nil {
		t.Error("Expected nil frames for nil sprite")
	}
}

// Test enemy animation with zero frames
func TestEnemyAnimationZeroFrames(t *testing.T) {
	gen := NewAnimationGenerator(12345)
	baseSprite := createTestSprite(1)
	
	idleFrames := gen.GenerateEnemyIdleFrames(baseSprite, 0)
	if idleFrames != nil {
		t.Error("Expected nil frames for 0 frame count")
	}
	
	patrolFrames := gen.GenerateEnemyPatrolFrames(baseSprite, 0)
	if patrolFrames != nil {
		t.Error("Expected nil frames for 0 frame count")
	}
}

// Test deterministic animation generation
func TestEnemyAnimationDeterminism(t *testing.T) {
	seed := int64(99999)
	baseSprite := createTestSprite(1)
	
	gen1 := NewAnimationGenerator(seed)
	frames1 := gen1.GenerateEnemyIdleFrames(baseSprite, 4)
	
	gen2 := NewAnimationGenerator(seed)
	frames2 := gen2.GenerateEnemyIdleFrames(baseSprite, 4)
	
	if len(frames1) != len(frames2) {
		t.Errorf("Frame counts differ: %d vs %d", len(frames1), len(frames2))
	}
	
	// Both generators with same seed should produce same number of frames
	for i := range frames1 {
		if frames1[i].Width != frames2[i].Width || frames1[i].Height != frames2[i].Height {
			t.Errorf("Frame %d dimensions differ", i)
		}
	}
}

