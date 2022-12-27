package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashes(t *testing.T) {
	assert.Equal(t, "BgwAA2Zv", Hash("foo"))
}

func TestToKey(t *testing.T) {
	k := ToKey("t1", "t2", "t3")
	assert.Equal(t, "dDF0MnQz", k)
}
