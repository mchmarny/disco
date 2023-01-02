package metric

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	safeLabelNameExp = regexp.MustCompile(`[^a-zA-Z-]+`)
	safeTypeNameExp  = regexp.MustCompile(`[^a-zA-Z0-9-/]+`)
)

// Counter is the interface for metric counter.
type Counter interface {
	Count(ctx context.Context, metric string, count int64, labels map[string]string) error
	CountAll(ctx context.Context, records ...*Record) error
}

// MakeMetricType creates a metric type string.
func MakeMetricType(v string) string {
	return fmt.Sprintf("custom.googleapis.com/%s", safeTypeNameExp.ReplaceAllString(strings.ToLower(v), ""))
}

// MakeMetricLabelSafe makes a metric label safe.
func MakeMetricLabelSafe(v string) string {
	return safeLabelNameExp.ReplaceAllString(v, "")
}

// Record is the metric record type.
type Record struct {
	MetricType  string
	MetricValue int64
	Labels      map[string]string
}
