package luatest

import (
	"strings"
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

// Module is a lua module to be tested.
type Module struct {
	Module      module.Module
	TestScript  string
	PreloadDeps []*Module
}

// Preload runs the loader to register the module name.
func (m *Module) Preload(t testing.TB) *lua.LState {
	L := lua.NewState()
	module.Preload(L, m.Module)
	for _, m := range m.PreloadDeps {
		module.Preload(L, m.Module)
	}
	err := L.DoFile(m.TestScript)
	if err != nil {
		t.Error(err)
	}
	return L
}

// Test runs the specified test function
func (m *Module) Test(t testing.TB) {
	L := m.Preload(t)
	defer L.Close()

	testFuncs := getTestFuncs(L)
	errFn := L.NewFunction(errTraceback)
	for _, fname := range testFuncs {
		t.Logf("ENTER  %s", fname)
		L.Push(L.GetGlobal(fname))
		err := L.PCall(0, 0, errFn)
		if err != nil {
			t.Errorf("FAIL   %s\n%s", fname, err)
		} else {
			t.Logf("PASS   %s", fname)
		}
	}
}

func getTestFuncs(L *lua.LState) []string {
	return getGlobals(L, selTestFuncs)
}

func selTestFuncs(L *lua.LState, k, v lua.LValue) bool {
	s := string(k.(lua.LString))
	if !strings.HasPrefix(s, "test_") {
		return false
	}
	return true
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
	L.Push(lua.LNumber(2))
	L.Call(2, 1)
	return 1
}
