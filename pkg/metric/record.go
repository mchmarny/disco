package metric

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	safeLabelKeyExp   = regexp.MustCompile(`[^a-zA-Z_]+`)
	safeLabelValueExp = regexp.MustCompile(`[^a-zA-Z0-9-_]+`)
	safeTypeNameExp   = regexp.MustCompile(`[^a-zA-Z0-9-/]+`)
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

// MakeMetricLabelKeySafe makes a metric label key safe.
func MakeMetricLabelKeySafe(v string) string {
	return safeLabelKeyExp.ReplaceAllString(v, "_")
}

// MakeMetricLabelValueSafe makes a metric label value safe.
func MakeMetricLabelValueSafe(v string) string {
	return safeLabelValueExp.ReplaceAllString(v, "_")
}

// Record is the metric record type.
type Record struct {
	MetricType  string
	MetricValue int64
	Labels      map[string]string
}
