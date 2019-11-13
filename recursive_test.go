package parsing

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type parseTest struct {
	name    string
	grammar *Symbol
	tks     []*Token
	r       *Tree
	err     error
}

var aTimesB = &parseTest{
	name:    "a*b",
	grammar: basicGrammar(),
	tks: []*Token{
		identifier("a"),
		fixedString(times),
		identifier("b"),
	},
	r: &Tree{
		Value: expression,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{Value: factor, Children: []*Tree{idt("a")}},
					termSymbol(times),
					{Value: factor, Children: []*Tree{idt("b")}},
				},
			},
		},
	},
}

var aPlusB = &parseTest{
	name:    "a+b",
	grammar: basicGrammar(),
	tks: []*Token{
		identifier("a"),
		fixedString(plus),
		identifier("b"),
	},
	r: &Tree{
		Value: expression,
		Children: []*Tree{
			termIdt("a"),
			termSymbol(plus),
			termIdt("b"),
		},
	},
}

var aPlusBTimesC = &parseTest{
	name:    "a+b*c",
	grammar: basicGrammar(),
	tks: []*Token{
		identifier("a"),
		fixedString(plus),
		identifier("b"),
		fixedString(times),
		identifier("c"),
	},
	r: &Tree{
		Value: expression,
		Children: []*Tree{
			termIdt("a"),
			termSymbol(plus),
			{
				Value: term,
				Children: []*Tree{
					{Value: factor, Children: []*Tree{idt("b")}},
					termSymbol(times),
					{Value: factor, Children: []*Tree{idt("c")}},
				},
			},
		},
	},
}

var aMinusBPlusC = &parseTest{
	name:    "a-b+c",
	grammar: basicGrammar(),
	tks: []*Token{
		identifier("a"),
		fixedString(minus),
		identifier("b"),
		fixedString(plus),
		identifier("c"),
	},
	r: &Tree{
		Value: expression,
		Children: []*Tree{
			termIdt("a"),
			termSymbol(minus),
			termIdt("b"),
			termSymbol(plus),
			termIdt("c"),
		},
	},
}

var aPlusPar = &parseTest{
	name:    "a+(b+c-d)",
	grammar: basicGrammar(),
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
		Value: expression,
		Children: []*Tree{
			termIdt("a"),
			termSymbol(plus),
			{
				Value: term,
				Children: []*Tree{
					{
						Value: factor,
						Children: []*Tree{
							termSymbol(opar),
							{
								Value: expression,
								Children: []*Tree{
									termIdt("b"),
									termSymbol(plus),
									termIdt("c"),
									termSymbol(minus),
									termIdt("d"),
								},
							},
							termSymbol(cpar),
						},
					},
				},
			},
		},
	},
}

var aPar = &parseTest{
	name:    "(a)",
	grammar: basicGrammar(),
	tks: []*Token{
		fixedString(opar),
		identifier("a"),
		fixedString(cpar),
	},
	r: &Tree{
		Value: expression,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value: factor,
						Children: []*Tree{
							termSymbol(opar),
							{
								Value:    expression,
								Children: []*Tree{termIdt("a")},
							},
							termSymbol(cpar),
						},
					},
				},
			},
		},
	},
}

var emptyGrammar = &parseTest{
	name:    "∅",
	grammar: Empty,
	r:       &Tree{Value: "∅"},
}

var eofError = &parseTest{
	name:    "EOF found",
	grammar: basicGrammar(),
	tks:     []*Token{fixedString(opar)},
	r:       &Tree{Value: expression},
	err:     &UnexpectedEOFErr{After: opar},
}

var expectingExpr = &parseTest{
	name:    "Expecting expression",
	grammar: basicGrammar(),
	tks:     []*Token{fixedString("X")},
	r:       &Tree{Value: expression},
	err: &ExpectingErr{
		Expected: opar,
		Actual:   "X",
	},
}

var remainingToken = &parseTest{
	name:    "remaining token",
	grammar: basicGrammar(),
	tks:     []*Token{identifier("a"), fixedString("X")},
	r: &Tree{
		Value:    expression,
		Children: []*Tree{termIdt("a")},
	},
	err: &RemainingTokenErr{Token: fixedString("X")},
}

func TestT0(t *testing.T) {
	tests := []*parseTest{
		aTimesB,
		aPlusB,
		aPlusBTimesC,
		aMinusBPlusC,
		aPlusPar,
		aPar,
		emptyGrammar,
		eofError,
		expectingExpr,
		remainingToken,
	}

	for i, j := range tests {
		tkf := &tkStream{tokens: j.tks}
		r, e := Parse(j.grammar, tkf)
		require.Equal(t, j.err, e, "At %d, %s", i, j.name)
		require.Equal(t, j.r, r, "At %d, %s", i, j.name)
	}
}

func TestT1(t *testing.T) {
	g := &Symbol{
		Name: "letter",
		Header: &Symbol{
			Name:       "id",
			IsTerminal: true,
		},
	}
	tks := []*Token{
		{Name: "id", Value: "a"},
	}
	tkf := &tkStream{tokens: tks}
	r, e := Parse(g, tkf)
	require.NoError(t, e)
	rt := &Tree{
		Value:    "letter",
		Children: []*Tree{idt("a")},
	}
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

func basicGrammar() (exp *Symbol) {
	// expression = term {("+"|"-") term}.
	// term = factor {("*"|"/") factor}.
	// factor = id | "(" expression ")".
	exp = &Symbol{
		Name: expression,
	}

	fact := &Symbol{
		Name: factor,
		Header: &Symbol{
			Name:       id,
			IsTerminal: true,
		},
		Alt: &Symbol{
			Name:       opar,
			IsTerminal: true,
			Next: &Symbol{
				Header: exp,
				Next: &Symbol{
					Name:       cpar,
					IsTerminal: true,
				},
			},
		},
	}

	trm := &Symbol{
		Name:   term,
		Header: fact,
		Next: &Symbol{
			Name:       times,
			IsTerminal: true,
			Alt: &Symbol{
				Name:       div,
				IsTerminal: true,
				Alt:        Empty,
			},
		},
	}
	trm.Next.Next = trm
	trm.Next.Alt.Next = trm
	exp.Header = trm
	exp.Next = &Symbol{
		Name:       plus,
		IsTerminal: true,
		Next:       exp,
		Alt: &Symbol{
			Name:       minus,
			IsTerminal: true,
			Next:       exp,
			Alt:        Empty,
		},
	}
	return
}

func fixedString(s string) *Token {
	return &Token{Name: s, Value: s}
}

func identifier(s string) *Token {
	return &Token{Name: id, Value: s}
}

func idt(s string) *Tree {
	return &Tree{Value: id, Token: &Token{Name: id, Value: s}}
}

func termIdt(s string) *Tree {
	return &Tree{
		Value: term,
		Children: []*Tree{
			{
				Value:    factor,
				Children: []*Tree{idt(s)},
			},
		},
	}
}

func termSymbol(s string) *Tree {
	return &Tree{Value: s, Token: &Token{Name: s, Value: s}}
}
