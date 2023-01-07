package disco

import (
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	vulnCmd = &c.Command{
		Name:    "vulnerability",
		Aliases: []string{"vul"},
		Usage:   "Check for OS and package level exposures in images (supports CVE filter)",
		Action:  runVulnsCmd,
		Flags: []c.Flag{
			imageListPathFlag,
			imageURIFlag,
			outputPathFlag,
			outputFormatFlag,
			projectIDFlag,
			minSeverityFlag,
			cveFlag,
		},
	}
)

func runVulnsCmd(c *c.Context) error {
	in := &types.VulnsQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindVulnerability
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.CVE = c.String(cveFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.ImageFile = c.String(imageListPathFlag.Name)
	in.ImageURI = c.String(imageURIFlag.Name)
	in.MinVulnSev = types.ParseMinVulnSeverityOrDefault(c.String(minSeverityFlag.Name))

	printVersion(c)

	if err := in.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	if err := disco.DiscoverVulns(c.Context, in); err != nil {
		return errors.Wrap(err, "error excuting command")
	}

	return nil
}
