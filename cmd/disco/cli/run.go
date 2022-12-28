package cli

import (
	"context"

	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

var (
	projectIDFlag = &c.StringFlag{
		Name:     "project",
		Aliases:  []string{"p"},
		Usage:    "project ID",
		Required: false,
	}

	outputPathFlag = &c.StringFlag{
		Name:     "output",
		Aliases:  []string{"o"},
		Usage:    "path where to save the output",
		Required: false,
	}

	outputFormatFlag = &c.StringFlag{
		Name:     "format",
		Aliases:  []string{"f"},
		Usage:    "output format (json, yaml)",
		Required: false,
	}

	outputURIOnlyFlag = &c.BoolFlag{
		Name:  "uri",
		Usage: "output only image URI",
		Value: false,
	}

	caAPIExecFlag = &c.BoolFlag{
		Name:  "use-ca",
		Usage: "invokes Container Analysis API instead of the local scanner",
		Value: false,
	}

	cveFlag = &c.StringFlag{
		Name:     "cve",
		Aliases:  []string{"e"},
		Usage:    "exposure ID (CVE number, e.g. CVE-2019-19378)",
		Required: false,
	}

	minSeverityFlag = &c.StringFlag{
		Name:     "min-severity",
		Aliases:  []string{"min-sev", "ms"},
		Usage:    "minimum severity to report (e.g. low, medium, high, critical)",
		Required: false,
	}

	imageListPathFlag = &c.StringFlag{
		Name:     "source",
		Aliases:  []string{"s", "src"},
		Usage:    "image list input file to serve as source (instead of discovery)",
		Required: false,
	}

	imageURIFlag = &c.StringFlag{
		Name:     "image",
		Aliases:  []string{"i", "img"},
		Usage:    "specific image URI to scan",
		Required: false,
	}

	runCmd = &c.Command{
		Name:  "run",
		Usage: "Cloud Run commands",
		Subcommands: []*c.Command{
			{
				Name:    "images",
				Aliases: []string{"img", "i"},
				Usage:   "List deployed container images",
				Action:  runImagesCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					outputURIOnlyFlag,
				},
			},
			{
				Name:    "vulnerabilities",
				Aliases: []string{"vul", "v"},
				Usage:   "Check for OS-level exposures in deployed images (supports specific CVE filter)",
				Action:  runVulnsCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					imageListPathFlag,
					imageURIFlag,
					minSeverityFlag,
					cveFlag,
					caAPIExecFlag,
				},
			},
			{
				Name:    "licenses",
				Aliases: []string{"lic", "l"},
				Usage:   "Scans images for license types (requires OSS scanner, e.g. Trivy)",
				Action:  runLicenseCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					imageListPathFlag,
					imageURIFlag,
				},
			},
		},
	}

	imgCmd imgCmdFunc = disco.DiscoverImages
	vulCmd vulCmdFunc = disco.DiscoverVulns
	licCmd licCmdFunc = disco.DiscoverLicenses
)

type imgCmdFunc func(ctx context.Context, in *types.ImagesQuery) error
type vulCmdFunc func(ctx context.Context, in *types.VulnsQuery) error
type licCmdFunc func(ctx context.Context, in *types.SimpleQuery) error

func printVersion(c *c.Context) {
	log.Info().Msgf(c.App.Version)
}

func getVersionFromContext(c *c.Context) string {
	val, ok := c.App.Metadata[metaKeyVersion]
	if !ok {
		log.Debug().Msg("version not found in app context")
		return "unknown"
	}
	return val.(string)
}

func runImagesCmd(c *c.Context) error {
	in := &types.ImagesQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindImage
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.URIOnly = c.Bool(outputURIOnlyFlag.Name)

	printVersion(c)
	if err := imgCmd(c.Context, in); err != nil {
		return errors.Wrap(err, "error discovering images")
	}

	return nil
}

func runVulnsCmd(c *c.Context) error {
	in := &types.VulnsQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindVulnerability
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.CVE = c.String(cveFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.CAAPI = c.Bool(caAPIExecFlag.Name)
	in.ImageFile = c.String(imageListPathFlag.Name)
	in.ImageURI = c.String(imageURIFlag.Name)
	in.MinVulnSev = types.ParseMinVulnSeverityOrDefault(c.String(minSeverityFlag.Name))

	printVersion(c)

	if err := in.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	if in.CAAPI {
		log.Info().Msg("note: Container Analysis scans currently are limited to base OS only")
	}

	if err := vulCmd(c.Context, in); err != nil {
		return errors.Wrap(err, "error excuting command")
	}

	return nil
}

func runLicenseCmd(c *c.Context) error {
	in := &types.SimpleQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindLicense
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.ImageFile = c.String(imageListPathFlag.Name)
	in.ImageURI = c.String(imageURIFlag.Name)

	printVersion(c)

	if err := in.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	if err := licCmd(c.Context, in); err != nil {
		return errors.Wrap(err, "error discovering licenses")
	}

	return nil
}
