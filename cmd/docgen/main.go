// Command docgen is an experimental development command for generating static
// documentation for lua modules.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bmatsuo/lark/lib"
	"github.com/bmatsuo/lark/lib/doc"
	"github.com/bmatsuo/lark/project"
	"github.com/yuin/gopher-lua"
)

func main() {
	l := lua.NewState()
	defer l.Close()

	err := dump(l)
	if err != nil {
		log.Fatal(err)
	}
}

func dump(l *lua.LState) error {
	err := project.InitLib(l, nil)
	if err != nil {
		return err
	}

	var modules []string
	if len(os.Args[1:]) > 0 {
		modules = os.Args[1:]
	} else {
	mloop:
		for _, m := range lib.Modules {
			name := m.Name()
			for _, hidden := range lib.InternalModules {
				if name == hidden.Name() {
					continue mloop
				}
			}
			modules = append(modules, name)
		}
	}

	return dumpDocs(l, modules)
}

func dumpDocs(l *lua.LState, names []string) error {
	gen := &textGenerator{}
	out := os.Stdout
	for _, m := range names {
		l.Push(l.GetGlobal("require"))
		l.Push(lua.LString(m))
		err := l.PCall(1, 1, nil)
		if err != nil {
			return fmt.Errorf("%s: %s", m, err)
		}

		mod := l.Get(-1)
		l.Pop(1)

		mdocs, err := doc.Get(l, mod, m)
		if err != nil {
			return fmt.Errorf("module %s: documentation error: %v", m, err)
		}

		header := &doc.Header{
			DocsType: "Module",
			Name:     m,
		}
		gen.GenerateDocs(out, header, mdocs)
	}
	return nil
}

// Generator formats doc.Docs objects and writes them to an output stream.
type Generator interface {
	GenerateDocs(out io.Writer, h *doc.Header, d *doc.Docs) error
}

type textGenerator struct {
}

func (g *textGenerator) GenerateDocs(out io.Writer, h *doc.Header, d *doc.Docs) error {
	s, err := doc.NewFormatter().Format(h, d, "  ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, s)
	return err
}
