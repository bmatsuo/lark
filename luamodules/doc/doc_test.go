package doc

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var _Module = &luatest.Module{
	Name:       "doc",
	Loader:     (&doc{}).LoaderNative,
	TestScript: "doc_test.lua",
}

func TestModule(t *testing.T) {
	_Module.Test(t)
}
