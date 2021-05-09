package main

import (
	"os"
	"runtime"
	"syscall"

	"mysql2xlsx/common"

	"golang.org/x/sys/unix"
)

func main() {
	// limit cpu, memory usage
	resourceLimit()

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

func resourceLimit() {
	// limit cpu usage
	runtime.GOMAXPROCS(1) // 1 core

	// ulimit -Sv, virtual memory
	var maxVirtualMemoryBytes uint64 = 1024 * 1024 * 1024 // 1GB
	err := unix.Setrlimit(syscall.RLIMIT_AS, &unix.Rlimit{
		Cur: maxVirtualMemoryBytes,
		Max: maxVirtualMemoryBytes,
	})
	if err != nil {
		panic(err.Error())
	}
}
