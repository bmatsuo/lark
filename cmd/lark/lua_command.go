package main

import (
	"io"
	"log"
	"os"

	"github.com/bmatsuo/lark/larkmeta"
	"github.com/codegangsta/cli"
	"github.com/yuin/gopher-lua"
)

// CommandLua implements the "lua" command and runs a lark-flavored Lua
// interpreter.
var CommandLua = Command(func(lark *Context, cmd *cli.Command) {
	cmd.Name = "lua"
	cmd.Usage = "Run an interactive intepreter"
	cmd.Action = lark.Action(REPL)
})

// Lua loads a lua vm with the lark library and executes Lua scripts or
// expressions.
func Lua(c *Context) {
	luaFiles, err := FindTaskFiles("")
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
	args := c.Args()

	// when connected to a terminal and executing with no arguments a special
	// interactive interpreter is launched which has a unique procedure for
	// reading and evaluating input.
	if IsTTY && len(args) == 0 {
		err := LuaInteractive(c)
		if err != nil {
			return err
		}
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
// from the terminal.
func LuaInteractive(c *Context) error {
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

	log.Printf("Lark %-10s Copyright (C) 2016 The Lark authors", larkmeta.Version)
	log.Println(lua.PackageCopyRight)
	log.Println()
	log.Printf("This environment simulates that of a lark task.")
	log.Printf("For information about any object use the help() \n" +
		"function.")

	return RunREPL(c.Lua)
}
