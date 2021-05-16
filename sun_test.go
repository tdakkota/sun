package sun

import (
	"fmt"
	"path/filepath"
	"testing"

	"go.starlark.net/starlark"
)

// load implements the 'load' operation as used in the evaluator tests.
func load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	if module == "assert.star" {
		return LoadAssertModule()
	}
	return nil, fmt.Errorf("load not implemented")
}

func runTestData(t *testing.T, name string) {
	thread := &starlark.Thread{Load: load}
	SetReporter(thread, t)

	filename := filepath.Join("testdata", name)

	// "Merge" ItertoolsModule and Module so ExecFile can use both
	// predeclared names.
	// The below method is surely not the correct way to accomplish
	// the task but works for testing purposes.
	// TODO(algebra8): Find and implement correct method for merging modules
	for k, v := range ItertoolsModule.Members {
		Module.Members[k] = v
	}

	if _, err := starlark.ExecFile(thread, filename, nil, Module.Members); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}
