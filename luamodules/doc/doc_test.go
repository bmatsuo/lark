package doc

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var TestModule = &luatest.Module{
	Name:       "doc",
	Loader:     (&doc{}).LoaderNative,
	TestScript: "doc_test.lua",
}

func TestDoc(t *testing.T) {
	TestModule.Test(t, "test_doc")
}
