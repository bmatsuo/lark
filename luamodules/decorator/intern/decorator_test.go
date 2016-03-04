package intern

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var testModule = &luatest.Module{
	Module:     Module,
	TestScript: "decorator_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
