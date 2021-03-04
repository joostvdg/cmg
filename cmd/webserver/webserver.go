package webserver

import (
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joostvdg/cmg/pkg/rollout"
	"github.com/joostvdg/cmg/pkg/webserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	roxServer "github.com/rollout/rox-go/server"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
)

const (
	envPort             = "PORT"
	envLogFormatter     = "LOG_FORMAT"
	envSentry           = "SENTRY_DSN"
	envDeploymentTarget = "DEPLOYMENT_TARGET"
	envRolloutKey       = "ROLLOUT_APP"
	envLogLevel         = "LOG_LEVEL"

	defaultPort             = "8080"
	debugLogLevel           = "DEBUG"
	defaultLogLevel         = "INFO"
	defaultDeploymentTarget = "LOCAL"
	defaultLogFormatter     = "PLAIN"
	jsonLogFormatter        = "JSON"
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

	deploymentTarget, deploymentOk := os.LookupEnv(envDeploymentTarget)
	if !deploymentOk {
		deploymentTarget = defaultDeploymentTarget
	}

	rolloutKey, rollOutOk := os.LookupEnv(envRolloutKey)

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
		"Rollout Enabled": rollOutOk,
		"Sentry Enabled":  sentryOk,
	}).Info("Webserver started")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Rollout
	if rollOutOk {
		options := roxServer.NewRoxOptions(roxServer.RoxOptionsBuilder{})

		rollout.Rox = roxServer.NewRox()
		rollout.Rox.Register("", rollout.RoxContainer)
		rollout.Rox.SetCustomStringProperty("DEPLOYMENT_TARGET", deploymentTarget)
		<-rollout.Rox.Setup(rolloutKey, options)
	}

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
