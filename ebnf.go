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

const (
	production = "production"
	syntax     = "syntax"
	EqualOp    = "="
	dot        = "."
	bar        = "|"
)

var ebnfGrammar = []Symbol{
	{},
	Empty,
	{Name: syntax, Header: 3, Next: 2, Alt: 1},    // 2
	{Name: production, Header: 4, Next: 5},        // 3
	{Name: Identifier, IsTerminal: true, Next: 4}, // 4
	{Name: EqualOp, IsTerminal: true, Next: 5},    // 5
	{Name: expression, Header: 6, Next: 7},        //
	{Name: term, Header: 8, Next: 9},              //
	{},

	{Name: "", Header: 8, Next: 9}, // 6
	{Name: dot, IsTerminal: true},  // 7
	{Name: bar, IsTerminal: true, Next: 7},
}

func Translate(rd io.Reader) (r []Symbol, e error) {
	return
}
