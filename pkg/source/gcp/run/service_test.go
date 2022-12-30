package run

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	client = &testAPIClient{}
	expectedProjects := 4
	list, err := getServices(context.Background(), "cloudy-demos", "us-west1")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, expectedProjects, len(list))
}
