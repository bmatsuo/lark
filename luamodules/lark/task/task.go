package task

import (
	"github.com/bmatsuo/lark/internal/module"
	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/yuin/gopher-lua"
)

// Module loads the "task" module.
var Module = module.New("lark.task", Loader)

// Loader loads the Lua module.
func Loader(l *lua.LState) int {
	mod := l.NewTable()

	setmt, ok := l.GetGlobal("setmetatable").(*lua.LFunction)
	if !ok {
		l.RaiseError("unexpected type for setmetatable")
	}
	anonTasks := weakTable(l, setmt, "k")
	namedTasks := weakTable(l, setmt, "kv")
	patterns := weakTable(l, setmt, "k")

	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString("decorator"))
	l.Call(1, 1)
	decorator, ok := l.GetField(l.Get(-1), "create").(*lua.LFunction)
	if !ok {
		l.RaiseError("unexpected type for decorator.create")
	}
	metatable, ok := l.GetField(l.Get(-1), "metatable").(*lua.LFunction)
	if !ok {
		l.RaiseError("unexpected type for decorator.create")
	}
	l.Pop(1)

	nameFunc := l.NewClosure(luaName(decorator, namedTasks), decorator, namedTasks)
	l.Push(decorator)
	l.Push(nameFunc)
	l.Call(1, 1)
	name := l.Get(-1)
	l.Pop(1)
	doc.Go(l, name, &doc.GoDocs{
		Desc: "A decorator that names its task.",
	})

	patternFunc := l.NewClosure(
		luaPattern(setmt, decorator, patterns),
		setmt, decorator, patterns,
	)
	l.Push(decorator)
	l.Push(patternFunc)
	l.Call(1, 1)
	pattern := l.Get(-1)
	l.Pop(1)
	doc.Go(l, pattern, &doc.GoDocs{
		Desc: "A decorator that defines a pattern for its task.",
	})

	createFunc := l.NewClosure(luaCreate(anonTasks), anonTasks)
	l.Push(decorator)
	l.Push(createFunc)
	l.Call(1, 1)
	create := l.Get(-1)
	l.Pop(1)
	doc.Go(l, create, &doc.GoDocs{
		Desc: "A decorator that defines an anonymous task.",
	})

	find := l.NewClosure(
		luaFind(anonTasks, namedTasks, patterns),
		anonTasks, namedTasks, patterns,
	)
	doc.Go(l, find, &doc.GoDocs{
		Desc: "Find the task by the given name.",
	})

	l.SetField(mod, "create", create)
	l.SetField(mod, "with_name", name)
	l.SetField(mod, "with_pattern", pattern)
	l.SetField(mod, "find", find)

	// setmetatable and return mod
	l.Push(setmt)
	l.Push(mod)
	l.Push(metatable)
	l.Push(l.NewFunction(luaDecorator))
	l.Call(1, 1)
	l.Call(2, 1)
	return 1
}

func luaDecorator(l *lua.LState) int {
	mod := l.CheckAny(1)
	create := l.GetField(mod, "create")
	l.Replace(1, create)
	narg := l.GetTop() - 1
	l.Call(narg, lua.MultRet)
	return l.GetTop()
}

func luaFind(anonTasks, namedTasks, patterns *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)

		val := l.GetField(namedTasks, name)
		if val != lua.LNil {
			l.Push(val)
			return 1
		}

		val = l.GetGlobal(name)
		if val != lua.LNil {
			isTask, ok := l.GetTable(anonTasks, val).(lua.LBool)
			if ok && bool(isTask) {
				l.Push(val)
				return 1
			}
		}

		allPatterns := l.NewTable()
		l.ForEach(patterns, func(k, v lua.LValue) {
			allPatterns.Append(v)
		})
		l.Push(l.GetField(l.GetGlobal("table"), "sort"))
		l.Push(allPatterns)
		l.Push(l.NewClosure(func(l *lua.LState) int {
			t1 := l.CheckTable(1)
			t2 := l.CheckTable(2)
			i1 := l.GetField(t1, "index")
			i2 := l.GetField(t2, "index")
			l.SetTop(0)
			l.Push(lua.LBool(l.LessThan(i1, i2)))
			return 1
		}))
		var found lua.LValue
		find := l.GetField(l.GetGlobal("string"), "find")
		l.ForEach(allPatterns, func(k, v lua.LValue) {
			patt := l.GetField(v, "pattern")
			l.Push(find)
			l.Push(lua.LString(name))
			l.Push(patt)
			l.Call(2, 1)
			val := l.Get(-1)
			l.Pop(1)
			if val != lua.LNil && found == nil {
				found = patt
			}
		})
		if found != nil {
			rec := l.GetTable(patterns, found)
			l.Push(l.GetField(rec, "value"))
			return 1
		}

		return 0
	}
}

func luaCreate(t lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		val := l.CheckAny(1)
		l.SetTable(t, val, lua.LBool(true))
		return 1
	}
}

func luaName(decorator *lua.LFunction, t lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)

		fn := l.NewClosure(func(l *lua.LState) int {
			val := l.CheckAny(1)
			l.SetField(t, name, val)
			return 1
		}, t)

		l.Push(decorator)
		l.Push(fn)
		l.Call(1, 1)
		return 1
	}
}

func luaPattern(setmt, decorator *lua.LFunction, t lua.LValue) lua.LGFunction {
	var numPatt int64
	return func(l *lua.LState) int {
		patt := l.CheckString(1)

		fn := l.NewClosure(func(l *lua.LState) int {
			val := l.CheckAny(1)
			rec := l.NewTable()
			numPatt++
			l.SetField(rec, "index", lua.LNumber(numPatt))
			l.SetField(rec, "pattern", lua.LString(patt))
			l.SetField(rec, "value", val)
			l.SetField(t, patt, rec)
			mt := l.NewTable()
			l.SetField(mt, "__mode", lua.LString("v"))
			l.Push(setmt)
			l.Push(rec)
			l.Push(mt)
			l.Call(2, 1)
			return 1
		}, t)

		l.Push(decorator)
		l.Push(fn)
		l.Call(1, 1)
		return 1
	}
}

func weakTable(l *lua.LState, setmt *lua.LFunction, mode string) *lua.LTable {
	mt := l.NewTable()
	l.SetField(mt, "__mode", lua.LString(mode))

	l.Push(setmt)
	l.Push(l.NewTable())
	l.Push(mt)
	l.Call(2, 1)
	val := l.Get(l.GetTop()).(*lua.LTable)
	l.Pop(1)
	return val
}
