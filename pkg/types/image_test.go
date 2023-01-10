package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {
	i := NewItemReport[LicenseReport](&SimpleQuery{
		Kind:    KindLicense,
		Version: "v1.2.3",
	})
	assert.NotNil(t, i)
	assert.NotNil(t, i.Meta)
	assert.Equal(t, i.Meta.Kind, KindLicense.String())
	assert.Equal(t, i.Meta.Version, "v1.2.3")
}
