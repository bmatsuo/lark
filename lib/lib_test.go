package lib

import (
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/gluatest"
	"github.com/yuin/gopher-lua"
)

func BenchmarkRequireModule(b *testing.B) {
	testRequireAll.BenchmarkRequireModule(b)
}

var testRequireAll = &gluatest.File{
	Module: gluamodule.New("requireall", loaderRequireAll, Modules...),
}

func loaderRequireAll(l *lua.LState) int {
	require := l.GetGlobal("require")
	for _, m := range Modules {
		l.Push(require)
		l.Push(lua.LString(m.Name()))
		l.Call(1, 0)
	}
	return 0
}
