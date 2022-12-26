package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsage(t *testing.T) {
	APIClient = &TestAPIClient{}
	ok, err := IsAPIEnabled(context.Background(), "799736955886", CloudRunAPI)
	assert.NoError(t, err)
	assert.True(t, ok)
}
