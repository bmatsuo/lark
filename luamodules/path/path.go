package path

import (
	"os"
	"path/filepath"

	"github.com/yuin/gopher-lua"
)

// ModuleLoader defines the path module so that it can be required.
func ModuleLoader(l *lua.LState) int {
	t := l.NewTable()
	mod := l.SetFuncs(t, Exports)
	l.Push(mod)
	return 1
}

// Exports defines the exported functions in the path module.
var Exports = map[string]lua.LGFunction{
	"glob":   LuaGlob,
	"base":   LuaBase,
	"dir":    LuaDir,
	"ext":    LuaExt,
	"join":   LuaJoin,
	"exists": LuaExists,
	"is_dir": LuaIsDir,
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

// LuaExists returns true if the provided path segment exists.
func LuaExists(state *lua.LState) int {
	path := state.CheckString(1)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		state.Push(lua.LBool(false))
	} else if err == nil {
		state.Push(lua.LBool(true))
	} else {
		state.RaiseError("%s", err.Error())
		return 0
	}
	return 1
}

// LuaIsDir returns true if the provided path segment exists and is a
// directory.
func LuaIsDir(state *lua.LState) int {
	path := state.CheckString(1)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		state.Push(lua.LBool(false))
		return 1
	}
	if err != nil {
		state.RaiseError("%s", err.Error())
		return 0
	}
	state.Push(lua.LBool(info.IsDir()))
	return 1
}
