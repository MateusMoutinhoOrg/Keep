# LibFunctions Specification

## Description
Defines the required shape of a public library function in `pkg/keep/` and `pkg/database/`. A library function is pure logic that reaches storage only through the injected `Deps`.

### Rules
- Functions hang off the object they belong to: `func (l KeepLib) ...` for entry-point behavior, `func (d *KeepDatabase) ...`, `func (si *SchemaInstance) ...`, or `func (s *SchemaItem) ...` for database behavior.
- Storage is touched **only** through the injected contract — `l.Deps.<Field>()`, `d.Deps.<Field>()`, or the back-pointer of the owning object (`si.db.Deps`, `s.db.Deps`). Never construct or import a concrete implementation.
- `pkg/` must never import anything from `adapters/`.
- Every storage access must respect the invariants of the [Dense Record Pattern](/docs/Explanation/DenseRecordPattern.md): single-key reads and writes only, no key listing, and the documented write orderings.
- Expected failures are returned as a typed `*database.Error`; storage failures are wrapped as `Internal`.
- Exported functions must have a doc comment and be listed in [PublicApi.md](/docs/Reference/PublicApi.md).

## Structure
1. **Package clause**: `package keep` or `package database`.
2. **Method on the owning object**: takes its own parameters, reaches storage through the injected `Deps`, and returns the composed result plus an error where it can fail.
3. **Doc comment**: one sentence describing what the function does.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
