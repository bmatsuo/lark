package doc

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var testModule = &luatest.Module{
	Name:       "doc",
	Loader:     docLoader,
	TestScript: "doc_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
