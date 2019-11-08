package parsing

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type parseTest struct {
	name string
	tks  []*Token
	r    *Tree
	err  error
}

var aTimesB = &parseTest{
	name: "a*b",
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
					{Value: times},
					{Value: factor, Children: []*Tree{idt("b")}},
				},
			},
		},
	},
}

var aPlusB = &parseTest{
	name: "a+b",
	tks: []*Token{
		identifier("a"),
		fixedString(plus),
		identifier("b"),
	},
	r: &Tree{
		Value: expression,
		Children: []*Tree{
			termIdt("a"),
			{Value: plus},
			termIdt("b"),
		},
	},
}

var aPlusBTimesC = &parseTest{
	name: "a+b*c",
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
			{Value: plus},
			{
				Value: term,
				Children: []*Tree{
					{Value: factor, Children: []*Tree{idt("b")}},
					{Value: times},
					{Value: factor, Children: []*Tree{idt("c")}},
				},
			},
		},
	},
}

var aMinusBPlusC = &parseTest{
	name: "a-b+c",
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
			{Value: minus},
			termIdt("b"),
			{Value: plus},
			termIdt("c"),
		},
	},
}

var aPlusPar = &parseTest{

	name: "a+(b+c-d)",
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
			{Value: plus},
			{
				Value: term,
				Children: []*Tree{
					{
						Value: factor,
						Children: []*Tree{
							{Value: opar},
							{
								Value: expression,
								Children: []*Tree{
									termIdt("b"),
									{Value: plus},
									termIdt("c"),
									{Value: minus},
									termIdt("d"),
								},
							},
							{Value: cpar},
						},
					},
				},
			},
		},
	},
}

func TestT0(t *testing.T) {
	exp := basicGrammar()
	tests := []*parseTest{
		aTimesB,
		aPlusB,
		aPlusBTimesC,
		aMinusBPlusC,
		aPlusPar,
	}

	for i, j := range tests {
		tkf := &tkStream{tokens: j.tks}
		r, e := Parse(exp, tkf)
		require.Equal(t, j.err, e, "At %d, %s", i, j.name)
		require.Equal(t, j.r, r, "At %d, %s", i, j.name)
	}
}

func TestT4(t *testing.T) {
	tks := []*Token{
		identifier("a"),
		fixedString(plus),
		fixedString(opar),
		identifier("b"),
		fixedString(plus),
		identifier("c"),
		fixedString(minus),
		identifier("d"),
		fixedString(cpar),
	}

	r := &Tree{
		Value: expression,
		Children: []*Tree{
			termIdt("a"),
			{Value: plus},
			{
				Value: term,
				Children: []*Tree{
					{
						Value: factor,
						Children: []*Tree{
							{Value: opar},
							{
								Value: expression,
								Children: []*Tree{
									termIdt("b"),
									{Value: plus},
									termIdt("c"),
									{Value: minus},
									termIdt("d"),
								},
							},
							{Value: cpar},
						},
					},
				},
			},
		},
	}
	tkf := &tkStream{tokens: tks}
	rt, e := Parse(basicGrammar(), tkf)
	require.NoError(t, e)

	require.Equal(t, r, rt)
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

func TestT2(t *testing.T) {
	exp := &Symbol{
		Name: expression,
	}
	fact := &Symbol{
		Name: factor,
		Header: &Symbol{
			Name:       id,
			IsTerminal: true,
		},
	}
	ftail := &Symbol{
		Name:       times,
		IsTerminal: true,
		Alt: &Symbol{
			Name:       div,
			IsTerminal: true,
			Alt:        Empty,
		},
	}
	exp.Header = fact
	fact.Next = ftail
	ftail.Next = fact
	ftail.Alt.Next = fact

	tkf := &tkStream{
		tokens: []*Token{{Value: "a", Name: "id"}},
	}
	rt, e := Parse(exp, tkf)
	require.NoError(t, e)
	r := &Tree{
		Value: expression,
		Children: []*Tree{
			{
				Value:    factor,
				Children: []*Tree{idt("a")},
			},
		},
	}
	require.Equal(t, r, rt)
}

func TestT3(t *testing.T) {
	tks := &tkStream{}
	g := Empty
	rt, e := Parse(g, tks)
	require.NoError(t, e)
	nt := &Tree{Value: "∅"}
	require.Equal(t, nt, rt)
}

func TestT5(t *testing.T) {
	tks := []*Token{
		fixedString(opar),
		identifier("a"),
		fixedString(cpar),
	}
	tkf := &tkStream{tokens: tks}
	r := &Tree{
		Value: expression,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{
						Value: factor,
						Children: []*Tree{
							{Value: opar},
							{
								Value: expression,
								Children: []*Tree{
									{
										Value: term,
										Children: []*Tree{
											{
												Value:    factor,
												Children: []*Tree{idt("a")},
											},
										},
									},
								},
							},
							{Value: cpar},
						},
					},
				},
			},
		},
	}
	rt, e := Parse(basicGrammar(), tkf)
	require.NoError(t, e)
	require.Equal(t, r, rt)
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
	id         = "id"
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
	return &Tree{Value: id, Children: []*Tree{{Value: s}}}
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
