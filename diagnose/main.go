package main

import (
	"fmt"
	"os"

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
	message.Banner("Diagnose v" + Version)
	fmt.Println("This program diagnoses eqemu's configuration, looking for things that may be wrong")
	err := eqemuConfig()
	if err != nil {
		return fmt.Errorf("eqemu_config.json %w", err)
	}
	message.OK("Success")

	return nil
}
