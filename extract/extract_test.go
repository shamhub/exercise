//go:build unit
// +build unit

package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadEntry(t *testing.T) {
	reader := NewJSONReader("../data.json")
	assert.NotNil(t, reader)
	actualData := reader.ReadEntry()
	expectedDateTime := "2023-06-27 22:22:19.62710.192501"
	assert.Equal(t, actualData.DateTime, expectedDateTime)
}
