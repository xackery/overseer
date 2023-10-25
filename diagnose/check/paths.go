package check

import (
	"os"
	"runtime"
	"strings"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
)

func Paths(cfg *config.OverseerConfiguration, emuCfg *config.EQEmuConfiguration) error {
	message.OKReset()
	if cfg.ServerPath == "" {
		message.Badf("overseer.ini server_path is empty ")
		message.Link("https://o.eqcodex.com/100")

	}
	if cfg.BinPath == "" {

		message.Badf("overseer.ini bin_path is empty ")
		message.Link("https://o.eqcodex.com/101")

	}
	essentialFiles := []string{
		"world",
		"zone",
		"shared_memory",
	}
	winExe := ""
	if runtime.GOOS == "windows" {
		winExe = ".exe"
	}
	for _, essentialFile := range essentialFiles {
		_, err := os.Stat(cfg.BinPath + "/" + essentialFile + winExe)
		if err != nil {
			message.Badf("overseer.ini bin_path/%s not found ", essentialFile+winExe)
			message.Link("https://o.eqcodex.com/101")

		}
	}
	if cfg.ZoneCount == 0 {
		message.Badf("overseer.ini zone_count is 0 ")
		message.Link("https://o.eqcodex.com/102")

	}

	switch cfg.Setup {
	case "default":
	case "docker":
		if cfg.DockerNetwork == "" {
			message.Badf("overseer.ini docker_network is empty ")
			message.Link("https://o.eqcodex.com/103")
		}
		badValueCheck := strings.ToLower(cfg.DockerNetwork)
		switch badValueCheck {
		case "bridge":
			message.Badf("overseer.ini docker_network is set to bridge (duplicate) ")
			message.Link("https://o.eqcodex.com/104")
		case "host":
			message.Badf("overseer.ini docker_network is set to host (duplicate) ")
			message.Link("https://o.eqcodex.com/104")
		case "none":
			message.Badf("overseer.ini docker_network is set to none (duplicate) ")
			message.Link("https://o.eqcodex.com/104")
		}
	}

	if emuCfg.Server.Files.MailOpcodes == "" {
		message.Badf("eqemu_config.json files.mail_opcodes is empty ")
		message.Link("https://o.eqcodex.com/106")
	} else {
		_, err := os.Stat(emuCfg.Server.Files.MailOpcodes)
		if err != nil {
			_, err = os.Stat(cfg.ServerPath + "/" + emuCfg.Server.Files.MailOpcodes)
			if err != nil {
				message.Badf("eqemu_config.json files.mail_opcodes %s not found ", emuCfg.Server.Files.MailOpcodes)
				message.Link("https://o.eqcodex.com/106")

			}
		}
	}

	if emuCfg.Server.Files.Opcodes == "" {
		message.Badf("eqemu_config.json files.opcodes is empty ")
		message.Link("https://o.eqcodex.com/106")

	} else {
		_, err := os.Stat(emuCfg.Server.Files.Opcodes)
		if err != nil {
			_, err = os.Stat(cfg.ServerPath + "/" + emuCfg.Server.Files.Opcodes)
			if err != nil {
				message.Badf("eqemu_config.json files.opcodes %s not found ", emuCfg.Server.Files.Opcodes)
				message.Link("https://o.eqcodex.com/106")

			}
		}
	}

	if emuCfg.Server.Directories.Patches == "" {
		message.Badf("eqemu_config.json directories.patches is empty ")
		message.Link("https://o.eqcodex.com/106")

	} else {
		_, err := os.Stat(emuCfg.Server.Directories.Patches)
		if err != nil {
			_, err = os.Stat(cfg.ServerPath + "/" + emuCfg.Server.Directories.Patches)
			if err != nil {
				message.Badf("eqemu_config.json directories.patches %s not found ", emuCfg.Server.Directories.Patches)
				message.Link("https://o.eqcodex.com/106")

			}
		}
	}

	if emuCfg.Server.Directories.Opcodes == "" {
		message.Badf("eqemu_config.json directories.opcodes is empty ")
		message.Link("https://o.eqcodex.com/106")

	} else {
		_, err := os.Stat(emuCfg.Server.Directories.Opcodes)
		if err != nil {
			_, err = os.Stat(cfg.ServerPath + "/" + emuCfg.Server.Directories.Opcodes)
			if err != nil {
				message.Badf("eqemu_config.json directories.opcodes %s not found ", emuCfg.Server.Directories.Opcodes)
				message.Link("https://o.eqcodex.com/106")

			}
		}
	}

	if message.IsOK() {
		message.OK("Assets OK")
	}

	return nil
}
