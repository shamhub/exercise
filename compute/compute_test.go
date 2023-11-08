//go:build unit
// +build unit

package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleDigit(t *testing.T) {
	result, err := SumDigits(9)

	assert.Equal(t, result, 9)
	assert.Nil(t, err)
}

func TestZeroDigit(t *testing.T) {
	result, err := SumDigits(0)

	assert.Equal(t, result, 0)
	assert.Nil(t, err)
}

func TestNegativeN(t *testing.T) {
	result, err := SumDigits(-10)

	assert.Equal(t, result, -1)
	assert.NotNil(t, err)
}

func TestSplit(t *testing.T) {
	allButLast, last := split(10)
	assert.Equal(t, 0, last)
	assert.Equal(t, 1, allButLast)

	allButLast, last = split(9)
	assert.Equal(t, 9, last)
	assert.Equal(t, 0, allButLast)
}

func TestFactZero(t *testing.T) {
	result, err := Fact(0)
	assert.Equal(t, 1, result)
	assert.Nil(t, err)
}
