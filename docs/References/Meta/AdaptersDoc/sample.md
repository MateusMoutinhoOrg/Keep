# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. Every adapter exposes a `New(...) deps.Deps` factory ready to be passed to [`lib.New`](/docs/References/PublicApi/lib.New.md), and honors the whole contract described in [RequiredApi.md](/docs/References/RequiredApi.md). To build a new adapter, follow [AddAdapter.md](/docs/Tutorials/AddAdapter.md).

---

## Available Adapters

| Adapter | Factory | Storage | Use When |
|---------|---------|---------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | One file per key under the working directory, path segments escaped | You want the zero-config default, with records surviving across runs |
| `frozen` | [frozen.New](/docs/References/PublicApi/frozen.New.md) | In-memory map that rejects every write | You want to guarantee a read-only consumer never mutates the database |
