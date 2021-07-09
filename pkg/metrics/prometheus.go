package metrics

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		}, []string{"code", "method", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
		}, []string{"code", "method", "path"},
	)
)

func Middleware(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		statusCode := strconv.Itoa(ctx.Writer.Status())
		method := ctx.Request.Method

		requestDuration.WithLabelValues(statusCode, method, path).Observe(f)
		totalRequests.WithLabelValues(statusCode, method, path).Inc()
	}))

	defer timer.ObserveDuration()

	ctx.Next()
}

func init() {
	if err := prometheus.Register(totalRequests); err != nil {
		log.Fatalf("prometheus error: %v", err)
	}
	if err := prometheus.Register(requestDuration); err != nil {
		log.Fatalf("prometheus error: %v", err)
	}
}
