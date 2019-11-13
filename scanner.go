package parsing

import (
	"bufio"
	"errors"
	"io"
	"unicode"
	"unicode/utf8"
)

// Token is the object recognized by the scanner
type Token struct {
	Name  string
	Value string
}

const (
	// 0x3 is the end of file character
	eof = 0x3
	id  = "id"

	// SupportedToken is the constant representing all supported
	// tokens when none of them is recognized
	SupportedToken = "supported token"
)

// ReaderTokens implements the TokenStream interface
type ReaderTokens struct {
	rd *bufio.Reader
	ss []Scanner

	curr *Token
	e    error
	rn   rune
	sc   func(rune) (*Token, bool, bool)
	n    int

	end, search, read, scan bool
}

// NewReaderTokens creates a new ReaderTokens
func NewReaderTokens(rd io.Reader,
	ss []Scanner) (r *ReaderTokens) {
	r = &ReaderTokens{
		rd: bufio.NewReader(rd),
		ss: ss,
	}
	r.end, r.read, r.search, r.scan = false, true, false, false
	return
}

// Next is the TokenStream.Next implementation
func (r *ReaderTokens) Next() {
	for !r.end {
		if r.read {
			r.rn, _, r.e = r.rd.ReadRune()
			if errors.Is(r.e, io.EOF) {
				r.rn = eof
			} else if r.e != nil {
				r.end = true
			}
			r.read, r.search = false, !r.scan
		} else if r.search {
			if r.n == len(r.ss) {
				r.e, r.end =
					&ExpectingErr{
						Actual:   string(r.rn),
						Expected: SupportedToken,
					},
					true
			} else {
				r.sc, r.n, r.search = r.ss[r.n](), r.n+1, false
			}
		} else if !r.search {
			r.curr, r.read, r.end = r.sc(r.rn)
			r.search = !r.read || r.end
			r.scan = !r.search
		}
	}
	r.n, r.end = 0, r.e != nil
}

// Current is the TokenStream.Current implementation
func (r *ReaderTokens) Current() (t *Token, e error) {
	if errors.Is(r.e, io.EOF) && r.curr != nil {
		t, r.curr = r.curr, nil
	} else {
		t, e = r.curr, r.e
	}
	return
}

// Scanner is
type Scanner func() func(rune) (*Token, bool, bool)

// IdentScan scans an identifier
func IdentScan() func(rune) (*Token, bool, bool) {
	var ident string
	return func(rn rune) (t *Token, cont, prod bool) {
		cont = unicode.IsLetter(rn) ||
			(ident != "" && unicode.IsDigit(rn))
		if cont {
			ident = ident + string(rn)
		} else if ident != "" {
			t, prod = &Token{Value: ident, Name: id}, true
		}
		return
	}
}

// StrScan scans a specific string
func StrScan(strScan string) (s Scanner) {
	s = func() func(rune) (*Token, bool, bool) {
		str := strScan
		return func(rn rune) (t *Token, cont, prod bool) {
			sr, size := utf8.DecodeRuneInString(str)
			cont = sr != utf8.RuneError && sr == rn
			if cont {
				str = str[size:]
			}
			prod = len(str) == 0
			if prod {
				t, cont = &Token{Value: strScan, Name: strScan}, true
			}
			return
		}
	}
	return
}

// SpaceScan scans spaces
func SpaceScan() func(rune) (*Token, bool, bool) {
	start := false
	return func(rn rune) (t *Token, cont, prod bool) {
		cont = unicode.IsSpace(rn)
		if cont {
			start = true
		}
		prod = start && !cont
		if prod {
			t, start = new(Token), false
		}
		return
	}
}

// EOFScan scans EOF character
func EOFScan() func(rune) (*Token, bool, bool) {
	return func(r rune) (t *Token, cont, prod bool) {
		if r == eof {
			eofStr := string(eof)
			t, prod = &Token{Value: eofStr, Name: eofStr}, true
		}
		return
	}
}
