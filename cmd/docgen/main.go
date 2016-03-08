// Command docgen is an experimental development command for generating static
// documentation for lua modules.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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
	var headers []*doc.Header
	var docs []*doc.Docs
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

		docs = append(docs, mdocs)
		headers = append(headers, &doc.Header{
			DocsType: "Module",
			Name:     m,
		})
	}

	gen := &mdGenerator{root: "docs", index: func(typ string) string { return "lua" }}
	hindex := &doc.Header{
		DocsType:    "index",
		DocsSubType: "Module",
	}
	path := gen.DocsPath(hindex)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	err = generateIndexFile(gen, path, hindex, headers, docs)
	if err != nil {
		return err
	}
	for i, h := range headers {
		d := docs[i]
		path := gen.DocsPath(h)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		err = generateDocsFile(gen, path, h, d)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateIndexFile(gen Generator, path string, h *doc.Header, headers []*doc.Header, docs []*doc.Docs) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = gen.GenerateIndex(f, h, headers, docs)
	if err != nil {
		return err
	}
	return f.Close()
}

func generateDocsFile(gen Generator, path string, h *doc.Header, d *doc.Docs) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = gen.GenerateDocs(f, h, d)
	if err != nil {
		return err
	}
	return f.Close()
}

// Generator formats doc.Docs objects and writes them to an output stream.
type Generator interface {
	FileFormat() string
	DocsPath(h *doc.Header) string
	GenerateDocs(out io.Writer, h *doc.Header, d *doc.Docs) error
	GenerateIndex(out io.Writer, h *doc.Header, headers []*doc.Header, docs []*doc.Docs) error
}

type textGenerator struct {
	root string
}

func (g *textGenerator) FileFormat() string {
	return "text/plain"
}

func (g *textGenerator) DocsPath(h *doc.Header) string {
	docstype := strings.ToLower(h.DocsType)
	subtype := strings.ToLower(h.DocsSubType)
	if docstype == "index" {
		if subtype != "" {
			return subtype + ".txt"
		}
		return "index.txt"
	}
	path := filepath.Join(strings.Split(h.Name, ".")...)
	dir := "modules"
	if docstype == "object" {
		dir = "objects"
	}
	path = filepath.Join(g.root, dir, path)
	path += ".txt"
	return path
}

func (g *textGenerator) GenerateDocs(out io.Writer, h *doc.Header, d *doc.Docs) error {
	s, err := doc.NewFormatter().Format(h, d, "  ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, s)
	return err
}

func (g *textGenerator) GenerateIndex(out io.Writer, h *doc.Header, headers []*doc.Header, docs []*doc.Docs) error {
	w := &failWriter{Writer: out}
	printf := func(format string, v ...interface{}) {
		fmt.Fprintf(w, format, v...)
	}
	if len(headers) != len(docs) {
		return fmt.Errorf("unequal number of headers and docs")
	}

	printf("%s index\n\n", h.DocsSubType)

	for i, h := range headers {
		d := docs[i]
		printf("  %s\n", h.Name)
		if d == nil {
			printf("\n")
		} else {
			text := textutil.Synopsis(d.Desc)
			if text == "" {
				printf("\n")
			} else {
				text = textutil.Wrap(text, 66)
				text = textutil.Indent(text, "      ")
				printf("%s\n\n", text)
			}
		}
	}

	return w.err
}

type mdGenerator struct {
	root  string
	index func(string) string
}

func (g *mdGenerator) FileFormat() string {
	return "text/markdown"
}

func (g *mdGenerator) DocsPath(h *doc.Header) string {
	path := g.path(h)
	if path == "" {
		return ""
	}
	return filepath.Join(g.root, path)
}

func (g *mdGenerator) path(h *doc.Header) string {
	docstype := strings.ToLower(h.DocsType)
	subtype := strings.ToLower(h.DocsSubType)
	if docstype == "index" {
		if g.index != nil {
			name := g.index(subtype)
			name += ".md"
			return name
		}
		filename := "index.md"
		if subtype != "" {
			return subtype + ".md"
		}
		return filename
	}
	path := filepath.Join(strings.Split(h.Name, ".")...)
	dir := "modules"
	if docstype == "object" {
		dir = "objects"
	}
	path = filepath.Join(dir, path)
	path += ".md"
	return path

}

func (g *mdGenerator) GenerateDocs(out io.Writer, h *doc.Header, d *doc.Docs) error {
	s, err := newMarkdownFormatter().Format(h, d, "  ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, s)
	return err
}

func (g *mdGenerator) GenerateIndex(out io.Writer, h *doc.Header, headers []*doc.Header, docs []*doc.Docs) error {
	w := &failWriter{Writer: out}
	printf := func(format string, v ...interface{}) {
		fmt.Fprintf(w, format, v...)
	}
	if len(headers) != len(docs) {
		return fmt.Errorf("unequal number of headers and docs")
	}

	printf("#%s index\n\n", h.DocsSubType)

	for i, h := range headers {
		d := docs[i]

		mdpath := g.path(h)
		printf("##[%s](%s)\n\n", h.Name, mdpath)
		if d == nil {
			continue
		}
		text := textutil.Synopsis(d.Desc)
		if text != "" {
			text = textutil.Wrap(text, 72)
			printf("%s\n\n", text)
		}
	}

	return w.err
}

type failWriter struct {
	io.Writer
	err error
}

func (w *failWriter) Write(b []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.Writer.Write(b)
	w.err = err
	return n, err
}

type mdFormatter struct {
	buf *bytes.Buffer
}

func newMarkdownFormatter() *mdFormatter {
	return &mdFormatter{}
}

func (g *mdFormatter) Format(h *doc.Header, d *doc.Docs, indentdocs string) (string, error) {
	if g.buf == nil {
		g.buf = new(bytes.Buffer)
	} else {
		g.buf.Reset()
	}
	g.writeDocs(h, d)
	return g.buf.String(), nil
}

func (g *mdFormatter) printf(format string, v ...interface{}) {
	fmt.Fprintf(g.buf, format, v...)
}

func (g *mdFormatter) writeP(text string, wrap int) {
	text = textutil.Unindent(text)
	text = strings.TrimSpace(text)
	ps := strings.Split(text, "\n\n")
	for _, p := range ps {
		if p == "" {
			continue
		}
		space := strings.IndexFunc(p, unicode.IsSpace)
		if space == 0 {
			p = textutil.Unindent(p)
			p = textutil.Indent(p, "    ")
		}
		text = textutil.Wrap(p, wrap)
		g.printf("%s\n\n", p)
	}
}

func (g *mdFormatter) writeDocs(h *doc.Header, d *doc.Docs) {
	sec := ""
	if h != nil {
		if h.DocsType == "Module" {
			fmt.Fprintf(g.buf, "#%s %s\n\n", h.DocsType, h.Name)
		} else {
			sec = "#"
			fmt.Fprintf(g.buf, "%s#%s %s\n\n", sec, h.DocsType, h.Name)
		}
	}
	if d == nil {
		return
	}
	text := d.Usage
	text = textutil.Unindent(text)
	text = strings.TrimSpace(text)
	if text != "" {
		text = textutil.Indent(text, "    ")
		g.printf("%s\n\n", text)
	}
	text = d.Sig
	text = textutil.Unindent(text)
	text = strings.TrimSpace(text)
	if text != "" {
		g.printf("%s##Signature\n\n", sec)
		g.printf("%s\n\n", text)
	}
	text = textutil.Unindent(d.Desc)
	text = strings.TrimSpace(text)
	if text != "" {
		g.printf("%s##Description\n\n", sec)
		g.writeP(text, 72)
	}
	numvar := d.NumVar()
	if numvar > 0 {
		g.printf("%s##Variables\n\n", sec)
		for i := 0; i < numvar; i++ {
			g.printf("**%s**", d.Var(i))
			typ := d.VarType(i)
			if typ != "" {
				g.printf(" _%s_\n\n", typ)
			} else {
				g.printf("\n\n")
			}
			text = d.VarDesc(i)
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			if text != "" {
				g.writeP(text, 72)
			}
		}
	}
	numparam := d.NumParam()
	if numparam > 0 {
		g.printf("%s##Parameters\n\n", sec)
		for i := 0; i < numparam; i++ {
			g.printf("**%s**", d.Param(i))
			typ := d.ParamType(i)
			if typ != "" {
				g.printf(" _%s_\n\n", typ)
			} else {
				g.printf("\n\n")
			}
			text = d.ParamDesc(i)
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			if text != "" {
				g.writeP(text, 72)
			}
		}
	}

	funcs := d.Funcs()
	if len(funcs) > 0 {
		g.printf("%s##Functions\n\n", sec)
		for _, sub := range funcs {
			if sub.Name == "" {
				continue
			}
			if sec == "" {
				g.printf("**[%s](#function-%s)**\n\n", sub.Name, strings.Replace(sub.Name, ".", "", -1))
			} else {
			}
			if sub.Docs == nil {
				continue
			}
			text = sub.Desc
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			text = textutil.Synopsis(text)
			text = textutil.Wrap(text, 72)
			if text != "" {
				g.printf("%s\n\n", text)
			}
		}
	}

	others := d.Others()
	if len(others) > 0 {
		g.printf("%s##Subtopics\n\n", sec)
		for _, sub := range others {
			if sub.Name == "" {
				continue
			}
			g.printf("**%s**\n\n", sub.Name)
			if sub.Docs == nil {
				continue
			}
			text = sub.Desc
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			text = textutil.Synopsis(text)
			text = textutil.Wrap(text, 72)
			if text != "" {
				g.printf("%s\n\n", text)
			}
		}
	}

	for _, sub := range others {
		hsub := &doc.Header{}
		*hsub = *h
		hsub.DocsType = "Object"
		hsub.Name += "." + sub.Name
		g.writeDocs(hsub, sub.Docs)
	}

	for _, sub := range funcs {
		hsub := &doc.Header{}
		*hsub = *h
		hsub.DocsType = "Function"
		hsub.Name += "." + sub.Name
		g.writeDocs(hsub, sub.Docs)
	}
}
