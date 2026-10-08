package timesleep

import (
	"time"

	sleepdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/sleepdeps"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

// sleep fills sleepdeps.Contract.Sleep, pausing the calling goroutine for the
// given number of nanoseconds.
func sleep(nanoseconds int64) {
	if nanoseconds <= 0 {
		return
	}
	time.Sleep(time.Duration(nanoseconds))
}

// Bind fills deps.Deps.SleepDeps with the standard library's time.Sleep.
func Bind(deps *deps.Deps) {
	deps.SleepDeps = sleepdeps.Contract{
		Sleep: func(nanoseconds int64) {
			sleep(nanoseconds)
		},
	}
}
