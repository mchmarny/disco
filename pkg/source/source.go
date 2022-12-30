package source

import (
	"context"

	"github.com/mchmarny/disco/pkg/source/gcp/run"
	"github.com/mchmarny/disco/pkg/types"
)

var (
	ImageProvider SourceImageProvider = run.GetImages
)

type SourceImageProvider func(ctx context.Context, in *types.ImagesQuery) ([]*types.ImageItem, error)
