package core

import (
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/gluatest"
)

var testModule = gluatest.Module{
	Module:     module.New("lark.core", Loader),
	TestScript: "core_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
