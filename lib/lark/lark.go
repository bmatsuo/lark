//go:generate ./build_lib.sh

package lark

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

// Module is a gluamodule.Module that loads the lark module.
var Module = gluamodule.New("lark", Loader)

// Loader loads the default lark module.
func Loader(l *lua.LState) int {
	fn, err := l.LoadString(LarkLib)
	if err != nil {
		l.RaiseError("%s", err)
	}
	l.Push(fn)
	l.Call(0, 1)
	return 1
}
