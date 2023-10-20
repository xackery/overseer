package deep

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xackery/overseer/lib/connect"
	"github.com/xackery/overseer/share/config"
	"github.com/xackery/overseer/share/message"
)

// Table runs a deep diagnostic on the server's tables
func Table(cfg *config.OverseerConfiguration, ecfg *config.EQEmuConfiguration) error {
	db, err := connect.DB(ecfg)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	type tableCheck struct {
		name string
		rows int
	}

	tables := []tableCheck{
		{name: "account", rows: 1},
	}

	message.OKReset()
	for _, table := range tables {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		exists, err := tableExists(ctx, db, table.name)
		if err != nil {
			message.Badf("table %s: %s\n", table.name, err)
			continue
		}
		if !exists {
			message.Badf("table %s does not exist\n", table.name)
			continue
		}
		count, err := tableRowCount(ctx, db, table.name)
		if err != nil {
			message.Badf("table query row %s: %s\n", table.name, err)
			continue
		}
		if count < table.rows {
			message.Badf("table %s has %d rows, expected greater than %d\n", table.name, count, table.rows)
		}
	}
	if message.IsOK() {
		message.OK("Tables OK")
	}

	return nil
}

func tableExists(ctx context.Context, db *sql.DB, table string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?)", "eqemu", table).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query row: %w", err)
	}
	return exists, nil
}

func tableRowCount(ctx context.Context, db *sql.DB, table string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}
	return count, nil
}
