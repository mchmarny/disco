package disco

import (
	"fmt"
	"time"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

const (
	name           = "disco"
	metaKeyVersion = "version"
	metaKeyCommit  = "commit"
	metaKeyDate    = "date"
)

func Execute(version, commit, date string, args []string) error {
	app, err := newApp(version, commit, date)
	if err != nil {
		return err
	}

	if err := app.Run(args); err != nil {
		return errors.Wrap(err, "error running app")
	}
	return nil
}

func newApp(version, commit, date string) (*c.App, error) {
	if version == "" || commit == "" || date == "" {
		return nil, errors.New("version, commit, and date must be set")
	}

	compileTime, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse date")
	}
	dateStr := compileTime.UTC().Format("2006-01-02 15:04 UTC")

	app := &c.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 name,
		Version:              fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, dateStr),
		Usage:                `Discover container images, vulnerabilities, packages, and licenses`,
		Compiled:             compileTime,
		Metadata: map[string]interface{}{
			metaKeyVersion: version,
			metaKeyCommit:  commit,
			metaKeyDate:    date,
		},
		Commands: []*c.Command{
			imgCmd,
			vulnCmd,
			licCmd,
			pkgCmd,
			importCmd,
		},
	}

	return app, nil
}

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

func addOptionalImportFlags(c *c.Context, req *types.ImportRequest) {
	datasetID := c.String(datasetIDFlag.Name)
	if datasetID != "" {
		req.DatasetID = datasetID
	}
}
