//go:build ignore

// This file is an illustrative sample, not part of the build.
package dense

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// Count returns the number of live records in the collection. Storage is
// reached through the deps passed in, with a single key read and no
// listing, and it references no object type — so any object can call it
// without creating an import cycle.
func Count(d deps.Deps, prefix string) (int64, api.Error) {
	size, err := ReadCount(d, SizeKey(prefix))
	if err != nil {
		return 0, InternalError(err)
	}
	return size, nil
}
