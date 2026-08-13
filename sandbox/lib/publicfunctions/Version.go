package publicfunctions

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/config"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

// VersionFactory fills api.Lib.Version with a closure reporting the release
// the consumer linked against. The number itself lives in its own file under
// sandbox/config, so a release bump is a one-line edit touching no logic; the
// constant is trimmed of any surrounding whitespace it is written with.
func VersionFactory(l *api.Lib) func() string {
	return func() string {
		return strings.TrimSpace(config.Version)
	}
}
