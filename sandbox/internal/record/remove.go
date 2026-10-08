package record

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
)

// The removal procedure of the dense record pattern. Every function here
// runs under the write lock its caller holds, and takes none itself.
//
// A removal is written so that running it again after it stopped part-way —
// a crash, or a backend failing between two of its writes — finishes the job
// rather than starting another one. It never trusts a key on its own: the
// victim's back-pointer is believed only when the slot it points at agrees,
// and the last slot is believed only when the record it names points back at
// it.

// remove deletes the record id, live or not: when it is still in the list it
// is unlinked by swapping the last record into its slot, and whatever keys it
// left behind are deleted either way.
func remove(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, id int64) *api.Error {
	// Step 1: read the victim's back-pointer. Absent means already gone.
	position, found, err := dense.ReadPosition(sandbox, prefix, id)
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	if !found {
		return nil
	}

	// Step 2: read the size, dropping any dead slot a previous failure left
	// at the end of the list, so the last slot is a live record.
	size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	size, failure := trimDeadTail(sandbox, prefix, size)
	if failure != nil {
		return failure
	}

	// Steps 3 and 4: unlink the victim. Its slot is refilled when it still
	// holds the victim, and also when it is dead — a removal that stopped
	// after writing the last record into it. A slot a live record owns
	// means the victim was unlinked already.
	if position >= 1 && position <= size {
		slotID, slotLive, err := dense.SlotLive(sandbox, prefix, position)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if slotID == id || !slotLive {
			if failure := detach(sandbox, prefix, position, size); failure != nil {
				return failure
			}
		}
	}

	// Steps 5 and 6: delete the keys the record owns.
	return cleanup(sandbox, fields, prefix, scope, id)
}

// detach empties position by moving the record at the last position into it,
// then shrinks the list by one. The last slot has to be live — trimDeadTail
// guarantees it — and the writes go slot first, back-pointer second: until
// the slot is written the moved record is still live at the end, and once it
// is, the old slot at the end is dead and is dropped by the next removal if
// this one stops here.
func detach(sandbox *api.Sandbox, prefix []string, position int64, size int64) *api.Error {
	storage := sandbox.Deps.StorageDeps
	if position != size {
		lastID, found, err := dense.ReadSlot(sandbox, prefix, size)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if found {
			if err := dense.WriteInt(sandbox, dense.ListKey(sandbox, prefix, position), lastID); err != nil {
				return dense.InternalError(sandbox, err)
			}
			if err := dense.WriteInt(sandbox, dense.PositionKey(sandbox, prefix, lastID), position); err != nil {
				return dense.InternalError(sandbox, err)
			}
		}
	}
	if err := storage.Delete(dense.ListKey(sandbox, prefix, size)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	if err := dense.WriteInt(sandbox, dense.SizeKey(sandbox, prefix), size-1); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}

// trimDeadTail drops every dead slot at the end of the list — a slot left
// behind by a removal that moved its record and stopped before shrinking the
// list — and hands back the size that remains. The slot goes before the size
// shrinks, so stopping between the two leaves one more dead slot to drop next
// time, never a live record outside the list.
func trimDeadTail(sandbox *api.Sandbox, prefix []string, size int64) (int64, *api.Error) {
	trimmed := size
	for trimmed > 0 {
		_, live, err := dense.SlotLive(sandbox, prefix, trimmed)
		if err != nil {
			return 0, dense.InternalError(sandbox, err)
		}
		if live {
			break
		}
		if err := sandbox.Deps.StorageDeps.Delete(dense.ListKey(sandbox, prefix, trimmed)); err != nil {
			return 0, dense.InternalError(sandbox, err)
		}
		trimmed--
	}
	if trimmed != size {
		if err := dense.WriteInt(sandbox, dense.SizeKey(sandbox, prefix), trimmed); err != nil {
			return 0, dense.InternalError(sandbox, err)
		}
	}
	return trimmed, nil
}

// cleanup deletes every key a record that is out of the list still owns:
// its data first (cleanupData), its back-pointer last — while the
// back-pointer exists, running cleanup again picks up where it stopped.
func cleanup(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, id int64) *api.Error {
	if failure := cleanupData(sandbox, fields, prefix, scope, id); failure != nil {
		return failure
	}
	if err := sandbox.Deps.StorageDeps.Delete(dense.PositionKey(sandbox, prefix, id)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}

// cleanupData deletes everything a record owns but its back-pointer: the
// index entries that still name it, then its values and every record of its
// nested collections.
func cleanupData(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, id int64) *api.Error {
	storage := sandbox.Deps.StorageDeps

	// Step 5: drop the unique index entries the record owned.
	for _, field := range fields {
		if field.Type != api.Key {
			continue
		}
		raw, found, err := storage.Read(dense.ValueKey(sandbox, prefix, id, field.Name))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if !found {
			continue
		}
		if failure := deleteOwnEntry(sandbox, prefix, field.Name, dense.HashIndexValue(sandbox, string(raw)), id); failure != nil {
			return failure
		}
	}

	// Step 6: remove the record's own data, nested collections included.
	for _, field := range fields {
		if field.Type == api.Nested {
			nested := dense.SubPrefix(sandbox, prefix, id, field.Name)
			nestedScope := dense.NestedScope(sandbox, scope, prefix, id)
			if failure := ClearCollection(sandbox, field.Fields, nested, nestedScope); failure != nil {
				return failure
			}
			continue
		}
		if err := storage.Delete(dense.ValueKey(sandbox, prefix, id, field.Name)); err != nil {
			return dense.InternalError(sandbox, err)
		}
	}
	return nil
}

// ClearCollection removes every record of a collection, and is what a
// removal runs over each nested collection of the record it deletes. It runs
// under the write lock its caller holds.
//
// Records go from the last position backwards, so no swap is ever needed,
// and each one is emptied while it is still in the list: its data, then its
// back-pointer — which turns its slot dead — then the slot and the size. A
// clear that stops part-way and runs again therefore always finds what it
// has left to delete, at the end of the list; nothing it started on drops
// out of the list with keys still under it. A dead slot at the end is
// dropped rather than read as a record.
func ClearCollection(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope) *api.Error {
	storage := sandbox.Deps.StorageDeps
	for {
		size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if size <= 0 {
			break
		}
		id, live, err := dense.SlotLive(sandbox, prefix, size)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if live {
			if failure := cleanup(sandbox, fields, prefix, scope, id); failure != nil {
				return failure
			}
		}
		if err := storage.Delete(dense.ListKey(sandbox, prefix, size)); err != nil {
			return dense.InternalError(sandbox, err)
		}
		if err := dense.WriteInt(sandbox, dense.SizeKey(sandbox, prefix), size-1); err != nil {
			return dense.InternalError(sandbox, err)
		}
	}
	if err := storage.Delete(dense.SizeKey(sandbox, prefix)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	if err := storage.Delete(dense.LastIDKey(sandbox, prefix)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}
