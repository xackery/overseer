package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/service"
)

func overseerSetup() (*OverseerConfiguration, error) {
	message.Banner("Initial Setup")
	fmt.Println("Since no overseer.ini file was found, let's do some quick setup")

	config := &OverseerConfiguration{}
	err := ConfigSetup(config)
	if err != nil {
		return nil, fmt.Errorf("config setup: %w", err)
	}

	return config, nil
}

func ConfigSetup(cfg *OverseerConfiguration) error {

	if cfg.Expansion != "" {
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
	cfg.Expansion = strings.ToLower(choice)

	choice, err = selection.New("Setup? (default recommended, docker experimental)", []string{
		"default",
		"docker",
	}).RunPrompt()
	if err != nil {
		return fmt.Errorf("select setup: %w", err)
	}
	cfg.Setup = strings.ToLower(choice)

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
		cfg.PortableDatabase = 1
	} else {
		cfg.PortableDatabase = 0
	}

	isChoice, err = confirmation.New("Do you want to auto update everything on start run?", confirmation.No).RunPrompt()
	if err != nil {
		return fmt.Errorf("select auto update: %w", err)
	}
	if isChoice {
		cfg.AutoUpdate = 1
	} else {
		cfg.AutoUpdate = 0
	}

	choice, err = selection.New("How many zones should be started?", []string{
		"1",
		"2",
		"3",
		"5",
		"10",
		"15",
		"20",
		"50",
	}).RunPrompt()
	if err != nil {
		return fmt.Errorf("zone setup: %w", err)
	}
	count, err := strconv.Atoi(choice)
	if err != nil {
		return fmt.Errorf("parse zone count: %w", err)
	}

	cfg.ZoneCount = count

	cfg.BinPath = "bin"
	cfg.ServerPath = "server"

	err = cfg.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
