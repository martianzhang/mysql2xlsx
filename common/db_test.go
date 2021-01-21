package common

import (
	"testing"
)

func TestGetRows(t *testing.T) {
	_, err := GetRows()
	if err != nil {
		panic(err.Error())
	}
}
