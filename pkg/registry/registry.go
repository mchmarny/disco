package registry

import (
	"context"
	"net/http"

	"github.com/mchmarny/disco/pkg/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	digestHeaderKey = "docker-content-digest"
)

func GetImageInfo(ctx context.Context, image string) (*ImageInfo, error) {
	if image == "" {
		return nil, errors.New("image cannot be empty")
	}

	info, err := ParseImageInfo(image)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing image")
	}

	// if already has a digest, return as is
	if info.Digest != "" {
		return info, nil
	}

	u := info.ManifestURL()
	req, err := http.NewRequest(http.MethodHead, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating registry request")
	}

	// add accept header for DOcker API v2
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	log.Debug().Msgf("getting image digest from %s", u)

	val, err := client.RequestHead(ctx, req, digestHeaderKey)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding registry response")
	}

	if val == "" {
		return nil, errors.Errorf("no digest found in response from: %s", u)
	}

	info.Digest = val

	return info, nil
}
