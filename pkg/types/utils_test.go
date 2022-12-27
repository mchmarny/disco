package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashes(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("BgwAA2Zv", string(Hash("foo")))
}
