package disco

import (
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	pkgCmd = &c.Command{
		Name:    "package",
		Aliases: []string{"pkg"},
		Usage:   "Scan images for packages",
		Action:  runPackageCmd,
		Flags: []c.Flag{
			imageListPathFlag,
			imageURIFlag,
			outputPathFlag,
			outputFormatFlag,
			projectIDFlag,
			targetFlag,
			pkgNameFlag,
		},
	}
)

func runPackageCmd(c *c.Context) error {
	in := &types.PackageQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindPackage
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.ImageFile = c.String(imageListPathFlag.Name)
	in.ImageURI = c.String(imageURIFlag.Name)
	in.TargetRaw = c.String(targetFlag.Name)
	in.NamePart = c.String(pkgNameFlag.Name)

	printVersion(c)

	if err := in.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	ir, err := validateTarget(&in.SimpleQuery)
	if err != nil {
		return errors.Wrap(err, "invalid target")
	}

	if err := disco.DiscoverPackages(c.Context, in, ir); err != nil {
		return errors.Wrap(err, "error discovering packages")
	}

	return nil
}
