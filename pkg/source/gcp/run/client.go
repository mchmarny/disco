package run

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

const (
	scopeDefault    = "https://www.googleapis.com/auth/cloud-platform"
	digestHeaderKey = "docker-content-digest"
)

var (
	apiClient GoogleAPIClient = &googleAPIClient{
		credProvider: getDefaultCredentials,
	}
	httpClientProvider ClientProvider = newHTTPClientWithCredentials
)

type CredentialProvider func(ctx context.Context) (*google.Credentials, error)
type ClientProvider func(ctx context.Context, credProvider CredentialProvider) (*http.Client, error)

type GoogleAPIClient interface {
	Get(ctx context.Context, req *http.Request, v any) error
	Head(ctx context.Context, req *http.Request, key string) (string, error)
}

type googleAPIClient struct {
	credProvider CredentialProvider
}

func (g *googleAPIClient) Get(ctx context.Context, req *http.Request, v any) error {
	c, err := httpClientProvider(ctx, g.credProvider)
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

func (g *googleAPIClient) Head(ctx context.Context, req *http.Request, key string) (string, error) {
	c, err := httpClientProvider(ctx, g.credProvider)
	if err != nil {
		return "", errors.Wrap(err, "error creating client")
	}

	r, err := c.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error executing request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		log.Error().Msgf("error getting head: %s", r.Status)
		return "", errors.Wrapf(err, "error getting projects: %s", r.Status)
	}

	v := r.Header.Get(key)
	log.Debug().Msgf("key: %s", key)

	if v == "" {
		list := r.Header.Values(key)
		log.Debug().Msgf("values: %v", list)
		if len(list) > 0 {
			v = list[0]
		}
	}

	return v, nil
}

func newHTTPClientWithCredentials(ctx context.Context, credProvider CredentialProvider) (*http.Client, error) {
	var ops []option.ClientOption
	var client *http.Client

	if credProvider != nil {
		creds, err := credProvider(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credentials")
		}
		ops = append(ops, option.WithCredentials(creds))
		c, _, err := htransport.NewClient(ctx, ops...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create http client")
		}
		client = c
	} else {
		client = &http.Client{}
	}

	return client, nil
}

func getDefaultCredentials(ctx context.Context) (*google.Credentials, error) {
	credentials, err := google.FindDefaultCredentials(ctx, scopeDefault)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create default credentials")
	}
	return credentials, nil
}
