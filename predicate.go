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

var predGrammar = []Symbol{
	{},                               // 0
	Empty,                            // 1
	{Name: pred, Header: 3},          // 2
	{Name: "", Header: 4, Next: 5},   // 3
	{Name: term, Header: 6, Next: 7}, // 4
	{
		Name:       EquivalesOp,
		IsTerminal: true,
		Next:       3,
		Alt:        8,
	}, // 5
	{
		Name:   junction,
		Header: 9,
	}, // 6
	{
		Name:       ImpliesOp,
		IsTerminal: true,
		Next:       10,
		Alt:        11,
	}, // 7
	{
		Name:       NotEquivalesOp,
		IsTerminal: true,
		Next:       3,
		Alt:        1,
	}, // 7

}

const (
	pred     = "predicate"
	junction = "junction"
	id       = "id"
	negation = "negation"
	tail     = "tail"
)

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
	t, e = Parse(predGrammar, tk)
	return
}
