package syft

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestLicenseFilterParsing(t *testing.T) {
	expectedResults := 100

	filter := func(v interface{}) bool {
		return false
	}

	rep, err := testLicenseParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Packages))

	for _, p := range rep.Packages {
		assert.NotEmpty(t, p.Package)
		assert.NotEmpty(t, p.PackageVersion)
		assert.NotEmpty(t, p.Format)
		assert.NotEmpty(t, p.Provider)
		assert.NotEmpty(t, p.Originator)
		assert.NotEmpty(t, p.Source)
		assert.NotEmpty(t, p.License)
	}
}

func testLicenseParsing(t *testing.T, filter types.ItemFilter) (*types.PackageReport, error) {
	src := "../../../etc/data/spdx23.json"
	rep, err := ParsePackages(src, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	return rep, err
}
