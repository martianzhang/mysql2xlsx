package common

import (
	"testing"
)

func TestSaveRows2ELSX(t *testing.T) {
	Cfg.File = testPath + "/test/TestSaveRows2ELSX.xlsx"

	rows, err := GetRows()
	if err != nil {
		panic(err.Error())
	}

	err = saveRows2XLSX(rows)
	if err != nil {
		panic(err.Error())
	}
}

func TestPreviewXlsxFile(t *testing.T) {
	Cfg.File = testPath + "/test/TestSaveRows2ELSX.xlsx"
	err := previewXlsx()
	if err != nil {
		panic(err.Error())
	}
}
