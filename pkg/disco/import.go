package disco

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	spdx_json "github.com/spdx/tools-golang/json"
)

func ImportSBOM(ctx context.Context, path string) error {
	if path == "" {
		return errors.New("path is required")
	}

	r, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", path)
	}
	defer r.Close()

	doc, err := spdx_json.Load2_3(r)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", path)
	}

	log.Debug().Msgf("Importing %s", doc.DocumentName)

	return nil
}
