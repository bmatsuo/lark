package doc

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/internal/lx"
	"github.com/bmatsuo/lark/internal/textutil"
	"github.com/bmatsuo/lark/lib/decorator/_intern"
	"github.com/yuin/gopher-lua"
)

// Module is a gluamodule.Module that loads the doc module.
var Module = gluamodule.New("doc", docLoader,
	intern.Module,
)

// Get loads documentation about lv from l.
func Get(l *lua.LState, lv lua.LValue, name string) (*Docs, error) {
	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString(Module.Name()))
	err := l.PCall(1, 1, nil)
	if err != nil {
		return nil, err
	}

	doc := l.Get(-1)
	l.Pop(1)

	l.Push(l.GetField(doc, "get"))
	l.Push(lv)
	l.Push(lua.LNumber(-1))
	err = l.PCall(2, 1, nil)
	if err != nil {
		return nil, err
	}

	ld := l.Get(-1)
	l.Pop(1)

	return decodeDocs(l, ld, name)
}

func decodeDocs(l *lua.LState, lv lua.LValue, name string) (*Docs, error) {
	if lv == lua.LNil {
		return nil, nil
	}
	if lv.Type() != lua.LTTable {
		return nil, fmt.Errorf("not a table: %s", lv.Type())
	}

	ldocs := &lx.NamedValue{
		Name:  fmt.Sprintf("%s docs", name),
		Value: lv,
	}
	if name == "" {
		ldocs.Name = name
	}

	d := &Docs{}

	d.Usage = lx.OptStringField(l, ldocs, "usage", "")
	d.Sig = lx.OptStringField(l, ldocs, "sig", "")
	d.Desc = lx.OptStringField(l, ldocs, "desc", "")
	lparams := lx.OptTableField(l, ldocs, "params", nil)
	lvars := lx.OptTableField(l, ldocs, "vars", nil)
	lsubs := lx.OptTableField(l, ldocs, "sub", nil)
	if lparams != nil {
		l.ForEach(lparams, func(k, v lua.LValue) {
			s, ok := v.(lua.LString)
			if !ok {
				return
			}
			d.Params = append(d.Params, string(s))
		})
	}
	if lvars != nil {
		l.ForEach(lvars, func(k, v lua.LValue) {
			s, ok := v.(lua.LString)
			if !ok {
				return
			}
			d.Vars = append(d.Vars, string(s))
		})
	}
	if lsubs != nil {
		var suberr error
		l.ForEach(lsubs, func(k, v lua.LValue) {
			_, ok := k.(lua.LNumber)
			if !ok {
				suberr = fmt.Errorf("field sub of %s docs has non-numeric", name)
				return
			}
			_, ok = v.(*lua.LTable)
			if !ok {
				suberr = fmt.Errorf("field sub of %s docs is not a table: %s", name, v.Type())
				return
			}
			lsubs := &lx.NamedValue{
				Name:  fmt.Sprintf("%s docs sub", name),
				Value: v,
			}
			var err error
			s := &Sub{}
			s.Name = lx.StringField(l, lsubs, "name")
			s.Type = lx.StringField(l, lsubs, "type")
			docs := lx.OptTableField(l, lsubs, "docs", nil)
			if docs != nil {
				fullName := name + "." + s.Name
				s.Docs, err = decodeDocs(l, docs, fullName)
				if err != nil {
					suberr = err
					return
				}
			}
			d.Subs = append(d.Subs, s)
		})
		if suberr != nil {
			return nil, suberr
		}
	}
	return d, nil
}

// Docs represents documentation for a Lua object.
type Docs struct {
	Usage  string
	Sig    string
	Desc   string
	Params []string
	Vars   []string
	Subs   []*Sub
}

// NumVar returns the number of variables declared for d.
func (d *Docs) NumVar() int {
	return len(d.Vars)
}

func splitNamed(named string) (name, rest string) {
	text := strings.TrimSpace(named)
	index := strings.IndexFunc(text, unicode.IsSpace)
	if index < 0 {
		return text, ""
	}
	return text[:index], text[index:]
}

// Var returns the name and description variable i in d.
func (d *Docs) Var(i int) (name string) {
	name, _ = splitNamed(d.Vars[i])
	return name
}

// VarDesc returns the description of variable i in d.
func (d *Docs) VarDesc(i int) (desc string) {
	_, rest := splitNamed(d.Vars[i])
	return rest
}

// VarType returns the type of variable i in d if it can be inferred from the
// documentation.
//
// BUG:
// VarType is not implemented.  No convention has been settled on.
func (d *Docs) VarType(i int) (typ string) {
	return ""
}

// NumParam returns the number of parameters declared for d.
func (d *Docs) NumParam() int {
	return len(d.Params)
}

// Param returns the name and description parameter i in d.
func (d *Docs) Param(i int) (name string) {
	name, _ = splitNamed(d.Params[i])
	return name
}

// ParamDesc returns the description of parameter i in d.
func (d *Docs) ParamDesc(i int) (desc string) {
	_, rest := splitNamed(d.Params[i])
	return rest
}

// ParamType returns the type of parameter i in d if it can be inferred from the
// documentation.
//
// BUG:
// ParamType is not implemented.  No convention has been settled on.
func (d *Docs) ParamType(i int) (typ string) {
	return ""
}

// Funcs returns function subtopics of d.
func (d *Docs) Funcs() []*Sub {
	var sub []*Sub
	for _, s := range d.Subs {
		if s.Type == "function" {
			sub = append(sub, s)
		}
	}
	return sub
}

// Others returns non-function subtopics of d.
func (d *Docs) Others() []*Sub {
	var sub []*Sub
	for _, s := range d.Subs {
		if s.Type != "function" {
			sub = append(sub, s)
		}
	}
	return sub
}

// Sub is subtopic documentation for a Lua object.
type Sub struct {
	Name string
	Type string
	*Docs
}

// Go sets the description for obj to desc.  Go ignores doc.Subs, functions and
// documented variables must have their documentation declared separately.
func Go(l *lua.LState, obj lua.LValue, doc *Docs) {
	require := l.GetGlobal("require")
	l.Push(require)
	l.Push(lua.LString("doc"))
	l.Call(1, 1)
	mod := l.CheckTable(-1)
	l.Pop(1)

	ndec := 0
	if doc.Usage != "" {
		l.Push(l.GetField(mod, "usage"))
		l.Push(lua.LString(doc.Usage))
		err := l.PCall(1, 1, nil)
		if err != nil {
			l.RaiseError("%s", err)
		}
		ndec++
	}
	if doc.Sig != "" {
		l.Push(l.GetField(mod, "sig"))
		l.Push(lua.LString(doc.Sig))
		err := l.PCall(1, 1, nil)
		if err != nil {
			l.RaiseError("%s", err)
		}
		ndec++
	}
	if doc.Desc != "" {
		l.Push(l.GetField(mod, "desc"))
		l.Push(lua.LString(doc.Desc))
		err := l.PCall(1, 1, nil)
		if err != nil {
			l.RaiseError("%s", err)
		}
		ndec++
	}
	if len(doc.Vars) > 0 {
		_var := l.GetField(mod, "var")
		for _, v := range doc.Vars {
			l.Push(_var)
			l.Push(lua.LString(v))
			err := l.PCall(1, 1, nil)
			if err != nil {
				l.RaiseError("%s", err)
			}
			ndec++
		}
	}
	if len(doc.Params) > 0 {
		param := l.GetField(mod, "param")
		for _, p := range doc.Params {
			l.Push(param)
			l.Push(lua.LString(p))
			err := l.PCall(1, 1, nil)
			if err != nil {
				l.RaiseError("%s", err)
			}
			ndec++
		}
	}
	l.Push(obj)
	for i := 0; i < ndec; i++ {
		err := l.PCall(1, 1, nil)
		if err != nil {
			l.RaiseError("%s", err)
		}
	}
}

func docLoader(l *lua.LState) int {
	mod := l.NewTable()

	setmt, ok := l.GetGlobal("setmetatable").(*lua.LFunction)
	if !ok {
		l.RaiseError("unexpected type for setmetatable")
	}
	usages := weakTable(l, setmt, "kv")
	signatures := weakTable(l, setmt, "kv")
	descriptions := weakTable(l, setmt, "kv")
	parameters := weakTable(l, setmt, "k")
	variables := weakTable(l, setmt, "k")

	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString("decorator.intern"))
	l.Call(1, 1)
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

	usage := newAnnotator(usages, false)
	sig := newAnnotator(signatures, false)
	desc := newAnnotator(descriptions, false)
	param := newAnnotator(parameters, true)
	_var := newAnnotator(variables, true)

	dodoc := func(obj lua.LValue, u, s, d string, ps ...string) {
		ncall := 3 + len(ps)
		if u != "" {
			l.Push(usage)
			l.Push(lua.LString(u))
			l.Call(1, 1)
		} else {
			ncall--
		}
		if s != "" {
			l.Push(sig)
			l.Push(lua.LString(s))
			l.Call(1, 1)
		} else {
			ncall--
		}
		if d != "" {
			l.Push(desc)
			l.Push(lua.LString(d))
			l.Call(1, 1)
		} else {
			ncall--
		}
		for _, p := range ps {
			l.Push(param)
			l.Push(lua.LString(p))
			l.Call(1, 1)
		}
		l.Push(obj)
		for i := 0; i < ncall; i++ {
			l.Call(1, 1)
		}
	}

	dodoc(mod,
		"local doc = require('doc')",
		"",
		`
		The doc module contains utilities for documenting Lua objects using
		decorators.  Sections of documentation are declared separately using
		small idiomatically named decorators.  Decorators are defined for
		documenting (module) table descriptions, variables, and functions.  For
		function decorators are defined to document signatures and parameter
		values.
		`,
	)
	dodoc(usage,
		"",
		"s => fn => fn",
		"A decorator that documents the usage of an object.",
		`s  string -- Text describing usage.`,
	)
	dodoc(sig,
		"",
		"s => fn => fn",
		"A decorator that documents a function's signature.",
		`s  string -- The function signature.`,
	)
	dodoc(desc,
		"",
		"s => fn => fn",
		"A decorator that describes an object.",
		`s  string -- The object description.`,
	)
	dodoc(param,
		"",
		"s => fn => fn",
		"A decorator that describes a function parameter.",
		`s  string -- The parameter name and description separated by white space.`,
	)
	dodoc(_var,
		"",
		"s => fn => fn",
		"A decorator that describes module variable (table field).",
		`s  string -- The variable name and description separated by white space.`,
	)

	get := l.NewClosure(
		luaGet(usages, signatures, descriptions, parameters, variables),
		usages, signatures, descriptions, parameters, variables,
	)
	dodoc(get,
		"",
		"obj => table",
		"Retrieve a table containing documentation for obj.",
		`obj   table, function, or userdata -- The object to retrieve documentation for.`,
	)

	help := l.NewClosure(
		luaHelp(mod, get),
		mod, get,
	)
	dodoc(help,
		"",
		"obj => ()",
		"Print the documentation for obj.",
		`obj   table, function, or userdata -- The object to retrieve documentation for.`,
	)

	// decorators
	l.SetField(mod, "usage", usage)
	l.SetField(mod, "sig", sig)
	l.SetField(mod, "desc", desc)
	l.SetField(mod, "var", _var)
	l.SetField(mod, "param", param)

	// accessors
	l.SetField(mod, "get", get)
	l.SetField(mod, "help", help)

	l.Push(mod)
	return 1
}

func luaHelp(mod lua.LValue, get lua.LValue) lua.LGFunction {
	return func(l *lua.LState) int {
		print := l.GetGlobal("print")
		if l.GetTop() == 0 {
			def := l.GetField(mod, "default")
			if def == lua.LNil {
				return 0
			}
			lstr, ok := l.ToStringMeta(def).(lua.LString)
			if !ok {
				l.RaiseError("default is not a string")
			}
			str := textutil.Unindent(string(lstr))
			str = textutil.Wrap(str, 72)
			l.Push(print)
			l.Push(lua.LString(str))
			l.Call(1, 0)
			return 0
		}

		val := l.Get(1)
		l.SetTop(0)
		l.Push(get)
		l.Push(val)
		l.Call(1, 1)

		docs := l.Get(-1)
		l.SetTop(0)

		godocs, err := decodeDocs(l, docs, "")
		if err != nil {
			l.RaiseError("%s", err)
		}
		text, err := NewFormatter().Format(nil, godocs, "")
		if err != nil {
			l.RaiseError("%s", err)
		}
		_, err = io.Copy(os.Stderr, strings.NewReader(text))
		if err != nil {
			l.RaiseError("%s", err)
		}
		return 0
	}
}

func luaGet(usages, signatures, descriptions, parameters, variables lua.LValue) lua.LGFunction {
	var rec lua.LGFunction
	rec = func(l *lua.LState) int {
		val := l.Get(1)
		depth := l.OptInt(2, 1)

		l.SetTop(0)
		usage := l.GetTable(usages, val)
		sig := l.GetTable(signatures, val)
		desc := l.GetTable(descriptions, val)
		params := l.GetTable(parameters, val)
		vars := l.GetTable(variables, val)
		tab, ok := val.(*lua.LTable)
		if sig == lua.LNil && desc == lua.LNil && params == lua.LNil && vars == lua.LNil && !ok {
			l.Push(lua.LNil)
			return 1
		}

		t := l.NewTable()
		l.SetField(t, "usage", usage)
		l.SetField(t, "sig", sig)
		l.SetField(t, "desc", desc)
		l.SetField(t, "params", params)
		l.SetField(t, "vars", vars)

		if tab != nil {
			topics := l.NewTable()

			l.ForEach(tab, func(k, v lua.LValue) {
				_, ok := k.(lua.LString)
				if !ok {
					return
				}
				_, ok = v.(*lua.LFunction)
				if !ok {
					return
				}

				subTopic := l.NewTable()
				l.SetField(subTopic, "name", k)
				l.SetField(subTopic, "type", lua.LString("function"))

				if depth > 0 || depth < 0 {
					l.Push(v)
					l.Push(lua.LNumber(depth - 1))
					if rec(l) != 1 {
						l.RaiseError("oh no my hack failed!")
					}

					subDocs := l.Get(-1)
					l.Pop(1)

					if subDocs != lua.LNil {
						l.SetField(subTopic, "docs", subDocs)
					}
				}

				topics.Append(subTopic)
			})

			l.SetField(t, "sub", topics)
		}

		l.Push(t)
		return 1
	}
	return rec
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

func splitParam(s string) (name, desc string) {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	i := strings.IndexFunc(s, unicode.IsSpace)
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i:]
}
