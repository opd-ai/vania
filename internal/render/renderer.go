// Package render provides the game rendering system using Ebiten to display
// procedurally generated sprites, tilesets, and game world on screen with
// camera controls and visual effects.
package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/vania/internal/graphics"
	"github.com/opd-ai/vania/internal/particle"
	"github.com/opd-ai/vania/internal/world"
)

const (
	// Screen dimensions
	ScreenWidth  = 960
	ScreenHeight = 640

	// Tile size in pixels
	TileSize = 32

	// Camera settings
	CameraSpeed = 4.0

	// UI Layout Constants
	UIMargin = 10

	// Health bar layout
	HealthBarWidth  = 200
	HealthBarHeight = 20
	HealthBarX      = UIMargin
	HealthBarY      = UIMargin
	HealthBarBorder = 2

	// Ability icons layout
	AbilityIconSize    = 30
	AbilityIconSpacing = 5
	AbilityIconY       = HealthBarY + HealthBarHeight + UIMargin

	// Enemy health bar layout
	EnemyHealthBarHeight = 4
	EnemyHealthBarOffset = 8

	// Message layout
	MessageWidth      = 200
	MessageHeight     = 40
	ProgressBarHeight = 3
)

// Camera represents the game camera
type Camera struct {
	X, Y   float64
	Width  int
	Height int
}

// Renderer handles all game rendering
type Renderer struct {
	screen      *ebiten.Image
	camera      *Camera
	tileImages  map[string]*ebiten.Image
	bgColor     color.Color
	textManager *TextRenderManager
}

// NewRenderer creates a new renderer
func NewRenderer() *Renderer {
	return &Renderer{
		camera: &Camera{
			X:      0,
			Y:      0,
			Width:  ScreenWidth,
			Height: ScreenHeight,
		},
		tileImages:  make(map[string]*ebiten.Image),
		bgColor:     color.RGBA{20, 20, 30, 255}, // Dark blue background
		textManager: NewTextRenderManager(true),  // Enable color rendering by default
	}
}

// RenderWorld draws the game world to the screen
func (r *Renderer) RenderWorld(screen *ebiten.Image, currentRoom *world.Room, tilesets map[string]*graphics.Tileset) {
	r.screen = screen

	// Clear screen
	screen.Fill(r.bgColor)

	if currentRoom == nil {
		ebitenutil.DebugPrint(screen, "No room to render")
		return
	}

	// Render room background
	r.renderRoomBackground(screen, currentRoom, tilesets)

	// Render platforms
	r.renderPlatforms(screen, currentRoom, tilesets)

	// Render hazards
	r.renderHazards(screen, currentRoom)

	// Render doors
	r.renderDoors(screen, currentRoom)
}

// renderRoomBackground draws the room's background tiles
func (r *Renderer) renderRoomBackground(screen *ebiten.Image, room *world.Room, tilesets map[string]*graphics.Tileset) {
	if room.Biome == nil {
		return
	}

	// Get tileset for this biome
	tileset, ok := tilesets[room.Biome.Name]
	if !ok || tileset == nil {
		return
	}

	// Calculate room dimensions in tiles
	roomWidthTiles := ScreenWidth / TileSize
	roomHeightTiles := ScreenHeight / TileSize

	// Render background tiles
	// Get background tile
	bgTile, ok := tileset.Tiles[graphics.BackgroundTile]
	if !ok || bgTile == nil || bgTile.Image == nil {
		return
	}

	bgImage := ebiten.NewImageFromImage(bgTile.Image)

	for y := 0; y < roomHeightTiles; y++ {
		for x := 0; x < roomWidthTiles; x++ {
			// Draw tile
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x*TileSize), float64(y*TileSize))
			screen.DrawImage(bgImage, opts)
		}
	}
}

// renderPlatforms draws platforms in the room
func (r *Renderer) renderPlatforms(screen *ebiten.Image, room *world.Room, tilesets map[string]*graphics.Tileset) {
	if room.Biome == nil {
		return
	}

	// Get tileset for this biome
	tileset, ok := tilesets[room.Biome.Name]
	if !ok || tileset == nil {
		return
	}

	// Select platform tile (use solid tile if available)
	platformTile, ok := tileset.Tiles[graphics.SolidTile]
	if !ok || platformTile == nil || platformTile.Image == nil {
		return
	}

	platformImg := ebiten.NewImageFromImage(platformTile.Image)

	// Render each platform
	for _, platform := range room.Platforms {
		// Draw platform tiles
		for px := 0; px < platform.Width; px++ {
			for py := 0; py < platform.Height; py++ {
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(
					float64(platform.X+px*TileSize),
					float64(platform.Y+py*TileSize),
				)
				screen.DrawImage(platformImg, opts)
			}
		}
	}
}

// renderHazards draws hazards in the room
func (r *Renderer) renderHazards(screen *ebiten.Image, room *world.Room) {
	for _, hazard := range room.Hazards {
		// Choose hazard color based on type
		var hazardColor color.Color
		switch hazard.Type {
		case "spike":
			hazardColor = color.RGBA{150, 150, 150, 255} // Gray
		case "lava":
			hazardColor = color.RGBA{255, 100, 0, 255} // Orange-red
		case "electric":
			hazardColor = color.RGBA{100, 200, 255, 255} // Electric blue
		default:
			hazardColor = color.RGBA{200, 0, 0, 255} // Red
		}

		// Draw hazard as a colored rectangle
		hazardImg := ebiten.NewImage(hazard.Width, hazard.Height)
		hazardImg.Fill(hazardColor)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(hazard.X), float64(hazard.Y))
		screen.DrawImage(hazardImg, opts)
	}
}

// renderDoors draws doors/exits in the room
func (r *Renderer) renderDoors(screen *ebiten.Image, room *world.Room) {
	for _, door := range room.Doors {
		// Choose door color based on whether it's locked
		var doorColor color.Color
		if door.Locked {
			doorColor = color.RGBA{150, 50, 50, 255} // Dark red for locked
		} else {
			doorColor = color.RGBA{100, 150, 200, 255} // Blue for unlocked
		}

		// Draw door frame
		doorImg := ebiten.NewImage(door.Width, door.Height)
		doorImg.Fill(doorColor)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(door.X), float64(door.Y))
		screen.DrawImage(doorImg, opts)

		// Draw inner part (lighter)
		innerColor := color.RGBA{150, 200, 255, 200}
		if door.Locked {
			innerColor = color.RGBA{200, 100, 100, 200}
		}
		innerImg := ebiten.NewImage(door.Width-8, door.Height-8)
		innerImg.Fill(innerColor)

		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(door.X+4), float64(door.Y+4))
		screen.DrawImage(innerImg, opts)
	}
}

// RenderPlayer draws the player sprite
func (r *Renderer) RenderPlayer(screen *ebiten.Image, x, y float64, sprite *graphics.Sprite) {
	if sprite == nil || sprite.Image == nil {
		// Draw a simple colored square as fallback
		playerImg := ebiten.NewImage(32, 32)
		playerImg.Fill(color.RGBA{100, 200, 100, 255}) // Green

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, y)
		screen.DrawImage(playerImg, opts)
		return
	}

	// Convert sprite to ebiten image
	playerImg := ebiten.NewImageFromImage(sprite.Image)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	screen.DrawImage(playerImg, opts)
}

// RenderUI draws the user interface (health, abilities, etc.)
func (r *Renderer) RenderUI(screen *ebiten.Image, health, maxHealth int, abilities map[string]bool) {
	barX, barY, barHeight := r.renderEnhancedHealthBar(screen, health, maxHealth)
	r.renderAbilityIcons(screen, abilities, barX, barY+barHeight+10)
}

// renderEnhancedHealthBar draws an improved health bar with segments and color coding
func (r *Renderer) renderEnhancedHealthBar(screen *ebiten.Image, health, maxHealth int) (int, int, int) {
	// Use layout constants
	barWidth := HealthBarWidth
	barHeight := HealthBarHeight
	barX := HealthBarX
	barY := HealthBarY
	borderWidth := HealthBarBorder

	// Validate health bounds
	if health < 0 {
		health = 0
	}
	if health > maxHealth {
		health = maxHealth
	}

	// Draw border (white outline)
	borderImg := ebiten.NewImage(barWidth+borderWidth*2, barHeight+borderWidth*2)
	borderImg.Fill(color.RGBA{255, 255, 255, 255})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(barX-borderWidth), float64(barY-borderWidth))
	screen.DrawImage(borderImg, opts)

	// Draw background (dark gray)
	bgImg := ebiten.NewImage(barWidth, barHeight)
	bgImg.Fill(color.RGBA{40, 40, 40, 255})
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(bgImg, opts)

	// Calculate health percentage and color
	healthPercent := 0.0
	if maxHealth > 0 {
		healthPercent = float64(health) / float64(maxHealth)
	}

	var healthColor color.Color
	if healthPercent > 0.66 {
		// Green for healthy (>66%)
		healthColor = color.RGBA{100, 200, 100, 255}
	} else if healthPercent > 0.33 {
		// Yellow for wounded (33-66%)
		healthColor = color.RGBA{200, 200, 100, 255}
	} else {
		// Red for critical (<33%)
		healthColor = color.RGBA{200, 50, 50, 255}
	}

	// Draw health fill with calculated color
	if maxHealth > 0 && health > 0 {
		fillWidth := int(float64(barWidth) * healthPercent)
		if fillWidth > 0 {
			fillImg := ebiten.NewImage(fillWidth, barHeight)
			fillImg.Fill(healthColor)
			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(barX), float64(barY))
			screen.DrawImage(fillImg, opts)
		}
	}

	// Draw health segments (dividers every 10 HP or maxHealth/10, whichever is smaller)
	segmentCount := maxHealth / 10
	if segmentCount > 10 {
		segmentCount = 10 // Max 10 segments for readability
	}
	if segmentCount > 0 {
		segmentWidth := float64(barWidth) / float64(segmentCount)
		for i := 1; i < segmentCount; i++ {
			segmentX := barX + int(float64(i)*segmentWidth)
			segmentImg := ebiten.NewImage(1, barHeight)
			segmentImg.Fill(color.RGBA{60, 60, 60, 255}) // Dark divider
			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(segmentX), float64(barY))
			screen.DrawImage(segmentImg, opts)
		}
	}

	return barX, barY, barHeight
}

// renderAbilityIcons draws ability indicators with procedural icons
func (r *Renderer) renderAbilityIcons(screen *ebiten.Image, abilities map[string]bool, startX, startY int) {
	abilitySize := AbilityIconSize
	abilitySpacing := AbilityIconSpacing

	abilityNames := []string{"double_jump", "dash", "wall_jump", "glide"}
	for i, abilityName := range abilityNames {
		hasAbility := abilities[abilityName]
		x := startX + i*(abilitySize+abilitySpacing)
		r.renderAbilityIcon(screen, abilityName, x, startY, abilitySize, hasAbility)
	}
}

// renderAbilityIcon draws a single ability icon with symbolic representation
func (r *Renderer) renderAbilityIcon(screen *ebiten.Image, ability string, x, y, size int, unlocked bool) {
	// Create icon background
	iconImg := ebiten.NewImage(size, size)

	// Choose colors based on unlock status
	var bgColor, iconColor color.Color
	if unlocked {
		bgColor = color.RGBA{50, 100, 150, 255}    // Blue background when unlocked
		iconColor = color.RGBA{200, 220, 255, 255} // Light blue icon
	} else {
		bgColor = color.RGBA{30, 30, 30, 255}   // Dark background when locked
		iconColor = color.RGBA{80, 80, 80, 255} // Dark gray icon
	}

	// Fill background
	iconImg.Fill(bgColor)

	// Draw procedural icon based on ability type
	switch ability {
	case "double_jump":
		r.drawJumpIcon(iconImg, iconColor, size)
	case "dash":
		r.drawDashIcon(iconImg, iconColor, size)
	case "wall_jump":
		r.drawWallJumpIcon(iconImg, iconColor, size)
	case "glide":
		r.drawGlideIcon(iconImg, iconColor, size)
	}

	// Draw border for unlocked abilities
	if unlocked {
		r.drawIconBorder(iconImg, size, color.RGBA{255, 255, 255, 150})
	}

	// Draw to screen
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(iconImg, opts)
}

// drawJumpIcon draws upward arrows to represent double jump
func (r *Renderer) drawJumpIcon(img *ebiten.Image, col color.Color, size int) {
	// Draw two upward arrows
	mid := size / 2

	// First arrow (bottom)
	r.drawArrowUp(img, mid, size-8, 6, col)

	// Second arrow (top)
	r.drawArrowUp(img, mid, size-16, 4, col)
}

// drawDashIcon draws horizontal lines to represent dash/speed
func (r *Renderer) drawDashIcon(img *ebiten.Image, col color.Color, size int) {
	mid := size / 2
	lineHeight := 2

	// Draw horizontal speed lines
	for i := 0; i < 3; i++ {
		y := mid - 4 + i*4
		lineImg := ebiten.NewImage(size-8, lineHeight)
		lineImg.Fill(col)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(4, float64(y))
		img.DrawImage(lineImg, opts)
	}
}

// drawWallJumpIcon draws a wall and figure to represent wall jump
func (r *Renderer) drawWallJumpIcon(img *ebiten.Image, col color.Color, size int) {
	// Draw wall on left side
	wallImg := ebiten.NewImage(3, size-6)
	wallImg.Fill(col)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(3, 3)
	img.DrawImage(wallImg, opts)

	// Draw figure jumping away from wall
	figureImg := ebiten.NewImage(6, 8)
	figureImg.Fill(col)
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(size-12), float64(size/2-4))
	img.DrawImage(figureImg, opts)

	// Draw jump arc
	r.drawArrowUp(img, size-6, size/2+2, 3, col)
}

// drawGlideIcon draws wing-like shape to represent glide
func (r *Renderer) drawGlideIcon(img *ebiten.Image, col color.Color, size int) {
	mid := size / 2

	// Draw wing shape (triangle-like)
	wingPoints := []struct{ x, y int }{
		{mid, 6},            // Top center
		{6, mid + 4},        // Bottom left
		{size - 6, mid + 4}, // Bottom right
	}

	// Draw wing outline
	for i := 0; i < len(wingPoints); i++ {
		start := wingPoints[i]
		end := wingPoints[(i+1)%len(wingPoints)]
		r.drawLine(img, start.x, start.y, end.x, end.y, col)
	}

	// Add wing details (inner lines)
	r.drawLine(img, mid, 6, 8, mid, col)
	r.drawLine(img, mid, 6, size-8, mid, col)
}

// drawArrowUp draws an upward pointing arrow
func (r *Renderer) drawArrowUp(img *ebiten.Image, x, y, size int, col color.Color) {
	// Arrow head
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			if x-j >= 0 && x+j < img.Bounds().Dx() && y-i >= 0 {
				pixelImg := ebiten.NewImage(1, 1)
				pixelImg.Fill(col)
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x-j), float64(y-i))
				img.DrawImage(pixelImg, opts)
				if j > 0 {
					opts = &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(float64(x+j), float64(y-i))
					img.DrawImage(pixelImg, opts)
				}
			}
		}
	}
}

// drawLine draws a line between two points
func (r *Renderer) drawLine(img *ebiten.Image, x1, y1, x2, y2 int, col color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := dx
	if dy > dx {
		steps = dy
	}
	if steps < 0 {
		steps = -steps
	}

	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)

	x := float64(x1)
	y := float64(y1)

	for i := 0; i <= steps; i++ {
		if int(x) >= 0 && int(x) < img.Bounds().Dx() && int(y) >= 0 && int(y) < img.Bounds().Dy() {
			pixelImg := ebiten.NewImage(1, 1)
			pixelImg.Fill(col)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(x, y)
			img.DrawImage(pixelImg, opts)
		}
		x += xInc
		y += yInc
	}
}

// drawIconBorder draws a border around an icon
func (r *Renderer) drawIconBorder(img *ebiten.Image, size int, col color.Color) {
	// Top and bottom borders
	topImg := ebiten.NewImage(size, 1)
	topImg.Fill(col)
	bottomImg := ebiten.NewImage(size, 1)
	bottomImg.Fill(col)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, 0)
	img.DrawImage(topImg, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, float64(size-1))
	img.DrawImage(bottomImg, opts)

	// Left and right borders
	leftImg := ebiten.NewImage(1, size)
	leftImg.Fill(col)
	rightImg := ebiten.NewImage(1, size)
	rightImg.Fill(col)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, 0)
	img.DrawImage(leftImg, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(size-1), 0)
	img.DrawImage(rightImg, opts)
}

// UpdateCamera updates camera position to follow target
func (r *Renderer) UpdateCamera(targetX, targetY float64) {
	// Simple camera that centers on target
	r.camera.X = targetX - float64(r.camera.Width)/2
	r.camera.Y = targetY - float64(r.camera.Height)/2
}

// GetCameraOffset returns the camera offset for positioning
func (r *Renderer) GetCameraOffset() (float64, float64) {
	return -r.camera.X, -r.camera.Y
}

// RenderText renders text using the text rendering abstraction
func (r *Renderer) RenderText(screen *ebiten.Image, text string, x, y int, col color.Color) {
	if r.textManager != nil {
		r.textManager.DrawText(screen, text, x, y, col)
	} else {
		// Fallback to debug text if no text manager
		ebitenutil.DebugPrintAt(screen, text, x, y)
	}
}

// MeasureText measures text dimensions using the current text renderer
func (r *Renderer) MeasureText(text string) (width, height int) {
	if r.textManager != nil {
		return r.textManager.MeasureText(text)
	}
	// Fallback measurements for debug text
	return len(text) * 6, 16
}

// SetTextColorMode enables or disables colored text rendering
func (r *Renderer) SetTextColorMode(enabled bool) {
	if r.textManager != nil {
		r.textManager.SetColorMode(enabled)
	}
}

// RenderEnemy draws an enemy to the screen
func (r *Renderer) RenderEnemy(screen *ebiten.Image, x, y, width, height float64, health, maxHealth int, isInvulnerable bool, sprite *graphics.Sprite) {
	// Apply camera offset
	screenX := x + r.camera.X
	screenY := y + r.camera.Y

	// Don't render if off screen
	if screenX+width < 0 || screenX > float64(ScreenWidth) ||
		screenY+height < 0 || screenY > float64(ScreenHeight) {
		return
	}

	// Draw enemy sprite
	if sprite != nil && sprite.Image != nil {
		// Use the animated sprite
		enemyImg := ebiten.NewImageFromImage(sprite.Image)

		// Apply transparency when invulnerable
		opts := &ebiten.DrawImageOptions{}
		if isInvulnerable {
			opts.ColorM.Scale(1, 1, 1, 0.5) // Half transparency
		}
		opts.GeoM.Translate(screenX, screenY)
		screen.DrawImage(enemyImg, opts)
	} else {
		// Fallback to colored rectangle if no sprite
		enemyImg := ebiten.NewImage(int(width), int(height))

		// Enemy color (red with transparency when invulnerable)
		enemyColor := color.RGBA{200, 50, 50, 255}
		if isInvulnerable {
			enemyColor = color.RGBA{200, 50, 50, 128} // Transparent when hit
		}
		enemyImg.Fill(enemyColor)

		// Draw enemy sprite
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(screenX, screenY)
		screen.DrawImage(enemyImg, opts)
	}

	// Draw health bar above enemy
	if maxHealth > 0 {
		barWidth := width
		barHeight := float64(EnemyHealthBarHeight)
		barY := screenY - EnemyHealthBarOffset

		// Background
		bgImg := ebiten.NewImage(int(barWidth), int(barHeight))
		bgImg.Fill(color.RGBA{50, 50, 50, 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(screenX, barY)
		screen.DrawImage(bgImg, opts)

		// Health fill
		fillWidth := barWidth * float64(health) / float64(maxHealth)
		if fillWidth > 0 {
			fillImg := ebiten.NewImage(int(fillWidth), int(barHeight))
			fillImg.Fill(color.RGBA{100, 200, 100, 255}) // Green
			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(screenX, barY)
			screen.DrawImage(fillImg, opts)
		}
	}
}

// RenderAttackEffect draws player attack visual effect
func (r *Renderer) RenderAttackEffect(screen *ebiten.Image, x, y, width, height float64) {
	if width <= 0 || height <= 0 {
		return
	}

	// Apply camera offset
	screenX := x + r.camera.X
	screenY := y + r.camera.Y

	// Create attack effect image (semi-transparent yellow)
	attackImg := ebiten.NewImage(int(width), int(height))
	attackImg.Fill(color.RGBA{255, 255, 100, 128})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(screenX, screenY)
	screen.DrawImage(attackImg, opts)
}

// RenderTransitionEffect renders a fade transition effect
func (r *Renderer) RenderTransitionEffect(screen *ebiten.Image, progress float64) {
	if progress <= 0 {
		return
	}

	// Fade to black during transition
	alpha := uint8(progress * 255)
	fadeImg := ebiten.NewImage(ScreenWidth, ScreenHeight)
	fadeImg.Fill(color.RGBA{0, 0, 0, alpha})

	screen.DrawImage(fadeImg, &ebiten.DrawImageOptions{})
}

// RenderParticles draws all particles on the screen
func (r *Renderer) RenderParticles(screen *ebiten.Image, particles []*particle.Particle) {
	for _, p := range particles {
		if p == nil || !p.IsAlive() {
			continue
		}

		// Calculate screen position relative to camera
		screenX := p.X - r.camera.X
		screenY := p.Y - r.camera.Y

		// Skip if particle is outside screen bounds
		if screenX < -10 || screenX > float64(ScreenWidth)+10 ||
			screenY < -10 || screenY > float64(ScreenHeight)+10 {
			continue
		}

		// Create a simple circle image for the particle
		size := int(p.Size)
		if size < 1 {
			size = 1
		}

		particleImg := ebiten.NewImage(size*2, size*2)

		// Apply alpha to color
		col := p.Color
		col.A = p.Alpha

		// Draw a simple filled circle (square for simplicity)
		particleImg.Fill(col)

		// Draw the particle
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(screenX-float64(size), screenY-float64(size))

		// Apply rotation if any
		if p.Rotation != 0 {
			opts.GeoM.Translate(-float64(size), -float64(size))
			opts.GeoM.Rotate(p.Rotation)
			opts.GeoM.Translate(float64(size), float64(size))
		}

		screen.DrawImage(particleImg, opts)
	}
}

// RenderItem draws a collectible item to the screen
func (r *Renderer) RenderItem(screen *ebiten.Image, x, y, width, height float64, collected bool, sprite *graphics.Sprite) {
	// Don't render if collected
	if collected {
		return
	}

	// If sprite is available, use it
	if sprite != nil && sprite.Image != nil {
		itemImg := ebiten.NewImageFromImage(sprite.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, y)
		screen.DrawImage(itemImg, opts)
		return
	}

	// Fallback: Draw a simple colored box with glow effect
	// Create glow effect (larger, semi-transparent)
	glowSize := int(width * 1.5)
	glowImg := ebiten.NewImage(glowSize, glowSize)
	glowImg.Fill(color.RGBA{255, 215, 0, 60}) // Golden glow

	glowOpts := &ebiten.DrawImageOptions{}
	glowOpts.GeoM.Translate(x-float64(glowSize-int(width))/2, y-float64(glowSize-int(height))/2)
	screen.DrawImage(glowImg, glowOpts)

	// Draw main item box
	itemImg := ebiten.NewImage(int(width), int(height))
	itemImg.Fill(color.RGBA{255, 215, 0, 255}) // Gold

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	screen.DrawImage(itemImg, opts)

	// Draw inner detail (smaller box)
	innerSize := int(width * 0.6)
	innerImg := ebiten.NewImage(innerSize, innerSize)
	innerImg.Fill(color.RGBA{255, 255, 200, 255}) // Light yellow

	innerOpts := &ebiten.DrawImageOptions{}
	innerOpts.GeoM.Translate(x+float64(int(width)-innerSize)/2, y+float64(int(height)-innerSize)/2)
	screen.DrawImage(innerImg, innerOpts)
}
