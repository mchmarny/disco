package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoKind(t *testing.T) {
	f := KindUndefined
	assert.Equal(t, f.String(), KindUndefinedName)
	f = KindImage
	assert.Equal(t, f.String(), KindImageName)
	f = KindVulnerability
	assert.Equal(t, f.String(), KindVulnerabilityName)
}
