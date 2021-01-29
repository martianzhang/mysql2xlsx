package common

import (
	"database/sql"
	"fmt"

	"github.com/tealeg/xlsx/v3"
)

// Excel limits
const (
	ExcelMaxRows      = 1048576
	ExcelMaxColumns   = 16384
	ExcelMaxCellChars = 32767
)

// saveRows2XLSX save rows result into xlsx format file
func saveRows2XLSX(rows *sql.Rows) error {
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
	if len(columns) > ExcelMaxColumns {
		return fmt.Errorf("excel max columns(%d) exceeded", ExcelMaxColumns)
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
	for i := 1; rows.Next(); i++ {
		if i > ExcelMaxRows {
			return fmt.Errorf("excel max rows(%d) exceeded", ExcelMaxRows)
		}
		rows.Scan(scanArgs...)
		sheetRow := sheet.AddRow()
		for _, v := range values {
			cell := sheetRow.AddCell()
			if len(v) > ExcelMaxCellChars {
				return fmt.Errorf("excel max cell characters(%d) exceeded", ExcelMaxCellChars)
			}
			cell.Value = string(v)
		}
	}

	// save to file
	err = file.Save(Cfg.File)
	if err != nil {
		return err
	}

	// preview xlsx file
	if Cfg.Preview > 0 {
		err = previewXlsx()
	}
	return err
}

// PreviewXlsx ...
func previewXlsx() error {
	if Cfg.Preview == 0 {
		return nil
	}

	opts := xlsx.RowLimit(Cfg.Preview)
	wb, err := xlsx.OpenFile(Cfg.File, opts)
	if err != nil {
		return err
	}

	for _, sh := range wb.Sheets {
		for i := 0; i < Cfg.Preview && i < sh.MaxRow; i++ {
			row, err := sh.Row(i)
			if err != nil {
				return err
			}
			for j := 0; j < sh.MaxCol; j++ {
				fmt.Print(row.GetCell(j), "\t")
			}
			fmt.Println() // add line feed
		}
		break
	}
	return nil
}
