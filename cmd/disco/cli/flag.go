package cli

import c "github.com/urfave/cli/v2"

var (
	projectIDFlag = &c.StringFlag{
		Name:    "project",
		Aliases: []string{"p"},
		Usage:   "project ID",
	}

	targetFlag = &c.StringFlag{
		Name:  "target",
		Usage: "Data store to where the results should be saved (e.g. bq://my-project)",
	}

	outputPathFlag = &c.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "path where to save the output",
	}

	outputFormatFlag = &c.StringFlag{
		Name:  "format",
		Usage: "output format (json, yaml)",
	}

	cveFlag = &c.StringFlag{
		Name:  "cve",
		Usage: "exposure ID (CVE number, e.g. CVE-2019-19378)",
	}

	pkgNameFlag = &c.StringFlag{
		Name:  "name",
		Usage: "package name contains (e.g. libgcc, gobinary, express, etc.)",
	}

	licenseTypeFlag = &c.StringFlag{
		Name:  "type",
		Usage: "license type (supports prefix: e.g. apache, bsd, mit, etc.)",
	}

	minSeverityFlag = &c.StringFlag{
		Name:    "min-severity",
		Aliases: []string{"min-sev"},
		Usage:   "minimum severity to report (e.g. low, medium, high, critical)",
	}

	imageListPathFlag = &c.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "image list input file to serve as source (instead of discovery)",
	}

	imageURIFlag = &c.StringFlag{
		Name:    "image",
		Aliases: []string{"img"},
		Usage:   "specific image URI to scan",
	}
)
