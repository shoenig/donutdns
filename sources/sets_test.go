package sources

import (
	"testing"

	"github.com/hashicorp/go-set"
	"github.com/shoenig/test/must"
)

func TestSets_BlockBySuffix(t *testing.T) {
	suffixes := set.From[string]([]string{"evil.com", "ads.good.com"})

	cases := []struct {
		domain string
		exp    bool
	}{
		{"evil.com", true},
		{"a.evil.com", true},
		{"evil.net", false},

		{"good.com", false},
		{"a.good.com", false},
		{"ads.good.com", true},
		{"b.ads.good.com", true},
		{"c.b.ads.good.com", true},

		{"good.com.", false},
		{"evil.com.", true},
		{".good.com", false},
		{".evil.com", true},

		{"evil.com.good.com", false},
		{"good.com.evil.com", true},
	}

	for _, tc := range cases {
		t.Run(tc.domain, func(t *testing.T) {
			s := &Sets{suffix: suffixes}
			result := s.BlockBySuffix(tc.domain)
			must.Eq(t, tc.exp, result)
		})
	}
}
