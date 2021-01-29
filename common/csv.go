package common

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// UTF8BOM BOM HEADER
var UTF8BOM = []byte{0xEF, 0xBB, 0xBF}

// saveRows2CSV save rows result into csv format file
func saveRows2CSV(rows *sql.Rows, comma rune) error {
	file, err := os.Create(Cfg.File)
	if err != nil {
		return err
	}
	defer file.Close()

	// 兼容 Windows 系统，文件头写入 UTF8 BOM，防止中文乱码。
	// windows 环境下导出的 csv 文件默认添加 UTF8 BOM。
	// 添加 BOM 对 less, awk 等 *nix 系统命令并不友好，因此仅对特定的文件名生效。
	// Linux 环境删除文件 UTF8 BOM 头命令：dos2unix xxx.csv
	if Cfg.BOM {
		_, err = file.Write(UTF8BOM)
		if err != nil {
			return err
		}
	}

	w := csv.NewWriter(file)
	w.Comma = comma
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

	if err != nil {
		return err
	}

	// preview file
	if Cfg.Preview > 0 {
		err = previewCSV()
	}
	return err
}

func previewCSV() error {
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
		if line > Cfg.Preview {
			break
		}
		fmt.Println(s.Text())
		line++
	}

	return s.Err()
}
