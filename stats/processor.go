package stats

import (
	"bytes"
	"github.com/criteo/graphite-writer-stats/prometheus"
	"go.uber.org/zap"
	"strconv"
)

type Stats struct {
	Logger         *zap.Logger
	MetricMetadata MetricMetadata
}

func (stats *Stats) process(metricPath string) {
	metric := stats.getMetric(metricPath)
	stats.Logger.Debug("metrics", zap.Any("metric", metric))
	prometheus.IncMetricPathCounter(metric.ExtractedMetric, metric.ApplicationName, string(metric.ApplicationType), strconv.FormatBool(metric.IsPrometheus))
}

func (stats *Stats) Process(dataPoint []byte) bool {
	metricPath, succeed := extractMetricPath(dataPoint)
	if succeed {
		stats.process(metricPath)
	} else {
		stats.Logger.Error("fail to convert datapoint to metricpath", zap.ByteString("datapoint", dataPoint))
		prometheus.IncDataPointToMetricErrorCounter()
	}
	return succeed
}

func extractMetricPath(metric []byte) (string, bool) {
	indexSpace := bytes.IndexByte(metric, ' ')
	return string(metric[:indexSpace+1]), indexSpace != -1
}
