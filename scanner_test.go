package parsing

import (
	"strings"
	"testing"

	alg "github.com/lamg/algorithms"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	txt := "true++-*/bla9   x3  (Abla)true"
	tks := []string{"true",
		plus, plus, minus, times, div,
		"bla9", "", "x3", "",
		opar, "Abla", cpar, "true"}
	ss := []Scanner{
		StrScan(plus),
		StrScan(minus),
		StrScan(times),
		StrScan(div),
		StrScan(opar),
		StrScan(cpar),
		IdentScan,
		SpaceScan,
	}
	rt := NewReaderTokens(strings.NewReader(txt), ss)
	inf := func(i int) {
		rt.Next()
		tk, e := rt.Current()
		require.NoError(t, e, "At %s", tks[i])
		require.Equal(t, tks[i], tk.Value)
	}
	alg.Forall(inf, len(tks))
}

func TestStrScan(t *testing.T) {
	ns := StrScan(minus)()
	tk, cont, prod := ns('-')
	require.True(t, prod)
	require.True(t, cont)
	require.Equal(t, "-", tk.Value)
}

func TestIdentScan(t *testing.T) {
	ids := IdentScan()
	rs := []rune{'a', 'b', 'c', '0'}
	var tk *Token
	var cont, prod bool
	inf := func(i int) {
		_, cont, prod = ids(rs[i])
		require.True(t, cont)
		require.False(t, prod)
	}
	alg.Forall(inf, len(rs))
	tk, cont, prod = ids(' ')
	require.False(t, cont)
	require.True(t, prod)
	require.Equal(t, "abc0", tk.Value)

	ids0 := IdentScan()
	_, cont, prod = ids0(' ')
	require.False(t, cont)
	require.False(t, prod)
}
