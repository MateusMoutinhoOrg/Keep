package deps

import (
	folddeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/folddeps"
	hashdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/hashdeps"
	sleepdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/sleepdeps"
	stddeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/stddeps"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
	stringsdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/stringsdeps"
)

// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
	FoldDeps    folddeps.Contract
	HashDeps    hashdeps.Contract
	SleepDeps   sleepdeps.Contract
	StdDeps     stddeps.Contract
	StorageDeps storagedeps.Contract
	StringsDeps stringsdeps.Contract
}
