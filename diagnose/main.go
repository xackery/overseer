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
	err = check.EqemuConfig()
	if err != nil {
		return fmt.Errorf("eqemu_config.json %w", err)
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
