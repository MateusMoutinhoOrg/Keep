package info

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	config "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/config"
)

// NameFactory fills api.Info.Name.
func NameFactory(sandbox *api.Sandbox, info *api.Info) func() string {
	return func() string {
		return sandbox.Deps.Stringsdeps.TrimSpace(config.ProjectName)
	}
}

// VersionFactory fills api.Info.Version. The constant is generated from
// AgnosConfig/project.yaml, and is trimmed of whatever whitespace it is
// written with so a caller can compare it as it stands.
func VersionFactory(sandbox *api.Sandbox, info *api.Info) func() string {
	return func() string {
		return sandbox.Deps.Stringsdeps.TrimSpace(config.Version)
	}
}

// NewInfo builds the api.Info contract, running every factory over it to
// fill its function fields. Adding a function field to api.Info means
// adding its factory call here.
func NewInfo(sandbox *api.Sandbox) api.Info {
	info := api.Info{}
	info.Name = NameFactory(sandbox, &info)
	info.Version = VersionFactory(sandbox, &info)
	return info
}
