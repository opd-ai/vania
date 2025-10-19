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
)

// Camera represents the game camera
type Camera struct {
	X, Y   float64
	Width  int
	Height int
}

// Renderer handles all game rendering
type Renderer struct {
	screen     *ebiten.Image
	camera     *Camera
	tileImages map[string]*ebiten.Image
	bgColor    color.Color
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
		tileImages: make(map[string]*ebiten.Image),
		bgColor:    color.RGBA{20, 20, 30, 255}, // Dark blue background
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
	// Draw health bar
	barWidth := 200
	barHeight := 20
	barX := 10
	barY := 10
	
	// Background
	bgImg := ebiten.NewImage(barWidth, barHeight)
	bgImg.Fill(color.RGBA{50, 50, 50, 255})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(bgImg, opts)
	
	// Health fill
	if maxHealth > 0 {
		fillWidth := int(float64(barWidth) * float64(health) / float64(maxHealth))
		if fillWidth > 0 {
			fillImg := ebiten.NewImage(fillWidth, barHeight)
			fillImg.Fill(color.RGBA{200, 50, 50, 255}) // Red
			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(barX), float64(barY))
			screen.DrawImage(fillImg, opts)
		}
	}
	
	// Draw ability indicators
	abilityY := barY + barHeight + 10
	abilitySize := 30
	abilitySpacing := 5
	abilityX := barX
	
	abilityNames := []string{"double_jump", "dash", "wall_jump", "glide"}
	for i, abilityName := range abilityNames {
		hasAbility := abilities[abilityName]
		
		// Draw ability icon
		var abilityColor color.Color
		if hasAbility {
			abilityColor = color.RGBA{100, 150, 255, 255} // Blue when unlocked
		} else {
			abilityColor = color.RGBA{30, 30, 30, 255} // Dark when locked
		}
		
		abilityImg := ebiten.NewImage(abilitySize, abilitySize)
		abilityImg.Fill(abilityColor)
		
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(abilityX+i*(abilitySize+abilitySpacing)), float64(abilityY))
		screen.DrawImage(abilityImg, opts)
	}
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
		barHeight := 4.0
		barY := screenY - 8
		
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
