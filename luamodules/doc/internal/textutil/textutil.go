package textutil

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Unindent removes indentation from non-blank lines in text.
func Unindent(text string) string {
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
