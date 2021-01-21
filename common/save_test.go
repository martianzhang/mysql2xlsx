package common

import (
	"testing"
)

func TestSaveRows(t *testing.T) {
	files := []string{
		"stdout",
		testPath + "/test/TestSaveRows.csv",
		testPath + "/test/TestSaveRows.xlsx",
	}

	for _, file := range files {
		Cfg.File = file
		rows, err := GetRows()
		if err != nil {
			panic(err.Error())
		}

		err = SaveRows(rows)
		if err != nil {
			panic(err.Error())
		}
	}
}
