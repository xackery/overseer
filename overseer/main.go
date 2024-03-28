package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/flog"
	"github.com/xackery/overseer/pkg/gui"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
	"github.com/xackery/overseer/pkg/signal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xackery/overseer/pkg/dashboard"
	"github.com/xackery/overseer/pkg/manager"
	"github.com/xackery/overseer/pkg/reporter"
)

var (
	Version = "0.0.0"
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6ZmFsc2UsImZvbnRGYW1pbHkiOiJBYmhheWEgTGlicmUiLCJmb250UG9zaXRpb24iOiI2NSIsImZvbnRTaXplIjoiNDUiLCJmb250V2VpZ2h0Ijo2MDAsImltYWdlIjoiIiwiaW1hZ2VNYXNrIjoiIiwiaW1hZ2VTaXplIjoiNDAiLCJzaGFwZSI6ImNpcmNsZSIsInRleHQiOiLwn5GB77iPIn0
func main() {
	start := time.Now()
	err := run()
	//if isInitialized {
	//	fmt.Print("\033[H\033[2J") // clear screen
	//}
	if err != nil {
		flog.Printf("Overseer failed: %s\n", err)
		message.Badf("Overseer failed: %s\n", err)
		operation.Exit(1)
	}
	message.OKf("Overseer exited after %0.2f seconds\n", time.Since(start).Seconds())
	operation.Exit(0)
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, err := NewMainWindow(ctx, cancel, Version)
	if err != nil {
		return fmt.Errorf("new main window: %w", err)
	}
	gui.New(g)

	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}

	err = flog.New("overseer.log")
	if err != nil {
		return fmt.Errorf("new flog: %w", err)
	}
	defer flog.Close()

	err = parseManager(config)
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	if runtime.GOOS == "windows" {
		return runWindows(ctx, g)
	}
	time.Sleep(10 * time.Millisecond)
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

func parseManager(cfg *config.OverseerConfiguration) error {
	var err error
	winExt := ".exe"
	if runtime.GOOS != "windows" {
		winExt = ""
	}

	setupType := manager.SetupDefault
	switch cfg.Setup {
	case "docker":
		setupType = manager.SetupDocker
		err = manager.InitializeDockerNetwork(cfg.DockerNetwork)
		if err != nil {
			return fmt.Errorf("initialize docker network: %w", err)
		}
	case "default":
		setupType = manager.SetupDefault
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	fi, err := os.Stat(cwd + "/" + cfg.ServerPath)
	if err == nil {
		cfg.ServerPath = cwd + "/" + cfg.ServerPath
	} else {
		fi, err = os.Stat(cfg.ServerPath)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("server path is not a directory")
	}

	fi, err = os.Stat(cwd + "/" + cfg.BinPath)
	if err == nil {
		cfg.BinPath = cwd + "/" + cfg.BinPath
	} else {
		fi, err = os.Stat(cfg.BinPath)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
	}
	if !fi.IsDir() {
		return fmt.Errorf("bin path is not a directory")
	}

	wdPath, err := filepath.Abs(cfg.ServerPath)
	if err != nil {
		return fmt.Errorf("abs wdPath: %w", err)
	}

	exePath, err := filepath.Abs(cfg.BinPath)
	if err != nil {
		return fmt.Errorf("abs exePath: %w", err)
	}

	err = manager.Manage(setupType, "world", cfg.IsOverseerVerboseLog, wdPath, exePath, "world"+winExt)
	if err != nil {
		return fmt.Errorf("manage world: %w", err)
	}
	time.Sleep(50 * time.Millisecond)

	for i := 0; i < cfg.ZoneCount; i++ {
		err = manager.Manage(setupType, fmt.Sprintf("zone%d", i), cfg.IsOverseerVerboseLog, wdPath, exePath, "zone"+winExt)
		if err != nil {
			return fmt.Errorf("manage zone%d: %w", i, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)

	err = manager.Manage(setupType, "ucs", cfg.IsOverseerVerboseLog, wdPath, exePath, "ucs"+winExt)
	if err != nil {
		return fmt.Errorf("manage ucs: %w", err)
	}
	time.Sleep(50 * time.Millisecond)
	//manager.Manage(setupType, "queryserv", wdPath, exePath, "queryserv"+winExt)
	//time.Sleep(50 * time.Millisecond)

	//manager.Manage(setupType, "loginserver", wdPath, exePath, "loginserver"+winExt)
	//	time.Sleep(50 * time.Millisecond)

	for _, app := range cfg.Apps {
		nonExt := strings.TrimSuffix(app, filepath.Ext(app))
		err = manager.Manage(setupType, nonExt, cfg.IsOverseerVerboseLog, wdPath, exePath, app)
		if err != nil {
			return fmt.Errorf("manage %s: %w", nonExt, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
