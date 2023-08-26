package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load loads a configuration file
func Load(path string) (*Configuration, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer r.Close()

	var config Configuration
	err = json.NewDecoder(r).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &config, nil
}

type Configuration struct {
	Server   ServerConfig   `json:"server"`
	WebAdmin WebAdminConfig `json:"web-admin"`
}

type ServerConfig struct {
	Zones       ZonesConfig       `json:"zones"`
	QSDatabase  QSDatabaseConfig  `json:"qsdatabase"`
	ChatServer  ChatServerConfig  `json:"chatserver"`
	MailServer  MailServerConfig  `json:"mailserver"`
	World       WorldConfig       `json:"world"`
	Database    DatabaseConfig    `json:"database"`
	Files       FilesConfig       `json:"files"`
	Directories DirectoriesConfig `json:"directories"`
}

type ZonesConfig struct {
	DefaultStatus string      `json:"defaultstatus"`
	Ports         PortsConfig `json:"ports"`
}

type PortsConfig struct {
	Low  string `json:"low"`
	High string `json:"high"`
}

type QSDatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

type ChatServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type MailServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type WorldConfig struct {
	Locked       bool              `json:"locked"`
	LoginServer1 LoginServerConfig `json:"loginserver1"`
	LoginServer2 LoginServerConfig `json:"loginserver2"`
	TCP          TCPConfig         `json:"tcp"`
	Telnet       TelnetConfig      `json:"telnet"`
	Key          string            `json:"key"`
	ShortName    string            `json:"shortname"`
	LongName     string            `json:"longname"`
	LocalAddress string            `json:"localaddress"`
	Address      string            `json:"address"`
	LoginServer3 LoginServerConfig `json:"loginserver3"`
}

type LoginServerConfig struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Legacy   string `json:"legacy"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type TCPConfig struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type TelnetConfig struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Enabled string `json:"enabled"`
}

type DatabaseConfig struct {
	DB       string `json:"db"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FilesConfig struct {
	Opcodes     string `json:"opcodes"`
	MailOpcodes string `json:"mail_opcodes"`
}

type DirectoriesConfig struct {
	Patches string `json:"patches"`
	Opcodes string `json:"opcodes"`
}

type WebAdminConfig struct {
	Application    ApplicationConfig `json:"application"`
	Launcher       LauncherConfig    `json:"launcher"`
	Quests         QuestsConfig      `json:"quests"`
	ServerCodePath string            `json:"serverCodePath"`
}

type ApplicationConfig struct {
	Key   string      `json:"key"`
	Admin AdminConfig `json:"admin"`
}

type AdminConfig struct {
	Password string `json:"password"`
}

type LauncherConfig struct {
	RunLoginServer   bool   `json:"runLoginserver"`
	RunQueryServ     bool   `json:"runQueryServ"`
	IsRunning        bool   `json:"isRunning"`
	StaticZones      string `json:"staticZones"`
	MinZoneProcesses int    `json:"minZoneProcesses"`
	RunSharedMemory  bool   `json:"runSharedMemory"`
}

type QuestsConfig struct {
	HotReload bool `json:"hotReload"`
}
