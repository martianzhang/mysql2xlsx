package common

import (
	"database/sql"
	"encoding/csv"
	"os"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// saveRows2CSV save rows result into csv format file
func saveRows2CSV(filePath string, rows *sql.Rows) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// set table header with column name
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	w.Write(columns)

	// set every table rows
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		columns := make([]sql.NullString, len(columns))
		cols := make([]interface{}, len(columns))
		for i := range columns {
			cols[i] = &columns[i]
		}

		if err := rows.Scan(cols...); err != nil {
			return err
		}

		values := make([]string, len(columns))
		for i, col := range columns {
			values[i] = col.String
		}
		w.Write(values)
	}
	return err
}
