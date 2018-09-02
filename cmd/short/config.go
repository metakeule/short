package main

import (
	"github.com/metakeule/config"
)

var cfg = config.MustNew("short", "1.0.3", "short is shortcut tool for commands. Commands will be stored inside $HOME/.short.json")

var (
	cmdAdd = cfg.MustCommand("add", "adds a shortcut")

	cmdAddArgName = cmdAdd.NewString("name", "name of the shortcut. may contain dots to seperate namespaces",
		config.Required, config.Shortflag('n'))

	cmdAddArgCommand = cmdAdd.NewString("command", "command of the shortcut. refer to a parent shortcut by starting with ':' ",
		config.Required, config.Shortflag('c'))

	cmdAddArgDefaults = cmdAdd.NewString("defaults", "defaults of the shortcut. syntax: 'param1=x,param2=y'",
		config.Shortflag('d'))
)

var (
	cmdRm = cfg.MustCommand("rm", "removes a shortcut")

	cmdRmArgName = cmdRm.NewString("name", "name of the shortcut.",
		config.Required, config.Shortflag('n'))
)

var (
	cmdExec = cfg.MustCommand("exec", "runs a shortcut")

	cmdExecArgName = cmdExec.NewString("name", "name of the shortcut.",
		config.Required, config.Shortflag('n'))

	cmdExecArgParams = cmdExec.NewString("params", "parameters of the shortcut. syntax: 'param1=x,param2=y'",
		config.Shortflag('p'))
)

var (
	cmdLs = cfg.MustCommand("ls", "lists all shortcuts")
)

var (
	cmdShell = cfg.MustCommand("shell", "interactive shell to exec shortcuts")
)
