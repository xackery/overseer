package message

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func Bannerf(format string, a ...interface{}) {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	titleWidth := 77
	if physicalWidth < titleWidth {
		titleWidth = physicalWidth
	}

	fmt.Printf("%s", lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#43BF6D")). // dark green #43BF6D
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}).Width(titleWidth).Render(fmt.Sprintf(format, a...)))
}

func Banner(msg string) {
	Bannerf("%s", fmt.Sprint(msg))
	fmt.Printf("\n")
}
