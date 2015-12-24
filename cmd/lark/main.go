package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "lark"
	app.Usage = "Run repeated tasks for this project"
	app.Action = func(c *cli.Context) {
		fmt.Println("LARK!")
	}
	app.Commands = Commands

	app.Run(os.Args)
}
