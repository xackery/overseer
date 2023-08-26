package dashboard

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/xackery/overseer/pkg/reporter"
)

var (
	green  = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	red    = lipgloss.AdaptiveColor{Light: "#F24C3D", Dark: "#FF8B92"}
	yellow = lipgloss.AdaptiveColor{Light: "#F2A93B", Dark: "#FFBE5C"}
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	titleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#43BF6D")). // green #8BE9FD
			Background(subtle)

	list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Height(8).
		Width(27)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render
)

func renderState(state reporter.AppState, msg string) string {
	switch state {
	case reporter.AppStateRunning:
		return lipgloss.NewStyle().SetString("✔ ").
			Foreground(green).
			PaddingRight(1).
			String() + lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(msg+" Running")
	case reporter.AppStateSleeping:
		return lipgloss.NewStyle().SetString("💤").
			Foreground(subtle).
			PaddingRight(1).
			String() + lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(msg+" Sleeping")
	case reporter.AppStateStopped:
		return lipgloss.NewStyle().SetString("✖ "). //crossmark
								Foreground(red).
								PaddingRight(1).
								String() + lipgloss.NewStyle().
								Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
								Render(msg+" Stopped")
	case reporter.AppStateStarting:
		return lipgloss.NewStyle().SetString("⚡"). //zap

								Foreground(yellow).
								PaddingRight(1).
								String() + lipgloss.NewStyle().
								Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
								Render(msg+" Starting")
	case reporter.AppStateErroring:
		return lipgloss.NewStyle().SetString("⚠ "). //warning
								Foreground(red).
								PaddingRight(1).
								String() + lipgloss.NewStyle().
								Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
								Render(msg+" Erroring")
	case reporter.AppStateRestarting:
		return lipgloss.NewStyle().SetString("🔄"). //refresh

								Foreground(yellow).
								PaddingRight(1).
								String() + lipgloss.NewStyle().
								Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
								Render(msg+" Restarting")
	default:
		return lipgloss.NewStyle().SetString("? ").
			Foreground(yellow).
			PaddingRight(1).
			String() + lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(msg+" Unknown")
	}
}

func renderIcon(icon string, msg string) string {
	return lipgloss.NewStyle().SetString(icon).
		Foreground(subtle).
		PaddingRight(1).
		String() + lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
		Render(msg)
}
