package disco

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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
		log.Info().Msgf("scanning all projects accessible to current user for: '%s'", subject)
	}
}

func readImageList(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening file: %s", path)
	}
	defer f.Close()

	var images []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		uri := scanner.Text()
		if uri != "" {
			images = append(images, uri)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "error reading file: %s", path)
	}

	return images, nil
}

func getImagesURIs(ctx context.Context, in *types.SimpleQuery) ([]string, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	if in.ImageURI != "" {
		log.Debug().Msgf("using image URI: '%s'", in.ImageURI)
		return []string{in.ImageURI}, nil
	}

	if in.ImageFile != "" {
		log.Info().Msgf("reading image list from: '%s'", in.ImageFile)
		return readImageList(in.ImageFile)
	}

	log.Debug().Msg("discovering images from API...")
	list, err := discoverImageURIs(ctx, &types.ImagesQuery{
		SimpleQuery: *in,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error discovering images")
	}

	log.Debug().Msgf("found %d images", len(list))
	return list, nil
}

type itemHandler func(dir, uri string) error

func handleImages(ctx context.Context, in *types.SimpleQuery, handler itemHandler) error {
	list, err := getImagesURIs(ctx, in)
	if err != nil {
		return errors.Wrap(err, "error getting images")
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

	for _, img := range list {
		if handler(dir, img) != nil {
			return errors.Wrap(err, "error handling image")
		}
	}

	return nil
}
