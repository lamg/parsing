package parsing

import (
	"fmt"
)

type Token struct {
	Value string
}

type Tree struct {
	Value    string
	Children []*Tree
}

type Grammar struct {
	Value      string
	IsTerminal bool
	IsEmpty    bool
	Next       *Grammar
	Alt        *Grammar
}

func Parse(g *Grammar, tks func() (*Token, error)) (t *Tree,
	e error) {
	tk, e := tks()
	curr, t := g, &Tree{Value: g.Value}
	for curr != nil && e == nil && tk.Value != "eof" {
		var nt *Tree
		if curr.IsTerminal {
			if g.Value == tk.Value {
				nt = &Tree{Value: g.Value}
				tk, e = tks()
			} else if !g.IsEmpty {
				e = &ExpectingErr{Expected: g.Value, Actual: tk.Value}
			}
		} else {
			nt, e = Parse(curr, tks)
		}
		if e == nil {
			curr = curr.Next
			if nt != nil {
				t.Children = append(t.Children, nt)
			}
		} else {
			curr = curr.Alt
		}
	}
	return
}

type ExpectingErr struct {
	Expected string
	Actual   string
}

func (x *ExpectingErr) Error() string {
	return fmt.Sprintf("Expected: %s, Actual: %s",
		x.Expected, x.Actual)
}

type EndOfFile struct{}

func (f *EndOfFile) Error() string {
	return "EOF"
}
