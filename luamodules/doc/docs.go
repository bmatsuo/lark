//go:generate ./build_lib.sh

package doc

import (
	"github.com/bmatsuo/lark/internal/module"
	"github.com/yuin/gopher-lua"
)

// Module returns an instance of a Lua module.
func Module() module.Module {
	return defaultDocs.Module()
}

var defaultDocs = &Doc{}

// Doc creates lua modules that provides the doc API.
type Doc struct {
}

// Module returns a lua module.
func (d *Doc) Module() module.Module {
	return &doc{}
}

type doc struct {
	desc   *lua.LTable
	params *lua.LTable
}

// Loader implements module.Module.
func (d *doc) Loader(l *lua.LState) int {
	err := l.DoString(DocLib)
	if err != nil {
		l.RaiseError("%s", err)
		return 0
	}
	return 1
	/*
		mt := l.NewTable()
		l.SetField(mt, "__mode", "kv")

		d.desc = l.NewTable()
		l.SetMetatable(d.desc, mt)

		d.params = l.NewTable()
		l.SetMetatable(d.params, mt)

		module := l.NewTable()
		l.SetFuncs(module, d.Exports(mt))
		l.Push(module)
		return 1
	*/
}

func (d *doc) Exports(l *lua.LState, mt *lua.LTable) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"help":      d.LuaHelp(mt),
		"signature": d.LuaSig(mt),
		"param":     d.LuaParam(mt),
	}
}

func (d *doc) LuaSig(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		return 0
	}
}

func (d *doc) LuaParam(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		return 0
	}
}

func (d *doc) LuaHelp(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		return 0
	}
}
