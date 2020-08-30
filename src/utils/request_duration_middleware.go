package utils

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber"
	"github.com/prometheus/client_golang/prometheus"
)

func RequestDurationMonitor() func(c *fiber.Ctx) {

	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "lumiere",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
		Buckets:   buckets,
	}, []string{"method", "status_code"})

	prometheus.MustRegister(responseTimeHistogram)

	return func(c *fiber.Ctx) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		statusCode := strconv.Itoa(c.Fasthttp.Response.StatusCode())
		responseTimeHistogram.WithLabelValues(c.Path(), statusCode).Observe(duration.Seconds())
	}
}
