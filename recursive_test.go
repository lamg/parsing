package parsing

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestT0(t *testing.T) {
	g := &Grammar{
		Value: "factor",
		Next: &Grammar{
			Value: "letter",
			Next: &Grammar{
				Value:      "a",
				IsTerminal: true,
				Next: &Grammar{
					Value:      "*",
					IsTerminal: true,
					Next: &Grammar{
						Value: "factor",
					},
				},
				Alt: &Grammar{Value: "b", Alt: &Grammar{IsEmpty: true}},
			},
		},
	}
	tokens, n := []*Token{
		{Value: "a"},
		{Value: "*"},
		{Value: "b"},
		{Value: "eof"},
	}, 0
	tks := func() (t *Token, e error) {
		if n == len(tokens) {
			e = new(EndOfFile)
		} else {
			t, n = tokens[n], n+1
		}
		return
	}
	r, e := Parse(g, tks)
	require.NoError(t, e)
	rt := &Tree{
		Value: "factor",
		Children: []*Tree{
			{Value: "letter", Children: []*Tree{{Value: "a"}}},
			{Value: "*"},
			{Value: "letter", Children: []*Tree{{Value: "b"}}},
		},
	}
	require.Equal(t, rt, r)
}
