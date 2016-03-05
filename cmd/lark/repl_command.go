package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bmatsuo/lark/larkmeta"
	"github.com/bmatsuo/lark/lib"
	"github.com/chzyer/readline"
	"github.com/codegangsta/cli"
	"github.com/yuin/gopher-lua"
)

// CommandREPL implements the "repl" action and launches an interactive
// interpreter.
var CommandREPL = Command(func(lark *Context, cmd *cli.Command) {
	cmd.Name = "repl"
	cmd.Usage = "Run an interactive intepreter"
	cmd.Action = lark.Action(REPL)
})

// REPL loads a lua vm and runs an interactive read, evaluate, print loop.
func REPL(c *Context) {
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

	c.Lua.Push(c.Lua.GetGlobal("require"))
	c.Lua.Push(lua.LString("doc"))
	err = c.Lua.PCall(1, 1, nil)
	if err != nil {
		log.Fatal(err)
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

	err = RunREPL(c.Lua)
	if err != nil {
		log.Fatal(err)
	}
}

// REPLHelp returns the help text for the lark REPL.
func REPLHelp() string {
	var modules []string
	for _, m := range lib.Modules {
		modules = append(modules, m.Name())
	}

	internal := map[string]bool{}
	for _, m := range lib.InternalModules {
		internal[m.Name()] = true
	}

	sort.Strings(modules)
	var help bytes.Buffer
	help.WriteString(REPLHelpDefault)
	help.WriteString("\n\nBuiltin Modules\n\n")
	for _, m := range modules {
		if internal[m] {
			continue
		}
		help.WriteString("\t")
		help.WriteString(m)
		if m == "lark" {
			help.WriteString(" (stored in global ``lark'')")
		}
		help.WriteString("\n")
	}
	return help.String()
}

// REPLHelpDefault describes the REPL environment and points the user at the
// available modules for additional help.
const REPLHelpDefault = `
The REPL environment mimics a lark task runtime with the additional global
function help().

	> help(lark)

Data in the REPL can be stored in global variables.  Local variables will be
out of scope on subsequent lines of input.

	> x = 1
	> return x
	1
	> local y = 2
	> return y
	nil

In order to print a value computed in the REPL it must be returned or otherwise
printed explicitly using the print() builtin.

	> function x() return "xyz" end
	> x()
	> return x()
	xyz
	> print(x())
	xyz
	>

Builtin modules, and custom modules in the lark_modules/ directory can be
imported using the require() function.  Modules can the be used or the
documentation may be inspected

	> path = require('path')
	> help(path)
	> help(path.exists)


Available Modules

	doc
	lark (required by default)
	lark.core
	path
`

// RunREPL runs the main loop of the REPL command.
func RunREPL(state *lua.LState) error {
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for {
		state.SetTop(0)
		content := ""
		prompt := "> "

	read:
		rl.SetPrompt(prompt)
		line, err := rl.Readline()
		if err != nil {
			break // io.EOF
		}
		if content == "" {
			content = line
		} else {
			content += "\n" + line
		}

		fn, err := state.Load(strings.NewReader(content), "stdin")
		if err != nil {
			switch err := err.(type) {
			case *lua.ApiError:
				if err.Type == lua.ApiErrorSyntax {
					if strings.HasPrefix(string(err.Object.(lua.LString)), "stdin at EOF:") {
						prompt = ">> "
						goto read
					}
				}
				fmt.Fprintln(os.Stderr, err)
			case *lua.CompileError:
				fmt.Fprintf(os.Stderr, "compile error: %v\n", err)
			default:
				fmt.Fprintf(os.Stderr, "%T: %v\n", err, err)
			}
			continue
		}

		state.Push(fn)
		err = state.PCall(0, lua.MultRet, nil)
		if err != nil {
			log.Print(err)
			continue
		}

		nret := state.GetTop()
		for i := 1; i <= nret; i++ {
			retval := state.Get(i)
			fmt.Println(retval)
		}
	}

	return nil
}
