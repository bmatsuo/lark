package decorator

import (
	"testing"

	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/bmatsuo/lark/luatest"
)

var testModule = &luatest.Module{
	Module:     Module,
	TestScript: "decorator_test.lua",
	PreloadDeps: []*luatest.Module{
		{Module: doc.Module},
	},
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
