package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

const (
	// ScopeDefault is the default scope for the client.
	ScopeDefault = "https://www.googleapis.com/auth/cloud-platform"
)

// NewClient creates a new http client using the default credentials.
func NewClient(ctx context.Context) (*http.Client, error) {
	creds, err := getCredentials(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create credentials")
	}

	client, _, err := htransport.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http client")
	}

	return client, nil
}

func Request(ctx context.Context, req *http.Request, v any) error {
	c, err := NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating client")
	}

	r, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "error getting projects: %s", r.Status)
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.Wrap(err, "error decoding response")
	}
	return nil
}

func RequestHead(ctx context.Context, req *http.Request, key string) (string, error) {
	c, err := NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error creating client")
	}

	r, err := c.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error executing request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", errors.Wrapf(err, "error getting projects: %s", r.Status)
	}

	v := r.Header.Get(key)

	return v, nil
}

func getCredentials(ctx context.Context) (*google.Credentials, error) {
	credentials, err := google.FindDefaultCredentials(ctx, ScopeDefault)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create default credentials")
	}
	return credentials, nil
}
