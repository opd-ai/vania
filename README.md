# VANIA - Procedural Metroidvania Game Engine

A fully procedurally-generated retro Metroidvania game written in pure Go, where **ALL assets** (graphics, audio, story, levels) are generated algorithmically at runtime.

## 🎮 Overview

VANIA is an advanced procedural content generation (PCG) project that eliminates traditional asset creation entirely. Instead of loading pre-made sprites, sounds, or story text, the game generates:

- **Pixel art sprites** through algorithmic drawing
- **Sound effects and music** through synthesis
- **Narrative elements** through procedural story generation
- **World layouts and progression** through algorithmic level design

Each game is generated from a single seed value, providing infinite unique, playable experiences.

## ✨ Features

### Procedural Graphics Generation
- Sprite generation using cellular automata and symmetry
- Multiple color schemes (complementary, triadic, analogous, monochromatic)
- Tileset generation with biome-specific themes
- Automatic palette generation based on HSV color theory
- Outline and shading effects for visual clarity

### Procedural Audio Synthesis
- Waveform generators (sine, square, sawtooth, triangle, noise)
- ADSR envelope system for realistic sound shaping
- Sound effect generation (jump, attack, hit, pickup, etc.)
- Music composition with chord progressions and multiple layers
- Biome-specific music with appropriate mood and scale

### Procedural Narrative Generation
- Story themes (fantasy, sci-fi, horror, mystical, post-apocalyptic)
- Character and faction generation
- Item and room descriptions
- World lore and player motivations
- Dynamic narrative elements tied to gameplay

### Advanced World Generation
- Graph-based world structure ensuring playability
- 80-150 procedurally generated rooms per world
- 4-6 distinct biomes with unique characteristics
- Ability-gated progression (Metroidvania-style)
- Critical path generation with side branches
- Backtracking shortcuts

### Entity & Combat Generation
- Procedural enemy generation scaled to danger level
- Boss fights with multiple phases
- Behavior patterns (patrol, chase, flee, flying, etc.)
- Attack types (melee, ranged, area, contact damage)
- Ability progression system
- Item generation (weapons, consumables, key items, upgrades)

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/opd-ai/vania.git
cd vania

# Install dependencies
go mod tidy

# Build the game
go build -o vania ./cmd/game
```

### Running the Game

```bash
# Generate a random game (uses current timestamp as seed)
./vania

# Generate a specific game from a seed
./vania --seed 42

# Play the game with rendering (NEW!)
./vania --seed 42 --play

# Share seeds with friends to play the same generated game!
./vania --seed 1337 --play
```

**Note**: The `--play` flag launches the full game with rendering, physics, controls, enemies, and combat. See [docs/RENDERING.md](docs/RENDERING.md) for detailed setup instructions and [docs/systems/COMBAT_SYSTEM.md](docs/systems/COMBAT_SYSTEM.md) for combat mechanics.

## 📊 Example Output

```
╔════════════════════════════════════════════════════════╗
║                                                        ║
║         VANIA - Procedural Metroidvania                ║
║         Pure Go Procedural Generation Demo             ║
║                                                        ║
╚════════════════════════════════════════════════════════╝

Master Seed: 42

📖 NARRATIVE
  Theme:              horror
  Mood:               epic
  Civilization:       haunted asylum
  Catastrophe:        a plague transformed people into monsters
  Player Motivation:  break the curse binding you to this place

🌍 WORLD
  Total Rooms:        85
  Boss Rooms:         10
  Biomes:             5

👾 ENTITIES
  Regular Enemies:    10
  Boss Enemies:       10
  Items:              43
  Abilities:          8
```

## 🏗️ Architecture

```
/cmd/game              - Entry point and CLI
/internal/
  ├── pcg/             - Core PCG framework (seed management, caching, validation)
  ├── graphics/        - Sprite and tileset generation
  ├── audio/           - Sound synthesis and music generation
  ├── narrative/       - Story and text generation
  ├── world/           - Level and biome generation
  ├── entity/          - Enemy, boss, and item generation + AI behaviors
  ├── render/          - Ebiten rendering system
  ├── input/           - Input handling
  ├── physics/         - Collision detection and physics
  └── engine/          - Game engine, integration, and combat system
```

## 🎯 Key Technical Achievements

### Zero External Assets
- No sprite files, sound files, or text files
- Everything generated from mathematical algorithms
- Reproducible: same seed = same game

### Deterministic Generation
- Single master seed derives all subsystem seeds
- Hash-based seed derivation for independent subsystems
- Consistent output across runs and platforms

### Quality Assurance
- Generation validation ensures playability
- Visual coherence scoring
- Audio harmony validation
- Narrative consistency verification

### Performance
- Generation completes in ~0.3 seconds
- Efficient caching system
- Optimized for low memory usage

## 🔧 Generation Pipeline

```
Master Seed (user input or timestamp)
    │
    ├─→ Narrative Seed → Story/Theme/Lore (influences all systems)
    │
    ├─→ Graphics Seed → Sprites/Tilesets
    │       ├─→ Player appearance
    │       ├─→ Enemy designs per biome
    │       ├─→ Item appearances
    │       └─→ Tile textures
    │
    ├─→ Audio Seed → Sound/Music
    │       ├─→ SFX library
    │       ├─→ Music tracks per biome
    │       └─→ Boss themes
    │
    ├─→ World Seed → Level Generation
    │       ├─→ Biome layout
    │       ├─→ Room graph
    │       └─→ Ability gate placement
    │
    └─→ Entity Seed → Characters/Items
            ├─→ Enemy roster per biome
            ├─→ Boss designs per region
            └─→ Ability unlock order
```

## 🎨 Procedural Content Generation Techniques

### Graphics
- **Cellular Automata**: Creates organic sprite shapes
- **Symmetry Transforms**: Ensures visual balance
- **HSV Color Theory**: Generates harmonious palettes
- **Flood Fill Algorithms**: Assigns colors to regions
- **Edge Detection**: Adds outlines for clarity

### Audio
- **Additive Synthesis**: Combines waveforms
- **ADSR Envelopes**: Shapes sound over time
- **Low-pass Filtering**: Smooths harsh sounds
- **Music Theory**: Generates valid chord progressions
- **Layered Composition**: Mixes melody, harmony, bass, drums

### Narrative
- **Template-based Generation**: Fills story templates
- **Constraint-based Selection**: Ensures thematic consistency
- **Markov-like Chains**: Creates varied text
- **Name Generation**: Produces pronounceable names

### World
- **Graph Theory**: Ensures connected, solvable worlds
- **Ability Gating**: Creates Metroidvania-style progression
- **Biome Assignment**: Distributes environments evenly
- **Procedural Platforms**: Generates traversable layouts

## 📈 Quality Metrics

```go
type QualityMetrics struct {
    Completability   float64 // % of seeds beatable
    GenerationTime   int64   // milliseconds
    VisualCoherence  float64 // 0-10 aesthetic score
    AudioHarmony     float64 // 0-10 music quality
    NarrativeScore   float64 // 0-10 story coherence
    ContentDiversity float64 // Difference between seeds
}
```

## 🛠️ Development Status

### Implemented ✅
- [x] Core PCG framework (seed management, caching, validation)
- [x] Procedural sprite generation
- [x] Procedural tileset generation
- [x] Color palette generation
- [x] Audio synthesis engine
- [x] Sound effect generation
- [x] Music composition system
- [x] Narrative generation
- [x] World graph generation
- [x] Biome system
- [x] Enemy generation
- [x] Boss generation
- [x] Item and ability generation
- [x] Integration pipeline
- [x] CLI interface
- [x] **Ebiten-based rendering system** ✨
- [x] **Player movement and physics** ✨
- [x] **Collision detection and platforming** ✨
- [x] **Input handling system** ✨
- [x] **Camera system** ✨
- [x] **UI/HUD rendering** ✨
- [x] **Enemy AI system** ✨
- [x] **Combat system** ✨
- [x] **Room transitions** ✨
- [x] **Animation system** ✨ 
- [x] **Save/load system** ✨
- [x] **Particle effects system** ✨
- [x] **Door/Key ability-gating system** ✨
- [x] **Advanced enemy animations** ✨
- [x] **Adaptive music system** ✨
- [x] **Achievement system** ✨ NEW

### In Progress 🚧
- None

### Recently Completed ✨
- [x] **Achievement system** - 19 unique achievements across 6 categories with progress tracking and persistence ✨ NEW
- [x] **Adaptive music system** - Dynamic multi-layer music that responds to gameplay (exploration, combat, boss fights)
- [x] **Advanced enemy animations** - Procedural animation frames for enemies (idle, patrol, attack, death, hit)
- [x] **Item collection system** - Visible items in treasure rooms with collection feedback
- [x] **Door/Key system** - Ability-gated progression with locked doors and automatic unlocking
- [x] **Particle effects system** - Visual feedback for combat and movement
- [x] **Save/load system** - Multiple save slots with automatic checkpoints

### Planned 📋
- [ ] Advanced enemy AI (learning behaviors, coordinated attacks)
- [ ] Puzzle generation
- [ ] Speedrun timer
- [ ] Seed leaderboards

## 🤝 Contributing

This is a demonstration project showcasing advanced PCG techniques. Contributions are welcome!

## 📄 License

MIT License - See LICENSE file for details

## 🙏 Acknowledgments

- Inspired by classic Metroidvania games (Castlevania, Metroid, Hollow Knight)
- PCG techniques from academic research in procedural content generation
- Audio synthesis concepts from digital signal processing

## 📚 Documentation

### System Documentation
Detailed documentation for each game system:
- [Achievement System](docs/systems/ACHIEVEMENT_SYSTEM.md) - Comprehensive achievement tracking with 19 unique achievements
- [Adaptive Music System](docs/systems/ADAPTIVE_MUSIC_SYSTEM.md) - Dynamic music that responds to gameplay
- [Animation System](docs/systems/ANIMATION_SYSTEM.md) - Frame-based sprite animations
- [Combat System](docs/systems/COMBAT_SYSTEM.md) - Player attacks, damage, AI behaviors
- [Door System](docs/systems/DOOR_SYSTEM.md) - Ability-gated progression mechanics
- [Enemy Animation System](docs/systems/ENEMY_ANIMATION_SYSTEM.md) - Procedural enemy animations
- [Item System](docs/systems/ITEM_SYSTEM.md) - Collectible items and inventory
- [Particle System](docs/systems/PARTICLE_SYSTEM.md) - Visual effects for combat and movement
- [Save System](docs/systems/SAVE_SYSTEM.md) - Persistent game state and save slots
- [Room Transitions](docs/systems/ROOM_TRANSITIONS.md) - Seamless room-to-room movement

### Additional Documentation
- [Rendering System](docs/RENDERING.md) - Ebiten-based rendering and graphics
- [Build Notes](docs/BUILD_NOTES.md) - Build and compilation information

## 📚 Technical References

- Procedural Content Generation in Games (PCG Book)
- Audio synthesis: Karplus-Strong algorithm, FM synthesis
- Cellular Automata: Conway's Game of Life variants
- Graph theory: Dijkstra's algorithm for pathfinding
- Music theory: Circle of fifths, chord progressions

---

**Note**: This is a technical demonstration of procedural generation with a fully playable game engine. Use `./vania --play` to launch the game.