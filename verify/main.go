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
		message.Badf("Verification failed: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	message.Banner("Verify v" + Version)
	fmt.Println("This program verifies eqemu as it runs, looking for things that may be wrong")
	err := eqemuConfig()
	if err != nil {
		return fmt.Errorf("eqemu_config.json %w", err)
	}
	message.OK("Success")

	return nil
}
