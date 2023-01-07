package disco

import c "github.com/urfave/cli/v2"

var (
	projectIDFlag = &c.StringFlag{
		Name:    "project",
		Aliases: []string{"p"},
		Usage:   "project ID",
	}

	targetFlag = &c.StringFlag{
		Name:  "target",
		Usage: "Target data store to save the results (e.g. bq://project.dataset.table)",
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

	outputURIOnlyFlag = &c.BoolFlag{
		Name:  "uri",
		Usage: "output only image URI",
	}

	cveFlag = &c.StringFlag{
		Name:  "cve",
		Usage: "exposure ID (CVE number, e.g. CVE-2019-19378)",
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
