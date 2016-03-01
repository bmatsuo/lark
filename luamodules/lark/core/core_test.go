package core

import (
	"testing"

	"github.com/bmatsuo/lark/luatest"
)

var Module = luatest.Module{
	Name:       "lark.core",
	Loader:     Loader,
	TestScript: "core_test.lua",
}

func TestLog(t *testing.T) {
	Module.Test(t, "test_log")
}

func TestEnviron(t *testing.T) {
	Module.Test(t, "test_environ")
}

func TestExec(t *testing.T) {
	Module.Test(t, "test_exec")
}

func TestCapture(t *testing.T) {
	Module.Test(t, "test_capture")
}
