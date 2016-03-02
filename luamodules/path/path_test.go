package path

import (
	"testing"

	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/bmatsuo/lark/luatest"
)

var Module = &luatest.Module{
	Name:       "path",
	Loader:     Loader,
	TestScript: "path_test.lua",
	PreloadDeps: []*luatest.Module{
		{Name: "doc", Loader: doc.Module().Loader},
	},
}

func TestDir(t *testing.T) {
	Module.Test(t, "test_dir")
}

func TestBase(t *testing.T) {
	Module.Test(t, "test_base")
}

func TestExt(t *testing.T) {
	Module.Test(t, "test_ext")
}

func TestExists(t *testing.T) {
	Module.Test(t, "test_exists")
}

func TestIsDir(t *testing.T) {
	Module.Test(t, "test_is_dir")
}

func TestGlob(t *testing.T) {
	Module.Test(t, "test_glob")
}

func TestJoin(t *testing.T) {
	Module.Test(t, "test_join")
}
