package core

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/yuin/gopher-lua"
)

var defaultCore = newCore(os.Stderr)

type core struct {
	logger *log.Logger
	isTTY  bool
}

func istty(w io.Writer) bool {
	type fd interface {
		Fd() uintptr
	}
	wfd, ok := w.(fd)
	if ok {
		if isatty.IsTerminal(wfd.Fd()) {
			return true
		}
	}
	return false
}

func newCore(logfile io.Writer) *core {
	c := &core{
		isTTY: istty(logfile),
	}
	flags := log.LstdFlags
	if c.isTTY {
		flags &^= log.Ldate | log.Ltime
	}
	c.logger = log.New(logfile, "", flags)

	return c
}

// Loader loads the lark.core module.
func Loader(l *lua.LState) int {
	t := l.NewTable()
	mod := l.SetFuncs(t, Exports)
	l.Push(mod)
	return 1
}

// Exports contains the API for the lark.core lua module.
var Exports = map[string]lua.LGFunction{
	"log":     defaultCore.LuaLog,
	"environ": defaultCore.LuaEnviron,
	"exec":    defaultCore.LuaExecRaw,
}

// LuaLog logs a message from lua.
func (c *core) LuaLog(state *lua.LState) int {
	opt := &LogOpt{}
	var msg string
	v1 := state.Get(1)
	if v1.Type() == lua.LTTable {
		arr := luaTableArray(state, v1.(*lua.LTable), nil)
		if len(arr) > 0 {
			msg = fmt.Sprint(arr[0])
		}
	} else if v1.Type() == lua.LTString {
		msg = string(string(v1.(lua.LString)))
	}

	lcolor, ok := state.GetField(v1, "color").(lua.LString)
	if ok {
		opt.Color = string(lcolor)
	}

	c.log(msg, opt)

	return 0
}

func (c *core) LuaEnviron(state *lua.LState) int {
	rt := state.NewTable()

	for _, env := range os.Environ() {
		pieces := strings.SplitN(env, "=", 2)
		if len(pieces) == 2 {
			state.SetField(rt, pieces[0], lua.LString(pieces[1]))
		} else {
			state.SetField(rt, pieces[0], lua.LString(""))
		}
	}

	state.Push(rt)

	return 1
}

func luaTableArray(state *lua.LState, t *lua.LTable, vals []lua.LValue) []lua.LValue {
	t.ForEach(func(kv, vv lua.LValue) {
		if kv.Type() == lua.LTNumber {
			vals = append(vals, vv)
		}
	})
	return vals
}

// LuaExecRaw makes executes a program.  LuaExecRaw expects one table argument
// and returns one table.
func (c *core) LuaExecRaw(state *lua.LState) int {
	v1 := state.Get(1)
	if v1.Type() != lua.LTTable {
		state.ArgError(1, "first argument must be a table")
		return 0
	}

	largs := flattenTable(state, v1.(*lua.LTable))
	if len(largs) == 0 {
		state.ArgError(1, "missing positional values")
		return 0
	}
	args := make([]string, len(largs))
	for i, larg := range largs {
		arg, ok := larg.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("positional values are not strings: %s", larg.Type())
			state.ArgError(1, msg)
			return 0
		}
		args[i] = string(arg)
	}

	opt := &ExecRawOpt{}

	ldir := state.GetField(v1, "dir")
	if ldir != lua.LNil {
		dir, ok := ldir.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'dir' is not a string: %s", ldir.Type())
			state.ArgError(1, msg)
			return 0
		}
		opt.Dir = string(dir)
	}

	var env []string
	lenv := state.GetField(v1, "env")
	if lenv != lua.LNil {
		t, ok := lenv.(*lua.LTable)
		if !ok {
			msg := fmt.Sprintf("env is not a table: %s", lenv.Type())
			state.ArgError(1, msg)
			return 0
		}
		var err error
		env, err = tableEnv(t)
		if err != nil {
			state.ArgError(1, err.Error())
			return 0
		}
	}
	opt.Env = env

	result := c.execRaw(args[0], args[1:], opt)
	rt := state.NewTable()
	if result.Err != nil {
		state.SetField(rt, "error", lua.LString(result.Err.Error()))
	}
	state.Push(rt)

	return 1
}

func tableEnv(t *lua.LTable) ([]string, error) {
	var env []string
	msg := ""
	t.ForEach(func(kv, vv lua.LValue) {
		if kv.Type() == lua.LTString && vv.Type() == lua.LTString {
			defn := fmt.Sprintf("%s=%s", kv.(lua.LString), vv.(lua.LString))
			env = append(env, defn)
		} else if msg == "" {
			msg = fmt.Sprintf("invalid %s-%s environment variable pair", kv.Type(), vv.Type())
		}
	})
	if msg != "" {
		return nil, errors.New(msg)
	}
	return env, nil
}

func flattenTable(state *lua.LState, val *lua.LTable) []lua.LValue {
	var flat []lua.LValue
	largs := luaTableArray(state, val, nil)
	for _, arg := range largs {
		switch t := arg.(type) {
		case *lua.LTable:
			flat = append(flat, flattenTable(state, t)...)
		default:
			flat = append(flat, arg)
		}
	}
	return flat
}

// ExecRawResult is returned from ExecRaw
type ExecRawResult struct {
	Err error
	// Output string
}

// ExecRawOpt contains options for ExecRaw.
type ExecRawOpt struct {
	Env []string
	Dir string

	StdinFile    string
	StdoutFile   string
	StdoutAppend bool
	StderrFile   string
	// StderrAppend is ignored if StderrFile equals StdoutFile.
	StderrAppend bool

	// TODO: Output capture. This will interact interestingly with command
	// result caching and file redirection.  It should be thought out more.
	//CaptureOutput  bool
}

// ExecRaw executes the named command with the given arguments.
func ExecRaw(name string, args []string, opt *ExecRawOpt) *ExecRawResult {
	return defaultCore.execRaw(name, args, opt)
}

func (c *core) execRaw(name string, args []string, opt *ExecRawOpt) *ExecRawResult {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if opt == nil {
		err := cmd.Run()
		return &ExecRawResult{Err: err}
	}

	cmd.Env = opt.Env
	cmd.Dir = opt.Dir

	if opt.StdinFile != "" {
		f, err := os.Open(opt.StdinFile)
		if err != nil {
			return &ExecRawResult{Err: err}
		}
		defer f.Close()
		cmd.Stdin = f
	}

	if opt.StdoutFile != "" {
		f, err := getOutFile(opt.StdoutFile, opt.StdoutAppend)
		if err != nil {
			return &ExecRawResult{Err: err}
		}
		defer f.Close()
		cmd.Stdout = f
	}

	if opt.StderrFile != "" && opt.StderrFile == opt.StdoutFile {
		cmd.Stderr = cmd.Stdout
	} else if opt.StderrFile != "" {
		f, err := getOutFile(opt.StderrFile, opt.StderrAppend)
		if err != nil {
			return &ExecRawResult{Err: err}
		}
		defer f.Close()
		cmd.Stderr = f
	}

	result := &ExecRawResult{}
	result.Err = cmd.Run()
	return result
}

func getOutFile(name string, a bool) (*os.File, error) {
	if !strings.HasPrefix(name, "&") {
		flag := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
		if a {
			flag |= os.O_APPEND
		}
		f, err := os.OpenFile(name, flag, 0644)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if name == "&1" {
		return os.Stdout, nil
	}
	if name == "&2" {
		return os.Stderr, nil
	}
	return nil, fmt.Errorf("invalid file descriptor")
}

// LogOpt contains options for the Log function
type LogOpt struct {
	Color string
}

// Log logs a message to standard error.
func Log(msg string, opt *LogOpt) {
	defaultCore.log(msg, opt)
}

func (c *core) log(msg string, opt *LogOpt) {
	if opt == nil {
		opt = &LogOpt{}
	}

	var esc func(format string, v ...interface{}) string
	if opt.Color != "" && c.isTTY {
		esc = colorMap[opt.Color]
	}

	if esc != nil {
		msg = esc("%s", msg)
	}
	c.logger.Print(msg)
}

var colorMap = map[string]func(format string, v ...interface{}) string{
	"black":   color.BlackString,
	"blue":    color.BlueString,
	"cyan":    color.CyanString,
	"green":   color.GreenString,
	"magenta": color.MagentaString,
	"red":     color.RedString,
	"white":   color.WhiteString,
	"yellow":  color.YellowString,
}
