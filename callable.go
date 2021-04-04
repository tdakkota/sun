package sun

import (
	"go.starlark.net/starlark"
)

func callable(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	if err := wantArgs(b.Name(), args, kwargs, 1); err != nil {
		return nil, err
	}

	if _, ok := args[0].(starlark.Callable); ok {
		return starlark.True, nil
	}
	return starlark.False, nil
}
