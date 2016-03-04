package lib

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib/decorator"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/bmatsuo/lark/lib/lark"
)

func init() {
	Modules = joinModules(
		decorator.Module,
		doc.Module,
		lark.Module,
	)
}

// Modules contains all the modules defined in the library.
var Modules []gluamodule.Module

func joinModules(m ...gluamodule.Module) []gluamodule.Module {
	var mods []gluamodule.Module
	for _, m := range m {
		mods = append(mods, gluamodule.Collect(m)...)
	}
	return mods
}
