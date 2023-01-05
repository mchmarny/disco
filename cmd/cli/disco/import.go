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

	sourceFileFlag = &c.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "Path to the source file to import",
		Required: true,
	}

	sbomFormatFlag = &c.StringFlag{
		Name:     "format",
		Usage:    "SBOM Format of the source file to import (spdx or cyclonedx)",
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
				Usage:   "Import vulnerabilities from trivy report (e.g. trivy image python:3.4-alpine --format json --security-checks vuln > report.json)",
				Action:  runVulnImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					sourceFileFlag,
				},
			},
			{
				Name:    "licenses",
				Aliases: []string{"lic", "l"},
				Usage:   "Import licenses from trivy report (e.g. trivy image python:3.4-alpine --format json --security-checks license > report.json)",
				Action:  runLicenseImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					sourceFileFlag,
				},
			},
			{
				Name:    "packages",
				Aliases: []string{"pkg", "p"},
				Usage:   "Import packages from SBOM file (SPDX 2.3 and CycloneDX 1.3 supported, e.g. syft packages -o spdx-json python:3.4-alpine > report.json)",
				Action:  runPackageImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					sourceFileFlag,
					sbomFormatFlag,
				},
			},
		},
	}
)

func runVulnImportCmd(c *c.Context) error {
	req := types.NewVulnerabilityImportRequest(c.String(projectIDFlag.Name),
		c.String(sourceFileFlag.Name))
	addOptionalImportFlags(c, req)
	printVersion(c)

	if err := target.VulnerabilityImporter(c.Context, req); err != nil {
		return errors.Wrap(err, "error importing vulnerabilities")
	}

	return nil
}

func runLicenseImportCmd(c *c.Context) error {
	req := types.NewLicenseImportRequest(c.String(projectIDFlag.Name),
		c.String(sourceFileFlag.Name))
	addOptionalImportFlags(c, req)
	printVersion(c)

	if err := target.LicenseImporter(c.Context, req); err != nil {
		return errors.Wrap(err, "error importing licenses")
	}

	return nil
}

func runPackageImportCmd(c *c.Context) error {
	req := types.NewPackageImportRequest(c.String(projectIDFlag.Name),
		c.String(sourceFileFlag.Name), c.String(sbomFormatFlag.Name))
	addOptionalImportFlags(c, req)
	printVersion(c)

	if err := target.PackageImporter(c.Context, req); err != nil {
		return errors.Wrap(err, "error importing licenses")
	}

	return nil
}

func addOptionalImportFlags(c *c.Context, req *types.ImportRequest) {
	datasetID := c.String(datasetIDFlag.Name)
	if datasetID != "" {
		req.DatasetID = datasetID
	}
}
