package doc

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var luaDocTest = &gluatest.File{
	Module: Module,
	Path:   "doc_test.lua",
}

func TestModule(t *testing.T) {
	luaDocTest.Test(t)
}

func BenchmarkRequireModule(b *testing.B) {
	luaDocTest.BenchmarkRequireModule(b)
}
