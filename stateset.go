package automaton

import "slices"

var _ IntSet = &StateSet{}

type StateSet struct {
	inner       map[int]int
	hashUpdated bool
	hashCode    uint64
}

func NewStateSet() *StateSet {
	return &StateSet{
		inner: make(map[int]int),
	}
}

func (s *StateSet) Hash() uint64 {
	if s.hashUpdated {
		return s.hashCode
	}
	s.hashCode = uint64(len(s.inner))

	s.hashCode = uint64(len(s.inner))
	for k := range s.inner {
		s.hashCode += uint64(mix(k))
	}
	s.hashUpdated = true
	return s.hashCode
}

func (s *StateSet) Equals(other Hashable) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateSet) GetArray() []int {
	keys := make([]int, 0, len(s.inner))

	for k := range s.inner {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func (s *StateSet) Size() int {
	return len(s.inner)
}

func (s *StateSet) keyChanged() {
	s.hashUpdated = false
}

func (s *StateSet) Incr(state int) {
	s.inner[state]++
	if s.inner[state] == 1 {
		s.keyChanged()
	}
}

func (s *StateSet) Decr(state int) {
	count, ok := s.inner[state]
	if !ok {
		return
	}
	if count == 0 {
		delete(s.inner, state)
	} else {
		s.inner[state]--
	}
}

func (s *StateSet) Freeze(state int) *FrozenIntSet {
	return NewFrozenIntSet(s.GetArray(), state, s.hashCode)
}
