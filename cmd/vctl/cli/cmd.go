package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

func Execute(name, version string) error {
	app := &c.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 name,
		Version:              fmt.Sprintf("%s (%s)", name, version),
		Usage:                "Vulnerability scanner for Cloud Run.",
		UsageText:            "vctl - Vulnerability scanner for Cloud Run.",
		Compiled:             time.Now(),
		Commands: []*c.Command{
			runCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		return errors.Wrap(err, "failed to run app")
	}
	return nil
}
