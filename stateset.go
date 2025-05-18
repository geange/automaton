package automaton

var _ IntSet = &StateSet{}

type StateSet struct {
	inner        map[int]int
	hashUpdated  bool
	arrayUpdated bool
	arrayCache   []int
	hashCode     uint64
}

func (s *StateSet) Hash() uint64 {
	if s.hashUpdated {
		return s.hashCode
	}
	s.hashCode = uint64(len(s.inner))

	// TODO:
	return s.hashCode
}

func (s *StateSet) Equals(other Hashable) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateSet) GetArray() []int {
	//TODO implement me
	panic("implement me")
}

func (s *StateSet) Size() int {
	//TODO implement me
	panic("implement me")
}
