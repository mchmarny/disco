package bq

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/assert"
)

func TestCycloneDXPackageImport(t *testing.T) {
	b, err := os.ReadFile("../../../etc/data/cyclonedx12.json")
	assert.NoError(t, err)

	var sbom cyclonedx.BOM
	err = json.Unmarshal(b, &sbom)
	assert.NoError(t, err)
	assert.NotNil(t, sbom)

	rows := MakeCycloneDXPackageRows(&sbom)
	assert.NotNil(t, rows)
	assert.Equal(t, 292, len(rows))
	for i, r := range rows {
		assert.NotNil(t, r)
		assert.NotEmpty(t, r.Image, "row %d", i)
		assert.NotEmpty(t, r.Sha, "row %d", i)
		assert.NotEmpty(t, r.Format, "row %d", i)
		assert.NotEmpty(t, r.Provider, "row %d", i)
		assert.NotEmpty(t, r.Package, "row %d", i)
		assert.NotEmpty(t, r.PackageVersion, "row %d", i)
		assert.NotEmpty(t, r.Updated, "row %d", i)
	}
}
