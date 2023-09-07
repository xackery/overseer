package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
	"github.com/xackery/overseer/pkg/sanity"
	"github.com/xackery/overseer/pkg/service"
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

	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, err = selection.New("Start which program?", []string{
			"overseer (all)",
			"shared_memory",
			"world",
			"zone",
			"queryserv",
			"ucs",
			"loginserver",
		}).RunPrompt()
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}
	}

	choice := path
	dir := cfg.ServerPath
	switch path {
	case "overseer (all)":
		path = "overseer" + winExt
		dir = "."
	case "shared_memory":
		path = cwd + "/" + cfg.BinPath + "/shared_memory" + winExt
	case "world":
		path = cwd + "/" + cfg.BinPath + "/world" + winExt
	case "zone":
		path = cwd + "/" + cfg.BinPath + "/zone" + winExt
	case "queryserv":
		path = cwd + "/" + cfg.BinPath + "/queryserv" + winExt
	case "ucs":
		path = cwd + "/" + cfg.BinPath + "/ucs" + winExt
	case "loginserver":
		path = cwd + "/" + cfg.BinPath + "/loginserver" + winExt
	default:
		return fmt.Errorf("unknown argument %s", path)
	}

	if !service.IsDatabaseUp() {
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
	}

	start := time.Now()
	cmd := exec.Command(path)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	message.OKf("%s exited after %0.2f seconds\n", choice, time.Since(start).Seconds())
	if err != nil {
		return fmt.Errorf("%s run: %w", path, err)
	}

	return nil
}
