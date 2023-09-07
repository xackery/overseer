package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
	"github.com/xackery/overseer/pkg/signal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xackery/overseer/pkg/dashboard"
	"github.com/xackery/overseer/pkg/manager"
	"github.com/xackery/overseer/pkg/reporter"
)

var (
	Version       = "0.0.0"
	isInitialized = false
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6ZmFsc2UsImZvbnRGYW1pbHkiOiJBYmhheWEgTGlicmUiLCJmb250UG9zaXRpb24iOiI2NSIsImZvbnRTaXplIjoiNDUiLCJmb250V2VpZ2h0Ijo2MDAsImltYWdlIjoiIiwiaW1hZ2VNYXNrIjoiIiwiaW1hZ2VTaXplIjoiNDAiLCJzaGFwZSI6ImNpcmNsZSIsInRleHQiOiLwn5GB77iPIn0
func main() {
	start := time.Now()
	err := run()
	if isInitialized {
		fmt.Print("\033[H\033[2J") // clear screen
	}
	if err != nil {
		message.Badf("Overseer failed: %s\n", err)
		operation.Exit(1)
	}
	message.OKf("Overseer exited after %0.2f seconds\n", time.Since(start).Seconds())
}

func run() error {
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}

	err = parseManager(config)
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	isInitialized = true
	p := tea.NewProgram(dashboard.New(Version))
	go func() {
		for {
			select {
			case <-reporter.SendUpdateChan:
				p.Send(dashboard.RefreshRequest{})
			case <-signal.Ctx().Done():
				return
			case <-time.After(5 * time.Second):
				p.Send(dashboard.RefreshRequest{})
			}
		}

	}()

	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}

func parseManager(config *config.OverseerConfiguration) error {
	var err error
	winExt := ".exe"
	if runtime.GOOS != "windows" {
		winExt = ""
	}

	setupType := manager.SetupDefault
	switch config.Setup {
	case "docker":
		setupType = manager.SetupDocker
		err = manager.InitializeDockerNetwork(config.DockerNetwork)
		if err != nil {
			return fmt.Errorf("initialize docker network: %w", err)
		}
	case "default":
		setupType = manager.SetupDefault
	}

	for i := 0; i < config.ZoneCount; i++ {
		manager.Manage(setupType, fmt.Sprintf("zone%d", i), "zone"+winExt, fmt.Sprintf("%d", i))
	}

	manager.Manage(setupType, "world", "world"+winExt)
	manager.Manage(setupType, "ucs", "ucs"+winExt)
	manager.Manage(setupType, "queryserv", "queryserv"+winExt)
	manager.Manage(setupType, "loginserver", "loginserver"+winExt)
	//for k, v := range config.Other {
	//	manager.Manage(k, v)
	//}
	return nil
}
