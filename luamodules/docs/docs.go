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

type docs struct{}

// Loader implements module.Module.
func (d *docs) Loader(l *lua.LState) int {
	concat := l.NewFunction(mtconcat)
	mt := l.NewTable()
	l.SetField(mt, "__concat", concat)

	module := l.NewTable()
	l.SetFuncs(module, d.Exports(mt))
	l.Push(module)
	return 1
}

func mtconcat(*lua.LState) int {

}

func (d *docs) Exports(l *lua.LState, mt *lua.LTable) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"name":    d.LuaName,
		"param":   d.LuaParam,
		"returns": d.LuaReturns,
	}
}

func (d *docs) LuaName(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		l.SetField(mt)
	}
}
