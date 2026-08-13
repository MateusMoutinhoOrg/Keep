- Always follow the rules in [docs/References/RULES.md](/docs/References/RULES.md).

- Read README.md, open the theme index matching your goal (`docs/Index/<Theme>.md`), and check whether one of its Tutorials matches your action; if there is one, follow it.

- Before creating or editing any file, locate its specification in [docs/References/Specs.md](/docs/References/Specs.md). Never browse `docs/References/Specs/` looking for a spec — the index is the entry point.

- [sandbox/](/sandbox/) is a closed sandbox: it may not import `adapters/`, `examples/libraryExamples/`, `tests/`, a third-party module, or an OS-bound stdlib package. Every such effect is a `Deps` function field — see [docs/References/SandboxIsolation.md](/docs/References/SandboxIsolation.md).
