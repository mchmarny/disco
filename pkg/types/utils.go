package types

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Filter is a filter function.
type ItemFilter func(v string) bool

// HashStr returns hash of the given value as string.
// If it can't, it returns string representation of the value.
func HashStr(v any) string {
	b := Hash(v)
	if b == nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

// Hash returns hash of the given value.
// If it can't, it logs error and returns nil.
func Hash(v any) []byte {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(v); err != nil {
		log.Error().Err(err).Msg("error encoding")
		return nil
	}
	return b.Bytes()
}

func ToHash(v ...any) string {
	s := fmt.Sprintf("%v", v)
	b := Hash(v)
	if b == nil {
		return s
	}
	return string(b)
}
