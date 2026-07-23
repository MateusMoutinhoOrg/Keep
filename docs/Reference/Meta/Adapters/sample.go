//go:build ignore

// This file is an illustrative sample, not part of the build.
package standard

import (
	"github.com/MateusMoutinhoOrg/Agnos/pkg/deps"
)

// New builds a deps.Deps by implementing every field of the Deps
// contract as a closure. depPropA is adapter-specific configuration.
func New(depPropA int) deps.Deps {
	return deps.Deps{
		ExampleDepFunctionA: func() int {
			return depPropA + 20
		},
	}

}
