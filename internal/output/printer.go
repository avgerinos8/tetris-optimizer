package output

import (
	"fmt"
	"log/slog"
)

type colorCodes string

const (
	Reset colorCodes = "\033[0m"

	// Format: Background RGB + Foreground Black (0,0,0)
	Red          colorCodes = "\033[48;2;255;22;22;38;2;255;255;255m"
	Orange       colorCodes = "\033[48;2;234;138;0;38;2;0;0;0m"
	Amber        colorCodes = "\033[48;2;255;168;0;38;2;0;0;0m"
	Yellow       colorCodes = "\033[48;2;255;255;0;38;2;0;0;0m"
	Lime         colorCodes = "\033[48;2;0;184;48;38;2;255;255;255m"
	Green        colorCodes = "\033[48;2;0;255;42;38;2;0;0;0m"
	Emerald      colorCodes = "\033[48;2;0;159;94;38;2;255;255;255m"
	Cyan         colorCodes = "\033[48;2;0;255;255;38;2;0;0;0m"
	SkyBlue      colorCodes = "\033[48;2;0;190;255;38;2;0;0;0m"
	ElectricBlue colorCodes = "\033[48;2;0;82;255;38;2;255;255;255m"
	HotPink      colorCodes = "\033[48;2;187;0;79;38;2;255;255;255m"
	Magenta      colorCodes = "\033[48;2;255;0;255;38;2;255;255;255m"
	Purple       colorCodes = "\033[48;2;166;0;166;38;2;0;0;0m"
	Gray         colorCodes = "\033[48;2;112;112;112;38;2;0;0;0m"
	Gold         colorCodes = "\033[48;2;163;133;35;38;2;0;0;0m"
)

// ── public entry point ─────────────────────────────────────────────────────

// PrintSolution renders the solved board to stdout.
// Each cell is the letter of the piece that covers it, or '.' if empty.
func PrintSolution(winningboard [][]rune, color bool, colorOnly bool) {
	// Map letters A-Z to our defined colors
	palette := []colorCodes{
		Red, Orange, ElectricBlue, Magenta, Yellow, Cyan, Purple, Lime, Green, Emerald,
		SkyBlue, HotPink, Amber, Gold, Gray,
	}

	for _, row := range winningboard {
		var rowLog string

		for _, cell := range row {
			char := "."
			if cell != 0 && cell != '.' {
				char = string(cell)
			}

			if char == "." {
				if color {
					fmt.Print("░░░")
				} else if colorOnly {
					fmt.Print("   ")
				} else {
					fmt.Print(".")
				}
			} else {
				if color {
					colorIdx := int(cell-'A') % len(palette)
					fmt.Printf("%s %s %s", palette[colorIdx], char, Reset)
				} else if colorOnly {
					colorIdx := int(cell-'A') % len(palette)
					fmt.Printf("%s   %s", palette[colorIdx], Reset)
				} else {
					fmt.Printf("%s", char)
				}
			}

			rowLog += char + " "
		}

		slog.Info(rowLog)
		fmt.Println()
	}
}
