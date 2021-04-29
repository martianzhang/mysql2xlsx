package common

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// saveRows2SQL save rows result into sql file
func saveRows2SQL(rows *sql.Rows) error {
	file, err := os.Create(Cfg.File)
	if err != nil {
		return err
	}
	defer file.Close()

	// get table name
	tableName := strings.Split(filepath.Base(Cfg.File), ".")[0]

	// get column name
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// write every row into sql
	w := bufio.NewWriter(file)
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
			if !col.Valid {
				values[i] = "NULL"
			} else {
				values[i] = strconv.Quote(col.String)
			}
			// hex-blob
			// values[i] = strconv.Quote(fmt.Sprintf("%x", col))
		}
		if _, err := w.WriteString(
			fmt.Sprintf("INSERT INTO `%s` VALUES (%s);\n",
				tableName, strings.Join(values, ", "))); err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	// preview file
	if Cfg.Preview > 0 {
		err = previewSQL()
	}
	return err
}

func previewSQL() error {
	if Cfg.Preview == 0 {
		return nil
	}

	fd, err := os.Open(Cfg.File)
	if err != nil {
		return err
	}
	defer fd.Close()

	var line int
	s := bufio.NewScanner(fd)
	for s.Scan() {
		if line >= Cfg.Preview {
			break
		}
		fmt.Println(s.Text())
		line++
	}

	return s.Err()
}
