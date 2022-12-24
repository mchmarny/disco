package vctl

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type SimpleQuery struct {
	ProjectID  string
	OutputPath string
}

func writeOutput(path string, data any) error {
	fmt.Println()
	w := os.Stdout
	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			return errors.Wrapf(err, "error creating file: %s", path)
		}
		defer f.Close()
		w = f
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return errors.Wrap(err, "error encoding")
	}
	return nil
}

func printProjectScope(projectID string) {
	if projectID != "" {
		log.Info().Msgf("Scanning project: '%s'.", projectID)
	} else {
		log.Info().Msgf("Scanning all projects accessible to current user.")
	}
}
