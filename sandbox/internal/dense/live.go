package dense

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// What makes a record, a list slot and an index entry live. A write that
// stops part-way — a crash, or a backend failing in the middle of an
// operation — leaves keys behind that look like data; these three readings
// are what tell the data from the debris, so every operation that trusts a
// key reads it through one of them. They are reads only, and cost a fixed
// number of them.
//
// Every key holding an id is read tolerantly: a key that is absent or holds
// something that is not a decimal id reads as "nothing here", never as id 0
// and never as a failure. Only a backend error is one.

// ReadPosition reads a record's position back-pointer. found is false when
// the key is absent or unreadable.
func ReadPosition(sandbox *api.Sandbox, prefix []string, id int64) (position int64, found bool, err error) {
	return readID(sandbox, PositionKey(sandbox, prefix, id))
}

// ReadSlot reads the id stored at one position of the dense list. found is
// false when the slot is absent or unreadable.
func ReadSlot(sandbox *api.Sandbox, prefix []string, position int64) (id int64, found bool, err error) {
	return readID(sandbox, ListKey(sandbox, prefix, position))
}

// ReadIndex reads the id one index entry names, live or not. found is false
// when the entry is absent or unreadable.
func ReadIndex(sandbox *api.Sandbox, prefix []string, field string, hash string) (id int64, found bool, err error) {
	return readID(sandbox, IndexKey(sandbox, prefix, field, hash))
}

// SlotLive reads the slot at position and reports whether it is live: it
// holds an id whose back-pointer points at this very position. A slot that
// is absent, unreadable or names a record that points elsewhere is dead —
// what a removal that stopped half-way through its swap leaves. It does not
// compare position with size: a caller walks positions 1 to size.
func SlotLive(sandbox *api.Sandbox, prefix []string, position int64) (id int64, live bool, err error) {
	id, found, err := ReadSlot(sandbox, prefix, position)
	if err != nil || !found {
		return 0, false, err
	}
	back, found, err := ReadPosition(sandbox, prefix, id)
	if err != nil || !found {
		return id, false, err
	}
	return id, back == position, nil
}

// Live reports whether id names a live record: its back-pointer exists, lies
// in [1, size], and the list slot there holds id back. A record whose insert
// never reached its commit point, and one whose removal got past its first
// write, are not live, whatever keys they left behind. Three reads.
func Live(sandbox *api.Sandbox, prefix []string, id int64) (bool, error) {
	position, found, err := ReadPosition(sandbox, prefix, id)
	if err != nil || !found {
		return false, err
	}
	if position < 1 {
		return false, nil
	}
	size, err := ReadInt(sandbox, SizeKey(sandbox, prefix))
	if err != nil {
		return false, err
	}
	if position > size {
		return false, nil
	}
	slot, found, err := ReadSlot(sandbox, prefix, position)
	if err != nil || !found {
		return false, err
	}
	return slot == id, nil
}

// IndexOwner reads the index entry of one value of one Key field and
// reports whether it is valid: it names a live record whose current value of
// that field hashes back to this very entry. An entry that is not valid is
// free — an insert or an update claims it by writing over it, and a lookup
// through it finds nothing. That is what makes an entry left behind by a
// failed write harmless: it can neither resolve to the wrong record nor block
// the next writer.
func IndexOwner(sandbox *api.Sandbox, prefix []string, field string, hash string) (id int64, valid bool, err error) {
	id, found, err := ReadIndex(sandbox, prefix, field, hash)
	if err != nil || !found {
		return 0, false, err
	}
	live, err := Live(sandbox, prefix, id)
	if err != nil || !live {
		return id, false, err
	}
	raw, found, err := sandbox.Deps.StorageDeps.Read(ValueKey(sandbox, prefix, id, field))
	if err != nil || !found {
		return id, false, err
	}
	return id, HashIndexValue(sandbox, string(raw)) == hash, nil
}

// readID reads a key holding one decimal id. found is false when the key is
// absent or does not hold one.
func readID(sandbox *api.Sandbox, key []string) (int64, bool, error) {
	raw, found, err := sandbox.Deps.StorageDeps.Read(key)
	if err != nil || !found {
		return 0, false, err
	}
	id, err := ParseID(sandbox, raw)
	if err != nil {
		return 0, false, nil
	}
	return id, true, nil
}
