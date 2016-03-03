package doc

import (
	"testing"

	"github.com/bmatsuo/lark/luamodules/decorator"
	"github.com/bmatsuo/lark/luatest"
)

var testModule = &luatest.Module{
	Module:     Module,
	TestScript: "doc_test.lua",
	PreloadDeps: []*luatest.Module{
		{Module: decorator.Module},
	},
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
