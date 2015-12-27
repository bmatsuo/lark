package main

import (
	"os"

	"github.com/codegangsta/cli"
)

// MainHelp is the top-level hop documentation.
var MainHelp = `

    If no builtin command is provided the "run" command is executed with any
    task arguments provided.  For more information see the command
    documentation

        lark run -h
`

func main() {
	cli.VersionFlag.Name = "version"

	app := cli.NewApp()
	app.Name = "lark"
	app.Usage = "Run repeated project tasks"
	app.ArgsUsage = MainHelp

	app.Action = func(c *cli.Context) {
		args := []string{os.Args[0], "run"}
		args = append(args, c.Args()...)
		app.Run(args)
	}
	app.Commands = Commands

	app.Run(os.Args)
}
