package context

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/segmentio/analytics-go.v3"
)

type CMGContext struct {
	echo.Context
	analytics.Client
}