package textutil

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Wrap wraps lines in text to width.  Indented lines are considered
// preformatted and are not modified by Wrap.
//
// BUG: does not properly handle multibyte runes, accents and such.
func Wrap(text string, width int) string {
	if text == "" {
		return text
	}

	var lead string
	var buf bytes.Buffer
	s := bufio.NewScanner(strings.NewReader(text))
	for s.Scan() {
		line := s.Text()
		if strings.IndexFunc(line, unicode.IsSpace) == 0 {
			if lead != "" {
				buf.WriteString(lead)
				buf.WriteString("\n")
			}
			buf.WriteString(line)
			buf.WriteString("\n")
			continue
		}

		if lead != "" {
			lead = ""
			if lead[len(lead)-1] == '.' {
				line = lead + "  " + line
			} else {
				line = lead + " " + line
			}
		}

		if len(line) < width {
			buf.WriteString(line)
			buf.WriteString("\n")
			continue
		}

		w := 0
		for i := range line {
			w = i
			if i >= width {
				break
			}
		}

		ispace := strings.LastIndexFunc(line[:w], unicode.IsSpace)
		if ispace < 0 {
			ispace = strings.IndexFunc(line[w:], unicode.IsSpace)
			if ispace >= 0 {
				ispace += w
			}
		}
		if ispace < 0 {
			buf.WriteString(line)
		} else {
			lead = strings.TrimSpace(line[ispace:])
			line = strings.TrimSpace(line[:ispace])
			buf.WriteString(line)
		}
		buf.WriteString("\n")
	}
	if lead != "" {
		buf.WriteString(lead)
		buf.WriteString("\n")
	}

	result := buf.String()
	if text[len(text)-1] != '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// Indent prepends indent to each line of text.
func Indent(text, indent string) string {
	if text == "" {
		return text
	}
	var buf bytes.Buffer
	s := bufio.NewScanner(strings.NewReader(text))
	for s.Scan() {
		line := s.Text()
		if line != "" {
			buf.WriteString(indent)
			buf.WriteString(line)
		}
		buf.WriteString("\n")
	}
	result := buf.String()
	if text[len(text)-1] != '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// Unindent removes indentation from non-blank lines in text.
func Unindent(text string) string {
	if text == "" {
		return text
	}
	var found bool
	var indent string
	s := bufio.NewScanner(strings.NewReader(text))
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		if !found {
			found = true
			i := strings.IndexFunc(line, func(c rune) bool { return !unicode.IsSpace(c) })
			if i == 0 {
				return text
			} else if i < 0 {
				indent = line
			} else {
				indent = line[:i]
			}
		} else if !strings.HasPrefix(line, indent) {
			indent = commonPrefix(indent, line)
		}
	}

	var buf bytes.Buffer
	s = bufio.NewScanner(strings.NewReader(text))
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, indent) {
			buf.WriteString(line[len(indent):])
		} else {
			buf.WriteString(line)
		}
		buf.WriteString("\n")
	}

	result := buf.String()
	if text[len(text)-1] == '\n' {
		return result
	}
	return result[:len(result)-1]
}

func commonPrefix(s1, s2 string) string {
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	s2 = s2[:len(s1)]

	var prefix string
	a, b := s1, s2
	for a != "" {
		r1, size1 := utf8.DecodeRuneInString(a)
		r2, size2 := utf8.DecodeRuneInString(b)
		if r1 != r2 {
			return prefix
		}
		if r1 == utf8.RuneError && size1 == 1 {
			return prefix
		}
		prefix = s1[:len(prefix)+size1]
		a, b = a[size1:], b[size2:]
	}
	return prefix
}
