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
func SaveRows(file string, rows *sql.Rows) error {
	var err error
	var suffix string
	file = strings.TrimSpace(file)
	if file == "" {
		suffix = "stdout"
	} else {
		tup := strings.Split(file, ".")
		suffix = strings.ToLower(tup[len(tup)-1])
	}
	switch suffix {
	case "stdout":
		printRowsAsASCII(rows)
	case "csv":
		err = saveRows2CSV(file, rows)
	case "xlsx":
		err = saveRows2XLSX(file, rows)
	default:
		err = errors.New("unknown file extension: " + suffix)
	}

	if err != nil && suffix != "stdout" {
		fmt.Println("save result into: ", file)
	}

	return err
}
