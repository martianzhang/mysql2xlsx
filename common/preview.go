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
	case "tsv", "txt": // tab-separated values
		err = previewCSV() // TODO:
	case "psv": // pipe-separated values
		err = previewCSV() // TODO:
	case "csv": // comma-separated values
		err = previewCSV()
	case "xlsx":
		err = previewXlsx()
	default:
		err = errors.New("unknown file extension: " + suffix)
	}

	return err
}
