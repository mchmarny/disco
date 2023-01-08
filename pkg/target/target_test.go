package target

import (
	"testing"

	"github.com/mchmarny/disco/pkg/target/bq"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestTargetParsing(t *testing.T) {
	req := &types.SimpleQuery{
		ProjectID: "test",
		TargetRaw: "bq://project.dataset.packages",
		Kind:      types.KindPackage,
	}

	ir, err := ParseImportRequest(req)
	assert.NoError(t, err)
	assert.NotNil(t, ir)
	assert.Equal(t, "project", ir.ProjectID)
	assert.Equal(t, "dataset", ir.DatasetID)
	assert.Equal(t, "packages", ir.TableID)
	assert.Equal(t, importDefaultLocation, ir.Location)
}

func TestTargetParsingWithoutTable(t *testing.T) {
	req := &types.SimpleQuery{
		ProjectID: "test",
		TargetRaw: "bq://project.dataset",
		Kind:      types.KindLicense,
	}

	ir, err := ParseImportRequest(req)
	assert.NoError(t, err)
	assert.NotNil(t, ir)
	assert.Equal(t, "project", ir.ProjectID)
	assert.Equal(t, "dataset", ir.DatasetID)
	assert.Equal(t, types.TableKindLicenseName, ir.TableID)
	assert.Equal(t, importDefaultLocation, ir.Location)
}

func TestTargetParsingWithoutDataset(t *testing.T) {
	req := &types.SimpleQuery{
		ProjectID: "test",
		TargetRaw: "bq://project",
		Kind:      types.KindVulnerability,
	}

	ir, err := ParseImportRequest(req)
	assert.NoError(t, err)
	assert.NotNil(t, ir)
	assert.Equal(t, "project", ir.ProjectID)
	assert.Equal(t, bq.DatasetNameDefault, ir.DatasetID)
	assert.Equal(t, types.TableKindVulnerabilityName, ir.TableID)
	assert.Equal(t, importDefaultLocation, ir.Location)
}
