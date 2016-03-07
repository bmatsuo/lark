package project

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestPackagePath(t *testing.T) {
	path := PackagePath(".")
	expect := "lark_modules/?.lua;lark_modules/?/init.lua"
	if path != expect {
		t.Errorf("package path: %q (!= %q)", path, expect)
	}
}

func TestSetPackagePath(t *testing.T) {
	dir := "x"
	expect := "x/lark_modules/?.lua;x/lark_modules/?/init.lua"

	l := lua.NewState()
	defer l.Close()

	err := SetPackagePath(l, dir)
	if err != nil {
		t.Error(err)
	}

	lvpath := l.GetField(l.GetGlobal("package"), "path")
	path, ok := lvpath.(lua.LString)
	if !ok {
		t.Errorf("package.path is not a string: %s", lvpath.Type())
	}

	if string(path) != expect {
		t.Errorf("package path: %q (!= %q)", path, expect)
	}
}
