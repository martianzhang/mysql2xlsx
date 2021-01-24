package common

import (
	"io"
	"path/filepath"
	"runtime"
	"testing"
)

var testConfig Config
var testPath string

func init() {

	_, filename, _, _ := runtime.Caller(0)
	testPath = filepath.Dir(filepath.Dir(filename))

	Cfg = Config{
		User:     "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "information_schema",
		Charset:  "utf8mb4",

		Query: "select 'test', '中文'",
	}
}

func TestParseDefaultsExtraFile(t *testing.T) {
	err := parseDefaultsExtraFile(testPath + "/test/my.cnf")
	if err != nil {
		panic(err.Error())
	}
}

func TestParseFlag(t *testing.T) {
	err := ParseFlag()
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
}
