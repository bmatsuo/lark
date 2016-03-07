package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/bmatsuo/lark/lib"
	"github.com/bmatsuo/lark/lib/lark/core"
	"github.com/codegangsta/cli"
)

func TestREPLHelp(t *testing.T) {
	help := REPLHelp()

	marker := "\nBuiltin Modules\n\n"
	index := strings.Index(help, marker)
	if index < 0 {
		t.Fatalf("cannot find module list")
	}

	modules := help[index+len(marker):]
mloop:
	for i, m := range lib.Modules {
		name := m.Name()
		for _, intern := range lib.InternalModules {
			if intern.Name() == name {
				continue mloop
			}
		}
		patt := fmt.Sprintf(`\s*%s\s`, name)
		re, err := regexp.Compile(patt)
		if err != nil {
			t.Errorf("pattern %d: %v", i, err)
			continue
		}
		if !re.MatchString(modules) {
			t.Errorf("couldn't find module: %q\n%s", name, modules)
		}
	}
}

func TestLua(t *testing.T) {
	rec := testLogging(t, true)
	defer rec.Reset()

	app := Init(cli.NewApp())
	app.Run([]string{"lark", "lua", "-c", "lark.log('testok')"})

	output := rec.Output()
	testString := " testok\n"
	if !strings.Contains(output, testString) {
		t.Errorf("output:\n\t%q\n\tdoes not contain\n\t%q", output, testString)
	}
}

func testLogging(t *testing.T, rec bool) *testLoggerOutput {
	r := &testLoggerOutput{t, nil}
	if rec {
		r.buf = &bytes.Buffer{}
	}
	r.TakeOver()
	return r
}

type testLoggerOutput struct {
	t   *testing.T
	buf *bytes.Buffer
}

func (r *testLoggerOutput) TakeOver() {
	log.SetOutput(r)
	core.InitModule(r, 0)
}

func (r *testLoggerOutput) Reset() {
	log.SetOutput(os.Stderr)
	core.InitModule(os.Stderr, 0)
}

func (r *testLoggerOutput) Output() string {
	return r.buf.String()
}

func (r *testLoggerOutput) Write(b []byte) (int, error) {
	if r.buf != nil {
		r.buf.Write(b)
	}
	r.t.Log(string(b))

	return len(b), nil
}
