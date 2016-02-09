package docs

import (
	"github.com/bmatsuo/lark/internal/module"
	"github.com/yuin/gopher-lua"
)

// Module returns an instance of a Lua module.
func Module() module.Module {
	return defaultDocs.Module()
}

var defaultDocs = &Docs{}

// Docs creates lua modules that provides the docs API.
type Docs struct {
}

// Module returns a lua module.
func (d *Docs) Module() module.Module {
	return &docs{}
}

type docs struct {
	desc   *lua.Table
	params *lua.Table
}

// Loader implements module.Module.
func (d *docs) Loader(l *lua.LState) int {
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
}

func (d *docs) Exports(l *lua.LState, mt *lua.LTable) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"help":      d.LuaHelp,
		"signature": d.LuaSig,
		"param":     d.LuaParam,
	}
}

func (d *docs) LuaSig(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		l.SetField(mt)
	}
}

func (d *docs) LuaParam(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		l.SetField(mt)
	}
}

func (d *docs) LuaHelp(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		l.SetField(mt)
	}
}
