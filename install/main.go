package main

import (
	"fmt"
	"math/rand"
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
	winExt  = ""
)

// icon link: https://prefinem.com/simple-icon-generator/#eyJiYWNrZ3JvdW5kQ29sb3IiOiIjMDA4MGZmIiwiYm9yZGVyQ29sb3IiOiIjMDAwMDAwIiwiYm9yZGVyV2lkdGgiOiI0IiwiZXhwb3J0U2l6ZSI6IjI1NiIsImV4cG9ydGluZyI6dHJ1ZSwiZm9udEZhbWlseSI6IkFiaGF5YSBMaWJyZSIsImZvbnRQb3NpdGlvbiI6IjU2IiwiZm9udFNpemUiOiIyMiIsImZvbnRXZWlnaHQiOjYwMCwiaW1hZ2UiOiIiLCJpbWFnZU1hc2siOiIiLCJpbWFnZVNpemUiOiI1MCIsInNoYXBlIjoic3F1YXJlIiwidGV4dCI6IklOU1RBTEwifQ
func main() {
	err := run()
	if err != nil {
		message.Badf("Install failed: %s\n", err)
		operation.Exit(1)
	}
	operation.Exit(0)
}

func run() error {
	if runtime.GOOS == "windows" {
		winExt = ".exe"
	}

	cfg, err := config.LoadOverseerConfig("overseer.ini")
	if err != nil {
		return fmt.Errorf("load overseer config: %w", err)
	}
	message.Banner("Install v" + Version)
	fmt.Println("This program installs eqemu, creating a usable environment from scratch")

	if cfg.Expansion != "" {
		choice, err := confirmation.New("It looks like install has been ran before. Would you like to reconfigure the install?", confirmation.No).RunPrompt()
		if err != nil {
			return fmt.Errorf("select reconfigure: %w", err)
		}
		if choice {
			fmt.Println("OK, flushing config! You'll be prompted new install options.")
			cfg.Expansion = ""
			cfg.PortableDatabase = 0
		}
	} else {
		err = config.ConfigSetup(cfg)
		if err != nil {
			return fmt.Errorf("install config setup: %w", err)
		}
	}

	start := time.Now()

	type opEntry struct {
		name string
		op   func(*config.OverseerConfiguration) error
	}
	ops := []opEntry{
		{"download binaries", downloadBinaries},
		{"configure eqemu_config.json", configureEqemuConfig},
		{"download quests", downloadQuests},
		{"download maps", downloadMaps},
		{"download portable database", downloadPortableDatabase},
		{"download assets", downloadAssets},
	}

	for _, op := range ops {
		err = op.op(cfg)
		if err != nil {
			return fmt.Errorf("%s: %w", op.name, err)
		}
	}

	seconds := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	message.OK("Successfully installed in " + seconds + " seconds")
	return nil
}

func downloadBinaries(cfg *config.OverseerConfiguration) error {
	if cfg.Setup == "docker" {
		return fmt.Errorf("docker setup not yet supported")
	}

	files := []string{
		"zone",
		"world",
	}

	areBinariesDownloaded := true
	for _, file := range files {
		path := "bin/" + file
		_, err := os.Stat(path + winExt)
		if err == nil {
			continue
		}
		areBinariesDownloaded = false
		break
	}

	if areBinariesDownloaded {
		message.Skip("Skipping binaries download, already exists")
		return nil
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

func downloadQuests(cfg *config.OverseerConfiguration) error {
	_, err := os.Stat("server/quests/airplane")
	if err == nil {
		message.Skip("Skipping quests download, already exists")
		return nil
	}

	err = os.MkdirAll("server/quests", 0755)
	if err != nil {
		return fmt.Errorf("mkdir server/quests: %w", err)
	}

	url := "https://github.com/eqemu-pack/" + cfg.ExpansionURI() + "/releases/download/latest/quests.zip"

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

func downloadMaps(cfg *config.OverseerConfiguration) error {
	_, err := os.Stat("server/maps/base")
	if err == nil {
		message.Skip("Skipping maps download, already exists")
		return nil
	}

	err = os.MkdirAll("server/maps", 0755)
	if err != nil {
		return fmt.Errorf("mkdir server/maps: %w", err)
	}

	url := "https://github.com/eqemu-pack/" + cfg.ExpansionURI() + "/releases/download/latest/maps.zip"

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

func downloadPortableDatabase(cfg *config.OverseerConfiguration) error {
	if cfg.PortableDatabase == 0 {
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

func configureEqemuConfig(cfg *config.OverseerConfiguration) error {
	path := cfg.ServerPath + "/eqemu_config.json"
	_, err := os.Stat(path)
	if err == nil {
		message.Skip("Skipping eqemu_config.json configuration, already exists")
		return nil
	}

	ecfg := config.EQEmuConfiguration{
		Server: config.ServerConfig{
			Zones: config.ZonesConfig{
				DefaultStatus: "0",
				Ports: config.PortsConfig{
					Low:  "7000",
					High: "7400",
				},
			},
			QSDatabase: config.QSDatabaseConfig{
				Host:     "localhost",
				Port:     "3306",
				Username: "root",
				Password: "eqemu",
				DB:       "peq",
			},
			ChatServer: config.ChatServerConfig{
				Port: "7778",
				Host: "",
			},
			MailServer: config.MailServerConfig{
				Host: "",
				Port: "7778",
			},
			World: config.WorldConfig{
				LoginServer1: config.LoginServerConfig{
					Account:  "",
					Password: "",
					Legacy:   "1",
					Host:     "login.eqemulator.net",
					Port:     "5998",
				},
				LoginServer2: config.LoginServerConfig{
					Port: "5998",
					Host: "login.projecteq.net",
				},
				TCP: config.TCPConfig{
					IP:   "127.0.0.1",
					Port: "9001",
				},
				Telnet: config.TelnetConfig{
					IP:      "0.0.0.0",
					Port:    "9000",
					Enabled: "true",
				},
				// generate a 32 character random string
				Key:       randomString(32),
				ShortName: "unk",
				// generate a 8 character random string
				LongName: fmt.Sprintf("Overseer [%s]", randomString(8)),
			},
			Database: config.DatabaseConfig{
				DB:   "peq",
				Host: "127.0.0.1",
				Port: "3306",
			},
			Files: config.FilesConfig{
				Opcodes:     "assets/opcodes/opcodes.conf",
				MailOpcodes: "assets/opcodes/mail_opcodes.conf",
			},
			Directories: config.DirectoriesConfig{
				Patches: "assets/patches/",
				Opcodes: "assets/opcodes/",
			},
		},
	}
	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	err = ecfg.Save(w)
	if err != nil {
		return fmt.Errorf("save %s: %w", path, err)
	}
	message.OK("Created eqemu_config.json")
	return nil
}

func downloadAssets(cfg *config.OverseerConfiguration) error {
	_, err := os.Stat("server/assets/opcodes")
	if err == nil {
		message.Skip("Skipping assets download, already exists")
		return nil
	}

	err = os.MkdirAll("server/assets", 0755)
	if err != nil {
		return fmt.Errorf("mkdir server/assets: %w", err)
	}

	url := "https://github.com/eqemu-pack/assets/releases/download/latest/assets.zip"
	cachePath := "server/cache/assets.zip"

	isCache, err := download.Save(url, cachePath)
	if err != nil {
		return fmt.Errorf("download %s: %w", filepath.Base(cachePath), err)
	}

	if isCache {
		fmt.Println("Using cached download at", cachePath)
	}

	err = zip.Unpack(cachePath, "server/assets/")
	if err != nil {
		return fmt.Errorf("unpack %s: %w", filepath.Base(cachePath), err)
	}

	message.OK("Downloaded assets")

	return nil
}

func randomString(length int) string {
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
