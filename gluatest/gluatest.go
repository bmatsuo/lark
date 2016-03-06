package gluatest

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

// FuncSetup is the name of a global variable containing a "setup" function.
// The setup function will be executed before each test function.  If any setup
// function fails the test function that would follow it will not execute.
var FuncSetup = "__test_setup"

// FuncTeardown is the name of a global variable containing a "teardown"
// function.  A teardown function will execute after each test function or
// after a setup function has failed.
var FuncTeardown = "__test_teardown"

// File is a test file for a lua module.
type File struct {
	Module gluamodule.Module
	Path   string
}

// preload runs the loader to register the module name.
func (m *File) preload(t testing.TB) (*lua.LState, *lua.LFunction) {
	l := lua.NewState()
	gluamodule.Preload(l, gluamodule.Resolve(m.Module)...)
	lfunc, err := l.LoadFile(m.Path)
	if err != nil {
		t.Fatal(err)
	}
	return l, lfunc
}

// Load runs the loader to register module and then executes the test file.
func (m *File) Load(t testing.TB) *lua.LState {
	l, lfunc := m.preload(t)
	l.Push(lfunc)
	err := l.PCall(0, 0, l.NewFunction(errTraceback))
	if err != nil {
		l.Close()
		t.Fatal(err)
	}
	return l
}

// BenchmarkRequireModule benchmarks the execution of the preload function (not
// the act of registering it).
func (m *File) BenchmarkRequireModule(b *testing.B) {
	b.StopTimer()
	for i := 0; i <= b.N; i++ {
		l, fn := m.preload(b)
		l.Push(fn)
		b.StartTimer()
		err := l.PCall(0, 0, nil)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test runs the specified test function
func (m *File) Test(t testing.TB) {
	testFuncs := m.getTestFuncs(t)
	for _, fname := range testFuncs {
		m.runTest(t, fname, getGlobalFunction(fname))
	}
}

func (m *File) getTestFuncs(t testing.TB) []string {
	l := m.Load(t)
	defer l.Close()

	return getTestFuncs(l)
}

type lfuncGetter func(*lua.LState) (*lua.LFunction, error)

func (m *File) runTest(t testing.TB, name string, getfunc lfuncGetter) {
	var fatal bool
	defer func() {
		if fatal {
			t.Fatal("FATAL")
		}
	}()

	l := m.Load(t)
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
	setup := l.GetGlobal(FuncSetup)
	teardown := l.GetGlobal(FuncTeardown)

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
