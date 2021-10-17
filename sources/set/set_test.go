package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Add(t *testing.T) {
	s := New()
	require.Equal(t, 0, s.Len())

	s.Add("foo")
	require.Equal(t, 1, s.Len())

	s.Add("bar")
	require.Equal(t, 2, s.Len())

	s.Add("foo")
	require.Equal(t, 2, s.Len())

	s.Add("")
	require.Equal(t, 2, s.Len())
}

func Test_Union(t *testing.T) {
	s1 := New()
	s1.Add("a")
	s1.Add("b")
	s1.Add("c")

	s2 := New()
	s2.Add("c")
	s2.Add("d")

	s1.Union(s2)
	require.Equal(t, 4, s1.Len())
	require.Equal(t, 2, s2.Len())
}

func Test_Has(t *testing.T) {
	s := New()
	s.Add("a")
	s.Add("b")
	require.True(t, s.Has("a"))
	require.True(t, s.Has("b"))
	require.False(t, s.Has("c"))
}
