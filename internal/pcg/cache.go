package pcg

import (
	"sync"
)

// AssetCache stores generated assets to avoid regeneration
type AssetCache struct {
	Sprites   map[string]interface{}
	Sounds    map[string]interface{}
	Music     map[string]interface{}
	Narrative map[string]interface{}
	mu        sync.RWMutex
}

// NewAssetCache creates a new asset cache
func NewAssetCache() *AssetCache {
	return &AssetCache{
		Sprites:   make(map[string]interface{}),
		Sounds:    make(map[string]interface{}),
		Music:     make(map[string]interface{}),
		Narrative: make(map[string]interface{}),
	}
}

// GetSprite retrieves a cached sprite
func (c *AssetCache) GetSprite(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Sprites[key]
	return val, ok
}

// SetSprite caches a sprite
func (c *AssetCache) SetSprite(key string, sprite interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Sprites[key] = sprite
}

// GetSound retrieves a cached sound
func (c *AssetCache) GetSound(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Sounds[key]
	return val, ok
}

// SetSound caches a sound
func (c *AssetCache) SetSound(key string, sound interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Sounds[key] = sound
}

// GetMusic retrieves cached music
func (c *AssetCache) GetMusic(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Music[key]
	return val, ok
}

// SetMusic caches music
func (c *AssetCache) SetMusic(key string, music interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Music[key] = music
}

// GetNarrative retrieves cached narrative
func (c *AssetCache) GetNarrative(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Narrative[key]
	return val, ok
}

// SetNarrative caches narrative
func (c *AssetCache) SetNarrative(key string, narrative interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Narrative[key] = narrative
}

// Clear removes all cached assets
func (c *AssetCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Sprites = make(map[string]interface{})
	c.Sounds = make(map[string]interface{})
	c.Music = make(map[string]interface{})
	c.Narrative = make(map[string]interface{})
}
