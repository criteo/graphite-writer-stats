package stats

import (
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestProcess(t *testing.T) {
	logger := zaptest.NewLogger(t)
	aggregatedRule := Rule{"aggregated", []string{"criteo", "aggregated"}, 2}
	aggloRule := Rule{"agglo", []string{"criteo", "agglo"}, 2}
	aggregatedAllRule := Rule{"aggregated-all", []string{"criteo", "aggregated-all"}, 2}
	legacyHostingRule := Rule{"legacy-hosting", []string{"prometheus", "hosting"}, 1}
	startWithCriteoRule := Rule{"start-by-criteo", []string{"criteo"}, 1}
	startbyAppRule := Rule{"start-by-app", []string{}, 0}
	rulesTab := []Rule{aggregatedRule, aggloRule, aggregatedAllRule, legacyHostingRule, startWithCriteoRule, startbyAppRule}

	stats := Stats{Logger: logger, MetricMetadata: MetricMetadata{
		Rules:        Rules{Rules: rulesTab},
		ComponentsNb: 3,
	}}

	metricDatapoint := []byte("criteo.aggregated.cas.value 3.2 1498887")
	success := stats.Process(metricDatapoint)
	if !success {
		t.Errorf("failed to process '%v'", string(metricDatapoint))
	}
	metricDatapoint = nil
	success = stats.Process(metricDatapoint)
	if success {
		t.Errorf("process a nil metric should return false")
	}
	metricDatapoint = []byte("criteo.498887")
	success = stats.Process(metricDatapoint)
	if success {
		t.Errorf("process a mal formed metric should return false '%v'", string(metricDatapoint))
	}
}
