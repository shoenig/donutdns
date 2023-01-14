//go:build linux

package agent

import (
	"github.com/shoenig/go-landlock"
)

var sysPaths = []*landlock.Path{
	landlock.Certs(),
}
