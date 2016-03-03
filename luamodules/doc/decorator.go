package doc

import "github.com/yuin/gopher-lua"

func luaDecorator(l *lua.LState) int {
	fn := l.CheckFunction(1)
	dec := l.NewTable()
	l.SetField(dec, "fn", fn)
	l.Push(l.GetGlobal("setmetatable"))
	l.Push(dec)
	l.Push(l.NewFunction(luaDecoratorMetatable))
	l.Call(0, 1)
	l.Call(2, 1)
	return 1
}

func luaDecoratorMetatable(l *lua.LState) int {
	mt := l.NewTable()
	l.SetField(mt, "__concat", l.NewFunction(luaDecoratorConcat))
	l.SetField(mt, "__call", l.NewFunction(luaDecoratorCall))
	l.Push(mt)
	return 1
}

func luaDecoratorConcat(l *lua.LState) int {
	dec := l.Get(1)
	if l.GetTop() < 2 {
		l.RaiseError("nothing to concatenate")
	}
	l.SetTop(2)
	fn := l.GetField(dec, "fn")
	l.Replace(1, fn)
	l.Call(1, 1)
	return 1
}

func luaDecoratorCall(l *lua.LState) int {
	dec := l.Get(1)
	if l.GetTop() < 1 {
		l.RaiseError("nothing to call")
	}
	l.Replace(1, l.GetField(dec, "fn"))
	narg := l.GetTop() - 1
	l.Call(narg, 1)
	return 1
}
