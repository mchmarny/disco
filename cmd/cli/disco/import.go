package disco

import (
	"github.com/mchmarny/disco/pkg/target"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	c "github.com/urfave/cli/v2"
)

var (
	importCmd = &c.Command{
		Name:    "import",
		Aliases: []string{"imp"},
		Usage:   "Import data from open source scanner reports",
		Subcommands: []*c.Command{
			{
				Name:    "vulnerability",
				Aliases: []string{"vul"},
				Usage: `import vulnerabilities from exiting report
e.g. result of trivy image python:3.4-alpine \
	    --format json --security-checks vuln > report.json`,
				Action: runVulnImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					sourceFileFlag,
				},
			},
			{
				Name:    "license",
				Aliases: []string{"lic"},
				Usage: `import licenses from exiting report 
e.g. result of trivy image python:3.4-alpine \
	    --format json --security-checks license > report.json`,
				Action: runLicenseImportCmd,
				Flags: []c.Flag{
					importProjectIDFlag,
					datasetIDFlag,
					sourceFileFlag,
				},
			},
			{
				Name:    "packages",
				Aliases: []string{"pkg"},
				Usage: `import packages from SBOM file
(SPDX 2.3 and CycloneDX 1.3 supported)
e.g. result of syft packages -o spdx-json python:3.4-alpine > report.json`,
				Action: runPackageImportCmd,
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
