package db

import (
	"database/sql"
	"fmt"
	"os"
)

func ApplySQLFile(conn *sql.DB, path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", path, err)
	}
	if _, err := conn.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("execute sql file %s: %w", path, err)
	}
	return nil
}
