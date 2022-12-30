package run

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsage(t *testing.T) {
	client = &testAPIClient{}
	ok, err := isAPIEnabled(context.Background(), "799736955886", cloudRunAPI)
	assert.NoError(t, err)
	assert.True(t, ok)
}
