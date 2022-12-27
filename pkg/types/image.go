package types

import (
	"fmt"
)

type ImageReport struct {
	Image    string `json:"image"`
	Service  string `json:"service"`
	Project  string `json:"project"`
	Location string `json:"location"`
}

type ImagesQuery struct {
	SimpleQuery
	URIOnly bool
}

func (q *ImagesQuery) String() string {
	return fmt.Sprintf("project:%s, output:%s, format:%s, uri-only:%t",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.URIOnly)
}
