package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/xackery/overseer/pkg/message"
)

var (
	Version = "0.0.0"
)

func main() {
	err := run()
	if err != nil {
		badf("Bootstrap failed: %s", err)
		os.Exit(1)
	}
}

func run() error {
	message.Banner("Bootstrap v" + Version)
	fmt.Println("This program bootstraps eqemu, doing quick sanity checks, downloading and upgrading binaries, then running overseer")
	type fileEntry struct {
		name string
		path string
	}
	files := []fileEntry{
		{name: "overseer", path: "overseer"},
		{name: "config.yaml", path: "config.yaml"},
		{name: "zone", path: "zone"},
		{name: "world", path: "world"},
		{name: "ucs", path: "ucs"},
		{name: "loginserver", path: "loginserver"},
		{name: "queryserv", path: "queryserv"},
	}
	for _, file := range files {
		fi, err := os.Stat(file.path)
		if err != nil {
			return fmt.Errorf("%s not found", file.name)
		}
		if fi.IsDir() {
			return fmt.Errorf("%s is a directory", file.name)
		}
	}

	goodf("Success")

	cmd := exec.Command("./overseer")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
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
