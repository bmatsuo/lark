package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFindTaskFiles_precedence(t *testing.T) {
	root, err := ioutil.TempDir("", "lark-project-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)

	tasksDir := filepath.Join(root, "lark_tasks")
	os.MkdirAll(tasksDir, 0755)

	luaFiles := []string{
		"lark.lua",
		filepath.Join(TaskDir, "a.lua"),
		filepath.Join(TaskDir, "b.lua"),
	}

	for _, f := range luaFiles {
		msg := fmt.Sprintf("in %s!", f)
		content := fmt.Sprintf("print(%q)\n", msg)
		path := filepath.Join(root, f)
		err := ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Errorf("file %s: %v", f, err)
		}
	}

	files, err := FindTaskFiles(root)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) != len(luaFiles) {
		t.Errorf("found %d task files: %q", len(files), files)
	}

	for i, f := range luaFiles {
		baseExpect := filepath.Base(f)
		if filepath.Base(files[i]) != baseExpect {
			t.Errorf("task file at index %d: %q (!= %q)", i, filepath.Base(files[i]), baseExpect)
		}
		dirExpect := filepath.Join(root, filepath.Dir(f))
		if filepath.Dir(files[i]) != dirExpect {
			t.Errorf("task file dirname at index %d: %q (!= %q)", i, filepath.Dir(files[i]), dirExpect)
		}
	}
}

func TestFindTaskFiles_TaskDir(t *testing.T) {
	root, err := ioutil.TempDir("", "lark-project-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)

	tasksDir := filepath.Join(root, "lark_tasks")
	os.MkdirAll(tasksDir, 0755)

	luaFiles := []string{
		"test1.lua",
		"test2.lua",
		"test3.lua",
	}
	nonLuaFiles := []string{
		"data.txt",
	}

	for _, f := range luaFiles {
		msg := fmt.Sprintf("in %s!", f)
		content := fmt.Sprintf("print(%q)\n", msg)
		path := filepath.Join(tasksDir, f)
		err := ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Errorf("file %s: %v", f, err)
		}
	}

	for _, f := range nonLuaFiles {
		content := "data"
		path := filepath.Join(tasksDir, f)
		err := ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Errorf("file %s: %v", f, err)
		}
	}

	files, err := FindTaskFiles(root)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) != len(luaFiles) {
		t.Errorf("found %d task files: %q", len(files), files)
	}

	for i, f := range luaFiles {
		if filepath.Base(files[i]) != f {
			t.Errorf("task file at index %d: %q (!= %q)", i, filepath.Base(files[i]), f)
		}
		if filepath.Dir(files[i]) != tasksDir {
			t.Errorf("task file dirname at index %d: %q (!= %q)", i, filepath.Dir(files[i]), f)
		}
	}
}
