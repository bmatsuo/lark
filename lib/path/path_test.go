package path

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var luaPathTest = &gluatest.File{
	Module: Module,
	Path:   "path_test.lua",
}

func TestModule(t *testing.T) {
	luaPathTest.Test(t)
}

func BenchmarkRequireModule(b *testing.B) {
	luaPathTest.BenchmarkRequireModule(b)
}
