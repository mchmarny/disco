package disco

import (
	"github.com/mchmarny/disco/pkg/target"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	importProjectIDFlag = &c.StringFlag{
		Name:     "project",
		Aliases:  []string{"p"},
		Usage:    "ID of the project to import into",
		Required: true,
	}

	datasetIDFlag = &c.StringFlag{
		Name:     "dataset",
		Aliases:  []string{"d"},
		Usage:    "ID of the dataset to import into",
		Required: false,
	}

	tableIDFlag = &c.StringFlag{
		Name:     "table",
		Aliases:  []string{"t"},
		Usage:    "ID of the table to import into",
		Required: false,
	}

	sourceFileFlag = &c.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "Path to the source file to import",
		Required: true,
	}

	importCmd = &c.Command{
		Name:    "import",
		Aliases: []string{"imp"},
		Usage:   "Data import commands (currently supports only BigQuery)",
		Subcommands: []*c.Command{
			{
				Name:    "vulnerabilities",
				Aliases: []string{"vul", "v"},
				Usage:   "Import vulnerabilities from output report file",
				Action:  runVulnImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					tableIDFlag,
					sourceFileFlag,
				},
			},
			{
				Name:    "licenses",
				Aliases: []string{"lic", "l"},
				Usage:   "Import licenses from output report file",
				Action:  runLicenseImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					tableIDFlag,
					sourceFileFlag,
				},
			},
		},
	}
)

func runVulnImportCmd(c *c.Context) error {
	req := types.NewVulnerabilityImportRequest(c.String(projectIDFlag.Name), c.String(sourceFileFlag.Name))
	datasetID := c.String(datasetIDFlag.Name)
	if datasetID != "" {
		req.DatasetID = datasetID
	}
	tableID := c.String(tableIDFlag.Name)
	if tableID != "" {
		req.TableID = tableID
	}

	printVersion(c)

	if err := target.SetupConfigurer(c.Context, req); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	if err := target.VulnerabilityImporter(c.Context, req); err != nil {
		return errors.Wrap(err, "error importing vulnerabilities")
	}

	return nil
}

func runLicenseImportCmd(c *c.Context) error {
	req := types.NewLicenseImportRequest(c.String(projectIDFlag.Name), c.String(sourceFileFlag.Name))
	datasetID := c.String(datasetIDFlag.Name)
	if datasetID != "" {
		req.DatasetID = datasetID
	}
	tableID := c.String(tableIDFlag.Name)
	if tableID != "" {
		req.TableID = tableID
	}

	printVersion(c)

	if err := target.SetupConfigurer(c.Context, req); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	if err := target.LicenseImporter(c.Context, req); err != nil {
		return errors.Wrap(err, "error importing licenses")
	}

	return nil
}
