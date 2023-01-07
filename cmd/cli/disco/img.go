package disco

import (
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	imgCmd = &c.Command{
		Name:    "image",
		Aliases: []string{"img"},
		Usage:   "List deployed container images in specific runtime",
		Action:  runImagesCmd,
		Flags: []c.Flag{
			outputPathFlag,
			outputFormatFlag,
			outputURIOnlyFlag,
			projectIDFlag,
		},
	}
)

func runImagesCmd(c *c.Context) error {
	in := &types.ImagesQuery{}
	in.Version = getVersionFromContext(c)
	in.Kind = types.KindImage
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = types.ParseOutputFormatOrDefault(c.String(outputFormatFlag.Name))
	in.URIOnly = c.Bool(outputURIOnlyFlag.Name)

	printVersion(c)
	if err := disco.DiscoverImages(c.Context, in); err != nil {
		return errors.Wrap(err, "error discovering images")
	}

	return nil
}
