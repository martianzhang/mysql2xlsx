package common

import (
	"database/sql"
	"errors"
	"strings"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// SaveRows ...
func SaveRows(rows *sql.Rows) error {
	var err error
	var suffix string

	tup := strings.Split(Cfg.File, ".")
	suffix = strings.ToLower(tup[len(tup)-1])

	switch suffix {
	case "stdout", "":
		printRowsAsASCII(rows)
	case "tsv", "txt": // tab-separated values
		err = saveRows2CSV(rows, '\t')
	case "psv": // pipe-separated values
		err = saveRows2CSV(rows, '|')
	case "csv": // comma-separated values
		err = saveRows2CSV(rows, ',')
	case "xlsx":
		err = saveRows2XLSX(rows)
	case "sql":
		err = saveRows2SQL(rows)
	default:
		err = errors.New("not support extension: " + suffix)
	}

	return err
}
