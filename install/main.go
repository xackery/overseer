package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/download"
	"github.com/xackery/overseer/pkg/message"
	"github.com/xackery/overseer/pkg/operation"
	"github.com/xackery/overseer/pkg/zip"
)

var (
	Version = "0.0.0"
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjMDA4MGZmIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6IjU2IiwiZm9udFNpemUiOiIyMiIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOiI1MCIsInNoYXBlIjoic3F1YXJlIiwidGV4dCI6IklOU1RBTEwifQ
func main() {
	err := run()
	if err != nil {
		message.Badf("Install failed: %s\n", err)

		if runtime.GOOS == "windows" {
			fmt.Println("Press any key to continue...")
			fmt.Scanln()
		}
		operation.Exit(1)
	}

	if runtime.GOOS == "windows" {
		fmt.Println("Press any key to continue...")
		fmt.Scanln()
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
		choice, err := confirmation.New("It looks like install has been ran before. Would you like to reconfigure the install?", confirmation.No).RunPrompt()
		if err != nil {
			return fmt.Errorf("select reconfigure: %w", err)
		}
		if choice {
			fmt.Println("OK, flushing config! You'll be prompted new install options.")
			config.Expansion = ""
			config.PortableDatabase = 0
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

	err = downloadPortableDatabase(config)
	if err != nil {
		return fmt.Errorf("download portable database: %w", err)
	}

	seconds := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	message.OK("Successfully installed in " + seconds + " seconds")
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

func downloadPortableDatabase(config *config.OverseerConfiguration) error {
	if config.PortableDatabase == 0 {
		return nil
	}
	err := os.MkdirAll("server/database", 0755)
	if err != nil {
		return fmt.Errorf("mkdir server/database: %w", err)
	}

	//url := "https://archive.mariadb.org/mariadb-10.6.10/winx64-packages/mariadb-10.6.10-winx64.zip"
	url := "https://cdn.mysql.com//Downloads/MySQL-8.1/mysql-8.1.0-linux-glibc2.17-x86_64-minimal.tar.xz"
	//cachePath := "server/cache/mariadb-10.6.10-winx64.zip"
	cachePath := "server/cache/mysql-8.1.0-linux-glibc2.17-x86_64-minimal.tar.xz"
	if runtime.GOOS != "windows" {
		//url = "https://archive.mariadb.org/mariadb-10.6.10/bintar-linux-systemd-x86_64/mariadb-10.6.10-linux-systemd-x86_64.tar.gz"
		url = "https://cdn.mysql.com//Downloads/MySQL-8.1/mysql-8.1.0-winx64.zip"
		//cachePath = "server/cache/mariadb-10.6.10-linux-systemd-x86_64.tar.gz"
		cachePath = "server/cache/mysql-8.1.0-winx64.zip"
	}

	isCache, err := download.Save(url, cachePath)
	if err != nil {
		return fmt.Errorf("download %s: %w", filepath.Base(cachePath), err)
	}

	if isCache {
		fmt.Println("Using cached download at", cachePath)
	}

	err = zip.Unpack(cachePath, "server/database/")
	if err != nil {
		return fmt.Errorf("unpack %s: %w", filepath.Base(cachePath), err)
	}

	message.OK("Downloaded db")

	return nil
}
