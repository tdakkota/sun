package sun

import (
	"fmt"

	"github.com/google/uuid"
	"go.starlark.net/starlark"
)

/* count object ************************************************************/

// float or int type to allow mixed inputs.
type floatOrInt struct {
	value starlark.Value
}

// Unpacker for floatOrInt.
func (p *floatOrInt) Unpack(v starlark.Value) error {
	switch v := v.(type) {
	case starlark.Int:
		p.value = v
		return nil
	case starlark.Float:
		p.value = v
		return nil
	}
	return fmt.Errorf("got %s, want float or int", v.Type())
}

func (f *floatOrInt) add(n floatOrInt) error {
	switch _f := f.value.(type) {
	case starlark.Int:
		switch _n := n.value.(type) {
		// int + int
		case starlark.Int:
			f.value = _f.Add(_n)
			return nil
		// int + float
		case starlark.Float:
			_n += _f.Float()
			f.value = _n
			return nil
		}
	case starlark.Float:
		switch _n := n.value.(type) {
		// float + int
		case starlark.Int:
			_f += _n.Float()
			f.value = _f
			return nil
		// float + float
		case starlark.Float:
			_f += _n
			f.value = _f
			return nil
		}
	}

	return fmt.Errorf("error with addition: types are not int, float combos")
}

func (f *floatOrInt) String() string {
	return f.value.String()
}

// Iterator implementation for countObject.
type countIter struct {
	co *countObject
}

func (c *countIter) Next(p *starlark.Value) bool {
	if c.co.frozen {
		return false
	}

	*p = c.co.cnt.value

	if e := c.co.cnt.add(c.co.step); e != nil {
		return false
	}

	return true
}

func (c *countIter) Done() {}

// countObject implementation as a starlark.Value.
type countObject struct {
	cnt, step floatOrInt
	frozen    bool
	id        uint32
}

func newCountObject(start, step floatOrInt) *countObject {
	return &countObject{cnt: start, step: step, id: uuid.New().ID()}
}

func (co countObject) String() string {
	// As with the cpython implementation, we don't display
	// step when it is an integer equal to 1 (default step value).
	step, ok := co.step.value.(starlark.Int)
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
	// Cpython's count object inherits tp_hash from object,
	// where object's hash is calculated by:
	// 	id() >> 4
	// Here, a UUID.ID() should suffice.
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
		start        floatOrInt
		step         floatOrInt
	)

	if err := starlark.UnpackPositionalArgs(
		"count", args, kwargs, 0, &start, &step,
	); err != nil {
		return nil, fmt.Errorf(
			"Got %v but expected no args, or one or two valid numbers",
			args.String(),
		)
	}

	// Check if start or step require default values.
	if start.value == nil {
		start.value = defaultStart
	}
	if step.value == nil {
		step.value = defaultStep
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
	// Python's itertools.islice seems to make a copy of any
	// underlying iterable, e.g.:
	// 	>>> a = [1, 2, 3]
	// 	>>> s = itertools.islice(a, 3)
	// 	>>> a = []
	// 	>>> list(s)
	// 	[1, 2, 3]
	// And since the islice object itself is immutable and hashable,
	// we consider the entire object immutable.
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
			return fmt.Errorf("expected int or None, got %s\n", v.String())
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
			starlark.AsInt(a, &start)
		}
		if b != starlark.None {
			starlark.AsInt(b, &stop)
		}
		if c != starlark.None {
			starlark.AsInt(c, &step)
		}
	} else { // 2 args; itertools.islice(iterable, stop)
		if a != starlark.None {
			starlark.AsInt(a, &stop)
		}
	}

	return isliceObject{
		iterator: iterable.Iterate(),
		next:     start,
		stop:     stop,
		step:     step,
	}, nil
}
