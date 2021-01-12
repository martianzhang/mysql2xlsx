package common

import (
	"database/sql"

	"github.com/tealeg/xlsx"
)

// saveRows2XLSX save rows result into xlsx format file
func saveRows2XLSX(filePath string, rows *sql.Rows) error {
	file := xlsx.NewFile()
	// create new sheet
	sheet, err := file.AddSheet("result")
	if err != nil {
		return err
	}

	// set table header with column name
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	sheetHeader := sheet.AddRow()
	for _, name := range columns {
		cell := sheetHeader.AddCell()
		cell.Value = name
	}

	// set every table rows
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		sheetRow := sheet.AddRow()
		for _, v := range values {
			cell := sheetRow.AddCell()
			cell.Value = string(v)
		}
	}

	// save to file
	return file.Save(filePath)
}
