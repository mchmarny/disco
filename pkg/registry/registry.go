package registry

import (
	"context"
	"net/http"

	"github.com/mchmarny/vctl/pkg/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type manifest struct {
	Config struct {
		Digest string `json:"digest"`
		Size   int64  `json:"size"`
	} `json:"config"`
}

func GetImageInfo(ctx context.Context, image string) (*ImageInfo, error) {
	if image == "" {
		return nil, errors.New("Image cannot be empty.")
	}

	info, err := ParseImageInfo(image)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing image.")
	}

	// if already has a digest, return as is
	if info.Digest != "" {
		return info, nil
	}

	u := info.ManifestURL()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating registry request.")
	}

	// add accept header for DOcker API v2
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	log.Debug().Msgf("Getting image digest from %s.", u)

	var m manifest
	if err := client.Request(ctx, req, &m); err != nil {
		return nil, errors.Wrap(err, "Error decoding registry response.")
	}

	if m.Config.Digest == "" {
		return nil, errors.Errorf("Unable to find digest for this image using %s.", u)
	}

	info.Digest = m.Config.Digest

	return info, nil
}
