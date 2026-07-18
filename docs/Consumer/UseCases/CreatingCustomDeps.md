# Creating Custom Dependencies

## Description
This guide describes how to create a custom `deps.Deps` to inject into the library, running Keep over any storage backend without using a pre-built adapter. The contract each function must honor is specified in [Required API](../RequiredApi.md).

### Rules
- Return the sentinel errors from `pkg/deps` (wrapped with `fmt.Errorf("%w: ...", ...)`) so the database layer can distinguish expected conditions from real failures.

## Workflow

1. Create a function that constructs a `deps.Deps` object. The `deps.Deps` struct is defined in [pkg/deps/deps.go](../../../pkg/deps/deps.go).
2. Provide your own implementation for each function field, following [Required API](../RequiredApi.md).
Like this:

```go
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

func createMyCustomDeps() deps.Deps {
	store := map[string][]byte{} // your storage client goes here
	return deps.Deps{
		Write: func(key string, value []byte) error {
			store[key] = value
			return nil
		},
		Read: func(key string) ([]byte, error) {
			value, found := store[key]
			if !found {
				// Sentinel errors let the lib tell "missing" from "broken"
				return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
			}
			return value, nil
		},
		// ...implement the remaining fields of deps.Deps...
	}
}

func main() {
	// 1. Create your custom deps
	myDeps := createMyCustomDeps()

	// 2. Inject them into the library
	keep := keep_lib.New(myDeps)

	// 3. Use the library normally
	_ = keep
}
```

3. Use the [native adapter](../../../adapters/native/native.go) as a complete reference implementation to copy from.
4. Run your application to ensure your custom dependencies are correctly utilized by the library.
