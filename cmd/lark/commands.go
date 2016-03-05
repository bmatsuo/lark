package main

import (
	"github.com/codegangsta/cli"
	"github.com/yuin/gopher-lua"
)

// Verbose causes verbose error logging when set to true.
var Verbose = new(bool)

// Commands contains the list of commands available in lark.
var Commands = []cli.Command{
	CommandRun,
	CommandList,
	CommandREPL,
	CommandLua,
}

// Command is a helper for creating a cli.Command that relies on a Context for
// setup.
func Command(fn func(lark *Context, cmd *cli.Command)) cli.Command {
	lark := NewContext(nil)
	cmd := new(cli.Command)
	fn(lark, cmd)
	return *cmd
}

// Context is a context for a lark command.
type Context struct {
	*cli.Context
	Lua     *lua.LState
	verbose *bool
}

// Verbose returns true if verbose output has been enabled.
func (c *Context) Verbose() bool {
	return c.verbose != nil && *c.verbose
}

// Action returns a function usable as the action for a cli.Command.
func (c *Context) Action(fn func(*Context)) func(*cli.Context) {
	return func(_c *cli.Context) {
		c.Context = _c
		fn(c)
	}
}

// NewContext returns a new context that wraps c.
func NewContext(c *cli.Context) *Context {
	return &Context{
		Context: c,
		verbose: Verbose,
	}
}
