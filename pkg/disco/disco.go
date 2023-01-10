package disco

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mchmarny/disco/pkg/source"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

const yamlIndent = 2

func writeOutput(path string, format types.OutputFormat, data any) error {
	if data == nil {
		return errors.New("nil data")
	}

	var w io.Writer
	w = os.Stdout

	if path != "" {
		log.Info().Msgf("writing report to: '%s'", path)
		f, err := os.Create(path)
		if err != nil {
			return errors.Wrapf(err, "error creating file: %s", path)
		}
		defer f.Close()
		w = f
	}

	fmt.Println() // add a new line before

	switch format {
	case types.JSONFormat:
		je := json.NewEncoder(w)
		je.SetIndent("", "  ")
		if err := je.Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	case types.YAMLFormat:
		ye := yaml.NewEncoder(w)
		ye.SetIndent(yamlIndent)
		if err := ye.Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	default:
		return errors.Errorf("unsupported output format: %d", format)
	}

	return nil
}

func printProjectScope(projectID, subject string) {
	if projectID != "" {
		log.Info().Msgf("scanning project: '%s' for: '%s'", projectID, subject)
	} else {
		log.Info().Msgf("scanning all projects accessible to current user for %s", subject)
	}
}

const (
	maxItemHandler = 5
)

type itemHandler func(dir, uri string, labels map[string]string) error

func handleImages(ctx context.Context, in *types.SimpleQuery, handler itemHandler) error {
	var list []*types.ImageItem
	if in.ImageFile != "" {
		var rep types.ItemReport[types.ImageItem]
		if err := types.UnmarshalFromFile(in.ImageFile, &rep); err != nil {
			return errors.Wrapf(err, "error reading image list from file: %s", in.ImageFile)
		}
		list = rep.Items
	} else {
		var err error
		list, err = source.ImageProvider(ctx, in)
		if err != nil {
			return errors.Wrap(err, "error getting images")
		}
	}

	dir, err := os.MkdirTemp(os.TempDir(), in.Kind.String())
	if err != nil {
		return errors.Wrap(err, "error creating temp dir")
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			log.Error().Err(err).Msgf("error deleting context: %s", dir)
		}
	}()

	var g errgroup.Group
	g.SetLimit(maxItemHandler)

	for _, img := range list {
		uri := img.URI
		labels := img.Context
		g.Go(func() error {
			return handler(dir, uri, labels)
		})
	}

	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "error handling images")
	}

	return nil
}
