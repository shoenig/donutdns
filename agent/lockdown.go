package agent

import "github.com/shoenig/go-landlock"

func Lockdown(cc *CoreConfig) error {
	paths := make([]*landlock.Path, 0, 4)
	paths = append(paths, sysPaths...)
	paths = append(paths, readable(cc)...)
	locker := landlock.New(paths...)
	return locker.Lock(landlock.OnlySupported)
}

func readable(cc *CoreConfig) []*landlock.Path {
	var paths []*landlock.Path
	if cc.AllowFile != "" {
		paths = append(paths, landlock.File(cc.AllowFile, "r"))
	}
	if cc.BlockFile != "" {
		paths = append(paths, landlock.File(cc.BlockFile, "r"))
	}
	if cc.SuffixFile != "" {
		paths = append(paths, landlock.File(cc.SuffixFile, "r"))
	}
	return paths
}
