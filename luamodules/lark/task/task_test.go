package task

import (
	"testing"

	"github.com/bmatsuo/lark/luamodules/decorator"
	"github.com/bmatsuo/lark/luamodules/decorator/intern"
	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/bmatsuo/lark/luatest"
)

var testModule = &luatest.Module{
	Module:     Module,
	TestScript: "task_test.lua",
	PreloadDeps: []*luatest.Module{
		{Module: doc.Module},
		{Module: decorator.Module},
		{Module: intern.Module},
	},
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
