package types

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

type PackageQuery struct {
	SimpleQuery
}

func (q *PackageQuery) Validate() error {
	if err := q.SimpleQuery.Validate(); err != nil {
		return errors.Wrap(err, "invalid package query")
	}

	return nil
}

func (q *PackageQuery) String() string {
	return fmt.Sprintf("project: %s, output: %s, format: %s, source: %s, uri: %s",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.ImageFile, q.ImageURI)
}

type PackageReport struct {
	Image    string     `json:"image"`
	Packages []*Package `json:"packages"`
}

type Package struct {
	Package        string `json:"package"`
	PackageVersion string `json:"version"`
	Format         string `json:"format"`
	Provider       string `json:"provider"`
	Originator     string `json:"originator"`
	Source         string `json:"source"`
	License        string `json:"license"`
	Image          string `json:"-"`
}

func (l *Package) String() string {
	return fmt.Sprintf("package: %s, version: %s, format: %s, provider: %s, originator: %s, source: %s, license: %s", l.Package, l.PackageVersion, l.Format, l.Provider, l.Originator, l.Source, l.License)
}

const spdxToolKey = "Tool"

func SPDXCreatorInfo(in *v2_3.CreationInfo) string {
	if in == nil {
		return ""
	}

	var sb strings.Builder

	for _, c := range in.Creators {
		if c.CreatorType == spdxToolKey {
			return c.Creator
		} else {
			sb.WriteString(c.Creator)
			sb.WriteString(" ")
		}
	}

	return strings.TrimSpace(sb.String())
}
