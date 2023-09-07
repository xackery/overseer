package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
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
		os.Exit(1)
	}
}

func run() error {
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Start v" + Version)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	fmt.Printf("Working directory: %s\n", cwd)

	err = sanity.EssentialFiles(config)
	if err != nil {
		return fmt.Errorf("essential files: %w", err)
	}

	if !service.IsDatabaseUp() {
		if config.PortableDatabase != 1 {
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

	cmd := exec.Command("./overseer")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("overseer run: %w", err)
	}
	return nil
}
