package metric

import (
	"context"
	"fmt"
)

// Counter is the interface for metric implementations.
type ConsoleCounter struct {
}

// Count records the metric to the console.
func (r *ConsoleCounter) Count(ctx context.Context, metric string, count int64, labels map[string]string) error {
	fmt.Printf("console counter - %s:%d", metric, count)
	return nil
}

// CountAll records multiple metrics to the console.
func (r *ConsoleCounter) CountAll(ctx context.Context, records ...*Record) error {
	fmt.Println("console counter:")
	for _, d := range records {
		fmt.Printf("%s:%d", d.MetricType, d.MetricValue)
	}
	return nil
}
