//go:build ignore

// This file is an illustrative sample, not part of the build.
package database

// Count returns the number of live records in the collection. Storage is
// reached through the injected Deps carried by the owning database, with
// a single key read and no listing.
func (si *SchemaInstance) Count() (int64, *Error) {
	size, err := readCount(si.db.Deps, sizeKey(si.prefix))
	if err != nil {
		return 0, internalError(err)
	}
	return size, nil
}
