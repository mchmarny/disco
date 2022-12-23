package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

func Run(name, version string) error {
	app := &c.App{
		Name:     name,
		Version:  fmt.Sprintf("%s - %s", name, version),
		Compiled: time.Now(),
		Usage:    name,
		Commands: []*c.Command{
			runCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		return errors.Wrap(err, "failed to run app")
	}
	return nil
}
