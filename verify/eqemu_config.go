package main

import (
	"fmt"
	"os"

	"github.com/xackery/overseer/share/config"
	"github.com/xackery/overseer/share/message"
)

func eqemuConfig(cfg *config.OverseerConfiguration) error {
	fi, err := os.Stat(cfg.ServerPath + "/eqemu_config.json")
	if err != nil {
		return fmt.Errorf("not found")
	}
	if fi.IsDir() {
		return fmt.Errorf("is a directory")
	}

	r, err := os.Open(cfg.ServerPath + "/eqemu_config.json")
	if err != nil {
		return err
	}
	defer r.Close()

	config, err := config.LoadEQEmuConfig(cfg.ServerPath + "/eqemu_config.json")
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	if config.Server.Database.DB == "" {
		return fmt.Errorf("database.db is empty")
	}

	message.OK(cfg.ServerPath + "/eqemu_config.json found")
	return nil
}
