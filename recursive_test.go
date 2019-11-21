package parsing

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type parseTest struct {
	name    string
	grammar []Symbol
	tks     []*Token
	r       *Tree
	err     error
}

var aTimesB = &parseTest{
	name:    "a*b",
	grammar: basicGrammar,
	tks: []*Token{
		identifier("a"),
		fixedString(times),
		identifier("b"),
	},
	r: &Tree{
		Token: identifier("a"),
		Children: []*Tree{
			TermTree(times),
			idt("b"),
		},
	},
}

var aPlusB = &parseTest{
	name:    "a+b",
	grammar: basicGrammar,
	tks: []*Token{
		identifier("a"),
		fixedString(plus),
		identifier("b"),
	},
	r: &Tree{
		Token: identifier("a"),
		Children: []*Tree{
			TermTree(plus),
			idt("b"),
		},
	},
}

var aPlusBTimesC = &parseTest{
	name:    "a+b*c",
	grammar: basicGrammar,
	tks: []*Token{
		identifier("a"),
		fixedString(plus),
		identifier("b"),
		fixedString(times),
		identifier("c"),
	},
	r: &Tree{
		Token: identifier("a"),
		Children: []*Tree{
			TermTree(plus),
			{
				Token: identifier("b"),
				Children: []*Tree{
					TermTree(times),
					idt("c"),
				},
			},
		},
	},
}

var aMinusBPlusC = &parseTest{
	name:    "a-b+c",
	grammar: basicGrammar,
	tks: []*Token{
		identifier("a"),
		fixedString(minus),
		identifier("b"),
		fixedString(plus),
		identifier("c"),
	},
	r: &Tree{
		Token: identifier("a"),
		Children: []*Tree{
			TermTree(minus),
			idt("b"),
			TermTree(plus),
			idt("c"),
		},
	},
}

var aPlusPar = &parseTest{
	name:    "a+(b+c-d)",
	grammar: basicGrammar,
	tks: []*Token{
		identifier("a"),
		fixedString(plus),
		fixedString(opar),
		identifier("b"),
		fixedString(plus),
		identifier("c"),
		fixedString(minus),
		identifier("d"),
		fixedString(cpar),
	},
	r: &Tree{
		Token: identifier("a"),
		Children: []*Tree{
			TermTree(plus),
			{
				Token: fixedString(opar),
				Children: []*Tree{
					{
						Token: identifier("b"),
						Children: []*Tree{
							TermTree(plus),
							IdentTree("c"),
							TermTree(minus),
							IdentTree("d"),
						},
					},
					TermTree(cpar),
				},
			},
		},
	},
}

var aPar = &parseTest{
	name:    "(a)",
	grammar: basicGrammar,
	tks: []*Token{
		fixedString(opar),
		identifier("a"),
		fixedString(cpar),
	},
	r: &Tree{
		Token: fixedString(opar),
		Children: []*Tree{
			IdentTree("a"),
			TermTree(cpar),
		},
	},
}

var eofError = &parseTest{
	name:    "EOF found",
	grammar: basicGrammar,
	tks:     []*Token{fixedString(opar)},
	err:     &UnexpectedEOFErr{Expected: opar},
}

var expectingExpr = &parseTest{
	name:    "Expecting expression",
	grammar: basicGrammar,
	tks:     []*Token{fixedString("X")},
	r:       nil,
	err: &ExpectingErr{
		Expected: opar,
		Actual:   "X",
	},
}

var remainingToken = &parseTest{
	name:    "remaining token",
	grammar: basicGrammar,
	tks:     []*Token{identifier("a"), fixedString("X")},
	err:     &RemainingTokenErr{Token: fixedString("X")},
}

func TestT0(t *testing.T) {
	tests := []*parseTest{
		aTimesB,
		aPlusB,
		aPlusBTimesC,
		aMinusBPlusC,
		aPlusPar,
		aPar,
		eofError,
		expectingExpr,
		remainingToken,
	}

	for i, j := range tests {
		tkf := &tkStream{tokens: j.tks, n: -1}
		r, e := Parse(j.grammar, tkf)
		require.Equal(t, j.err, e, "At %d, %s: %v != %v",
			i, j.name, j.err, e)
		require.Equal(t, j.r, r, "At %d, %s", i, j.name)
	}
}

func TestT1(t *testing.T) {
	g := []Symbol{
		{},
		Empty,
		{Name: factor, Header: 3},
		{Name: Identifier, IsTerminal: true},
	}
	tks := []*Token{
		{Name: "id", Value: "a"},
	}
	tkf := &tkStream{tokens: tks, n: -1}
	r, e := Parse(g, tkf)
	require.NoError(t, e)
	rt := idt("a")
	require.Equal(t, rt, r)
}

type tkStream struct {
	tokens []*Token
	n      int
}

func (s *tkStream) Current() (t *Token, e error) {
	if s.n == len(s.tokens) {
		e = io.EOF
	} else {
		t = s.tokens[s.n]
	}
	return
}

func (s *tkStream) Next() {
	if s.n != len(s.tokens) {
		s.n = s.n + 1
	}
}

const (
	expression = "expression"
	factor     = "factor"
	term       = "term"
	plus       = "+"
	minus      = "-"
	times      = "*"
	div        = "/"
	opar       = "("
	cpar       = ")"
)

// expression = term {("+"|"-") term}.
// term = factor {("*"|"/") factor}.
// factor = id | "(" expression ")".

var basicGrammar = []Symbol{
	{},                             // 0
	Empty,                          // 1
	{Name: expression, Header: 3},  // 2
	{Name: "", Header: 4, Next: 5}, // 3
	{Name: term, Header: 6},        // 4
	{Name: plus, IsTerminal: true, Next: 3, Alt: 7},   // 5
	{Name: "", Header: 8, Next: 9},                    // 6
	{Name: minus, IsTerminal: true, Next: 3, Alt: 1},  // 7
	{Name: factor, Header: 10},                        // 8
	{Name: times, IsTerminal: true, Next: 6, Alt: 11}, // 9
	{Name: Identifier, IsTerminal: true, Alt: 12},     // 10
	{Name: div, IsTerminal: true, Next: 6, Alt: 1},    // 11
	{Name: opar, IsTerminal: true, Next: 13},          // 12
	{Name: "", Header: 2, Next: 14},                   // 13
	{Name: cpar, IsTerminal: true},                    // 14
}

func fixedString(s string) *Token {
	return &Token{Name: s, Value: s}
}

func identifier(s string) *Token {
	return &Token{Name: Identifier, Value: s}
}

func idt(s string) *Tree {
	return &Tree{Token: &Token{Name: Identifier, Value: s}}
}

func termSymbol(s string) *Tree {
	return &Tree{Token: &Token{Name: s, Value: s}}
}
