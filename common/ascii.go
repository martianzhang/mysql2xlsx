package common

import (
	"database/sql"
	"fmt"
	"os"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)

// printRowsAsASCII print rows as ascii table
func printRowsAsASCII(rows *sql.Rows) {

	table := tablewriter.NewWriter(os.Stdout)

	// set table header
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	table.SetHeader(columns)

	// set every rows
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
			fmt.Println(err.Error())
			return
		}

		values := make([]string, len(columns))
		for i, col := range columns {
			values[i] = col.String
		}
		table.Append(values)
	}

	// print table
	table.Render()
}
