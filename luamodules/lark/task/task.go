package task

import (
	"log"

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
	patterns := weakTable(l, setmt, "kv")

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
	annotator := l.GetField(l.Get(-1), "annotator")
	l.Pop(1)

	newAnnotator := func(t lua.LValue, prepend bool) lua.LValue {
		l.Push(annotator)
		l.Push(t)
		l.Push(lua.LBool(prepend))
		l.Call(2, 1)
		val := l.Get(-1)
		l.Pop(1)
		return val
	}

	nameFunc := l.NewClosure(luaName(decorator, namedTasks), decorator, namedTasks)
	l.Push(decorator)
	l.Push(nameFunc)
	l.Call(1, 1)
	name := l.Get(-1)
	l.Pop(1)
	doc.Go(l, name, &doc.GoDocs{
		Desc: "A decorator that names its task.",
	})

	pattern := newAnnotator(patterns, false)
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
	l.SetField(mod, "name", name)
	l.SetField(mod, "pattern", pattern)
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

func luaFind(anonTasks, namedTasks, patterns lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)

		l.ForEach(namedTasks.(*lua.LTable), func(k, v lua.LValue) {
			log.Printf("KEY %q", k)
		})

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

func weakTable(l *lua.LState, setmt *lua.LFunction, mode string) lua.LValue {
	mt := l.NewTable()
	l.SetField(mt, "__mode", lua.LString(mode))

	l.Push(setmt)
	l.Push(l.NewTable())
	l.Push(mt)
	l.Call(2, 1)
	val := l.Get(l.GetTop())
	l.Pop(1)
	return val
}
