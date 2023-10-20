package confirm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xackery/overseer/lib/connect"
	"github.com/xackery/overseer/share/config"
)

func DBConnects(cfg *config.OverseerConfiguration, emuCfg *config.EQEmuConfiguration) error {
	db, err := connect.DB(emuCfg)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
