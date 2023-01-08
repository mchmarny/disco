package types

import (
	"fmt"

	"github.com/pkg/errors"
)

type PackageQuery struct {
	SimpleQuery
	NamePart string
}

func (q *PackageQuery) Validate() error {
	if err := q.SimpleQuery.Validate(); err != nil {
		return errors.Wrap(err, "invalid package query")
	}

	return nil
}

func (q *PackageQuery) String() string {
	return fmt.Sprintf("project: %s, output: %s, format: %s, source: %s, uri: %s, package: %s",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.ImageFile, q.ImageURI, q.NamePart)
}

type PackageReport struct {
	Image    string     `json:"image"`
	Packages []*Package `json:"packages"`
}

type Package struct {
	Package        string `json:"package"`
	PackageVersion string `json:"version"`
	Source         string `json:"source,omitempty"`
	License        string `json:"license,omitempty"`
	Format         string `json:"format"`
	Provider       string `json:"provider"`
	Image          string `json:"-"`
}

func (l *Package) String() string {
	return fmt.Sprintf("package: %s, version: %s, format: %s, provider: %s, source: %s", l.Package, l.PackageVersion, l.Format, l.Provider, l.Source)
}
