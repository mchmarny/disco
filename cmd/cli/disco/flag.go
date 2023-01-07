package disco

import c "github.com/urfave/cli/v2"

var (
	projectIDFlag = &c.StringFlag{
		Name:     "project",
		Aliases:  []string{"p"},
		Usage:    "project ID",
		Required: false,
	}

	importProjectIDFlag = &c.StringFlag{
		Name:     "project",
		Aliases:  []string{"p"},
		Usage:    "ID of the project to import into",
		Required: true,
	}

	outputPathFlag = &c.StringFlag{
		Name:     "output",
		Aliases:  []string{"o"},
		Usage:    "path where to save the output",
		Required: false,
	}

	outputFormatFlag = &c.StringFlag{
		Name:     "format",
		Usage:    "output format (json, yaml)",
		Required: false,
	}

	sbomFormatFlag = &c.StringFlag{
		Name:     "format",
		Usage:    "SBOM Format of the source file to import (spdx or cyclonedx)",
		Required: true,
	}

	reportFormatFlag = &c.StringFlag{
		Name:     "format",
		Usage:    "Report Format of the source file to import (disco or trivy)",
		Required: true,
	}

	outputURIOnlyFlag = &c.BoolFlag{
		Name:  "uri",
		Usage: "output only image URI",
		Value: false,
	}

	cveFlag = &c.StringFlag{
		Name:     "cve",
		Usage:    "exposure ID (CVE number, e.g. CVE-2019-19378)",
		Required: false,
	}

	licenseTypeFlag = &c.StringFlag{
		Name:     "type",
		Aliases:  []string{"t"},
		Usage:    "license type (supports prefix: e.g. apache, bsd, mit, etc.)",
		Required: false,
	}

	minSeverityFlag = &c.StringFlag{
		Name:     "min-severity",
		Aliases:  []string{"min-sev"},
		Usage:    "minimum severity to report (e.g. low, medium, high, critical)",
		Required: false,
	}

	imageListPathFlag = &c.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "image list input file to serve as source (instead of discovery)",
		Required: false,
	}

	imageURIFlag = &c.StringFlag{
		Name:     "image",
		Aliases:  []string{"img"},
		Usage:    "specific image URI to scan",
		Required: false,
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
)
