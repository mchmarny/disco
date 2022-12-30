package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOutputFormat(t *testing.T) {
	f := ParseOutputFormatOrDefault("bad")
	assert.Equal(t, f, DefaultOutputFormat)
	f = ParseOutputFormatOrDefault("json")
	assert.Equal(t, f, JSONFormat)
	f = ParseOutputFormatOrDefault("yaml")
	assert.Equal(t, f, YAMLFormat)
}
