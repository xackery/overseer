package check

import (
	"os"

	"github.com/xackery/overseer/pkg/config"
)

// IsQuestInstalled checks if quest is installed
func IsQuestInstalled(cfg *config.OverseerConfiguration) bool {
	questPath := cfg.ServerPath + "/quests"
	_, err := os.Stat(questPath)
	return err == nil
}
