package agent

import (
	"strings"

	"github.com/shoenig/go-landlock"
)

func Lockdown(cc *CoreConfig) error {
	paths := make([]*landlock.Path, 0, 4)
	paths = append(paths, sysPaths...)
	paths = append(paths, readable(cc)...)
	locker := landlock.New(paths...)
	return locker.Lock(landlock.OnlySupported)
}

func readable(cc *CoreConfig) []*landlock.Path {
	var paths []*landlock.Path
	add := func(path string, f func(string, string) *landlock.Path) {
		if nonempty(path) {
			paths = append(paths, f(path, "r"))
		}
	}
	add(cc.AllowFile, landlock.File)
	add(cc.BlockFile, landlock.File)
	add(cc.SuffixFile, landlock.File)
	add(cc.AllowDir, landlock.Dir)
	add(cc.BlockDir, landlock.Dir)
	add(cc.SuffixDir, landlock.Dir)
	return paths
}

func nonempty(s string) bool {
	return strings.TrimSpace(s) != ""
}
