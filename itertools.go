package sun

import (
	"fmt"

	"go.starlark.net/starlark"
)

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
}

func (co countObject) String() string {
	// As with the cpython implementation, we don't display
	// step when it is an integer equal to 1 (default step value).
	step, ok := co.step.value.(starlark.Int)
	if ok {
		if x, ok := step.Int64(); ok && x == 1 {
			return "count(1)"
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
	// TODO(algebra8): Implement inherited type object hash.
	return uint32(10), nil
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

	return &countObject{cnt: start, step: step}, nil
}
