package folddeps

// This package is the sandbox's *copy* of the api a Unicode text library
// exposes — the same mechanic as hashdeps, stddeps and stringsdeps, for the
// same reason: the case-folding and normalization tables live in a library
// outside the sandbox, so it may not import them. The contract is restated
// here, and the adapter — which lives outside the sandbox — is what fills it.
//
// Keep needs it for one thing: the unique index of a Key field compares
// values without regard to case, and "without regard to case" has to mean the
// same thing for every script, not only for ASCII.

// Contract is the folding library injected whole as the Deps.FoldDeps field.
type Contract struct {
	// Fold returns the canonical caseless form of s: canonically decomposed,
	// fully case-folded, then decomposed again — Unicode's canonical caseless
	// match, NFD(toCasefold(NFD(s))). Two strings that differ only in case or
	// in how an accented letter is encoded fold to the same bytes ("STRASSE"
	// and "straße", "ΟΔΟΣ" and "οδος", "café" precomposed and decomposed). An
	// ASCII string folds to its lower-case form.
	Fold func(s string) string
}
