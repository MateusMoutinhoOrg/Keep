- Always follow Rules in [docs/References/RULES.md](/docs/References/RULES.md)

- Read README.md and verify if there is no Tutorial that matches your action, if it is, follow it.

- Before creating or editing any file, locate its specification in [docs/References/Specs.md](/docs/References/Specs.md). Never browse `docs/References/Meta/` looking for a spec — the index is the entry point.

- [sandbox/](/sandbox/) is a closed sandbox: it may not import `adapters/`, `examples/`, `tests/`, a third-party module, or an OS-bound stdlib package. Every such effect is a `Deps` function field — see [docs/Explanations/SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
