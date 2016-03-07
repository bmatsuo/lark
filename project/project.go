package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var pathSeparator = string([]rune{filepath.Separator})

// LarkFile is contains possible names for the primary (highest
// priority) lark file.
var LarkFile = []string{
	"lark.lua",
	"Larkfile",
}

// TaskDir is the auxiliary task directory in a lark project.
var TaskDir = "lark_tasks"

// ModuleDir is the root directory for modules in a lark project.
var ModuleDir = "lark_modules"

// LuaPath returns the dir project LUA_PATH value, referencing only ModuleDir
// inside dir.
func LuaPath(dir string) string {
	root := filepath.Join(dir, ModuleDir)
	luaFiles := filepath.Join(root, "?.lua")
	luaInits := filepath.Join(root, "?", "init.lua")
	return fmt.Sprintf("%s;%s", luaFiles, luaInits)
}

// FindTaskFiles locates task scripts in the project dir.
func FindTaskFiles(dir string) ([]string, error) {
	var luaFiles []string
	join := filepath.Join

	for _, possible := range LarkFile {
		path := join(dir, possible)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %s", possible, err)
		}
		luaFiles = append(luaFiles, path)
		break
	}

	subpatt := join(dir, TaskDir, "*.lua")
	files, err := filepath.Glob(subpatt)
	luaFiles = append(luaFiles, files...)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", TaskDir, err)
	}

	return luaFiles, nil
}

// FindModules locates modules in the ModuleDir of project dir.
// The modules names returned by FindModules match what would be
// passed to Lua's require() function.
//
// BUG:
// Handling of symbolic links is undefined.
func FindModules(dir string) ([]string, error) {
	var modules []string
	root := filepath.Join(dir, ModuleDir)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".lua" {
			return nil
		}
		relpath, err := filepath.Rel(root, path)
		if err != nil {
			// really unexpected
			return err
		}

		if filepath.Base(path) == "init.lua" {
			if filepath.Dir(path) == root {
				return fmt.Errorf("directory %s cannot be a module: %s", ModuleDir, path)
			}

			mpath := filepath.Dir(relpath)
			m := strings.Replace(mpath, pathSeparator, ".", -1)
			modules = append(modules, m)

			return nil
		}

		mpath := relpath[:len(relpath)-4] // trim .lua extension
		m := strings.Replace(mpath, pathSeparator, ".", -1)
		modules = append(modules, m)

		return nil
	})
	return nil, nil
}
