package message

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Linkf(format string, a ...interface{}) {
	fmt.Printf("%s", lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#0077FF")). // blue link #0077FF
		Render(fmt.Sprintf(format, a...)))
}

func Link(msg string) {
	Linkf("%s", fmt.Sprint(msg))
	fmt.Printf("\n")
}
