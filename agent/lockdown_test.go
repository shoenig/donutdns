package agent

import (
	"testing"

	"github.com/shoenig/go-landlock"
	"github.com/shoenig/test/must"
)

func Test_readable(t *testing.T) {
	cases := []struct {
		name string
		cc   *CoreConfig
		exp  []*landlock.Path
	}{
		{
			name: "none",
			cc:   new(CoreConfig),
			exp:  nil,
		},
		{
			name: "partial",
			cc:   &CoreConfig{BlockFile: "/opt/blocks.txt"},
			exp:  []*landlock.Path{landlock.File("/opt/blocks.txt", "r")},
		},
		{
			name: "all",
			cc: &CoreConfig{
				AllowFile:  "/opt/allows.txt",
				BlockFile:  "/opt/blocks.txt",
				SuffixFile: "/opt/suffix.txt",
			},
			exp: []*landlock.Path{
				landlock.File("/opt/blocks.txt", "r"),
				landlock.File("/opt/allows.txt", "r"),
				landlock.File("/opt/suffix.txt", "r"),
			},
		},
	}

	for _, tc := range cases {
		result := readable(tc.cc)
		must.SliceContainsAll(t, tc.exp, result)
	}
}
