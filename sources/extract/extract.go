package extract

import (
	"bufio"
	"io"
	"strings"

	"gophers.dev/cmds/donutdns/sources/set"
	"gophers.dev/pkgs/loggy"
)

type Extractor interface {
	Extract(io.Reader) (*set.Set, error)
}

type extractor struct {
	log loggy.Logger
}

func New() Extractor {
	return &extractor{
		log: loggy.New("extract"),
	}
}

func (e *extractor) Extract(r io.Reader) (*set.Set, error) {
	scanner := bufio.NewScanner(r)
	single := set.New()
	for scanner.Scan() {
		line := scanner.Text()
		if domain := parse(line); domain != "" {
			single.Add(domain)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return single, nil
}

// parse returns a domain, or empty string if no domain is found
func parse(line string) string {
	suffix := func(s string) string {
		start := strings.LastIndexAny(s, " \t")
		if start == -1 {
			return s
		}
		return s[start+1:]
	}

	clean := strings.TrimSpace(line)
	switch {
	case clean == "":
		return ""
	case strings.HasPrefix(clean, "#"):
		return ""
	default:
		return suffix(clean)
	}
}
