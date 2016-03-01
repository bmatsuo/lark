package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	helpFunc := c.Lua.GetField(docModule, "help")
	c.Lua.SetGlobal("help", helpFunc)

	err = RunREPL(c.Lua)
	if err != nil {
		log.Fatal(err)
	}
}

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
