package parsing

import (
	"errors"
	"fmt"
	"io"
)

// Token is the object recognized by the scanner
type Token struct {
	Name  string
	Value string
}

// Tree is the abstract syntax tree generated by the parser
type Tree struct {
	Value    string
	Children []*Tree
}

// Symbol is the data structure representing the grammar
// to be parsed
type Symbol struct {
	Name       string
	IsTerminal bool
	IsEmpty    bool
	Next       *Symbol
	Alt        *Symbol
	Header     *Symbol
}

// Parse parses the grammar in g using the tokens in tks
func Parse(g *Symbol, tks TokenStream) (t *Tree, e error) {
	curr := g
	t = &Tree{Value: curr.Name}
	for curr != nil && e == nil {
		if curr.IsTerminal {
			if !curr.IsEmpty {
				var tk *Token
				tk, e = tks.Current()
				if e == nil {
					if curr.Name == tk.Name {
						if g.IsTerminal {
							t.Value = tk.Value
						} else {
							t.Children = append(t.Children,
								&Tree{Value: tk.Value})
						}
						tks.Next()
					} else {
						e = &ExpectingErr{
							Expected: curr.Name,
							Actual:   tk.Name,
						}
					}
				} else if errors.Is(e, io.EOF) {
					e = &ExpectingErr{
						Expected: curr.Name,
						Actual:   e.Error(),
					}
				}
			}
		} else {
			var nt *Tree
			nt, e = Parse(curr.Header, tks)
			if e == nil {
				t.Children = append(t.Children, nt)
			}
		}
		if e == nil {
			curr = curr.Next
		} else {
			curr, e = curr.Alt, nil
		}
	}
	if errors.Is(e, io.EOF) {
		e = nil
	}
	return
}

// TokenStream is the interface for supplying tokens to the
// parser
type TokenStream interface {
	Current() (*Token, error)
	Next()
}

// ExpectingErr is the error signaled when an unexpected
// token is received
type ExpectingErr struct {
	Expected string
	Actual   string
}

func (x *ExpectingErr) Error() string {
	return fmt.Sprintf("Expected: %s, Actual: %s",
		x.Expected, x.Actual)
}

// Empty is the empty symbol used for closing loops in
// the Symbol data structure
var Empty = &Symbol{
	Name:       "∅",
	IsEmpty:    true,
	IsTerminal: true,
}
