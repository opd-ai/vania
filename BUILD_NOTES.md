# Build Notes

## CI/Headless Environment Limitations

The VANIA rendering system uses Ebiten, which requires graphics libraries (X11 on Linux, native on macOS/Windows). This means:

### ✅ What Works in CI/Headless Environments
- All existing tests (PCG, graphics generation, audio synthesis)
- Physics tests (no graphics required)
- Input/render data structure tests
- Code analysis and linting
- Security scanning

### ❌ What Requires Graphics Environment
- Building with Ebiten (`go build`)
- Running with rendering mode (`./vania --play`)
- Full integration tests with rendering

## Building Locally

To build and run the game with rendering on your local machine:

### Linux
```bash
# Install dependencies
sudo apt-get install gcc libc6-dev libgl1-mesa-dev libxcursor-dev \
  libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev \
  libasound2-dev pkg-config

# Build and run
go build -o vania ./cmd/game
./vania --seed 42 --play
```

### macOS
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Build and run
go build -o vania ./cmd/game
./vania --seed 42 --play
```

### Windows
```bash
# Install mingw-w64 (if needed)
# No additional graphics dependencies

# Build and run
go build -o vania.exe ./cmd/game
.\vania.exe --seed 42 --play
```

## Testing in CI

For CI environments without graphics:

```bash
# Run tests that don't require graphics
go test ./internal/pcg -v
go test ./internal/graphics -v
go test ./internal/audio -v
go test ./internal/physics -v
go test ./internal/narrative -v
go test ./internal/world -v
go test ./internal/entity -v

# Skip render/input/engine tests (require graphics context)
```

## Development Workflow

1. **Local Development**: Full build with graphics
2. **CI Testing**: Tests that don't require graphics
3. **Integration Testing**: On machines with graphics capabilities

## Alternative: Cross-Platform Build

For CI/CD pipelines, consider:
- Build on platform-specific runners (GitHub Actions: ubuntu-latest, macos-latest, windows-latest)
- Use Docker with X11 forwarding for Linux builds
- Separate build artifacts for each platform

See `.github/workflows/` for CI configuration examples.
