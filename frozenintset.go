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
	//if other == nil {
	//	if f == nil {
	//		return true
	//	}
	//	return false
	//}
	//
	//fs, ok := other.(*FrozenIntSet)
	//if !ok {
	//	return false
	//}
	//
	//if fs == nil {
	//	if f == nil {
	//		return true
	//	}
	//	return false
	//}
	//
	//if !slices.Equal(f.values, fs.values) {
	//	return false
	//}
	//
	//return f.state == fs.state && f.hashCode == fs.hashCode
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
