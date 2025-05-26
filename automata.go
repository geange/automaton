package automaton

import (
	"bytes"
	"errors"
	"math"
	"unicode"
)

type Automata struct {
}

// MakeEmpty
// Returns a new (deterministic) automaton with the empty language.
func (*Automata) MakeEmpty() *Automaton {
	a := NewAutomaton()
	a.FinishState()
	return a
}

// MakeEmptyString
// Returns a new (deterministic) automaton that accepts only the empty string.
func (*Automata) MakeEmptyString() *Automaton {
	a := NewAutomaton()
	a.CreateState()
	a.SetAccept(0, true)
	return a
}

// MakeAnyString
// Returns a new (deterministic) automaton that accepts all strings.
func (*Automata) MakeAnyString() (*Automaton, error) {
	a := NewAutomaton()
	s := a.CreateState()
	a.SetAccept(s, true)
	if err := a.AddTransition(s, s, 0, unicode.MaxRune); err != nil {
		return nil, err
	}
	a.FinishState()
	return a, nil
}

func (*Automata) MakeAnyBinary() (*Automaton, error) {
	a := NewAutomaton()
	s := a.CreateState()
	a.SetAccept(s, true)
	if err := a.AddTransition(s, s, 0, math.MaxUint8); err != nil {
		return nil, err
	}
	a.FinishState()
	return a, nil
}

func (*Automata) MakeNonEmptyBinary() (*Automaton, error) {
	a := NewAutomaton()
	s1 := a.CreateState()
	s2 := a.CreateState()
	a.SetAccept(s2, true)
	if err := a.AddTransition(s1, s2, 0, 255); err != nil {
		return nil, err
	}
	if err := a.AddTransition(s2, s2, 0, 255); err != nil {
		return nil, err
	}
	a.FinishState()
	return a, nil
}

func (r *Automata) MakeAnyChar() (*Automaton, error) {
	return r.MakeCharRange(0, unicode.MaxRune)
}

func (r *Automata) MakeChar(c int32) (*Automaton, error) {
	return r.MakeCharRange(c, c)
}

func (r *Automata) MakeCharRange(min, max int32) (*Automaton, error) {
	if min > max {
		return r.MakeEmpty(), nil
	}
	a := NewAutomaton()
	s1 := a.CreateState()
	s2 := a.CreateState()
	a.SetAccept(s2, true)
	if err := a.AddTransition(s1, s2, int(min), int(max)); err != nil {
		return nil, err
	}
	a.FinishState()
	return a, nil
}

func (r *Automata) MakeBinaryInterval(min []byte, minInclusive bool,
	max []byte, maxInclusive bool) (*Automaton, error) {

	if len(min) == 0 && minInclusive == false {
		return nil, errors.New("minInclusive must be true when min is null (open ended)")
	}

	if len(max) == 0 && maxInclusive == false {
		return nil, errors.New("maxInclusive must be true when max is null (open ended)")
	}

	if len(min) == 0 {
		//min = new BytesRef();
		minInclusive = true
	}

	var cmp int
	if len(max) != 0 {
		cmp = bytes.Compare(min, max)
	} else {
		cmp = -1
		if len(min) == 0 {
			if minInclusive {
				return r.MakeAnyBinary()
			} else {
				return r.MakeNonEmptyBinary()
			}
		}
	}

	if cmp == 0 {
		if minInclusive == false || maxInclusive == false {
			return r.MakeEmpty(), nil
		} else {
			return r.MakeBinary(min)
		}
	} else if cmp > 0 {
		// max < min
		return r.MakeEmpty(), nil
	}

	if len(max) != 0 &&
		bytes.HasPrefix(max, min) &&
		suffixIsZeros(max, len(min)) {

		// Finite case: no sink state!

		maxLength := len(max)

		// the == case was handled above
		//assert maxLength > min.length;

		//  bar -> bar\0+
		if maxInclusive == false {
			maxLength--
		}

		if maxLength == len(min) {
			if minInclusive == false {
				return r.MakeEmpty(), nil
			} else {
				return r.MakeBinary(min)
			}
		}

		a := NewAutomaton()
		lastState := a.CreateState()
		for i := 0; i < len(min); i++ {
			state := a.CreateState()
			label := int(min[i])
			if err := a.AddTransitionLabel(lastState, state, label); err != nil {
				return nil, err
			}
			lastState = state
		}

		if minInclusive {
			a.SetAccept(lastState, true)
		}

		for i := len(min); i < maxLength; i++ {
			state := a.CreateState()
			if err := a.AddTransitionLabel(lastState, state, 0); err != nil {
				return nil, err
			}
			a.SetAccept(state, true)
			lastState = state
		}
		a.FinishState()
		return a, nil
	}

	a := NewAutomaton()
	startState := a.CreateState()

	sinkState := a.CreateState()
	a.SetAccept(sinkState, true)

	// This state accepts all suffixes:
	if err := a.AddTransition(sinkState, sinkState, 0, 255); err != nil {
		return nil, err
	}

	equalPrefix := true
	lastState := startState
	firstMaxState := -1
	sharedPrefixLength := 0
	for i := 0; i < len(min); i++ {
		minLabel := int(min[i])

		var maxLabel int
		if len(max) != 0 && equalPrefix && i < len(max) {
			maxLabel = int(max[i])
		} else {
			maxLabel = -1
		}

		var nextState int
		if minInclusive && i == len(min)-1 && (equalPrefix == false || minLabel != maxLabel) {
			nextState = sinkState
		} else {
			nextState = a.CreateState()
		}

		if equalPrefix {

			if minLabel == maxLabel {
				// Still in shared prefix
				if err := a.AddTransitionLabel(lastState, nextState, minLabel); err != nil {
					return nil, err
				}
			} else if len(max) == 0 {
				equalPrefix = false
				sharedPrefixLength = 0
				if err := a.AddTransition(lastState, sinkState, minLabel+1, 0xff); err != nil {
					return nil, err
				}
				if err := a.AddTransitionLabel(lastState, nextState, minLabel); err != nil {
					return nil, err
				}
			} else {
				// This is the first point where min & max diverge:
				//assert maxLabel > minLabel;

				if err := a.AddTransitionLabel(lastState, nextState, minLabel); err != nil {
					return nil, err
				}

				if maxLabel > minLabel+1 {
					if err := a.AddTransition(lastState, sinkState, minLabel+1, maxLabel-1); err != nil {
						return nil, err
					}
				}

				// Now fork off path for max:
				if maxInclusive || i < len(max)-1 {
					firstMaxState = a.CreateState()
					if i < len(max)-1 {
						a.SetAccept(firstMaxState, true)
					}
					if err := a.AddTransitionLabel(lastState, firstMaxState, maxLabel); err != nil {
						return nil, err
					}
				}
				equalPrefix = false
				sharedPrefixLength = i
			}
		} else {
			// OK, already diverged:
			if err := a.AddTransitionLabel(lastState, nextState, minLabel); err != nil {
				return nil, err
			}
			if minLabel < 255 {
				if err := a.AddTransition(lastState, sinkState, minLabel+1, 255); err != nil {
					return nil, err
				}
			}
		}
		lastState = nextState
	}

	// Accept any suffix appended to the min term:
	if equalPrefix == false && lastState != sinkState && lastState != startState {
		if err := a.AddTransition(lastState, sinkState, 0, 255); err != nil {
			return nil, err
		}
	}

	if minInclusive {
		// Accept exactly the min term:
		a.SetAccept(lastState, true)
	}

	if len(max) != 0 {

		// Now do max:
		if firstMaxState == -1 {
			// Min was a full prefix of max
			sharedPrefixLength = len(min)
		} else {
			lastState = firstMaxState
			sharedPrefixLength++
		}
		for i := sharedPrefixLength; i < len(max); i++ {
			maxLabel := int(max[i])
			if maxLabel > 0 {
				if err := a.AddTransition(lastState, sinkState, 0, maxLabel-1); err != nil {
					return nil, err
				}
			}
			if maxInclusive || i < len(max)-1 {
				nextState := a.CreateState()
				if i < len(max)-1 {
					a.SetAccept(nextState, true)
				}
				if err := a.AddTransitionLabel(lastState, nextState, maxLabel); err != nil {
					return nil, err
				}
				lastState = nextState
			}
		}

		if maxInclusive {
			a.SetAccept(lastState, true)
		}
	}

	a.FinishState()

	//assert a.isDeterministic(): a.toDot();

	return a, nil
}

func suffixIsZeros(bs []byte, size int) bool {
	for _, v := range bs[size:] {
		if v != 0 {
			return false
		}
	}
	return true
}

func (r *Automata) MakeDecimalInterval(min, max, digits int) (*Automaton, error) {
	panic("")
}

func (r *Automata) MakeString(s string) (*Automaton, error) {
	a := NewAutomaton()
	lastState := a.CreateState()

	for _, v := range s {
		state := a.CreateState()
		if err := a.AddTransitionLabel(lastState, state, int(v)); err != nil {
			return nil, err
		}
		lastState = state
	}

	a.SetAccept(lastState, true)
	a.FinishState()

	return a, nil
}

func (r *Automata) MakeBinary(term []byte) (*Automaton, error) {
	a := NewAutomaton()
	lastState := a.CreateState()
	for i := 0; i < len(term); i++ {
		state := a.CreateState()
		label := int(term[i])
		if err := a.AddTransition(lastState, state, label, label); err != nil {
			return nil, err
		}
		lastState = state
	}

	a.SetAccept(lastState, true)
	a.FinishState()

	return a, nil
}
