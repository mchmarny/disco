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
		Usage:    "ID of the project.",
		Required: false,
	}

	cveFlag = &c.StringFlag{
		Name:     "cve",
		Aliases:  []string{"e"},
		Usage:    "ID of the exposure.",
		Required: false,
	}

	digestFlag = &c.StringFlag{
		Name:     "digest",
		Aliases:  []string{"d"},
		Usage:    "Image digest.",
		Required: true,
	}

	runCmd = &c.Command{
		Name:  "run",
		Usage: "Cloud Run commands.",
		Subcommands: []*c.Command{
			{
				Name:    "images",
				Aliases: []string{"img", "i"},
				Usage:   "List container images used in Cloud Run.",
				Action:  runImagesCmd,
				Flags: []c.Flag{
					projectIDFlag,
				},
			},
			{
				Name:    "vulnerabilities",
				Aliases: []string{"vul", "v"},
				Usage:   "Check if any of the currently deployed images have that specific CVE.",
				Action:  runVulnsCmd,
				Flags: []c.Flag{
					projectIDFlag,
					cveFlag,
				},
			}, {
				Name:    "licenses",
				Aliases: []string{"lic", "l"},
				Usage:   "List a unique list of licenses used in this image.",
				Action:  runLicenseCmd,
				Flags: []c.Flag{
					digestFlag,
				},
			},
		},
	}
)

func printVersion(c *c.Context) {
	log.Info().Msgf(c.App.Version)
}

func runImagesCmd(c *c.Context) error {
	projectID := c.String(projectIDFlag.Name)

	printVersion(c)
	if err := vctl.DiscoverImages(c.Context, projectID); err != nil {
		return errors.Wrap(err, "Error discovering images.")
	}

	return nil
}

func runVulnsCmd(c *c.Context) error {
	projectID := c.String(projectIDFlag.Name)
	cveID := c.String(cveFlag.Name)

	printVersion(c)
	if err := vctl.DiscoverVulns(c.Context, projectID, cveID); err != nil {
		return errors.Wrap(err, "Error discovering vulnerabilities.")
	}

	return nil
}

func runLicenseCmd(c *c.Context) error {
	digest := c.String(digestFlag.Name)
	printVersion(c)
	log.Debug().Msgf("Digest: %s.", digest)
	return nil
}
