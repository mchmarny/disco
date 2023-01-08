package trivy

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPackageFilterParsing(t *testing.T) {
	expectedResults := 568

	filter := func(v interface{}) bool {
		return false
	}

	rep, err := testPackageParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Packages))

	for _, p := range rep.Packages {
		assert.NotEmpty(t, p.Package, "package in %s is empty", p.Package)
		assert.NotEmpty(t, p.Format, "format in %s is empty", p.Package)
		assert.NotEmpty(t, p.Provider, "provider in %s is empty", p.Package)
		// assert.NotEmpty(t, p.Source, "source in %s is empty", p.Package)
	}
}

func testPackageParsing(t *testing.T, filter types.ItemFilter) (*types.PackageReport, error) {
	src := "../../../etc/data/trivy-spdx-sbom.json"
	rep, err := ParsePackages(src, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	return rep, err
}
