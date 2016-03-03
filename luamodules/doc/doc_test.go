package doc

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var testModule = &luatest.Module{
	Module:     Module(),
	TestScript: "doc_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
