package check

import (
	"os"

	"github.com/xackery/overseer/pkg/config"
)

// IsQuestInstalled checks if quest is installed
func IsQuestInstalled(config *config.OverseerConfiguration) bool {
	questPath := config.ServerPath + "/quests"
	_, err := os.Stat(questPath)
	return err == nil
}
