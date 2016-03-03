package path

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var Module = &luatest.Module{
	Name:       "path",
	Loader:     Loader,
	TestScript: "path_test.lua",
}

func TestModule(t *testing.T) {
	Module.Test(t)
}
