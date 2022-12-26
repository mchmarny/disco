package disco

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/mchmarny/disco/pkg/gcp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDisco(t *testing.T) {
	getProjectsFunc = getTestProjects
	getLocationsFunc = getTestLocations
	getServicesFunc = getTestServices
	getImageInfoFunc = getTestImageInfo
	getCVEVulnsFunc = getTestCVEVulns
	getImageVulnsFunc = getTestImageVulns
	isAPIEnabledFunc = isTestAPIEnabled

	ctx := context.Background()

	err := DiscoverImages(ctx, &ImagesQuery{})
	assert.NoError(t, err, "error discovering images")

	err = DiscoverLicense(ctx, &SimpleQuery{})
	assert.NoError(t, err, "error discovering license")

	err = DiscoverVulns(ctx, &VulnsQuery{})
	assert.NoError(t, err, "error discovering vulns")
}

func getTestProjects(ctx context.Context) ([]*gcp.Project, error) {
	list := []*gcp.Project{
		{
			Number: "123456789",
			ID:     "test-project",
			State:  "ACTIVE",
		},
	}

	return list, nil
}

func getTestLocations(ctx context.Context, projectNumber string) ([]*gcp.Location, error) {
	list := []*gcp.Location{
		{
			ID:   "us-central1",
			Name: "us-central1",
		},
		{
			ID:   "us-east1",
			Name: "us-east1",
		},
	}
	return list, nil
}

func getTestServices(ctx context.Context, projectNumber string, region string) ([]*gcp.Service, error) {
	var list []*gcp.Service
	if err := loadTestData("../../etc/test-service.json", &list); err != nil {
		return nil, err
	}
	return list, nil
}

func getTestImageInfo(ctx context.Context, image string) (*gcp.ImageInfo, error) {
	img, err := gcp.ParseImageInfo("us-docker.pkg.dev/cloudy-demos/art/artomator@sha256:1234567890")
	if err != nil {
		return nil, errors.Wrap(err, "error parsing image info")
	}
	return img, nil
}

func getTestCVEVulns(ctx context.Context, projectID string, cveID string) ([]*gcp.Occurrence, error) {
	var list []*gcp.Occurrence
	if err := loadTestData("../../etc/test-occurrence.json", &list); err != nil {
		return nil, err
	}
	return list, nil
}

func getTestImageVulns(ctx context.Context, projectID string, imageURL string) ([]*gcp.Occurrence, error) {
	var list []*gcp.Occurrence
	if err := loadTestData("../../etc/test-occurrence.json", &list); err != nil {
		return nil, err
	}
	return list, nil
}

func isTestAPIEnabled(ctx context.Context, projectNumber string, uri string) (bool, error) {
	return false, nil
}

func loadTestData(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "error reading test data")
	}

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(v); err != nil {
		return errors.Wrap(err, "error decoding test data")
	}
	return nil
}
