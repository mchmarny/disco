package target

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestTargetParsing(t *testing.T) {
	req := &types.SimpleQuery{
		ProjectID: "test",
		TargetRaw: "bq://project.dataset.packages",
		Kind:      types.KindPackage,
	}

	// cloudy-demos.disco.packages

	ir, err := ParseImportRequest(req)
	assert.NoError(t, err)
	assert.NotNil(t, ir)
	assert.Equal(t, "project", ir.ProjectID)
	assert.Equal(t, "dataset", ir.DatasetID)
	assert.Equal(t, "packages", ir.TableID)
	assert.Equal(t, importDefaultLocation, ir.Location)
}
