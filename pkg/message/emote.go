package message

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Emotef(emote string, format string, a ...interface{}) {
	fmt.Printf(lipgloss.NewStyle().SetString(emote).
		Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
		PaddingRight(1).
		String()+" "+format, a...)
}

func Emote(emote string, msg string) {
	Emotef(emote, "%s\n", fmt.Sprint(msg))
}
