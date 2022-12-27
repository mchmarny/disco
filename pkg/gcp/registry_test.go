package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	client = &testAPIClient{}
	info, err := GetImageInfo(context.Background(), "us-west1-docker.pkg.dev/cloudy-demos/artomator/artomator@sha256:b4a094e55244bc442bdaf2a5cd06a589f754ffc8ce09183868acaa79419cd88d")
	assert.NoError(t, err)
	assert.NotNil(t, info)

	_, err = GetImageInfo(context.Background(), "")
	assert.Error(t, err)
}
