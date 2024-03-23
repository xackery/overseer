package main

import (
	"context"
	"fmt"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/xackery/overseer/diagnose/check"
	"github.com/xackery/overseer/diagnose/deep"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/gui"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
)

var (
	Version = "0.0.0"
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjZmZmZjAwIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6IjY3IiwiZm9udFNpemUiOiIzMiIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOiI1MCIsInNoYXBlIjoidHJpYW5nbGUiLCJ0ZXh0IjoiPyJ9
func main() {
	err := run()
	if err != nil {
		message.Badf("Diagnostics failed: %s\n", err)
		operation.Exit(1)
	}
	operation.Exit(0)
}

func run() error {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, err := NewMainWindow(ctx, cancel, Version)
	if err != nil {
		return fmt.Errorf("new main window: %w", err)
	}
	gui.New(g)

	cfg, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}

	message.Banner("Diagnose v" + Version)

	err = check.OverseerConfig()
	if err != nil {
		return fmt.Errorf("overseer.ini %w", err)
	}

	//fmt.Println("This program diagnoses eqemu's configuration, looking for things that may be wrong")
	err = check.EqemuConfig(cfg)
	if err != nil {
		return fmt.Errorf("parse %s: %w", cfg.ServerPath+"/eqemu_config.json", err)
	}

	emuCfg, err := config.LoadEQEmuConfig(cfg.ServerPath + "/eqemu_config.json")
	if err != nil {
		return fmt.Errorf("load eqemu config: %w", err)
	}

	err = check.Paths(cfg, emuCfg)
	if err != nil {
		return fmt.Errorf("paths %w", err)
	}

	message.OK("Completed quick diagnose")
	choice, err := confirmation.New("Run deep diagnostics?", confirmation.Yes).RunPrompt()
	if err != nil {
		return fmt.Errorf("select deep diagnose: %w", err)
	}
	if !choice {
		message.OK("OK, exiting\n")
		return nil
	}

	type op struct {
		name string
		op   func(*config.OverseerConfiguration, *config.EQEmuConfiguration) error
	}

	ops := []op{
		{"Table diagnose", deep.Table},
	}

	for _, op := range ops {
		err = op.op(cfg, emuCfg)
		if err != nil {
			message.Badf("%s: %s\n", op.name, err)
		} else {
			message.OKf("%s: OK\n", op.name)
		}
	}

	return nil
}
