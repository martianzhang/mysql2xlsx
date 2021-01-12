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

	testConfig = Config{
		User:     "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "test",
		Charset:  "utf8mb4",

		Query: "select 1",
	}
}

func TestParseDefaultsExtraFile(t *testing.T) {
	_, err := parseDefaultsExtraFile(testPath + "/test/my.cnf")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestParseFlagParseFlag(t *testing.T) {
	_, err := ParseFlag()
	if err != nil && err != io.EOF {
		t.Error(err.Error())
	}
}
