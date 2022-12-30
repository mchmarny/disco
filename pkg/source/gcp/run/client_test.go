package run

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	urlExpLocations   = regexp.MustCompile(`/locations$`)
	urlEpxProjects    = regexp.MustCompile(`/projects$`)
	urlExpServices    = regexp.MustCompile(`/locations/us-west1/services$`)
	urlExpRevision    = regexp.MustCompile(`/revisions/`)
	urlExpUsage       = regexp.MustCompile(`/projects/799736955886/services$`)
	urlEpxOccurrences = regexp.MustCompile(`/occurrences$`)
	urlExpTestGet     = regexp.MustCompile(`/users/mchmarny$`)
	urlExpTestHead    = regexp.MustCompile(`api/v2/status.json$`)
)

type testAPIClient struct{}

func (t *testAPIClient) Get(ctx context.Context, req *http.Request, v any) error {
	var testFile string

	switch u := req.URL.Path; {
	case urlExpLocations.MatchString(u):
		testFile = "../../../../etc/data/test-locations.json"
	case urlEpxProjects.MatchString(u):
		testFile = "../../../../etc/data/test-projects.json"
	case urlExpServices.MatchString(u):
		testFile = "../../../../etc/data/test-services.json"
	case urlExpRevision.MatchString(u):
		testFile = "../../../../etc/data/test-revision.json"
	case urlExpUsage.MatchString(u):
		testFile = "../../../../etc/data/test-usage.json"
	case urlEpxOccurrences.MatchString(u):
		testFile = "../../../../etc/data/test-occurrences.json"
	case urlExpTestGet.MatchString(u):
		testFile = "../../../../etc/data/test-map.json"
	case urlExpTestHead.MatchString(u):
		testFile = "../../../../etc/data/test-locations.json"
	default:
		return errors.Errorf("unknown request path: %s", u)
	}

	b, err := os.ReadFile(testFile)
	if err != nil {
		return errors.Wrapf(err, "error reading test data using: %s", testFile)
	}

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(v); err != nil {
		return errors.Wrap(err, "error decoding test data")
	}
	return nil
}

func (t *testAPIClient) Head(ctx context.Context, req *http.Request, key string) (string, error) {
	return "test", nil
}

func testHTTPClientProvider(ctx context.Context, credProvider CredentialProvider) (*http.Client, error) {
	return &http.Client{}, nil
}

func TestClientGet(t *testing.T) {
	ctx := context.Background()
	httpClientProvider = testHTTPClientProvider
	r, err := http.NewRequest(http.MethodGet, "https://api.github.com/users/mchmarny", nil)
	assert.NoError(t, err)
	var d map[string]interface{}
	err = apiClient.Get(ctx, r, &d)
	assert.NoError(t, err)
}

func TestClientHead(t *testing.T) {
	ctx := context.Background()
	httpClientProvider = testHTTPClientProvider
	r, err := http.NewRequest(http.MethodHead, "https://www.githubstatus.com/api/v2/status.json", nil)
	assert.NoError(t, err, "error creating request")
	v, err := apiClient.Head(ctx, r, "x-cache")
	assert.NoError(t, err, "error executing request")
	assert.NotEmpty(t, v)
}
