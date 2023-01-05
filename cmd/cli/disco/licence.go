package disco

import (
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	licCmd = &c.Command{
		Name:    "license",
		Aliases: []string{"lic"},
		Usage:   "Scan images for licenses (supports type filter)",
		Action:  runLicenseCmd,
		Flags: []c.Flag{
			imageListPathFlag,
			imageURIFlag,
			runtimeFlag,
			outputPathFlag,
			outputFormatFlag,
			projectIDFlag,
			licenseTypeFlag,
		},
	}
)

func runLicenseCmd(c *c.Context) error {
	in := &types.LicenseQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindLicense
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.ImageFile = c.String(imageListPathFlag.Name)
	in.ImageURI = c.String(imageURIFlag.Name)
	in.TypeFilter = c.String(licenseTypeFlag.Name)

	printVersion(c)

	if err := in.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	if err := disco.DiscoverLicenses(c.Context, in); err != nil {
		return errors.Wrap(err, "error discovering licenses")
	}

	return nil
}
