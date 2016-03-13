package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/bmatsuo/lark/execgroup"
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/yuin/gopher-lua"
)

// Module is a gluamodule.Module that loads the lark.core module.
var Module = gluamodule.New("lark.core", Loader,
	doc.Module,
)

// InitModule changes the configuration of the module.  It is not safe to call
// InitModule after the module has been loaded.
func InitModule(logWriter io.Writer, limit int) {
	if logWriter == nil {
		logWriter = os.Stderr
	}
	if limit == 0 {
		limit = runtime.NumCPU()
	}
	defaultCore = newCore(logWriter, limit)
}

var defaultCore = newCore(os.Stderr, runtime.NumCPU())

type core struct {
	logger     *log.Logger
	isTTY      bool
	groups     map[string]*execgroup.Group
	limit      chan struct{}
	grouplimit map[string]chan struct{}
}

func istty(w io.Writer) bool {
	type fd interface {
		Fd() uintptr
	}
	switch w := w.(type) {
	case fd:
		return isatty.IsTerminal(w.Fd())
	default:
		return false
	}
}

func newCore(logfile io.Writer, limit int) *core {
	c := &core{
		isTTY:  istty(logfile),
		groups: make(map[string]*execgroup.Group),
	}
	if limit > 0 {
		c.limit = make(chan struct{}, limit)
	}

	logFlags := log.LstdFlags
	if c.isTTY {
		logFlags &^= log.Ldate | log.Ltime
	}
	c.logger = log.New(logfile, "", logFlags)

	return c
}

// Loader preloads the lark.core module so it may be required in lua scripts.
func Loader(l *lua.LState) int {
	t := l.NewTable()
	mod := l.SetFuncs(t, defaultCore.exports())
	l.Push(mod)
	return 1
}

// exports contains the API for the lark.core lua module.
func (c *core) exports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"log":        c.LuaLog,
		"environ":    c.LuaEnviron,
		"exec":       c.LuaExecRaw,
		"start":      c.LuaStartRaw,
		"make_group": c.LuaMakeGroup,
		"wait":       c.LuaWait,
	}
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

// LuaMakeGroup makes creates a group with dependencies.  LuaMakeGroup expects
// one table argument.
func (c *core) LuaMakeGroup(state *lua.LState) int {
	v1 := state.Get(1)
	if v1.Type() != lua.LTTable {
		state.ArgError(1, "first argument must be a table")
		return 0
	}

	var groupname string
	var lname = state.GetField(v1, "name")
	if lname == lua.LNil {
		state.ArgError(1, "missing named value 'name'")
		return 0
	}
	if lname.Type() != lua.LTString {
		msg := fmt.Sprintf("named value 'name' is not a string: %s", lname.Type())
		state.ArgError(1, msg)
		return 0
	}
	groupname = string(lname.(lua.LString))

	var limit int
	llimit := state.GetField(v1, "limit")
	if llimit != lua.LNil {
		_limit, ok := llimit.(lua.LNumber)
		if !ok {
			msg := fmt.Sprintf("named value 'limit' is not a number: %s", llimit.Type())
			state.ArgError(1, msg)
			return 0
		}
		limit = int(_limit)
	}

	var follows []string
	lfollows := state.GetField(v1, "follows")
	if lfollows != lua.LNil {
		switch val := lfollows.(type) {
		case lua.LString:
			follows = append(follows, string(val))
		case *lua.LTable:
			tvals := flattenTable(state, val)
			for _, tv := range tvals {
				s, ok := tv.(lua.LString)
				if !ok {
					msg := fmt.Sprintf("named value 'follows' may only contain strings: %s", tv.Type())
					state.ArgError(1, msg)
					return 0
				}
				follows = append(follows, string(s))
			}
		default:
			msg := fmt.Sprintf("named value 'follows' is not a table: %s", lfollows.Type())
			state.ArgError(1, msg)
			return 0
		}
	}

	var gfollows []*execgroup.Group
	for _, name := range follows {
		g, ok := c.groups[name]
		if !ok {
			g = execgroup.NewGroup(nil)
			c.groups[name] = g
		}
		gfollows = append(gfollows, g)
	}

	_, ok := c.groups[groupname]
	if ok {
		msg := fmt.Sprintf("group already exists: %q", groupname)
		state.ArgError(1, msg)
		return 0
	}

	c.groups[groupname] = execgroup.NewGroup(gfollows)
	if limit < 0 {
		c.grouplimit[groupname] = nil
	} else if limit > 0 {
		c.grouplimit[groupname] = make(chan struct{}, limit)
	}

	return 0
}

func (c *core) LuaWait(state *lua.LState) int {
	var names []string
	n := state.GetTop()
	if n == 0 {
		for name := range c.groups {
			names = append(names, name)
		}
	} else {
		for i := 1; i <= n; i++ {
			names = append(names, state.CheckString(1))
		}
	}

	rt := state.NewTable()

	var err error
	for _, name := range names {
		group := c.groups[name]
		if group != nil {
			if err == nil {
				err = group.Wait()
			} else {
				group.Wait()
			}
		}
	}

	if err != nil {
		msg := fmt.Sprintf("asynchronous error: %v", err)
		state.SetField(rt, "error", lua.LString(msg))
	}
	state.Push(rt)

	return 1
}

// LuaStartRaw makes executes a program.  LuaStartRaw expects one table
// argument and returns one table.
func (c *core) LuaStartRaw(state *lua.LState) int {
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

	ignore := false
	lignore := state.GetField(v1, "ignore")
	if lignore != lua.LNil {
		_ignore, ok := lignore.(lua.LBool)
		if !ok {
			msg := fmt.Sprintf("named value 'ignore' is not boolean: %s", lignore.Type())
			state.ArgError(1, msg)
			return 0
		}
		ignore = bool(_ignore)
	}

	var groupname string
	lgroup := state.GetField(v1, "group")
	if lgroup != lua.LNil {
		group, ok := lgroup.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'group' is not a string: %s", lgroup.Type())
			state.ArgError(1, msg)
			return 0
		}
		groupname = string(group)
	}

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

	lstdin := state.GetField(v1, "stdin")
	if lstdin != lua.LNil {
		stdin, ok := lstdin.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'stdin' is not a string: %s", lstdin.Type())
			state.ArgError(1, msg)
			return 0
		}
		opt.StdinFile = string(stdin)
	}

	linput := state.GetField(v1, "input")
	if linput != lua.LNil {
		input, ok := linput.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'input' is not a string: %s", linput.Type())
			state.ArgError(1, msg)
			return 0
		}
		if opt.StdinFile != "" {
			msg := fmt.Sprintf("conflicting named values 'stdin' and 'input' both provided")
			state.ArgError(1, msg)
			return 0
		}
		opt.Input = []byte(input)
	}

	lstdout := state.GetField(v1, "stdout")
	if lstdout != lua.LNil {
		stdout, ok := lstdout.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'stdout' is not a string: %s", lstdout.Type())
			state.ArgError(1, msg)
			return 0
		}
	stdoutsigloop:
		for {
			switch {
			case strings.HasPrefix(string(stdout), "+"):
				opt.StdoutAppend = true
				stdout = stdout[1:]
			case strings.HasPrefix(string(stdout), "&"):
				opt.StdoutTee = true
				stdout = stdout[1:]
			case strings.HasPrefix(string(stdout), "$"):
				state.RaiseError("output capture not allowed for 'start'")
			default:
				break stdoutsigloop
			}
		}
		opt.StdoutFile = string(stdout)
		if !opt.StdoutCapture && stdout == "" {
			opt.StdoutTee = false
		}
	}

	lstderr := state.GetField(v1, "stderr")
	if lstderr != lua.LNil {
		stderr, ok := lstderr.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'stderr' is not a string: %s", lstderr.Type())
			state.ArgError(1, msg)
			return 0
		}
	stderrsigloop:
		for {
			switch {
			case strings.HasPrefix(string(stderr), "+"):
				opt.StderrAppend = true
				stderr = stderr[1:]
			case strings.HasPrefix(string(stderr), "&"):
				opt.StderrTee = true
				stderr = stderr[1:]
			case strings.HasPrefix(string(stderr), "$"):
				state.RaiseError("output capture not allowed for 'start'")
			default:
				break stderrsigloop
			}
		}
		opt.StderrFile = string(stderr)
		if !opt.StderrCapture && stderr == "" {
			opt.StderrTee = false
		}
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

	lstr := state.GetField(v1, "_str")
	str, _ := lstr.(lua.LString)
	lecho, ok := state.GetField(v1, "echo").(lua.LBool)
	if !ok {
		lecho = true
	}

	group, ok := c.groups[groupname]
	if !ok {
		group = execgroup.NewGroup(nil)
		c.groups[groupname] = group
	}

	limit := c.limit
	glimit, ok := c.grouplimit[groupname]
	if ok && glimit == nil {
		// if the group has specifically been unilimited then remove the global
		// limit as well.
		limit = nil
	}
	err := group.Exec(func() error {
		if glimit != nil {
			glimit <- struct{}{}
			defer func() { <-glimit }()
		}
		if limit != nil {
			limit <- struct{}{}
			defer func() { <-limit }()
		}
		if str != "" && lecho {
			opt := &LogOpt{Color: "green"}
			c.log(string(str), opt)
		}
		result := c.execRaw(args[0], args[1:], opt)
		if ignore {
			return nil
		}
		return result.Err
	})

	rt := state.NewTable()

	if err != nil {
		msg := fmt.Sprintf("asynchronous error: %v", err)
		state.SetField(rt, "error", lua.LString(msg))
	}
	state.Push(rt)

	return 1
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

	lstdin := state.GetField(v1, "stdin")
	if lstdin != lua.LNil {
		stdin, ok := lstdin.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'stdin' is not a string: %s", lstdin.Type())
			state.ArgError(1, msg)
			return 0
		}
		opt.StdinFile = string(stdin)
	}

	linput := state.GetField(v1, "input")
	if linput != lua.LNil {
		input, ok := linput.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'input' is not a string: %s", linput.Type())
			state.ArgError(1, msg)
			return 0
		}
		if opt.StdinFile != "" {
			msg := fmt.Sprintf("conflicting named values 'stdin' and 'input' both provided")
			state.ArgError(1, msg)
			return 0
		}
		opt.Input = []byte(input)
	}

	lstdout := state.GetField(v1, "stdout")
	if lstdout != lua.LNil {
		stdout, ok := lstdout.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'stdout' is not a string: %s", lstdout.Type())
			state.ArgError(1, msg)
			return 0
		}
	stdoutsigloop:
		for {
			switch {
			case strings.HasPrefix(string(stdout), "+"):
				opt.StdoutAppend = true
				stdout = stdout[1:]
			case strings.HasPrefix(string(stdout), "&"):
				opt.StdoutTee = true
				stdout = stdout[1:]
			case strings.HasPrefix(string(stdout), "$"):
				opt.StdoutCapture = true
				stdout = stdout[1:]
			default:
				break stdoutsigloop
			}
		}
		opt.StdoutFile = string(stdout)
		if !opt.StdoutCapture && stdout == "" {
			opt.StdoutTee = false
		}
	}

	lstderr := state.GetField(v1, "stderr")
	if lstderr != lua.LNil {
		stderr, ok := lstderr.(lua.LString)
		if !ok {
			msg := fmt.Sprintf("named value 'stderr' is not a string: %s", lstderr.Type())
			state.ArgError(1, msg)
			return 0
		}
	stderrsigloop:
		for {
			switch {
			case strings.HasPrefix(string(stderr), "+"):
				opt.StderrAppend = true
				stderr = stderr[1:]
			case strings.HasPrefix(string(stderr), "&"):
				opt.StderrTee = true
				stderr = stderr[1:]
			case strings.HasPrefix(string(stderr), "$"):
				opt.StderrCapture = true
				stderr = stderr[1:]
			default:
				break stderrsigloop
			}
		}
		opt.StderrFile = string(stderr)
		if !opt.StderrCapture && stderr == "" {
			opt.StderrTee = false
		}
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

	lecho, ok := state.GetField(v1, "echo").(lua.LBool)
	if !ok {
		lecho = true
	}

	lstr := state.GetField(v1, "_str")
	str, _ := lstr.(lua.LString)

	if str != "" && lecho {
		opt := &LogOpt{Color: "green"}
		c.log(string(str), opt)
	}
	result := c.execRaw(args[0], args[1:], opt)
	rt := state.NewTable()
	if result.Err != nil {
		state.SetField(rt, "error", lua.LString(result.Err.Error()))
	}
	if opt.StdoutCapture || opt.StderrCapture {
		state.SetField(rt, "output", lua.LString(result.Output))
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
	Err    error
	Output string
}

// ExecRawOpt contains options for ExecRaw.
type ExecRawOpt struct {
	Env []string
	Dir string

	Input        []byte
	StdinFile    string
	StdoutFile   string
	StdoutAppend bool
	StderrFile   string
	// StderrAppend is ignored if StderrFile equals StdoutFile.
	StderrAppend bool

	// TODO: Output capture. This will interact interestingly with command
	// result caching and file redirection.  It should be thought out more.
	StdoutCapture bool
	StderrCapture bool

	StdoutTee bool
	StderrTee bool
}

// ExecRaw executes the named command with the given arguments.
func ExecRaw(name string, args []string, opt *ExecRawOpt) *ExecRawResult {
	return defaultCore.execRaw(name, args, opt)
}

func (c *core) execRaw(name string, args []string, opt *ExecRawOpt) *ExecRawResult {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin

	var buf io.Writer
	if opt.StdoutCapture || opt.StderrCapture {
		buf = &syncBuffer{}
	}
	var stdout io.Writer
	var stderr io.Writer

	if opt == nil {
		err := cmd.Run()
		return &ExecRawResult{Err: err}
	}

	cmd.Env = opt.Env
	cmd.Dir = opt.Dir

	if len(opt.Input) != 0 {
		cmd.Stdin = bytes.NewReader(opt.Input)
	} else if opt.StdinFile != "" {
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
		stdout = f
	} else if opt.StdoutCapture {
		stdout = nil
	}

	if opt.StderrFile != "" && opt.StderrFile == opt.StdoutFile {
		stderr = stdout
	} else if opt.StderrFile != "" {
		f, err := getOutFile(opt.StderrFile, opt.StderrAppend)
		if err != nil {
			return &ExecRawResult{Err: err}
		}
		defer f.Close()
		stderr = f
	} else if opt.StderrCapture {
		stderr = nil
	}

	if opt.StdoutCapture {
		if stdout != nil {
			stdout = io.MultiWriter(buf, stdout)
		} else {
			stdout = buf
		}
	}
	if opt.StdoutTee && stdout != nil {
		stdout = io.MultiWriter(stdout, os.Stdout)
	}
	if opt.StderrCapture {
		if stderr != nil {
			stderr = io.MultiWriter(buf, stderr)
		} else {
			stderr = buf
		}
	}
	if opt.StderrTee && stderr != nil {
		stderr = io.MultiWriter(stderr, os.Stderr)
	}

	ioerr := make(chan error, 2)
	read := func(w io.Writer, r io.Reader, e chan<- error) {
		_, err := io.Copy(w, r)
		e <- err
	}

	var pout, perr io.ReadCloser
	var closers []io.ReadCloser
	doclose := func() {
		for _, f := range closers {
			f.Close()
		}
	}
	if stdout != nil {
		var err error
		pout, err = cmd.StdoutPipe()
		if err != nil {
			return &ExecRawResult{Err: err}
		}
		closers = append(closers, pout)
	} else {
		cmd.Stdout = os.Stdout
	}
	if stderr != nil {
		var err error
		perr, err = cmd.StderrPipe()
		if err != nil {
			doclose()
			return &ExecRawResult{Err: err}
		}
		closers = append(closers, perr)
	} else {
		cmd.Stderr = os.Stderr
	}

	result := &ExecRawResult{}
	result.Err = cmd.Start()
	if result.Err != nil {
		doclose()
		return result
	}

	defer func() {
		if buf != nil {
			result.Output = string(buf.(*syncBuffer).Bytes())
		}
	}()

	n := 0
	if stdout != nil {
		n++
		read(stdout, pout, ioerr)
	}
	if stderr != nil {
		n++
		read(stderr, perr, ioerr)
	}
	for i := 0; i < n; i++ {
		result.Err = <-ioerr
		if result.Err != nil {
			go func(k int) {
				for j := 0; j < n; j++ {
					<-ioerr
				}
				cmd.Wait()
			}(n - i - 1)
			return result
		}
	}

	result.Err = cmd.Wait()

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

type syncBuffer struct {
	mut sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Read(p []byte) (int, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.Read(p)
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) WriteString(s string) (int, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.WriteString(s)
}

func (b *syncBuffer) ReadFrom(r io.Reader) (int64, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.ReadFrom(r)
}

func (b *syncBuffer) WriteTo(w io.Writer) (int64, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.WriteTo(w)
}

func (b *syncBuffer) Bytes() []byte {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.Bytes()
}

func (b *syncBuffer) String() string {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.String()
}
