package main

import (
	"fmt"

	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
)

var (
	Version = "0.0.0"
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjODBmZmZmIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6Ijg2IiwiZm9udFNpemUiOiIxOCIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOiI1MCIsInNoYXBlIjoidHJpYW5nbGUiLCJ0ZXh0IjoiVVBEQVRFIn0
func main() {
	err := run()
	if err != nil {
		message.Badf("Update failed: %s\n", err)
		operation.Exit(1)
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
