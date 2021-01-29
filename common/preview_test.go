package common

import (
	"fmt"
	"testing"
)

func TestPreview(t *testing.T) {
	files := []string{
		testPath + "/test/TestSaveRows.csv",
		testPath + "/test/TestSaveRows.tsv",
		testPath + "/test/TestSaveRows.psv",
		testPath + "/test/TestSaveRows.xlsx",
	}
	oldPreview := Cfg.Preview
	Cfg.Preview = 10
	for _, file := range files {
		Cfg.File = file
		fmt.Println("# Preview: ", Cfg.File)
		err := Preview()
		if err != nil {
			panic(err)
		}
	}
	Cfg.Preview = oldPreview
}
