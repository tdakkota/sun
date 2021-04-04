package sun_test

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/tdakkota/sun"
)

func Example_filter() {
	code := "list(filter(lambda x: x % 2 == 0, range(10)))"
	thread := &starlark.Thread{Name: "main"}
	result, err := starlark.Eval(thread, "example.star", code, sun.Module.Members)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// Output:
	// [0, 2, 4, 6, 8]
}
