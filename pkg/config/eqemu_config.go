package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// EQEmuConfiguration is the configuration for the EQEmu server
type EQEmuConfiguration struct {
	Server   ServerConfig   `json:"server"`
	WebAdmin WebAdminConfig `json:"web-admin"`
}

// ServerConfig is the configuration for the EQEmu server
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

// ZonesConfig is the configuration for the EQEmu server zones
type ZonesConfig struct {
	DefaultStatus string      `json:"defaultstatus"`
	Ports         PortsConfig `json:"ports"`
}

// PortsConfig is the configuration for the EQEmu server ports
type PortsConfig struct {
	Low  string `json:"low"`
	High string `json:"high"`
}

// QSDatabaseConfig is the configuration for the EQEmu server qsdatabase
type QSDatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

// ChatServerConfig is the configuration for the EQEmu server chatserver
type ChatServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// MailServerConfig is the configuration for the EQEmu server mailserver
type MailServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// WorldConfig is the configuration for the EQEmu server world
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

// LoginServerConfig is the configuration for the EQEmu server loginserver
type LoginServerConfig struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Legacy   string `json:"legacy"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// TCPConfig is the configuration for the EQEmu server tcp
type TCPConfig struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// TelnetConfig is the configuration for the EQEmu server telnet
type TelnetConfig struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Enabled string `json:"enabled"`
}

// DatabaseConfig is the configuration for the EQEmu server database
type DatabaseConfig struct {
	DB       string `json:"db"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// FilesConfig is the configuration for the EQEmu server files
type FilesConfig struct {
	Opcodes     string `json:"opcodes"`
	MailOpcodes string `json:"mail_opcodes"`
}

// DirectoriesConfig is the configuration for the EQEmu server directories
type DirectoriesConfig struct {
	Patches string `json:"patches"`
	Opcodes string `json:"opcodes"`
}

// WebAdminConfig is the configuration for the EQEmu server web-admin
type WebAdminConfig struct {
	Application    ApplicationConfig `json:"application"`
	Launcher       LauncherConfig    `json:"launcher"`
	Quests         QuestsConfig      `json:"quests"`
	ServerCodePath string            `json:"serverCodePath"`
}

// ApplicationConfig is the configuration for the EQEmu server application
type ApplicationConfig struct {
	Key   string      `json:"key"`
	Admin AdminConfig `json:"admin"`
}

// AdminConfig is the configuration for the EQEmu server admin
type AdminConfig struct {
	Password string `json:"password"`
}

// LauncherConfig is the configuration for the EQEmu server launcher
type LauncherConfig struct {
	RunLoginServer   bool   `json:"runLoginserver"`
	RunQueryServ     bool   `json:"runQueryServ"`
	IsRunning        bool   `json:"isRunning"`
	StaticZones      string `json:"staticZones"`
	MinZoneProcesses int    `json:"minZoneProcesses"`
	RunSharedMemory  bool   `json:"runSharedMemory"`
}

// QuestsConfig is the configuration for the EQEmu server quests
type QuestsConfig struct {
	HotReload bool `json:"hotReload"`
}

// LoadEQEmuConfig loads a configuration file
func LoadEQEmuConfig(path string) (*EQEmuConfiguration, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var config EQEmuConfiguration
	err = json.NewDecoder(r).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &config, nil
}

// Save saves a configuration file
func (e *EQEmuConfiguration) Save(w *os.File) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(e)
}
