package main

import (
	"fmt"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/service"
)

func installConfigSetup(config *config.OverseerConfiguration) error {
	if config.Expansion != "" {
		return nil
	}

	choice, err := selection.New("Expansion?", []string{
		"Classic",
		"Kunark",
		"Velious",
		"Luclin",
		"PoP",
		"Ykesha",
		"Gates",
		"Omens",
		"Dragons",
	}).RunPrompt()
	if err != nil {
		return fmt.Errorf("select expansion: %w", err)
	}
	config.Expansion = strings.ToLower(choice)

	choice, err = selection.New("Setup? (default recommended, docker experimental)", []string{
		"default",
		"docker",
	}).RunPrompt()
	if err != nil {
		return fmt.Errorf("select setup: %w", err)
	}
	config.Setup = strings.ToLower(choice)

	dbOptions := []string{"Yes", "No"}
	if service.IsDatabaseUp() {
		fmt.Println("It looks like a MySQL server is already running")
		dbOptions = []string{"No", "Yes"}
	}
	choice, err = selection.New("Would you like to install and use a portable database service?", dbOptions).RunPrompt()
	if err != nil {
		return fmt.Errorf("select portable database: %w", err)
	}
	if strings.Contains(strings.ToLower(choice), "yes") {
		config.PortableDatabase = 1
	} else {
		config.PortableDatabase = 0
	}

	err = config.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
