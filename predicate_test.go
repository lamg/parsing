// Copyright © 2019 Luis Ángel Méndez Gort

// This file is part of Predicate.

// Predicate is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General
// Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your
// option) any later version.

// Predicate is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the
// implied warranty of MERCHANTABILITY or FITNESS FOR A
// PARTICULAR PURPOSE. See the GNU Lesser General Public
// License for more details.

// You should have received a copy of the GNU Lesser General
// Public License along with Predicate.  If not, see
// <https://www.gnu.org/licenses/>.

package parsing

import (
	"testing"

	alg "github.com/lamg/algorithms"
	"github.com/stretchr/testify/require"
)

func TestFuncPointer(t *testing.T) {
	var f func()
	g := func() {
		f()
	}
	a := false
	f = func() { a = true }
	g()
	require.True(t, a)
}

type testParse struct {
	p string
	t *Tree
	e error
}

var empty = &testParse{
	p: "",
	e: &UnexpectedEOFErr{
		Expected: opar,
	},
}

var trueAndFalse = &testParse{
	p: "true ∧ false",
}

var trueAnd = &testParse{
	p: "true ∧",
	e: &UnexpectedEOFErr{Expected: opar},
	t: &Tree{},
}

var notA = &testParse{
	p: "¬A",
	t: &Tree{},
}

var notAAndBOrC = &testParse{
	p: "¬A ∧ (B ∨ C)",
	t: &Tree{},
}

var aOrNotBAndC = &testParse{
	p: "A ∨ ¬(B ∧ C)",
	t: &Tree{},
}

func TestParse(t *testing.T) {
	ps := []*testParse{
		empty,
		trueAndFalse,
		trueAnd,
		notA,
		notAAndBOrC,
		aOrNotBAndC,
		//{"A ≡ B ≢ ¬C ⇒ D", nil},
		//{"A ≡ B ≡ ¬C ⇐ D", nil},
		//{"A ≡ B ≡ ¬(C ⇐ D)", nil},
		//{"A ∨ B ∨ C", nil},
		//{"A ∨ B ∧ C", errAnd},
		//{"A ⇒ B ⇐ C", errFl},
		//{"A ∨ (B ∧ C)", nil},
		//{"A ⇒ (B ⇐ C)", nil},
		//{"a ≡ b ≢ c ≡ ¬x ∧ (¬z ≡ y) ≢ true", nil},
	}
	inf := func(i int) {
		println("pred:", ps[i].p)
		np, e := ParsePred(ps[i].p)
		require.Equal(t, ps[i].e, e,
			"At %d: %s", i, ps[i].p)
		if e == nil {
			require.Equal(t, ps[i].t, np, "At %d: %s", i, ps[i].p)
		}
	}
	alg.Forall(inf, len(ps))
}
