package fun

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var luaFunTest = &gluatest.File{
	Module: Module,
	Path:   "fun_test.lua",
}

func TestFun(t *testing.T) {
	luaFunTest.Test(t)
}

func BenchmarkRequireModule(b *testing.B) {
	file := &gluatest.File{Module: Module}
	file.BenchmarkRequireModule(b)
}
