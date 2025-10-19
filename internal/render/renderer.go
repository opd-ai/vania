// Package render provides the game rendering system using Ebiten to display
// procedurally generated sprites, tilesets, and game world on screen with
// camera controls and visual effects.
package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/vania/internal/graphics"
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
func (r *Renderer) RenderEnemy(screen *ebiten.Image, x, y, width, height float64, health, maxHealth int, isInvulnerable bool) {
	// Apply camera offset
	screenX := x + r.camera.X
	screenY := y + r.camera.Y
	
	// Don't render if off screen
	if screenX+width < 0 || screenX > float64(ScreenWidth) ||
		screenY+height < 0 || screenY > float64(ScreenHeight) {
		return
	}
	
	// Create enemy sprite image
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
	
	// Draw health bar above enemy
	if maxHealth > 0 {
		barWidth := width
		barHeight := 4.0
		barY := screenY - 8
		
		// Background
		bgImg := ebiten.NewImage(int(barWidth), int(barHeight))
		bgImg.Fill(color.RGBA{50, 50, 50, 255})
		opts = &ebiten.DrawImageOptions{}
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
