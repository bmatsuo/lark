package task

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var testModule = &gluatest.Module{
	Module:     Module,
	TestScript: "task_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
