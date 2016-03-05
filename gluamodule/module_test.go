package gluamodule

import (
	"reflect"
	"sort"
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestNew(t *testing.T) {
	mname := "foo"
	m := New(mname, basicTestLoader)
	if m.Name() != mname {
		t.Errorf("name: %q (!= %q)", m.Name(), mname)
	}
}

func TestPreload(t *testing.T) {
	mname := "foo"
	m := New(mname, basicTestLoader)

	l := lua.NewState()
	defer l.Close()

	Preload(l, m)

	errh := l.NewFunction(func(l *lua.LState) int { return 1 })

	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString(mname))
	err := l.PCall(1, 1, errh)
	if err != nil {
		t.Error(err)
		return
	}
	l.SetGlobal(mname, l.Get(1))
	l.Pop(1)

	fn, err := l.LoadString(`return foo.double(3)`)
	if err != nil {
		t.Error(err)
		return
	}

	l.Push(fn)
	err = l.PCall(0, 0, errh)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestResolve(t *testing.T) {
	m1 := New("module1", basicTestLoader)
	_ = New("module2", basicTestLoader)
	m3 := New("module3", basicTestLoader, m1)
	m4 := New("module4", basicTestLoader, m3, m1)

	mods := Resolve(m4)
	var names []string
	for _, m := range mods {
		names = append(names, m.Name())
	}
	sort.Strings(names)

	expect := []string{
		m1.Name(),
		m3.Name(),
		m4.Name(),
	}
	if !reflect.DeepEqual(expect, names) {
		t.Errorf("modules: %q (!= %q)", names, expect)
	}
}

func basicTestLoader(l *lua.LState) int {
	mod := l.NewTable()

	double := l.NewClosure(func(l *lua.LState) int {
		x := l.CheckNumber(1)
		l.Push(lua.LNumber(2 * x))
		return 1
	}, mod)
	l.SetField(mod, "double", double)

	l.Push(mod)
	return 1
}
