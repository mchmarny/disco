package gcp

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
	APIClient Client = &GCPClient{}
)

type Client interface {
	Get(ctx context.Context, req *http.Request, v any) error
	Head(ctx context.Context, req *http.Request, key string) (string, error)
}

type GCPClient struct {
	Anon bool
}

func (g *GCPClient) Get(ctx context.Context, req *http.Request, v any) error {
	c, err := g.newClient(ctx)
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

func (g *GCPClient) Head(ctx context.Context, req *http.Request, key string) (string, error) {
	c, err := g.newClient(ctx)
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

func (g *GCPClient) newClient(ctx context.Context) (*http.Client, error) {
	var ops []option.ClientOption

	if !g.Anon {
		creds, err := g.getCredentials(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credentials")
		}
		ops = append(ops, option.WithCredentials(creds))
	}

	client, _, err := htransport.NewClient(ctx, ops...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http client")
	}

	return client, nil
}

func (g *GCPClient) getCredentials(ctx context.Context) (*google.Credentials, error) {
	credentials, err := google.FindDefaultCredentials(ctx, scopeDefault)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create default credentials")
	}
	return credentials, nil
}
