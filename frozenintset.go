package automaton

var _ IntSet = &FrozenIntSet{}

type FrozenIntSet struct {
	values   []int
	state    int
	hashCode uint64
}

func (f *FrozenIntSet) Hash() uint64 {
	return f.hashCode
}

func (f *FrozenIntSet) Equals(other Hashable) bool {
	if f == nil {
		switch other.(type) {
		case *FrozenIntSet:
			ptr := other.(*FrozenIntSet)
			if ptr == nil {
				return true
			}
		case *StateSet:
			ptr := other.(*StateSet)
			if ptr == nil {
				return true
			}
		default:
			return false
		}
	}

	iset, ok := other.(IntSet)
	if !ok {
		return false
	}
	return iset.Hash() == f.Hash()
}

func NewFrozenIntSet(values []int, state int, hashCode uint64) *FrozenIntSet {
	return &FrozenIntSet{values: values, state: state, hashCode: hashCode}
}

func (f *FrozenIntSet) GetArray() []int {
	return f.values
}

func (f *FrozenIntSet) Size() int {
	return len(f.values)
}
