package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/flog"
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
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}

	err = parseManager(config)
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	err = flog.New("overseer.log")
	if err != nil {
		return fmt.Errorf("new flog: %w", err)
	}
	defer flog.Close()

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

	flog.Printf("Setup type: %s\n", cfg.Setup)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	dir, err := filepath.Abs(cwd + "/" + cfg.ServerPath)
	if err != nil {
		return fmt.Errorf("abs: %w", err)
	}

	command, err := filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/zone"+winExt)
	if err != nil {
		return fmt.Errorf("rel: %w", err)
	}

	for i := 0; i < cfg.ZoneCount; i++ {
		manager.Manage(setupType, fmt.Sprintf("zone%d", i), dir, command)
	}

	command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/world"+winExt)
	if err != nil {
		return fmt.Errorf("rel: %w", err)
	}
	manager.Manage(setupType, "world", dir, command)
	command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/ucs"+winExt)
	if err != nil {
		return fmt.Errorf("rel: %w", err)
	}
	manager.Manage(setupType, "ucs", dir, command)
	//command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/queryserv"+winExt)
	//if err != nil {
	//	return fmt.Errorf("rel: %w", err)
	//}
	//manager.Manage(setupType, "queryserv", dir, command)
	//command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/loginserver"+winExt)
	//if err != nil {
	//	return fmt.Errorf("rel: %w", err)
	//}
	//manager.Manage(setupType, "loginserver", dir, command)

	for _, app := range cfg.Apps {
		nonExt := strings.TrimSuffix(app, filepath.Ext(app))
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/"+app)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
		manager.Manage(setupType, nonExt, dir, command)
	}
	return nil
}
