package gluatest

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

// Module is a lua module to be tested.
type Module struct {
	Module     gluamodule.Module
	TestScript string
}

// Preload runs the loader to register the module name.
func (m *Module) Preload(t testing.TB) *lua.LState {
	L := lua.NewState()
	gluamodule.Preload(L, gluamodule.Resolve(m.Module)...)
	err := L.DoFile(m.TestScript)
	if err != nil {
		t.Error(err)
	}
	return L
}

// Test runs the specified test function
func (m *Module) Test(t testing.TB) {
	testFuncs := m.getTestFuncs(t)
	for _, fname := range testFuncs {
		m.runTest(t, fname, getGlobalFunction(fname))
	}
}

func (m *Module) getTestFuncs(t testing.TB) []string {
	l := m.Preload(t)
	defer l.Close()

	return getTestFuncs(l)
}

type lfuncGetter func(*lua.LState) (*lua.LFunction, error)

func (m *Module) runTest(t testing.TB, name string, getfunc lfuncGetter) {
	var fatal bool
	defer func() {
		if fatal {
			t.Fatal("FATAL")
		}
	}()

	l := m.Preload(t)
	defer l.Close()

	failmsg := func(name string, err error) string {
		return fmt.Sprintf("FAIL   %s\n%s", name, err)
	}

	errFn := l.NewFunction(errTraceback)

	lfunc, err := getfunc(l)
	if err != nil {
		t.Errorf("FAIL\n%s", err)
		return
	}
	setup := l.GetGlobal("__test_setup")
	teardown := l.GetGlobal("__test_teardown")

	if setup != lua.LNil {
		l.Push(setup)
		err := l.PCall(0, 0, errFn)
		if err != nil {
			fatal = true
			t.Error(failmsg("SETUP", err))
		}
	}

	if !fatal {
		t.Logf("ENTER  %s", name)
		l.Push(lfunc)
		err = l.PCall(0, 0, errFn)
		if err != nil {
			t.Error(failmsg(name, err))
		} else {
			t.Logf("PASS   %s", name)
		}
	}

	if teardown != lua.LNil {
		l.Push(teardown)
		err := l.PCall(0, 0, errFn)
		if err != nil {
			fatal = true
			t.Error(failmsg("TEARDOWN", err))
			return
		}
	}
}

func getTestFuncs(L *lua.LState) []string {
	funcs := getGlobals(L, selTestFuncs)
	sort.Strings(funcs)
	return funcs
}

func selTestFuncs(L *lua.LState, k, v lua.LValue) bool {
	s := string(k.(lua.LString))
	if !strings.HasPrefix(s, "test_") {
		return false
	}
	return true
}

func getGlobalFunction(fname string) lfuncGetter {
	return func(l *lua.LState) (*lua.LFunction, error) {
		lfunc, ok := l.GetGlobal(fname).(*lua.LFunction)
		if !ok {
			return nil, fmt.Errorf("not a function: %s", lfunc)
		}
		return lfunc, nil
	}
}

func getGlobals(L *lua.LState, sel func(L *lua.LState, k, v lua.LValue) bool) []string {
	var sels []string
	globals := L.Get(lua.GlobalsIndex).(*lua.LTable)
	L.ForEach(globals, func(k, v lua.LValue) {
		s, ok := k.(lua.LString)
		if !ok {
			return
		}
		if sel(L, k, v) {
			sels = append(sels, string(s))
		}
	})
	return sels
}

func errTraceback(L *lua.LState) int {
	msg := L.Get(1)
	L.SetTop(0)
	L.Push(L.GetField(L.GetGlobal("debug"), "traceback"))
	L.Push(msg)
	L.Push(lua.LNumber(3))
	L.Call(2, 1)
	return 1
}
