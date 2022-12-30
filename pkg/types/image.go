package types

import (
	"fmt"
	"time"
)

func NewItemReport[T any](in *SimpleQuery, items ...*T) *ItemReport[T] {
	itemCount := len(items)
	return &ItemReport[T]{
		Meta: Meta{
			Kind:    in.Kind.String(),
			Version: in.Version,
			Created: time.Now().UTC().Format(time.RFC3339),
			Count:   &itemCount,
		},
		Items: items,
	}
}

type ItemReport[T any] struct {
	Meta  Meta `json:"meta"`
	Items []*T `json:"items"`
}

type ImageItem struct {
	URI     string                 `json:"uri"`
	Context map[string]interface{} `json:"context"`
}

type ImagesQuery struct {
	SimpleQuery
	URIOnly bool
}

func (q *ImagesQuery) String() string {
	return fmt.Sprintf("project:%s, output:%s, format:%s, uri-only:%t",
		q.ProjectID, q.OutputPath, q.OutputFmt, q.URIOnly)
}
