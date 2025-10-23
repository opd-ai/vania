package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/vania/internal/engine"
	"github.com/opd-ai/vania/internal/menu"
)

// GameApp represents the main application with menu integration
type GameApp struct {
	menuManager *menu.MenuManager
	gameRunner  *engine.GameRunner
	currentGame *engine.Game
	inMenu      bool

	// Command line options
	directPlay bool
	fixedSeed  int64
}

// NewGameApp creates a new game application
func NewGameApp(directPlay bool, fixedSeed int64) *GameApp {
	app := &GameApp{
		menuManager: menu.NewMenuManager(),
		inMenu:      !directPlay,
		directPlay:  directPlay,
		fixedSeed:   fixedSeed,
	}

	// Set up menu callbacks
	app.menuManager.SetCallbacks(
		app.onNewGame,    // New game
		app.onLoadGame,   // Load game
		app.onSettings,   // Settings
		app.onQuitGame,   // Quit
		app.onResumeGame, // Resume
	)

	// If direct play mode, start game immediately
	if directPlay {
		seed := fixedSeed
		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		if err := app.startGame(seed); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting game: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Show main menu
		app.menuManager.ShowMainMenu()
	}

	return app
}

// Update implements ebiten.Game interface
func (app *GameApp) Update() error {
	if app.inMenu {
		return app.menuManager.Update()
	} else if app.gameRunner != nil {
		err := app.gameRunner.Update()

		// Check if player died (health <= 0)
		if app.currentGame != nil && app.currentGame.Player.Health <= 0 {
			app.showGameOver()
			return nil
		}

		return err
	}
	return nil
}

// Draw implements ebiten.Game interface
func (app *GameApp) Draw(screen *ebiten.Image) {
	if app.inMenu {
		app.menuManager.Draw(screen)
	} else if app.gameRunner != nil {
		app.gameRunner.Draw(screen)
	}
}

// Layout implements ebiten.Game interface
func (app *GameApp) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 960, 640
}

// onNewGame handles new game creation
func (app *GameApp) onNewGame(seed int64) error {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return app.startGame(seed)
}

// onLoadGame handles game loading
func (app *GameApp) onLoadGame(slot int) error {
	if app.gameRunner == nil {
		// Need to create a dummy game first for save system to work
		if err := app.startGame(42); err != nil {
			return err
		}
	}

	// Load the game from save slot
	if err := app.gameRunner.LoadGame(slot); err != nil {
		return fmt.Errorf("failed to load game: %v", err)
	}

	app.inMenu = false
	app.menuManager.Hide()
	return nil
}

// onSettings handles settings menu
func (app *GameApp) onSettings() error {
	app.menuManager.ShowSettingsMenu()
	return nil
}

// onQuitGame handles game quit
func (app *GameApp) onQuitGame() error {
	return ebiten.Termination
}

// onResumeGame handles game resume from pause
func (app *GameApp) onResumeGame() error {
	app.inMenu = false
	app.menuManager.Hide()
	return nil
}

// startGame creates and starts a new game
func (app *GameApp) startGame(seed int64) error {
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                        ║")
	fmt.Println("║         VANIA - Procedural Metroidvania                ║")
	fmt.Println("║         Pure Go Procedural Generation Demo             ║")
	fmt.Println("║                                                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Master Seed: %d\n", seed)
	fmt.Println("Generating game world...")

	// Create game generator
	generator := engine.NewGameGenerator(seed)

	// Generate complete game
	game, err := generator.GenerateCompleteGame()
	if err != nil {
		return fmt.Errorf("error generating game: %v", err)
	}

	fmt.Println("Generation complete! Starting game...")

	// Create game runner
	app.currentGame = game
	app.gameRunner = engine.NewGameRunner(game)

	// Switch to game mode
	app.inMenu = false
	app.menuManager.Hide()

	return nil
}

// showGameOver displays game over screen
func (app *GameApp) showGameOver() {
	app.inMenu = true
	app.menuManager.ShowGameOverMenu()
}

// showPauseMenu displays pause menu
func (app *GameApp) showPauseMenu() {
	app.inMenu = true
	app.menuManager.ShowPauseMenu()
}

// Run starts the application
func (app *GameApp) Run() error {
	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("VANIA - Procedural Metroidvania")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(app)
}

func main() {
	// Parse command line arguments
	seedFlag := flag.Int64("seed", 0, "Master seed for generation (0 = use timestamp)")
	playFlag := flag.Bool("play", false, "Launch the game with rendering (default: show main menu)")
	noMenuFlag := flag.Bool("no-menu", false, "Skip menu and go directly to gameplay")
	statsOnlyFlag := flag.Bool("stats-only", false, "Generate and show stats only (original behavior)")
	flag.Parse()

	// Handle legacy stats-only mode
	if *statsOnlyFlag {
		runStatsOnlyMode(*seedFlag)
		return
	}

	// Determine if we should skip menus
	directPlay := *playFlag && *noMenuFlag

	// Create and run the application
	app := NewGameApp(directPlay, *seedFlag)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Game error: %v\n", err)
		os.Exit(1)
	}
}

// runStatsOnlyMode provides the original stats-only behavior
func runStatsOnlyMode(seedFlag int64) {
	var masterSeed int64
	if seedFlag == 0 {
		masterSeed = time.Now().UnixNano()
	} else {
		masterSeed = seedFlag
	}

	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                        ║")
	fmt.Println("║         VANIA - Procedural Metroidvania                ║")
	fmt.Println("║         Pure Go Procedural Generation Demo             ║")
	fmt.Println("║                                                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Master Seed: %d\n", masterSeed)
	fmt.Println()
	fmt.Println("Generating game world...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Create game generator
	generator := engine.NewGameGenerator(masterSeed)

	// Generate complete game
	game, err := generator.GenerateCompleteGame()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating game: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Generation Complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Display game statistics
	displayGameStats(game)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Game ready to play!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Display achievement info
	if game.Achievements != nil {
		fmt.Println("🏆 ACHIEVEMENTS")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Total Achievements: %d\n", len(game.Achievements.GetAllAchievements()))
		fmt.Printf("  Max Points:         %d\n", game.Achievements.GetMaxPoints())
		fmt.Println("  Play the game to unlock achievements!")
		fmt.Println()
	}

	fmt.Println("(Use --play flag to launch the game with rendering)")
	game.Run()
}

func displayGameStats(game *engine.Game) {
	fmt.Println("📖 NARRATIVE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Theme:              %s\n", game.Narrative.Theme)
	fmt.Printf("  Mood:               %s\n", game.Narrative.Mood)
	fmt.Printf("  Civilization:       %s\n", game.Narrative.CivilizationType)
	fmt.Printf("  Catastrophe:        %s\n", game.Narrative.Catastrophe)
	fmt.Printf("  Player Motivation:  %s\n", game.Narrative.PlayerMotivation)
	fmt.Println()

	fmt.Println("🌍 WORLD")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Total Rooms:        %d\n", len(game.World.Rooms))
	fmt.Printf("  Boss Rooms:         %d\n", len(game.World.BossRooms))
	fmt.Printf("  Biomes:             %d\n", len(game.World.Biomes))
	fmt.Printf("  Grid Size:          %dx%d\n", game.World.Width, game.World.Height)
	fmt.Println()
	fmt.Println("  Biome List:")
	for i, biome := range game.World.Biomes {
		fmt.Printf("    %d. %s (Danger: %d, Temp: %d°C)\n",
			i+1, biome.Name, biome.DangerLevel, biome.Temperature)
	}
	fmt.Println()

	fmt.Println("👾 ENTITIES")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Regular Enemies:    %d\n", len(game.Entities))
	fmt.Printf("  Boss Enemies:       %d\n", len(game.Bosses))
	fmt.Printf("  Items:              %d\n", len(game.Items))
	fmt.Printf("  Abilities:          %d\n", len(game.Abilities))
	fmt.Println()

	if len(game.Bosses) > 0 {
		fmt.Println("  Boss Preview:")
		for i, boss := range game.Bosses {
			if i >= 3 {
				break
			}
			fmt.Printf("    - %s (HP: %d, Phases: %d)\n",
				boss.Name, boss.Health, len(boss.Phases))
		}
		fmt.Println()
	}

	if len(game.Abilities) > 0 {
		fmt.Println("  Ability Progression:")
		for i, ability := range game.Abilities {
			if i >= 5 {
				fmt.Println("    ...")
				break
			}
			fmt.Printf("    %d. %s (%s)\n",
				ability.UnlockOrder+1, ability.Name, ability.Description)
		}
		fmt.Println()
	}

	fmt.Println("🎨 GRAPHICS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Sprites Generated:  %d\n", len(game.Graphics.Sprites))
	fmt.Printf("  Tilesets:           %d\n", len(game.Graphics.Tilesets))
	fmt.Println("  All graphics procedurally generated at runtime!")
	fmt.Println()

	fmt.Println("🎵 AUDIO")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Sound Effects:      %d\n", len(game.Audio.Sounds))
	fmt.Printf("  Music Tracks:       %d\n", len(game.Audio.Music))
	fmt.Println("  All audio synthesized at runtime!")
	fmt.Println()

	fmt.Println("🏛️ FACTIONS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for i, faction := range game.Narrative.Factions {
		fmt.Printf("  %d. %s\n", i+1, faction.Name)
		fmt.Printf("     %s (%s)\n", faction.Description, faction.Relationship)
	}
}

func displayAchievementSummary(game *engine.Game) {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🏆 ACHIEVEMENT SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	unlocked := game.Achievements.GetUnlockedAchievements()
	fmt.Printf("  Unlocked:           %d / %d (%.1f%%)\n",
		len(unlocked),
		len(game.Achievements.GetAllAchievements()),
		game.Achievements.GetCompletionPercentage())
	fmt.Printf("  Points Earned:      %d / %d\n",
		game.Achievements.GetTotalPoints(),
		game.Achievements.GetMaxPoints())
	fmt.Println()

	if len(unlocked) > 0 {
		fmt.Println("  Unlocked Achievements:")
		for i, u := range unlocked {
			if i >= 5 {
				fmt.Printf("    ... and %d more\n", len(unlocked)-5)
				break
			}
			achievement := game.Achievements.GetAchievement(u.AchievementID)
			if achievement != nil {
				fmt.Printf("    ✓ %s - %s (%d pts)\n",
					achievement.Name,
					achievement.Description,
					achievement.Points)
			}
		}
		fmt.Println()
	}

	// Show some progress on locked achievements
	fmt.Println("  In Progress:")
	shown := 0
	for _, achievement := range game.Achievements.GetAllAchievements() {
		if game.Achievements.IsUnlocked(achievement.ID) || achievement.Hidden {
			continue
		}

		progress := game.Achievements.GetProgress(achievement.ID)
		if progress != nil && progress.Progress > 0 {
			fmt.Printf("    ⋯ %s - %.0f%%\n",
				achievement.Name,
				progress.Progress*100)
			shown++
			if shown >= 3 {
				break
			}
		}
	}

	if shown == 0 {
		fmt.Println("    Keep playing to unlock more achievements!")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
