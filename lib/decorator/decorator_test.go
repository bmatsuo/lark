package decorator

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var luaDecoratorTest = &gluatest.File{
	Module: Module,
	Path:   "decorator_test.lua",
}

func TestModule(t *testing.T) {
	luaDecoratorTest.Test(t)
}

func BenchmarkRequireModule(b *testing.B) {
	luaDecoratorTest.BenchmarkRequireModule(b)
}
