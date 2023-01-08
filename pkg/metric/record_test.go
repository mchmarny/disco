package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPICounter(t *testing.T) {
	ctx := context.TODO()
	c := &APICounter{}
	records := []*Record{}
	err := c.CountAll(ctx, records...)
	assert.NoError(t, err)
}

func TestMetricType(t *testing.T) {
	v := "vulnerability/high"
	c := MakeMetricType(v)
	assert.Equal(t, "custom.googleapis.com/vulnerability/high", c)

	v = "one-two-three"
	c = MakeMetricType(v)
	assert.Equal(t, "custom.googleapis.com/one-two-three", c)
}

func TestMetricRecord(t *testing.T) {
	ctx := context.TODO()
	c := &ConsoleCounter{}
	err := c.Count(ctx, "test", 1, nil)
	assert.NoError(t, err)
	err = c.Count(ctx, "test", 1, map[string]string{"foo": "bar"})
	assert.NoError(t, err)

	records := []*Record{
		{
			MetricType:  "test1",
			Labels:      map[string]string{"foo": "bar"},
			MetricValue: 1,
		},
		{
			MetricType:  "test2",
			Labels:      map[string]string{"foo": "bar"},
			MetricValue: 2,
		},
	}

	err = c.CountAll(ctx, records...)
	assert.NoError(t, err)
}
