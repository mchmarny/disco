package cli

import (
	"fmt"
	"time"

	"github.com/mchmarny/disco/pkg/target"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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
		log.Warn().Err(err).Msg("compile time not set")
		compileTime = time.Now()
	}
	dateStr := compileTime.UTC().Format("2006-01-02 15:04 UTC")

	app := &c.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 name,
		Version:              fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, dateStr),
		Usage:                `Discover container images, vulnerabilities, packages, and licenses`,
		Compiled:             compileTime,
		Flags: []c.Flag{
			&c.BoolFlag{
				Name:  "debug",
				Usage: "verbose output",
				Action: func(c *c.Context, debug bool) error {
					if debug {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}
					return nil
				},
			},
			&c.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "suppress output unless error",
				Action: func(c *c.Context, quiet bool) error {
					if quiet {
						c.App.Metadata["quiet"] = true
					}
					return nil
				},
			},
		},
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
		},
	}

	return app, nil
}

func isQuiet(c *c.Context) bool {
	_, ok := c.App.Metadata["quiet"]
	if ok {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	return ok
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

func validateTarget(req *types.SimpleQuery) (*types.ImportRequest, error) {
	if req.TargetRaw != "" {
		ir, err := target.ParseImportRequest(req)
		if err != nil {
			return nil, errors.Wrap(err, "invalid target definition")
		}
		return ir, nil
	}
	return nil, nil
}
