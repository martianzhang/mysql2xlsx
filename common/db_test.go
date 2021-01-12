package common

import (
	"testing"
)

func TestGetRows(t *testing.T) {
	_, err := GetRows(testConfig)
	if err != nil {
		t.Error(err.Error())
	}
}
