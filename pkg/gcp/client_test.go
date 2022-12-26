package gcp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

var (
	urlExpLocations   = regexp.MustCompile(`/locations$`)
	urlEpxProjects    = regexp.MustCompile(`/projects$`)
	urlExpServices    = regexp.MustCompile(`/locations/us-central1/services$`)
	urlExpUsage       = regexp.MustCompile(`/projects/799736955886/services$`)
	urlEpxOccurrences = regexp.MustCompile(`/occurrences$`)
)

type testClient struct{}

func (t *testClient) Get(ctx context.Context, req *http.Request, v any) error {
	var testFile string

	switch u := req.URL.Path; {
	case urlExpLocations.MatchString(u):
		testFile = "../../etc/test-locations.json"
	case urlEpxProjects.MatchString(u):
		testFile = "../../etc/test-projects.json"
	case urlExpServices.MatchString(u):
		testFile = "../../etc/test-services.json"
	case urlExpUsage.MatchString(u):
		testFile = "../../etc/test-usage.json"
	case urlEpxOccurrences.MatchString(u):
		testFile = "../../etc/test-occurrences.json"
	default:
		return errors.Errorf("unknown request path: %s", u)
	}

	b, err := os.ReadFile(testFile)
	if err != nil {
		return errors.Wrap(err, "error reading test data")
	}

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(v); err != nil {
		return errors.Wrap(err, "error decoding test data")
	}
	return nil
}

func (t *testClient) Head(ctx context.Context, req *http.Request, key string) (string, error) {
	return "", nil
}
