package common

import (
	"testing"
)

func TestSaveRows2ELSX(t *testing.T) {

	rows, err := GetRows(testConfig)
	if err != nil {
		t.Error(err.Error())
	}

	err = saveRows2XLSX(testPath+"/test/1.xlsx", rows)
	if err != nil {
		t.Error(err.Error())
	}
}
