package bq

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicenseImport(t *testing.T) {
	ctx := context.TODO()
	obj := "gs://disco-cloudy-demos/lic-2023-01-03T00-03-00.json"
	req := &ImportRequest{
		ProjectID: "cloudy-demos",
		DatasetID: "disco",
		TableID:   "licenses",
		ObjectURI: obj,
	}
	err := ImportJSON(ctx, req)
	assert.NoError(t, err)
}
