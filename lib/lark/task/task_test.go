package task

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var luaTaskTest = &gluatest.File{
	Module: Module,
	Path:   "task_test.lua",
}

func TestModule(t *testing.T) {
	luaTaskTest.Test(t)
}

func BenchmarkRequireModule(b *testing.B) {
	luaTaskTest.BenchmarkRequireModule(b)
}
