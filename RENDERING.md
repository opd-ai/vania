# Building and Running VANIA with Rendering

## Prerequisites

To build and run VANIA with the graphical rendering engine, you need:

- Go 1.21 or higher
- C compiler (gcc or clang)
- Graphics libraries (for Linux: X11, for macOS: included, for Windows: included)

### Installing Dependencies

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install -y gcc libc6-dev libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libasound2-dev pkg-config
```

#### Linux (Fedora/RHEL)
```bash
sudo dnf install -y gcc mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
```

#### macOS
```bash
# Xcode command line tools (includes required libraries)
xcode-select --install
```

#### Windows
No additional dependencies needed beyond Go and a C compiler (e.g., mingw-w64)

## Building

```bash
# Clone the repository
git clone https://github.com/opd-ai/vania.git
cd vania

# Install dependencies
go mod tidy

# Build the game
go build -o vania ./cmd/game
```

## Running

### Mode 1: Generation Only (Original Behavior)

Generate game content and display statistics without rendering:

```bash
# Random seed (uses current timestamp)
./vania

# Specific seed
./vania --seed 42
```

### Mode 2: Play with Rendering (NEW!)

Launch the full game with graphical rendering:

```bash
# Random seed with rendering
./vania --play

# Specific seed with rendering
./vania --seed 42 --play
```

## Controls

- **Movement**: WASD or Arrow Keys
- **Jump**: Space, W, or Up Arrow
- **Dash**: K or X (requires dash ability)
- **Attack**: J or Z
- **Pause**: P or Escape
- **Quit**: Ctrl+Q

## Features in Rendering Mode

- **Procedurally generated tilesets** rendered in real-time
- **Player sprite** with physics-based movement
- **Platform collision detection** with gravity and jumping
- **Health bar** and ability indicators
- **Camera system** that follows the player
- **Biome-specific visual themes**
- **Debug information** showing FPS, position, and state

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/physics -v
go test ./internal/graphics -v
go test ./internal/audio -v
```

Note: Tests for rendering, input, and engine runner require a graphical environment and may not run in headless CI environments.

## Architecture

The new rendering system adds three key packages:

1. **`internal/render/`** - Ebiten-based rendering system
   - Camera management
   - Tile and sprite rendering
   - UI rendering (health, abilities)

2. **`internal/input/`** - Input handling
   - Keyboard input processing
   - Action mapping
   - Input state management

3. **`internal/physics/`** - Physics and collision
   - Gravity simulation
   - Platform collision detection
   - Player movement mechanics
   - Jump, dash, and wall-jump abilities

4. **`internal/engine/runner.go`** - Game loop integration
   - Ebiten game interface implementation
   - Update/Draw cycle
   - Integration with procedural generation

## Performance

The rendering system runs at 60 FPS with:
- Generation time: ~0.3 seconds
- Runtime performance: 60 FPS on modern hardware
- Memory usage: ~50-100 MB

## Troubleshooting

### Build Errors on Linux

If you get errors about missing X11 headers:
```bash
sudo apt-get install libx11-dev
```

### Black Screen on Launch

Ensure your graphics drivers are up to date and OpenGL is supported.

### Low FPS

Check that vsync is enabled and your system isn't running in power-saving mode.

## Development Roadmap

✅ **Implemented**
- Core rendering system with Ebiten
- Player physics and movement
- Collision detection
- Input handling
- Camera system
- UI rendering

🚧 **In Progress**
- Enemy rendering and AI
- Combat system
- Room transitions
- Animation system

📋 **Planned**
- Particle effects
- Advanced lighting
- Minimap
- Save/load system
- Audio playback integration
