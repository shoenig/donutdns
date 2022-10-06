package sources

import (
	"os"
	"strings"

	"github.com/hashicorp/go-set"
	"github.com/shoenig/donutdns/agent"
	"github.com/shoenig/donutdns/output"
	"github.com/shoenig/donutdns/sources/extract"
	"github.com/shoenig/ignore"
)

// Sets enables efficient look-ups of whether a domain should be allowable or blocked.
type Sets struct {
	allow  *set.Set[string]
	block  *set.Set[string]
	suffix *set.Set[string]
}

// New returns a Sets pre-filled according to cc.
func New(logger output.Logger, cc *agent.CoreConfig) *Sets {
	allow := set.New[string](100)
	block := set.New[string](100)
	suffix := set.New[string](100)

	// initialize defaults if enabled
	if !cc.NoDefaults {
		defaults(block, logger)
	}

	// insert individual custom allowable domains
	allow.InsertAll(cc.Allows)

	// insert file of custom allowable domains
	custom(cc.AllowFile, allow)

	// insert individual custom block domains
	block.InsertAll(cc.Blocks)

	// insert file of custom block domains
	custom(cc.BlockFile, block)

	// insert individual block domain suffixes
	suffix.InsertAll(cc.Suffix)

	// insert file of custom block domain suffixes
	custom(cc.SuffixFile, suffix)

	return &Sets{
		allow:  allow,
		block:  block,
		suffix: suffix,
	}
}

// Size returns the number of items in the allow, block, suffix sets.
func (s *Sets) Size() (int, int, int) {
	allow := s.allow.Size()
	block := s.block.Size()
	suffix := s.suffix.Size()
	return allow, block, suffix
}

// Allow indicates whether domain is on the explicit allow-list.
func (s *Sets) Allow(domain string) bool {
	return s.allow.Contains(domain)
}

// BlockByMatch indicates whether domain is on the explicit block-list.
func (s *Sets) BlockByMatch(domain string) bool {
	return s.block.Contains(domain)
}

// BlockBySuffix indicates whether domain is on the suffix block-list.
func (s *Sets) BlockBySuffix(domain string) bool {
	if s.suffix.Size() == 0 {
		return false
	}

	domain = strings.Trim(domain, ".")
	if domain == "" {
		return false
	}

	if s.suffix.Contains(domain) {
		return true
	}

	idx := strings.Index(domain, ".")
	if idx <= 0 {
		return false
	}

	return s.BlockBySuffix(domain[idx+1:])
}

func defaults(set *set.Set[string], logger output.Logger) {
	d := NewDownloader(logger)
	s, err := d.Download(Defaults())
	if err != nil {
		panic(err)
	}
	set.InsertSet(s)
}

func custom(filename string, set *set.Set[string]) {
	if filename == "" {
		return // nothing to do
	}

	// for now, everything uses the generic domain extractor
	ex := extract.New(extract.Generic)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer ignore.Close(f)

	s, err := ex.Extract(f)
	if err != nil {
		panic(err)
	}
	set.InsertSet(s)
}