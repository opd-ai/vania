package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/opd-ai/vania/internal/engine"
)

func main() {
	// Parse command line arguments
	seedFlag := flag.Int64("seed", 0, "Master seed for generation (0 = use timestamp)")
	flag.Parse()

	// Determine seed
	var masterSeed int64
	if *seedFlag == 0 {
		masterSeed = time.Now().UnixNano()
	} else {
		masterSeed = *seedFlag
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
	fmt.Println("(Full game engine implementation in progress)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Run the game (currently just prints stats)
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
