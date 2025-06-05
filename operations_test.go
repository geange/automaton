package automaton

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_concatenate(t *testing.T) {
	automata := &Automata{}

	a1, err := automata.MakeString("m")
	assert.Nil(t, err)
	a2, err := automata.MakeAnyString()
	assert.Nil(t, err)
	a3, err := automata.MakeString("n")
	assert.Nil(t, err)
	a4, err := automata.MakeAnyString()
	assert.Nil(t, err)

	a, err := concatenate(a1, a2, a3, a4)
	assert.Nil(t, err)
	//a, err = determinize(a, 10000)
	//assert.Nil(t, err)

	assert.True(t, Run(a, "mn"))
	assert.True(t, Run(a, "mone"))
	assert.False(t, Run(a, "m"))
}
