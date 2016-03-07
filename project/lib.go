package project

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib"
	"github.com/yuin/gopher-lua"
)

// Init configure how the library is initialized.
type Init struct {
	Load    bool
	Verbose bool
}

// InitLib preloads the lark library into l
func InitLib(l *lua.LState, init *Init) error {
	gluamodule.Preload(l, lib.Modules...)
	if init == nil || !init.Load {
		return nil
	}
	if init == nil {
		init = &Init{}
	}

	l.Push(l.GetGlobal("require"))
	l.Push(lua.LString("lark"))
	err := l.PCall(1, 1, nil)
	if err != nil {
		return err
	}
	lark := l.Get(-1)
	l.Pop(1)
	l.SetGlobal("lark", lark)

	if init.Verbose {
		l.SetField(lark, "verbose", lua.LBool(true))
	}
	return nil
}
