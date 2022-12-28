package gcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// serviceAPIBaseURL is the base URL for Cloud Run service.
	// The first parameter is project NUMBER (not ID),
	// and the second parameter is the region name.
	serviceAPIBaseURL = "https://run.googleapis.com/v1/projects/%s/locations/%s/services"

	// revisionAPIBaseURL is the base URL for Cloud Run revision.
	revisionAPIBaseURL = "https://run.googleapis.com/v2/projects/%s/locations/%s/services/%s/revisions/%s"
)

type serviceList struct {
	Services []*Service `json:"items"`
}

type Service struct {
	Metadata struct {
		Name string `json:"name"`
		ID   string `json:"uid"`
	} `json:"metadata"`
	Containers []*Container `json:"containers"`
	Status     struct {
		Revision string `json:"latestReadyRevisionName"`
	}
}

type Revision struct {
	Conditions []*Container `json:"containers"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func GetServices(ctx context.Context, projectNumber, region string) ([]*Service, error) {
	if projectNumber == "" {
		return nil, errors.New("project number is empty")
	}
	if region == "" {
		return nil, errors.New("region is empty")
	}

	u := fmt.Sprintf(serviceAPIBaseURL, projectNumber, region)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating service request")
	}

	var list serviceList
	if err := client.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding service response")
	}

	if len(list.Services) == 0 {
		return list.Services, nil
	}

	// add revision images
	for _, s := range list.Services {
		if s.Status.Revision == "" {
			return nil, errors.Errorf("service %s has no revision", s.Metadata.Name)
		}
		u := fmt.Sprintf(revisionAPIBaseURL, projectNumber, region, s.Metadata.Name, s.Status.Revision)
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error client creating revision request")
		}
		var rev Revision
		if err := client.Get(ctx, req, &rev); err != nil {
			return nil, errors.Wrapf(err, "error decoding revision response from: %s", u)
		}
		s.Containers = rev.Conditions
	}

	return list.Services, nil
}
