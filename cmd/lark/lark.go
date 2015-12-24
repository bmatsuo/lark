package main

import (
	"os"

	"golang.org/x/net/context"

	"github.com/bmatsuo/lark/cmd"
)

// MainCmd is the root lark command that will call subcommands as necessary.
var MainCmd = &cmd.Cmd{
	Name: os.Args[0],
	Run:  RunLark,
}

// RunLark performs logic for the default invocation `lark`.
func RunLark(c context.Context, prefix string, args []string) error {
	return nil
}
