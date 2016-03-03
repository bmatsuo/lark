package core

import (
	"testing"

	"github.com/bmatsuo/lark/internal/module"
	"github.com/bmatsuo/lark/luatest"
)

var Module = luatest.Module{
	Module:     module.New("lark.core", Loader),
	TestScript: "core_test.lua",
}

func TestModule(t *testing.T) {
	Module.Test(t)
}
