package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests.",
		},
		[]string{"method", "status", "authType"})

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"method", "status", "authType"})
)

func init() {
	prometheus.MustRegister(totalRequests, responseStatus)
}

func RecordMetrics(dataAuth, method, status string) {

	defer func() {
		totalRequests.WithLabelValues(method, status, dataAuth).Inc()
		responseStatus.WithLabelValues(method, status, dataAuth).Inc()
	}()

}
