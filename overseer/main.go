package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/xackery/overseer/pkg/signal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xackery/overseer/pkg/dashboard"
	"github.com/xackery/overseer/pkg/manager"
	"github.com/xackery/overseer/pkg/reporter"
)

var (
	Version = "0.0.0"
)

func main() {
	start := time.Now()
	err := run()
	// clear the screen on exit
	fmt.Print("\033[H\033[2J")
	if err != nil {
		fmt.Println("Overseer failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Overseer exited after %0.2f seconds\n", time.Since(start).Seconds())
}

func run() error {

	err := parseManager()
	if err != nil {
		return err
	}

	p := tea.NewProgram(dashboard.New(Version))
	go func() {
		for {
			select {
			case <-reporter.SendUpdateChan:
				p.Send(dashboard.RefreshRequest{})
			case <-signal.Ctx().Done():
				return
			case <-time.After(5 * time.Second):
				p.Send(dashboard.RefreshRequest{})
			}
		}

	}()

	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}

func parseManager() error {

	r, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer r.Close()

	type yamlConfig struct {
		World struct {
			IsLocked   bool   `yaml:"is_locked"`
			ShortName  string `yaml:"short_name"`
			LongName   string `yaml:"long_name"`
			WanAddress string `yaml:"wan_address"`
			LanAddress string `yaml:"lan_address"`

			MaxClients      int    `yaml:"max_clients"`
			IntranetIP      string `yaml:"intranet_ip"`
			IntranetPort    int    `yaml:"intranet_port"`
			TelnetIsEnabled bool   `yaml:"telnet_is_enabled"`
			TelnetIP        string `yaml:"telnet_ip"`
			TelnetPort      int    `yaml:"telnet_port"`
			SharedKey       string `yaml:"shared_key"`
		} `yaml:"world"`
		Ucs struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"ucs"`
		Database struct {
			DB       string `yaml:"db"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"database"`
		QueryServ struct {
			DB       string `yaml:"db"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"query_serv"`
		Zone struct {
			DefaultStatus int `yaml:"default_status"`
			PortMin       int `yaml:"port_min"`
			PortMax       int `yaml:"port_max"`
			MaxZones      int `yaml:"max_zones"`
		} `yaml:"zone"`
		LoginServer []struct {
			Port     int    `yaml:"port"`
			Account  string `yaml:"account"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Type     int    `yaml:"type"`
		} `yaml:"login_server"`
		Dir struct {
			Patches      string `yaml:"patches"`
			Opcodes      string `yaml:"opcodes"`
			SharedMemory string `yaml:"shared_memory"`
			LuaModules   string `yaml:"lua_modules"`
			Quests       string `yaml:"quests"`
			Plugins      string `yaml:"plugins"`
			Maps         string `yaml:"maps"`
			Logs         string `yaml:"logs"`
		} `yaml:"dir"`
		Other map[string]string `yaml:"other"`
	}
	config := yamlConfig{}
	dec := yaml.NewDecoder(r)
	err = dec.Decode(&config)
	if err != nil {
		return fmt.Errorf("decode config.yaml: %w", err)
	}

	zoneCount := 0
	for i := config.Zone.PortMin; i < config.Zone.PortMax; i++ {
		zoneCount++
		if zoneCount > config.Zone.MaxZones {
			break
		}
		manager.Manage(fmt.Sprintf("zone%d", i), "./zone", fmt.Sprintf("%d", i))
	}

	manager.Manage("world", "./world")
	manager.Manage("ucs", "./ucs")
	manager.Manage("queryserv", "./queryserv")
	manager.Manage("loginserver", "./loginserver")
	for k, v := range config.Other {
		manager.Manage(k, v)
	}
	return nil
}
