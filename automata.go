package automaton

import (
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
