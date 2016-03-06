package gluatest

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/bmatsuo/lark/gluamodule"
	"github.com/yuin/gopher-lua"
)

var testFile = &File{
	Module: gluamodule.New("test.module", basicTestLoader),
	Path:   "gluatest_test.lua",
}

func TestFile(t *testing.T) {
	l := testFile.Load(t)
	defer l.Close()

	l.Push(l.GetGlobal("helper"))
	err := l.PCall(0, 1, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if l.Get(1).String() != "HELP!" {
		t.Errorf("return value: %q (!= %q)", l.Get(1), "HELP!")
	}
}

func TestFile_sorted(t *testing.T) {
	funcs := testFile.getTestFuncs(t)
	sorted := make([]string, len(funcs))
	copy(sorted, funcs)
	for i := range sorted {
		if sorted[i] != funcs[i] {
			t.Errorf("function %d: %q (!= %q)", i, funcs[i], sorted[i])
		}
	}
}

func TestFile_runTest(t *testing.T) {
	state := map[string]interface{}{}
	testFile := &File{
		Module: gluamodule.New("test.module", statefulLoader(state)),
		Path:   "gluatest_test.lua",
	}
	testfn := "test_a"
	testFile.runTest(t, testfn, getGlobalFunction(testfn))
	if !reflect.DeepEqual(state["__test_setup"], true) {
		t.Errorf("%s: %v (!= %v)", "__test_setup", state["__test_setup"], true)
	}
	if !reflect.DeepEqual(state["__test_teardown"], true) {
		t.Errorf("%s: %v (!= %v)", "__test_teardown", state["__test_teardown"], true)
	}
	if !reflect.DeepEqual(state["test_a"], true) {
		t.Errorf("%s: %v (!= %v)", "test_a", state["test_a"], true)
	}
	if !reflect.DeepEqual(state["test_b"], nil) {
		t.Errorf("%s: %v (!= %v)", "test_b", state["test_b"], true)
	}
}

func TestFile_runTest_failure(t *testing.T) {
	state := map[string]interface{}{}
	testFile := &File{
		Module: gluamodule.New("test.module", statefulLoader(state)),
		Path:   "gluatest_test.lua",
	}
	var testFailed bool
	testfn := "test_fail"
	recorder := &TBFailRecorder{
		FailFunc: func() { testFailed = true },
		TB:       t,
	}
	testFile.runTest(recorder, testfn, getGlobalFunction(testfn))
	if !testFailed {
		t.Errorf("test did not fail")
	}
	if !reflect.DeepEqual(state["test_a"], nil) {
		t.Errorf("%s: %v (!= %v)", "test_a", state["test_a"], nil)
	}
	if !reflect.DeepEqual(state["test_fail"], true) {
		t.Errorf("%s: %v (!= %v)", "test_fail", state["test_fail"], true)
	}
	if !reflect.DeepEqual(state["__test_setup"], true) {
		t.Errorf("%s: %v (!= %v)", "__test_setup", state["__test_setup"], true)
	}
	if !reflect.DeepEqual(state["__test_teardown"], true) {
		t.Errorf("%s: %v (!= %v)", "__test_teardown", state["__test_teardown"], true)
	}
}

func BenchmarkRequireModule(b *testing.B) {
	emptyModuleFile := &File{
		Module: gluamodule.New("test.module", basicTestLoader),
		Path:   "gluatest_test.lua",
	}
	emptyModuleFile.BenchmarkRequireModule(b)
}

type TBFailRecorder struct {
	ReallyFail bool
	FailFunc   func()
	testing.TB
}

func (t *TBFailRecorder) Fail() {
	t.FailFunc()
	if t.ReallyFail {
		t.TB.Fail()
	}
}

func (t *TBFailRecorder) FailNow() {
	t.Fail()
	runtime.Goexit()
}

func (t *TBFailRecorder) Fatal(v ...interface{}) {
	args := []string{"(witheld fatal) "}
	for _, v := range v {
		args = append(args, fmt.Sprint(v))
	}
	t.Log(strings.Join(args, ""))
	t.FailNow()
}

func (t *TBFailRecorder) Fatalf(format string, v ...interface{}) {
	args := []string{"(witheld fatal) "}
	args = append(args, fmt.Sprintf(format, v...))
	t.Log(strings.Join(args, ""))
	t.FailNow()
}

func (t *TBFailRecorder) Error(v ...interface{}) {
	args := []string{"(witheld error) "}
	for _, v := range v {
		args = append(args, fmt.Sprint(v))
	}
	t.Log(strings.Join(args, ""))
	t.Fail()
}

func (t *TBFailRecorder) Errorf(format string, v ...interface{}) {
	args := []string{"(witheld error) "}
	args = append(args, fmt.Sprintf(format, v...))
	t.Log(strings.Join(args, ""))
	t.Fail()
}

func statefulLoader(state map[string]interface{}) lua.LGFunction {
	return func(l *lua.LState) int {
		mod := l.NewTable()

		l.SetField(mod, "set_value", l.NewFunction(func(l *lua.LState) int {
			name := l.CheckString(1)
			val := l.CheckAny(2)
			l.SetTop(0)
			if val == lua.LNil {
				state[name] = nil
			} else {
				switch val := val.(type) {
				case lua.LString:
					state[name] = string(val)
				case lua.LNumber:
					state[name] = float64(val)
				case lua.LBool:
					state[name] = bool(val)
				default:
					l.RaiseError("argument #2 is not a string, number, or boolean: %s", val.Type())
				}
			}
			return 0
		}))

		l.Push(mod)
		return 1
	}
}

func emptyLoader(l *lua.LState) int {
	l.Push(l.NewTable())
	return 1
}

func basicTestLoader(l *lua.LState) int {
	mod := l.NewTable()

	double := l.NewClosure(func(l *lua.LState) int {
		x := l.CheckNumber(1)
		l.Push(lua.LNumber(2 * x))
		return 1
	}, mod)
	l.SetField(mod, "double", double)

	l.Push(mod)
	return 1
}
