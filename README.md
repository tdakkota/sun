# sun
Sun is a Starlark module with Python builtins.

## Installation
```bash
go get github.com/tdakkota/sun
```

## Usage
```go
package main

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/tdakkota/sun"
)

func main() {
	code := "list(filter(lambda x: x % 2 == 0, range(10)))"

	// Eval Starlark expresion.
	thread := &starlark.Thread{Name: "main"}
	result, err := starlark.Eval(thread, "example.star", code, sun.Module.Members)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// Output:
	// [0, 2, 4, 6, 8]
}
```