//go:build unit
// +build unit

package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroeth(t *testing.T) {
	actualResult := Fib(1)
	assert.Equal(t, 1, actualResult)
}
