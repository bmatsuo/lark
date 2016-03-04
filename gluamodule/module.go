package gluamodule

import "github.com/yuin/gopher-lua"

// Preload preloads m into l.
func Preload(l *lua.LState, m Module) {
	l.PreloadModule(m.Name(), m.Loader)
}

// Module is a Lua module that can be loaded into a VM.
type Module interface {
	Name() string
	Loader(*lua.LState) int
}

// New creates a new module.
func New(name string, loader lua.LGFunction) Module {
	return &simpleModule{name, loader}
}

// Simple returns a simple module that only contains the functions provided by
// export.
func Simple(name string, export map[string]lua.LGFunction) Module {
	return &simpleModule{name, funcLoader(export)}
}

type simpleModule struct {
	name   string
	loader lua.LGFunction
}

var _ Module = (*simpleModule)(nil)

func (m *simpleModule) Name() string {
	return m.name
}

func (m *simpleModule) Loader(l *lua.LState) int {
	return m.loader(l)
}

func funcLoader(export map[string]lua.LGFunction) lua.LGFunction {
	return func(l *lua.LState) int {
		t := l.NewTable()
		mod := l.SetFuncs(t, export)
		l.Push(mod)
		return 1
	}
}
