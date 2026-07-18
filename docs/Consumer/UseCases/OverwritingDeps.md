# Overwriting Dependencies

## Description
This guide describes how to overwrite specific dependencies after initialization. You can change individual functions of the `Deps` struct and alter the library's behavior without creating a full adapter from scratch.

## Workflow

1. Initialize the deps using an existing adapter, which will provide a fully populated `deps.Deps` struct.
2. Store the returned `deps.Deps` object in a variable before passing it to `keep_lib.New()`.
3. Overwrite specific function fields in the `Deps` struct to inject your custom behavior.
Like this:

```go
package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

func main() {
	// 1. Get default dependencies from an adapter
	myDeps := keep_deps.New()

	// 2. Overwrite only the specific dependency you want to change
	defaultWrite := myDeps.Write
	myDeps.Write = func(key string, value []byte) error {
		fmt.Println("writing key:", key) // custom logging on every write
		return defaultWrite(key, value)
	}

	// 3. Inject the modified dependencies into the library
	keep := keep_lib.New(myDeps)

	// Now every write performed by the database is logged
	_ = keep
}
```

4. Run your application to confirm the overwritten logic is being executed instead of the default adapter logic.
