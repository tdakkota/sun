package sun

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var ItertoolsModule = &starlarkstruct.Module{
	Name: "itertools",
	Members: starlark.StringDict{
		"count": starlark.NewBuiltin("itertools.count", count_),
	},
}

type countObject struct {
	cnt    int
	step   int
	frozen bool
	value  starlark.Value
}

func newCountObject(cnt int, stepValue int) *countObject {
	return &countObject{cnt: cnt, step: stepValue, value: starlark.MakeInt(cnt)}
}

func (co *countObject) String() string {
	// As with the cpython implementation, we don't display
	// step when it is an integer equal to 1.
	if co.step == 1 {
		return fmt.Sprintf("count(%v)", co.cnt)
	}
	return fmt.Sprintf("count(%v, %v)", co.cnt, co.step)
}

func (co *countObject) Type() string {
	return "itertools.count"
}

func (co *countObject) Freeze() {
	if !co.frozen {
		co.frozen = true
		co.value.Freeze()
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
	*p = starlark.MakeInt(c.co.cnt)
	c.co.cnt += c.co.step
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
		start int
		step  int
	)

	if err := starlark.UnpackPositionalArgs(
		"count", args, kwargs, 0, &start, &step,
	); err != nil {
		return nil, fmt.Errorf(
			`Got %v but expected NoneType or valid
	integer values for start and step, such as (0, 1).`, args,
		)
	}

	const (
		defaultStart = 0
		defaultStep  = 1
	)
	// The rules for populating the count object based on the number
	// of args passed is as follows:
	// 	0 args -> default values for start and step
	// 	1 args -> arg defines start, default for step
	// 	2 args -> both start and step are defined by args
	var co_ *countObject
	switch nargs := len(args); {
	case nargs == 0:
		co_ = newCountObject(defaultStart, defaultStep)
	case nargs == 1:
		co_ = newCountObject(start, defaultStep)
	default: // nargs == 2
		co_ = newCountObject(start, step)
	}

	return co_, nil
}
