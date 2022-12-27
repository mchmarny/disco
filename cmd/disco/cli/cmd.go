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

func Execute(version, commit string) error {
	app.Version = fmt.Sprintf("%s (commit: %s)", version, commit)
	if err := app.Run(os.Args); err != nil {
		return errors.Wrap(err, "failed to run app")
	}
	return nil
}
