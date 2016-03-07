// Command docgen is an experimental development command for generating static
// documentation for lua modules.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bmatsuo/lark/internal/textutil"
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
	gen := &bufferedTextGenerator{}
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

		header := &DocsHeader{
			DocsType: "Module",
			Name:     m,
			Usage:    fmt.Sprintf("local %s = require(%q)", m, m),
		}
		gen.GenerateDocs(out, header, mdocs)
	}
	return nil
}

// DocsHeader describes a page of documentation
type DocsHeader struct {
	DocsType string
	Name     string
	Usage    string
}

// Generator formats doc.Docs objects and writes them to an output stream.
type Generator interface {
	GenerateDocs(out io.Writer, h *DocsHeader, d *doc.Docs) error
}

type bufferedTextGenerator textGenerator

func (g *bufferedTextGenerator) GenerateDocs(out io.Writer, h *DocsHeader, d *doc.Docs) (err error) {
	_out := bufio.NewWriter(out)
	defer func() {
		if err != nil {
			err = _out.Flush()
		}
	}()
	return (*textGenerator)(g).GenerateDocs(out, h, d)
}

type textGenerator struct {
}

func (g *textGenerator) GenerateDocs(out io.Writer, h *DocsHeader, d *doc.Docs) error {
	var err error
	printf := func(format string, v ...interface{}) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(out, format, v...)
	}

	printf("%s %s\n\n", h.DocsType, h.Name)
	if d.Sig != "" {
		text := d.Sig
		text = textutil.Unindent(text)
		text = strings.TrimSpace(text)
		text = textutil.Indent(text, "  ")
		printf("Signature:\n\n%s\n\n", d.Sig)
	}
	if h.Usage != "" {
		printf("  %s\n\n", h.Usage)
	}
	if d.Desc != "" {
		text := textutil.Unindent(d.Desc)
		text = textutil.Wrap(text, 72)
		text = strings.TrimSpace(text)
		printf("%s\n\n", text)
	}
	numvar := d.NumVar()
	if numvar > 0 {
		printf("Variables:\n\n")
		for i := 0; i < numvar; i++ {
			printf("  %s", d.Var(i))
			typ := d.VarType(i)
			if typ != "" {
				printf("  %s\n", typ)
			} else {
				printf("\n")
			}
			text := d.VarDesc(i)
			text = textutil.Unindent(text)
			text = textutil.Wrap(text, 66)
			text = strings.TrimSpace(text)
			if text != "" {
				text = textutil.Indent(text, "      ")
				printf("%s\n\n", text)
			} else {
				printf("\n")
			}
		}
	}
	numparam := d.NumParam()
	if numparam > 0 {
		printf("Parameters:\n\n")
		for i := 0; i < numparam; i++ {
			printf("  %s", d.Param(i))
			typ := d.ParamType(i)
			if typ != "" {
				printf("  %s\n", typ)
			} else {
				printf("\n")
			}
			text := d.ParamDesc(i)
			text = textutil.Unindent(text)
			text = textutil.Wrap(text, 66)
			text = strings.TrimSpace(text)
			if text != "" {
				text = textutil.Indent(text, "      ")
				printf("%s\n\n", text)
			} else {
				printf("\n")
			}
		}
	}

	funcs := d.Funcs()
	if len(funcs) > 0 {
		printf("Functions:\n\n")
		for _, sub := range funcs {
			if sub.Name == "" {
				continue
			}
			printf("  %s\n", sub.Name)
			if sub.Docs == nil {
				printf("\n")
				continue
			}
			text := sub.Desc
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			text = textutil.Synopsis(text)
			text = textutil.Wrap(text, 66)
			if text != "" {
				text = textutil.Indent(text, "      ")
				printf("%s\n\n", text)
			} else {
				printf("\n")
			}
		}
	}

	others := d.Others()
	if len(others) > 0 {
		printf("Subtopics:\n\n")
		for _, sub := range others {
			if sub.Name == "" {
				continue
			}
			printf("  %s\n", sub.Name)
			if sub.Docs == nil {
				printf("\n")
				continue
			}
			text := sub.Desc
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			text = textutil.Synopsis(text)
			text = textutil.Wrap(text, 66)
			if text != "" {
				text = textutil.Indent(text, "      ")
				printf("%s\n\n", text)
			} else {
				printf("\n")
			}
		}
	}

	return err
}
