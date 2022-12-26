package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	api = &testClient{}
	expectedProjects := 4
	list, err := GetServices(context.Background(), "799736955886", "us-central1")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, expectedProjects, len(list))
}
