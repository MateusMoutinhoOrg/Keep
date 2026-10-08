package sleepdeps

// This package is the sandbox's *copy* of the one blocking primitive of the
// clock — the same mechanic as hashdeps, stddeps and stringsdeps, for the same
// reason: pausing the calling goroutine is an OS-bound effect, so `time` may
// not appear inside the sandbox. The contract is restated here, and the
// adapter — which lives outside the sandbox — is what fills it.
//
// Keep needs it for one thing: a writer that finds the write lease of a
// collection held by another process waits a little and asks again, rather
// than spinning.

// Contract is the pause library injected whole as the Deps.SleepDeps field.
type Contract struct {
	// Sleep pauses the calling goroutine for at least nanoseconds
	// nanoseconds. A zero or negative duration returns at once. The sandbox
	// may not name a `time.Duration`, so a duration crosses this boundary
	// as a plain integer.
	Sleep func(nanoseconds int64)
}
