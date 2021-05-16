package sun

import (
	"fmt"
	"runtime"

	"go.starlark.net/starlark"
)

type mapIter struct {
	thread    *starlark.Thread
	function  starlark.Callable
	iterators []starlark.Iterator
	buf       starlark.Tuple
}

func (f *mapIter) Next(p *starlark.Value) bool {
	f.buf = f.buf[:0]

	var x starlark.Value
	for _, iter := range f.iterators {
		if !iter.Next(&x) {
			return false
		}
		f.buf = append(f.buf, x)
	}

	v, err := starlark.Call(f.thread, f.function, f.buf, nil)
	if err != nil {
		return false
	}

	*p = v
	return true
}

func (f *mapIter) Done() {
	for i := range f.iterators {
		f.iterators[i].Done()
	}
}

type mapObject struct {
	thread    *starlark.Thread
	function  starlark.Callable
	iterables []starlark.Iterable
	iterators []starlark.Iterator
}

func (f mapObject) String() string {
	return "<map object>"
}

func (f mapObject) Type() string {
	return "map"
}

func (f mapObject) Freeze() {
	f.function.Freeze()
	for _, iterable := range f.iterables {
		iterable.Freeze()
	}
}

func (f mapObject) Truth() starlark.Bool {
	return starlark.True
}

func (f mapObject) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: map")
}

func (f mapObject) Iterate() starlark.Iterator {
	// TODO(tdakkota): specialize iterator if there is only one iterable.
	return &mapIter{
		thread:    f.thread,
		function:  f.function,
		iterators: f.iterators,
		buf:       make([]starlark.Value, 0, len(f.iterators)),
	}
}

func map_(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		function  starlark.Callable
		iterables = make([]starlark.Iterable, len(args)-1)
	)

	unpack := make([]interface{}, 0, len(args))
	unpack = append(unpack, &function)
	for i := range iterables {
		unpack = append(unpack, &iterables[i])
	}
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, unpack...); err != nil {
		return nil, err
	}

	iterators := make([]starlark.Iterator, len(iterables))
	for i := range iterators {
		iterators[i] = iterables[i].Iterate()
	}

	obj := &mapObject{
		thread:    thread,
		function:  function,
		iterables: iterables,
		iterators: iterators,
	}
	// TODO(tdakkota): find better way to release iterators
	runtime.SetFinalizer(obj, func(m *mapObject) {
		for _, iterator := range m.iterators {
			iterator.Done()
		}
	})

	return obj, nil
}
