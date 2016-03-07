// DO NOT EDIT
// THIS IS A GENERATED FILE

package lib

import (
	"github.com/bmatsuo/lark/gluamodule"
	"github.com/bmatsuo/lark/lib/decorator"
	"github.com/bmatsuo/lark/lib/decorator/_intern"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/bmatsuo/lark/lib/lark"
	"github.com/bmatsuo/lark/lib/lark/core"
	"github.com/bmatsuo/lark/lib/lark/task"
	"github.com/bmatsuo/lark/lib/path"
)

// Modules lists every module in the library.
var Modules = []gluamodule.Module{
	decorator.Module,
	intern.Module,
	doc.Module,
	lark.Module,
	core.Module,
	task.Module,
	path.Module,
}

// InteralModules modules that are not general purpose and should not be imported by scripts.
var InternalModules = []gluamodule.Module{
	intern.Module,
}
