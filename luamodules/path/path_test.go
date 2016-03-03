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

func TestModule(t *testing.T) {
	Module.Test(t)
}
