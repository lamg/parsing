package parsing

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestT0(t *testing.T) {
	exp := &Symbol{
		Name: "expression",
	}

	factor := &Symbol{
		Name:   "factor",
		Header: &Symbol{Name: "id", IsTerminal: true},
		Alt: &Symbol{
			Header: &Symbol{
				Name:       "(",
				IsTerminal: true,
			},
			Next: &Symbol{
				Header: exp,
				Next: &Symbol{
					Name:       ")",
					IsTerminal: true,
				},
			},
		},
	}

	term := &Symbol{
		Name:   "term",
		Header: factor,
		Next: &Symbol{
			Name:       "*",
			IsTerminal: true,
			Alt: &Symbol{
				Name:       "/",
				IsTerminal: true,
				Alt:        Empty,
			},
		},
	}
	term.Next.Next = term
	term.Next.Alt.Next = term
	exp.Header = term
	exp.Next = &Symbol{
		Name:       "+",
		IsTerminal: true,
		Next:       term,
		Alt: &Symbol{
			Name:       "-",
			IsTerminal: true,
			Next:       term,
			Alt:        Empty,
		},
	}

	tks := []*Token{
		{Name: "id", Value: "a"},
		{Name: "*", Value: "*"},
		{Name: "id", Value: "b"},
	}
	tkf := &tkStream{tokens: tks}
	r, e := Parse(exp, tkf)
	require.NoError(t, e)
	tms := []*Tree{
		{Value: "factor", Children: []*Tree{{Value: "a"}}},
		{Value: "*"},
		{Value: "factor", Children: []*Tree{{Value: "b"}}},
	}
	rt := &Tree{
		Value: "expression",
		Children: []*Tree{
			{
				Value:    "term",
				Children: tms,
			},
		},
	}
	require.Equal(t, rt, r)

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
