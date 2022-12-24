package vctl

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	// JSONFormat is JSON output format.
	JSONFormat OutputFormat = iota
	// YAMLFormat is YAML output format.
	YAMLFormat OutputFormat = iota
	// RawFormat is list output format.
	RawFormat OutputFormat = iota
	// DefaultOutputFormat is default output format.
	DefaultOutputFormat = JSONFormat
)

type OutputFormat int64

type SimpleQuery struct {
	ProjectID  string
	OutputPath string
	OutputFmt  OutputFormat
}

// ParseOutputFormat parses output format.
func ParseOutputFormat(format string) (OutputFormat, error) {
	if format == "" {
		return DefaultOutputFormat, nil
	}

	switch format {
	case "json":
		return JSONFormat, nil
	case "yaml":
		return YAMLFormat, nil
	case "raw":
		return RawFormat, nil
	default:
		return DefaultOutputFormat, errors.Errorf("unsupported output format: %s", format)
	}
}

func writeOutput(path string, format OutputFormat, data any) error {
	var w io.Writer
	w = os.Stdout

	if path != "" {
		log.Info().Msgf("writing output to: '%s'", path)
		f, err := os.Create(path)
		if err != nil {
			return errors.Wrapf(err, "error creating file: %s", path)
		}
		defer f.Close()
		w = f
	}

	fmt.Println() // add a new line before

	switch format {
	case JSONFormat:
		if err := json.NewEncoder(w).Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	case YAMLFormat:
		if err := yaml.NewEncoder(w).Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	case RawFormat:
		os.Stdout.Write([]byte(fmt.Sprintf("%v", data)))
	default:
		return errors.Errorf("unsupported output format: %d", format)
	}

	return nil
}

func printProjectScope(projectID string) {
	if projectID != "" {
		log.Info().Msgf("scanning project: '%s'", projectID)
	} else {
		log.Info().Msgf("scanning all projects accessible to current user")
	}
}
