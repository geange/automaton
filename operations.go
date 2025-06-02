package automaton

import (
	"bytes"
	"errors"
	"slices"
	"sync/atomic"
	"unicode"

	"github.com/bits-and-blooms/bitset"
)

// DeterminizeAutomaton Determinizes the given automaton.
// Worst case complexity: exponential in number of states.
// Params: 	workLimit – Maximum amount of "work" that the powerset construction will spend before throwing
//
//	TooComplexToDeterminizeException. Higher numbers allow this operation to consume more memory and
//	CPU but allow more complex automatons. Use DEFAULT_DETERMINIZE_WORK_LIMIT as a decent default
//	if you don't otherwise know what to specify.
//
// Throws: TooComplexToDeterminizeException – if determinizing requires more than workLimit "effort"
func DeterminizeAutomaton(a *Automaton, workLimit int) *Automaton {
	if a.IsDeterministic() {
		return a
	}
	if a.GetNumStates() <= 1 {
		// Already determinized
		return a
	}

	// subset construction
	b := NewBuilder()

	// Same initial values and state will always have the same hashCode
	initialSet := NewFrozenIntSet([]int{0}, mix32(0)+1, 0)
	// Create state 0:
	b.CreateState()

	worklist := make([]*FrozenIntSet, 0)
	newState := NewHashMap[int](WithCapacity(1))

	worklist = append(worklist, initialSet)
	b.SetAccept(0, a.IsAccept(0))
	newState.Set(initialSet, 0)

	// TODO:

	return a
}

// IsEmptyAutomaton
// Returns true if the given automaton accepts no strings.
func IsEmptyAutomaton(a *Automaton) bool {
	if a.GetNumStates() == 0 {
		// Common case: no states
		return true
	}

	if a.IsAccept(0) == false && a.GetNumTransitionsWithState(0) == 0 {
		// Common case: just one initial state
		return true
	}
	if a.IsAccept(0) == true {
		// Apparently common case: it accepts the damned empty string
		return false
	}

	workList := make([]int, 0)
	seen := bitset.New(uint(a.GetNumStates()))
	workList = append(workList, 0)
	seen.Set(0)

	t := NewTransition()
	for len(workList) > 0 {
		state := workList[0]
		workList = workList[1:]

		if a.IsAccept(state) {
			return false
		}

		count := a.InitTransition(state, t)
		for i := 0; i < count; i++ {
			a.GetNextTransition(t)
			if seen.Test(uint(t.Dest)) == false {
				workList = append(workList, t.Dest)
				seen.Set(uint(t.Dest))
			}
		}
	}
	return true
}

// IsTotalAutomaton
// Returns true if the given automaton accepts all strings. The automaton must be minimized.
func IsTotalAutomaton(a *Automaton) bool {
	return IsTotalAutomatonRange(a, 0, 0x10FFFF)
}

// IsTotalAutomatonRange
// Returns true if the given automaton accepts all strings for the specified min/max range of the alphabet.
// The automaton must be minimized.
func IsTotalAutomatonRange(a *Automaton, minAlphabet, maxAlphabet int) bool {
	if a.IsAccept(0) && a.GetNumTransitionsWithState(0) == 1 {
		t := NewTransition()
		a.getTransition(0, 0, t)
		return t.Dest == 0 && t.Min == minAlphabet && t.Max == maxAlphabet
	}
	return false
}

func GetSingletonAutomaton(a *Automaton) ([]int, error) {
	if a.IsDeterministic() == false {
		return nil, errors.New("input automaton must be deterministic")
	}

	ints := make([]int, 0)
	visited := make(map[int]struct{})
	s := 0
	t := NewTransition()
	for {
		visited[s] = struct{}{}

		if a.IsAccept(s) == false {
			if a.GetNumTransitionsWithState(s) == 1 {
				a.getTransition(s, 0, t)
				if _, ok := visited[t.Dest]; t.Min == t.Max && ok {
					ints = append(ints, t.Min)
					s = t.Dest
					continue
				}
			}
		} else if a.GetNumTransitionsWithState(s) == 0 {
			return ints, nil
		}

		// Automaton accepts more than one string:
		return nil, nil
	}
}

func IsFiniteAutomaton(a *Automaton) *atomic.Bool {
	flag := &atomic.Bool{}

	if a.GetNumStates() == 0 {
		flag.Store(true)
		return flag
	}

	b1 := bitset.New(uint(a.GetNumStates()))
	b2 := bitset.New(uint(a.GetNumStates()))

	return isFinite(NewTransition(), a, 0, b1, b2, 0)
}

// Checks whether there is a loop containing state. (This is sufficient since there are never transitions to dead states.)
// TODO: not great that this is recursive... in theory a
// large automata could exceed java's stack so the maximum level of recursion is bounded to 1000
func isFinite(scratch *Transition, a *Automaton, state int, path, visited *bitset.BitSet, level int) *atomic.Bool {
	flag := &atomic.Bool{}

	// if (level > MAX_RECURSION_LEVEL) {
	//      throw new IllegalArgumentException("input automaton is too large: " +  level);
	//    }
	path.Set(uint(state))
	numTransitions := a.InitTransition(state, scratch)
	for t := 0; t < numTransitions; t++ {
		a.getTransition(state, t, scratch)
		if path.Test(uint(scratch.Dest)) || (!visited.Test(uint(scratch.Dest)) && !isFinite(scratch, a, scratch.Dest, path, visited, level+1).Load()) {
			flag.Store(false)
			return flag
		}
	}
	path.Clear(uint(state))
	visited.Set(uint(state))
	flag.Store(true)
	return flag
}

// getCommonSuffixBytesRef
// Returns the longest BytesRef that is a suffix of all accepted strings. Worst case complexity: quadratic with the number of states+transitions.
// Returns: common suffix, which can be an empty (length 0) BytesRef (never null)
func getCommonSuffixBytesRef(a *Automaton) ([]byte, error) {
	// reverse the language of the automaton, then reverse its common prefix.
	ra, err := reverse(a)
	if err != nil {
		return nil, err
	}
	r, err := removeDeadStates(ra)
	if err != nil {
		return nil, err
	}

	ref, err := getCommonPrefixBytesRef(r)
	if err != nil {
		return nil, err
	}
	slices.Reverse(ref)
	return ref, nil
}

//func reverse[T cmp.Ordered](ref []T) {
//	i, j := 0, len(ref)-1
//	for i < j {
//		ref[i], ref[j] = ref[j], ref[i]
//	}
//}

func BitsetAndNot(a, b *bitset.BitSet) {

}

func hasDeadStatesFromInitial(a *Automaton) bool {
	reachableFromInitial := getLiveStatesFromInitial(a)
	reachableFromAccept := getLiveStatesToAccept(a)
	reachableFromInitial.Difference(reachableFromAccept)
	return reachableFromInitial.Count() == 0
}

func getCommonPrefix(a *Automaton) (string, error) {

	if hasDeadStatesFromInitial(a) {
		return "", errors.New("input automaton has dead states")
	}
	if isEmpty(a) {
		return "", nil
	}
	builder := new(bytes.Buffer)
	scratch := NewTransition()
	visited := bitset.New(uint(a.GetNumStates()))
	current := bitset.New(uint(a.GetNumStates()))
	next := bitset.New(uint(a.GetNumStates()))
	current.Set(0) // start with initial state
OUT:
	for {
		label := -1
		// do a pass, stepping all current paths forward once
		state, ok := current.NextSet(0)
		for ok {
			visited.Set(state)
			// if it is an accept state, we are done
			if a.IsAccept(int(state)) {
				break OUT
			}
			for transition := 0; transition < a.GetNumTransitionsWithState(int(state)); transition++ {
				a.getTransition(int(state), transition, scratch)
				if label == -1 {
					label = scratch.Min
				}
				// either a range of labels, or label that doesn't match all the other paths this round
				if scratch.Min != scratch.Max || scratch.Min != label {
					break OUT
				}
				// mark target state for next iteration
				next.Set(uint(scratch.Dest))
			}
			state++
			if state >= current.Len() {
				ok = false
			} else {
				current.Set(state)
			}
		}

		// add the label to the prefix
		builder.WriteRune(rune(label))
		// swap "current" with "next", clear "next"
		tmp := current
		current = next
		next = tmp
		next.ClearAll()
	}
	return builder.String(), nil
}

func isEmpty(a *Automaton) bool {
	if a.GetNumStates() == 0 {
		// Common case: no states
		return true
	}
	if a.IsAccept(0) == false && a.GetNumTransitionsWithState(0) == 0 {
		// Common case: just one initial state
		return true
	}
	if a.IsAccept(0) == true {
		// Apparently common case: it accepts the damned empty string
		return false
	}

	workList := make([]int, 0)

	seen := bitset.New(uint(a.GetNumStates()))
	workList = append(workList, 0)
	seen.Set(0)

	t := NewTransition()
	for len(workList) > 0 {
		state := workList[0]
		workList = workList[1:]
		if a.IsAccept(state) {
			return false
		}
		count := a.InitTransition(state, t)
		for i := 0; i < count; i++ {
			a.GetNextTransition(t)
			if seen.Test(uint(t.Dest)) == false {
				workList = append(workList, t.Dest)
				seen.Set(uint(t.Dest))
			}
		}
	}

	return true
}

func getCommonPrefixBytesRef(a *Automaton) ([]byte, error) {
	prefix, err := getCommonPrefix(a)
	if err != nil {
		return nil, err
	}
	builder := new(bytes.Buffer)

	for _, ch := range prefix {
		if ch > 255 {
			return nil, errors.New("automaton is not binary")
		}
		builder.WriteRune(ch)
	}

	return builder.Bytes(), nil
}

func reverse(a *Automaton) (*Automaton, error) {
	return reverseStates(a, nil)
}

func reverseStates(a *Automaton, initialStates map[int]struct{}) (*Automaton, error) {
	panic("")
}

func reverseAutomaton(a *Automaton) *Automaton {
	return reverseAutomatonIntSet(a, nil)
}

func removeDeadStates(a *Automaton) (*Automaton, error) {
	numStates := a.GetNumStates()
	liveSet := getLiveStates(a)

	mp := make([]int, numStates)

	result := NewAutomaton()
	for i := 0; i < numStates; i++ {
		if liveSet.Test(uint(i)) {
			mp[i] = result.CreateState()
			result.SetAccept(mp[i], a.IsAccept(i))
		}
	}

	t := NewTransition()

	for i := 0; i < numStates; i++ {
		if liveSet.Test(uint(i)) {
			numTransitions := a.InitTransition(i, t)
			// filter out transitions to dead states:
			for j := 0; j < numTransitions; j++ {
				a.GetNextTransition(t)
				if liveSet.Test(uint(t.Dest)) {
					err := result.AddTransition(mp[i], mp[t.Dest], t.Min, t.Max)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	result.FinishState()
	//assert hasDeadStates(result) == false;
	return result, nil
}

func getLiveStates(a *Automaton) *bitset.BitSet {
	live := getLiveStatesFromInitial(a)
	live.Union(getLiveStatesToAccept(a))
	return live
}

func getLiveStatesFromInitial(a *Automaton) *bitset.BitSet {
	numStates := a.GetNumStates()
	live := bitset.New(uint(numStates))
	if numStates == 0 {
		return live
	}
	workList := make([]int, 0)
	live.Set(0)
	workList = append(workList, 0)

	t := NewTransition()
	for len(workList) == 0 {
		s := workList[0]
		count := a.InitTransition(s, t)
		for i := 0; i < count; i++ {
			a.GetNextTransition(t)
			if live.Test(uint(t.Dest)) == false {
				live.Set(uint(t.Dest))
				workList = append(workList, t.Dest)
			}
		}
	}

	return live
}

func getLiveStatesToAccept(a *Automaton) *bitset.BitSet {
	panic("")
}

func reverseAutomatonIntSet(a *Automaton, initialStates map[int]struct{}) *Automaton {
	if IsEmptyAutomaton(a) {
		return NewAutomaton()
	}

	numStates := a.GetNumStates()

	// Build a new automaton with all edges reversed
	builder := NewBuilder()

	// Initial node; we'll add epsilon transitions in the end:
	builder.CreateState()

	for s := 0; s < numStates; s++ {
		builder.CreateState()
	}

	// Old initial state becomes new accept state:
	builder.SetAccept(1, true)

	t := NewTransition()
	for s := 0; s < numStates; s++ {
		numTransitions := a.GetNumTransitionsWithState(s)
		a.InitTransition(s, t)
		for i := 0; i < numTransitions; i++ {
			a.GetNextTransition(t)
			builder.AddTransition(t.Dest+1, s+1, t.Min, t.Max)
		}
	}

	result := builder.Finish()

	s := 0
	acceptStates := a.getAcceptStates()
	for {
		if _, ok := acceptStates.NextSet(uint(s)); !(ok && s < numStates) {
			break
		}

		result.AddEpsilon(0, s+1)
		if initialStates != nil {
			initialStates[s+1] = struct{}{}
		}
		s++
	}

	result.FinishState()

	return result
}

func union(automatons ...*Automaton) (*Automaton, error) {
	result := NewAutomaton()

	// Create initial state:
	result.CreateState()

	// Copy over all automata
	for _, a := range automatons {
		result.Copy(a)
	}

	// Add epsilon transition from new initial state
	stateOffset := 1
	for _, a := range automatons {
		if a.GetNumStates() == 0 {
			continue
		}
		result.AddEpsilon(0, stateOffset)
		stateOffset += a.GetNumStates()
	}

	result.FinishState()

	return removeDeadStates(result)
}

func concatenate(automatons ...*Automaton) (*Automaton, error) {
	result := NewAutomaton()

	// First pass: create all states
	for _, a := range automatons {
		if a.GetNumStates() == 0 {
			result.FinishState()
			return result, nil
		}
		numStates := a.GetNumStates()
		for s := 0; s < numStates; s++ {
			result.CreateState()
		}
	}

	// Second pass: add transitions, carefully linking accept
	// states of A to init state of next A:
	stateOffset := 0
	t := NewTransition()

	for i, a := range automatons {
		numStates := a.GetNumStates()

		var nextA *Automaton
		if i == len(automatons)-1 {
			nextA = nil
		} else {
			nextA = automatons[i+1]
		}

		for s := 0; s < numStates; s++ {
			numTransitions := a.InitTransition(s, t)
			for j := 0; j < numTransitions; j++ {
				a.GetNextTransition(t)
				err := result.AddTransition(stateOffset+s, stateOffset+t.Dest, t.Min, t.Max)
				if err != nil {
					return nil, err
				}
			}

			if a.IsAccept(s) {
				followA := nextA
				followOffset := stateOffset
				upto := i + 1
				for {
					if followA != nil {
						// Adds a "virtual" epsilon transition:
						numTransitions = followA.InitTransition(0, t)
						for j := 0; j < numTransitions; j++ {
							followA.GetNextTransition(t)
							err := result.AddTransition(stateOffset+s, followOffset+numStates+t.Dest, t.Min, t.Max)
							if err != nil {
								return nil, err
							}
						}
						if followA.IsAccept(0) {
							// Keep chaining if followA accepts empty string
							followOffset += followA.GetNumStates()
							if upto == len(automatons)-1 {
								followA = nil
							} else {
								followA = automatons[upto+1]
							}
							upto++
						} else {
							break
						}
					} else {
						result.SetAccept(stateOffset+s, true)
						break
					}
				}
			}
		}

		stateOffset += numStates
	}

	if result.GetNumStates() == 0 {
		result.CreateState()
	}

	result.FinishState()

	return result, nil
}

func totalize(a *Automaton) (*Automaton, error) {
	result := NewAutomaton()
	numStates := a.GetNumStates()
	for i := 0; i < numStates; i++ {
		result.CreateState()
		result.SetAccept(i, a.IsAccept(i))
	}

	deadState := result.CreateState()
	err := result.AddTransition(deadState, deadState, 0, unicode.MaxRune)
	if err != nil {
		return nil, err
	}

	t := NewTransition()
	for i := 0; i < numStates; i++ {
		maxi := 0
		count := a.InitTransition(i, t)
		for j := 0; j < count; j++ {
			a.GetNextTransition(t)
			err := result.AddTransition(i, t.Dest, t.Min, t.Max)
			if err != nil {
				return nil, err
			}
			if t.Min > maxi {
				err := result.AddTransition(i, deadState, maxi, t.Min-1)
				if err != nil {
					return nil, err
				}
			}
			if t.Max+1 > maxi {
				maxi = t.Max + 1
			}
		}

		if maxi <= unicode.MaxRune {
			err := result.AddTransition(i, deadState, maxi, unicode.MaxRune)
			if err != nil {
				return nil, err
			}
		}
	}

	result.FinishState()
	return result, nil
}
func complement(a *Automaton, determinizeWorkLimit int) (*Automaton, error) {
	a, err := determinize(a, determinizeWorkLimit)
	if err != nil {
		return nil, err
	}
	a, err = totalize(a)
	if err != nil {
		return nil, err
	}
	numStates := a.GetNumStates()
	for p := 0; p < numStates; p++ {
		a.SetAccept(p, !a.IsAccept(p))
	}
	return removeDeadStates(a)
}

func determinize(a *Automaton, workLimit int) (*Automaton, error) {
	panic("")
}

func repeat(a *Automaton) (*Automaton, error) {
	if a.GetNumStates() == 0 {
		// Repeating the empty automata will still only accept the empty automata.
		return a, nil
	}
	builder := NewBuilder()
	builder.CreateState()
	builder.SetAccept(0, true)
	builder.CopyStates(a)

	t := NewTransition()
	count := a.InitTransition(0, t)
	for i := 0; i < count; i++ {
		a.GetNextTransition(t)
		builder.AddTransition(0, t.Dest+1, t.Min, t.Max)
	}

	numStates := a.GetNumStates()
	for s := 0; s < numStates; s++ {
		if a.IsAccept(s) {
			count = a.InitTransition(0, t)
			for i := 0; i < count; i++ {
				a.GetNextTransition(t)
				builder.AddTransition(s+1, t.Dest+1, t.Min, t.Max)
			}
		}
	}

	return builder.Finish(), nil
}

func repeatCount(a *Automaton, count int) (*Automaton, error) {
	if count == 0 {
		return repeat(a)
	}
	as := make([]*Automaton, 0)
	for count > 0 {
		count--
		as = append(as, a)
	}

	ra, err := repeat(a)
	if err != nil {
		return nil, err
	}
	as = append(as, ra)

	return concatenate(as...)
}

func repeatRange(a *Automaton, min, max int) (*Automaton, error) {
	if min > max {
		return defaultAutomata.MakeEmpty(), nil
	}

	var b *Automaton
	var err error
	if min == 0 {
		b = defaultAutomata.MakeEmptyString()
	} else if min == 1 {
		b = NewAutomaton()
		b.Copy(a)
	} else {
		as := make([]*Automaton, 0)
		for i := 0; i < min; i++ {
			as = append(as, a)
		}
		b, err = concatenate(as...)
		if err != nil {
			return nil, err
		}
	}

	prevAcceptStates := toSet(b, 0)
	builder := NewBuilder()
	builder.Copy(b)
	for i := min; i < max; i++ {
		numStates := builder.GetNumStates()
		builder.Copy(a)
		for s := range prevAcceptStates {
			builder.AddEpsilon(s, numStates)
		}
		prevAcceptStates = toSet(a, numStates)
	}

	return builder.Finish(), nil
}

func toSet(a *Automaton, offset int) map[int]struct{} {
	numStates := uint(a.GetNumStates())
	isAccept := a.getAcceptStates()
	result := make(map[int]struct{})
	upto := uint(0)
	var ok bool
	for upto < numStates {
		upto, ok = isAccept.NextSet(upto)
		if !ok {
			break
		}
		result[offset+int(upto)] = struct{}{}
		upto++
	}

	return result
}

func intersection(a1, a2 *Automaton) (*Automaton, error) {
	panic("")
}

func optional(a *Automaton) (*Automaton, error) {
	panic("")
}
