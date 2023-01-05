package trivy

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestLicenseFilterParsing(t *testing.T) {
	expectedResults := 3

	filter := func(v interface{}) bool {
		return false
	}

	rep, err := testLicenseParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Licenses))
}

func testLicenseParsing(t *testing.T, filter types.ItemFilter) (*types.LicenseReport, error) {
	src := "../../../etc/data/test-license.json"
	rep, err := ParseLicenses(src, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	return rep, err
}
