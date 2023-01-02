package metric

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	metricPageSize = 150
)

// NewAPICounter creates new API counter instance.
func NewAPICounter(project string) (Counter, error) {
	return &APICounter{
		projectID: project,
		labels: map[string]string{
			"project_id": project,
		},
	}, nil
}

// APICounter is the API counter type.
type APICounter struct {
	projectID string
	labels    map[string]string
}

// Count records the metric to the monitoring API.
func (r *APICounter) Count(ctx context.Context, metricType string, count int64, labels map[string]string) error {
	items := make(map[string]int64)
	items[metricType] = count

	rec := &Record{
		MetricType:  metricType,
		MetricValue: count,
		Labels:      labels,
	}

	if err := r.CountAll(ctx, rec); err != nil {
		return err
	}
	return nil
}

// CountAll records multiple metrics to the monitoring API.
func (r *APICounter) CountAll(ctx context.Context, records ...*Record) error {
	log.Debug().Msgf("recording %d metrics", len(records))
	if len(records) < 1 {
		log.Debug().Msg("no metrics to record")
		return nil
	}
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating client")
	}
	defer c.Close()
	now := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
	}

	list := make([]*monitoringpb.TimeSeries, 0)
	for _, d := range records {
		if d.Labels == nil {
			d.Labels = map[string]string{}
		}

		// HACK: prevents time series from being overwritten \
		// for timespan which leads to errors on write.
		d.Labels["nanos"] = fmt.Sprintf("e-%d", now.AsTime().UnixMilli())

		s := &monitoringpb.TimeSeries{
			Resource: &monitoredres.MonitoredResource{
				Type:   "global",
				Labels: r.labels,
			},
			Metric: &metricpb.Metric{
				Type:   d.MetricType,
				Labels: d.Labels,
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: d.MetricValue,
					},
				},
			}},
		}
		list = append(list, s)
	}

	log.Debug().Msgf("creating %d metrics...", len(list))
	pager, err := types.GetPager(list, metricPageSize)
	if err != nil {
		return errors.Wrapf(err, "getting pager of size %d from %d items",
			metricPageSize, len(list))
	}

	for {
		items := pager.Next()
		if len(items) < 1 {
			break
		}

		if err = c.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
			Name:       "projects/" + r.projectID,
			TimeSeries: items,
		}); err != nil {
			return errors.Wrap(err, "writing time series request")
		}
	}

	return nil
}
