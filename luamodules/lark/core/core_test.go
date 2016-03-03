package core

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var Module = luatest.Module{
	Name:       "lark.core",
	Loader:     Loader,
	TestScript: "core_test.lua",
}

func TestModule(t *testing.T) {
	Module.Test(t)
}
