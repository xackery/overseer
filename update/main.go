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
		message.Badf("Update failed: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	message.Banner("Update v" + Version)
	fmt.Println("This program updates eqemu and all dependencies where applicable")
	err := eqemuConfig()
	if err != nil {
		return fmt.Errorf("eqemu_config.json %w", err)
	}
	message.OK("Success")

	return nil
}
