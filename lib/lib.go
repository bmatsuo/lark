//go:generate bash collect.sh modules.go

package lib

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

// Preload preloads Modules into l.
func Preload(l *lua.LState) {
	gluamodule.Preload(l, Modules...)
}
