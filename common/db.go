package common

import (
	"database/sql"
	"fmt"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// GetRows mysql query get rows result
func GetRows(cfg Config) (*sql.Rows, error) {

	// init database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.Query(cfg.Query)
}
