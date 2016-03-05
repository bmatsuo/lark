/*
Command lark executes project tasks defined with lua scripts.  Lark isolates
the lua modules available to scripts to modules in the relative directory path
./lark_modules/ to ensure portability of project tasks across developer
machines.

The lark command locates tasks defined in the ./lark.lua file or otherwise
directly under the directory ./lark_tasks/.  Tasks can have names (either
explicitly given or otherwise inferred) or patterns.  Pattern matching tasks
will match a set of names defined by a regular expression.

Names are matched against available tasks with a strict precedence.  Explicitly
named tasks will match the same name with the highest priority.  Any task with
an inferred name will match the same name with the second highest priority.
Pattern matching tasks have the lowest priority and will match names in the
order they were defined.

Tasks can be executed by calling the lua function lark.run() in a script, using
the lark subcommand "run".  When given no arguments, run will execute the first
named task that was defined, or a task specified by setting the "default"
variable in the "lark.task" lua module.

	local task = require('lark.task')
	task1 = task .. function() print('task1') end
	task2 = task .. function() print('task2') end
	lark.run()
	task.default = 'task2'
	lark.run()
	lark.run('task1')

The above script will print a line containing text "task1" followed by a line
containing "task2" and finally a line containing "task1" again.


Command Reference

Command reference documentation is available through the "help" subcommand.

	lark help

The documentation for a specific subcommand is available through the help
command or by passing the subcommand the -h (or --help) flag.

	lark run -h
	lark help run


Lua Reference

Lua API documentation is available through the help() function in the embedded REPL.

	lark repl
	> help()
	> help(lark)
*/
package main

import (
	"log"
	"os"

	"github.com/bmatsuo/lark/larkmeta"
	"github.com/codegangsta/cli"
	"github.com/mattn/go-isatty"
	"github.com/yuin/gopher-lua"
)

// IsTTY is true if standard error is connected to a terminal. This is taken to
// mean that lark was executed from the command line and is not being logged to
// a file.
//
// BUG: The assumptions made due to IsTTY cannot be overridden.
var IsTTY = isatty.IsTerminal(os.Stderr.Fd())

// MainHelp is the top-level hop documentation.
var MainHelp = `

    If no builtin command is provided the "run" command is executed with any
    task arguments provided.  For more information see the command
    documentation

        lark run -h
`

func main() {
	if IsTTY {
		logflags := log.Flags()
		logflags &^= log.Ldate | log.Ltime
		log.SetFlags(logflags)
	}

	// Set search path for lua modules.  The search path must be completely
	// contained by the working directory to help ensure repeatable builds
	// across machines.
	lua.LuaPathDefault = "./lark_modules/?.lua;./lark_modules/?/init.lua"
	os.Setenv(lua.LuaPath, "")

	cli.VersionFlag.Name = "version"

	app := cli.NewApp()
	app.Name = "lark"
	app.Usage = "Run repeated project tasks"
	app.ArgsUsage = MainHelp
	app.Version = larkmeta.Version
	app.Authors = []cli.Author{
		{
			Name:  "Bryan Matsuo",
			Email: "bryan.matsuo@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) {
		args := []string{os.Args[0], "run"}
		args = append(args, c.Args()...)
		app.Run(args)
	}
	app.Commands = Commands

	app.Run(os.Args)
}
