package common

import (
	"testing"
)

func TestSaveRows2CSV(t *testing.T) {
	rows, err := GetRows(testConfig)
	if err != nil {
		t.Error(err.Error())
	}

	err = saveRows2CSV(testPath+"/test/1.csv", rows)
	if err != nil {
		t.Error(err.Error())
	}
}
