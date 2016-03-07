package lx

import "github.com/yuin/gopher-lua"

// NamedValue is a lua.LValue that with a common name
type NamedValue struct {
	Name  string
	Value lua.LValue
}

// StringField retrieves a string field from nv
func StringField(l *lua.LState, nv *NamedValue, field string) string {
	lvf := l.GetField(nv.Value, field)
	ls, ok := lvf.(lua.LString)
	if !ok {
		l.RaiseError("field %s of %s is not a string: %s", field, nv.Name, nv.Value.Type())
	}
	return string(ls)
}

// BoolField retrieves a bool field from nv
func BoolField(l *lua.LState, nv *NamedValue, field string) bool {
	lvf := l.GetField(nv.Value, field)
	lb, ok := lvf.(lua.LBool)
	if !ok {
		l.RaiseError("field %s of %s is not a boolean: %s", field, nv.Name, nv.Value.Type())
	}
	return bool(lb)
}

// IntField retrieves an integer field from nv
func IntField(l *lua.LState, nv *NamedValue, field string) int {
	lvf := l.GetField(nv.Value, field)
	lx, ok := lvf.(lua.LNumber)
	if !ok {
		l.RaiseError("field %s of %s is not a number: %s", field, nv.Name, nv.Value.Type())
	}
	return int(lx)
}

// TableField retrieves a *lua.LTable field from nv
func TableField(l *lua.LState, nv *NamedValue, field string) *lua.LTable {
	lvf := l.GetField(nv.Value, field)
	lt, ok := lvf.(*lua.LTable)
	if !ok {
		l.RaiseError("field %s of %s is not a table: %s", field, nv.Name, nv.Value.Type())
	}
	return lt
}

// OptStringField retrieves a string field from nv
func OptStringField(l *lua.LState, nv *NamedValue, field string, def string) string {
	lvf := l.GetField(nv.Value, field)
	ls, ok := lvf.(lua.LString)
	if !ok || lvf == lua.LNil {
		return def
	}
	return string(ls)
}

// OptBoolField retrieves a bool field from nv
func OptBoolField(l *lua.LState, nv *NamedValue, field string, def bool) bool {
	lvf := l.GetField(nv.Value, field)
	lb, ok := lvf.(lua.LBool)
	if !ok || lvf == lua.LNil {
		return def
	}
	return bool(lb)
}

// OptIntField retrieves an integer field from nv
func OptIntField(l *lua.LState, nv *NamedValue, field string, def int) int {
	lvf := l.GetField(nv.Value, field)
	lx, ok := lvf.(lua.LNumber)
	if !ok || lvf == lua.LNil {
		return def
	}
	return int(lx)
}

// OptTableField retrieves a *lua.LTable field from nv
func OptTableField(l *lua.LState, nv *NamedValue, field string, def *lua.LTable) *lua.LTable {
	lvf := l.GetField(nv.Value, field)
	lt, ok := lvf.(*lua.LTable)
	if !ok || lvf == lua.LNil {
		return def
	}
	return lt
}

// StringGlobal retrieves a string global value
func StringGlobal(l *lua.LState, global string) string {
	lvg := l.GetGlobal(global)
	ls, ok := lvg.(lua.LString)
	if !ok {
		l.RaiseError("global %s is not a string: %s", global, lvg.Type())
	}
	return string(ls)
}

// BoolGlobal retrieves a bool global value
func BoolGlobal(l *lua.LState, global string) bool {
	lvg := l.GetGlobal(global)
	lb, ok := lvg.(lua.LBool)
	if !ok {
		l.RaiseError("global %s is not a boolean: %s", global, lvg.Type())
	}
	return bool(lb)
}

// IntGlobal retrieves an integer global value
func IntGlobal(l *lua.LState, global string) int {
	lvg := l.GetGlobal(global)
	lx, ok := lvg.(lua.LNumber)
	if !ok {
		l.RaiseError("global %s is not a number: %s", global, lvg.Type())
	}
	return int(lx)
}

// TableGlobal retrieves a *lua.LTable global value
func TableGlobal(l *lua.LState, global string) *lua.LTable {
	lvg := l.GetGlobal(global)
	lt, ok := lvg.(*lua.LTable)
	if !ok {
		l.RaiseError("global %s is not a table: %s", global, lvg.Type())
	}
	return lt
}
