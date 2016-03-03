package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/bmatsuo/lark/internal/module"
	"github.com/bmatsuo/lark/luamodules/decorator"
	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/bmatsuo/lark/luamodules/lark"
	"github.com/bmatsuo/lark/luamodules/lark/core"
	"github.com/bmatsuo/lark/luamodules/lark/task"
	"github.com/bmatsuo/lark/luamodules/path"
	"github.com/yuin/gopher-lua"
)

// PreloadModules defines the (ordered) set of modules to preload and their
// loader functions.
var PreloadModules = []module.Module{
	task.Module,
	decorator.Module,
	doc.Module,
	path.Module,
	core.Module,
	lark.Module,
}

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
	for _, mod := range PreloadModules {
		module.Preload(c.Lua, mod)
	}

	trace := c.Lua.NewFunction(errTraceback)

	c.Lua.Push(c.Lua.GetGlobal("require"))
	c.Lua.Push(lua.LString("lark"))
	err := c.Lua.PCall(1, 1, trace)
	if err != nil {
		return err
	}
	lark := c.Lua.Get(-1)
	c.Lua.Pop(1)
	c.Lua.SetGlobal("lark", lark)

	c.Lua.SetField(lark, "verbose", lua.LBool(c.Verbose()))

	// Load files after lark has been loaded with require().
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
