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
	Apps             []string
	IsScreenStart    bool
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
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			switch key {
			case "bin_path":
				config.BinPath = value
				if config.BinPath == "" {
					config.BinPath = "bin"
				}
			case "server_path":
				config.ServerPath = value
				if config.ServerPath == "" {
					config.ServerPath = "server"
				}
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
			case "app":
				config.Apps = append(config.Apps, value)
			case "is_screen_start":
				val, err := strconv.Atoi(value)
				if err != nil {
					config.IsScreenStart = false
					if strings.EqualFold(value, "true") {
						config.IsScreenStart = true
					}
				}
				if val == 1 {
					config.IsScreenStart = true
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
	fi, err := os.Stat("overseer.ini")
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat overseer.ini: %w", err)
		}
		w, err := os.Create("overseer.ini")
		if err != nil {
			return fmt.Errorf("create overseer.ini: %w", err)
		}
		w.Close()
	}
	if fi != nil && fi.IsDir() {
		return fmt.Errorf("overseer.ini is a directory")
	}

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
			if tmpConfig.BinPath == "1" {
				continue
			}
			out += fmt.Sprintf("%s = %s\n", key, c.BinPath)
			tmpConfig.BinPath = "1"
			continue
		case "server_path":
			if tmpConfig.ServerPath == "1" {
				continue
			}
			out += fmt.Sprintf("%s = %s\n", key, c.ServerPath)
			tmpConfig.ServerPath = "1"
			continue
		case "zone_count":
			if tmpConfig.ZoneCount == 1 {
				continue
			}
			out += fmt.Sprintf("%s = %d\n", key, c.ZoneCount)
			tmpConfig.ZoneCount = 1
			continue
		case "setup":
			if tmpConfig.Setup == "1" {
				continue
			}
			out += fmt.Sprintf("%s = %s\n", key, c.Setup)
			tmpConfig.Setup = "1"
			continue
		case "docker_network":
			if tmpConfig.DockerNetwork == "1" {
				continue
			}
			out += fmt.Sprintf("%s = %s\n", key, c.DockerNetwork)
			tmpConfig.DockerNetwork = "1"
			continue
		case "expansion":
			if tmpConfig.Expansion == "1" {
				continue
			}
			out += fmt.Sprintf("%s = %s\n", key, c.Expansion)
			tmpConfig.Expansion = "1"
			continue
		case "portable_database":
			if tmpConfig.PortableDatabase == 1 {
				continue
			}

			out += fmt.Sprintf("%s = %d\n", key, c.PortableDatabase)
			tmpConfig.PortableDatabase = 1
			continue
		case "auto_update":
			if tmpConfig.AutoUpdate == 1 {
				continue
			}

			out += fmt.Sprintf("%s = %d\n", key, c.AutoUpdate)
			tmpConfig.AutoUpdate = 1
			continue
		case "is_screen_start":
			if tmpConfig.IsScreenStart {
				continue
			}

			val := 0
			if c.IsScreenStart {
				val = 1
			}

			out += fmt.Sprintf("%s = %d\n", key, val)
			tmpConfig.IsScreenStart = true
			continue
		}
		line = fmt.Sprintf("%s = %s", key, value)
		out += line + "\n"
	}

	if tmpConfig.BinPath != "1" {
		out += fmt.Sprintf("bin_path = %s\n", c.BinPath)
	}
	if tmpConfig.ServerPath != "1" {
		out += fmt.Sprintf("server_path = %s\n", c.ServerPath)
	}
	if tmpConfig.ZoneCount != 1 {
		out += fmt.Sprintf("zone_count = %d\n", c.ZoneCount)
	}
	if tmpConfig.Setup != "1" {
		out += fmt.Sprintf("setup = %s\n", c.Setup)
	}
	if tmpConfig.DockerNetwork != "1" {
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
	val := 0
	if tmpConfig.IsScreenStart {
		val = 1
	}
	out += fmt.Sprintf("is_screen_start = %d\n", val)

	err = os.WriteFile("overseer.ini", []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
