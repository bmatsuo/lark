package lark

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var luaLarkTest = &gluatest.File{
	Module: Module,
	Path:   "lark_test.lua",
}

func TestLark(t *testing.T) {
	luaLarkTest.Test(t)
}
