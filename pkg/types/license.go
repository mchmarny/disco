package types

import (
	"fmt"

	"github.com/pkg/errors"
)

type LicenseQuery struct {
	SimpleQuery
	TypeFilter string
}

func (q *LicenseQuery) Validate() error {
	if err := q.SimpleQuery.Validate(); err != nil {
		return errors.Wrap(err, "invalid license query")
	}

	return nil
}

func (q *LicenseQuery) String() string {
	return fmt.Sprintf("project: %s, output: %s, format: %s, source: %s, uri: %s, filter: %s",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.ImageFile, q.ImageURI, q.TypeFilter)
}

type LicenseReport struct {
	Image    string     `json:"image"`
	Licenses []*License `json:"licenses"`
}

type License struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Image  string `json:"-"`
}

func (l *License) String() string {
	return fmt.Sprintf("%s: %s", l.Source, l.Name)
}
