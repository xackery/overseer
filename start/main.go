package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
	"github.com/xackery/overseer/pkg/sanity"
)

var (
	Version = "0.0.0"
)

func main() {
	err := run()
	if err != nil {
		message.Badf("Start failed: %s\n", err)
		operation.Exit(1)
	}
	operation.Exit(0)
}

func run() error {
	winExt := ""
	if runtime.GOOS == "windows" {
		winExt = ".exe"
	}

	cfg, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Start v" + Version)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	fmt.Printf("Working directory: %s\n", cwd)

	err = sanity.EssentialFiles(cfg)
	if err != nil {
		return fmt.Errorf("essential files: %w", err)
	}

	optionList := []string{
		"overseer (all)",
		"shared_memory",
		"world",
		"zone",
		"queryserv",
		"ucs",
		"loginserver",
	}

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		command, err = selection.New("Start which program?", optionList).RunPrompt()
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}
	}

	choice := command
	dir, err := filepath.Abs(cwd + "/" + cfg.ServerPath)
	if err != nil {
		return fmt.Errorf("abs: %w", err)
	}
	switch choice {
	case "overseer (all)":
		command = "./overseer" + winExt
		dir = cwd
	case "overseer":
		command = "./overseer" + winExt
		dir = cwd
	case "shared_memory":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/shared_memory"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "world":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/world"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "zone":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/zone"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "queryserv":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/queryserv"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "ucs":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/ucs"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	case "loginserver":
		command, err = filepath.Rel(dir, cwd+"/"+cfg.BinPath+"/loginserver"+winExt)
		if err != nil {
			return fmt.Errorf("rel: %w", err)
		}
	default:
		return fmt.Errorf("unknown argument %s", command)
	}

	if choice == "overseer (all)" {
		choice = "overseer"
	}
	/*if !service.IsDatabaseUp() {
		if cfg.PortableDatabase != 1 {
			message.Bad("Database is not running and we aren't portable. Please run it first.")
			return fmt.Errorf("database not running")
		}
		err = service.DatabaseStart()
		if err != nil {
			return fmt.Errorf("database start: %w", err)
		}

		if !service.IsDatabaseUp() {
			message.Bad("Database is not running even after trying to start it.")
			return fmt.Errorf("database failed to start")
		}
	}*/

	fmt.Println("running", command, "from", dir)
	start := time.Now()
	cmd := exec.Command(command)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	err = cmd.Run()
	if err != nil {
		message.Badf("%s exited after %0.2f seconds\n", choice, time.Since(start).Seconds())
		return fmt.Errorf("%s run: %w", command, err)
	}
	message.OKf("%s exited after %0.2f seconds\n", choice, time.Since(start).Seconds())

	return nil
}
