package gluamodule

import "github.com/yuin/gopher-lua"

// Preload preloads m into l.
func Preload(l *lua.LState, m Module) {
	l.PreloadModule(m.Name(), m.Loader)
}

// Collect recursively collects submodules of m
func Collect(m Module) []Module {
	mods := []Module{m}
	for i := 0; i < len(mods); i++ {
		mods = append(mods, m.Submodules()...)
	}
	return mods
}

// Module is a Lua module that can be loaded into a VM.
type Module interface {
	Name() string
	Loader(*lua.LState) int
	Submodules() []Module
}

// New creates a new module.
func New(name string, loader lua.LGFunction) Module {
	return &simpleModule{name, loader, nil}
}

// NewSub creates a new module with submodules.
func NewSub(name string, loader lua.LGFunction, sub ...Module) Module {
	return &simpleModule{name, loader, copyModules(sub)}
}

// Simple returns a simple module that only contains the functions provided by
// export.
func Simple(name string, export map[string]lua.LGFunction) Module {
	return &simpleModule{name, funcLoader(export), nil}
}

type simpleModule struct {
	name   string
	loader lua.LGFunction
	sub    func() []Module
}

var _ Module = (*simpleModule)(nil)

func (m *simpleModule) Name() string {
	return m.name
}

func (m *simpleModule) Loader(l *lua.LState) int {
	return m.loader(l)
}

func (m *simpleModule) Submodules() []Module {
	if m.sub != nil {
		return m.sub()
	}
	return nil
}

func copyModules(subs []Module) func() []Module {
	return func() []Module {
		_subs := make([]Module, len(subs))
		copy(_subs, subs)
		return _subs
	}
}

func funcLoader(export map[string]lua.LGFunction) lua.LGFunction {
	return func(l *lua.LState) int {
		t := l.NewTable()
		mod := l.SetFuncs(t, export)
		l.Push(mod)
		return 1
	}
}
