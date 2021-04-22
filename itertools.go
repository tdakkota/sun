package sun

import (
	"fmt"

	"go.starlark.net/starlark"
)

// Type that attempts to allow operations between numerics,
// i.e. float and int.
type floatOrInt struct {
	f_ *starlark.Float
	i_ *starlark.Int
}

// Unpacker for float or int type. This allows int types and float
// types to interact with one another, e.g. count(0, 0.1).
// type floatOrInt float6
func (p *floatOrInt) Unpack(v starlark.Value) error {
	errorMsg := "floatOrInt must have default initialization"

	switch v := v.(type) {
	case starlark.Int:
		if p.f_ != nil {
			return fmt.Errorf(errorMsg)
		}
		p.i_ = &v
		return nil
	case starlark.Float:
		if p.i_ != nil {
			return fmt.Errorf(errorMsg)
		}
		p.f_ = &v
		return nil
	}
	return fmt.Errorf("got %s, want float or int", v.Type())
}

func (fi *floatOrInt) add(n floatOrInt) error {
	switch {
	case fi.i_ != nil && n.i_ != nil:
		x := fi.i_.Add(*n.i_)
		fi.i_ = &x
		return nil
	case fi.i_ != nil && n.f_ != nil:
		x := starlark.Float(float64(fi.i_.Float()) + float64(*n.f_))
		fi.f_ = &x
		fi.i_ = nil
		return nil
	case fi.f_ != nil && n.i_ != nil:
		x := starlark.Float(float64(*fi.f_) + float64(n.i_.Float()))
		fi.f_ = &x
		fi.i_ = nil
		return nil
	case fi.f_ != nil && n.f_ != nil:
		x := starlark.Float(float64(*fi.f_) + float64(*n.f_))
		fi.f_ = &x
		return nil
	}
	return fmt.Errorf("float to int addition not possible")
}

func (fi floatOrInt) string() string {
	switch {
	case fi.i_ != nil:
		return fi.i_.String()
	case fi.f_ != nil:
		return fi.f_.String()
	default:
		// This block should not be reached.
		// starlark's String() method is being replicated
		// so an error is not raised.
		return ""
	}
}

// Equality operator between floatOrInt and starlark's Int, Float
// and Golang's int.
func (fi *floatOrInt) eq(v interface{}) bool {
	switch v := v.(type) {
	case starlark.Int:
		if fi.i_ != nil && *fi.i_ == v {
			return true
		} else {
			return false
		}
	case starlark.Float:
		if fi.f_ != nil && *fi.f_ == v {
			return true
		} else {
			return false
		}
	case int:
		if fi.i_ == nil {
			return false
		}
		var x int
		starlark.AsInt(*fi.i_, &x)
		return x == v
	}

	return false
}

type countObject struct {
	cnt    floatOrInt
	step   floatOrInt
	frozen bool
}

func (co *countObject) String() string {
	// As with the cpython implementation, we don't display
	// step when it is an integer equal to 1 (default step value).
	if co.step.eq(1) {
		return fmt.Sprintf("count(%v)", co.cnt.string())
	}
	return fmt.Sprintf("count(%v, %v)", co.cnt.string(), co.step.string())
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

type countIter struct {
	co *countObject
}

func (c *countIter) Next(p *starlark.Value) bool {
	if c.co.frozen {
		return false
	}

	switch {
	case c.co.cnt.i_ != nil:
		*p = c.co.cnt.i_
	case c.co.cnt.f_ != nil:
		*p = c.co.cnt.f_
	}

	c.co.cnt.add(c.co.step)

	return true
}

func (c *countIter) Done() {}

func count_(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		defaultStart            = starlark.MakeInt(0)
		defaultStep             = starlark.MakeInt(1)
		start        floatOrInt = floatOrInt{}
		step         floatOrInt = floatOrInt{}
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
	if start.f_ == nil && start.i_ == nil {
		start.i_ = &defaultStart
	}
	if step.f_ == nil && step.i_ == nil {
		step.i_ = &defaultStep
	}

	return &countObject{cnt: start, step: step}, nil
}
