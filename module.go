package sun

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module builtins is a Starlark module of Python-like builtins functions.
var Module = &starlarkstruct.Module{
	Name: "builtins",
	Members: starlark.StringDict{
		"filter":   starlark.NewBuiltin("filter", filter),
		"callable": starlark.NewBuiltin("callable", callable),
	},
}
