package metric

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Counter is the interface for metric counter.
type Counter interface {
	Count(ctx context.Context, metric string, count int64, labels map[string]string) error
	CountAll(ctx context.Context, records ...*Record) error
}

// MakeMetricType creates a metric type string.
func MakeMetricType(v string) string {
	return fmt.Sprintf("custom.googleapis.com/%s", strings.ToLower(v))
}

// Record is the metric record type.
type Record struct {
	MetricType  string
	MetricValue int64
	Labels      map[string]string
}

// NewRecorder creates new recorder instance.
func NewRecorder(client Counter, labels map[string]string) *Recorder {
	r := &Recorder{
		client: client,
		items:  make(map[string]*Record),
		labels: labels,
	}
	if r.labels == nil {
		r.labels = make(map[string]string)
	}
	return r
}

// Recorder is the metric recorder type.
type Recorder struct {
	client Counter
	items  map[string]*Record
	lock   sync.Mutex
	labels map[string]string
}

const flushOnREcorderItems = 100

var recorderFlushError = "error flushing metric records"

// Add adds a metric record to the recorder.
func (r *Recorder) Add(ctx context.Context, name string, labels map[string]string) error {
	if name == "" {
		return errors.New("name required")
	}
	if labels == nil {
		labels = make(map[string]string)
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	for k, v := range r.labels {
		labels[k] = v
	}

	d := &Record{
		MetricType: MakeMetricType(name),
		Labels:     labels,
	}

	k := fmt.Sprintf("%+v", d)
	if _, ok := r.items[k]; !ok {
		r.items[k] = d
	}

	r.items[k].MetricValue++

	if len(r.items) > 0 && len(r.items) >= flushOnREcorderItems {
		if err := r.client.CountAll(ctx, r.slice()...); err != nil {
			return errors.Wrap(err, recorderFlushError)
		}
		r.items = make(map[string]*Record)
	}
	return nil
}

// Flush flushes the recorder.
func (r *Recorder) Flush(ctx context.Context) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.client.CountAll(ctx, r.slice()...); err != nil {
		return errors.Wrap(err, recorderFlushError)
	}

	r.items = make(map[string]*Record)

	return nil
}

func (r *Recorder) slice() []*Record {
	v := make([]*Record, 0, len(r.items))
	for _, value := range r.items {
		v = append(v, value)
	}
	return v
}

func (r *Recorder) Size() int {
	return len(r.items)
}
