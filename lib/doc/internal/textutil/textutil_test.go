package textutil

import "testing"

func TestSynopsis(t *testing.T) {
	tests := []struct {
		text   string
		expect string
	}{
		{"", ""},
		{"\nfour pinapples", "four pinapples"},
		{"four pinapples", "four pinapples"},
		{"four. pinapples", "four."},
		{"four pinapples\n\nfive oranges", "four pinapples"},
		{"four pinapples.\n\nfive oranges", "four pinapples."},
	}

	for i, test := range tests {
		out := Synopsis(test.text)
		if out != test.expect {
			t.Errorf("test %d: %q (!= %q)", i, out, test.expect)
		}
	}
}

func TestUnindent(t *testing.T) {
	tests := []struct {
		text   string
		expect string
	}{
		{"", ""},
		{"\nfour pinapples", "\nfour pinapples"},
		{"four pinapples\n", "four pinapples\n"},
		{"four pinapples", "four pinapples"},
		{"\n four\n pinapples\n ", "\nfour\npinapples\n"},
		{"\nfour\n pinapples\n ", "\nfour\n pinapples\n "},
		{"\n four\npinapples\n ", "\n four\npinapples\n "},
		{"\n four\n pinapples\n", "\nfour\npinapples\n"},
		{"\n\n four\n pinapples", "\n\nfour\npinapples"},
		{"\n  \n four\n pinapples", "\n \nfour\npinapples"},
		{" four\n pinapples\n  \n ", "four\npinapples\n \n"},
	}

	for i, test := range tests {
		out := Unindent(test.text)
		if out != test.expect {
			t.Errorf("test %d: %q (!= %q)", i, out, test.expect)
		}
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		width  int
		text   string
		expect string
	}{
		{5, "", ""},
		{5, "\nfour pinapples", "\nfour\npinapples"},
		{5, "four pinapples\n\n\n", "four\npinapples\n\n\n"},
		{5, "four pinapples", "four\npinapples"},
		{5, "four \npinapples", "four\npinapples"},
		{5, "four pinapples\n\nfive oranges", "four\npinapples\n\nfive\noranges"},
		{100, "four pinapples", "four pinapples"},
	}

	for i, test := range tests {
		out := Wrap(test.text, test.width)
		if out != test.expect {
			t.Errorf("test %d: %q (!= %q)", i, out, test.expect)
		}
	}
}
