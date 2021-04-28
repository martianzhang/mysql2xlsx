package common

import (
	"errors"
	"strings"
)

// Preview preview export file
func Preview() error {
	var err error
	var suffix string

	tup := strings.Split(Cfg.File, ".")
	suffix = strings.ToLower(tup[len(tup)-1])

	switch suffix {
	case "stdout", "":
	case "csv", "psv", "tsv", "txt":
		err = previewCSV()
	case "xlsx":
		err = previewXlsx()
	case "sql":
		err = previewSQL()
	default:
		err = errors.New("not support extension: " + suffix)
	}

	return err
}
