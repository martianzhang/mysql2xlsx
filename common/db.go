package common

import (
	"database/sql"
	"fmt"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// GetRows mysql query get rows result
func GetRows() (*sql.Rows, error) {
	var dsn string
	// init database connection
	if Cfg.Socket != "" {
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=%s",
			Cfg.User, Cfg.Password, Cfg.Socket, Cfg.Database, Cfg.Charset)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
			Cfg.User, Cfg.Password, Cfg.Host, Cfg.Port, Cfg.Database, Cfg.Charset)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.Query(Cfg.Query)
}
