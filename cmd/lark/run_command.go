package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/bmatsuo/lark/lib/lark/core"
	"github.com/codegangsta/cli"
	"github.com/yuin/gopher-lua"
)

// CommandRun implements the "run" action (the default action)
var CommandRun = Command(func(lark *Context, cmd *cli.Command) {
	cmd.Name = "run"
	cmd.Aliases = []string{"make"}
	cmd.Usage = "Run lark project task(s)"
	cmd.ArgsUsage = `task ...

    The arguments are the names of tasks from lark.lua.`
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "C",
			Usage:  "Change the working directory before loading files and running tasks",
			EnvVar: "LARK_RUN_DIRECTORY",
		},
		cli.IntFlag{
			Name:   "j",
			Usage:  "Number of parallel processes.",
			EnvVar: "LARK_RUN_PARALLEL",
		},
		cli.BoolFlag{
			Name:        "v",
			Usage:       "Enable verbose reporting of errors.",
			EnvVar:      "LARK_VERBOSE",
			Destination: Verbose,
		},
	}
	cmd.Action = lark.Action(Run)
})

// Run loads a lua vm and runs tasks specified in the command line.
func Run(c *Context) {
	chdir := c.String("C")
	if chdir != "" {
		err := os.Chdir(chdir)
		if err != nil {
			log.Fatal(err)
		}
	}

	args := c.Args()
	var tasks []*Task
	for {
		t, n, err := ParseTask(args)
		if err != nil {
			log.Fatalf("task %d: %v", len(tasks), err)
		}
		if n == 0 {
			break
		}
		tasks = append(tasks, t)
		args = args[n:]
	}
	if len(tasks) == 0 {
		tasks = []*Task{{}}
	}

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

	for _, task := range tasks {
		err := RunTask(c, task)
		if err != nil {
			os.Exit(1)
		}
	}
}

func normTasks(args []string) ([]string, error) {
	if len(args) == 0 {
		return []string{""}, nil
	}
	for _, task := range args {
		if task == "" {
			return nil, fmt.Errorf("invalid task name")
		}
	}
	return args, nil
}

// RunTask calls lark.run in state to execute task.
func RunTask(c *Context, task *Task) error {
	lark := c.Lua.GetGlobal("lark")
	run := c.Lua.GetField(lark, "run")
	trace := c.Lua.NewFunction(errTraceback)

	narg := 1
	c.Lua.Push(run)
	if task.Name == "" {
		c.Lua.Push(lua.LNil)
	} else {
		c.Lua.Push(lua.LString(task.Name))
	}
	if len(task.Params) > 0 {
		params := c.Lua.NewTable()
		for k, v := range task.Params {
			c.Lua.SetField(params, k, lua.LString(v))
		}

		c.Lua.Push(params)
		narg++
	}
	err := c.Lua.PCall(narg, 0, trace)
	if err != nil {
		handleErr(c, err)
	}

	wait := c.Lua.GetField(lark, "wait")
	for {
		c.Lua.Push(wait)
		errwait := c.Lua.PCall(0, 0, trace)
		if errwait == nil {
			break
		}
		if err == nil {
			handleErr(c, errwait)

			// prevent handleErr from being called multiple times.
			err = errwait
		}
	}

	return err
}

func handleErr(c *Context, err error) {
	core.Log(fmt.Sprint(err), &core.LogOpt{
		Color: "red",
	})
}

// Task is a task invocation from the command line.
type Task struct {
	Name   string
	Params map[string]string
}

// ToLua returns a string representing the task in lua table syntax.
func (t *Task) ToLua() string {
	var name string
	if t.Name == "" {
		name = "nil"
	} else {
		name = fmt.Sprintf("%q", t.Name)
	}
	return fmt.Sprintf("{name=%s,params=%s}", name, luamap(t.Params))
}

func luamap(m map[string]string) string {
	buf := bytes.NewBuffer(nil)
	io.WriteString(buf, "{")
	for k, v := range m {
		io.WriteString(buf, k)
		io.WriteString(buf, "=")
		fmt.Fprintf(buf, "%q", v)
	}
	io.WriteString(buf, "}")
	return buf.String()
}

// ParseTask parses a task from command line arguments and returns it along
// with the number of args consumed.
func ParseTask(args []string) (*Task, int, error) {
	t := &Task{}
	if len(args) == 0 {
		return t, 0, nil
	}
	if args[0] == "--" {
		return t, 1, nil
	}
	if !strings.Contains(args[0], "=") {
		t.Name = args[0]
		args = args[1:]
		if len(t.Name) == 0 {
			return nil, 0, fmt.Errorf("missing name")
		}
	}
	t.Params = make(map[string]string)
	for _, p := range args {
		pieces := strings.SplitN(p, "=", 2)
		if len(pieces) == 1 {
			break
		}
		err := sanitizeParam(pieces[0])
		if err != nil {
			return nil, 0, err
		}
		t.Params[pieces[0]] = pieces[1]
	}
	return t, 1 + len(t.Params), nil
}

func sanitizeParam(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("missing param name")
	}
	badchars := strings.TrimFunc(name, func(c rune) bool {
		return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_'
	})
	if len(badchars) > 0 {
		return fmt.Errorf("invalid character in param: %q", badchars[0])
	}
	return nil
}

var reLoc = regexp.MustCompile(`^[^:]+:\d+:\s*`)

func trimLoc(msg string) string {
	return reLoc.ReplaceAllString(msg, "")
}

func errTraceback(L *lua.LState) int {
	msg := L.Get(1)
	L.SetTop(0)
	L.Push(L.GetField(L.GetGlobal("debug"), "traceback"))
	L.Push(msg)
	L.Push(lua.LNumber(2))
	L.Call(2, 1)
	return 1
}
