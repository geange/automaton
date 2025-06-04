package automaton

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRegExp(t *testing.T) {
	regExp, err := NewRegExp("+-*(A|.....|BC)*]", WithSyntaxFlags(NONE))
	assert.Nil(t, err)
	fmt.Println(regExp)

	automaton, err := regExp.ToAutomaton(1000000)
	assert.Nil(t, err)

	fmt.Println(automaton)

	//e2, err := NewRegExp("")
	//assert.Nil(t, err)
	//fmt.Println(e2)
}
