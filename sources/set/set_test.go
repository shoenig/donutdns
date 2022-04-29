package set

import (
	"testing"

	"github.com/shoenig/test/must"
)

func Test_Add(t *testing.T) {
	s := New()
	must.EqOp(t, 0, s.Len())

	s.Add("foo")
	must.EqOp(t, 1, s.Len())

	s.Add("bar")
	must.EqOp(t, 2, s.Len())

	s.Add("foo")
	must.EqOp(t, 2, s.Len())

	s.Add("")
	must.EqOp(t, 2, s.Len())
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
	must.EqOp(t, 4, s1.Len())
	must.EqOp(t, 2, s2.Len())
}

func Test_Has(t *testing.T) {
	s := New()
	s.Add("a")
	s.Add("b")
	must.True(t, s.Has("a"))
	must.True(t, s.Has("b"))
	must.False(t, s.Has("c"))
}
