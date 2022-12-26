package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVulns(t *testing.T) {
	APIClient = &TestAPIClient{}
	expectedProjects := 1
	list, err := GetCVEVulnerabilities(context.Background(), "799736955886", "CVE-2022-3996")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, len(list), expectedProjects)
}

func TestImages(t *testing.T) {
	APIClient = &TestAPIClient{}
	expectedProjects := 1
	list, err := GetImageVulnerabilities(context.Background(), "799736955886", "https://gcr.io/cloudy-demos/hello-broken@sha256:0900c08e7d40f9485c8497c035de07391ba3c274a1035f504f8602531b2314e6")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.GreaterOrEqual(t, len(list), expectedProjects)
}
