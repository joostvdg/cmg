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
	"strings"
	"time"
)

const (
	envPort         = "PORT"
	envLogFormatter = "LOG_FORMAT"
	envSentry       = "SENTRY_DSN"
	envSegmentKey   = "SEGMENT_KEY"
	envLogLevel     = "LOG_LEVEL"
	envRootPath     = "ROOT_PATH"

	defaultRootPath     = "/"
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

	rootPath, rootPathOk := os.LookupEnv(envRootPath)
	if !rootPathOk || rootPath == "" {
		rootPath = defaultRootPath
	}
	// ensure we end with a "/" so all derived paths will work
	if !strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
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
		"RootPath":        rootPath,
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
	g := e.Group(rootPath)
	g.GET("", handleRoutes)
	g.GET("routes", handleRoutes)
	g.GET("rollout", rolloutDemo)

	g.GET("api/map", webserver.GetMap)
	g.GET("api/v1/map", webserver.GetMapViaCodeGeneration)
	g.GET("api/map/code", webserver.GetMapCode)
	g.GET("api/map/code/:code", webserver.GetMapByCode)
	g.GET("api/legend", webserver.GetMapLegend)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

// handleRoutes shows the routes that are handled
func handleRoutes(c echo.Context) error {
	content := c.Echo().Routes()

	callback := c.QueryParam("callback")
	jsonp := c.QueryParam("jsonp")

	if jsonp == "true" {
		return c.JSONP(http.StatusOK, callback, &content)
	}
	return c.JSON(http.StatusOK, &content)
}

func rolloutDemo(c echo.Context) error {
	message := "Stay put"
	if rollout.RoxContainer.EnableTutorial.IsEnabled(nil) {
		message = "Lets Rollout"
	}
	return c.String(http.StatusOK, message)
}
