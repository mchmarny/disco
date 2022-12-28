package types

import (
	"fmt"
	"time"
)

func NewItemReport(in *SimpleQuery) *ItemReport {
	return &ItemReport{
		Meta: Meta{
			Kind:    in.Kind.String(),
			Version: in.Version,
			Created: time.Now().UTC().Format(time.RFC3339),
		},
		Items: make([]interface{}, 0),
	}
}

type ItemReport struct {
	Meta  Meta          `json:"meta"`
	Items []interface{} `json:"items"`
}

type ImageItem struct {
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
