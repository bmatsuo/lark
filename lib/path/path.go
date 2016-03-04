package path

import (
	"os"
	"path/filepath"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/yuin/gopher-lua"
)

// Module is a gluamodule.Module that loads the path module.
var Module = gluamodule.New("path", Loader,
	doc.Module,
)

// Loader preloads the path module so that it can be required in lua scripts.
func Loader(l *lua.LState) int {
	l.Pop(1) // first argument is the module name

	mod := l.NewTable()
	doc.Go(l, mod, &doc.GoDocs{
		Desc: `
		The path module provides utilities for working with filesystem paths.
		`,
	})
	l.SetFuncs(mod, Exports)

	glob := l.NewClosure(LuaGlob)
	doc.Go(l, glob, &doc.GoDocs{
		Sig:  "patt => [string]",
		Desc: "Returns an array of paths that match the given pattern.",
		Params: []string{
			"patt  Pattern using star '*' as a wildcard.",
		},
	})
	l.SetField(mod, "glob", glob)

	base := l.NewClosure(LuaBase)
	doc.Go(l, base, &doc.GoDocs{
		Sig:  "path => string",
		Desc: "Returns the basename of the given path.",
		Params: []string{
			"path  A file path that may not exist",
		},
	})
	l.SetField(mod, "base", base)

	dir := l.NewClosure(LuaDir)
	doc.Go(l, dir, &doc.GoDocs{
		Sig:  "path => string",
		Desc: "Returns the directory containing the given path.",
		Params: []string{
			"path  A file path that may not exist",
		},
	})
	l.SetField(mod, "dir", dir)

	ext := l.NewClosure(LuaExt)
	doc.Go(l, ext, &doc.GoDocs{
		Sig:  "path => string",
		Desc: "Returns the file extension of the given path.",
		Params: []string{
			"path  A file path that may not exist",
		},
	})
	l.SetField(mod, "ext", ext)

	join := l.NewClosure(LuaJoin)
	doc.Go(l, join, &doc.GoDocs{
		Sig:  "[path] => string",
		Desc: "Joins the given paths using the filesystem path separator and returns the result.",
		Params: []string{
			"path  A file path that may not exist",
		},
	})
	l.SetField(mod, "join", join)

	exists := l.NewClosure(LuaExists)
	doc.Go(l, exists, &doc.GoDocs{
		Sig:  "path => bool",
		Desc: "Returns true if and only if path exists",
		Params: []string{
			"path  A file path that may not exist",
		},
	})
	l.SetField(mod, "exists", exists)

	isDir := l.NewClosure(LuaIsDir)
	doc.Go(l, isDir, &doc.GoDocs{
		Sig:  "path => bool",
		Desc: "Returns true if and only if path exists and is a directory",
		Params: []string{
			"path  A file path that may not exist",
		},
	})
	l.SetField(mod, "is_dir", isDir)

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
