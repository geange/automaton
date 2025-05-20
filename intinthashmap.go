package automaton

import (
	"errors"
	"iter"
	"sync"
)

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

func (m *IntIntHashmap) AddTo(key int32, incrementValue int32) int {
	return m.PutOrAdd(key, incrementValue, incrementValue)
}

func (m *IntIntHashmap) PutOrAdd(key, putValue, incrementValue int32) int {
	keyIndex, exist := m.IndexOf(key)
	if exist {
		putValue = m.values[keyIndex] + incrementValue
		m.IndexReplace(keyIndex, putValue)
	} else {
		_ = m.IndexInsert(keyIndex, key, incrementValue)
	}
	return keyIndex
}

func (m *IntIntHashmap) IndexOf(key int32) (int, bool) {
	mask := m.mask
	if key == 0 {
		if m.hasEmptyKey {
			return int(m.mask + 1), true
		}
		return int(m.mask + 1), false
	}

	slot := int(uint32(m.hashKey(key)) & mask)

	existing := m.keys[slot]

	for existing != 0 {
		if existing == key {
			return slot, true
		}
		slot = int(uint32(slot+1) & mask)
		existing = m.keys[slot]
	}
	return slot, false
}

func (m *IntIntHashmap) IndexGet(idx int) (int32, bool) {
	if idx < 0 {
		return 0, false
	}
	if !(idx <= int(m.mask) || (idx == int(m.mask+1) && m.hasEmptyKey)) {
		return 0, false
	}
	return m.values[idx], true
}

func (m *IntIntHashmap) IndexRemove(idx int) (int32, bool) {
	if idx < 0 {
		return 0, false
	}
	if !(idx <= int(m.mask) || (idx == int(m.mask+1) && m.hasEmptyKey)) {
		return 0, false
	}

	previousValue := m.values[idx]
	if idx > int(m.mask) {
		m.hasEmptyKey = false
		m.values[idx] = 0
	} else {
		m.shiftConflictingKeys(idx)
	}
	return previousValue, true
}

func (m *IntIntHashmap) shiftConflictingKeys(gapSlot int) {
	keys := m.keys
	values := m.values
	mask := m.mask

	// Perform shifts of conflicting keys to fill in the gap.
	distance := 0
	for {
		distance++
		slot := int(uint32(gapSlot+distance) & mask)
		existing := keys[slot]
		if (existing) == 0 {
			break
		}

		idealSlot := m.hashKey(existing)
		shift := int(uint32(slot-int(idealSlot)) & mask)
		if shift >= distance {
			// Entry at this position was originally at or before the gap slot.
			// Move the conflict-shifted entry to the gap's position and repeat the procedure
			// for any entries to the right of the current position, treating it
			// as the new gap.
			keys[gapSlot] = existing
			values[gapSlot] = values[slot]
			gapSlot = slot
			distance = 0
		}
	}

	// Mark the last found gap slot without a conflict as empty.
	keys[gapSlot] = 0
	values[gapSlot] = 0
	m.assigned--
}

func (m *IntIntHashmap) validateIndex(idx int) bool {
	if idx < 0 {
		return false
	}
	if idx <= int(m.mask) {
		return true
	}

	if idx == int(m.mask+1) && m.hasEmptyKey {
		return true
	}

	return false
}

func (m *IntIntHashmap) IndexReplace(idx int, value int32) bool {
	panic("implement me")
}

func (m *IntIntHashmap) IndexInsert(idx int, key, value int32) error {
	if idx < 0 {
		return errors.New("invalid index")
	}
	if key == 0 {
		if idx != int(m.mask+1) {
			return errors.New("invalid index")
		}
		m.values[idx] = value
		m.hasEmptyKey = true
		return nil
	}

	if m.keys[idx] != 0 {
		return errors.New("current index is already in use")
	}
	if m.assigned == m.resizeAt {

	} else {
		m.keys[idx] = key
		m.values[idx] = value
	}
	m.assigned++
	return nil
}

func (m *IntIntHashmap) allocateThenInsertThenRehash(slot int, pendingKey, pendingValue int32) {
	panic("implement me")
}

func (m *IntIntHashmap) allocateBuffers(arraySize int) {
	panic("implement me")
}

func nextBufferSize(arraySize, elements int, loadFactor float64) int {
	panic("implement me")
}

func (m *IntIntHashmap) rehash(fromKeys, fromValues []int32) {
	panic("implement me")
}

func (m *IntIntHashmap) Size() int {
	empty := 0
	if m.hasEmptyKey {
		empty = 1
	}
	return m.assigned + empty
}

func (m *IntIntHashmap) Keys() iter.Seq[int32] {
	panic("")
}

func (m *IntIntHashmap) hashKey(key int32) int32 {
	return mixPhi(key)
}
