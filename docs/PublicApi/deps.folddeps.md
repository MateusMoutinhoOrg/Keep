# `deps.FoldDeps`

`sandbox/deps/folddeps`

## `Contract`

Contract is the folding library injected whole as the Deps.FoldDeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Fold` | `func(s string) string` | Fold returns the canonical caseless form of s: canonically decomposed, fully case-folded, then decomposed again — Unicode's canonical caseless match, NFD(toCasefold(NFD(s))). Two strings that differ only in case or in how an accented letter is encoded fold to the same bytes ("STRASSE" and "straße", "ΟΔΟΣ" and "οδος", "café" precomposed and decomposed). An ASCII string folds to its lower-case form. |

[every contract](doc.md)
