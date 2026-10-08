package record

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// The repair of a whole collection: what Collection.Repair runs, under the
// write lock its caller holds. It restores the invariants of the dense record
// pattern after any failure the write orderings allow, reading only the keys
// the layout lets it name — positions 1 to size and ids 1 to last-id — so it
// needs no listing either.
//
// The one kind of debris it cannot reach is an index entry under a hash no
// live record holds: nothing names it. Such an entry is not valid
// (dense.IndexOwner), so it resolves to nothing and the next insert of that
// value claims it.

// Repair restores a collection and every collection nested under its live
// records:
//
//  1. every dead slot of the position list is refilled from the end, so the
//     list is dense again and every page is full;
//  2. last-id is raised to the highest id the list holds, should it be below;
//  3. every id up to last-id that is not live loses every key it left behind;
//  4. every live record gets back every index entry it is missing — which is
//     also what indexes the records of a field turned from String into Key.
//
// It keeps going past a Key value two live records share, and reports the
// first such clash as a KeyConflict at the end; everything else is repaired.
func Repair(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope) *api.Error {
	// Step 1: close every hole of the position list.
	size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	size, failure := trimDeadTail(sandbox, prefix, size)
	if failure != nil {
		return failure
	}
	for position := int64(1); position <= size; {
		_, live, err := dense.SlotLive(sandbox, prefix, position)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if live {
			position++
			continue
		}
		// The last slot is live and this one is not, so position < size.
		if failure := detach(sandbox, prefix, position, size); failure != nil {
			return failure
		}
		size, failure = trimDeadTail(sandbox, prefix, size-1)
		if failure != nil {
			return failure
		}
	}

	// Step 2: collect the live ids, and keep last-id above every one of them.
	lastID, err := dense.ReadInt(sandbox, dense.LastIDKey(sandbox, prefix))
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	liveIDs := make([]int64, 0, size)
	isLive := make(map[int64]bool, size)
	highest := lastID
	for position := int64(1); position <= size; position++ {
		id, _, err := dense.ReadSlot(sandbox, prefix, position)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		liveIDs = append(liveIDs, id)
		isLive[id] = true
		if id > highest {
			highest = id
		}
	}
	if highest > lastID {
		if err := dense.WriteInt(sandbox, dense.LastIDKey(sandbox, prefix), highest); err != nil {
			return dense.InternalError(sandbox, err)
		}
		lastID = highest
	}

	// Step 3: delete what every id that is not live left behind.
	for id := int64(1); id <= lastID; id++ {
		if isLive[id] {
			continue
		}
		has, failure := hasKeys(sandbox, fields, prefix, id)
		if failure != nil {
			return failure
		}
		if !has {
			continue
		}
		if failure := cleanup(sandbox, fields, prefix, scope, id); failure != nil {
			return failure
		}
	}

	// Step 4: give every live record back its index entries, and repair
	// every collection nested under it.
	var conflict *api.Error
	for _, id := range liveIDs {
		for _, field := range fields {
			switch field.Type {
			case api.Key:
				failure := repairEntry(sandbox, prefix, field, id)
				if failure == nil {
					continue
				}
				if failure.Type != api.KeyConflict {
					return failure
				}
				if conflict == nil {
					conflict = failure
				}
			case api.Nested:
				nested := dense.SubPrefix(sandbox, prefix, id, field.Name)
				nestedScope := dense.NestedScope(sandbox, scope, prefix, id)
				failure := Repair(sandbox, field.Fields, nested, nestedScope)
				if failure == nil {
					continue
				}
				if failure.Type != api.KeyConflict {
					return failure
				}
				if conflict == nil {
					conflict = failure
				}
			}
		}
	}
	return conflict
}

// repairEntry makes the index entry of one Key value of the live record id
// name it, unless another live record holds the same value — a KeyConflict,
// with nothing written.
func repairEntry(sandbox *api.Sandbox, prefix []string, field api.Field, id int64) *api.Error {
	raw, found, err := sandbox.Deps.StorageDeps.Read(dense.ValueKey(sandbox, prefix, id, field.Name))
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	if !found {
		return nil
	}
	hash := dense.HashIndexValue(sandbox, string(raw))
	owner, valid, err := dense.IndexOwner(sandbox, prefix, field.Name, hash)
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	if valid && owner == id {
		return nil
	}
	if valid {
		return liberror.NewWithValue(sandbox, api.KeyConflict, field.Name, string(raw),
			sandbox.Deps.StdDeps.Sprintf("records %d and %d hold the same value for key %q", owner, id, field.Name))
	}
	if err := dense.WriteInt(sandbox, dense.IndexKey(sandbox, prefix, field.Name, hash), id); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}

// hasKeys reports whether the id still owns any key at all: a back-pointer, a
// value, or a nested collection. Checking first keeps a repair of a
// collection with a long history from deleting keys that are long gone.
func hasKeys(sandbox *api.Sandbox, fields []api.Field, prefix []string, id int64) (bool, *api.Error) {
	storage := sandbox.Deps.StorageDeps
	keys := [][]string{dense.PositionKey(sandbox, prefix, id)}
	for _, field := range fields {
		if field.Type == api.Nested {
			nested := dense.SubPrefix(sandbox, prefix, id, field.Name)
			keys = append(keys, dense.SizeKey(sandbox, nested), dense.LastIDKey(sandbox, nested))
			continue
		}
		keys = append(keys, dense.ValueKey(sandbox, prefix, id, field.Name))
	}
	for _, key := range keys {
		exists, err := storage.Exists(key)
		if err != nil {
			return false, dense.InternalError(sandbox, err)
		}
		if exists {
			return true, nil
		}
	}
	return false, nil
}
