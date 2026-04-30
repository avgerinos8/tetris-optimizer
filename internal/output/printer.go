package output

import (
	"fmt"
	"log/slog"
)

type colorCodes string

const (
	Reset colorCodes = "\033[0m"

	// Format: Background RGB + Foreground Black (0,0,0)
	Red          colorCodes = "\033[48;2;255;50;50;38;2;0;0;0m"
	Orange       colorCodes = "\033[48;2;255;150;0;38;2;0;0;0m"
	Amber        colorCodes = "\033[48;2;255;200;0;38;2;0;0;0m"
	Yellow       colorCodes = "\033[48;2;255;255;0;38;2;0;0;0m"
	Lime         colorCodes = "\033[48;2;180;255;0;38;2;0;0;0m"
	Green        colorCodes = "\033[48;2;50;255;50;38;2;0;0;0m"
	Emerald      colorCodes = "\033[48;2;0;255;150;38;2;0;0;0m"
	Cyan         colorCodes = "\033[48;2;0;255;255;38;2;0;0;0m"
	SkyBlue      colorCodes = "\033[48;2;0;190;255;38;2;0;0;0m"
	ElectricBlue colorCodes = "\033[48;2;100;150;255;38;2;0;0;0m"
	Violet       colorCodes = "\033[48;2;180;100;255;38;2;0;0;0m"
	Magenta      colorCodes = "\033[48;2;255;0;255;38;2;0;0;0m"
	HotPink      colorCodes = "\033[48;2;255;100;200;38;2;0;0;0m"
	Coral        colorCodes = "\033[48;2;255;127;80;38;2;0;0;0m"
	Gold         colorCodes = "\033[48;2;212;175;55;38;2;0;0;0m"
)

// ── public entry point ─────────────────────────────────────────────────────

// PrintSolution renders the solved board to stdout.
// Each cell is the letter of the piece that covers it, or '.' if empty.
func PrintSolution(winningboard [][]rune) {
	// Map letters A-Z to our defined colors
	palette := []colorCodes{
		Red, Orange, Amber, Yellow, Lime, Green, Emerald,
		Cyan, SkyBlue, ElectricBlue, Violet, Magenta, HotPink, Coral, Gold,
	}

	for _, row := range winningboard {
		var rowLog string

		for _, cell := range row {
			char := "."
			if cell != 0 && cell != '.' {
				char = string(cell)
			}

			if char == "." {
				fmt.Print(". ")
			} else {
				colorIdx := int(cell-'A') % len(palette)
				fmt.Printf("%s %s %s", palette[colorIdx], char, Reset)
			}

			rowLog += char + " "
		}

		slog.Info(rowLog)
		fmt.Println()
	}
}
