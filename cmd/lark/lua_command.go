package main

import (
	"io"
	"log"
	"os"

	"github.com/bmatsuo/lark/larkmeta"
	"github.com/bmatsuo/lark/project"
	"github.com/codegangsta/cli"
	"github.com/yuin/gopher-lua"
)

// CommandLua implements the "lua" command and runs a lark-flavored Lua
// interpreter.
var CommandLua = Command(func(lark *Context, cmd *cli.Command) {
	cmd.Name = "lua"
	cmd.Usage = "Run an interactive intepreter"
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Usage: "Code to evaluate instead of reading from a file (or stdin).",
		},
	}
	cmd.Action = lark.Action(Lua)
})

// Lua loads a lua vm with the lark library and executes Lua scripts or
// expressions.
func Lua(c *Context) {
	luaFiles, err := project.FindTaskFiles(".")
	if err != nil {
		log.Fatal(err)
	}

	luaConfig := &LuaConfig{}
	c.Lua, err = LoadVM(luaConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Lua.Close()

	err = InitLark(c, luaFiles)
	if err != nil {
		log.Fatal(err)
	}

	err = RunLua(c)
	if err != nil {
		log.Fatal(err)
	}
}

// RunLua executes lua code specified in c.Args().
func RunLua(c *Context) error {
	luaExpr := c.String("c")
	if luaExpr != "" {
		fn, err := c.Lua.LoadString(luaExpr)
		if err != nil {
			return err
		}
		c.Lua.Push(fn)
		args := c.Args()
		for _, arg := range args {
			c.Lua.Push(lua.LString(arg))
		}
		return c.Lua.PCall(len(args), 0, c.Lua.NewFunction(errTraceback))
	}

	args := c.Args()
	if len(args) == 0 {
		return LuaInteractive(c)
	}

	// determine the input source (stdin or a file)
	var luaName string
	var luaReader io.Reader
	var luaArgs []string
	if len(args) == 0 {
		luaName = "stdin"
		luaReader = os.Stdin
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()

		luaName = args[0]
		luaReader = f
		luaArgs = args[1:]
	}

	// read and prepare to execute the input
	fn, err := c.Lua.Load(luaReader, luaName)
	if err != nil {
		return err
	}
	c.Lua.Push(fn)
	for _, arg := range luaArgs {
		c.Lua.Push(lua.LString(arg))
	}

	return c.Lua.PCall(len(luaArgs), 0, c.Lua.NewFunction(errTraceback))
}

// LuaInteractive runs an interactive interpreter for the embedded Lua
// environment.  LuaInteractive performs extra setup for the interpereter and
// prints relevant copyright information before reading and interpreting input
// from a terminal device.
func LuaInteractive(c *Context) error {
	// when connected to a terminal perform special setup for an interactive
	// session.
	if IsTTYStdin {
		c.Lua.Push(c.Lua.GetGlobal("require"))
		c.Lua.Push(lua.LString("doc"))
		err := c.Lua.PCall(1, 1, nil)
		if err != nil {
			return err
		}
		docModule := c.Lua.Get(c.Lua.GetTop())
		c.Lua.Pop(1)
		c.Lua.SetGlobal("help", c.Lua.GetField(docModule, "help"))
		c.Lua.SetField(docModule, "default", lua.LString(REPLHelp()))

		var authors = "The Lark authors"
		if len(larkmeta.Authors) > 0 {
			authors = larkmeta.Authors[0].Name
		}
		log.Printf("Lark %-10s Copyright (C) 2016 %s", larkmeta.Version, authors)
		log.Println(lua.PackageCopyRight)
		log.Println()
		log.Printf("This environment simulates that of a lark task.")
		log.Printf("For information about any object use the help() \n" +
			"function.")
	}

	return RunREPL(c.Lua)
}
