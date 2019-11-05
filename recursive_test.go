package parsing

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestT0(t *testing.T) {
	exp := basicGrammar()
	tests := []struct {
		tks []*Token
		r   *Tree
		err error
	}{
		{
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
							{Value: factor, Children: []*Tree{{Value: "a"}}},
							{Value: times},
							{Value: factor, Children: []*Tree{{Value: "b"}}},
						},
					},
				},
			},
		},
		{
			tks: []*Token{
				identifier("a"),
				fixedString(plus),
				identifier("b"),
			},
			r: &Tree{
				Value: expression,
				Children: []*Tree{
					{
						Value: term,
						Children: []*Tree{
							{Value: factor, Children: []*Tree{{Value: "a"}}},
						},
					},
					{Value: plus},
					{
						Value: term,
						Children: []*Tree{
							{Value: factor, Children: []*Tree{{Value: "b"}}},
						},
					},
				},
			},
		},
		{
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
					{
						Value: term,
						Children: []*Tree{
							{Value: factor, Children: []*Tree{{Value: "a"}}},
						},
					},
					{Value: plus},
					{
						Value: term,
						Children: []*Tree{
							{Value: factor, Children: []*Tree{{Value: "b"}}},
							{Value: times},
							{Value: factor, Children: []*Tree{{Value: "c"}}},
						},
					},
				},
			},
		},
		{
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
					{
						Value: term,
						Children: []*Tree{
							{
								Value:    factor,
								Children: []*Tree{{Value: "a"}},
							},
							{Value: plus},
							{
								Value: factor,
								Children: []*Tree{
									{
										Value: expression,
										Children: []*Tree{
											{
												Value: term,
												Children: []*Tree{
													{
														Value:    factor,
														Children: []*Tree{{Value: "b"}},
													},
													{Value: plus},
													{
														Value:    factor,
														Children: []*Tree{{Value: "c"}},
													},
													{Value: minus},
													{
														Value:    factor,
														Children: []*Tree{{Value: "d"}},
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
	for i, j := range tests {
		tkf := &tkStream{tokens: j.tks}
		r, e := Parse(exp, tkf)
		require.Equal(t, j.err, e, "At %d", i)
		require.Equal(t, j.r, r, "At %d", i)
	}
}

func TestT4(t *testing.T) {
	tks := []*Token{
		identifier("a"),
		fixedString(plus),
		identifier("b"),
		fixedString(times),
		identifier("c"),
	}
	r := &Tree{
		Value: expression,
		Children: []*Tree{
			{
				Value: term,
				Children: []*Tree{
					{Value: factor, Children: []*Tree{{Value: "a"}}},
				},
			},
			{Value: plus},
			{
				Value: term,
				Children: []*Tree{
					{Value: factor, Children: []*Tree{{Value: "b"}}},
					{Value: times},
					{Value: factor, Children: []*Tree{{Value: "c"}}},
				},
			},
		},
	}
	tkf := &tkStream{tokens: tks}
	rt, e := Parse(basicGrammar(), tkf)
	require.NoError(t, e)
	require.Equal(t, term, rt.Children[0].Value)
	require.Equal(t, term, rt.Children[2].Value)
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
		Value: "letter", Children: []*Tree{{Value: "a"}},
	}
	require.Equal(t, rt, r)
}

func TestT2(t *testing.T) {
	exp := &Symbol{
		Name: "expression",
	}
	factor := &Symbol{
		Name: "factor",
		Header: &Symbol{
			Name:       "id",
			IsTerminal: true,
		},
	}
	ftail := &Symbol{
		Name:       "*",
		IsTerminal: true,
		Alt: &Symbol{
			Name:       "/",
			IsTerminal: true,
			Alt:        Empty,
		},
	}
	exp.Header = factor
	factor.Next = ftail
	ftail.Next = factor
	ftail.Alt.Next = factor

	tkf := &tkStream{
		tokens: []*Token{{Value: "a", Name: "id"}},
	}
	rt, e := Parse(exp, tkf)
	require.NoError(t, e)
	r := &Tree{
		Value: "expression",
		Children: []*Tree{
			{
				Value: "factor",
				Children: []*Tree{
					{Value: "a"},
				},
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
	exp = &Symbol{
		Name: expression,
	}

	fact := &Symbol{
		Name:   factor,
		Header: &Symbol{Name: id, IsTerminal: true},
		Alt: &Symbol{
			Header: &Symbol{
				Name:       opar,
				IsTerminal: true,
			},
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
