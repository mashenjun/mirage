package util

import "github.com/prometheus/client_golang/prometheus"

var (
	ServerRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "requests",
			Help:      "Number of currently processing RESTful requests on server side.",
		},
		[]string{"method", "uri"},
	)

	ServerRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "requests_total",
			Help:      "Total number of RESTful requests on server side.",
		},
		[]string{"method", "uri"},
	)
	ServerResponses4xx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "responses_4xx",
			Help:      "Total number of RESTful 4xx response on server side.",
		},
		[]string{"method", "uri", "status"},
	)
	ServerResponses5xx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "responses_5xx",
			Help:      "Total number of RESTful 5xx response on server side.",
		},
		[]string{"method", "uri", "status"},
	)
	ServerDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "request_duration_seconds",
			Help:      "The RESTful request latencies in milliseconds on server side.",
		},
		[]string{"method", "uri", "status"},
	)

	RPCError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "rpc",
			Name:      "error",
			Help:      "Total number of rpc error",
		},
		[]string{"rpc_name", "status"},
	)
)
