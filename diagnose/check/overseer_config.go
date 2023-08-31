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
	fi, err := os.Stat("eqemu_config.json")
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
			case "server_path":
				cfg.ServerPath = value
			case "zone_count":
				cfg.ZoneCount, err = strconv.Atoi(value)
				if err != nil {
					message.Badf("parse zone_count value %s: %w", value, err)
				}
			case "setup":
				cfg.Setup = strings.ToLower(value)
			case "docker_network":
				cfg.DockerNetwork = value
			case "expansion":
				cfg.Expansion = value
				if !config.IsValidExpansion(value) {
					message.Badf("overseer.ini unknown expansion value %s", value)
					message.Link("https://o.eqcodex.com/105")
				}
			default:
				message.Badf("overseer.ini unknown key in overseer.ini: %s", key)
			}
		}
	}
	return nil
}
