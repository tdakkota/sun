package sun

import (
	"fmt"

	"go.starlark.net/starlark"
)

func wantArgs(fnName string, args starlark.Tuple, kwargs []starlark.Tuple, count int) error {
	switch {
	case len(kwargs) > 0:
		return fmt.Errorf("%s does not accept keyword arguments", fnName)
	case len(args) != count:
		return fmt.Errorf("%s: got %d arguments, want exactly %d", fnName, len(args), count)
	default:
		return nil
	}
}
