package common

import "testing"

func TestSaveRows(t *testing.T) {
	formats := []string{
		"",
		"stdout",
		"../test/1.csv",
		"../test/1.xlsx",
	}
	for _, format := range formats {
		rows, err := GetRows(testConfig)
		if err != nil {
			t.Error(err.Error())
		}

		err = SaveRows(format, rows)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
