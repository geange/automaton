package automaton

import "sync"

const (
	DEFAULT_EXPECTED_ELEMENTS = 4
	DEFAULT_LOAD_FACTOR       = 0.75
	MIN_LOAD_FACTOR           = 1 / 100.0
	MAX_LOAD_FACTOR           = 99 / 100.0
	MIN_HASH_ARRAY_LENGTH     = 4
	MAX_HASH_ARRAY_LENGTH     = uint32(0x80000000) >> 1
)

type IntIntHashmap struct {
	sync.RWMutex

	keys   []int32
	values []int32

	assigned      int
	mask          uint32  // Mask for slot scans in {@link #keys}.
	resizeAt      int     // Expand (rehash) {@link #keys} when {@link #assigned} hits this value.
	hasEmptyKey   bool    // Special treatment for the "empty slot" key marker.
	loadFactor    float64 // The load factor for {@link #keys}.
	iterationSeed int     // Seed used to ensure the hash iteration order is different from an iteration to another.
}

func (m *IntIntHashmap) Set(key, value int32) int32 {
	panic("")
}

func (m *IntIntHashmap) Exists(key int32) bool {
	m.RLock()
	defer m.RUnlock()

	_, exist := m.get(key)
	return exist
}

func (m *IntIntHashmap) Get(key int32) (int32, bool) {
	m.RLock()
	defer m.RUnlock()

	return m.get(key)
}

func (m *IntIntHashmap) get(key int32) (int32, bool) {
	if key == 0 {
		if m.hasEmptyKey {
			return m.values[m.mask+1], true
		}
		return 0, false
	}

	idx := uint32(m.hashKey(key)) & m.mask

	keys := m.keys
	for {
		if keys[idx] == key {
			return m.values[idx], true
		}
		if keys[idx] == 0 {
			return 0, false
		}
		idx = (idx + 1) & m.mask
	}
}

func (m *IntIntHashmap) IndexOf(key int32) int {
	panic("")
}

func (m *IntIntHashmap) IndexGet(index int) int {
	panic("")
}

func (m *IntIntHashmap) IndexRemove(index int) int32 {
	panic("")
}

func (m *IntIntHashmap) hashKey(key int32) int32 {
	return mixPhi(key)
}

func (m *IntIntHashmap) indexGet(index int) (int32, bool) {
	panic("")
}

func (m *IntIntHashmap) indexRemove(index int) int32 {
	panic("")
}

func (m *IntIntHashmap) indexExists(index int) int32 {
	panic("")
}

func (m *IntIntHashmap) indexSet(index int, value int32) int32 {
	panic("")
}

type IIMapMatchFunc func(oldValue int32, exist bool) bool

func (m *IntIntHashmap) Remove(key int32, match IIMapMatchFunc) int32 {
	panic("")
}

func (m *IntIntHashmap) Replace(key int32, value int32, match IIMapMatchFunc) int32 {
	if match(m.Get(key)) {
		return m.Set(key, value)
	}
	return value
}

func (m *IntIntHashmap) IndexReplace(index int, value int32) int32 {
	panic("")
}

func (m *IntIntHashmap) Size() int {
	panic("")
}

func (m *IntIntHashmap) Keys() []int32 {
	panic("")
}

func (m *IntIntHashmap) IsEmpty() bool {
	return m.Size() == 0
}
