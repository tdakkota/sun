package sun

import (
	"fmt"

	"go.starlark.net/starlark"
)

func formatInt(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
	f string,
) (starlark.Value, error) {
	var i starlark.Int
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &i); err != nil {
		return nil, err
	}

	if v, ok := i.Int64(); ok {
		return starlark.String(fmt.Sprintf(f, v)), nil
	}

	return starlark.String(fmt.Sprintf(f, i.BigInt())), nil
}

func bin(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	return formatInt(thread, b, args, kwargs, "%#b")
}

func oct(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	return formatInt(thread, b, args, kwargs, "%O")
}

func hex(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	return formatInt(thread, b, args, kwargs, "%#x")
}
