package bq

import (
	"context"
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestLicenseImport(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.TODO()
	req := &types.ImportRequest{
		ProjectID: "cloudy-demos",
		DatasetID: "disco",
		TableID:   "licenses",
		FilePath:  "../../../etc/data/report-lic.json",
	}
	err := ImportLicenses(ctx, req)
	assert.NoError(t, err)
}
