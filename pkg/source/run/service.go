package run

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

	managedByGCF = "cloudfunctions"

	runtimeGCF = "gcf"
	runtimeRun = "run"
)

type serviceList struct {
	Services []*service `json:"services"`
}

type service struct {
	Name       string       `json:"name"`
	FullName   string       `json:"fullName"`
	Revision   string       `json:"latestReadyRevision"`
	Containers []*container `json:"containers"`
	Runtime    string       `json:"runtime"`
}

type revision struct {
	Labels struct {
		ManagedBy string `json:"goog-managed-by"` //nolint:tagliatelle
	} `json:"labels"`
	Conditions []*container `json:"containers"`
}

type container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func getServices(ctx context.Context, projectID, region string) ([]*service, error) {
	if projectID == "" {
		return nil, errors.New("projectID is empty")
	}
	if region == "" {
		return nil, errors.New("region is empty")
	}

	u := fmt.Sprintf(serviceAPIBaseURL, projectID, region)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error client creating service request")
	}

	var list serviceList
	if err := apiClient.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding service response")
	}

	if len(list.Services) == 0 {
		return list.Services, nil
	}

	// add revision images
	for _, s := range list.Services {
		if s.Revision == "" {
			log.Logger.Debug().Msgf("service %s has no revision", s.Name)
			continue
		}
		u := fmt.Sprintf(revisionAPIBaseURL, s.Revision)
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error client creating revision request")
		}
		var rev revision
		if err := apiClient.Get(ctx, req, &rev); err != nil {
			return nil, errors.Wrapf(err, "error decoding revision response from: %s", u)
		}
		s.FullName = s.Name
		s.Name = parseServiceName(s.Name)
		s.Containers = rev.Conditions
		if rev.Labels.ManagedBy == managedByGCF {
			s.Runtime = runtimeGCF
		} else {
			s.Runtime = runtimeRun
		}
	}

	return list.Services, nil
}

func parseServiceName(s string) string {
	return s[strings.LastIndex(s, "/")+1:]
}
