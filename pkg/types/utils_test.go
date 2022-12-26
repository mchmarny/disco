package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashes(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Hash("foo"), Hash("foo"))
	assert.NotEqual(Hash("foo"), Hash("bar"))
}
