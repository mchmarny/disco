package disco

import (
	"context"

	"github.com/mchmarny/disco/pkg/source"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DiscoverImages discovers all deployed images in the project.
func DiscoverImages(ctx context.Context, in *types.SimpleQuery) error {
	if in == nil {
		return errors.New("nil input")
	}
	log.Debug().Msgf("discovering images with: %s", in)
	printProjectScope(in.ProjectID, "images")

	images, err := source.ImageProvider(ctx, in)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	report := types.NewItemReport(in, images...)

	if err := writeOutput(in.OutputPath, in.OutputFmt, report); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}
