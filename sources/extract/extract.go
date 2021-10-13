package extract

import (
	"bufio"
	"io"
	"regexp"

	"gophers.dev/cmds/donutdns/sources/set"
)

const (
	// Generic is a basic domain regexp pattern
	// from: https://stackoverflow.com/a/30007882/221569
	//
	// It does very well, but matches ipv4 addresses in some cases.
	Generic = `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`
)

type Extractor interface {
	Extract(io.Reader) (*set.Set, error)
}

type extractor struct {
	re *regexp.Regexp
}

func New(re string) Extractor {
	return &extractor{
		re: regexp.MustCompile(re),
	}
}

func (e *extractor) Extract(r io.Reader) (*set.Set, error) {
	scanner := bufio.NewScanner(r)
	s := set.New()
	for scanner.Scan() {
		line := scanner.Text()
		if domain := e.parse(line); domain != "" {
			s.Add(domain)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

// parse returns a domain, or empty string if no domain is found
func (e *extractor) parse(line string) string {
	switch {
	case line == "":
		return ""
	case line[0] == '#':
		return ""
	}

	return e.re.FindString(line)
}
