package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber-go/zap"
	"net/http"
	"strconv"
)

var (
	metricPathCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "metrics_path_total",
		Help: "The total number of metrics paths events",
	}, []string{"metric_path", "application", "application_type"})
	dataPointTometricErrorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "metrics_error_total",
		Help: "The total number of bad parsed metrics paths",
	})
)

func SetupPrometheusHTTPServer(logger *zap.Logger,port int, endpoint string) {
	go func() {
		portBinding := ":" + strconv.Itoa(port)
		http.Handle(endpoint, promhttp.Handler())
		err:=http.ListenAndServe(portBinding, nil)
		if err!= nil {
			logger.Panic("could not set up prometheus endpoint",zap.Error(err))
		}
	}()
}
func IncDataPointToMetricErrorCounter() {
	dataPointTometricErrorCount.Inc()
}

func IncMetricPathCounter(extractedMetric string, applicationName string, applicationType string) {
	metricPathCount.WithLabelValues(extractedMetric, applicationName, applicationType).Inc()
}
