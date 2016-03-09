package fun

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/yuin/gopher-lua"
)

// Module provides the gopher-lua "fun" module.
var Module = gluamodule.New("fun", Loader,
	doc.Module,
)

// Loader loads the module
func Loader(l *lua.LState) int {
	mod := l.NewTable()
	tfuncs := l.NewTable()
	dtor := doc.Must(l)

	upvalues := []lua.LValue{
		upModule: mod,
		upFuncs:  tfuncs,
	}
	for k, v := range funcs {
		fn := l.NewClosure(v.fn, upvalues...)
		dtor.MustDoc(fn, v.doc)
		tfuncs.RawSetString(k, fn)
	}
	tfuncs.ForEach(func(k, v lua.LValue) {
		s := k.(lua.LString)
		if exports[string(s)] {
			mod.RawSetString(string(s), v)
		}
	})

	l.Push(mod)

	return 1
}

type funUpvalue int

const (
	upModule funUpvalue = iota + 1
	upFuncs
)

func (up funUpvalue) index() int {
	return lua.UpvalueIndex(int(up))
}

type funIndex int

var funcs = map[string]struct {
	fn  lua.LGFunction
	doc *doc.Docs
}{
	"flatten": {
		luaFlatten,
		&doc.Docs{
			Desc: "Returns flat array containing non-table elements of nested array values.",
			Sig:  "t => tmap",
			Params: []string{
				`
				t  array
				An array with possibly nested tables.
				`,
				`
				tmap  array
				A copy of t with nested arrays flattened.
				`,
				`
				d  (optional) number
				A depth at which to stop flattening.  A value of zero will
				return a copy of the array t.  A negative value or a nil value
				will flatten nested arrays at all depths.
				`,
			},
		},
	},
	"map": {
		luaMap,
		&doc.Docs{
			Desc: "Returns a copy of a table with its key-value pairs transformed by a given function.",
			Sig:  "(t, (k, v) => (kmap, vmap)) => tmap",
			Params: []string{
				`
				t  table
				A table.
				`,
				`
				tmap  table
				A copy of t with its key-value pairs transformed by the given function.
				`,
				`
				k   key
				A key contained in t.
				`,
				`
				v   key
				The value corresponding to k in t.
				`,
				`
				kmap   (optional) key
				A key to store in tmap corresponding to the input key-value
				pair.  If no value is returned then k will be used as a key in
				tmap.
				`,
				`
				vmap   (optional) any
				The value to store in tmap at key kmap. If nil, or no value is
				returned then kmap is removed from tmap.
				`,
			},
		},
	},
	"vmap": {
		luaVMap,
		&doc.Docs{
			Desc: "Returns a copy of a table with its values transformed by a given function.",
			Sig:  "(t, (v) => (vmap)) => tmap",
			Params: []string{
				`
				t     table
				A table.
				`,
				`
				tmap  table
				A copy of t with its values transformed by the given function.
				`,
				`
				v     any
				A value in t.
				`,
				`
				vmap  (optional) any
				The value to store in tmap corresponding to v.  If nil, or no
				value is returned then key which v is associated with in t is
				not included in tmap.
				`,
			},
		},
	},
	"sel": {
		luaSelect,
		nil,
	},
	"vsel": {
		luaVSelect,
		&doc.Docs{
			Desc: "Returns an table containing elements with values matched by a given function.",
			Sig:  "(t, (v) => keep) => tsel",
			Params: []string{
				`
				t     table
				A table.
				`,
				`
				v     any
				A value in t.
				`,
				`
				keep  boolean
				If true then the input value will be included in tsel under the
				same key it had in t.
				`,
				`
				tsel  table
				A table containing values from t for which keep was true.
				`,
			},
		},
	},
}

var exports = map[string]bool{
	"flatten": true,
	"map":     true,
	"vmap":    true,
	"sel":     true,
	"vsel":    true,
}

func luaFlatten(l *lua.LState) int {
	t := l.CheckTable(1)
	d := l.OptInt(2, -1)

	dest := l.NewTable()
	flattenArray(dest, t, d)
	l.Push(dest)
	return 1
}

func flatten(dest *lua.LTable, val lua.LValue, d int) {
	if d == 0 {
		dest.Append(val)
		return
	}

	t, ok := val.(*lua.LTable)
	if ok {
		flattenArray(dest, t, d-1)
	} else {
		dest.Append(val)
	}
}

func flattenArray(dest *lua.LTable, arr *lua.LTable, d int) {
	n := arr.Len()
	for i := 1; i <= n; i++ {
		flatten(dest, arr.RawGetInt(i), d)
	}
}

func luaVMap(l *lua.LState) int {
	a := l.CheckTable(1)
	fn := l.CheckFunction(2)
	narg := l.GetTop()

	b := l.NewTable()
	l.ForEach(a, func(k, v lua.LValue) {
		l.Push(fn)
		l.Push(v)
		for i := 3; i <= narg; i++ {
			l.Push(l.Get(i))
		}
		l.Call(narg-1, 1)
		vmap := l.Get(-1)
		l.Pop(1)
		b.RawSet(k, vmap)
	})

	l.Push(b)
	return 1
}

func luaMap(l *lua.LState) int {
	a := l.CheckTable(1)
	fn := l.CheckFunction(2)
	narg := l.GetTop()

	b := l.NewTable()
	l.ForEach(a, func(k, v lua.LValue) {
		l.Push(fn)
		l.Push(k)
		l.Push(v)
		for i := 3; i <= narg; i++ {
			l.Push(l.Get(i))
		}
		l.Call(narg, lua.MultRet)
		var kmap, vmap lua.LValue
		nret := l.GetTop() - narg
		switch nret {
		case 0:
			return
		case 1:
			vmap = l.Get(-1)
			kmap = k
		default:
			kmap = l.Get(-2)
			vmap = l.Get(-1)
		}
		l.Pop(nret)
		b.RawSet(kmap, vmap)
	})

	l.Push(b)
	return 1
}

func luaSelect(l *lua.LState) int {
	a := l.CheckTable(1)
	fn := l.CheckFunction(2)
	narg := l.GetTop()

	b := l.NewTable()
	l.ForEach(a, func(k, v lua.LValue) {
		l.Push(fn)
		l.Push(k)
		l.Push(v)
		for i := 3; i <= narg; i++ {
			l.Push(l.Get(i))
		}
		l.Call(narg, 1)
		if !lua.LVIsFalse(l.Get(-1)) {
			if isArrayIndex(k) {
				b.Append(v)
			} else {
				b.RawSet(k, v)
			}
		}
		l.Pop(1)
	})

	l.Push(b)
	return 1
}

func luaVSelect(l *lua.LState) int {
	a := l.CheckTable(1)
	fn := l.CheckFunction(2)
	narg := l.GetTop()

	b := l.NewTable()
	l.ForEach(a, func(k, v lua.LValue) {
		l.Push(fn)
		l.Push(v)
		for i := 3; i <= narg; i++ {
			l.Push(l.Get(i))
		}
		l.Call(narg-1, 1)
		if !lua.LVIsFalse(l.Get(-1)) {
			if isArrayIndex(k) {
				b.Append(v)
			} else {
				b.RawSet(k, v)
			}
		}
		l.Pop(1)
	})

	l.Push(b)
	return 1
}

func isArrayIndex(v lua.LValue) bool {
	x, ok := v.(lua.LNumber)
	return ok && isInt(x) && x > 0 && x < lua.LNumber(lua.MaxArrayIndex)
}

func isInt(x lua.LNumber) bool {
	return x == lua.LNumber(int64(x))
}
