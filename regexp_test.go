package automaton

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRegExp(t *testing.T) {
	regExp, err := NewRegExp("[A-Z][a-z]*")
	assert.Nil(t, err)
	fmt.Println(regExp)

	//e2, err := NewRegExp("")
	//assert.Nil(t, err)
	//fmt.Println(e2)
}
