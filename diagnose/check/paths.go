package check

import (
	"strings"

	"github.com/xackery/overseer/pkg/config"
	"github.com/xackery/overseer/pkg/message"
)

func Paths(config *config.OverseerConfiguration) error {
	if config.ServerPath == "" {
		message.Badf("overseer.ini server_path is empty ")
		message.Link("https://o.eqcodex.com/100")
	}
	if config.BinPath == "" {
		message.Badf("overseer.ini bin_path is empty ")
		message.Link("https://o.eqcodex.com/101")
	}
	if config.ZoneCount == 0 {
		message.Badf("overseer.ini zone_count is 0 ")
		message.Link("https://o.eqcodex.com/102")
	}

	switch config.Setup {
	case "default":
	case "docker":
		if config.DockerNetwork == "" {
			message.Badf("overseer.ini docker_network is empty ")
			message.Link("https://o.eqcodex.com/103")
		}
		badValueCheck := strings.ToLower(config.DockerNetwork)
		switch badValueCheck {
		case "bridge":
			message.Badf("overseer.ini docker_network is set to bridge (duplicate)")
			message.Link("https://o.eqcodex.com/104")
		case "host":
			message.Badf("overseer.ini docker_network is set to host (duplicate)")
			message.Link("https://o.eqcodex.com/104")
		case "none":
			message.Badf("overseer.ini docker_network is set to none (duplicate)")
			message.Link("https://o.eqcodex.com/104")
		}
	}

	return nil
}
