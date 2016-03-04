package gluamodule

import (
	"sync"

	"github.com/yuin/gopher-lua"
)

// Preload preloads m into l.
func Preload(l *lua.LState, m ...Module) {
	for _, m := range m {
		l.PreloadModule(m.Name(), m.Loader)
	}
}

// Resolve recursively collects the dependencies of m and returns a slice of
// unique modules.
func Resolve(m Module) []Module {
	all := []Module{m}
	names := map[int]string{m.moduleID(): m.Name()}
	for i := 0; i < len(all); i++ {
		for _, m := range all[i].Deps() {
			if _, ok := names[m.moduleID()]; ok {
				continue
			}
			names[m.moduleID()] = m.Name()
			all = append(all, m)
		}
	}
	return all
}

// Module is a Lua module that can be loaded into a VM.
type Module interface {
	Name() string
	Loader(*lua.LState) int
	Deps() []Module
	moduleID() int
}

// New creates a new module.
func New(name string, loader lua.LGFunction, deps ...Module) Module {
	return defaultRegistry.NewModule(name, loader, copyModules(deps))
}

// Simple returns a simple module that only contains the functions provided by
// export.
func Simple(name string, export map[string]lua.LGFunction) Module {
	return defaultRegistry.NewModule(name, funcLoader(export), nil)
}

// Registry handles synchronized module construction.
type Registry interface {
	NewModule(name string, loader lua.LGFunction, deps func() []Module) Module
}

var defaultRegistry Registry = &registry{}

type registry struct {
	mut   sync.Mutex
	count int
}

var _ Registry = &registry{}

func (r *registry) next() int {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.count++
	return r.count
}

func (r *registry) NewModule(name string, loader lua.LGFunction, deps func() []Module) Module {
	return &simpleModule{
		id:     r.next(),
		name:   name,
		loader: loader,
		deps:   deps,
	}
}

type simpleModule struct {
	id     int
	name   string
	loader lua.LGFunction
	deps   func() []Module
}

var _ Module = (*simpleModule)(nil)

func (m *simpleModule) Name() string {
	return m.name
}

func (m *simpleModule) Loader(l *lua.LState) int {
	return m.loader(l)
}

func (m *simpleModule) Deps() []Module {
	return m.deps()
}

func (m *simpleModule) moduleID() int {
	return m.id
}

func copyModules(m []Module) func() []Module {
	return func() []Module {
		cp := make([]Module, len(m))
		copy(cp, m)
		return cp
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
