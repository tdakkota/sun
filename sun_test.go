package sun

import (
	"fmt"
	"path/filepath"
	"testing"

	"go.starlark.net/starlark"
)

// load implements the 'load' operation as used in the evaluator tests.
func load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	if module == "assert.star" {
		return LoadAssertModule()
	}
	return nil, fmt.Errorf("load not implemented")
}

func runTestData(t *testing.T, name string) {
	thread := &starlark.Thread{Load: load}
	SetReporter(thread, t)

	filename := filepath.Join("testdata", name)
	if _, err := starlark.ExecFile(thread, filename, nil, Module.Members); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}
