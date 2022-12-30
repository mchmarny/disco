package run

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	client = &testAPIClient{}
	expectedProjects := 4
	list, err := getProjects(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, expectedProjects, len(list))
}
