package deps

import (
	hashdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/hashdeps"
	std "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/std"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
	stringsdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/stringsdeps"
)

// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
	Hashdeps    hashdeps.Sandbox
	Std         std.Sandbox
	Storagedeps storagedeps.Sandbox
	Stringsdeps stringsdeps.Sandbox
}
