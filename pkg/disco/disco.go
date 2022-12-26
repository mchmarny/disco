package disco

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mchmarny/disco/pkg/gcp"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/types"
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

var (
	getProjectsFunc   getProjects   = gcp.GetProjects
	getLocationsFunc  getLocations  = gcp.GetLocations
	getServicesFunc   getServices   = gcp.GetServices
	getImageInfoFunc  getImageInfo  = gcp.GetImageInfo
	isAPIEnabledFunc  isAPIEnabled  = gcp.IsAPIEnabled
	getCVEVulnsFunc   getCVEVulns   = gcp.GetCVEVulnerabilities
	getImageVulnsFunc getImageVulns = gcp.GetImageVulnerabilities
)

type getProjects func(ctx context.Context) ([]*gcp.Project, error)
type getLocations func(ctx context.Context, projectNumber string) ([]*gcp.Location, error)
type getServices func(ctx context.Context, projectNumber string, region string) ([]*gcp.Service, error)
type getImageInfo func(ctx context.Context, image string) (*gcp.ImageInfo, error)
type isAPIEnabled func(ctx context.Context, projectNumber string, uri string) (bool, error)
type getCVEVulns func(ctx context.Context, projectID string, cveID string) ([]*gcp.Occurrence, error)
type getImageVulns func(ctx context.Context, projectID string, imageURL string) ([]*gcp.Occurrence, error)

type OutputFormat int64

func (o OutputFormat) String() string {
	switch o {
	case JSONFormat:
		return "json"
	case YAMLFormat:
		return "yaml"
	case RawFormat:
		return "raw"
	default:
		return "unknown"
	}
}

type SimpleQuery struct {
	ProjectID  string
	OutputPath string
	OutputFmt  OutputFormat
}

func (q *SimpleQuery) String() string {
	return fmt.Sprintf("ProjectID:%s, Output:%s, Format:%s",
		q.ProjectID, q.OutputPath, q.OutputFmt)
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
	case "raw":
		return RawFormat
	default:
		log.Error().Msgf("unsupported output format: %s", format)
		return DefaultOutputFormat
	}
}

const yamlIndent = 2

func writeOutput(path string, format OutputFormat, data any) error {
	if data == nil {
		return errors.New("nil data")
	}

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
		je := json.NewEncoder(w)
		je.SetIndent("", "  ")
		if err := je.Encode(data); err != nil {
			return errors.Wrap(err, "error encoding")
		}
	case YAMLFormat:
		ye := yaml.NewEncoder(w)
		ye.SetIndent(yamlIndent)
		if err := ye.Encode(data); err != nil {
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

func scan(ctx context.Context, scan scanner.ScannerType, in *SimpleQuery, filter types.ItemFilter) error {
	if in == nil {
		return errors.New("nil input")
	}

	images, err := getDeployedImages(ctx, in.ProjectID)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	dir, err := os.MkdirTemp(os.TempDir(), scan.String())
	if err != nil {
		return errors.Wrap(err, "error creating temp dir")
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			log.Error().Err(err).Msgf("error deleting context: %s", dir)
		}
	}()

	list := make([]any, 0)
	for _, img := range images {
		log.Debug().Msgf("getting %s for %s", scan.String(), img.Image.Name)
		p := path.Join(dir, img.Image.Name)

		switch scan {
		case scanner.LicenseScanner:
			report, err := scanner.GetLicenses(img.Image.URI(), p, filter)
			if err != nil {
				return errors.Wrapf(err, "error getting licenses for %s", img.Image.Name)
			}
			log.Info().Msgf("found %d licenses in %s", len(report.Licenses), img.Image.Name)
			if len(report.Licenses) > 0 {
				list = append(list, report)
			}
		case scanner.VulnerabilityScanner:
			report, err := scanner.GetVulnerabilities(img.Image.URI(), p, filter)
			if err != nil {
				return errors.Wrapf(err, "error getting vulnerabilities for %s", img.Image.Name)
			}
			log.Info().Msgf("found %d vulnerabilities in %s", len(report.Vulnerabilities), img.Image.Name)
			if len(report.Vulnerabilities) > 0 {
				list = append(list, report)
			}
		default:
			return errors.Errorf("unsupported scanner: %s", scan)
		}
	}

	if err := writeOutput(in.OutputPath, in.OutputFmt, list); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}
