package dashboard

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xackery/overseer/share/flog"
	"github.com/xackery/overseer/share/reporter"
	"github.com/xackery/overseer/share/signal"
	"golang.org/x/term"
)

type Dashboard struct {
	version       string
	stateOrdering []string
}

// RefreshRequest is a message that tells the program to refresh the dashboard.
type RefreshRequest struct {
}

func (r RefreshRequest) String() string {
	return "RefreshRequest"
}

func New(version string) Dashboard {
	e := Dashboard{
		version: version,
	}
	state := reporter.AppStates()
	e.stateOrdering = []string{}
	builtinApps := []string{
		"world",
		"ucs",
		"zone",
	}

	for _, builtin := range builtinApps {
		if builtin == "zone" {
			continue
		}
		e.stateOrdering = append(e.stateOrdering, builtin)
	}

	flog.Printf("Number of services tracked: %d\n", len(e.stateOrdering))

	for appName := range state.States {
		isFound := false
		for _, builtinApp := range builtinApps {
			if !strings.HasPrefix(appName, builtinApp) {
				continue
			}
			isFound = true
			break

		}
		if !isFound {
			e.stateOrdering = append(e.stateOrdering, appName)
		}
	}
	return e
}

func (e Dashboard) Init() tea.Cmd {

	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (e Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RefreshRequest:
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			flog.Println("Received ctrl+c or q, exiting")
			signal.Cancel()
			signal.WaitWorker()
			return e, tea.Quit
		}
	}

	// Return the updated Model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return e, nil
}

func (e Dashboard) View() string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	titleWidth := 77
	if physicalWidth < titleWidth {
		titleWidth = physicalWidth
	}
	doc := strings.Builder{}

	height := 0
	select {
	case <-signal.Ctx().Done():
		doc.WriteString(titleStyle.Width(titleWidth).Render("Shutting down..."))
	default:
		doc.WriteString(titleStyle.Width(titleWidth).Render("Overseer v" + e.version))
	}
	doc.WriteString("\n\n")
	height += 2

	state := reporter.AppStates()

	renderStates := []string{
		listHeader("Services"),
	}
	for _, order := range e.stateOrdering {
		renderStates = append(renderStates, renderState(state.States[order], order))
	}

	//colHeight := physicalHeight - height - 2
	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		list.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				renderStates...,
			),
		),
		list.Copy().Width(16).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				listHeader(fmt.Sprintf("Zones (%d)", state.ZoneTotal)),
				renderState(reporter.AppStateRunning, fmt.Sprintf("%d", state.ZoneRunning)),
				renderState(reporter.AppStateStarting, fmt.Sprintf("%d", state.ZoneStarting)),
				renderState(reporter.AppStateSleeping, fmt.Sprintf("%d", state.ZoneSleeping)),
				renderState(reporter.AppStateErroring, fmt.Sprintf("%d", state.ZoneErroring)),
				renderState(reporter.AppStateRestarting, fmt.Sprintf("%d", state.ZoneRestarting)),
				renderState(reporter.AppStateStopped, fmt.Sprintf("%d", state.ZoneStopped)),
			),
		),
		list.Copy().Width(27).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				listHeader("Stats"),
				renderIcon("👤", fmt.Sprintf("%d Online", 3)), // person emoji: 👤
			),
		),
	))
	doc.WriteString("\n\n")

	return doc.String()
}
