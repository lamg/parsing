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
	t: opTree(AndOp, "true", "false"),
}

var trueAnd = &testParse{
	p: "true ∧",
	e: &UnexpectedEOFErr{Expected: opar},
	t: &Tree{
		Value: pred,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value: junction,
						Children: []*Tree{
							factId("true"),
						},
					},
				},
			},
		},
	},
}

var notA = &testParse{
	p: "¬A",
	t: &Tree{
		Value: pred,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value:    junction,
						Children: []*Tree{notFactor("A")},
					},
				},
			},
		},
	},
}

var notAAndBOrC = &testParse{
	p: "¬A ∧ (B ∨ C)",
	t: &Tree{
		Value: pred,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value: junction,
						Children: []*Tree{
							notFactor("A"),
							{
								Value: AndOp,
								Token: fixedString(AndOp),
								Children: []*Tree{
									{
										Value: factor,
										Children: []*Tree{
											TermTree(opar),
											opTree(OrOp, "B", "C"),
											TermTree(cpar),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var aOrNotBAndC = &testParse{
	p: "A ∨ ¬(B ∧ C)",
	t: &Tree{
		Value: pred,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value: junction,
						Children: []*Tree{
							factId("A"),
							{
								Value: OrOp,
								Token: fixedString(OrOp),
								Children: []*Tree{
									{
										Value: factor,
										Children: []*Tree{
											{
												Value: NotOp,
												Token: fixedString(NotOp),
												Children: []*Tree{
													TermTree(opar),
													opTree(AndOp, "B", "C"),
													TermTree(cpar),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
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

func factId(s string) *Tree {
	return &Tree{
		Value:    factor,
		Children: []*Tree{IdentTree(s)},
	}
}

func opTree(op, a, b string) *Tree {
	return &Tree{
		Value: pred,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value: junction,
						Children: []*Tree{
							factId(a),
							&Tree{
								Value: op,
								Token: fixedString(op),
								Children: []*Tree{
									factId(b),
								},
							},
						},
					},
				},
			},
		},
	}
}

func notFactor(id string) *Tree {
	return &Tree{
		Value: factor,
		Children: []*Tree{
			{
				Value: NotOp,
				Token: fixedString(NotOp),
				Children: []*Tree{
					IdentTree(id),
				},
			},
		},
	}
}

func TestTailSymbol(t *testing.T) {
	s := &Symbol{
		Name:       id,
		IsTerminal: true,
	}
	ts := tailSymbol(AndOp, s)
	ss := []Scanner{IdentScan, SpaceScan, StrScan(AndOp)}
	tests := []testParse{
		{
			p: "∧ b ∧ a",
			t: &Tree{
				Value: tail,
				Children: []*Tree{
					tailElem(AndOp, "b"),
					tailElem(AndOp, "a"),
				},
			},
		},
		{
			p: "∧ b ∧",
			e: &UnexpectedEOFErr{Expected: id},
		},
		{},
		{
			p: "a",
			e: &RemainingTokenErr{identifier("a")},
		},
	}
	testCases(t, ts, ss, tests)
}

func TestOpTail(t *testing.T) {
	s := &Symbol{
		Name:       id,
		IsTerminal: true,
	}
	ts := optTail(factor, AndOp, OrOp, s)
	ss := []Scanner{
		IdentScan,
		SpaceScan,
		StrScan(AndOp),
		StrScan(OrOp),
	}
	tests := []testParse{
		{
			p: "p ∧ b",
			t: &Tree{
				Value: factor,
				Children: []*Tree{
					idt("p"),
					TermTree(AndOp),
					idt("b"),
				},
			},
		},
		{
			p: "p ∧ b ∧ b",
			t: &Tree{
				Value: factor,
				Children: []*Tree{
					idt("p"),
					TermTree(AndOp),
					idt("b"),
					TermTree(AndOp),
					idt("b"),
				},
			},
		},
		{
			p: "p ∨ b",
			t: &Tree{
				Value: factor,
				Children: []*Tree{
					idt("p"),
					TermTree(OrOp),
					idt("b"),
				},
			},
		},
		{
			p: "p ∨ b ∧ b",
			t: &Tree{
				Value: factor,
				Children: []*Tree{
					idt("p"),
					TermTree(OrOp),
					idt("b"),
				},
			},
			e: &RemainingTokenErr{Token: fixedString(AndOp)},
		},
	}
	testCases(t, ts, ss, tests)
}

func testCases(t *testing.T, ts *Symbol, ss []Scanner,
	tests []testParse) {
	inf := func(i int) {
		tks := NewReaderTokens(strings.NewReader(tests[i].p), ss)
		r, e := Parse(ts, tks)
		require.Equal(t, tests[i].e, e, "At %d: %s, %v ≠ %v",
			i, tests[i].p, tests[i].e, e)
		require.Equal(t, tests[i].t, r, "At %d: %s", i, tests[i].p)
	}
	alg.Forall(inf, len(tests))
}

func tailElem(op, id string) *Tree {
	return &Tree{
		Value:    op,
		Token:    fixedString(op),
		Children: []*Tree{IdentTree(id)},
	}
}
