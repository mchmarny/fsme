package lighter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	id := GetNewID()
	assert.NotNil(t, id)
	assert.True(t, IsValidID(id))
}

func TestToID(t *testing.T) {
	id := ToID("1234567")
	assert.NotNil(t, id)
	assert.True(t, IsValidID(id))
}

func TestIsValidID(t *testing.T) {
	assert.False(t, IsValidID("1234567"))
	assert.True(t, IsValidID("a1234567"))
}
