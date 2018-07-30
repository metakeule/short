package main

import (
	"fmt"
	"github.com/metakeule/config"
	"os"
	"path/filepath"
)

var SHORTCUT_FILE string
var run runner

func init() {
	SHORTCUT_FILE = filepath.Join(os.Getenv("HOME"), ".short.json")
	run = map[*config.Config]func() error{
		cmdExec:  run.Exec,
		cmdAdd:   run.Add,
		cmdLs:    run.Ls,
		cmdRm:    run.Rm,
		cmdShell: run.Shell,
	}
}

func main() {
	err := run.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
		os.Exit(1)
	}
}
