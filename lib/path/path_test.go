package path

import (
	"testing"

	"github.com/bmatsuo/lark/gluatest"
)

var testModule = &gluatest.Module{
	Module:     Module,
	TestScript: "path_test.lua",
}

func TestModule(t *testing.T) {
	testModule.Test(t)
}
