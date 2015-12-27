package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/yuin/gopher-lua"
)

//go:generate ./build_lib.sh

// FindTaskFiles locates task scripts in the project dir.
func FindTaskFiles(dir string) ([]string, error) {
	var luaFiles []string
	join := filepath.Join
	files, err := filepath.Glob(join(dir, "lark.lua"))
	if err != nil {
		return nil, fmt.Errorf("lark.lua: %v", err)
	}
	luaFiles = append(luaFiles, files...)
	files, err = filepath.Glob(join(dir, "lark_tasks/*.lua"))
	if err != nil {
		return nil, fmt.Errorf("lark_tasks: %v", err)
	}
	luaFiles = append(luaFiles, files...)
	return luaFiles, nil

}

// LuaConfig contains options for a new Lua virtual machine.
type LuaConfig struct {
}

// LoadVM creates a lua.State from conf and returns it.
func LoadVM(conf *LuaConfig) (s *lua.LState, err error) {
	s = lua.NewState()
	defer func() {
		if err != nil {
			log.Print(err)
			s.Close()
		}
	}()

	return s, nil
}

// InitLark initializes the lark library and loads files.
func InitLark(c *Context, files []string) error {
	c.Lua.PreloadModule("path", PathLibLoader)
	c.Lua.PreloadModule("lark.core", LarkCoreLoader)

	err := LoadLarkLib(c)
	if err != nil {
		return err
	}

	// This needs to come after LoadLarkLib but can't because the primary
	// library is a file.
	if c.Verbose() && len(files) > 0 {
		log.Printf("loading files: %v", files)
	}
	for _, file := range files {
		err := c.Lua.DoFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadFiles loads the given files into state
func LoadFiles(state *lua.LState, files []string) error {
	for _, file := range files {
		err := state.DoFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// LarkCoreLoader loads the lark.core module.
func LarkCoreLoader(l *lua.LState) int {
	t := l.NewTable()
	mod := l.SetFuncs(t, LarkCoreExports)
	l.Push(mod)
	return 1
}

// LarkCoreExports contains the API for the lark.core lua module.
var LarkCoreExports = map[string]lua.LGFunction{
	"log":  LuaLog,
	"exec": LuaExecRaw,
}

// LoadLarkLib loads the default lark module.
func LoadLarkLib(c *Context) error {
	err := c.Lua.DoString(LarkLib)
	if err != nil {
		return nil
	}

	lark := c.Lua.GetGlobal("lark")

	c.Lua.SetField(lark, "log", c.Lua.NewFunction(LuaLog))
	c.Lua.SetField(lark, "exec_raw", c.Lua.NewFunction(LuaExecRaw))
	c.Lua.SetField(lark, "verbose", lua.LBool(c.Verbose()))

	return nil
}

// LuaLog logs a message from lua.
func LuaLog(state *lua.LState) int {
	opt := &logOpt{}
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

	logLark(msg, opt)

	return 0
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
func LuaExecRaw(state *lua.LState) int {
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

	opt := &execRawOpt{}

	result := execRawLark(args[0], args[1:], opt)
	rt := state.NewTable()
	if result.Err != nil {
		state.SetField(rt, "error", lua.LString(result.Err.Error()))
	}
	state.Push(rt)

	return 1
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

type execRawResult struct {
	Err error
	// Output string
}

type execRawOpt struct {
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

func execRawLark(name string, args []string, opt *execRawOpt) *execRawResult {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if opt == nil {
		err := cmd.Run()
		return &execRawResult{Err: err}
	}

	cmd.Env = opt.Env
	cmd.Dir = opt.Dir

	if opt.StdinFile != "" {
		f, err := os.Open(opt.StdinFile)
		if err != nil {
			return &execRawResult{Err: err}
		}
		defer f.Close()
		cmd.Stdin = f
	}

	if opt.StdoutFile != "" {
		f, err := getOutFile(opt.StdoutFile, opt.StdoutAppend)
		if err != nil {
			return &execRawResult{Err: err}
		}
		defer f.Close()
		cmd.Stdout = f
	}

	if opt.StderrFile != "" && opt.StderrFile == opt.StdoutFile {
		cmd.Stderr = cmd.Stdout
	} else if opt.StderrFile != "" {
		f, err := getOutFile(opt.StderrFile, opt.StderrAppend)
		if err != nil {
			return &execRawResult{Err: err}
		}
		defer f.Close()
		cmd.Stderr = f
	}

	result := &execRawResult{}
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

type logOpt struct {
	Color string
}

func logLark(msg string, opt *logOpt) {
	if opt == nil {
		opt = &logOpt{}
	}

	var esc func(format string, v ...interface{}) string
	if opt.Color != "" {
		esc = colorMap[opt.Color]
	}

	if esc != nil {
		msg = esc("%s", msg)
	}
	log.Print(msg)
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

// PathLibLoader defines the path module so that it can be required.
func PathLibLoader(l *lua.LState) int {
	t := l.NewTable()
	mod := l.SetFuncs(t, PathLibExports)
	l.Push(mod)
	return 1
}

// PathLibExports defines the exported functions in the path module.
var PathLibExports = map[string]lua.LGFunction{
	"glob": LuaGlob,
	"base": LuaBase,
	"dir":  LuaDir,
	"ext":  LuaDir,
	"join": LuaJoin,
}

// LuaGlob executes a file glob.
func LuaGlob(state *lua.LState) int {
	pattern := state.CheckString(1)

	files, err := filepath.Glob(pattern)
	if err != nil {
		state.RaiseError("%s", err.Error())
		return 0
	}
	t := state.NewTable()
	for i, file := range files {
		state.SetTable(t, lua.LNumber(i+1), lua.LString(file))
	}
	state.Push(t)
	return 1
}

// LuaBase returns the basename of the path arguent provided.
func LuaBase(state *lua.LState) int {
	path := state.CheckString(1)

	base := filepath.Base(path)
	state.Push(lua.LString(base))
	return 1
}

// LuaDir returns the parent directory of the path arguent provided.
func LuaDir(state *lua.LState) int {
	path := state.CheckString(1)

	dir := filepath.Dir(path)
	state.Push(lua.LString(dir))
	return 1
}

// LuaExt returns the file extension of the path arguent provided.
func LuaExt(state *lua.LState) int {
	path := state.CheckString(1)

	ext := filepath.Ext(path)
	state.Push(lua.LString(ext))

	return 1
}

// LuaJoin joins the provided path segments.
func LuaJoin(state *lua.LState) int {
	var segs []string

	n := state.GetTop()
	for i := 1; i <= n; i++ {
		str := state.CheckString(i)
		segs = append(segs, str)
	}
	path := filepath.Join(segs...)

	state.Push(lua.LString(path))

	return 1
}
