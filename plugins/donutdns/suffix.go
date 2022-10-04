package donutdns

import (
	"strings"

	"github.com/hashicorp/go-set"
)

func blockBySuffix(suffixes *set.Set[string], domain string) bool {
	if suffixes.Size() == 0 {
		return false
	}

	domain = strings.Trim(domain, ".")
	if domain == "" {
		return false
	}

	if suffixes.Contains(domain) {
		return true
	}

	idx := strings.Index(domain, ".")
	if idx <= 0 {
		return false
	}

	return blockBySuffix(suffixes, domain[idx+1:])
}
