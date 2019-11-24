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
	"strings"
	"testing"

	alg "github.com/lamg/algorithms"
	"github.com/stretchr/testify/require"
)

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

func TestTail(t *testing.T) {
	// '⇒' id {'⇒' id}.
	g := []Symbol{
		{},                                   // 0
		Empty,                                // 1
		{Name: junction, Header: 3, Next: 4}, // 2
		{Name: ImpliesOp, IsTerminal: true},  // 3
		{Name: "", Header: 5, Next: 6},       // 4
		{Name: Identifier, IsTerminal: true}, // 5
		{Name: ImpliesOp, IsTerminal: true, Next: 4, Alt: 1}, //6
	}
	ss := []Scanner{
		IdentScan,
		SpaceScan,
		StrScan(ImpliesOp),
	}
	ts := []testParse{
		{
			p: "⇒ a ⇒ b",
			t: &Tree{
				Token: fixedString(ImpliesOp),
				Children: []*Tree{
					IdentTree("a"),
					TermTree(ImpliesOp),
					IdentTree("b"),
				},
			},
		},
		{
			p: "⇒ a",
			t: &Tree{
				Token:    fixedString(ImpliesOp),
				Children: []*Tree{IdentTree("a")},
			},
		},
		{
			p: "a",
			e: &ExpectingErr{
				Expected: ImpliesOp,
				Actual:   Identifier,
			},
		},
		{
			p: "⇒",
			e: &UnexpectedEOFErr{Expected: Identifier},
		},
		{
			p: "⇒ a ⇒",
			e: &UnexpectedEOFErr{Expected: Identifier},
		},
	}
	for i, j := range ts {
		rd := strings.NewReader(j.p)
		tk := NewReaderTokens(rd, ss)
		r, e := Parse(g, tk)
		require.Equal(t, j.e, e)
		require.Equal(t, j.t, r, "At %d: %s", i, j.p)
	}
}
