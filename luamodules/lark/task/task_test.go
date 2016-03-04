package task

import (
	"testing"

	"github.com/bmatsuo/lark/luamodules/decorator"
	"github.com/bmatsuo/lark/luamodules/decorator/intern"
	"github.com/bmatsuo/lark/luamodules/doc"
	"github.com/bmatsuo/lark/gluatest"
)

var testModule = &gluatest.Module{
	Module:     Module,
	TestScript: "task_test.lua",
	PreloadDeps: []*gluatest.Module{
		{Module: doc.Module},
		{Module: decorator.Module},
		{Module: intern.Module},
	},
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
