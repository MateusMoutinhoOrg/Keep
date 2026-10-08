package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// The check Databases.New runs over a Props before it builds anything. A
// mistake in a schema is a mistake in the code, so it is refused once, where
// the database is declared, rather than surfacing later as a collection no
// insert can satisfy or a link that never resolves.

// Validate reports the first thing in props a database cannot be built
// from, as an InvalidSchema whose Field is the dotted path of the schema or
// field at fault ("user", "user.sessions.token"); nil when there is none.
func Validate(sandbox *api.Sandbox, props api.Props) *api.Error {
	for _, segment := range dense.PathSegments(sandbox, props.Path) {
		if segment == "." || segment == ".." {
			return invalidSchema(sandbox, "",
				sandbox.Deps.StdDeps.Sprintf("Path %q holds a %q segment: give the backend the base directory instead", props.Path, segment))
		}
	}
	seen := map[string]bool{}
	for _, schema := range props.Schemas {
		if schema.Name == "" {
			return invalidSchema(sandbox, "", "a schema has no Name")
		}
		if seen[schema.Name] {
			return invalidSchema(sandbox, schema.Name,
				sandbox.Deps.StdDeps.Sprintf("schema %q is declared twice", schema.Name))
		}
		seen[schema.Name] = true
	}
	for _, schema := range props.Schemas {
		if failure := validateFields(sandbox, props, schema.Name, schema.Fields); failure != nil {
			return failure
		}
	}
	return nil
}

// validateFields checks the fields of one schema, or of one Nested field,
// whose dotted path is path.
func validateFields(sandbox *api.Sandbox, props api.Props, path string, fields []api.Field) *api.Error {
	seen := map[string]bool{}
	for _, field := range fields {
		at := path + "." + field.Name
		if field.Name == "" {
			return invalidSchema(sandbox, path,
				sandbox.Deps.StdDeps.Sprintf("a field of %q has no Name", path))
		}
		if seen[field.Name] {
			return invalidSchema(sandbox, at,
				sandbox.Deps.StdDeps.Sprintf("field %q is declared twice", at))
		}
		seen[field.Name] = true

		if field.Type == 0 {
			return invalidSchema(sandbox, at,
				sandbox.Deps.StdDeps.Sprintf("field %q declares no Type", at))
		}
		if field.Type < api.Key || field.Type > api.Bytes {
			return invalidSchema(sandbox, at,
				sandbox.Deps.StdDeps.Sprintf("field %q declares Type %d, which is not a field type", at, field.Type))
		}
		if field.Type != api.Link && field.Target != "" {
			return invalidSchema(sandbox, at,
				sandbox.Deps.StdDeps.Sprintf("field %q declares a Target, and only a Link field has one", at))
		}
		if field.Type != api.Nested && len(field.Fields) > 0 {
			return invalidSchema(sandbox, at,
				sandbox.Deps.StdDeps.Sprintf("field %q declares Fields, and only a Nested field has them", at))
		}

		switch field.Type {
		case api.Link:
			if field.Target == "" {
				return invalidSchema(sandbox, at,
					sandbox.Deps.StdDeps.Sprintf("link field %q declares no Target", at))
			}
			if !declaresSchema(sandbox, props, field.Target) {
				return invalidSchema(sandbox, at,
					sandbox.Deps.StdDeps.Sprintf("link field %q targets %q, which is not a schema of the Props", at, field.Target))
			}
		case api.Nested:
			if dense.ReservedNestedName(sandbox, field.Name) {
				return invalidSchema(sandbox, at,
					sandbox.Deps.StdDeps.Sprintf("nested field %q takes a reserved name: a record already keeps its %q key", at, field.Name))
			}
			if failure := validateFields(sandbox, props, at, field.Fields); failure != nil {
				return failure
			}
		}
	}
	return nil
}

// declaresSchema reports whether props holds a schema of the given name.
func declaresSchema(sandbox *api.Sandbox, props api.Props, name string) bool {
	for _, schema := range props.Schemas {
		if schema.Name == name {
			return true
		}
	}
	return false
}

// invalidSchema builds the failure Validate reports.
func invalidSchema(sandbox *api.Sandbox, field string, message string) *api.Error {
	return liberror.New(sandbox, api.InvalidSchema, field, message)
}
