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
			outputPathFlag,
			outputFormatFlag,
			projectIDFlag,
			licenseTypeFlag,
			targetFlag,
		},
	}
)

func runLicenseCmd(c *c.Context) error {
	in := &types.LicenseQuery{}
	in.Quiet = isQuiet(c)
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindLicense
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.ImageFile = c.String(imageListPathFlag.Name)
	in.ImageURI = c.String(imageURIFlag.Name)
	in.TypeFilter = c.String(licenseTypeFlag.Name)
	in.TargetRaw = c.String(targetFlag.Name)

	printVersion(c)

	if err := in.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	ir, err := validateTarget(&in.SimpleQuery)
	if err != nil {
		return errors.Wrap(err, "invalid target")
	}

	if err := disco.DiscoverLicenses(c.Context, in, ir); err != nil {
		return errors.Wrap(err, "error discovering licenses")
	}

	return nil
}
