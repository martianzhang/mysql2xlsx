package main

import (
	"runtime"

	"mysql2xlsx/common"
)

func main() {
	// limit cpu usage
	runtime.GOMAXPROCS(1)

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
