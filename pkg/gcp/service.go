package gcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	// serviceAPIBaseURL is the base URL for GCP project API, and
	// the first parameter is project NUMBER (not ID),
	// and the second parameter is the region name.
	serviceAPIBaseURL = "https://run.googleapis.com/v1/projects/%s/locations/%s/services"
)

type serviceList struct {
	Services []*Service `json:"items"`
}

type Service struct {
	Metadata struct {
		Name string `json:"name"`
		ID   string `json:"uid"`
	} `json:"metadata"`
	Annotations map[string]string `json:"annotations"`
	Created     string            `json:"creationTimestamp"`
	Spec        struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Image string `json:"image"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			Type               string `json:"type"`
			Status             string `json:"status"`
			LastTransitionTime string `json:"lastTransitionTime"`
		} `json:"conditions"`
	} `json:"status"`
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
		return nil, errors.Wrap(err, "error client creating request")
	}

	var list serviceList
	if err := client.Get(ctx, req, &list); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return list.Services, nil
}
