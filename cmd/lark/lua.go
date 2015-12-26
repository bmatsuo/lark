package main

import (
	"fmt"
	"log"
	"path/filepath"

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
	err := LoadLarkLib(c)
	if err != nil {
		return err
	}

	// This needs to come after LoadLarkLib but can't because the primary
	// library is a file.
	if c.Verbose && len(files) > 0 {
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

// LoadLarkLib loads the default lark module.
func LoadLarkLib(c *Context) error {
	err := c.Lua.DoString(LarkLib)
	if err != nil {
		return nil
	}

	lark := c.Lua.GetGlobal("lark")
	c.Lua.SetField(lark, "log", c.Lua.NewFunction(LuaLog))
	c.Lua.SetField(lark, "verbose", lua.LBool(c.Verbose))
	return nil
}

// LuaLog logs a message from lua.
func LuaLog(state *lua.LState) int {
	opt := &logOpt{}
	var msg string
	v1 := state.Get(1)
	if v1.Type() == lua.LTTable {
		arr := luaTableArray(state, v1.(*lua.LTable))
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

func luaTableArray(state *lua.LState, t *lua.LTable) []lua.LValue {
	var vals []lua.LValue
	t.ForEach(func(kv, vv lua.LValue) {
		if kv.Type() == lua.LTNumber {
			vals = append(vals, vv)
		}
	})
	return vals
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
