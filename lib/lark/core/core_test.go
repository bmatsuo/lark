package core

import (
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/gluatest"
)

var luaCoreTest = &gluatest.File{
	Module: gluamodule.New("lark.core", Loader),
	Path:   "core_test.lua",
}

func TestModule(t *testing.T) {
	luaCoreTest.Test(t)
}
