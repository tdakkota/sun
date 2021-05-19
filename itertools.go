package sun

import (
	"fmt"

	"github.com/google/uuid"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// A note on Hash functions:
// In the CPython implementation of itertools, some itertools methods, such as
// count and islice, inherit tp_hash from object where object's hash is
// calculated by id() >> 4 and where id(), in some Python implementations,
// returns the memory address of the underlying object.
// In this itertools module, a UUID.ID() is used.

/* count object ************************************************************/

// Iterator implementation for countObject.
type countIter struct {
	co *countObject
}

func (c *countIter) Next(p *starlark.Value) bool {
	if c.co.frozen {
		return false
	}

	*p = c.co.cnt

	// Numeric types for count and step should be guaranteed in countObject
	// creation in count_ function.
	count, step := c.co.cnt, c.co.step
	c.co.cnt, _ = starlark.Binary(syntax.PLUS, count, step)

	return true
}

func (c *countIter) Done() {}

// countObject implementation as a starlark.Value.
type countObject struct {
	cnt, step starlark.Value
	frozen    bool
	id        uint32
}

func newCountObject(start, step starlark.Value) *countObject {
	return &countObject{cnt: start, step: step, id: uuid.New().ID()}
}

func (co countObject) String() string {
	// As with the cpython implementation, we don't display
	// step when it is an integer equal to 1 (default step value).
	step, ok := co.step.(starlark.Int)
	if ok {
		if x, ok := step.Int64(); ok && x == 1 {
			return fmt.Sprintf("count(%v)", co.cnt.String())
		}
	}
	return fmt.Sprintf("count(%v, %v)", co.cnt.String(), co.step.String())
}

func (co *countObject) Type() string {
	return "itertools.count"
}

func (co *countObject) Freeze() {
	if !co.frozen {
		co.frozen = true
	}
}

func (co *countObject) Truth() starlark.Bool {
	return starlark.True
}

func (co *countObject) Hash() (uint32, error) {
	return co.id, nil
}

func (co *countObject) Iterate() starlark.Iterator {
	return &countIter{co: co}
}

func count_(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		defaultStart = starlark.MakeInt(0)
		defaultStep  = starlark.MakeInt(1)
		start        starlark.Value
		step         starlark.Value
	)

	if err := starlark.UnpackPositionalArgs(
		"count", args, kwargs, 0, &start, &step,
	); err != nil {
		return nil, err
	}

	if start == nil {
		start = defaultStart
	}
	if step == nil {
		step = defaultStep
	}

	// Assert that count and step are numeric starlark Values.
	switch start.(type) {
	case starlark.Float, starlark.Int:
	default:
		return nil, fmt.Errorf("a number is required")
	}
	switch step.(type) {
	case starlark.Float, starlark.Int:
	default:
		return nil, fmt.Errorf("a number is required")
	}

	return newCountObject(start, step), nil
}

/* islice object ************************************************************/

// isliceObject as a starlark.Value
type isliceObject struct {
	// Store the iterator directly on the object; see
	// https://github.com/tdakkota/sun/issues/14
	iterator starlark.Iterator
	next     int
	stop     int
	step     int
	id       uint32
	// TODO(algebra8): Need itercount?
}

func (is isliceObject) String() string {
	return "<itertools.islice object>"
}

func (is isliceObject) Type() string {
	return "itertools.islice"
}

func (is isliceObject) Freeze() {
	// Since isliceObject does not hold onto the iterable,
	// the iterable value is not reachable from it so Freeze
	// doesn't need to do anything.
}

func (is isliceObject) Truth() starlark.Bool {
	return starlark.True
}

func (is isliceObject) Hash() (uint32, error) {
	return is.id, nil
}

// Iterator for islice object
type isliceIter struct {
	islice *isliceObject
	cnt    int
}

func (it *isliceIter) Next(p *starlark.Value) bool {
	var (
		x       starlark.Value
		oldNext int
	)

	stop := it.islice.stop

	// Get iterator up to the "next" iteration.
	for it.cnt < it.islice.next {
		if !it.islice.iterator.Next(&x) {
			return false
		}
		it.cnt += 1
	}

	if it.cnt >= stop {
		return false
	}

	if !it.islice.iterator.Next(&x) {
		return false
	}

	*p = x

	it.cnt += 1
	oldNext = it.islice.next
	it.islice.next += it.islice.step
	if it.islice.next < oldNext || it.islice.next > stop {
		it.islice.next = it.islice.stop
	}

	return true
}

func (it *isliceIter) Done() {
	it.islice.iterator.Done()
}

// isliceObject as a starlark.Iteratable
func (is isliceObject) Iterate() starlark.Iterator {
	return &isliceIter{islice: &is}
}

func assertPosIntOrNone(vs ...starlark.Value) error {
	for _, v := range vs {
		i, ok := v.(starlark.Int)
		if !ok && v != starlark.None {
			return fmt.Errorf("expected int or None, got %s\n", v.Type())
		}
		if ok && i.Sign() == -1 {
			return fmt.Errorf("expected non-negative values, got %s\n", v.String())
		}
	}
	return nil
}

func islice(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		iterable starlark.Iterable

		// Positional args from islice call
		a starlark.Value
		b starlark.Value
		c starlark.Value

		// islice values
		start int = 0
		stop  int = (1 << 63) - 1
		step  int = 1
	)

	if err := starlark.UnpackPositionalArgs(
		"islice", args, kwargs, 2, &iterable,
		&a, &b, &c,
	); err != nil {
		return nil, err
	}

	if a == nil {
		a = starlark.None
	}
	if b == nil {
		b = starlark.None
	}
	if c == nil {
		c = starlark.None
	}
	if err := assertPosIntOrNone(a, b, c); err != nil {
		return nil, err
	}

	if len(args) > 2 { // itertools.islice(iterable, start, stop[, step])
		if a != starlark.None {
			if err := starlark.AsInt(a, &start); err != nil {
				return nil, err
			}
		}
		if b != starlark.None {
			if err := starlark.AsInt(b, &stop); err != nil {
				return nil, err
			}
		}
		if c != starlark.None {
			if err := starlark.AsInt(c, &step); err != nil {
				return nil, err
			}
		}
	} else { // 2 args; itertools.islice(iterable, stop)
		if a != starlark.None {
			if err := starlark.AsInt(a, &stop); err != nil {
				return nil, err
			}
		}
	}

	return isliceObject{
		iterator: iterable.Iterate(),
		next:     start,
		stop:     stop,
		step:     step,
	}, nil
}
