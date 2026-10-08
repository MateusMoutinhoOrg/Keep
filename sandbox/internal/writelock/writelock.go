package writelock

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// The write lock of a top-level collection. Every operation that writes to a
// collection — or to any collection nested under one of its records — holds
// it from its first read to its last write, so the write orderings of the
// dense record pattern never interleave with another writer's.
//
// It is two locks taken one after the other:
//
//  1. an in-process lock, one per top-level collection per sandbox, so two
//     goroutines of one program hand the collection over without polling;
//  2. a storage lease, StorageDeps.Lock on the collection's own prefix, so a
//     writer in another process sharing the backend waits too. A backend with
//     no leases reports locked == true at once, and then only the first lock
//     protects anything.
//
// A lease expires on its own after LeaseSeconds, which is what frees a
// collection whose writer crashed holding it. The storage contract's Unlock
// names no holder, so an operation that outlives its lease could release a
// lease another writer has since taken: LeaseSeconds is set far above what
// one operation costs.
//
// Readers take no lock: every read of the record package tolerates a writer
// working beside it.

// LeaseSeconds is how long a storage lease is taken for. It is the longest a
// crashed writer can keep a collection locked.
const LeaseSeconds = 60

// waitNanoseconds is how long a writer waits for a lease held elsewhere
// before giving up with Internal: two leases' worth, so a lease left behind
// by a crash always expires within it.
const waitNanoseconds = int64(2 * LeaseSeconds * 1000 * 1000 * 1000)

// firstBackoffNanoseconds and maxBackoffNanoseconds bound the pause between
// two attempts at a lease held by another process.
const (
	firstBackoffNanoseconds = int64(1000 * 1000)
	maxBackoffNanoseconds   = int64(100 * 1000 * 1000)
)

// Registry holds the in-process lock of every top-level collection one
// sandbox has written to. It is built once, by the Databases the sandbox
// carries, and shared by every database and collection built from it. Its
// locks are buffered channels — a builtin — so the sandbox needs no `sync`.
type Registry struct {
	// guard is the lock of the locks map itself.
	guard chan struct{}
	// locks maps the flattened prefix of a top-level collection to its lock.
	locks map[string]chan struct{}
}

// NewRegistry builds an empty registry. A lock is created the first time its
// collection is written to.
func NewRegistry(sandbox *api.Sandbox) *Registry {
	return &Registry{
		guard: make(chan struct{}, 1),
		locks: map[string]chan struct{}{},
	}
}

// Acquire takes the write lock of the top-level collection whose prefix is
// root, and hands back the function that releases it. It waits for as long
// as another goroutine of this sandbox holds the collection, then for at
// most two leases' worth for a writer in another process. A backend failure
// or a lease that never frees reports Internal, and nothing is held.
func Acquire(sandbox *api.Sandbox, registry *Registry, root []string) (func(), *api.Error) {
	local := lockOf(sandbox, registry, root)
	local <- struct{}{}

	deadline := sandbox.Deps.StdDeps.Now() + waitNanoseconds
	backoff := firstBackoffNanoseconds
	for {
		locked, err := sandbox.Deps.StorageDeps.Lock(root, LeaseSeconds)
		if err != nil {
			<-local
			return nil, liberror.New(sandbox, api.Internal, "", err.Error())
		}
		if locked {
			break
		}
		if sandbox.Deps.StdDeps.Now() >= deadline {
			<-local
			return nil, liberror.New(sandbox, api.Internal, "",
				sandbox.Deps.StdDeps.Sprintf("keep: timed out waiting for the write lock of %q",
					sandbox.Deps.StringsDeps.Join(root, "/")))
		}
		sandbox.Deps.SleepDeps.Sleep(backoff)
		backoff *= 2
		if backoff > maxBackoffNanoseconds {
			backoff = maxBackoffNanoseconds
		}
	}

	release := func() {
		// A lease that cannot be released expires on its own; the writes it
		// guarded have already succeeded or failed on their own terms.
		_ = sandbox.Deps.StorageDeps.Unlock(root)
		<-local
	}
	return release, nil
}

// lockOf returns the in-process lock of the collection at root, creating it
// on first use.
func lockOf(sandbox *api.Sandbox, registry *Registry, root []string) chan struct{} {
	name := flatten(sandbox, root)
	registry.guard <- struct{}{}
	defer func() { <-registry.guard }()
	lock, ok := registry.locks[name]
	if !ok {
		lock = make(chan struct{}, 1)
		registry.locks[name] = lock
	}
	return lock
}

// flatten turns a prefix into the one string the registry indexes it by.
// Every segment is preceded by its length, so no character of a segment can
// be read back as a boundary and two prefixes never share a lock by
// accident.
func flatten(sandbox *api.Sandbox, prefix []string) string {
	parts := make([]string, 0, len(prefix))
	for _, segment := range prefix {
		parts = append(parts, sandbox.Deps.StringsDeps.FormatInt(int64(len(segment)), 10)+":"+segment)
	}
	return sandbox.Deps.StringsDeps.Join(parts, "")
}
