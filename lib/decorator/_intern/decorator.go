package intern

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

// Module is a gluamodule.Module that loads the "decorator" module.
var Module = gluamodule.New("decorator.intern", Loader)

// Loader is a lua.LGFunction that loads the module.
func Loader(l *lua.LState) int {
	mod := l.NewTable()

	alias := l.NewClosure(luaCallAlias)
	metatable := l.NewClosure(luaMetatable(alias), alias)

	l.Push(metatable)
	l.Push(lua.LNil)
	l.Call(1, 1)
	mt := l.Get(-1).(*lua.LTable)
	l.Pop(0)

	create := l.NewClosure(luaCreate(mt), mt)
	l.SetField(mod, "metatable", metatable)
	l.SetField(mod, "create", create)
	l.SetField(mod, "annotator", l.NewClosure(luaAnnotator(create), create))

	l.Push(mod)
	return 1
}

func luaMetatable(alias *lua.LFunction) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.OptFunction(1, nil)
		l.SetTop(0)

		mt := l.NewTable()
		l.SetField(mt, "__concat", alias)
		if fn != nil {
			l.SetField(mt, "__call", fn)
		}
		l.Push(mt)
		return 1
	}
}

func luaCallAlias(l *lua.LState) int {
	l.CheckAny(1)
	l.CheckAny(2)
	l.SetTop(2)
	l.Call(1, lua.MultRet)
	return l.GetTop()
}

func luaCreate(mt *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		fn := l.CheckFunction(1)
		l.SetTop(0)
		l.Push(l.GetGlobal("setmetatable"))
		l.Push(fn)
		l.Push(mt)
		l.Call(2, 1)
		return 1
	}
}

func luaAnnotator(create *lua.LFunction) lua.LGFunction {
	return func(l *lua.LState) int {
		t := l.CheckTable(1)
		prepend := l.OptBool(2, false)

		var fn *lua.LFunction
		if prepend {
			fn = l.NewClosure(prepender(create, t), create)
		} else {
			fn = l.NewClosure(setter(create, t), create)
		}

		l.Push(create)
		l.Push(fn)
		l.Call(1, 1)
		return 1
	}
}

func setter(create *lua.LFunction, table lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		s := l.CheckString(1)
		l.SetTop(0)
		fn := l.NewClosure(setFunc(table, s), table) // close variable ``s''?
		l.Push(create)
		l.Push(fn)
		l.Call(1, 1)
		return 1
	}
}

func prepender(create *lua.LFunction, table lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		s := l.CheckString(1)
		l.SetTop(0)
		fn := l.NewClosure(prependFunc(table, s), table)
		l.Push(create)
		l.Push(fn)
		l.Call(1, 1)
		return 1
	}
}

func setFunc(table lua.LValue, s string) lua.LGFunction {
	return func(l *lua.LState) int {
		val := l.Get(1)
		l.SetTable(table, val, lua.LString(s))
		return 1
	}
}

func prependFunc(table lua.LValue, s string) lua.LGFunction {
	return func(l *lua.LState) int {
		val := l.Get(1)
		t := l.GetTable(table, val)
		if t == lua.LNil {
			t = l.NewTable()
		}
		l.Push(l.GetField(l.GetGlobal("table"), "insert"))
		l.Push(t)
		l.Push(lua.LNumber(1))
		l.Push(lua.LString(s))
		l.Call(3, 0)
		l.SetTable(table, val, t)
		return 1
	}
}
