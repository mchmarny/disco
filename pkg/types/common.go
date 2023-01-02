package types

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// Filter is a filter function.
type ItemFilter func(v interface{}) bool

// UnmarshalFromFile unmarshals JSON from file.
func UnmarshalFromFile(path string, report interface{}) error {
	if path == "" {
		return errors.New("report path required")
	}

	reportContent, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "error reading report file: %s", path)
	}

	if err := json.Unmarshal(reportContent, report); err != nil {
		return errors.Wrapf(err, "error parsing report file: %s", path)
	}
	return nil
}
