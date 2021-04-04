package sun

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestFilter(t *testing.T) {
	runTestData(t, "filter.star", starlark.StringDict{
		"filter": starlark.NewBuiltin("filter", filter),
	})
}
