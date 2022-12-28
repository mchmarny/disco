package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

const (
	name = "disco"
)

var (
	app = &c.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 name,
		Usage:                `Discover container images, vulnerabilities, and licenses in currently deployed across your runtimes`,
		Compiled:             time.Now(),
		Commands: []*c.Command{
			runCmd,
		},
	}
)

func Execute(version, commit, date string) error {
	d, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err != nil {
		return errors.Wrap(err, "failed to parse date")
	}
	date = d.UTC().Format("2006-01-02 15:04 UTC")
	app.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	if err := app.Run(os.Args); err != nil {
		return errors.Wrap(err, "failed to run app")
	}
	return nil
}
