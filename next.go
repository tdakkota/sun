package sun

import (
	"errors"

	"go.starlark.net/starlark"
)

// ErrIterationDone denotes that iterator can't return values anymore.
var ErrIterationDone = errors.New("iteration done")

func next(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		iterable starlark.Iterable
		default_ starlark.Value
	)

	if err := starlark.UnpackArgs(
		b.Name(), args, kwargs,
		"iterator", &iterable, "default?", &default_,
	); err != nil {
		return nil, err
	}

	iter := iterable.Iterate()
	defer iter.Done()

	var x starlark.Value
	switch {
	case iter.Next(&x):
		return x, nil
	case default_ != nil:
		return default_, nil
	default:
		return nil, ErrIterationDone
	}
}
