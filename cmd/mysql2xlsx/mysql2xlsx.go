package main

import (
	"os"
	"runtime"

	"mysql2xlsx/common"
)

func main() {
	// limit cpu usage
	runtime.GOMAXPROCS(1)

	// parse config
	err := common.ParseFlag()
	if err != nil {
		panic(err.Error())
	}

	// xlsx file preview
	if common.Cfg.Preview != 0 && common.Cfg.File != "" {
		if _, err = os.Stat(common.Cfg.File); err == nil {
			err = common.Preview()
			if err != nil {
				panic(err.Error())
			}
			return
		}
	}

	// execute sql and get all result rows
	rows, err := common.GetRows()
	if err != nil {
		panic(err.Error())
	}

	// save rows result
	err = common.SaveRows(rows)
	if err != nil {
		panic(err.Error())
	}
}
