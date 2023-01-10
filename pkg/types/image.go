package types

import (
	"time"
)

func NewItemReport[T any](in *SimpleQuery, items ...*T) *ItemReport[T] {
	return &ItemReport[T]{
		Meta: Meta{
			Kind:    in.Kind.String(),
			Version: in.Version,
			Created: time.Now().UTC().Format(time.RFC3339),
		},
		Items: items,
	}
}

type ItemReport[T any] struct {
	Meta  Meta `json:"meta"`
	Items []*T `json:"items"`
}

type ImageItem struct {
	URI     string            `json:"uri"`
	Context map[string]string `json:"context"`
}
