//go:generate ./build_lib.sh

package lark

import "github.com/yuin/gopher-lua"

// Loader loads the default lark module.
func Loader(l *lua.LState) int {
	fn, err := l.LoadString(LarkLib)
	if err != nil {
		l.RaiseError("%s", err)
	}
	l.Push(fn)
	err = l.PCall(0, 0, nil)
	if err != nil {
		l.RaiseError("%s", err)
	}
	return 0
}
