package decorator

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/yuin/gopher-lua"
)

// Module is a gluamodule.Module that loads the "decorator" module.
var Module = gluamodule.New("decorator", Loader)

// Loader is a lua.LGFunction that loads the module.
func Loader(l *lua.LState) int {
	mod := l.NewTable()

	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString("decorator.intern"))
	l.Call(1, 1)
	internal := l.Get(-1)
	l.Pop(1)

	metatable := l.GetField(internal, "metatable")
	l.SetField(mod, "metatable", metatable)
	doc.Go(l, metatable, &doc.GoDocs{
		Sig: "call => mt",
		Desc: `
			Create a new metatable for a basic decorator with concat/call
			syntax.
			`,
		Params: []string{
			`call function
			-- The value of __call in the returned metatable
			`,
			`mt table
			-- A metamethod table with __call and __concat set to call.
			`,
		},
	})
	create := l.GetField(internal, "create")
	l.SetField(mod, "create", create)
	doc.Go(l, create, &doc.GoDocs{
		Sig: "dec => obj => obj",
		Desc: `
			Return a simple copy/concat decorator using the given decorating
			function.
			`,
		Params: []string{
			`dec function
			-- The decorating function.  Typically dec will return the same
			object it is given though it is free to wrap or transform the
			value.
			`,
			`obj any
			-- A value to be decorated.
			`,
		},
	})
	annotator := l.GetField(internal, "annotator")
	l.SetField(mod, "annotator", annotator)
	doc.Go(l, annotator, &doc.GoDocs{
		Sig: "(tab, prepend) => annot => obj => obj",
		Desc: `
			Return a copy/concat decorator that 
			`,
		Params: []string{
			`tab function
			-- The table in which annotations are stored.  The table may employ
			weak references but this is not a requirement.
			`,
			`prepend (optional) boolean
			-- When true multiple (chained) annotations on the same obj will be
			prepended in an array instead of being overwritten.  Prepending
			makes the apparent order equal to the insertion order (opposite
			call resolution order).
			`,
			`annot any
			-- An annotation, typically a string, that is associated with given
			objects.
			`,
			`obj any
			-- A value to be decorated.
			`,
		},
	})

	l.Push(mod)
	return 1
}
