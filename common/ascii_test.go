package common

import "testing"

func TestPrintRowsAsASCII(t *testing.T) {
	rows, err := GetRows()
	if err != nil {
		panic(err.Error())
	}

	printRowsAsASCII(rows)
}
