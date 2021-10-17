package set

import "strings"

type nothing struct{}

var null = nothing{}

// New creates a new Set.
func New() *Set {
	return &Set{
		items: make(map[string]nothing),
	}
}

// Set contains a set of unique strings.
type Set struct {
	items map[string]nothing
}

// Add item to the Set.
func (s *Set) Add(item string) {
	s.items[item] = null
}

// Len returns the number of items in the Set.
func (s *Set) Len() int {
	return len(s.items)
}

// Union includes all elements of o into the Set.
func (s *Set) Union(o *Set) {
	for k := range o.items {
		s.items[k] = null
	}
}

// Has returns whether the Set contains item.
//
// For convenience, strip the suffix '.' from item before checking.
// Optimized for domains, which may or may not contain a trailing period.
func (s *Set) Has(item string) bool {
	clean := strings.TrimRight(item, ".")
	_, exists := s.items[clean]
	return exists
}
