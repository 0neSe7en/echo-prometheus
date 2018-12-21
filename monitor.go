package echoprometheus

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type (
	PrometheusConfig struct {
		Skipper   middleware.Skipper
		Namespace string
	}
)

var (
	DefaultPrometheusConfig = PrometheusConfig{
		Skipper:   middleware.DefaultSkipper,
		Namespace: "echo",
	}
)

var (
	echoReqQps      *prometheus.CounterVec
	echoReqDuration *prometheus.SummaryVec
	echoOutBytes    prometheus.Summary
)

func initCollector(namespace string) {
	echoReqQps = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_total",
			Help:      "HTTP requests processed.",
		},
		[]string{"code", "method", "host", "url"},
	)
	echoReqDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		},
		[]string{"method", "host", "url"},
	)
	echoOutBytes = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response bytes.",
		},
	)
	prometheus.MustRegister(echoReqQps, echoReqDuration, echoOutBytes)
}

func NewMetric() echo.MiddlewareFunc {
	return NewMetricWithConfig(DefaultPrometheusConfig)
}

func NewMetricWithConfig(config PrometheusConfig) echo.MiddlewareFunc {
	initCollector(config.Namespace)
	if config.Skipper == nil {
		config.Skipper = DefaultPrometheusConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}
			uri := req.URL.Path
			status := strconv.Itoa(res.Status)
			elapsed := time.Since(start).Seconds()
			bytesOut := float64(res.Size)
			echoReqQps.WithLabelValues(status, req.Method, req.Host, uri).Inc()
			echoReqDuration.WithLabelValues(req.Method, req.Host, uri).Observe(elapsed)
			echoOutBytes.Observe(bytesOut)
			return nil
		}
	}
}
