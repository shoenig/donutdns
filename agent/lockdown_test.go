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
			name: "mix",
			cc: &CoreConfig{
				AllowFile:  "/opt/allows.txt",
				BlockDir:   "/opt/blocks",
				SuffixFile: "/opt/suffix.txt",
			},
			exp: []*landlock.Path{
				landlock.Dir("/opt/blocks", "r"),
				landlock.File("/opt/allows.txt", "r"),
				landlock.File("/opt/suffix.txt", "r"),
			},
		},
		{
			name: "all",
			cc: &CoreConfig{
				AllowFile:  "/opt/allows.txt",
				BlockFile: "/opt/block.txt",
				SuffixFile: "/opt/suffix.txt",
				AllowDir: "/opt/allow",
				BlockDir:   "/opt/blocks",
				SuffixDir: "/opt/suffix",
			},
			exp: []*landlock.Path{
				landlock.File("/opt/block.txt", "r"),
				landlock.File("/opt/allows.txt", "r"),
				landlock.File("/opt/suffix.txt", "r"),
				landlock.Dir("/opt/blocks", "r"),
				landlock.Dir("/opt/allow", "r"),
				landlock.Dir("/opt/suffix", "r"),
			},
		},
	}

	for _, tc := range cases {
		result := readable(tc.cc)
		must.SliceContainsAll(t, tc.exp, result)
	}
}

func Test_Lockdown(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := Lockdown(&CoreConfig{
			BlockFile: "../hack/social-media.list",
		})
		must.NoError(t, err)
	})

	t.Run("does not exist", func(t *testing.T) {
		err := Lockdown(&CoreConfig{
			BlockFile: "/does/not/exist",
		})
		must.ErrorContains(t, err, "no such file")
	})
}
