package check

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
)

func OverseerConfig() error {
	message.OKReset()
	fi, err := os.Stat("overseer.ini")
	if err != nil {
		return fmt.Errorf("not found")
	}
	if fi.IsDir() {
		return fmt.Errorf("is a directory")
	}
	r, err := os.Open("overseer.ini")
	if err != nil {
		return fmt.Errorf("open: %s", strings.TrimPrefix(err.Error(), "open overseer.ini: "))
	}
	defer r.Close()

	var cfg config.OverseerConfiguration
	tmpConfig := config.OverseerConfiguration{}

	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "bin_path":
				cfg.BinPath = value
				tmpConfig.BinPath = "1"
			case "server_path":
				cfg.ServerPath = value
				tmpConfig.ServerPath = "1"
			case "zone_count":
				cfg.ZoneCount, err = strconv.Atoi(value)
				if err != nil {
					message.Badf("parse zone_count value %s: %s", value, err)
				}
				tmpConfig.ZoneCount = 1
			case "setup":
				cfg.Setup = strings.ToLower(value)
				tmpConfig.Setup = "1"
			case "docker_network":
				cfg.DockerNetwork = value
				tmpConfig.DockerNetwork = "1"
			case "expansion":
				cfg.Expansion = value
				if !config.IsValidExpansion(value) {
					message.Badf("overseer.ini unknown expansion value %s", value)
					message.Link("https://o.eqcodex.com/105")
				}
				tmpConfig.Expansion = "1"
			case "auto_update":
				cfg.AutoUpdate, err = strconv.Atoi(value)
				if err != nil {
					message.Badf("parse auto_update value %s: %s", value, err)
				}
				if cfg.AutoUpdate != 0 && cfg.AutoUpdate != 1 {
					message.Badf("overseer.ini unknown auto_update value %s", value)
				}
				tmpConfig.AutoUpdate = 1
			case "portable_database":
				cfg.PortableDatabase, err = strconv.Atoi(value)
				if err != nil {
					message.Badf("parse portable_database value %s: %s", value, err)
				}
				if cfg.PortableDatabase != 0 && cfg.PortableDatabase != 1 {
					message.Badf("overseer.ini unknown portable_database value %s", value)
				}
				tmpConfig.PortableDatabase = 1
			case "is_screen_start":

			default:
				message.Badf("overseer.ini unknown key in overseer.ini: %s", key)
			}
		}
	}

	if tmpConfig.BinPath == "" {
		message.Bad("overseer.ini missing bin_path")
	}

	if tmpConfig.ServerPath == "" {
		message.Bad("overseer.ini missing server_path")
	}

	if tmpConfig.ZoneCount == 0 {
		message.Bad("overseer.ini missing zone_count")
	}

	if tmpConfig.Setup == "" {
		message.Bad("overseer.ini missing setup")
	}

	if tmpConfig.DockerNetwork == "" {
		message.Bad("overseer.ini missing docker_network")
	}

	if tmpConfig.Expansion == "" {
		message.Bad("overseer.ini missing expansion")
	}

	if tmpConfig.AutoUpdate == 0 {
		message.Bad("overseer.ini missing auto_update")
	}

	if tmpConfig.PortableDatabase == 0 {
		message.Bad("overseer.ini missing portable_database")
	}

	for _, app := range tmpConfig.Apps {
		fi, err := os.Stat(fmt.Sprintf("%s/%s", cfg.BinPath, app))
		if err != nil {
			message.Badf("overseer.ini app %s not found", app)
			continue
		}
		if fi.IsDir() {
			message.Badf("overseer.ini app %s is a directory", app)
			continue
		}
	}

	if message.IsOK() {
		message.OK("Overseer Config OK")
	}

	return nil
}
