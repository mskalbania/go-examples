package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

/*
Counter with labels here - used to count the number of requests
Example metric exposed:
rest_app_request_count{method="GET",path="/metrics",status="200"} 1
*/
var requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "rest_app",
	Name:      "request_count",
	Help:      "Counts the number of requests served by http server",
}, []string{"method", "path", "status"})

/*
Histogram with labels here - used to measure the duration of request and put it into buckets
Example of histogram metrics exposed:
rest_app_request_duration_bucket{path="/metrics",le="0.005"} 1
rest_app_request_duration_bucket{path="/metrics",le="0.01"} 1
rest_app_request_duration_bucket{path="/metrics",le="0.025"} 1
rest_app_request_duration_sum{path="/metrics"} 0.001185791
rest_app_request_duration_count{path="/metrics"} 1
*/
var requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "rest_app",
	Name:      "request_duration",
	Help:      "Duration of requests served by http server",
	//custom buckets - will count the amount of request <= 0.2s, <= 0.5s, <= 1s, <= 5s
	//Buckets: []float64{0.2, 0.5, 1, 5},
}, []string{"path"})

func RegisterMetrics() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func Metrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start).Seconds()

		requestCounter.With(prometheus.Labels{
			"method": ctx.Request.Method,
			"path":   ctx.FullPath(),
			"status": strconv.Itoa(ctx.Writer.Status()),
		}).Inc()

		requestDuration.With(prometheus.Labels{
			"path": ctx.FullPath(),
		}).Observe(duration)
	}
}

func MetricsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	}
}
