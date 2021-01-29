package common

import (
	"testing"
)

func TestSaveRows2CSV(t *testing.T) {
	Cfg.File = testPath + "/test/TestSaveRows2CSV.csv"

	rows, err := GetRows()
	if err != nil {
		panic(err.Error())
	}

	err = saveRows2CSV(rows, ',')
	if err != nil {
		panic(err.Error())
	}
}

func TestPreviewCSV(t *testing.T) {
	Cfg.File = testPath + "/test/TestSaveRows2CSV.csv"
	err := previewCSV()
	if err != nil {
		panic(err.Error())
	}
}
