package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/download"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/zip"
)

var (
	Version = "0.0.0"
)

func main() {
	err := run()
	if err != nil {
		message.Badf("Install failed: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Install v" + Version)
	fmt.Println("This program installs eqemu, creating a usable environment from scratch")

	if config.Expansion != "" {
		choice, err := selection.New("It looks like install has been ran before. Would you like to reconfigure the install?", []string{"No", "Yes"}).RunPrompt()
		if err != nil {
			return fmt.Errorf("select reconfigure: %w", err)
		}
		if strings.Contains(strings.ToLower(choice), "yes") {
			fmt.Println("OK, flushing config! You'll be prmopted new install options.")
			config.Expansion = ""
		}
	}

	err = installConfigSetup(config)
	if err != nil {
		return fmt.Errorf("install config setup: %w", err)
	}

	start := time.Now()
	err = downloadBinaries(config)
	if err != nil {
		return fmt.Errorf("download binaries: %w", err)
	}

	err = downloadQuests(config)
	if err != nil {
		return fmt.Errorf("download quests: %w", err)
	}

	err = downloadMaps(config)
	if err != nil {
		return fmt.Errorf("download quests: %w", err)
	}

	seconds := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	message.OK("Successfully installed in " + seconds + " seconds")

	return nil
}

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

	err = config.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	message.OKf("Expansion: %s\n", config.Expansion)

	return nil
}

func downloadBinaries(config *config.OverseerConfiguration) error {
	if config.Setup == "docker" {
		return fmt.Errorf("docker setup not yet supported")
	}

	url := "https://github.com/EQEmu/Server/releases/latest/download/eqemu-server-" + runtime.GOOS + "-x64.zip"
	if runtime.GOOS == "darwin" {
		url = "https://github.com/EQEmu/Server/releases/latest/download/eqemu-server-linux-x64.zip"
	}

	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return fmt.Errorf("mkdir bin: %w", err)
	}

	err = os.MkdirAll("server/cache", 0755)
	if err != nil {
		return fmt.Errorf("mkdir cache: %w", err)
	}

	cachePath := "server/cache/eqemu-server-" + runtime.GOOS + "-x64.zip"
	if runtime.GOOS == "darwin" {
		cachePath = "server/cache/eqemu-server-linux-x64.zip"
	}
	isCache, err := download.Save(url, cachePath)
	if err != nil {
		return fmt.Errorf("download %s: %w", filepath.Base(cachePath), err)
	}

	if isCache {
		fmt.Println("Using cached download at", cachePath)
	}

	err = zip.Unpack(cachePath, "bin")
	if err != nil {
		return fmt.Errorf("unpack %s: %w", filepath.Base(cachePath), err)
	}

	message.OK("Downloaded binaries")

	return nil
}

func downloadQuests(config *config.OverseerConfiguration) error {

	err := os.MkdirAll("server/quests", 0755)
	if err != nil {
		return fmt.Errorf("mkdir server/quests: %w", err)
	}

	url := "https://github.com/eqemu-pack/" + config.ExpansionURI() + "/releases/download/latest/quests.zip"

	cachePath := "server/cache/quests.zip"
	isCache, err := download.Save(url, cachePath)
	if err != nil {
		return fmt.Errorf("download %s: %w", filepath.Base(cachePath), err)
	}

	if isCache {
		fmt.Println("Using cached download at", cachePath)
	}

	err = zip.Unpack(cachePath, "server")
	if err != nil {
		return fmt.Errorf("unpack %s: %w", filepath.Base(cachePath), err)
	}

	message.OK("Downloaded quests")

	return nil
}

func downloadMaps(config *config.OverseerConfiguration) error {

	err := os.MkdirAll("server/maps", 0755)
	if err != nil {
		return fmt.Errorf("mkdir server/maps: %w", err)
	}

	url := "https://github.com/eqemu-pack/" + config.ExpansionURI() + "/releases/download/latest/maps.zip"

	cachePath := "server/cache/maps.zip"
	isCache, err := download.Save(url, cachePath)
	if err != nil {
		return fmt.Errorf("download %s: %w", filepath.Base(cachePath), err)
	}

	if isCache {
		fmt.Println("Using cached download at", cachePath)
	}

	err = zip.Unpack(cachePath, "server/maps/")
	if err != nil {
		return fmt.Errorf("unpack %s: %w", filepath.Base(cachePath), err)
	}

	message.OK("Downloaded maps")

	return nil
}
