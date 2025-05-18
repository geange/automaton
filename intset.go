package automaton

type IntSet interface {
	Hashable

	GetArray() []int

	Size() int
}

var _ IntSet = &MockIntSet{}

type MockIntSet struct {
}

func (m MockIntSet) Hash() uint64 {
	return 0
}

func (m MockIntSet) Equals(other Hashable) bool {
	return false
}

func (m MockIntSet) GetArray() []int {
	return make([]int, 0)
}

func (m MockIntSet) Size() int {
	return 0
}
