package core

import (
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/luatest"
)

var testModule = luatest.Module{
	Module:     module.New("lark.core", Loader),
	TestScript: "core_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
