package config

import (
	"fmt"
	"os"

	"github.com/xackery/overseer/pkg/message"
)

func overseerSetup() (*OverseerConfiguration, error) {
	message.Banner("Initial Setup")
	fmt.Println("Since no overseer.ini file was found, let's do some quick setup")

	os.WriteFile("overseer.ini", []byte(`# Overseer configuration file`), 0644)

	config := &OverseerConfiguration{}

	return config, nil
}
