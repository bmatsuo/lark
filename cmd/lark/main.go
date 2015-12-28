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
