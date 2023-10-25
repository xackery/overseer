package connect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xackery/overseer/pkg/config"
)

func DB(emuCfg *config.EQEmuConfiguration) (*sql.DB, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		emuCfg.Server.Database.Username,
		emuCfg.Server.Database.Password,
		emuCfg.Server.Database.Host,
		emuCfg.Server.Database.Port,
		emuCfg.Server.Database.DB))
	if err != nil {
		return nil, fmt.Errorf("open %s@%s: %w", emuCfg.Server.Database.Username, emuCfg.Server.Database.Host, err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping %s@%s: %w", emuCfg.Server.Database.Username, emuCfg.Server.Database.Host, err)
	}

	_, err = conn.Exec("USE " + emuCfg.Server.Database.DB)
	if err != nil {
		return nil, fmt.Errorf("use %s: %w", emuCfg.Server.Database.DB, err)
	}
	return conn, nil
}
