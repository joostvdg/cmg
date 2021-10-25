package webserver

import (
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joostvdg/cmg/cmd/context"
	"github.com/joostvdg/cmg/pkg/rollout"
	"github.com/joostvdg/cmg/pkg/webserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/segmentio/analytics-go.v3"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	envPort         = "PORT"
	envLogFormatter = "LOG_FORMAT"
	envSentry       = "SENTRY_DSN"
	envSegmentKey   = "SEGMENT_KEY"
	envLogLevel     = "LOG_LEVEL"

	defaultPort         = "8080"
	debugLogLevel       = "DEBUG"
	defaultLogLevel     = "INFO"
	defaultLogFormatter = "PLAIN"
	jsonLogFormatter    = "JSON"
)

// StartWebserver starts the Echo webserver
// Retrieves environment variable PORT for the server port to listen on
// Retrieves environment variable SENTRY_DSN for exporting Sentry.io events
func StartWebserver() {
	port, portOk := os.LookupEnv(envPort)
	if !portOk {
		port = defaultPort
	}

	logFormat, logFormatOk := os.LookupEnv(envLogFormatter)
	if logFormatOk && logFormat == jsonLogFormatter {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		logFormat = defaultLogFormatter
	}

	logLevel, logLevelFormatOk := os.LookupEnv(envLogLevel)
	if logLevelFormatOk && logLevel == debugLogLevel {
		log.SetLevel(log.DebugLevel)
	} else {
		logLevel = defaultLogLevel
	}

	segmentKey, segmentOk := os.LookupEnv(envSegmentKey)

	sentryDsn, sentryOk := os.LookupEnv(envSentry)
	if sentryOk {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDsn,
		})

		if err != nil {
			log.Warnf("Sentry initialization failed: %v\n", err)
			sentryOk = false
		}
	}

	// Echo instance
	e := echo.New()
	log.WithFields(log.Fields{
		"Port":            port,
		"LogFormatter":    logFormat,
		"LogLevel":        logLevel,
		"OS":              runtime.GOOS,
		"ARCH":            runtime.GOARCH,
		"CPUs":            runtime.NumCPU(),
		"Sentry Enabled":  sentryOk,
		"Segment Enabled": segmentOk,
	}).Info("Webserver started")
	runtime.GOMAXPROCS(runtime.NumCPU())

	var segmentClient analytics.Client
	if segmentOk {
		segmentClient, _ = analytics.NewWithConfig(segmentKey, analytics.Config{
			Interval:  5 * time.Second,
			BatchSize: 10,
			Verbose:   true,
		})
		defer segmentClient.Close()
	}

	// Segment for Custom Context
	e.Use(func(e echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cmgContext := &context.CMGContext{
				c,
				segmentClient,
			}
			return e(cmgContext)
		}
	})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(sentryecho.New(sentryecho.Options{}))

	// Routes
	e.GET("/", hello)
	e.GET("/rollout", rolloutDemo)
	e.GET("/api/map", webserver.GetMap)
	e.GET("/api/v1/map", webserver.GetMapViaCodeGeneration)
	e.GET("/api/map/code", webserver.GetMapCode)
	e.GET("/api/map/code/:code", webserver.GetMapByCode)
	e.GET("/api/legend", webserver.GetMapLegend)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!!")
}

func rolloutDemo(c echo.Context) error {
	message := "Stay put"
	if rollout.RoxContainer.EnableTutorial.IsEnabled(nil) {
		message = "Lets Rollout"
	}
	return c.String(http.StatusOK, message)
}
