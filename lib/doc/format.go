package doc

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bmatsuo/lark/internal/textutil"
)

// Header provides a header for Docs if the information is known.  Acceptable
// DocsType values are "module", "object", and "index".  Acceptable DocsSubType
// values are "module" and "object".  DocsType and DocsSubType are case
// insensitive.
type Header struct {
	DocsType    string
	DocsSubType string
	Name        string
}

// Formatter formats d as text with header h.  The Format method only formats
// the docs for d.  Any subtopics will only have their names printed with a
// synopsis if any exists.
//
// A Formatter does not need to be threadsafe.  The caller is responsible for
// synchronizing calls to a Formatter if concurrent access is expected.
type Formatter interface {
	Format(h *Header, d *Docs, indentdocs string) (string, error)
}

// NewFormatter returns a documentation formatter
func NewFormatter() Formatter {
	return newTextFormatter()
}

type textFormatter struct {
	indent string
	buf    *bytes.Buffer
}

func newTextFormatter() *textFormatter {
	return &textFormatter{}
}

func (g *textFormatter) Format(h *Header, d *Docs, indentdocs string) (string, error) {
	if g.buf == nil {
		g.buf = new(bytes.Buffer)
	} else {
		g.buf.Reset()
	}
	g.indent = indentdocs
	g.writeDocs(h, d)
	return g.buf.String(), nil
}

func (g *textFormatter) printf(format string, v ...interface{}) {
	if g.indent == "" {
		fmt.Fprintf(g.buf, format, v...)
	} else {
		text := fmt.Sprintf(format, v...)
		text = textutil.Indent(text, g.indent)
		g.buf.WriteString(text)
	}
}

func (g *textFormatter) writeDocs(h *Header, d *Docs) {
	if h != nil {
		fmt.Fprintf(g.buf, "%s %s\n\n", h.DocsType, h.Name)
	} else {
		g.printf("\n")
	}
	text := d.Usage
	text = textutil.Unindent(text)
	text = strings.TrimSpace(text)
	if text != "" {
		g.printf("  %s\n\n", textutil.Indent(text, "  "))
	}
	text = d.Sig
	text = textutil.Unindent(text)
	text = strings.TrimSpace(text)
	if text != "" {
		g.printf("Signature:\n\n")
		g.printf("%s\n\n", textutil.Indent(text, "  "))
	}
	if d.Desc != "" {
		text = textutil.Unindent(d.Desc)
		text = textutil.Wrap(text, 72)
		text = strings.TrimSpace(text)
		g.printf("%s\n\n", text)
	}
	numvar := d.NumVar()
	if numvar > 0 {
		g.printf("Variables:\n\n")
		for i := 0; i < numvar; i++ {
			g.printf("  %s", d.Var(i))
			typ := d.VarType(i)
			if typ != "" {
				g.printf("  %s\n", typ)
			} else {
				g.printf("\n")
			}
			text = d.VarDesc(i)
			text = textutil.Unindent(text)
			text = textutil.Wrap(text, 66)
			text = strings.TrimSpace(text)
			if text != "" {
				text = textutil.Indent(text, "      ")
				g.printf("%s\n\n", text)
			} else {
				g.printf("\n")
			}
		}
	}
	numparam := d.NumParam()
	if numparam > 0 {
		g.printf("Parameters:\n\n")
		for i := 0; i < numparam; i++ {
			g.printf("  %s", d.Param(i))
			typ := d.ParamType(i)
			if typ != "" {
				g.printf("  %s\n", typ)
			} else {
				g.printf("\n")
			}
			text = d.ParamDesc(i)
			text = textutil.Unindent(text)
			text = textutil.Wrap(text, 66)
			text = strings.TrimSpace(text)
			if text != "" {
				text = textutil.Indent(text, "      ")
				g.printf("%s\n\n", text)
			} else {
				g.printf("\n")
			}
		}
	}

	funcs := d.Funcs()
	if len(funcs) > 0 {
		g.printf("Functions:\n\n")
		for _, sub := range funcs {
			if sub.Name == "" {
				continue
			}
			g.printf("  %s\n", sub.Name)
			if sub.Docs == nil {
				g.printf("\n")
				continue
			}
			text = sub.Desc
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			text = textutil.Synopsis(text)
			text = textutil.Wrap(text, 66)
			if text != "" {
				text = textutil.Indent(text, "      ")
				g.printf("%s\n\n", text)
			} else {
				g.printf("\n")
			}
		}
	}

	others := d.Others()
	if len(others) > 0 {
		g.printf("Subtopics:\n\n")
		for _, sub := range others {
			if sub.Name == "" {
				continue
			}
			g.printf("  %s\n", sub.Name)
			if sub.Docs == nil {
				g.printf("\n")
				continue
			}
			text = sub.Desc
			text = textutil.Unindent(text)
			text = strings.TrimSpace(text)
			text = textutil.Synopsis(text)
			text = textutil.Wrap(text, 66)
			if text != "" {
				text = textutil.Indent(text, "      ")
				g.printf("%s\n\n", text)
			} else {
				g.printf("\n")
			}
		}
	}
}
