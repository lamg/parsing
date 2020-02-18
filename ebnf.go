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
	"io"
)

const (
	production = "production"
	syntax     = "syntax"
	EqualOp    = "="
	dot        = "."
	bar        = "|"
)

/*
syntax = {production}.
production = identifier "=" expression ".".
expression = term {"|" term}.
term = factor {factor}.
factor = identifier | string | "(" expression ")" | "[" expression "]" | "{" expression "}".

identifier = letter {letter | digit}.
string = """ {character} """.

TODO algoritmo para generar la tabla a partir de la gramática
- si es un factor entonces referencia a Next
- si es un término entonces referencia a Alt
- si es un no terminal entonces referencia a Header
- repetición: header: elemento a repetir, next: autoreferencia, alt: vacío
*/

var syntaxSym = []Symbol{
	{Name: "syntax", Header: 3, Next: 2, Alt: 1},
}

var productionSym = []Symbol{
	{Name: "production", Header: 3},
	{Name: "identifier", Next: 4},
	{Name: "=", Next: 5},
	{Name: "expression", Header: 7, Next: 6},
	{Name: "."},
}

var expressionSym = []Symbol{
	{Name: "term", Header: 5, Next: 3},
	{Name: "", Header: 4, Next: 3, Alt: 0},
	{Name: "|", Next: 2},
}

var termSym = []Symbol{
	{Name: "factor", Header: 3, Next: 2},
	{Name: "", Header: 1, Next: 2, Alt: 0},
}

var factorSym = []Symbol{
	{Name: "", Header: 3},
	{Name: "identifier", Alt: 4},
	{Name: "string", Alt: 5},
	{Name: "parenthesesExpression", Header: 8, Alt: 6},
	{Name: "bracketExpression", Header: 9, Alt: 7},
	{Name: "bracesExpression", Header: 10},
}

var parExp = groupExpr("parenthesesExpression", "(", ")")
var bracketExp = groupExpr("bracketExpression", "[", "]")
var bracesExp = groupExpr("bracesExpression", "{", "}")

func groupExpr(name, open, close string) []Symbol {
	return []Symbol{
		{Name: name, Header: 3},
		{Name: open, Next: 4},
		{Name: "expression", Header: 6, Next: 5},
		{Name: close},
	}
}

func assemble() (ebnf []Symbol) {
	subexp := [][]Symbol{
		syntaxSym,
		productionSym,
		expressionSym,
		termSym,
		factorSym,
		parExp,
		bracketExp,
		bracesExp,
	}
	ebnf = []Symbol{{}, Empty}
	currPos := 0
	for _, j := range subexp {
		var startPos int
		startPos, currPos = currPos, currPos+len(j)
		for k := range j {
			noffs := newOffsets(startPos, j[k].Header, j[k].Next, j[k].Alt)
			j[k].Header, j[k].Next, j[k].Alt = noffs[0], noffs[1], noffs[2]
		}
	}
	for _, j := range subexp {
		ebnf = append(ebnf, j...)
	}
	return
}

func newOffsets(startPos int, positions ...int) (r []int) {
	r = make([]int, len(positions))
	for i, j := range positions {
		if j != 0 && j != 1 {
			r[i] = startPos + j
		}
	}
	return
}

func Translate(rd io.Reader) (r []Symbol, e error) {
	return
}
