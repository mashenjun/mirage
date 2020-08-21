package prom

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(serverDuration)
	prometheus.MustRegister(serverRequests)
	prometheus.MustRegister(serverRequestsTotal)
	prometheus.MustRegister(serverResponses4xx)
	prometheus.MustRegister(serverResponses5xx)
}

var (
	serverRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "requests",
			Help:      "Number of currently processing RESTful requests on server side.",
		},
		[]string{"method", "uri"},
	)

	serverRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "requests_total",
			Help:      "Total number of RESTful requests on server side.",
		},
		[]string{"method", "uri"},
	)
	serverResponses4xx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "responses_4xx",
			Help:      "Total number of RESTful 4xx response on server side.",
		},
		[]string{"method", "uri", "status"},
	)
	serverResponses5xx = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "responses_5xx",
			Help:      "Total number of RESTful 5xx response on server side.",
		},
		[]string{"method", "uri", "status"},
	)
	serverDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "request_duration_seconds",
			Help:      "The RESTful request latencies in milliseconds on server side.",
		},
		[]string{"method", "uri", "status"},
	)
)

func MetricsMiddleware(ctx *gin.Context) {
	method := ctx.Request.Method
	uri := ctx.Request.URL.Path

	defer func(begin time.Time) {
		status := fmt.Sprintf("%d", ctx.Writer.Status())

		serverRequests.WithLabelValues(method, uri).Dec()

		timeElapsed := float64(time.Since(begin)) / float64(time.Millisecond)
		serverDuration.WithLabelValues(method, uri, status).Observe(timeElapsed)

		if ctx.Writer.Status() >= http.StatusBadRequest {
			serverResponses4xx.WithLabelValues(method, uri, status).Inc()
		}

		if ctx.Writer.Status() >= http.StatusInternalServerError {
			serverResponses5xx.WithLabelValues(method, uri, status).Inc()
		}

	}(time.Now())

	serverRequests.WithLabelValues(method, uri).Inc()
	serverRequestsTotal.WithLabelValues(method, uri).Inc()

	ctx.Next()
}

const (
	// DefaultMetricsPath url path of metrics
	DefaultMetricsPath = "/metrics"
)

func getPath(pathOptions ...string) string {
	path := DefaultMetricsPath
	if len(pathOptions) > 0 {
		path = pathOptions[0]
	}
	return path
}

func Register(r *gin.Engine, pathOptions ...string) {
	path := getPath(pathOptions...)
	r.GET(path, prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
