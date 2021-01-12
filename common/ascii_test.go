package common

import "testing"

func TestPrintRowsAsASCII(t *testing.T) {
	rows, err := GetRows(testConfig)
	if err != nil {
		t.Error(err.Error())
	}

	printRowsAsASCII(rows)
}
