package gcp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	// serviceAPIBaseURL is the base URL for Cloud Run service.
	serviceAPIBaseURL = "https://run.googleapis.com/v2/projects/%s/locations/%s/services"

	// revisionAPIBaseURL is the base URL for Cloud Run revision.
	revisionAPIBaseURL = "https://run.googleapis.com/v2/%s"
)

type serviceList struct {
	Services []*Service `json:"services"`
}

type Service struct {
	Name       string       `json:"name"`
	FullName   string       `json:"fullName"`
	Revision   string       `json:"latestReadyRevision"`
	Containers []*Container `json:"containers"`
}

type Revision struct {
	Conditions []*Container `json:"containers"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func GetServices(ctx context.Context, projectID, region string) ([]*Service, error) {
	if projectID == "" {
		return nil, errors.New("projectID is empty")
	}
	if region == "" {
		return nil, errors.New("region is empty")
	}

	u := fmt.Sprintf(serviceAPIBaseURL, projectID, region)
	log.Debug().Msgf("getting services from: %s", u)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating service request")
	}

	var list serviceList
	if err := client.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding service response")
	}

	if len(list.Services) == 0 {
		log.Debug().Msgf("no services found in project %s", projectID)
		return list.Services, nil
	}

	// add revision images
	for _, s := range list.Services {
		if s.Revision == "" {
			return nil, errors.Errorf("service %s has no revision", s.Name)
		}
		u := fmt.Sprintf(revisionAPIBaseURL, s.Revision)
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error client creating revision request")
		}
		var rev Revision
		if err := client.Get(ctx, req, &rev); err != nil {
			return nil, errors.Wrapf(err, "error decoding revision response from: %s", u)
		}
		s.FullName = s.Name
		s.Name = parseServiceName(s.Name)
		s.Containers = rev.Conditions
	}

	return list.Services, nil
}

func parseServiceName(s string) string {
	return s[strings.LastIndex(s, "/")+1:]
}
