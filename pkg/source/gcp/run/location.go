package run

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// regionAPIBaseURL is the base URL for GCP project API, and
	// the parameter is project NUMBER (not ID).
	regionAPIBaseURL = "https://run.googleapis.com/v1/projects/%s/locations"
)

type locaitonList struct {
	Locations []*location `json:"locations"`
}

type location struct {
	ID   string `json:"locationId"`
	Name string `json:"displayName"`
}

func getLocations(ctx context.Context, projectNumber string) ([]*location, error) {
	if projectNumber == "" {
		return nil, errors.New("project number is empty")
	}
	u := fmt.Sprintf(regionAPIBaseURL, projectNumber)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating request")
	}

	var list locaitonList
	if err := apiClient.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Locations, nil
}
