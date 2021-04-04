package sun

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestCallable(t *testing.T) {
	runTestData(t, "callable.star", starlark.StringDict{
		"callable": starlark.NewBuiltin("callable", callable),
	})
}
