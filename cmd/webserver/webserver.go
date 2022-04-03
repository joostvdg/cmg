package webserver

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	promecho "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/joostvdg/cmg/cmd/context"
	"github.com/joostvdg/cmg/pkg/webserver"
)

const (
	envPort              = "PORT"
	envLogFormatter      = "LOG_FORMAT"
	envSentry            = "SENTRY_DSN"
	envLogLevel          = "LOG_LEVEL"
	envRootPath          = "ROOT_PATH"
	envAnalyticsEndpoint = "ANALYTICS_API_ENDPOINT"

	defaultRootPath       = "/"
	defaultPort           = "8080"
	debugLogLevel         = "DEBUG"
	defaultLogLevel       = "INFO"
	defaultLogFormatter   = "PLAIN"
	jsonLogFormatter      = "JSON"
	prometheusMetricsPath = "metrics"
	prometheusEchoSystem  = "echo"
	prometheusCmgSystem   = "cmg"
	envSegmentKey         = "SEGMENT_KEY"
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

	cmgAnalyticsEndpoint, cmgAnalyticsEndpointOk := os.LookupEnv(envAnalyticsEndpoint)
	if !cmgAnalyticsEndpointOk {
		cmgAnalyticsEndpoint = ""
	}

	// Echo instance
	e := echo.New()
	log.WithFields(log.Fields{
		"RootPath":           rootPath,
		"Port":               port,
		"LogFormatter":       logFormat,
		"LogLevel":           logLevel,
		"OS":                 runtime.GOOS,
		"ARCH":               runtime.GOARCH,
		"CPUs":               runtime.NumCPU(),
		"Sentry Enabled":     sentryOk,
		"Segment Enabled":    segmentOk,
		"Analytics Endpoint": cmgAnalyticsEndpoint,
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(sentryecho.New(sentryecho.Options{}))
	// Enable metrics middleware
	p := promecho.NewPrometheus(prometheusEchoSystem, nil)
	p.MetricsPath = rootPath + prometheusMetricsPath
	p.Use(e)

	// Define and Register custom Prometheus Metrics
	var mapGenCollector prometheus.Collector = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Subsystem: prometheusCmgSystem,
			Name:      "map_generations",
			Help:      "Number of map generation attempts it took to generate a valid map (1)",
			Buckets: []float64{
				1,
				10,
				25,
				50,
				100,
				250,
				500,
				1000,
				2000,
			},
		},
	)
	var mapGenDurationCollector prometheus.Collector = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Subsystem: prometheusCmgSystem,
			Name:      "map_generation_duration",
			Help:      "Duration of map generation attempts it took to generate a valid map",
		},
	)

	if err := prometheus.Register(mapGenCollector); err != nil {
		log.Warnf("Could not register Prometheus Collector for Map Generations: %v", err)
	}

	if err := prometheus.Register(mapGenDurationCollector); err != nil {
		log.Warnf("Could not register Prometheus Collector for Map Generation Duration: %v", err)
	}

	// Segment for Custom Context
	e.Use(func(e echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cmgContext := &context.CMGContext{
				Context:        c,
				MapGenAttempts: mapGenCollector,
				MapGenDuration: mapGenDurationCollector,
				SegmentClient:  segmentClient,
			}
			return e(cmgContext)
		}
	})

	// Routes
	g := e.Group(rootPath)
	g.GET("", handleRoutes)
	g.GET("routes", handleRoutes)

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
