package set

import "strings"

type nothing struct{}

var null = nothing{}

func New() *Set {
	return &Set{
		data: make(map[string]nothing),
	}
}

type Set struct {
	data map[string]nothing
}

func (s *Set) Add(domain string) {
	s.data[domain] = null
}

func (s *Set) Len() int {
	return len(s.data)
}

func (s *Set) Union(o *Set) {
	for k := range o.data {
		s.data[k] = null
	}
}

func (s *Set) Has(name string) bool {
	clean := strings.TrimRight(name, ".")
	_, exists := s.data[clean]
	return exists
}
