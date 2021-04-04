package sun

import (
	"fmt"

	"go.starlark.net/starlark"
)

type filterFunc = func(x starlark.Value) (starlark.Value, error)

type filterIter struct {
	thread   *starlark.Thread
	function filterFunc
	iterator starlark.Iterator
}

func (f filterIter) Next(p *starlark.Value) bool {
	var x starlark.Value
	for {
		if !f.iterator.Next(&x) {
			return false
		}

		v, err := f.function(x)
		if err != nil {
			return false
		}

		if v.Truth() {
			*p = x
			return true
		}
	}
}

func (f filterIter) Done() {
	f.iterator.Done()
}

type filterObject struct {
	thread   *starlark.Thread
	function filterFunc
	iterable starlark.Iterable
	iterator starlark.Iterator
}

func (f filterObject) String() string {
	return fmt.Sprintf("<filter object>")
}

func (f filterObject) Type() string {
	return "filter"
}

func (f filterObject) Freeze() {
	f.iterable.Freeze()
}

func (f filterObject) Truth() starlark.Bool {
	return starlark.True
}

func (f filterObject) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: filter")
}

func (f filterObject) Iterate() starlark.Iterator {
	return filterIter{
		thread:   f.thread,
		function: f.function,
		iterator: f.iterator,
	}
}

func filter(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		function filterFunc
		iterable starlark.Iterable
	)
	if err := wantArgs(b.Name(), args, kwargs, 2); err != nil {
		return nil, err
	}

	switch fn := args[0].(type) {
	case starlark.Callable:
		function = func(x starlark.Value) (starlark.Value, error) {
			return starlark.Call(thread, fn, starlark.Tuple{x}, nil)
		}
	case starlark.NoneType:
		function = func(x starlark.Value) (starlark.Value, error) {
			return x.Truth(), nil
		}
	default:
		return nil, fmt.Errorf("got %s, want callable", fn.Type())
	}

	iterable, ok := args[1].(starlark.Iterable)
	if !ok {
		return nil, fmt.Errorf("got %s, want iterable", args[1].Type())
	}

	return &filterObject{
		thread:   thread,
		function: function,
		iterable: iterable,
		iterator: iterable.Iterate(),
	}, nil
}
