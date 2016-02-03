package path

import (
	"fmt"
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestDir(t *testing.T) {
	luaModuleTest(t, "test_dir")
}

func TestBase(t *testing.T) {
	luaModuleTest(t, "test_base")
}

func TestExt(t *testing.T) {
	luaModuleTest(t, "test_ext")
}

func TestExists(t *testing.T) {
	luaModuleTest(t, "test_exists")
}

func TestIsDir(t *testing.T) {
	luaModuleTest(t, "test_is_dir")
}

func TestGlob(t *testing.T) {
	luaModuleTest(t, "test_glob")
}

func TestJoin(t *testing.T) {
	luaModuleTest(t, "test_join")
}

var luaTestFile = "path_test.lua"

func loadLuaModule(t *testing.T) *lua.LState {
	L := lua.NewState()
	L.PreloadModule("path", ModuleLoader)
	err := L.DoFile(luaTestFile)
	if err != nil {
		t.Error(err)
	}
	return L
}

func luaModuleTest(t *testing.T, fn string) {
	L := loadLuaModule(t)
	defer L.Close()
	err := L.DoString(fmt.Sprintf("(%s)()", fn))
	if err != nil {
		t.Error(err)
	}
}
