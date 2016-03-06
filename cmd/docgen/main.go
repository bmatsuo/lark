package main

import (
	"fmt"
	"log"

	"github.com/bmatsuo/lark/lib"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/bmatsuo/lark/project"
	"github.com/yuin/gopher-lua"
)

func main() {
	l := lua.NewState()
	defer l.Close()

	err := dumpModules(l)
	if err != nil {
		log.Fatal(err)
	}
}

func dumpModules(l *lua.LState) error {
	err := project.InitLib(l, nil)
	if err != nil {
		return err
	}

	for _, m := range lib.Modules {
		log.Print(m.Name())

		l.Push(l.GetGlobal("require"))
		l.Push(lua.LString(m.Name()))
		err := l.PCall(1, 1, nil)
		if err != nil {
			return fmt.Errorf("%s: %s", m.Name(), err)
		}

		mod := l.Get(-1)
		l.Pop(1)

		mdocs, err := doc.Get(l, mod)
		if err != nil {
			return fmt.Errorf("module %s: documentation error: %v", m.Name(), err)
		}

		log.Printf("%q", mdocs)
	}

	return nil
}
