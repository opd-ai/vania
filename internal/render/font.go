// Package render provides bitmap font rendering utilities shared across the
// render and menu subsystems.  GetCharPattern is the canonical source for
// pixel-art character bitmaps used throughout the UI.
package render

// GetCharPattern returns an 8×12 bitmap pattern for the given rune.
// width and height control the pattern grid dimensions; a minimum of 8×12 is
// expected for correct rendering.  The same pattern is used by both the HUD
// text renderer and the menu system.
func GetCharPattern(char rune, width, height int) [][]bool {
	// Initialize empty pattern
	pattern := make([][]bool, height)
	for i := range pattern {
		pattern[i] = make([]bool, width)
	}

	// Simple 8x12 bitmap patterns for basic ASCII characters
	switch char {
	case 'A', 'a':
		// A pattern
		pattern[1][3] = true
		pattern[1][4] = true
		pattern[2][2] = true
		pattern[2][5] = true
		pattern[3][2] = true
		pattern[3][5] = true
		pattern[4][1] = true
		pattern[4][6] = true
		pattern[5][1] = true
		pattern[5][6] = true
		pattern[6][1] = true
		pattern[6][2] = true
		pattern[6][5] = true
		pattern[6][6] = true
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][1] = true
		pattern[8][6] = true

	case 'B', 'b':
		// B pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		pattern[1][2] = true
		pattern[1][3] = true
		pattern[1][4] = true
		pattern[1][5] = true
		pattern[2][6] = true
		pattern[3][6] = true
		pattern[4][6] = true
		pattern[5][2] = true
		pattern[5][3] = true
		pattern[5][4] = true
		pattern[5][5] = true
		pattern[6][6] = true
		pattern[7][6] = true
		pattern[8][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'C', 'c':
		// C pattern
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[3][6] = true
		for y := 4; y < 7; y++ {
			pattern[y][1] = true
		}
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'E', 'e':
		// E pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		for x := 1; x < 7; x++ {
			pattern[1][x] = true
			pattern[5][x] = true
			pattern[8][x] = true
		}

	case 'G', 'g':
		// G pattern
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[4][1] = true
		pattern[5][1] = true
		pattern[6][1] = true
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true
		pattern[5][4] = true
		pattern[5][5] = true
		pattern[5][6] = true
		pattern[6][6] = true

	case 'L', 'l':
		// L pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		for x := 1; x < 7; x++ {
			pattern[8][x] = true
		}

	case 'M', 'm':
		// M pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[2][2] = true
		pattern[2][5] = true
		pattern[3][3] = true
		pattern[3][4] = true

	case 'N', 'n':
		// N pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[2][2] = true
		pattern[3][2] = true
		pattern[4][3] = true
		pattern[5][4] = true
		pattern[6][5] = true
		pattern[7][5] = true

	case 'O', 'o':
		// O pattern
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[3][6] = true
		for y := 4; y < 7; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[7][1] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'P', 'p':
		// P pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		for x := 2; x < 6; x++ {
			pattern[1][x] = true
			pattern[5][x] = true
		}
		pattern[2][6] = true
		pattern[3][6] = true
		pattern[4][6] = true

	case 'Q', 'q':
		// Q pattern (O with tail)
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[3][6] = true
		for y := 4; y < 7; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[6][4] = true // Inner diagonal
		pattern[7][1] = true
		pattern[7][5] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][6] = true

	case 'R', 'r':
		// R pattern
		for y := 1; y < 9; y++ {
			pattern[y][1] = true
		}
		for x := 2; x < 6; x++ {
			pattern[1][x] = true
			pattern[5][x] = true
		}
		pattern[2][6] = true
		pattern[3][6] = true
		pattern[4][6] = true
		pattern[6][4] = true
		pattern[7][5] = true
		pattern[8][6] = true

	case 'S', 's':
		// S pattern
		pattern[2][2] = true
		pattern[2][3] = true
		pattern[2][4] = true
		pattern[2][5] = true
		pattern[3][1] = true
		pattern[4][1] = true
		pattern[5][2] = true
		pattern[5][3] = true
		pattern[5][4] = true
		pattern[5][5] = true
		pattern[6][6] = true
		pattern[7][6] = true
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'T', 't':
		// T pattern
		for x := 1; x < 7; x++ {
			pattern[1][x] = true
		}
		for y := 1; y < 9; y++ {
			pattern[y][3] = true
		}

	case 'U', 'u':
		// U pattern
		for y := 1; y < 8; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[8][2] = true
		pattern[8][3] = true
		pattern[8][4] = true
		pattern[8][5] = true

	case 'V', 'v':
		// V pattern
		for y := 1; y < 7; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[7][2] = true
		pattern[7][5] = true
		pattern[8][3] = true
		pattern[8][4] = true

	case 'W', 'w':
		// W pattern
		for y := 1; y < 8; y++ {
			pattern[y][1] = true
			pattern[y][6] = true
		}
		pattern[6][3] = true
		pattern[6][4] = true
		pattern[7][2] = true
		pattern[7][3] = true
		pattern[7][4] = true
		pattern[7][5] = true

	case 'X', 'x':
		// X pattern
		pattern[1][1] = true
		pattern[1][6] = true
		pattern[2][2] = true
		pattern[2][5] = true
		pattern[3][3] = true
		pattern[3][4] = true
		pattern[4][3] = true
		pattern[4][4] = true
		pattern[5][2] = true
		pattern[5][5] = true
		pattern[6][1] = true
		pattern[6][6] = true

	case 'Y', 'y':
		// Y pattern
		pattern[1][1] = true
		pattern[1][6] = true
		pattern[2][2] = true
		pattern[2][5] = true
		pattern[3][3] = true
		pattern[3][4] = true
		for y := 4; y < 9; y++ {
			pattern[y][3] = true
		}

	case 'Z', 'z':
		// Z pattern
		for x := 1; x < 7; x++ {
			pattern[1][x] = true
			pattern[8][x] = true
		}
		pattern[2][6] = true
		pattern[3][5] = true
		pattern[4][4] = true
		pattern[5][3] = true
		pattern[6][2] = true
		pattern[7][1] = true

	case ' ':
		// Space - already empty

	case ':':
		// Colon
		pattern[3][3] = true
		pattern[6][3] = true

	case '-':
		// Hyphen
		pattern[5][2] = true
		pattern[5][3] = true
		pattern[5][4] = true
		pattern[5][5] = true

	case '.':
		// Period
		pattern[8][3] = true

	case '!':
		// Exclamation
		for y := 1; y < 7; y++ {
			pattern[y][3] = true
		}
		pattern[8][3] = true

	case '?':
		// Question mark
		pattern[1][2] = true
		pattern[1][3] = true
		pattern[1][4] = true
		pattern[1][5] = true
		pattern[2][6] = true
		pattern[3][6] = true
		pattern[4][4] = true
		pattern[4][5] = true
		pattern[5][3] = true
		pattern[8][3] = true

	case '(':
		// Left parenthesis
		pattern[2][4] = true
		pattern[3][3] = true
		for y := 4; y < 7; y++ {
			pattern[y][2] = true
		}
		pattern[7][3] = true
		pattern[8][4] = true

	case ')':
		// Right parenthesis
		pattern[2][3] = true
		pattern[3][4] = true
		for y := 4; y < 7; y++ {
			pattern[y][5] = true
		}
		pattern[7][4] = true
		pattern[8][3] = true

	case '/':
		// Forward slash
		pattern[1][6] = true
		pattern[2][5] = true
		pattern[3][5] = true
		pattern[4][4] = true
		pattern[5][3] = true
		pattern[6][3] = true
		pattern[7][2] = true
		pattern[8][1] = true

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Numbers
		switch char {
		case '0':
			pattern[2][2] = true
			pattern[2][3] = true
			pattern[2][4] = true
			pattern[2][5] = true
			for y := 3; y < 8; y++ {
				pattern[y][1] = true
				pattern[y][6] = true
			}
			pattern[8][2] = true
			pattern[8][3] = true
			pattern[8][4] = true
			pattern[8][5] = true
		case '1':
			pattern[2][3] = true
			pattern[3][2] = true
			pattern[3][3] = true
			for y := 4; y < 9; y++ {
				pattern[y][3] = true
			}
			for x := 1; x < 7; x++ {
				pattern[8][x] = true
			}
		case '2':
			pattern[2][2] = true
			pattern[2][3] = true
			pattern[2][4] = true
			pattern[2][5] = true
			pattern[3][6] = true
			pattern[4][6] = true
			pattern[5][5] = true
			pattern[6][4] = true
			pattern[7][3] = true
			pattern[8][2] = true
			for x := 1; x < 7; x++ {
				pattern[8][x] = true
			}
		case '3':
			for x := 2; x < 6; x++ {
				pattern[2][x] = true
				pattern[5][x] = true
				pattern[8][x] = true
			}
			pattern[3][6] = true
			pattern[4][6] = true
			pattern[6][6] = true
			pattern[7][6] = true
		case '4':
			for y := 2; y < 6; y++ {
				pattern[y][1] = true
				pattern[y][5] = true
			}
			for x := 1; x < 7; x++ {
				pattern[5][x] = true
			}
			for y := 6; y < 9; y++ {
				pattern[y][5] = true
			}
		case '5':
			for x := 1; x < 7; x++ {
				pattern[2][x] = true
				pattern[5][x] = true
			}
			for y := 3; y < 5; y++ {
				pattern[y][1] = true
			}
			pattern[6][6] = true
			pattern[7][6] = true
			for x := 2; x < 6; x++ {
				pattern[8][x] = true
			}
		case '6':
			pattern[2][3] = true
			pattern[2][4] = true
			pattern[3][2] = true
			for y := 4; y < 8; y++ {
				pattern[y][1] = true
			}
			pattern[5][2] = true
			pattern[5][3] = true
			pattern[5][4] = true
			pattern[5][5] = true
			pattern[6][6] = true
			pattern[7][6] = true
			pattern[8][2] = true
			pattern[8][3] = true
			pattern[8][4] = true
			pattern[8][5] = true
		case '7':
			for x := 1; x < 7; x++ {
				pattern[2][x] = true
			}
			pattern[3][6] = true
			pattern[4][5] = true
			pattern[5][4] = true
			pattern[6][3] = true
			pattern[7][3] = true
			pattern[8][3] = true
		case '8':
			for x := 2; x < 6; x++ {
				pattern[2][x] = true
				pattern[5][x] = true
				pattern[8][x] = true
			}
			pattern[3][1] = true
			pattern[3][6] = true
			pattern[4][1] = true
			pattern[4][6] = true
			pattern[6][1] = true
			pattern[6][6] = true
			pattern[7][1] = true
			pattern[7][6] = true
		case '9':
			pattern[2][2] = true
			pattern[2][3] = true
			pattern[2][4] = true
			pattern[2][5] = true
			pattern[3][1] = true
			pattern[3][6] = true
			pattern[4][1] = true
			pattern[4][6] = true
			pattern[5][2] = true
			pattern[5][3] = true
			pattern[5][4] = true
			pattern[5][6] = true
			pattern[6][5] = true
			pattern[7][4] = true
			pattern[8][3] = true
			pattern[8][4] = true
		}

	default:
		// Unknown character - draw a simple box
		for y := 2; y < 9; y++ {
			pattern[y][2] = true
			pattern[y][5] = true
		}
		for x := 2; x < 6; x++ {
			pattern[2][x] = true
			pattern[8][x] = true
		}
	}

	return pattern
}
