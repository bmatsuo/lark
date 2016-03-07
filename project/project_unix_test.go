package project

import "testing"

func TestPackagePath(t *testing.T) {
	path := PackagePath(".")
	expect := "lark_modules/?.lua;lark_modules/?/init.lua"
	if path != expect {
		t.Errorf("package path: %q (!= %q)", path, expect)
	}
}
