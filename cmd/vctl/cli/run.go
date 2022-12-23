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
		Usage:    "ID of the project to scope discovery to.",
		Required: false,
	}

	cveFlag = &c.StringFlag{
		Name:     "cve",
		Usage:    "ID of the CVE to search for.",
		Required: true,
	}

	runCmd = &c.Command{
		Name:    "run",
		Aliases: []string{"s"},
		Usage:   "Cloud Run commands.",
		Subcommands: []*c.Command{
			{
				Name:    "disco",
				Aliases: []string{"D"},
				Usage:   "Discover CVEs in all the currently deployed images.",
				Action:  runDiscoCmd,
				Flags: []c.Flag{
					projectIDFlag,
				},
			},
			{
				Name:    "find",
				Aliases: []string{"F"},
				Usage:   "Check if any of the currently deployed images have a CVE.",
				Action:  runCVECmd,
				Flags: []c.Flag{
					projectIDFlag,
					cveFlag,
				},
			},
		},
	}
)

func runDiscoCmd(c *c.Context) error {
	projectID := c.String(projectIDFlag.Name)

	log.Debug().Msgf("projectID: %s", projectID)

	if err := vctl.DiscoverVulns(c.Context, projectID); err != nil {
		return errors.Wrapf(err, "error discovering vulnerabilities for project: %s", projectID)
	}

	return nil
}

func runCVECmd(c *c.Context) error {
	projectID := c.String(projectIDFlag.Name)
	log.Debug().Msgf("projectID: %s", projectID)

	return nil
}
