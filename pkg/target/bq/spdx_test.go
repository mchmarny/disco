package bq

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	spdx "github.com/spdx/tools-golang/spdx/v2_3"
	"github.com/stretchr/testify/assert"
)

func TestSPDXPackageImport(t *testing.T) {
	var sbom spdx.Document
	err := types.UnmarshalFromFile("../../../etc/data/spdx23.json", &sbom)
	assert.NoError(t, err)
	assert.NotNil(t, sbom)

	rows := MakeSPDXPackageRows(&sbom)
	assert.NotNil(t, rows)
	assert.Equal(t, 100, len(rows))
	for _, r := range rows {
		assert.NotNil(t, r)
		assert.NotEmpty(t, r.Image)
		assert.NotEmpty(t, r.Sha)
		assert.NotEmpty(t, r.Format)
		assert.NotEmpty(t, r.Provider)
		assert.NotEmpty(t, r.Originator)
		assert.NotEmpty(t, r.Package)
		assert.NotEmpty(t, r.PackageVersion)
		assert.NotEmpty(t, r.Source)
		assert.NotEmpty(t, r.License)
		assert.NotEmpty(t, r.Updated)
	}
}
