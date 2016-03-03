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

	nameFunc := l.NewClosure(
		luaName(decorator, namedTasks, mod),
		decorator, namedTasks, mod,
	)
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

	createFunc := l.NewClosure(
		luaCreate(anonTasks, mod),
		anonTasks, mod,
	)
	l.Push(decorator)
	l.Push(createFunc)
	l.Call(1, 1)
	create := l.Get(-1)
	l.Pop(1)
	doc.Go(l, create, &doc.GoDocs{
		Desc: "A decorator that defines an anonymous task.",
	})

	find := l.NewClosure(
		luaFind(anonTasks, namedTasks, patterns, mod),
		anonTasks, namedTasks, patterns, mod,
	)
	doc.Go(l, find, &doc.GoDocs{
		Desc: "Find the task by the given name.",
	})

	dump := l.NewClosure(
		luaDump(anonTasks, namedTasks, patterns, mod),
		anonTasks, namedTasks, patterns, mod,
	)
	doc.Go(l, dump, &doc.GoDocs{
		Desc: "Write all task names and patterns to standard output.",
	})

	l.SetField(mod, "create", create)
	l.SetField(mod, "with_name", name)
	l.SetField(mod, "with_pattern", pattern)
	l.SetField(mod, "find", find)
	l.SetField(mod, "dump", dump)
	l.SetField(mod, "run", l.NewClosure(
		luaRun(find),
		find,
	))
	l.SetField(mod, "get_name", l.NewClosure(luaGetName))
	l.SetField(mod, "get_pattern", l.NewClosure(luaGetPattern))
	l.SetField(mod, "get_param", l.NewClosure(luaGetParam))

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

func luaFind(anonTasks, namedTasks, patterns, mod *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		var noname bool
		var name string
		if l.GetTop() > 0 {
			name = l.CheckString(1)
		} else {
			noname = true
			var lname lua.LString
			var ok bool
			def := l.GetField(mod, "default")
			lname, ok = def.(lua.LString)
			if !ok {
				l.ForEach(l.Get(lua.GlobalsIndex).(*lua.LTable), func(k, v lua.LValue) {
					if !l.Equal(v, def) {
						return
					}
					lname, ok = k.(lua.LString)
					if !ok {
						l.RaiseError("unexpected global index")
					}
				})
			}
			name = string(lname)
			if name == "" {
				l.RaiseError("cannot determine name of task")
			}
		}

		val := l.GetField(namedTasks, name)
		if val != lua.LNil {
			l.Push(val)
			l.Push(lua.LString(name))
			return 2
		}

		val = l.GetGlobal(name)
		if val != lua.LNil {
			isTask, ok := l.GetTable(anonTasks, val).(lua.LBool)
			if ok && bool(isTask) {
				l.Push(val)
				l.Push(lua.LString(name))
				return 2
			}
		}

		if noname {
			return 0
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
			l.Push(lua.LString(name))
			l.Push(found)
			return 3
		}

		return 0
	}
}

func luaDump(anonTasks, namedTasks, patterns, mod *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		print := l.GetGlobal("print")
		def := l.GetField(mod, "default")

		set := l.NewTable()
		l.ForEach(namedTasks, func(k, v lua.LValue) {
			l.Push(print)
			l.Push(lua.LString("="))
			l.Push(k)
			if l.Equal(def, k) {
				l.Push(lua.LString(" (default)"))
				l.Call(3, 0)
			} else {
				l.Call(2, 0)
			}

			l.SetTable(set, k, lua.LBool(true))
		})
		l.ForEach(anonTasks, func(val, _ lua.LValue) {
			l.ForEach(l.Get(lua.GlobalsIndex).(*lua.LTable), func(k, v lua.LValue) {
				if !l.Equal(v, val) {
					return
				}
				lname, ok := k.(lua.LString)
				if !ok {
					return
				}
				if l.GetField(set, string(lname)) == lua.LNil {
					l.Push(print)
					l.Push(lua.LString("-"))
					l.Push(lname)
					if l.Equal(def, v) {
						l.Push(lua.LString(" (default)"))
						l.Call(3, 0)
					} else {
						l.Call(2, 0)
					}
				}
			})
		})
		l.ForEach(patterns, func(k, v lua.LValue) {
			l.Push(print)
			l.Push(lua.LString("~"))
			l.Push(l.GetField(v, "pattern"))
			l.Call(2, 0)
		})

		return 0
	}
}

func luaCreate(t lua.LValue, mod *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		val := l.CheckAny(1)
		if l.GetField(mod, "default") == lua.LNil {
			l.SetField(mod, "default", val)
		}
		l.SetTable(t, val, lua.LBool(true))
		return 1
	}
}

func luaName(decorator *lua.LFunction, t lua.LValue, mod *lua.LTable) lua.LGFunction {
	return func(l *lua.LState) int {
		name := l.CheckString(1)

		fn := l.NewClosure(func(l *lua.LState) int {
			val := l.CheckAny(1)
			if l.GetField(mod, "default") == lua.LNil {
				l.SetField(mod, "default", lua.LString(name))
			}
			l.SetField(t, name, val)
			return 1
		}, t, mod)

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

func luaRun(find *lua.LFunction) lua.LGFunction {
	return func(l *lua.LState) int {
		var name string
		lname, ok := l.Get(1).(lua.LString)
		if ok {
			name = string(lname)
		} else if l.Get(1) != lua.LNil {
			l.CheckString(1) // call will raise an error
		}
		params := lua.LValue(lua.LNil)
		if l.GetTop() > 1 {
			params = l.CheckTable(2)
		}
		l.SetTop(0)

		if name != "" {
			l.Push(find)
			l.Push(lua.LString(name))
			l.Call(1, 3)
			if l.Get(1) == lua.LNil {
				l.RaiseError("no task matching name: %s", name)
			}
		} else {
			var ok bool
			l.Push(find)
			l.Call(0, 2)
			lname, ok = l.Get(2).(lua.LString)
			if !ok {
				l.RaiseError("task name is not a string")
			}
			name = string(lname)
		}
		patt := lua.LValue(lua.LNil)
		if l.GetTop() > 2 {
			patt = l.Get(3)
		}
		l.SetTop(1)

		ctx := l.NewTable()
		l.SetField(ctx, "name", lua.LString(name))
		l.SetField(ctx, "pattern", patt)
		l.SetField(ctx, "params", params)
		l.Push(ctx)
		l.Call(1, 0)
		return 0
	}
}

func luaGetName(l *lua.LState) int {
	if l.GetTop() == 0 {
		return 0
	}
	ctx := l.CheckTable(1)
	l.Replace(1, l.GetField(ctx, "name"))
	return 1
}

func luaGetPattern(l *lua.LState) int {
	if l.GetTop() == 0 {
		return 0
	}
	ctx := l.CheckTable(1)
	l.Replace(1, l.GetField(ctx, "pattern"))
	return 1
}

func luaGetParam(l *lua.LState) int {
	ctx := l.CheckTable(1)
	name := l.CheckString(2)
	l.SetTop(0)
	params := l.GetField(ctx, "params")
	if params == lua.LNil {
		return 0
	}
	l.Push(l.GetField(params, name))
	return 1
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
