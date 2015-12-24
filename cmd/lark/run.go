package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// CommandRun implements the "run" action (the default action)
var CommandRun = cli.Command{
	Name:  "run",
	Usage: "run lark task(s) for the project",
	Action: func(c *cli.Context) {
		fmt.Println("running")
	},
}
