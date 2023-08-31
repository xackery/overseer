package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type OverseerConfiguration struct {
	BinPath    string
	ServerPath string
	ZoneCount  int
	// Setup represents the setup type, options include default (bare-metal), docker (docker run), docker-compose (akk-stack)
	Setup string
	// If setup is docker, this is the network to use, defaults to eqemu
	DockerNetwork    string
	Expansion        string
	PortableDatabase int
	AutoUpdate       int
}

// LoadOverseerConfig loads an overseer config file
func LoadOverseerConfig(path string) (*OverseerConfiguration, error) {
	_, err := os.Stat(path)
	if err != nil {
		return overseerSetup()
	}

	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %s", strings.TrimPrefix(err.Error(), "open overseer.ini: "))
	}
	defer r.Close()

	var config OverseerConfiguration

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
				config.BinPath = value
			case "server_path":
				config.ServerPath = value
			case "zone_count":
				config.ZoneCount, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("parse zone_count: %w", err)
				}
			case "auto_update":
				config.AutoUpdate, err = strconv.Atoi(value)
				if err != nil {
					config.AutoUpdate = 0
				}
			case "setup":
				config.Setup = strings.ToLower(value)
			case "docker_network":
				if value == "" {
					value = "eqemu"
				}
				badValue := strings.ToLower(value)
				switch badValue {
				case "bridge":
					value = "eqemu"
				case "host":
					value = "eqemu"
				case "none":
					value = "eqemu"
				}

				config.DockerNetwork = value
			case "expansion":
				value := strings.ToLower(value)
				if !IsValidExpansion(value) {
					value = ""
				}
				config.Expansion = value
			case "portable_database":
				config.PortableDatabase, err = strconv.Atoi(value)
				if err != nil {
					config.PortableDatabase = 0
				}
			default:
				return nil, fmt.Errorf("unknown key in overseer.ini: %s", key)
			}
		}
	}

	return &config, nil
}

func IsValidExpansion(name string) bool {
	switch strings.ToLower(name) {
	case "classic":
		return true
	case "kunark":
		return true
	case "velious":
		return true
	case "luclin":
		return true
	case "pop":
		return true
	case "ldon":
		return true
	case "ykesha":
		return true
	case "gates":
		return true
	case "omens":
		return true
	case "dragons":
		return true
	default:
		return false
	}
}

func (o *OverseerConfiguration) ExpansionURI() string {
	switch o.Expansion {
	case "classic":
		return "classic"
	case "kunark":
		return "kunark"
	case "velious":
		return "velious"
	case "luclin":
		return "luclin"
	case "pop":
		return "pop"
	case "ldon":
		return "ldon"
	case "ykesha":
		return "ykesha"
	case "gates":
		return "gates"
	case "omens":
		return "omens"
	case "dragons":
		return "dragons"
	default:
		return ""
	}
}

// Save saves the config
func (c *OverseerConfiguration) Save() error {
	r, err := os.Open("overseer.ini")
	if err != nil {
		return fmt.Errorf("open: %s", strings.TrimPrefix(err.Error(), "open overseer.ini: "))
	}
	defer r.Close()

	tmpConfig := OverseerConfiguration{}

	out := ""
	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "#") {
			out += line + "\n"
			continue
		}
		if !strings.Contains(line, "=") {
			out += line + "\n"
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "bin_path":
			out += fmt.Sprintf("%s = %s\n", key, c.BinPath)
			tmpConfig.BinPath = "1"
			continue
		case "server_path":
			out += fmt.Sprintf("%s = %s\n", key, c.ServerPath)
			tmpConfig.ServerPath = "1"
		case "zone_count":
			out += fmt.Sprintf("%s = %d\n", key, c.ZoneCount)
			tmpConfig.ZoneCount = 1
		case "setup":
			out += fmt.Sprintf("%s = %s\n", key, c.Setup)
			tmpConfig.Setup = "1"
		case "docker_network":
			out += fmt.Sprintf("%s = %s\n", key, c.DockerNetwork)
			tmpConfig.DockerNetwork = "1"
		case "expansion":
			out += fmt.Sprintf("%s = %s\n", key, c.Expansion)
			tmpConfig.Expansion = "1"
		case "portable_database":
			out += fmt.Sprintf("%s = %d\n", key, c.PortableDatabase)
			tmpConfig.PortableDatabase = 1
		case "auto_update":
			out += fmt.Sprintf("%s = %d\n", key, c.AutoUpdate)
			tmpConfig.AutoUpdate = 1
		}
		line = fmt.Sprintf("%s = %s", key, value)
		out += line + "\n"
	}

	if tmpConfig.BinPath == "" {
		out += fmt.Sprintf("bin_path = %s\n", c.BinPath)
	}
	if tmpConfig.ServerPath == "" {
		out += fmt.Sprintf("server_path = %s\n", c.ServerPath)
	}
	if tmpConfig.ZoneCount == 0 {
		out += fmt.Sprintf("zone_count = %d\n", c.ZoneCount)
	}
	if tmpConfig.Setup == "" {
		out += fmt.Sprintf("setup = %s\n", c.Setup)
	}
	if tmpConfig.DockerNetwork == "" {
		out += fmt.Sprintf("docker_network = %s\n", c.DockerNetwork)
	}

	if tmpConfig.Expansion == "" {
		out += fmt.Sprintf("expansion = %s\n", c.Expansion)
	}
	if tmpConfig.PortableDatabase == 0 {
		out += fmt.Sprintf("portable_database = %d\n", c.PortableDatabase)
	}
	if tmpConfig.AutoUpdate == 0 {
		out += fmt.Sprintf("auto_update = %d\n", c.AutoUpdate)
	}

	err = os.WriteFile("overseer.ini", []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
