package path

import (
	"testing"

	"github.com/bmatsuo/lark/internal/module"
	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/bmatsuo/lark/luatest"
)

var Module = &luatest.Module{
	Module:     module.New("path", Loader),
	TestScript: "path_test.lua",
	PreloadDeps: []*luatest.Module{
		{Module: doc.Module()},
	},
}

func TestModule(t *testing.T) {
	Module.Test(t)
}
