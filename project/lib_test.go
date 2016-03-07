package project

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestInitLib(t *testing.T) {
	l := lua.NewState()
	defer l.Close()
	err := InitLib(l, nil)
	if err != nil {
		t.Error(err)
	}
	lvlark := l.GetGlobal("lark")
	if lvlark != lua.LNil {
		t.Errorf("lark was loaded")
	}
	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString("lark"))
	err = l.PCall(1, 1, nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, ok := l.Get(-1).(*lua.LTable)
	if !ok {
		t.Errorf("require(\"lark\") did not return module table: %s", l.Get(-1).Type())
	}
}

func TestInitLib_load(t *testing.T) {
	l := lua.NewState()
	defer l.Close()
	err := InitLib(l, &Init{
		Load:    true,
		Verbose: true,
	})
	if err != nil {
		t.Error(err)
	}
	lvlark := l.GetGlobal("lark")
	if lvlark == lua.LNil {
		t.Errorf("lark was not loaded")
	}
	_, ok := lvlark.(*lua.LTable)
	if !ok {
		t.Errorf("lark is not a table: %s", lvlark.Type())
	}
}
