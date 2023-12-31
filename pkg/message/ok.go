package message

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/xackery/overseer/pkg/slog"
)

var (
	isOK bool
)

func OKReset() {
	isOK = true
}

func IsOK() bool {
	return isOK
}

func OKf(format string, a ...interface{}) {
	fmt.Printf(lipgloss.NewStyle().SetString("✅").
		Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
		PaddingRight(1).
		String()+format, a...)
}

func OK(msg string) {
	OKf("%s\n", fmt.Sprint(msg))
}

func Badf(format string, a ...interface{}) {
	fmt.Printf(lipgloss.NewStyle().SetString("❌").
		Foreground(lipgloss.AdaptiveColor{Light: "#FF5555", Dark: "#FF5555"}).
		PaddingRight(1).
		String()+format, a...)
	slog.Printf(format, a...)
	isOK = false
}

func Bad(msg string) {
	Badf("%s\n", fmt.Sprint(msg))
}

func Skipf(format string, a ...interface{}) {
	fmt.Printf(lipgloss.NewStyle().SetString("⏩").
		Foreground(lipgloss.AdaptiveColor{Light: "#FF5555", Dark: "#FF5555"}).
		PaddingRight(1).
		String()+format, a...)
}

func Skip(msg string) {
	Skipf("%s\n", fmt.Sprint(msg))
}
