package cli

import (
	"github.com/mchmarny/vctl/pkg/vctl"
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
		Usage:    "output format (json, yaml, raw)",
		Required: false,
	}

	outputDigestOnlyFlag = &c.BoolFlag{
		Name:  "digest-ids",
		Usage: "output only image digests",
		Value: false,
	}

	cveFlag = &c.StringFlag{
		Name:     "cve",
		Aliases:  []string{"e"},
		Usage:    "exposure ID (CVE number, e.g. CVE-2019-19378)",
		Required: false,
	}

	digestFlag = &c.StringFlag{
		Name:     "digest",
		Aliases:  []string{"d"},
		Usage:    "image digest",
		Required: true,
	}

	runCmd = &c.Command{
		Name:  "run",
		Usage: "Cloud Run commands",
		Subcommands: []*c.Command{
			{
				Name:    "images",
				Aliases: []string{"img", "i"},
				Usage:   "list deployed container images",
				Action:  runImagesCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					outputDigestOnlyFlag,
				},
			},
			{
				Name:    "vulnerabilities",
				Aliases: []string{"vul", "v"},
				Usage:   "check for exposures in deployed images (supports specific CVE filter)",
				Action:  runVulnsCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					cveFlag,
				},
			}, {
				Name:    "licenses",
				Aliases: []string{"lic", "l"},
				Usage:   "list licenses used in currently deployed image",
				Action:  runLicenseCmd,
				Flags: []c.Flag{
					digestFlag,
					outputPathFlag,
					outputFormatFlag,
				},
			},
		},
	}
)

func printVersion(c *c.Context) {
	log.Info().Msgf(c.App.Version)
}

func runImagesCmd(c *c.Context) error {
	fmtStr := c.String(outputFormatFlag.Name)
	outFmt, err := vctl.ParseOutputFormat(fmtStr)
	if err != nil {
		return errors.Wrapf(err, "error parsing output format: %s", fmtStr)
	}

	in := &vctl.ImagesQuery{}
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = outFmt
	in.OnlyDigest = c.Bool(outputDigestOnlyFlag.Name)

	printVersion(c)
	if err := vctl.DiscoverImages(c.Context, in); err != nil {
		return errors.Wrap(err, "error discovering images")
	}

	return nil
}

func runVulnsCmd(c *c.Context) error {
	fmtStr := c.String(outputFormatFlag.Name)
	outFmt, err := vctl.ParseOutputFormat(fmtStr)
	if err != nil {
		return errors.Wrapf(err, "error parsing output format: %s", fmtStr)
	}

	in := &vctl.VulnsQuery{}
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.CVE = c.String(cveFlag.Name)
	in.OutputFmt = outFmt

	printVersion(c)

	if in.CVE != "" {
		log.Info().Msg("Note: vulnerability scans currently limited to base OS only")
	}

	if err := vctl.DiscoverVulns(c.Context, in); err != nil {
		return errors.Wrap(err, "error excuting command")
	}

	return nil
}

func runLicenseCmd(c *c.Context) error {
	digest := c.String(digestFlag.Name)
	printVersion(c)
	log.Debug().Msgf("digest: %s", digest)
	return nil
}
