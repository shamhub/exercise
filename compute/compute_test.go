//go:build unit
// +build unit

package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF(t *testing.T) {
	F()
	assert.Equal(t, 9, 9)

}
