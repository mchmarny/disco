package disco

import (
	"context"
	"errors"
	"testing"

	"github.com/mchmarny/disco/pkg/gcp"
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
	return nil, errors.New("not implemented")
}

func getTestImageInfo(ctx context.Context, image string) (*gcp.ImageInfo, error) {
	return nil, errors.New("not implemented")
}

func getTestCVEVulns(ctx context.Context, projectID string, cveID string) ([]*gcp.Occurrence, error) {
	return nil, errors.New("not implemented")
}

func getTestImageVulns(ctx context.Context, projectID string, imageURL string) ([]*gcp.Occurrence, error) {
	return nil, errors.New("not implemented")
}

func isTestAPIEnabled(ctx context.Context, projectNumber string, uri string) (bool, error) {
	return false, errors.New("not implemented")
}
