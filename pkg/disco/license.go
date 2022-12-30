package disco

import (
	"context"
	"path"

	"github.com/google/uuid"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func DiscoverLicenses(ctx context.Context, in *types.SimpleQuery) error {
	if in == nil {
		return errors.New("nil input")
	}

	f := func(v interface{}) bool {
		return false
	}

	log.Debug().Msgf("discovering licenses with: %s", in)
	printProjectScope(in.ProjectID)

	if err := scanLicenses(ctx, in, f); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}

func scanLicenses(ctx context.Context, in *types.SimpleQuery, filter types.ItemFilter) error {
	results := make([]*types.LicenseReport, 0)
	h := func(dir, uri string) error {
		p := path.Join(dir, uuid.NewString())
		log.Debug().Msgf("getting %s for %s (file: %s)", types.KindLicenseName, uri, p)

		rez, err := scanner.GetLicenses(uri, p, filter)
		if err != nil {
			return errors.Wrapf(err, "error getting licenses for %s", uri)
		}
		log.Info().Msgf("found %d licenses in %s", len(rez.Licenses), uri)
		if len(rez.Licenses) > 0 {
			results = append(results, rez)
		}
		return nil
	}

	if err := handleImages(ctx, in, h); err != nil {
		return errors.Wrap(err, "error handling images")
	}

	report := types.NewItemReport(in, results...)

	if err := writeOutput(in.OutputPath, in.OutputFmt, report); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}
