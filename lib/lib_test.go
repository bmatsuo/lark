package lib

import (
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/gluatest"
	"github.com/yuin/gopher-lua"
)

func TestLib(t *testing.T) {
	if len(Modules) == 0 {
		t.Fatal("no modules")
	}

	l := lua.NewState()
	defer l.Close()

	Preload(l)

	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString("lark"))
	err := l.PCall(1, 0, nil)
	if err != nil {
		t.Error(err)
	}
}

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
