package types

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Filter is a filter function.
type ItemFilter func(v interface{}) bool

// ToKey returns base64 encoded string from the given parts.
func ToKey(parts ...string) string {
	b := bytes.Buffer{}
	for _, p := range parts {
		b.WriteString(p)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// Hash returns hash of the given value.
// If it can't, it logs error and returns nil.
func Hash(v interface{}) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer func() {
		if err := encoder.Close(); err != nil {
			log.Error().Err(err).Msgf("error encoding: %v", v)
		}
	}()
	if err := gob.NewEncoder(encoder).Encode(v); err != nil {
		log.Error().Err(err).Msg("error encoding")
		return fmt.Sprintf("%v", v)
	}
	return buf.String()
}
