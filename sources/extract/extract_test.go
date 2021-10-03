package extract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtract_parse(t *testing.T) {
	try := func(in, exp string) {
		result := parse(in)
		require.Equal(t, exp, result)
	}

	try("", "")
	try("foo", "foo")
	try("#foo", "")
	try(" foo ", "foo")
	try("foo bar", "bar")
	try(" foo bar baz\t", "baz")
}
