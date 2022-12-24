package region

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mchmarny/vctl/pkg/client"
	"github.com/pkg/errors"
)

const (
	// regionAPIBaseURL is the base URL for GCP project API, and
	// the parameter is project NUMBER (not ID).
	regionAPIBaseURL = "https://run.googleapis.com/v1/projects/%s/locations"
)

type RegionList struct {
	Regions []*Region `json:"locations"`
}

type Region struct {
	ID   string `json:"locationId"`
	Name string `json:"displayName"`
}

func GetRegions(ctx context.Context, projectNumber string) ([]*Region, error) {
	if projectNumber == "" {
		return nil, errors.New("project number is empty")
	}
	u := fmt.Sprintf(regionAPIBaseURL, projectNumber)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating request")
	}

	var list RegionList
	if err := client.Request(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Regions, nil
}
