package main

import (
	"fmt"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
)

var (
	Version = "0.0.0"
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjZmY4MDAwIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6IjU2IiwiZm9udFNpemUiOiIyMiIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOiI1MCIsInNoYXBlIjoiZGlhbW9uZCIsInRleHQiOiJWRVJJRlkifQ
func main() {
	err := run()
	if err != nil {
		message.Badf("Verification failed: %s\n", err)
		operation.Exit(1)
	}
}

func run() error {
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}

	message.Banner("Verify v" + Version)

	fmt.Println("This program verifies eqemu as it runs, looking for things that may be wrong")
	err = eqemuConfig(config)
	if err != nil {
		return fmt.Errorf("eqemu_config.json %w", err)
	}
	message.OK("Success")

	return nil
}
