package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/xackery/overseer/diagnose/check"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
)

var (
	Version = "0.0.0"
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjZmZmZjAwIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6IjY3IiwiZm9udFNpemUiOiIzMiIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOiI1MCIsInNoYXBlIjoidHJpYW5nbGUiLCJ0ZXh0IjoiPyJ9
func main() {
	err := run()
	if err != nil {
		message.Badf("Diagnostics failed: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var err error
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}

	message.Banner("Diagnose v" + Version)

	err = check.OverseerConfig()
	if err != nil {
		return fmt.Errorf("overseer.ini %w", err)
	}

	//fmt.Println("This program diagnoses eqemu's configuration, looking for things that may be wrong")
	err = check.EqemuConfig(config)
	if err != nil {
		return fmt.Errorf("parse %s: %w", config.ServerPath+"/eqemu_config.json", err)
	}

	err = check.Paths(config)
	if err != nil {
		return fmt.Errorf("paths %w", err)
	}

	message.OK("Completed diagnose")

	if runtime.GOOS == "windows" {
		fmt.Println("Press any key to continue...")
		fmt.Scanln()
	}

	return nil
}
