package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocations(t *testing.T) {
	client = &testAPIClient{}
	expectedRegions := 36
	list, err := GetLocations(context.Background(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, expectedRegions, len(list))
}
