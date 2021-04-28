package common

import (
	"testing"
)

func TestSaveRows2SQL(t *testing.T) {
	Cfg.File = testPath + "/test/TestSaveRows2SQL.sql"

	rows, err := GetRows()
	if err != nil {
		panic(err.Error())
	}

	err = saveRows2SQL(rows)
	if err != nil {
		panic(err.Error())
	}
}

func TestPreviewSQL(t *testing.T) {
	Cfg.File = testPath + "/test/TestSaveRows2SQL.sql"
	err := previewSQL()
	if err != nil {
		panic(err.Error())
	}
}
