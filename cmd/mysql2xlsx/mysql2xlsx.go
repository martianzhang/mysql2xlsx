package main

import (
	"mysql2xlsx/common"
)

func main() {
	// parse config
	cfg, err := common.ParseFlag()
	if err != nil {
		panic(err.Error())
	}

	// execute sql and get all result rows
	rows, err := common.GetRows(cfg)
	if err != nil {
		panic(err.Error())
	}

	// save rows result
	err = common.SaveRows(cfg.File, rows)
	if err != nil {
		panic(err.Error())
	}
}
