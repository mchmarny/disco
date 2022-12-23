package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

const (
	limitDefault = 100
)

var (
	stringFlag = &c.StringFlag{
		Name:     "username",
		Usage:    "Twitter username",
		Required: true,
	}

	intFlag = &c.IntFlag{
		Name:  "limit",
		Usage: "Limit the number of results",
		Value: limitDefault,
	}

	simpleCmd = &c.Command{
		Name:    "simple",
		Aliases: []string{"s"},
		Usage:   "Simple CLI command.",
		Subcommands: []*c.Command{
			{
				Name:    "one",
				Aliases: []string{"o"},
				Usage:   "First command.",
				Action:  cmdImplementation,
				Flags: []c.Flag{
					stringFlag,
				},
			},
			{
				Name:    "two",
				Aliases: []string{"t"},
				Action:  cmdImplementation,
				Usage:   "Second command.",
				Flags: []c.Flag{
					intFlag,
				},
			},
		},
	}
)

func Run(name, version string) error {
	app := &c.App{
		Name:     name,
		Version:  fmt.Sprintf("%s - %s", name, version),
		Compiled: time.Now(),
		Usage:    name,
		Commands: []*c.Command{
			simpleCmd,
		},
	}

	return app.Run(os.Args)
}

func cmdImplementation(c *c.Context) error {
	val := c.String(stringFlag.Name)

	limit := c.Int(intFlag.Name)
	if limit == 0 {
		limit = limitDefault
	}

	log.Debug().Msgf("value: %s and %d", val, limit)

	list := []string{"one", "two", "three"}

	if err := json.NewEncoder(os.Stdout).Encode(list); err != nil {
		return errors.Wrapf(err, "error encoding list: %v", list)
	}

	return nil
}
