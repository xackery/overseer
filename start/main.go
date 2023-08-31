package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/sanity"
)

var (
	Version = "0.0.0"
)

func main() {
	err := run()
	if err != nil {
		badf("Start failed: %s", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Start v" + Version)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	fmt.Printf("Working directory: %s\n", cwd)

	err = sanity.EssentialFiles(config)
	if err != nil {
		return fmt.Errorf("essential files: %w", err)
	}

	cmd := exec.Command("./overseer")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("overseer run: %w", err)
	}
	return nil
}

func goodf(format string, a ...interface{}) {
	fmt.Printf(lipgloss.NewStyle().SetString("✅").
		Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
		PaddingRight(1).
		String()+format+"\n", a...)
}

func badf(format string, a ...interface{}) {
	fmt.Printf(lipgloss.NewStyle().SetString("❌").
		Foreground(lipgloss.AdaptiveColor{Light: "#FF5555", Dark: "#FF5555"}).
		PaddingRight(1).
		String()+format+"\n", a...)
}
