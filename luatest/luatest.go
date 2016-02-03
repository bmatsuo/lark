package luatest

import (
	"fmt"
	"testing"

	"github.com/yuin/gopher-lua"
)

// Module is a lua module to be tested.
type Module struct {
	Name       string
	Loader     lua.LGFunction
	TestScript string
}

// Preload runs the loader to register the module name.
func (m *Module) Preload(t testing.TB) *lua.LState {
	L := lua.NewState()
	L.PreloadModule(m.Name, m.Loader)
	err := L.DoFile(m.TestScript)
	if err != nil {
		t.Error(err)
	}
	return L
}

// Test runs the specified test function
func (m *Module) Test(t testing.TB, fn string) {
	L := m.Preload(t)
	defer L.Close()
	err := L.DoString(fmt.Sprintf("(%s)()", fn))
	if err != nil {
		t.Error(err)
	}
}
