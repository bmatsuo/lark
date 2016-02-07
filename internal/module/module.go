package module

import "github.com/yuin/gopher-lua"

// Module is a Lua module that can be loaded into a VM.
type Module interface {
	Loader(*lua.LState) int
}

// Exporter returns the API for a Lua module.
type Exporter interface {
	Export() map[string]lua.LGFunction
}

// Simple returns a simple module that only contains the functions provided by
// export.
func Simple(export Exporter) Module {
	return &simpleModule{export}
}

type simpleModule struct {
	export Exporter
}

func (m *simpleModule) Loader(l *lua.LState) int {
	t := l.NewTable()
	mod := l.SetFuncs(t, m.exports.Export())
	l.Push(mod)
	return 1
}
