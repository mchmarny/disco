package trivy

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPackageFilterParsing(t *testing.T) {
	expectedResults := 100

	filter := func(v interface{}) bool {
		return false
	}

	rep, err := testPackageParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Packages))

	for _, p := range rep.Packages {
		assert.NotEmpty(t, p.Package)
		assert.NotEmpty(t, p.PackageVersion)
		assert.NotEmpty(t, p.Format)
		assert.NotEmpty(t, p.Provider)
		assert.NotEmpty(t, p.Source)
	}
}

func testPackageParsing(t *testing.T, filter types.ItemFilter) (*types.PackageReport, error) {
	src := "../../../etc/data/spdx23.json"
	rep, err := ParsePackages(src, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	return rep, err
}
