package api

// Info is the contract reporting the library's own identity, carried by the
// Sandbox as the field of the same name. Both values are compile-time
// constants of sandbox/internal/config, generated from
// AgnosConfig/project.yaml, so a release bump is a one-line edit touching no
// logic — and a caller can report which Keep it linked against without
// importing anything but this package.
type Info struct {
	// Name is the library's name, "Keep".
	Name func() string
	// Version is the release the caller linked against, in the "v0.0.0"
	// spelling the repository tags with.
	Version func() string
}
