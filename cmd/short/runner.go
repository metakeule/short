package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/metakeule/config"
	"github.com/metakeule/short"
)

type runner map[*config.Config]func() error

func (r runner) Exec() error {
	all, err := load()
	if err != nil {
		return err
	}

	var runtimeParams map[string]string
	runtimeParams, err = paramsStringToMap(cmdExecArgParams.Get())
	if err != nil {
		return err
	}

	return short.Exec(cmdExecArgName.Get(), all, runtimeParams)
}

func (r runner) Add() error {
	all, err := load()
	if err != nil {
		return err
	}

	var c short.Cut
	c.Name = cmdAddArgName.Get()
	c.Command = cmdAddArgCommand.Get()
	c.Defaults, err = paramsStringToMap(cmdAddArgDefaults.Get())
	if err != nil {
		return err
	}

	err = short.Add(c, all)

	if err != nil {
		return err
	}

	return save(all)
}

func (r runner) Ls() error {
	all, err := load()
	if err != nil {
		return err
	}

	cuts := short.Sort(all)

	if len(cuts) == 0 {
		fmt.Fprintf(os.Stdout, "no shortcuts defined")
		return nil
	}

	for _, c := range cuts {
		fmt.Fprintf(os.Stdout, "%s\n\t---> %s\n\n", c.Name, c.Command)
	}

	return nil
}

func (r runner) Rm() error {
	all, err := load()
	if err != nil {
		return err
	}

	name := strings.TrimSpace(cmdRmArgName.Get())
	if name == "" {
		return fmt.Errorf("empty name not allowed")
	}

	delete(all, name)

	return save(all)
}

func (r runner) Shell() error {
	s, err := NewScreen()

	if err != nil {
		return err
	}

	return s.Run()
}

func (r runner) Run() error {

	err := cfg.Run()

	// shortcut errbreak
	if err != nil {
		return err
	}

	activeCMD := cfg.ActiveCommand()
	if activeCMD == nil {
		activeCMD = cmdShell
	}

	fn := r[activeCMD]
	return fn()
}
