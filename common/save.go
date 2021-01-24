package common

import (
	"database/sql"
	"errors"
	"fmt"
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
	default:
		err = errors.New("unknown file extension: " + suffix)
	}

	if err != nil && suffix != "stdout" {
		fmt.Println("save result into: ", Cfg.File)
	}

	return err
}
