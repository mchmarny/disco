package cli

import (
	"fmt"
	"os"
	"time"

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

	return app.Run(os.Args)
}
