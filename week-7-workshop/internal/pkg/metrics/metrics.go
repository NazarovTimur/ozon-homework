package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "handler_request_total_counter",
			Help:      "Total amount of request by handler",
		},
		[]string{"handler"},
	)

	handlerHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "handler_duration_histogram",
			Help:      "Total duration of handler processing by request",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"handler"},
	)
)

func IncRequestCounter(handler string) {
	requestCounter.WithLabelValues(handler).Inc()
}

func StoreHandlerDuration(handler string, since time.Duration) {
	handlerHistogram.WithLabelValues(handler).Observe(since.Seconds())
}
