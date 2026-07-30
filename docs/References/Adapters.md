# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. Every adapter exposes a `New(...) deps.Deps` factory ready to be passed to [`lib.New`](/docs/References/PublicApi/lib.New.md), and honors the whole contract described in [RequiredApi.md](/docs/References/RequiredApi.md). To build a new adapter, follow [AddAdapter.md](/docs/Tutorials/AddAdapter.md).

---

## Available Adapters

| Adapter | Factory | Storage | Use When |
|---------|---------|---------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | One file per key on the filesystem; path segments split on `/` and escaped, so keys may hold arbitrary characters. `NewWithBase` picks the root directory | You want the zero-config default, with records surviving across process restarts |
| `native` | [native.New](/docs/References/PublicApi/native.New.md) | An in-memory map guarded by a mutex; data lives only for the lifetime of the process | You want the fastest backend and don't need records after the process exits — tests, prototypes, and throwaway runs |
