package disco

import (
	"context"

	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func DiscoverLicenses(ctx context.Context, in *types.SimpleQuery) error {
	if in == nil {
		return errors.New("nil input")
	}

	licenseFilter := func(v interface{}) bool {
		return false
	}

	log.Debug().Msgf("discovering licenses with: %s", in)
	printProjectScope(in.ProjectID)

	if err := scan(ctx, scanner.LicenseScanner, in, licenseFilter); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}
