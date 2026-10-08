package xtextfold

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"

	folddeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/folddeps"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

// fold fills folddeps.Contract.Fold with Unicode's canonical caseless match,
// NFD(toCasefold(NFD(s))). A cases.Caser is not safe for concurrent use, so
// every call builds its own.
func fold(s string) string {
	decomposed := norm.NFD.String(s)
	folded := cases.Fold().String(decomposed)
	return norm.NFD.String(folded)
}

// Bind fills deps.Deps.FoldDeps with golang.org/x/text's full case folding
// and canonical decomposition.
func Bind(deps *deps.Deps) {
	deps.FoldDeps = folddeps.Contract{
		Fold: func(s string) string {
			return fold(s)
		},
	}
}
