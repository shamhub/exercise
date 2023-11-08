//go:build unit
// +build unit

package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroeth(t *testing.T) {
	actualResult := Fib(0)
	assert.Equal(t, 0, actualResult)
}

func TestOneth(t *testing.T) {
	actualResult := Fib(1)
	assert.Equal(t, 1, actualResult)
}

func TestNth(t *testing.T) { // 0, 1,1, 2, 3, 5, 8
	actualResult := Fib(6)
	assert.Equal(t, 8, actualResult)
}
