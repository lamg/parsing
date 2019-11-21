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
)

/*
Grammar in EBNF syntax

predicate = term {('≡'|'≢') term}.
term = junction
	['⇒' junction {'⇒' junction} | '⇐' junction {'⇐' junction}].
junction = factor
	['∨' factor {'∨' factor} | '∧' factor {'∧' factor}].
factor =	[unaryOp] (identifier | '(' predicate ')').
unaryOp = '¬'.
*/
const (
	NotOp          = "¬" // C-k NO
	AndOp          = "∧" // C-k AN
	OrOp           = "∨" // C-k OR
	EquivalesOp    = "≡" // C-k 3=
	NotEquivalesOp = "≢" // C-k ne (custom def. `:digraph ne 8802`)
	ImpliesOp      = "⇒" // C-k =>
	FollowsOp      = "⇐" // C-k <=
)

func grammar() (s *Symbol) {
	s = &Symbol{Name: pred}
	fct := &Symbol{
		Name: factor,
		Header: &Symbol{
			Name:       NotOp,
			IsTerminal: true,
			Alt:        Empty,
		},
		Next: &Symbol{
			Name:       id,
			IsTerminal: true,
			Alt: &Symbol{
				Name:       opar,
				IsTerminal: true,
				Next: &Symbol{
					Header: s,
					Next: &Symbol{
						Name:       cpar,
						IsTerminal: true,
					},
				},
			},
		},
	}
	junct := optTail(junction, OrOp, AndOp, fct)
	println("ok")
	trm := optTail(term, ImpliesOp, FollowsOp, junct)
	s.Header = trm
	s.Next = &Symbol{
		Header: &Symbol{
			Header: &Symbol{
				Name:       EquivalesOp,
				IsTerminal: true,
				Alt: &Symbol{
					Name:       NotEquivalesOp,
					IsTerminal: true,
				},
			},
			Next: trm,
		},
		Alt: Empty,
	}
	s.Next.Next = s.Next
	return
}

const (
	pred     = "predicate"
	junction = "junction"
	id       = "id"
	negation = "negation"
	tail     = "tail"
)

func tailSymbol(op string,
	next *Symbol) (t *Symbol) {
	t = &Symbol{
		Name: tail,
		Header: &Symbol{
			Name:       op,
			IsTerminal: true,
			Next:       next,
			Alt:        Empty,
		},
	}
	t.Next = t
	return
}

func optTail(name, tail0, tail1 string,
	header *Symbol) (r *Symbol) {
	t0, t1 := tailSymbol(tail0, header), tailSymbol(tail1, header)
	t0.Alt = t1
	r = &Symbol{
		Name:   name,
		Header: header,
		Next:   t0,
	}
	return
}

// ParsePred parses a predicate
func ParsePred(s string) (t *Tree, e error) {
	rd := strings.NewReader(s)
	ss := []Scanner{
		IdentScan,
		SpaceScan,
		StrScan(NotOp),
		StrScan(AndOp),
		StrScan(OrOp),
		StrScan(EquivalesOp),
		StrScan(NotEquivalesOp),
		StrScan(ImpliesOp),
		StrScan(FollowsOp),
		StrScan(opar),
		StrScan(cpar),
	}
	tk := NewReaderTokens(rd, ss)
	t, e = Parse(grammar(), tk)
	return
}
