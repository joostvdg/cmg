package context

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/segmentio/analytics-go.v3"
)

type CMGContext struct {
	echo.Context
	MapGenAttempts prometheus.Collector
	MapGenDuration prometheus.Collector
	SegmentClient  analytics.Client
}
