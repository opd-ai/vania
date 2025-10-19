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
	playFlag := flag.Bool("play", false, "Launch the game with rendering (default: just generate and show stats)")
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

	// Run the game
	if *playFlag {
		// Launch with rendering
		fmt.Println("Launching game with rendering...")
		fmt.Println("Controls: WASD/Arrows=Move, Space=Jump, K=Dash, P=Pause, Ctrl+Q=Quit")
		fmt.Println()
		
		runner := engine.NewGameRunner(game)
		if err := runner.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
			os.Exit(1)
		}
		
		// Display achievement progress after game ends
		if game.Achievements != nil {
			displayAchievementSummary(game)
		}
	} else {
		// Just show stats (original behavior)
		fmt.Println("(Use --play flag to launch the game with rendering)")
		game.Run()
	}
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
