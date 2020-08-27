package prom

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/mashenjun/mirage/util"
)

func init() {
	prometheus.MustRegister(util.ServerDuration)
	prometheus.MustRegister(util.ServerRequests)
	prometheus.MustRegister(util.ServerRequestsTotal)
	prometheus.MustRegister(util.ServerResponses4xx)
	prometheus.MustRegister(util.ServerResponses5xx)
}

func MetricsMiddleware(ctx *gin.Context) {
	method := ctx.Request.Method
	uri := ctx.Request.URL.Path

	defer func(begin time.Time) {
		status := fmt.Sprintf("%d", ctx.Writer.Status())

		util.ServerRequests.WithLabelValues(method, uri).Dec()

		timeElapsed := float64(time.Since(begin)) / float64(time.Millisecond)
		util.ServerDuration.WithLabelValues(method, uri, status).Observe(timeElapsed)

		if ctx.Writer.Status() >= http.StatusBadRequest {
			util.ServerResponses4xx.WithLabelValues(method, uri, status).Inc()
		}

		if ctx.Writer.Status() >= http.StatusInternalServerError {
			util.ServerResponses5xx.WithLabelValues(method, uri, status).Inc()
		}

	}(time.Now())

	util.ServerRequests.WithLabelValues(method, uri).Inc()
	util.ServerRequestsTotal.WithLabelValues(method, uri).Inc()

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
