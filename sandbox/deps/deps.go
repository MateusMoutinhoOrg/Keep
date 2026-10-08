package deps

import (
	hashdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/hashdeps"
	stddeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/stddeps"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
	stringsdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/stringsdeps"
)

// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
	HashDeps    hashdeps.Contract
	StdDeps     stddeps.Contract
	StorageDeps storagedeps.Contract
	StringsDeps stringsdeps.Contract
}
