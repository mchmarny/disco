package types

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	// JSONFormat is JSON output format.
	JSONFormat OutputFormat = iota
	// YAMLFormat is YAML output format.
	YAMLFormat OutputFormat = iota
	// DefaultOutputFormat is default output format.
	DefaultOutputFormat = JSONFormat
)

type SimpleQuery struct {
	Version    string
	Kind       DiscoKind
	ProjectID  string
	OutputPath string
	OutputFmt  OutputFormat
	ImageFile  string
	ImageURI   string
	TargetRaw  string
	Target     *ImportRequest
	Quiet      bool
	Bucket     string
}

func (q *SimpleQuery) Validate() error {
	if q.ImageFile != "" && q.ImageURI != "" {
		return errors.New("only one of image file or image URI can be specified")
	}

	if q.ProjectID != "" && (q.ImageFile != "" || q.ImageURI != "") {
		return errors.New("project ID can only be used during image discovery")
	}

	return nil
}

func (q *SimpleQuery) String() string {
	return fmt.Sprintf("project: %s, output: %s, format: %s, source: %s, uri: %s, target: %s, backet: %s",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.ImageFile, q.ImageURI, q.Target, q.Bucket)
}

type OutputFormat int64

func (o OutputFormat) String() string {
	switch o {
	case JSONFormat:
		return "json"
	case YAMLFormat:
		return "yaml"
	default:
		return "unknown"
	}
}

// ParseOutputFormat parses output format.
func ParseOutputFormatOrDefault(format string) OutputFormat {
	if format == "" {
		return DefaultOutputFormat
	}

	switch format {
	case "json":
		return JSONFormat
	case "yaml":
		return YAMLFormat
	default:
		log.Error().Msgf("unsupported output format: %s", format)
		return DefaultOutputFormat
	}
}
