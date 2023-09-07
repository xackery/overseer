package main

import (
	"fmt"
	"strings"

	"github.com/erikgeiser/promptkit/confirmation"
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

	preChoice := confirmation.Yes
	if service.IsDatabaseUp() {
		fmt.Println("It looks like a MySQL server is already running")
		preChoice = confirmation.No
	}

	isChoice, err := confirmation.New("Would you like to install and use a portable database service?", preChoice).RunPrompt()
	if err != nil {
		return fmt.Errorf("select portable database: %w", err)
	}
	if isChoice {
		config.PortableDatabase = 1
	} else {
		config.PortableDatabase = 0
	}

	isChoice, err = confirmation.New("Do you want to auto update everything on start run?", confirmation.No).RunPrompt()
	if err != nil {
		return fmt.Errorf("select auto update: %w", err)
	}
	if isChoice {
		config.AutoUpdate = 1
	} else {
		config.AutoUpdate = 0
	}

	config.BinPath = "bin"
	config.ServerPath = "server"

	err = config.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
