package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	api = &testClient{}
	expectedProjects := 4
	list, err := GetProjects(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, expectedProjects, len(list))
}
