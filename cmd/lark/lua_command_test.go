package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bmatsuo/lark/lib/lark/core"
	"github.com/codegangsta/cli"
)

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
