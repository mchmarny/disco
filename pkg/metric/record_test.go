package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricType(t *testing.T) {
	v := "vulnerability/high"
	c := MakeMetricType(v)
	assert.Equal(t, "custom.googleapis.com/vulnerability/high", c)

	v = "one-two-three"
	c = MakeMetricType(v)
	assert.Equal(t, "custom.googleapis.com/one-two-three", c)
}
