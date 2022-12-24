package usage

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mchmarny/disco/pkg/client"
	"github.com/pkg/errors"
)

const (
	// usageAPIBaseURL is the base URL for GCP project API, and
	// the parameter is project NUMBER (not ID).
	usageAPIBaseURL = "https://serviceusage.googleapis.com/v1/projects/%s/services?filter=state:ENABLED"

	// CloudRunAPI is the name of the Cloud Run API.
	CloudRunAPI = "run.googleapis.com"
)

type serviceList struct {
	Services []struct {
		Config struct {
			Name string `json:"name"`
		} `json:"config"`
	} `json:"services"`
}

func IsAPIEnabled(ctx context.Context, projectNumber, api string) (bool, error) {
	if projectNumber == "" {
		return false, errors.New("project number is empty")
	}
	u := fmt.Sprintf(usageAPIBaseURL, projectNumber)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, errors.Wrap(err, "error client creating request")
	}

	var list serviceList
	if err := client.Request(ctx, req, &list); err != nil {
		return false, errors.Wrap(err, "error decoding response")
	}

	if len(list.Services) == 0 {
		return false, nil
	}

	for _, s := range list.Services {
		if strings.EqualFold(s.Config.Name, api) {
			return true, nil
		}
	}

	return false, nil
}
