package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIVersion(t *testing.T) {
	args := []string{"--version"}
	err := app.Run(args)
	assert.NoError(t, err)
}
