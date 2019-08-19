package webserver

import (
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joostvdg/cmg/pkg/webserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

// StartWebserver starts the Echo webserver
// Retrieves environment variable PORT for the server port to listen on
// Retrieves environment variable SENTRY_DSN for exporting Sentry.io events
func StartWebserver() {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}
	// Echo instance
	e := echo.New()
	e.Logger.Printf("Starting server on port %s\n", port)

	sentryDsn, ok := os.LookupEnv("SENTRY_DSN")
	if ok {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDsn,
		})

		if err != nil {
			e.Logger.Printf("Sentry initialization failed: %v\n", err)
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryecho.Options{}))

	// Routes
	e.GET("/", hello)
	e.GET("/api/map", webserver.GetMap)
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
